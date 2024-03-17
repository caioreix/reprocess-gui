package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"reprocess-gui/internal/api/adapter/driving/http"
	"reprocess-gui/internal/api/adapter/repository/mongodb"
	"reprocess-gui/internal/api/core/service"
)

func main() {
	ctx := context.Background()

	mongo, err := mongodb.New("mongodb://caio:secret@localhost:27017")
	if err != nil {
		panic(err)
	}
	defer mongo.Close(ctx)

	tableCollection := mongo.Database("api").Collection("tables")

	tableRepository := mongodb.NewTableRepository(tableCollection)
	tableService := service.NewTableService(tableRepository)
	tableHandler := http.NewTableHandler(tableService)

	router, err := http.NewRouter(":8080", tableHandler)
	if err != nil {
		panic(err)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("running...")

	go func() {
		router.Serve()
	}()

	fmt.Println("waiting...")

	<-done
	fmt.Println("graceful shutdown.")
}
