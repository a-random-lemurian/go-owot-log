package bot

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "owot_log"

var (
	worldConnections = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "world_connections_created",
		Help:      "The total number of connection attempts done",
	}, []string{"status"})
	worldDisconnects = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "world_disconnections",
		Help:      "The total number of world disconnections",
	})
	worldGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "connected_worlds",
		Help: "The total number of worlds being actively logged by the bot",
	})
)
