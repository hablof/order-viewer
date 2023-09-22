package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Payment struct {
	Transaction  string   `json:"transaction"`
	RequestID    string   `json:"request_id"`
	Currency     string   `json:"currency"`
	Provider     string   `json:"provider"`
	Amount       int      `json:"amount"`
	PaymentDT    UnixTime `json:"payment_dt"`
	Bank         string   `json:"bank"`
	DeliveryCost int      `json:"delivery_cost"`
	GoodsTotal   int      `json:"goods_total"`
	CustomFee    int      `json:"custom_fee"`
}

type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	u.Time = time.Unix(timestamp, 0)
	return nil
}

func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", (u.Time.Unix()))), nil
}
