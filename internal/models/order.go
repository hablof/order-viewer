package models

import (
	"errors"
	"time"
)

var (
	ErrNoItems          = errors.New("order has no items")
	ErrWrongPaymentID   = errors.New("payment has wrong id")
	ErrWrongTrackNumber = errors.New("item has wrong track number")
)

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`

	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`

	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SMID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (o Order) IsValid() error {
	var err error
	if len(o.Items) == 0 {
		err = errors.Join(err, ErrNoItems)
	}

	if o.OrderUID != o.Payment.Transaction {
		err = errors.Join(err, ErrWrongPaymentID)
	}

	for _, item := range o.Items {
		if item.TrackNumber != o.TrackNumber {
			err = errors.Join(err, ErrWrongTrackNumber)
		}
	}

	return err
}
