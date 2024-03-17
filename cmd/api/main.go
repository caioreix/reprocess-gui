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
)

func main() {
	configPath := flag.String("cpath", ".", "config path")
	flag.Parse()

	config, err := config.New(*configPath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	mongo, err := mongodb.New(config)
	if err != nil {
		panic(err)
	}
	defer mongo.Close(ctx)

	tableCollection := mongo.Database(config.Mongo.TableDatabase).Collection(config.Mongo.TableCollection)

	tableRepository := mongodb.NewTableRepository(config, tableCollection)
	tableService := service.NewTableService(config, tableRepository)
	tableHandler := http.NewTableHandler(config, tableService)

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	router, err := http.NewRouter(addr, tableHandler)
	if err != nil {
		panic(err)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Running on: %s\n", addr)
		router.Serve()
	}()

	<-done
	fmt.Println("graceful shutdown.")
}
