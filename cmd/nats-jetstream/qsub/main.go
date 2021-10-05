package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	StreamName   = "events"
	ConsumerName = "cmdr_durable"
)

func main() {
	// Connect to a server
	// opt := nats.Options{}
	nConn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	// JET STREAM
	natsJs, err := nConn.JetStream(
		nats.PublishAsyncMaxPending(1000),
	)
	if err != nil {
		log.Fatal(err)
	}
	if s, err := natsJs.StreamInfo(StreamName); err != nil || s == nil {
		log.Fatal("Stream ", StreamName, " not found: ", err)
	}
	if c, err := natsJs.ConsumerInfo(StreamName, ConsumerName); err != nil || c == nil {
		log.Println("Consumer ", StreamName, " not found: ", err)
		log.Println("Creating a consumer> ", ConsumerName)

		ci, err := natsJs.AddConsumer(StreamName, &nats.ConsumerConfig{
			Durable:        ConsumerName,
			DeliverPolicy:  nats.DeliverAllPolicy,
			AckPolicy:      nats.AckExplicitPolicy,
			AckWait:        5 * time.Second,
			MaxAckPending:  1000,
			MaxDeliver:     2,
			ReplayPolicy:   nats.ReplayInstantPolicy,
			DeliverSubject: ConsumerName + "." + StreamName + ".processed",
			FilterSubject:  StreamName + ".>",
		})
		log.Println("ConsumerInfo:", ci, "error", err)
	}
	// -----------

	// topic
	substr := flag.String("subj", "", "a string")
	flag.Parse()
	if *substr == "" {
		log.Fatal("subject required")
	}
	subj := *substr

	// -------------------------------------------------------------------------
	// steaming no loadbalncing
	// -------------------------------------------------------------------------
	// Simple Async Queue Subscriber
	asyncQSub, err := natsJs.QueueSubscribe(subj, ConsumerName, func(msg *nats.Msg) {
		log.Println("Msg>>>: ", string(msg.Data))
		// msg.Nak()
		if err := msg.Ack(); err != nil {
			log.Println("Ack err: ", err.Error())
		}
	},
		nats.Durable(ConsumerName),
		nats.DeliverAll(),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.AckWait(5*time.Second),
		nats.MaxAckPending(1000),
		nats.MaxDeliver(2),
	)
	if err != nil {
		log.Fatal("async_queue_sub err: " + err.Error())
	}
	log.Printf("Listening on [%s], queue group [%s]\n", subj, ConsumerName)
	if !asyncQSub.IsValid() {
		log.Fatal("async_queue_sub closed")
	}

	// wait forever
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Draining...")
	if err := nConn.Drain(); err != nil {
		log.Println("Drain err: ", err.Error())
	}
}
