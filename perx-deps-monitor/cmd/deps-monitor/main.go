package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"perx-deps-monitor/internal/httpapi"
	"perx-deps-monitor/internal/monitor"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDB := os.Getenv("MONGO_DB")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	if mongoHost == "" || mongoPort == "" || mongoDB == "" || mongoUser == "" || mongoPassword == "" {
		log.Fatal("Mongo env vars are not fully set")
	}

	mongoURI :=
		"mongodb://" +
			mongoUser + ":" + mongoPassword +
			"@" + mongoHost + ":" + mongoPort +
			"/" + mongoDB

	natsHost := os.Getenv("NATS_HOST")
	natsPort := os.Getenv("NATS_PORT")
	natsUser := os.Getenv("NATS_USER")
	natsPassword := os.Getenv("NATS_PASSWORD")

	if natsHost == "" || natsPort == "" || natsUser == "" || natsPassword == "" {
		log.Fatal("NATS env vars are not fully set")
	}

	natsURL :=
		"nats://" +
			natsUser + ":" + natsPassword +
			"@" + natsHost + ":" + natsPort

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	mon := monitor.New(monitor.Config{
		MongoURI: mongoURI,
		NATSURL:  natsURL,
		Interval: 2 * time.Second,
		Timeout:  2 * time.Second,
	}, prometheus.DefaultRegisterer)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()
	go mon.Run(appCtx)

	api := &httpapi.API{Mon: mon}
	mux := api.Mux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server is running on port %s ...", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	appCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	mon.Close(ctx)
}
