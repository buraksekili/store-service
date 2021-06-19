package rest

import "github.com/prometheus/client_golang/prometheus"

const (
	PID   = "productID"
	PNAME = "productName"
)

var prodCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:      "prod_count",
		Namespace: "store",
		Help:      "Count of created products",
	}, []string{PID, PNAME})

func init() {
	prometheus.MustRegister(prodCounter)
}
