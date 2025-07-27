package datatype

import (
	"encoding/json"
	"fmt"
)

type OrderLevel struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type PriceMessage struct {
	Ticker    string       `json:"ticker"`
	Timestamp int64        `json:"timestamp"` // microsecond
	Bids      []OrderLevel `json:"bids"`
	Asks      []OrderLevel `json:"asks"`
}

func (pm *PriceMessage) GetBytes() (data []byte, err error) {
	data, err = json.Marshal(pm)
	if err != nil {
		err = fmt.Errorf("PriceMessage.GetBytes: %w", err)
	}
	return
}

func ParsePriceMessageBytes(msg []byte) (*PriceMessage, error) {
	var (
		pm  PriceMessage
		err error
	)
	err = json.Unmarshal(msg, &pm)
	if err != nil {
		err = fmt.Errorf("ParsePriceMessageBytes: %w", err)
	}
	return &pm, err
}
