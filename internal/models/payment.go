package models

import "time"

type Payment struct {
	Transaction   string    `json:"transaction"`
	RequestID     string    `json:"request_id"`
	Currency      string    `json:"currency"`
	Provider      string    `json:"provider"`
	Amount        int       `json:"amount"`
	PaymentDT     time.Time `json:"payment_dt"`
	Bank          string    `json:"bank"`
	Delivery_cost int       `json:"delivery_cost"`
	GoodsTotal    int       `json:"goods_total"`
	CustomFee     int       `json:"custom_fee"`
}
