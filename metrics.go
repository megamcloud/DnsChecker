package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	DNSCounter *prometheus.CounterVec
}

func (metrics *metrics) incDNSCounter(found bool, nameServer, targetHost string) {
	metrics.DNSCounter.WithLabelValues(ifThenElse(found, "true", "false").(string), nameServer, targetHost).Inc()
}

func ifThenElse(condition bool, ifValue interface{}, elseValue interface{}) interface{} {
	if condition {
		return ifValue
	}

	return elseValue
}

func newMetrics(config *config) *metrics {
	metrics := metrics{
		DNSCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "dc",
				Subsystem: "dns",
				Name:      "requests_total",
				Help:      "Total dns request.",
				ConstLabels: map[string]string{
					"app": config.ApplicationName,
				},
			},
			[]string{"found", "nameServer", "targetHost"},
		),
	}

	return &metrics
}

func (app *app) serveMetrics() {
	prometheus.MustRegister(app.Metrics.DNSCounter)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(app.Config.ListenAddr, nil))
}
