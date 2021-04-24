package main

import (
	"flag"
	"log"

	nats "github.com/nats-io/nats.go"
)

func main() {
	// Connect to a server
	// opt := nats.Options{}
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	// topic
	substr := flag.String("subj", "", "a string")
	// channel
	chn := flag.String("chn", "", "a string")
	flag.Parse()
	if *substr == "" {
		log.Fatal("subject required")
	}
	subj := *substr

	// -------------------------------------------------------------------------
	// steaming no loadbalncing
	// -------------------------------------------------------------------------
	// Simple Async Queue Subscriber
	asyncQSub, err := nc.QueueSubscribe(subj, *chn, func(m *nats.Msg) {
		log.Println("Msg>>>: ", string(m.Data))
	})
	if err != nil {
		log.Fatal("async_queue_sub err: " + err.Error())
	}
	log.Println("Sub: ", subj, " Queue: ", *chn)
	if !asyncQSub.IsValid() {
		log.Fatal("async_queue_sub closed")
	}

	// wait forever
	ch := make(chan bool, 1)
	<-ch
}
