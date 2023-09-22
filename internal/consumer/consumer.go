package consumer

import (
	"context"
	"encoding/json"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"

	"github.com/hablof/order-viewer/internal/models"
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

	log := log.Logger.With().Str("func", "RegisterStanClient").Caller().Logger()

	opts := []nats.Option{nats.Name(clientName)}

	natsConnection, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to NATS")
		return nil, err
	}

	log.Info().Msg("connected to NATS")

	connLostOpt := stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		log.Error().Err(err).Msg("STAN Connection lost")
	})

	stanConnection, err := stan.Connect(cluster, clientID, stan.NatsConn(natsConnection), connLostOpt)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to STAN")
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

	log := log.Logger.With().Str("func", "SubscriberClient.RunNconsumers").Caller().Logger()

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
			log.Error().Err(err).Send()

			return err
		}
		log.Info().Int("consumer_id", c.id).Msg("is running")
		defer subscription.Close()
	}

	log.Info().Msg("subscribed")

	<-ctx.Done()

	log.Info().Msg("consumers recived shutdown signal")

	return nil
}

type consumer struct {
	service Service
	id      int
}

func (c *consumer) getMessageHandler(ctx context.Context) func(msg *stan.Msg) {

	return func(msg *stan.Msg) {

		log := log.Logger.With().Str("func", "consumer.MessageHandler").
			Caller().Int("consumer_id", c.id).Logger()

		ctx, cf := context.WithTimeout(ctx, saveOrdeTimeout)
		defer cf()

		log.Debug().Msg("recived msg")

		order := models.Order{}

		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal data")
			return
		}

		if err := c.service.SaveOrder(ctx, order); err != nil {
			log.Error().Err(err).Msg("failed to save order")
		}
	}
}
