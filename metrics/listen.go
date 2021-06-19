package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServe(addr string) {
	log.Println("Metrics are listening on ", addr)

	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Printf("Prometheus metrics error, %v", err)
	}
}
