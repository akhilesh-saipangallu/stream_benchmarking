package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/tick_reader"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	ctx        context.Context
	numClients int
	clients    map[string]*Client
	clientsMux sync.RWMutex
	tickReader *tick_reader.Reader
	upgrader   websocket.Upgrader
}

func NewServer(ctx context.Context) (*Server, error) {
	tr, trErr := tick_reader.NewReader(ctx)
	if trErr != nil {
		return nil, fmt.Errorf("NewServer: %w", trErr)
	}
	return &Server{
		ctx:        ctx,
		numClients: 0,
		clients:    map[string]*Client{},
		clientsMux: sync.RWMutex{},
		tickReader: tr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Accepting all requests
			},
		},
	}, nil
}

func (wss *Server) Start() (err error) {
	http.HandleFunc("/ws", wss.HandleWs)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("Server.Start: %w", err)
		return
	}
	log.Printf("Server.Start: server running on :%d\n", cfg.WsServer.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.WsServer.Port), nil)
	if err != nil {
		err = fmt.Errorf("Server.Start: %w", err)
	}
	return
}

func (wss *Server) HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Server.HandleWs: webSocket upgrade failed: %v\n", err)
		return
	}

	err = wss.startClient(conn)
	if err != nil {
		log.Printf("Server.HandleWs: %s\n", err)
	}
}

func (wss *Server) startClient(conn *websocket.Conn) error {
	ctx, cancel := context.WithCancel(wss.ctx)
	wss.clientsMux.Lock()
	defer wss.clientsMux.Unlock()

	cfg, err := config.Get()
	if err != nil {
		cancel()
		err = fmt.Errorf("startClient: %w", err)
		return err
	}

	wss.numClients += 1
	clientConId := fmt.Sprintf("cid_%d", wss.numClients)
	client := NewClient(
		ctx, cancel, conn, clientConId, cfg.WsServer.WriterBufferSize, wss.tickReader,
	)

	wss.clients[clientConId] = client
	log.Printf("Server.HandleWs: client %s connected\n", client.ConId)

	client.Start()

	return nil
}

func (wss *Server) removeClient(c *Client) {
	wss.clientsMux.Lock()
	defer wss.clientsMux.Unlock()

	delete(wss.clients, c.ConId)
}

func (wss *Server) Stop() {
	log.Println("Server.Stop: stopped server")
}
