package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/google/gops/agent"
	"github.com/yaput/gosample/page"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tokopedia/logging/tracer"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	// Page Index
	http.HandleFunc("/", page.Index)
	// API User List
	http.HandleFunc("/first", page.GetUsers)
	go logging.StatsLog()

	tracer.Init(&tracer.Config{Port: 8700, Enabled: true})

	log.Fatal(grace.Serve(":9000", nil))
}
