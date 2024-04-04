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

	log, err := logger.New(config.Log.Level)
	if err != nil {
		panic(err)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("running...")
	}()

	<-done
	log.Warn("finished graceful shutdown")
}
