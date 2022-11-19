package expect

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// DeepEqualOpt defines an interface for types that can be used as options
// for the DeepEqual matcher.
type DeepEqualOpt interface {
	deepEqualOpt()
}

// FloatPrecision is a DeepEqualOpt that customizes the float comparison
// behavior. The number given defines the number of significant floating point
// digits.
type FloatPrecision uint

func (FloatPrecision) deepEqualOpt() {}

// NilSlicesAreEmpty is a DeepEqualOpt that defines whether nil slices are
// treated as empty ones or differently.
type NilSlicesAreEmpty bool

func (NilSlicesAreEmpty) deepEqualOpt() {}

// NilMapsAreEmpty is a DeepEqualOpt that defines whether nil maps are treated
// as empty ones or differently.
type NilMapsAreEmpty bool

func (NilMapsAreEmpty) deepEqualOpt() {}

// ExcludeUnexportedStructFields is a DeepEqualOpt that defines whether
// unexported struct fields should be excluded from the equality check or
// not.
type ExcludeUnexportedStructFields bool

func (ExcludeUnexportedStructFields) deepEqualOpt() {}

// DeepEqual asserts that given and wanted value are deeply equal by using reflection to inspect and dive into
// nested structures.
func DeepEqual[T any](want T, opts ...DeepEqualOpt) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		// if diff := deep.Equal(want, got); diff != nil {
		// 	ctx.Failf("values are not deeply equal\nwant: %v\ngot: %v\ndiff: %s", want, got, diff)
		// }

		if diff := deepEquals(want, got, opts...); diff != nil {
			ctx.Failf("values are not deeply equal:%s", diff)
		}
	})
}

func deepEquals(want, got any, opts ...DeepEqualOpt) diff {
	ctx := &diffContext{
		floatFormat:       fmt.Sprintf("%%.%df", 10),
		nilSlicesAreEmpty: true,
		nilMapsAreEmpty:   true,
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case FloatPrecision:
			ctx.floatFormat = fmt.Sprintf("%%.%df", o)
		case NilSlicesAreEmpty:
			ctx.nilSlicesAreEmpty = bool(o)
		case NilMapsAreEmpty:
			ctx.nilMapsAreEmpty = bool(o)
		case ExcludeUnexportedStructFields:
			ctx.excludeUnexportedStructFields = bool(o)
		}
	}

	wv := reflect.ValueOf(want)
	gv := reflect.ValueOf(got)

	determineDiff(ctx, wv, gv)

	return ctx.diff
}

