package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var OrderStatusUpdatedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "order_status_updated_total",
	Help: "Total number of order status updates, labeled by new status",
}, []string{"status"})

var OrdersAddedCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "orders_added_total",
	Help: "Total number of orders added",
})

func init() {
	prometheus.MustRegister(OrderStatusUpdatedCounter)
	prometheus.MustRegister(OrdersAddedCounter)
}
