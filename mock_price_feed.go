package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/simulator"
)

func main() {
	suffixPtr := flag.String("suffix", "nifty", "Ticker suffix")
	separatorPtr := flag.String("sep", "_", "Ticker separator")
	numTickersPtr := flag.Int("count", 300, "Number of tickers to generate")
	gapMsPtr := flag.Int("gap", 100, "Interval between price updates in milliseconds")
	flag.Parse()

	cmdArgs := simulator.CmdArgs{
		Suffix:     suffixPtr,
		Separator:  separatorPtr,
		NumTickers: numTickersPtr,
		GapMs:      gapMsPtr,
	}

	config.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	simulatorEngine, err := simulator.NewEngine(ctx, cmdArgs)
	if err != nil {
		log.Fatalf("main: exiting with an error: %w", err)
	}
	go func() {
		simulatorEngine.Run()
	}()

	<-ctx.Done()
	log.Println("main: exit signal received")
	simulatorEngine.GetWg().Wait()
	simulatorEngine.Stop()
	log.Println("main: cleanup completed; shutting down")
}
