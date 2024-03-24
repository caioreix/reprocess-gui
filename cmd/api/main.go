package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"reprocess-gui/internal/apps/api/adapter/driving/http"
	"reprocess-gui/internal/apps/api/adapter/repository/mongodb"
	"reprocess-gui/internal/apps/api/config"
	"reprocess-gui/internal/apps/api/core/service"
	"reprocess-gui/internal/logger"
)

func main() {
	configPath := flag.String("cpath", ".", "config path")
	flag.Parse()

	config, err := config.New(*configPath)
	if err != nil {
		panic(err)
	}

	log, err := logger.New(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	mongo, err := mongodb.New(config)
	if err != nil {
		log.Fatal("failed creating mongo connection", []logger.Field{
			{Key: "error", Value: err},
		}...)
	}
	defer mongo.Close(ctx)

	tableCollection := mongo.Database(config.Mongo.TableDatabase).Collection(config.Mongo.TableCollection)

	tableRepository := mongodb.NewTableRepository(config, log, tableCollection)
	tableService := service.NewTableService(config, log, tableRepository)
	tableHandler := http.NewTableHandler(config, log, tableService)

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	router, err := http.NewRouter(addr, tableHandler)
	if err != nil {
		log.Fatal("failed creating the router", []logger.Field{
			{Key: "error", Value: err},
		}...)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(fmt.Sprintf("Running on: %s\n", addr))
		err := router.Serve()
		log.Fatal("failed server running", []logger.Field{
			{Key: "error", Value: err},
		}...)
	}()

	<-done
	log.Warn(fmt.Sprintf("Running on: %s\n", addr))
}
