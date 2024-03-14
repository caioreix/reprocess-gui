package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("running...")

	go func() {
		for i := range 100 {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("waiting...")

	<-done
	fmt.Println("graceful shutdown.")
}
