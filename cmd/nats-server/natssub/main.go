package main

import (
	"flag"
	"log"
	"nats/types"

	nats "github.com/nats-io/nats.go"
)

func main() {
	// Connect to a server
	// opt := nats.Options{}
	conn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

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

	// Simple Async Subscriber
	_, err = nc.Subscribe(subj, func(m *types.Incident) {
		log.Printf("'%+v'\n", m)
	})
	if err != nil {
		log.Fatal("async_sub err: " + err.Error())
	}
	// asyncSub.Queue = "test"

	// ---------
	/*
		// Simple sync Subscriber, available in normal conn
		sub, err := nc.SubscribeSync(subj)
		if err != nil {
			log.Fatal("sub err: " + err.Error())
		}
		sub.Queue = "test"

		for {
			fmt.Println(sub.Queue, sub.Subject)
			m, err := sub.NextMsg(time.Second * 5)
			if err != nil {
				fmt.Println("stream sub next msg: ", err)
				continue
			}

			fmt.Println("Stream Msg: ", string(m.Data))
		}
	*/

	// wait forever
	wait := make(chan bool, 1)
	<-wait

	// Responding to a request message
	/*
		nc.Subscribe("request", func(m *nats.Msg) {
			m.Respond([]byte("answer is 42"))
		})
	*/

	// Channel Subscriber
	/*
		ch := make(chan *nats.Msg, 1000)
		sub, err = nc.ChanSubscribe(subj, ch)
		msg := <-ch
		fmt.Println("Stream Msg: ", string(msg.Data))

		// Unsubscribe
		sub.Unsubscribe()

		// Drain
		sub.Drain()
	*/

	// Requests
	// msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)

	// Replies
	// nc.Subscribe("help", func(m *nats.Msg) {
	// 	nc.Publish(m.Reply, []byte("I can help!"))
	// })

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}
