package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dnsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dc",
			Subsystem: "dns",
			Name:      "requests_total",
			Help:      "Total dns request.",
			ConstLabels: map[string]string{
				"app": applicationName,
			},
		},
		[]string{"found", "nameServer", "targetHost"},
	)
)

func init() {
	prometheus.MustRegister(dnsCounter)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func exposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
