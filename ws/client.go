package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
	"github.com/akhilesh-saipangallu/stream_benchmarking/metrics"
	"github.com/akhilesh-saipangallu/stream_benchmarking/tick_reader"
	"github.com/gorilla/websocket"
)

type Client struct {
	ctx        context.Context
	cancel     context.CancelFunc
	lock       sync.Mutex
	Conn       *websocket.Conn
	ConId      string
	tokens     map[string]struct{}
	sendCh     chan []byte
	tickReader *tick_reader.Reader
	stopOnce   sync.Once
}

func NewClient(
	ctx context.Context,
	cancel context.CancelFunc,
	conn *websocket.Conn,
	conId string,
	sendChBufferSize int,
	tickReader *tick_reader.Reader,
) *Client {
	return &Client{
		ctx:        ctx,
		cancel:     cancel,
		lock:       sync.Mutex{},
		Conn:       conn,
		ConId:      conId,
		tokens:     make(map[string]struct{}),
		sendCh:     make(chan []byte, sendChBufferSize),
		tickReader: tickReader,
	}
}

func (c *Client) Start() {
	go c.Reader()
	go c.Writer()
}

func (c *Client) HandleSubscribe(tokens []string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, token := range tokens {
		_, exists := c.tokens[token]
		if exists {
			log.Printf("Client.HandleSubscribe: %s already subscribed to: %s\n", c.ConId, token)
			continue
		}

		c.tickReader.Subscribe(token, c)
		log.Printf("Client.HandleSubscribe: %s subscribed to: %s\n", c.ConId, token)
	}
}

func (c *Client) HandleUnsubscribe(tokens []string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, token := range tokens {
		_, exists := c.tokens[token]
		if !exists {
			log.Printf("Client.handleUnsubscribe: unable to unsubscribe as %s is not subscribed to: %s\n", c.ConId, token)
			continue
		}
		delete(c.tokens, token)
		// TODO: unsubscribe on NATS
		log.Printf("Client.handleUnsubscribe: %s unsubscribed to: %s\n", c.ConId, token)
	}
}

func (c *Client) Reader() {
	msgChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		for {
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Client.Reader: %s stopping reader", c.ConId)
			return

		case msg := <-msgChan:
			var subMsg SubscriptionMessage
			if err := json.Unmarshal(msg, &subMsg); err != nil {
				log.Printf("Client.Reader: %s: invalid subscription message: %s\n", c.ConId, string(msg))
				continue
			}

			switch subMsg.Action {
			case "subscribe":
				c.HandleSubscribe(subMsg.Tickers)
			case "unsubscribe":
				c.HandleUnsubscribe(subMsg.Tickers)
			default:
				log.Printf("Client.Reader: %s: unknown action: %s\n", c.ConId, subMsg.Action)
			}

		case err := <-errChan:
			log.Printf("Client.Reader: %s: read error: %v\n", c.ConId, err)
			return
		}
	}
}

func (c *Client) Writer() {
	defer func() {
		c.Stop()
	}()

	for {
		select {
		case <-c.ctx.Done():
			log.Printf("Client.Writer: stopped for client %s", c.ConId)
			return

		case msg, ok := <-c.sendCh:
			if !ok {
				return // channel closed
			}

			// set write deadline
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Client.Writer: error for client %s: %v", c.ConId, err)
				return
			}
		}
	}
}

func (c *Client) SendTick(pm *datatype.PriceMessage) {
	if pm == nil {
		return
	}

	msg, err := pm.GetBytes()
	if err != nil {
		log.Printf("Client.SendTick: %s\n", err)
		return
	}

	select {
	case <-c.ctx.Done():
		// client writer is gone; stop trying
		// TODO: Implement unsubscribe to avoid this block getting triggered
	case c.sendCh <- msg:
		// messaged pushed to the channel
	default:
		// channel is full; drop the message
		metrics.WsServerDroppedMessages.Inc()
		log.Printf("Send buffer full for client %s", c.ConId)
	}
}

func (c *Client) Stop() {
	c.stopOnce.Do(func() {
		c.cancel()
		c.Conn.Close()
	})
}
