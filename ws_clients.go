package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
	"github.com/akhilesh-saipangallu/stream_benchmarking/metrics"
	"github.com/akhilesh-saipangallu/stream_benchmarking/utils"
	"github.com/akhilesh-saipangallu/stream_benchmarking/ws"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type WsClientCmdArgs struct {
	suffix     *string
	separator  *string
	numTickers *int
}

func generateTickers(cmdArgs WsClientCmdArgs) (tickers []string) {
	for i := 0; i < *cmdArgs.numTickers; i++ {
		ticker := fmt.Sprintf("%s%s%d", *cmdArgs.suffix, *cmdArgs.separator, i)
		tickers = append(tickers, ticker)
	}
	return
}

type WsClient struct {
	conn               *websocket.Conn
	rLogger            utils.LogOnceEveryN
	tickersToSubscribe []string
}

func NewWsClient(tickersToSubscribe []string) *WsClient {
	return &WsClient{
		rLogger:            *utils.NewLogOnceEveryN(100),
		tickersToSubscribe: tickersToSubscribe,
	}
}

func (c *WsClient) Start(url string) (err error) {
	err = c.Connect(url)
	if err != nil {
		return err
	}

	err = c.Subscribe()
	if err != nil {
		return err
	}
	go c.Reader()

	return
}

func (c *WsClient) Connect(url string) (err error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("WsClient.Start: dial error: %w", err)
	}
	c.conn = conn
	return
}

func (c *WsClient) Subscribe() (err error) {
	subMsg := ws.SubscriptionMessage{
		Action:  "subscribe",
		Tickers: c.tickersToSubscribe,
	}
	err = c.conn.WriteJSON(subMsg)
	if err != nil {
		err = fmt.Errorf("WsClient.Subscribe: write error: %w", err)
	}
	return
}

func (c *WsClient) Reader() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}

		now := time.Now().UnixMicro()
		pm, err := datatype.ParsePriceMessageBytes(msg)
		if err != nil {
			log.Printf("WsClient.Reader: error: %s\n", err)
			continue
		}
		latency := float64(now - pm.Timestamp) // latency in µs
		metrics.EndToEndLatency.Observe(latency)

		// log.Printf("received: %s", msg)
		c.rLogger.Log("received: %s", msg)
	}
}

func (c WsClient) Stop() {
	c.conn.Close()
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	numWsClientsPtr := flag.Int("client_count", 100, "Number of WS clients to create")
	numTickersPtr := flag.Int("ticker_count", 300, "Number of tickers to subscribe")
	suffixPtr := flag.String("suffix", "nifty", "Ticker suffix")
	separatorPtr := flag.String("sep", "_", "Ticker separator")
	flag.Parse()

	cmdArgs := WsClientCmdArgs{
		suffix:     suffixPtr,
		separator:  separatorPtr,
		numTickers: numTickersPtr,
	}

	metrics.InitNatsMetrics()
	config.Init()
	cfg, err := config.Get()
	if err != nil {
		log.Println("main: %s", err)
		return
	}

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2212", nil)

	url := cfg.WsServer.WSEndpoint
	tickers := generateTickers(cmdArgs)

	var (
		totalConnectedClients int
		wsClients             []*WsClient
	)
	for i := 0; i < *numWsClientsPtr; i++ {
		wsClient := NewWsClient(tickers)
		wsErr := wsClient.Start(url)
		if wsErr != nil {
			log.Printf("==ERROR==: failed to initiate ws client id:%d\n", i)
			continue
		}
		totalConnectedClients += 1
		wsClients = append(wsClients, wsClient)
	}
	log.Println("==INFO==: totalConnectedClients:", totalConnectedClients)

	// wait for interrupt
	<-ctx.Done()
	log.Println("interrupt received; closing connection...")

	for _, wsClient := range wsClients {
		wsClient.Stop()
	}
	log.Println("stopped all WS clients; exiting")
}
