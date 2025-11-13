package botprom

import (
	"net/http"

	"github.com/a-random-lemurian/go-owot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "owot_log"

var (
	messagesReceived = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_received",
		Help:      "The total number of messages received",
	})
)

type BotPrometheus struct {
	metrics    *prometheus.Registry
	httpServer *http.Server
	httpMux    *http.ServeMux
	listen     string
}

func NewBotPrometheus(listen string) (*BotPrometheus, error) {
	bp := &BotPrometheus{
		metrics: prometheus.NewRegistry(),
		listen: listen,
	}

	bp.metrics.MustRegister(
		messagesReceived,

		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	bp.httpMux = http.NewServeMux()
	bp.httpMux.Handle(
		"/metrics",
		promhttp.HandlerFor(
			prometheus.Gatherers{
				bp.metrics,
			},
			promhttp.HandlerOpts{},
		),
	)
	bp.httpServer = &http.Server{
		Handler: bp.httpMux,
	}

	return bp, nil
}

func (bp *BotPrometheus) Run() error {
	bp.httpServer.Addr = bp.listen
	return bp.httpServer.ListenAndServe()
}

func (bp *BotPrometheus) IngestMessage(c *owot.MessageChat) {
	messagesReceived.Inc()
}
