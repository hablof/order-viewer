package service

import (
	"context"
	"errors"
	"time"

	"github.com/hablof/order-viewer/internal/models"
)

// cache / repository errors
var (
	ErrNotFound             = errors.New("not found")
	ErrDuplicatesNotAllowed = errors.New("not allowed to store duplicates")
)

// // service errors
// var ()

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

	s.c.LoadCache(orders)

	return s, nil
}

func (s *Service) SaveOrder(ctx context.Context, order models.Order) error {

	if err := order.Validate(); err != nil {
		return ErrValidationErr{
			Msg: "cannot save invalid order",
			Err: err,
		}
	}

	if err := s.r.InsertOrder(ctx, order); err != nil {
		return ErrRepositoryErr{
			Msg: "failed to save order",
			Err: err,
		}
	}

	if err := s.c.Set(order); err != nil {
		return ErrCacheErr{
			Msg: "failed to cache order",
			Err: err,
		}
	}

	return nil
}

func (s *Service) GetOrder(ctx context.Context, OrderUID string) (models.Order, error) {
	order, err := s.c.Get(OrderUID)
	if err != nil {
		return models.Order{},
			ErrCacheErr{
				Msg: "failed to get order from cache",
				Err: err,
			}
	}

	return order, nil
}
