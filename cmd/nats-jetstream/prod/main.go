package main

import (
	"encoding/json"
	"flag"
	"log"
	"nats/types"
	"time"

	"github.com/nats-io/nats.go"
)

const StreamName = "events"

func main() {
	// Connect to a server
	conn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	// JET STREAM
	natsJs, err := nc.Conn.JetStream(
		nats.PublishAsyncMaxPending(1000),
	)
	if err != nil {
		log.Fatal(err)
	}
	if s, err := natsJs.StreamInfo(StreamName); err != nil || s == nil {
		log.Println("Stream ", StreamName, " not found: ", err)
		log.Println("Creating a stream> ", StreamName)
		// Create a stream
		si, err := natsJs.AddStream(
			&nats.StreamConfig{
				Name: StreamName,
				Subjects: []string{
					StreamName + ".*",
				},
				Retention:    nats.WorkQueuePolicy,
				MaxConsumers: 5,
				Storage:      nats.FileStorage,
				Discard:      nats.DiscardOld,
				MaxMsgs:      -1, // unlimitted
				MaxBytes:     -1, // stream size unlimitted
				MaxAge:       365 * 24 * time.Hour,
				MaxMsgSize:   1024,
				Duplicates:   1 * time.Hour,
			},
			nats.PublishAsyncMaxPending(1000),
		)
		log.Println("StreamInfo:", si, "error", err)
	}
	// -----------
	sub := flag.String("subj", "", "a string")
	flag.Parse()
	if *sub == "" {
		log.Fatal("subject required")
	}
	subj := *sub

	msg := &types.Incident{
		Message: "nats fired",
		Time:    time.Now().Local().Format(time.RFC3339),
	}
	msgB, _ := json.Marshal(msg)

	// Simple Publisher
	// TODO: async publish
	pubAck, err := natsJs.Publish(subj, msgB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pubAck)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Published [%s] : '%+v'\n", subj, msg)
}
