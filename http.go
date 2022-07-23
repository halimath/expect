package expect

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// WithHeader defines a type constraint to use the HTTPHeader
type WithHeader interface {
	*http.Request | httptest.ResponseRecorder
}

// HTTPHeader returns a Matcher that asserts that the given HTTP entity contains a header with a given value.
func HTTPHeader[T WithHeader](name, value string) Matcher[T] {
	return MatcherFunc[T](func(ctx Context, got T) {
		var g any = got
		var header http.Header

		switch x := g.(type) {
		case *http.Request:
			header = x.Header
		case http.ResponseWriter:
			header = x.Header()
		default:
			panic(fmt.Sprintf("unexpected HTTP entity: %v", got))
		}

		values, ok := header[name]
		if !ok {
			ctx.Failf("expected %T to contain header %s", got, name)
			return
		}

		for _, v := range values {
			if v == value {
				return
			}
		}

		ctx.Failf("expected %T to contain header %s with value %s but no such value found (although header with that name exists)", got, name, value)
	})
}

// HTTPStatus returns a Matcher that asserts that the given HTTP ResponseRecorder has the given status code.
func HTTPStatus(status int) Matcher[*httptest.ResponseRecorder] {
	return MatcherFunc[*httptest.ResponseRecorder](func(ctx Context, got *httptest.ResponseRecorder) {
		if got.Result().StatusCode != status {
			ctx.Failf("expected HTTP status code to be %d but got %d", status, got.Result().StatusCode)
		}
	})
}
