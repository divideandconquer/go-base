package test

import (
	"net/http"

	"github.com/divideandconquer/negotiator"
	"github.com/go-martini/martini"
)

type pong struct {
	Message string `json:"message"`
}

func Ping() martini.Handler {
	return func(w http.ResponseWriter, r *http.Request, neg negotiator.Negotiator) (int, []byte) {
		ret := pong{Message: "pong!"}
		return http.StatusOK, negotiator.Must(neg.Negotiate(r, ret))
	}
}
