package models

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNoItems          = errors.New("order has no items")
	ErrWrongPaymentID   = errors.New("payment has wrong id")
	ErrWrongTrackNumber = errors.New("item has wrong track number")
	ErrWrongItemOrderID = errors.New("item has wrong order_id")
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

func (o *Order) UnmarshalJSON(b []byte) error {
	tempOrder := struct {
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
	}{}

	if err := json.Unmarshal(b, &tempOrder); err != nil {
		return err
	}

	o.OrderUID = tempOrder.OrderUID
	o.TrackNumber = tempOrder.TrackNumber
	o.Entry = tempOrder.Entry

	o.Delivery = tempOrder.Delivery
	o.Payment = tempOrder.Payment
	o.Items = tempOrder.Items

	o.Locale = tempOrder.Locale
	o.InternalSignature = tempOrder.InternalSignature
	o.CustomerID = tempOrder.CustomerID
	o.DeliveryService = tempOrder.DeliveryService
	o.ShardKey = tempOrder.ShardKey
	o.SMID = tempOrder.SMID
	o.DateCreated = tempOrder.DateCreated
	o.OofShard = tempOrder.OofShard

	for i := range o.Items {
		o.Items[i].OrderUID = o.OrderUID
	}

	return nil
}

// returns nil if order is valid
func (o Order) Validate() error {
	var err error
	if len(o.Items) == 0 {
		err = errors.Join(err, ErrNoItems)
	}

	if o.OrderUID != o.Payment.Transaction {
		err = errors.Join(err, ErrWrongPaymentID)
	}

	for _, item := range o.Items {

		// не уверен на счёт этой проверки
		if item.TrackNumber != o.TrackNumber {
			err = errors.Join(err, ErrWrongTrackNumber)
		}

		if item.OrderUID != o.OrderUID {
			err = errors.Join(err, ErrWrongItemOrderID)
		}
	}

	return err
}
