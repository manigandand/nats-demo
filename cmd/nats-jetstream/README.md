## NATS Streaming

- Close client connections
  > It is therefore strongly recommended for clients to close their connection when the application exits, otherwise the server will consider these clients connected (sending data, etc...) until it detects missing heartbeats.
- client ID
  > no two connections with the same client ID will be able to run concurrently.

1. Stream

   - Durable subscriptions
   - RetensionPolicy
     - WorkQueuePolicy
     - LimitPolicy
   - dupe-window
   - ack
     - AckExplict
     - AckAll
     - AckNone
   - Ack wait

2. Consumer

   - DeliveryPolicy
     - all
   - MaxInFlight

3. Publisher

   - MaxPubAcksInflight
     > publisher max in flight?

4. Subscriber
   - Durable
     > If an application wishes to resume message consumption from where it previously stopped, it needs to create a durable subscription. It does so by providing a `durable name`, which is combined with the `client ID` provided when the client created its connection.
   - Queue Group
     > When consumers want to consume from the same channel but each receive a different message, as opposed to all receiving the same messages, they need to create a queue subscription. (`queue group name` and `durable name`)

## Developing With NATS JetStreaming

- **Connecting to NATS JetStreaming**

  ```go
  conn, err := nats.Connect("nats://127.0.0.1:4222")
  natsJs, err := nc.Conn.JetStream(
  	nats.PublishAsyncMaxPending(1000),
  )
  ```

- **Publishing to a Channel**

  ```go
  // publish
  err := natsJs.Publish("foo", []byte("Hello World"))
  // async publich
  ackHandler := func(ackedNuid string, err error){ ... }
  nuid, err := natsJs.PublishAsync("foo", []byte("Hello World"), ackHandler)
  ```

- **Receiving Messages from a Channel**
  I need `Queue+Durable` subscription.
  ```go
  // subscribe to all available messages
  qsub, err := natsJs.QueueSubscribe(subject, queueName,
    func(m *stan.Msg) {
        //...
        m.Ack()
    },
    nats.Durable("durable-name"),
    nats.DeliverAll(),
    nats.AckExplicit(),
    nats.AckWait(5*time.Second),
    nats.MaxInflight(5),
    )
  ```
