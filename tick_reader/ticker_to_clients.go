package tick_reader

import (
	"sync"
)

type tickerClients struct {
	mu      *sync.RWMutex
	mapping map[string][]Client
}

func newTickerClients() tickerClients {
	return tickerClients{
		mu:      &sync.RWMutex{},
		mapping: map[string][]Client{},
	}
}

func (tc *tickerClients) tickerExists(ticker string) (exists bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	_, exists = tc.mapping[ticker]
	return
}

func (tc *tickerClients) addClient(ticker string, client Client) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.mapping[ticker] = append(tc.mapping[ticker], client)
}

func (tc *tickerClients) getClients(ticker string) []Client {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.mapping[ticker]
}
