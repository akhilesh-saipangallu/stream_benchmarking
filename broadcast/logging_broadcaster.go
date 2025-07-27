package broadcast

import "log"

type LoggingBroadCaster struct {}

func NewLoggingBroadCaster() (lb LoggingBroadCaster, err error) {
	return lb, err
}

func (lb LoggingBroadCaster) Broadcast(channel, payload string) error {
	log.Printf("LoggingBroadCaster.Broadcast: %s | %s\n", channel, payload)
	return nil
}

func (lb LoggingBroadCaster) Cleanup() {
	log.Println("LoggingBroadCaster.Broadcast: cleanup complete")
}
