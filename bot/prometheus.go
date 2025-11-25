package bot

import (
	"github.com/a-random-lemurian/go-owot-log/internal/util"
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
		Name:      "connected_worlds",
		Help:      "The total number of worlds being actively logged by the bot",
	})
)

func (b *Bot) initChatsizeMetric() prometheus.Collector {
	chatSize := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "data_directory_size_bytes",
		Help: "The total size of the data directory",
	}, b.chatsizeMetric)
	return chatSize
}

func (b *Bot) chatsizeMetric() float64 {
	size, err := util.DirectorySize(b.v.GetString("data_directory"))
	if err != nil {
		return 0
	}
	return float64(size)
}

func getPrometheusRegistry() prometheus.Collector {
	reg := prometheus.NewRegistry()
	reg.MustRegister(worldConnections, worldDisconnects, worldGauge)
	return reg
}
