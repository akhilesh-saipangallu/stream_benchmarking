package broadcast

type Broadcaster interface {
	Broadcast(channel, payload string) error
	Cleanup()
}
