# Stream Benchmarking

This project simulates a real-time price streaming system using WebSocket and NATS. It includes three components:

1. **Mock Price Feed Generator** : Publishes mock prices to NATS
2. **WebSocket Server** : Reads from NATS and broadcasts to clients
3. **WebSocket Clients** — Simulates clients subscribing to tickers and receiving updates


## 1. Generate Mock Prices

Run mock price feed to publish simulated ticker prices to NATS.

### Configuration
Make sure to update the the following configs in the `config/config.yaml` before running the mock price feed and WS server.
- nats.url

### Using Go
```bash
go run mock_price_feed.go -suffix=bn -count=200 -gap=100
```

### Using Docker
```bash
docker build -t mock_price -f docker_files/Dockerfile.mock_price_feed .
docker run --rm --name mock_price mock_price --suffix=bn --count=200 --gap=1000
```

## 2. Run WebSocket Server

Start the server to broadcast prices over WebSocket.

### Configuration
Make sure to update the the following configs in the `config/config.yaml` before running the WS server.
- nats.url

### Using Go
```bash
go run ws_server.go
```

### Using Docker
```bash
docker build -t ws_server -f docker_files/Dockerfile.ws_server .
docker run --rm -p 8080:8080 -p 2112:2112 --name ws_server ws_server
```

## 3. Run WebSocket Clients

Simulate clients connecting to the WebSocket server.

### Configuration
Make sure to update the the following configs in the `config/config.yaml` before running the WS clients.
- ws_server.endpoint

### Using Go
```bash
go run ws_clients.go --client_count=5 --ticker_count=10 --suffix=bn
```

### Using Docker
```bash
docker build -t ws_clients -f docker_files/Dockerfile.ws_clients .
docker run --rm -p 2212:2212 --name ws_clients ws_clients --client_count=5 --ticker_count=10 --suffix=bn
```

## 4. Run Monitoring Stack (Prometheus + Grafana)
Use Docker Compose to start Prometheus and Grafana for monitoring metrics exposed by the application(s).

### Configuration
Make sure to update the the following configs in the `metrics/prometheus.yml` before running the monitoring stack.
- scrape_configs.static_configs.targets

### Using Docker
```bash
cd metrics
docker-compose up
```