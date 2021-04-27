## Developing With NATS Streaming

- **Connecting to NATS Streaming**
  clusterID - defined by the server configuration
  clientID - defined by the client

  ```go
  sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
  sc, err := sc.Connect(clusterID, clientName, stan.MaxPubAcksInflight(1000))
  ```

- **Publishing to a Channel**

  ```go
  // publish
  err := sc.Publish("foo", []byte("Hello World"))
  // async publich
  ackHandler := func(ackedNuid string, err error){ ... }
  nuid, err := sc.PublishAsync("foo", []byte("Hello World"), ackHandler)
  // connect with publisher acks max in flights
  sc, err := sc.Connect(clusterID, clientName, stan.MaxPubAcksInflight(1000))
  ```

- **Receiving Messages from a Channel**
  I need `Queue+Durable` subscription.
  ```go
  // subscribe to all available messages
  qsub, err := sc.QueueSubscribe(channelName, queueName,
    func(m *stan.Msg) {
        //...
        m.Ack()
    },
    stan.DurableName("durable-name"),
    stan.DeliverAllAvailable(),
    stan.SetManualAckMode(),
    stan.AckWait(5*time.Second),
    stan.MaxInflight(5),
    )
  ```
