package test

import (
	"net/http"
	"testing"
)

func TestUnit_Ping_BasePath(t *testing.T) {
	handler := Ping()
	ti := NewTestInjector("GET", "")
	InvokeAndCheck(t, handler, ti, http.StatusOK, []byte(`{"message":"pong!"}`))
}
