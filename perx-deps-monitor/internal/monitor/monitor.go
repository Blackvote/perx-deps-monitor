package monitor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"perx-deps-monitor/internal/metrics"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatusSnapshot struct {
	Mongo      string    `json:"mongodb"`
	NATS       string    `json:"nats"`
	CheckedAt  time.Time `json:"checked_at"`
	StartedAt  time.Time `json:"started_at"`
	CheckCount int64     `json:"check_count"`
}

type Config struct {
	MongoURI string
	NATSURL  string
	Interval time.Duration
	Timeout  time.Duration
	Logger   *slog.Logger
}

type Monitor struct {
	mongoURI string
	natsURL  string

	interval time.Duration
	timeout  time.Duration

	mu     sync.RWMutex
	status StatusSnapshot

	mongo *mongo.Client
	nc    *nats.Conn

	m   *metrics.DepsMetrics
	log *slog.Logger
}

func New(cfg Config, reg prometheus.Registerer) *Monitor {
	now := time.Now()

	l := cfg.Logger
	if l == nil {
		l = slog.Default()
	}

	mon := &Monitor{
		mongoURI: cfg.MongoURI,
		natsURL:  cfg.NATSURL,
		interval: cfg.Interval,
		timeout:  cfg.Timeout,
		status: StatusSnapshot{
			Mongo:     "unknown",
			NATS:      "unknown",
			StartedAt: now,
		},
		m:   metrics.NewDepsMetrics(reg),
		log: l,
	}

	if mon.interval <= 0 {
		mon.interval = 2 * time.Second
	}
	if mon.timeout <= 0 {
		mon.timeout = 2 * time.Second
	}

	return mon
}

func (m *Monitor) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.checkOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkOnce(ctx)
		}
	}
}

func (m *Monitor) Snapshot() StatusSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

func (m *Monitor) HasCheckedAtLeastOnce() bool {
	return !m.Snapshot().CheckedAt.IsZero()
}

func (m *Monitor) Close(ctx context.Context) {
	if m.mongo != nil {
		_ = m.mongo.Disconnect(ctx)
	}
	if m.nc != nil {
		m.nc.Close()
	}
}

func (m *Monitor) checkOnce(parent context.Context) {
	start := time.Now()
	defer func() {
		m.m.CheckDuration.Observe(time.Since(start).Seconds())
	}()

	ctx, cancel := context.WithTimeout(parent, m.timeout)
	defer cancel()

	mongoUp, mongoErr := m.checkMongo(ctx)
	natsUp, natsErr := m.checkNATS()

	mongoStatus := "error"
	if mongoUp {
		mongoStatus = "ok"
	}
	natsStatus := "error"
	if natsUp {
		natsStatus = "ok"
	}

	checkedAt := time.Now()

	var prevMongo, prevNats string
	m.mu.Lock()
	prevMongo = m.status.Mongo
	prevNats = m.status.NATS
	m.status.Mongo = mongoStatus
	m.status.NATS = natsStatus
	m.status.CheckedAt = checkedAt
	m.status.CheckCount++
	m.mu.Unlock()

	if prevMongo != mongoStatus {
		if mongoUp {
			m.log.Info("mongodb status changed", "from", prevMongo, "to", mongoStatus)
		} else {
			m.log.Warn("mongodb status changed", "from", prevMongo, "to", mongoStatus, "err", mongoErr)
		}
	}

	if prevNats != natsStatus {
		if natsUp {
			m.log.Info("nats status changed", "from", prevNats, "to", natsStatus)
		} else {
			m.log.Warn("nats status changed", "from", prevNats, "to", natsStatus, "err", natsErr)
		}
	}

	if mongoUp {
		m.m.MongoUp.Set(1)
	} else {
		m.m.MongoUp.Set(0)
		m.m.MongoErrors.Inc()
	}

	if natsUp {
		m.m.NatsUp.Set(1)
	} else {
		m.m.NatsUp.Set(0)
		m.m.NatsErrors.Inc()
	}

	m.m.LastCheckUnix.Set(float64(checkedAt.Unix()))
	m.m.CheckCount.Inc()
}

func (m *Monitor) checkMongo(ctx context.Context) (bool, error) {
	if m.mongo == nil {
		m.log.Info("mongodb connect attempt")
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(m.mongoURI))
		if err != nil {
			m.log.Warn("mongodb connect failed", "err", err)
			return false, err
		}
		m.mongo = c
	}

	if err := m.mongo.Ping(ctx, nil); err != nil {
		m.log.Warn("mongodb ping failed", "err", err)
		_ = m.mongo.Disconnect(context.Background())
		m.mongo = nil
		return false, err
	}

	return true, nil
}

func (m *Monitor) checkNATS() (bool, error) {
	if m.nc == nil || m.nc.IsClosed() {
		m.log.Info("nats connect attempt")
		nc, err := nats.Connect(
			m.natsURL,
			nats.Timeout(2*time.Second),
			nats.MaxReconnects(-1),
			nats.ReconnectWait(time.Second),
		)
		if err != nil {
			m.log.Warn("nats connect failed", "err", err)
			m.nc = nil
			return false, err
		}
		m.nc = nc
	}

	ok := m.nc.IsConnected() && !m.nc.IsClosed()
	if !ok {
		return false, m.nc.LastError()
	}
	return true, nil
}
