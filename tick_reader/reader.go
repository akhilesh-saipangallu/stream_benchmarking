package tick_reader

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
	"github.com/akhilesh-saipangallu/stream_benchmarking/metrics"
)

type Reader struct {
	ctx             context.Context
	tickSource      TickSource
	tickerToClients tickerClients
}

func NewReader(ctx context.Context) (*Reader, error) {
	nts, err := newNatsTickSource()
	if err != nil {
		err = fmt.Errorf("NewReader: %w", err)
		return nil, err
	}
	return &Reader{
		ctx:             ctx,
		tickSource:      nts,
		tickerToClients: newTickerClients(),
	}, err
}

func (tr Reader) Subscribe(ticker string, client Client) {
	tickerExists := tr.tickerToClients.tickerExists(ticker)
	tr.tickerToClients.addClient(ticker, client)

	if !tickerExists {
		tr.tickSource.Subscribe(ticker, func(s string, b []byte) { tr.deliverTick(ticker, b) })
	}
}

func (tr Reader) deliverTick(ticker string, msg []byte) {
	now := time.Now().UnixMicro()
	metrics.NatsMessagesReceived.Inc()

	// log.Println(string(msg))
	pm, err := datatype.ParsePriceMessageBytes(msg)
	if err != nil {
		log.Printf("deliverTick: %s\n", err)
		return
	}

	latency := float64(now - pm.Timestamp) // latency in µs
	metrics.NatsMessageLatency.Observe(latency)

	clients := tr.tickerToClients.getClients(ticker)

	for _, client := range clients {
		client.SendTick(pm)
	}
}

func (r Reader) Stop(ctx context.Context) {
	r.tickSource.Stop(ctx)
}
