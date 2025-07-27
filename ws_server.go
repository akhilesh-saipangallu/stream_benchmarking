package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/metrics"
	"github.com/akhilesh-saipangallu/stream_benchmarking/ws"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config.Init()
	metrics.InitNatsMetrics()

	server, err := ws.NewServer(ctx)
	if err != nil {
		log.Printf("main: %s\n", err)
		stop()
	} else {
		go func() {
			err := server.Start()
			if err != nil {
				log.Printf("main: %v\n", err)
				stop()
			}
			return
		}()
	}

	<-ctx.Done()
	log.Println("main: exit signal received")
	server.Stop()
	log.Println("main: cleanup completed; shutting down")
}
