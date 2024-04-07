package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"reprocess-gui/internal/apps/worker/config"
	"reprocess-gui/internal/logger"
)

func main() {
	configPath := flag.String("cpath", ".", "config path")
	flag.Parse()

	config, err := config.New(*configPath)
	if err != nil {
		panic(err)
	}

	loggerCfg := logger.Config{
		Level:       config.Log.Level,
		Environment: config.Environment,
	}
	log, err := loggerCfg.New()
	if err != nil {
		panic(err)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("running...")
	}()

	log.Info("Waiting for shutdown signal...")
	<-done
	log.Warn("Finished graceful shutdown")
}
