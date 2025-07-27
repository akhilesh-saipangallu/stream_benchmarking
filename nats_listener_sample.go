package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to local NATS server
	nc, err := nats.Connect(nats.DefaultURL) // "nats://localhost:4222"
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	subject := "bn_0"

	// Subscribe to subject
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("Received message on [%s]: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}

	fmt.Printf("Listening for messages on [%s]...\n", subject)

	// Keep the process running
	select {}
}
