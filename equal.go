package expect

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

// DeepEqual asserts that given and wanted value are deeply equal by using reflection to inspect and dive into
// nested structures.
func DeepEqual[T any](want T) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		if !reflect.DeepEqual(want, got) {
			ws := fmt.Sprintf("%v", want)
			gs := fmt.Sprintf("%v", got)
			var ds strings.Builder

			var wstart, gstart int

			for {
				if wstart >= len(ws) || gstart >= len(gs) {
					break
				}

				w, wl := utf8.DecodeRuneInString(ws[wstart:])
				g, gl := utf8.DecodeRuneInString(gs[gstart:])

				if w == g {
					ds.WriteRune(' ')
				} else {
					ds.WriteRune('▲')
				}

				wstart += wl
				gstart += gl
			}

			ctx.Failf("\nvalues are not deeply equal\n\nwant: %s\ngot:  %s\n      %s\n", ws, gs, ds.String())
		}
	})
}

// Equal asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func Equal[G comparable](want G) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		if want != got {
			ctx.Failf("values are not equal\nwant: %v\ngot:  %v", want, got)
		}
	})
}
