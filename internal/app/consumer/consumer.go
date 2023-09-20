package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hablof/order-viewer/internal/models"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	natsURL = "nats://127.0.0.1:4222"
	cluster = "my_cluster"

	clientName = "order-viewer-sub"
	clientID   = "order-viewer-id"
)

const (
	saveOrdeTimeout = 5 * time.Second
)

type Service interface {
	SaveOrder(ctx context.Context, order models.Order) error
}

type Consumer struct {
	con     stan.Conn
	service Service
	id      int
}

func NewConsumer(con stan.Conn, s Service, id int) (*Consumer, error) {

	return &Consumer{
		con:     con,
		service: s,
		id:      id,
	}, nil
}

func (c *Consumer) getMessageHandler(ctx context.Context) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {

		ctx, cf := context.WithTimeout(ctx, saveOrdeTimeout)
		defer cf()

		log.Println("consumer#", c.id, " recived msg")

		data := msg.Data

		order := models.Order{}

		if err := json.Unmarshal(data, &order); err != nil {
			log.Println("failed to unmarshal data: ", err)
			return
		}

		if err := c.service.SaveOrder(ctx, order); err != nil {
			log.Println("consumer#", c.id, "failed to save order")
		}

	}
}

func StanConn() (stan.Conn, error) {
	opts := []nats.Option{nats.Name(clientName)}

	natsConnection, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Println("failed to connect to NATS: ", err)
		return nil, err
	}

	connLostOpt := stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Println("STAN Connection lost, reason:", reason)
	})

	stanConnection, err := stan.Connect(cluster, clientID, stan.NatsConn(natsConnection), connLostOpt)
	if err != nil {
		log.Fatalf("failed to connect to STAN: %v", err)
		return nil, err
	}
	return stanConnection, nil
}

// Run is blocking operation
// returns error if unable to run
// ctx.Done stops consumer, returning nil
func (c *Consumer) Run(ctx context.Context) error {
	ctx, cf := context.WithCancel(ctx)
	defer cf()

	msgHandler := c.getMessageHandler(ctx)

	subscription, err := c.con.QueueSubscribe("orders", "order-viewer-group", msgHandler, stan.StartWithLastReceived(), stan.DurableName("order-viewer-group"))
	if err != nil {
		c.con.Close()
		log.Println(err)

		return err
	}

	log.Println("subscribed")

	<-ctx.Done()
	log.Println("consumer recived shutdown signal")

	if err := subscription.Close(); err != nil {
		log.Println("failed to close STAN subscription: ", err)
	}

	if err := c.con.Close(); err != nil {
		log.Println("failed to close STAN connection: ", err)
	}

	c.con.NatsConn().Close()

	return nil
}
