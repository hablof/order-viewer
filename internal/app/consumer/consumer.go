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

	durableQueueGroup = "order-viewer-group"
	subject           = "orders"
)

const (
	saveOrdeTimeout = 5 * time.Second
)

type Service interface {
	SaveOrder(ctx context.Context, order models.Order) error
}

type SubscriberClient struct {
	service Service
	con     stan.Conn
}

func RegisterStanClient(s Service) (*SubscriberClient, error) {
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

	return &SubscriberClient{
		service: s,
		con:     stanConnection,
	}, nil
}

// RunNconsumers is blocking operation
// returns error if unable to run
// ctx.Done stops consumers, returning nil
func (sc *SubscriberClient) RunNconsumers(ctx context.Context, number int) error {
	ctx, cf := context.WithCancel(ctx)
	defer cf()

	for i := 0; i < number; i++ {
		c := consumer{
			service: sc.service,
			id:      i,
		}
		msgHandler := c.getMessageHandler(ctx)

		subscription, err := sc.con.QueueSubscribe(subject, durableQueueGroup, msgHandler, stan.StartWithLastReceived(), stan.DurableName(durableQueueGroup))
		if err != nil {
			sc.con.Close()
			log.Println(err)

			return err
		}
		defer subscription.Close()
	}

	log.Println("subscribed")

	<-ctx.Done()

	log.Println("consumer recived shutdown signal")

	return nil
}

type consumer struct {
	service Service
	id      int
}

func (c *consumer) getMessageHandler(ctx context.Context) func(msg *stan.Msg) {

	return func(msg *stan.Msg) {

		ctx, cf := context.WithTimeout(ctx, saveOrdeTimeout)
		defer cf()

		log.Printf("consumer #%d recived msg", c.id)

		order := models.Order{}

		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Println("failed to unmarshal data: ", err)
			return
		}

		if err := c.service.SaveOrder(ctx, order); err != nil {
			log.Printf("consumer #%d failed to save order", c.id)
		}
	}
}
