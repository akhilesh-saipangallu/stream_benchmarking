package tick_reader

import "github.com/akhilesh-saipangallu/stream_benchmarking/datatype"

type Client interface {
	SendTick(*datatype.PriceMessage)
}