package main

import (
	"net/http"

	"github.com/divideandconquer/go-base/src/test"
	"github.com/divideandconquer/negotiator"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()

	r.Group("/v1", func(v1router martini.Router) {

		//Setup v1 routes

		v1router.Group("/test", func(v1router martini.Router) {
			r.Get("/ping", test.Ping())
		})
	})

	// Add the router action
	m.Action(r.Handle)

	// Inject dependencies
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		enc := negotiator.JsonEncoder{true}
		cn := negotiator.NewContentNegotiator(enc, w)
		cn.AddEncoder(negotiator.MimeJSON, enc)
		c.MapTo(cn, (*negotiator.Negotiator)(nil))
	})

	// Start up the server
	m.RunOnAddr(":8080")
}
