package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mcbobke/ynab-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	modVersion    string = "0.0.1"
	ynabCollector collector.YnabCollector
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)

	token, success := os.LookupEnv("YNAB_API_TOKEN")
	if success {
		log.Println("YNAB token found")
	} else {
		log.Fatalln("YNAB token not found")
	}

	ynabCollector = collector.New(token)
	prometheus.MustRegister(ynabCollector)
	http.Handle("/metrics", promhttp.Handler())
}

func main() {
	log.Printf("Starting ynab-exporter version %s", modVersion)
	// TODO: Configure signal handling
	http.ListenAndServe("localhost:9090", nil)
	// TODO: Handle signals
}
