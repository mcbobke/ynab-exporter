package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"

	"github.com/mcbobke/ynab-exporter/cmd/ynab-exporter/version"
	"github.com/mcbobke/ynab-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	bindAddr, success := os.LookupEnv("BIND_ADDR")
	if success {
		sugar.Debugf("Bind address will be set to %s", bindAddr)
	} else {
		sugar.Debug("No bind address set, will bind to 0.0.0.0")
		bindAddr = "0.0.0.0"
	}

	port, success := os.LookupEnv("PORT")
	if success {
		sugar.Debugf("Port will be set to %s", port)
	} else {
		sugar.Debug("No port specified, will run on 9090")
		port = "9090"
	}

	sugar.Infof("Starting ynab-exporter version %s, bound to %s on port %s", version.BuildVersion, bindAddr, port)

	go http.ListenAndServe(fmt.Sprintf("%s:%s", bindAddr, port), nil)

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
