package service

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hablof/order-viewer/internal/models"
)

// cache / repository errors
var (
	ErrNotFound             = errors.New("not found")
	ErrDuplicatesNotAllowed = errors.New("not allowed to store duplicates")
)

// service errors
var (
	ErrOrderNotFound = errors.New("order not found")
)

type Repository interface {
	ReadAll(ctx context.Context) (map[string]models.Order, error)
	InsertOrder(ctx context.Context, order models.Order) error
}

type Cache interface {
	LoadCache(cache map[string]models.Order)
	Get(OrderId string) (models.Order, error)
	Set(order models.Order) error
}

type Service struct {
	c Cache
	r Repository
}

func NewService(c Cache, r Repository) (*Service, error) {

	log := log.Logger.With().Str("func", "NewService").Caller().Logger()

	s := &Service{
		c: c,
		r: r,
	}

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	orders, err := s.r.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	log.Info().Int("records count", len(orders)).Msg("records fetched from repository")

	s.c.LoadCache(orders)

	log.Info().Msg("records loaded to cache")

	return s, nil
}

func (s *Service) SaveOrder(ctx context.Context, order models.Order) error {

	log := log.Logger.With().Str("func", "Service.SaveOrder").Caller().Logger()

	log.Debug().Msg("validating order struct")

	if err := order.Validate(); err != nil {
		return ErrValidationErr{
			Msg: "cannot save invalid order",
			Err: err,
		}
	}

	log.Debug().Msg("saving order to repository")

	if err := s.r.InsertOrder(ctx, order); err != nil {
		return ErrRepositoryErr{
			Msg: "failed to save order",
			Err: err,
		}
	}

	log.Debug().Msg("saving order to cache")

	if err := s.c.Set(order); err != nil {
		return ErrCacheErr{
			Msg: "failed to cache order",
			Err: err,
		}
	}

	log.Debug().Msg("order saved")

	return nil
}

func (s *Service) GetOrder(ctx context.Context, OrderUID string) (models.Order, error) {
	order, err := s.c.Get(OrderUID)
	if err != nil {
		return models.Order{}, ErrOrderNotFound
	}

	return order, nil
}
