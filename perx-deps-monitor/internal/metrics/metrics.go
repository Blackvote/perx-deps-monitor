package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DepsMetrics struct {
	MongoUp       prometheus.Gauge
	NatsUp        prometheus.Gauge
	LastCheckUnix prometheus.Gauge
	CheckCount    prometheus.Counter
	CheckDuration prometheus.Histogram
	MongoErrors   prometheus.Counter
	NatsErrors    prometheus.Counter
}

func NewDepsMetrics(reg prometheus.Registerer) *DepsMetrics {
	m := &DepsMetrics{
		MongoUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dep_mongodb_up",
			Help: "MongoDB availability (1=up, 0=down).",
		}),
		NatsUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dep_nats_up",
			Help: "NATS availability (1=up, 0=down).",
		}),
		LastCheckUnix: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "deps_last_check_timestamp_seconds",
			Help: "Unix timestamp of the last dependency check.",
		}),
		CheckCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "deps_check_total",
			Help: "Total number of dependency checks performed.",
		}),
		CheckDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "deps_check_duration_seconds",
			Help:    "Duration of a dependency check.",
			Buckets: prometheus.DefBuckets,
		}),
		MongoErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "dep_mongodb_errors_total",
			Help: "Total number of MongoDB check/connect errors.",
		}),
		NatsErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "dep_nats_errors_total",
			Help: "Total number of NATS check/connect errors.",
		}),
	}

	reg.MustRegister(
		m.MongoUp,
		m.NatsUp,
		m.LastCheckUnix,
		m.CheckCount,
		m.CheckDuration,
		m.MongoErrors,
		m.NatsErrors,
	)

	return m
}
