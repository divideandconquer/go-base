package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/divideandconquer/go-base/src/v1/test"
	"github.com/divideandconquer/go-consul-client/src/client"
	"github.com/divideandconquer/negotiator"
	"github.com/go-martini/martini"
)

func main() {
	consulAddress := mustGetEnvVar("CONSUL_HTTP_ADDR")
	environment := mustGetEnvVar("ENVIRONMENT")
	serviceName := mustGetEnvVar("SERVICE_NAME")

	conf := setupConfig(serviceName, environment, consulAddress)

	//setup martini
	m := martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	router := martini.NewRouter()
	router.Group("/v1", func(v1router martini.Router) {
		//Setup v1 routes
		v1router.Group("/test", func(r martini.Router) {
			r.Get("/ping", test.Ping())
		})
	})

	// Add the router action
	m.Action(router.Handle)

	// Inject dependencies
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		enc := negotiator.JsonEncoder{false}
		cn := negotiator.NewContentNegotiator(enc, w)
		cn.AddEncoder(negotiator.MimeJSON, enc)
		c.MapTo(cn, (*negotiator.Negotiator)(nil))
	})

	// Start up the server
	m.RunOnAddr(":8080")
}

func setupConfig(service string, env string, consulAddr string) client.Loader {
	appNamespace := fmt.Sprintf("%s/%s", env, service)

	// create a cached loader
	conf, err := client.NewCachedLoader(appNamespace, consulAddr)
	if err != nil {
		log.Fatalf("Could not create a config cached loader: %v", err)
	}

	// initialize the cache
	err = conf.Initialize()
	if err != nil {
		log.Fatalf("Could not initialize the config cached loader: %v", err)
	}
}

func mustGetEnvVar(key string) string {
	ret := os.Getenv(key)
	if ret == "" {
		log.Fatalf("Could not find %s environment variable", key)
	}
	return ret
}
