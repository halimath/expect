package expect

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	var tm contextMock

	var r httptest.ResponseRecorder

	HTTPStatus(http.StatusOK).Match(&tm, &r)
	HTTPStatus(http.StatusNotFound).Match(&tm, &r)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected HTTP status code to be 404 but got 200"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestHTTPHeader(t *testing.T) {
	var tm contextMock

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("X-Foo", "bar")

	HTTPHeader[*http.Request]("X-Spam", "eggs").Match(&tm, r)
	HTTPHeader[*http.Request]("X-Foo", "foo").Match(&tm, r)
	HTTPHeader[*http.Request]("X-Foo", "bar").Match(&tm, r)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected *http.Request to contain header X-Spam",
			"expected *http.Request to contain header X-Foo with value foo but no such value found (although header with that name exists)",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
