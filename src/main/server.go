package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/divideandconquer/go-base/src/v1/test"
	"github.com/divideandconquer/go-consul-client/src/balancer"
	"github.com/divideandconquer/go-consul-client/src/client"
	"github.com/divideandconquer/negotiator"
	"github.com/go-martini/martini"
)

// config keys
const (
	configKeyBalancerTTL = "balancerTTL"
)

func main() {
	consulAddress := mustGetEnvVar("CONSUL_HTTP_ADDR")
	environment := mustGetEnvVar("ENVIRONMENT")
	serviceName := mustGetEnvVar("SERVICE_NAME")

	conf, _ := setupConfig(serviceName, environment, consulAddress)
	//pull config and pass it around
	conf.MustGetDuration(configKeyBalancerTTL)

	// use loadbalancer to lookup services and db locations if necessary
	// Each request should repeat the lookup to make sure that this app
	// follows any services that move
	// s, err := loadbalancer.FindService("foo")

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
		enc := negotiator.JsonEncoder{PrettyPrint: false}
		cn := negotiator.NewContentNegotiator(enc, w)
		cn.AddEncoder(negotiator.MimeJSON, enc)
		c.MapTo(cn, (*negotiator.Negotiator)(nil))
	})

	// Start up the server
	m.RunOnAddr(":8080")
}

func setupConfig(service string, env string, consulAddr string) (client.Loader, balancer.DNS) {
	appNamespace := fmt.Sprintf("%s/%s", env, service)
	log.Printf("Initializing config for %s with consul %s", appNamespace, consulAddr)

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

	loadbalancer, err := balancer.NewRandomDNSBalancer(env, consulAddr, conf.MustGetDuration(configKeyBalancerTTL))
	if err != nil {
		log.Fatalf("Could not create loadbalancer: %v", err)
	}

	return conf, loadbalancer
}

func mustGetEnvVar(key string) string {
	ret := os.Getenv(key)
	if ret == "" {
		log.Fatalf("Could not find %s environment variable", key)
	}
	return ret
}
