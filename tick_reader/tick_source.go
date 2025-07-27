package tick_reader

import "context"

type TickSource interface {
	Subscribe(ticker string, msgHandler func(string, []byte)) error
	Stop(ctx context.Context)
}
