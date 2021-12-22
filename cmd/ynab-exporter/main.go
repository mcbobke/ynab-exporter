package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"

	"github.com/mcbobke/ynab-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	modVersion string = "0.0.1"
)

func main() {
	loggerConfig := zap.NewDevelopmentConfig()

	level, success := os.LookupEnv("LOG_LEVEL")
	if success {
		switch strings.ToLower(level) {
		case "debug":
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "info":
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		case "warn":
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		case "fatal":
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
		default:
			loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		}
	} else {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, _ := loggerConfig.Build()
	sugar := logger.Sugar()
	defer sugar.Sync()

	token, success := os.LookupEnv("YNAB_API_TOKEN")
	if success {
		sugar.Info("YNAB token found")
	} else {
		sugar.Fatal("YNAB token not found")
	}

	ynabCollector := collector.New(token, sugar)
	prometheus.MustRegister(ynabCollector)
	http.Handle("/metrics", promhttp.Handler())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sugar.Infof("Starting ynab-exporter version %s", modVersion)

	go http.ListenAndServe("localhost:9090", nil)

	signal := <-quit
	switch signal {
	case syscall.SIGINT:
		sugar.Info("Interrupt received, shutting down")
		os.Exit(0)
	case syscall.SIGTERM:
		sugar.Info("Termination received, shutting down")
		os.Exit(0)
	default:
		sugar.Warnf("Unknown signal [%s] received, shutting down", signal)
		os.Exit(1)
	}
}
