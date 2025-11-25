package botprom

import (
	"net/http"

	"github.com/a-random-lemurian/go-owot"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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
	log        *zap.Logger
}

func NewBotPrometheus(listen string, log *zap.Logger) (*BotPrometheus, error) {
	bp := &BotPrometheus{
		metrics: prometheus.NewRegistry(),
		listen:  listen,
		log:     log.Named("botprom"),
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
	bp.log.Info("Running Prometheus HTTP server")
	bp.httpServer.Addr = bp.listen
	return bp.httpServer.ListenAndServe()
}

func (bp *BotPrometheus) IngestMessage(c *owot.MessageChat) {
	messagesReceived.Inc()
}

func (bp *BotPrometheus) MustRegister(cs ...prometheus.Collector) {
	bp.metrics.MustRegister(cs...)
}
