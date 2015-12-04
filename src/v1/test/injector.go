package test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/codegangsta/inject"
	"github.com/divideandconquer/negotiator"
	"github.com/go-martini/martini"
)

// TestInjector is a testing helper that can invoke martini handlers
type TestInjector struct {
	inject.Injector
	r *http.Request
}

// NewTestInjector creates a TestInjector
func NewTestInjector(method string, body string) *TestInjector {
	var t TestInjector
	t.r, _ = http.NewRequest(method, "http://localhost/v1/", strings.NewReader(body))

	w := httptest.NewRecorder()

	enc := negotiator.JsonEncoder{false}
	cn := negotiator.NewContentNegotiator(enc, w)
	cn.AddEncoder(negotiator.MimeJSON, enc)

	t.Injector = inject.New()
	t.Injector.Map(t.r)
	t.Injector.MapTo(w, (*http.ResponseWriter)(nil))
	t.Injector.MapTo(cn, (*negotiator.Negotiator)(nil))

	return &t
}

// SetParams sets up DI for martinit params
func (t *TestInjector) SetParams(p martini.Params) {
	t.Injector.Map(p)
}

// SetHeaders sets the headers on the internal request
func (t *TestInjector) SetHeaders(h http.Header) {
	t.r.Header = h
}

func injectAndInvoke(t *testing.T, handler martini.Handler, ti *TestInjector) (int, []byte) {
	val, err := ti.Invoke(handler)
	if err != nil {
		t.Fatalf("Unexpected error invoking handler: %v", err)
	}
	if len(val) != 2 {
		t.Fatalf("Unexpected number of return values: %d", len(val))
	}

	intVal := val[0]
	if intVal.Kind() != reflect.Int {
		t.Fatalf("Unexpected type for int returned: %s", intVal.Kind().String())
	}
	sliceVal := val[1]
	if sliceVal.Kind() != reflect.Slice {
		t.Fatalf("Unexpected body type returned: %s", sliceVal.Kind().String())
	}

	return int(intVal.Int()), sliceVal.Bytes()
}

// InvokeAndCheck will invoke a handler with the provided injector and check its return values
func InvokeAndCheck(t *testing.T, handler martini.Handler, ti *TestInjector, expectedStatus int, expectedBody []byte) {
	status, body := injectAndInvoke(t, handler, ti)
	if status != expectedStatus {
		t.Fatalf("Unexpected status returned %d instead of %d", status, expectedStatus)
	}

	if !reflect.DeepEqual(body, expectedBody) {
		t.Fatalf("Unexpected body returned %s instead of %s", body, expectedBody)
	}
}
