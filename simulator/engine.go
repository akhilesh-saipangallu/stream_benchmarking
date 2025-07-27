package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/akhilesh-saipangallu/stream_benchmarking/broadcast"
	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
	"github.com/akhilesh-saipangallu/stream_benchmarking/utils"
)

type CmdArgs struct {
	Suffix     *string
	Separator  *string
	NumTickers *int
	GapMs      *int
}

type Engine struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	cmdArgs     CmdArgs
	broadcaster broadcast.Broadcaster
}

func NewEngine(ctx context.Context, cmdArgs CmdArgs) (Engine, error) {
	var err error
	e := Engine{
		ctx:     ctx,
		wg:      &sync.WaitGroup{},
		cmdArgs: cmdArgs,
	}
	// e.broadcaster, err = broadcast.NewLoggingBroadCaster()
	e.broadcaster, err = broadcast.NewNatsBroadcaster()
	if err != nil {
		err = fmt.Errorf("newEngine: %w", err)
	}
	return e, err
}

func (e Engine) Run() {
	log.Println("Engine.Run: mock price feed started")
	e.simulateAndBroadcastPrices()
}

func (e Engine) Stop() {
	e.broadcaster.Cleanup()
	log.Println("Engine.Run: mock price feed stopped")
}

func (e Engine) GetWg() *sync.WaitGroup {
	return e.wg
}

func (e Engine) simulateAndBroadcastPrices() {
	interval := time.Duration(*e.cmdArgs.GapMs) * time.Millisecond
	for i := 0; i < *e.cmdArgs.NumTickers; i++ {
		ticker := fmt.Sprintf("%s%s%d", *e.cmdArgs.Suffix, *e.cmdArgs.Separator, i)
		e.GetWg().Add(1)

		go func(tkr string) {
			defer e.GetWg().Done()
			timer := time.NewTicker(interval)
			defer timer.Stop()
			price := 100 + (rand.Float64() * 400) // starting range [100 , 500]

			for {
				select {
				case <-e.ctx.Done():
					return
				case <-timer.C:
					price = utils.RandomWalkPrice(price)
					bids, asks := utils.GeneratePriceLevels(price, 6)
					payload := datatype.PriceMessage{
						Ticker:    ticker,
						Timestamp: time.Now().UnixMicro(),
						Bids:      bids,
						Asks:      asks,
					}
					payloadJson, _ := json.Marshal(payload)
					err := e.broadcaster.Broadcast(tkr, string(payloadJson))
					if err != nil {
						log.Printf("simulateAndBroadcastPrices: error in broadcasting: %s\n", err)
					}
				}
			}
		}(ticker)
	}
}
