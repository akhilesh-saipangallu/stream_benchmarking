package tick_reader

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/nats-io/nats.go"
)

type NatsTickSource struct {
	mu            *sync.Mutex
	conn          *nats.Conn
	subscriptions map[string]*nats.Subscription
}

func newNatsTickSource() (NatsTickSource, error) {
	nts := NatsTickSource{
		mu:            &sync.Mutex{},
		subscriptions: map[string]*nats.Subscription{},
	}
	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("newNatsTickSource: %w", err)
		return nts, err
	}
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		err = fmt.Errorf("newNatsTickSource: error connecting to NATS server: %w", err)
		return nts, err
	}

	nts.conn = nc
	return nts, err
}

func (nts NatsTickSource) Subscribe(ticker string, msgHandler func(string, []byte)) (err error) {
	nts.mu.Lock()
	defer nts.mu.Unlock()

	_, exists := nts.subscriptions[ticker]
	if exists {
		log.Printf("NatsTickSource.Subscribe: %s already subscribed\n", ticker)
		return
	}

	sub, err := nts.conn.Subscribe(
		ticker,
		func(msg *nats.Msg) {
			msgHandler(msg.Subject, msg.Data)
		},
	)
	if err != nil {
		return fmt.Errorf("NatsTickSource.Subscribe: %w", err)
	}
	nts.subscriptions[ticker] = sub

	return
}

func (nts NatsTickSource) Stop(ctx context.Context) {
	nts.mu.Lock()
	defer nts.mu.Unlock()

	for _, sub := range nts.subscriptions {
		sub.Unsubscribe()
	}
	nts.conn.Close()
	log.Println("NatsTickSource.Stop: nats tick source stopped")
}