func determineDiff(ctx *diffContext, want, got reflect.Value) {
	// If want has already been visited, determination ends here to not run into cycles.
	if ctx.hasVisisted(want) {
		return
	}
	// Otherwise mark want as visited.
	ctx.visit(want)

	// If neither want nor got are value (i.e. both are nil) there is no difference.
	if !want.IsValid() && !got.IsValid() {
		return
	}

	// If either want xor got is invalid we have a difference (nil vs. !nil).
	if !want.IsValid() && got.IsValid() {
		ctx.addDiff("<nil>", got)
		return
	}
	if want.IsValid() && !got.IsValid() {
		ctx.addDiff(want, "<nil>")
		return
	}

	// Compare both underlying types. If they are not the same the values cannot be equal so we report a
	// type difference.
	wantType := want.Type()
	gotType := got.Type()
	if wantType != gotType {
		ctx.addDiff(wantType, gotType)
		return
	}

	// Inspect the value's kinds.
	wantKind := want.Kind()
	gotKind := got.Kind()

	// Determine if either want or got refer to some other value (for being a
	// pointer of an interface value).
	wantHasElem := wantKind == reflect.Ptr || wantKind == reflect.Interface
	gotHasElem := gotKind == reflect.Ptr || gotKind == reflect.Interface
	// If so, follow the reference and compare the underlying values.
	if wantHasElem || gotHasElem {
		if wantHasElem {
			want = want.Elem()
		}
		if gotHasElem {
			got = got.Elem()
		}
		determineDiff(ctx, want, got)
		return
	}

	// Otherwise compare the values based on their kind.
	switch wantKind {
	case reflect.Struct:
		// Iterate over struct fields and compare the values. The structs
		// must be of the same type, otherwise the values wouldn't have passed
		// the type check.
		for i := 0; i < want.NumField(); i++ {
			f := wantType.Field(i)
			if f.PkgPath != "" && ctx.excludeUnexportedStructFields {
				continue
			}

			ctx.pushPathf(".%s", f.Name)
			wantVal := want.Field(i)
			gotVal := got.Field(i)

			determineDiff(ctx, wantVal, gotVal)

			ctx.popPath()
		}

	case reflect.Map:
		// If nil maps aren't considered empty and exactly one map is nil
		// we have a difference and return here.
		if !ctx.nilMapsAreEmpty && want.IsNil() != got.IsNil() {
			if want.IsNil() {
				ctx.addDiff("<nil map>", got)
			} else {
				ctx.addDiff(want, "<nil map>")
			}
			return
		}

		// Determine lengths
		var wantLen, gotLen int
		if want.IsNil() {
			wantLen = 0
		} else {
			wantLen = want.Len()
		}
		if got.IsNil() {
			gotLen = 0
		} else {
			gotLen = got.Len()
		}

		// Two empty maps are always equal.
		if wantLen == 0 && gotLen == 0 {
			return
		}

		// If the maps are of the same length and point to the same start element
		// they are equal.
		if wantLen == gotLen && want.Pointer() == got.Pointer() {
			return
		}

		// Iterate over wanted keys
		for _, wantKey := range want.MapKeys() {
			ctx.pushPathf("[%v]", wantKey)
			wantVal := want.MapIndex(wantKey)
			gotVal := got.MapIndex(wantKey)
			if !gotVal.IsValid() {
				ctx.addDiff(wantVal, "<missing map key>")
			} else {
				determineDiff(ctx, wantVal, gotVal)
			}
			ctx.popPath()
		}

		// Do the same with got keys
		for _, gotKey := range got.MapKeys() {
			ctx.pushPathf("[%v]", gotKey)
			wantVal := want.MapIndex(gotKey)
			gotVal := got.MapIndex(gotKey)
			// No need to handle a valid wantVal here as it has been handled
			// in the previous loop over wanted keys. Thus, we only check for
			// invalid wantVal here and report a missing key in want.
			if !wantVal.IsValid() {
				ctx.addDiff("<missing map key>", gotVal)
			}
			ctx.popPath()
		}

	case reflect.Slice:
		// Check, if nil slices should be treated as empty ones. If not and
		// exactly one slice is nil, there is a difference.
		if !ctx.nilSlicesAreEmpty && want.IsNil() != got.IsNil() {
			if want.IsNil() {
				ctx.addDiff("<nil slice>", got)
			} else {
				ctx.addDiff(want, "<nil slice>")
			}

			return
		}

		// Determine slice lengths. Treat nil slices as zero length here. This
		// is ok because if nilSlicesAreEmpty is set to false we only reach
		// this point if either both slices are nil or both are non-nil. Thus
		// we can safely assume a length of zero for nil slices.
		var wantLen, gotLen int
		if want.IsNil() {
			wantLen = 0
		} else {
			wantLen = want.Len()
		}
		if got.IsNil() {
			gotLen = 0
		} else {
			gotLen = got.Len()
		}

		// If lengths are equal and both slice's first element point to the
		// same address, the slices must be equal.
		if wantLen == gotLen && want.Pointer() == got.Pointer() {
			return
		}

		// Iterate over elements and compare them
		for i := 0; i < wantLen; i++ {
			ctx.pushPathf("[%d]", i)
			wantVal := want.Index(i)

			if i >= gotLen {
				ctx.addDiff(wantVal, "<missing slice index>")
			} else {
				gotVal := got.Index(i)
				determineDiff(ctx, wantVal, gotVal)
			}
			ctx.popPath()
		}

	case reflect.Array:
		// No need to compare lengths here; for arrays the length is part of
		// the type declaration, so [1]any is a different type then [2]any.
		// Thus, type difference has been checked before and we only get to
		// this point if both arrays share the same underlying type and length.

		l := want.Len()

		// Iterate over elements and compare them
		for i := 0; i < l; i++ {
			ctx.pushPath(fmt.Sprintf("[%d]", i))
			determineDiff(ctx, want.Index(i), got.Index(i))
			ctx.popPath()
		}

	case reflect.Float32, reflect.Float64:
		w := fmt.Sprintf(ctx.floatFormat, want.Float())
		g := fmt.Sprintf(ctx.floatFormat, got.Float())
		addDiffIfUnequal(ctx, w, g)

	case reflect.Bool:
		addDiffIfUnequal(ctx, want.Bool(), got.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		addDiffIfUnequal(ctx, want.Int(), got.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		addDiffIfUnequal(ctx, want.Uint(), got.Uint())

	case reflect.String:
		addDiffIfUnequal(ctx, want.String(), got.String())

	default:
		panic(fmt.Sprintf("unimplemented diff kind: %v", wantKind))
	}
}

type diffEntry struct {
	path      string
	want, got string
}

func (d diffEntry) writeTo(w io.Writer) {
	if len(d.path) == 0 {
		fmt.Fprintf(w, "  want: %s\n   got: %s", d.want, d.got)
		return
	}
	fmt.Fprintf(w, "  at %s\n    want: %s\n     got: %s", d.path, d.want, d.got)
}

type diff []diffEntry

func (d diff) String() string {
	var b strings.Builder

	for i := range d {
		b.WriteRune('\n')
		d[i].writeTo(&b)
	}

	return b.String()
}

type diffContext struct {
	floatFormat                   string
	nilSlicesAreEmpty             bool
	nilMapsAreEmpty               bool
	excludeUnexportedStructFields bool

	wantsSeen   set[reflect.Value]
	diff        diff
	nestingPath []string
}

func (c *diffContext) visit(want reflect.Value) {
	if c.wantsSeen == nil {
		c.wantsSeen = make(set[reflect.Value])
	}
	c.wantsSeen.Add(want)
}

func (c *diffContext) hasVisisted(want reflect.Value) bool {
	return c.wantsSeen.Contains(want)
}

func addDiffIfUnequal[T comparable](ctx *diffContext, want, got T) {
	if want == got {
		return
	}

	ctx.addDiff(want, got)
}

func (c *diffContext) addDiff(want, got any) {
	var w, g string

	if s, ok := want.(string); ok {
		w = s
	} else {
		w = fmt.Sprint(want)
	}

	if s, ok := got.(string); ok {
		g = s
	} else {
		g = fmt.Sprint(got)
	}

	c.diff = append(c.diff, diffEntry{
		path: c.path(),
		want: w,
		got:  g,
	})
}
func (c *diffContext) pushPathf(format string, args ...any) {
	c.pushPath(fmt.Sprintf(format, args...))
}

func (c *diffContext) pushPath(p string) {
	c.nestingPath = append(c.nestingPath, p)
}

func (c *diffContext) popPath() {
	c.nestingPath = c.nestingPath[:len(c.nestingPath)-1]
}

func (c *diffContext) path() string {
	if len(c.nestingPath) == 0 {
		return ""
	}

	return strings.Join(c.nestingPath, "")
}
