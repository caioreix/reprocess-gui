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
	configPath := flag.String("cpath", "", "config path")
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

	rowCollection := mongo.Database(config.Mongo.RowDatabase).Collection(config.Mongo.RowCollection)
	rowRepository := mongodb.NewRowRepository(config, log, rowCollection)
	rowService := service.NewRowService(config, log, rowRepository)
	rowHandler := http.NewRowHandler(config, log, rowService)

	consumerCollection := mongo.Database(config.Mongo.ConsumerDatabase).Collection(config.Mongo.ConsumerCollection)
	consumerRepository := mongodb.NewConsumerRepository(config, log, consumerCollection)
	consumerService := service.NewConsumerService(config, log, consumerRepository)
	consumerHandler := http.NewConsumerHandler(config, log, consumerService)

	r := http.Router{
		TableHandler:    tableHandler,
		RowHandler:      rowHandler,
		ConsumerHandler: consumerHandler,
	}
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	router, err := r.NewRouter(addr, config.Server.ReadHeaderTimeout)
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
	log.Warn("finished graceful shutdown")
}
