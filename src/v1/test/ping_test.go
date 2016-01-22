package test

import (
	"net/http"
	"testing"
)

func TestUnit_Ping_BasePath(t *testing.T) {
	handler := Ping()
	i := NewInjector("GET", "")
	InvokeAndCheck(t, handler, i, http.StatusOK, []byte(`{"message":"pong!"}`))
}

func TestFunctional_Ping_BasePath(t *testing.T) {
	res, err := DoTestRequest(t, "GET", "test/ping", "", nil)
	VerifyResponseBody(t, res, err, http.StatusOK, []byte(`{"message":"pong!"}`))
}
