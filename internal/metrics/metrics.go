package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PVZCreated        prometheus.Counter
	ReceptionsCreated prometheus.Counter
	ProductsAdded     prometheus.Counter
)

func init() {
	PVZCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "total number of pvz created",
		},
	)

	ReceptionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "total number of receptions created",
		},
	)

	ProductsAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "total number of products added",
		},
	)
	prometheus.MustRegister(PVZCreated)
	prometheus.MustRegister(ReceptionsCreated)
	prometheus.MustRegister(ProductsAdded)
}
