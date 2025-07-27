package broadcast

import (
	"fmt"
	"log"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/utils"
	"github.com/nats-io/nats.go"
)

type NatsBroadcaster struct {
	conn    *nats.Conn
	rLogger *utils.LogOnceEveryN
}

func NewNatsBroadcaster() (n NatsBroadcaster, err error) {
	cfg, err := config.Get()
	if err != nil {
		err = fmt.Errorf("newNatsBroadcaster: %w", err)
		return
	}
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		err = fmt.Errorf("newNatsBroadcaster: error connecting to NATS server: %w", err)
		return
	}
	n.conn = nc
	// todo: maybe change this rate limit
	n.rLogger = utils.NewLogOnceEveryN(100)
	return
}

func (n NatsBroadcaster) Broadcast(channel, payload string) error {
	n.rLogger.Log("NatsBroadcaster.Broadcast: %s | %s\n", channel, payload)
	err := n.conn.Publish(channel, []byte(payload))
	return err
}

func (n NatsBroadcaster) Cleanup() {
	n.conn.Close()
	log.Println("NatsBroadcaster.Cleanup: cleanup completed")
}
