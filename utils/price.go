package utils

import (
	"math/rand/v2"

	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
)

func RandomWalkPrice(current float64) float64 {
	delta := rand.Float64()*0.002 - 0.001 // range [-0.001, 0.001]
	return current * (1 + delta)
}

// returns random bids and asks
func GeneratePriceLevels(lastPrice float64, levels int) ([]datatype.OrderLevel, []datatype.OrderLevel) {
	step := 0.5

	bids := make([]datatype.OrderLevel, levels)
	asks := make([]datatype.OrderLevel, levels)

	for i := 0; i < levels; i++ {
		bids[i] = datatype.OrderLevel{
			Price:    lastPrice - step*float64(i+1),
			Quantity: rand.IntN(100) + 50, // 50–149
		}
		asks[i] = datatype.OrderLevel{
			Price:    lastPrice + step*float64(i+1),
			Quantity: rand.IntN(100) + 50,
		}
	}

	return bids, asks
	// return datatype.PriceMessage{
	// 	Ticker:    ticker,
	// 	Timestamp: time.Now().UnixMicro(),
	// 	Bids:      bids,
	// 	Asks:      asks,
	// }
}
