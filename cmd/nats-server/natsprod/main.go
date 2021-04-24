package main

import (
	"flag"
	"log"
	"nats/types"
	"time"

	"github.com/nats-io/nats.go"
)

func main() { // Connect to a server
	conn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

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

	// Simple Publisher
	err = nc.Publish(subj, msg)
	if err != nil {
		log.Fatal(err)
	}
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Published [%s] : '%+v'\n", subj, msg)
}
