package ws

type SubscriptionMessage struct {
	Action  string   `json:"action"` // "subscribe" or "unsubscribe"
	Tickers []string `json:"tickers"`
}
