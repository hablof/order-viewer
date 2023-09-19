package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hablof/order-viewer/internal/httpcontroller"
	"github.com/hablof/order-viewer/internal/models"
	"github.com/hablof/order-viewer/internal/templates"
)

type MockService struct{}

func (MockService) GetOrder(ctx context.Context, OrderUID string) (models.Order, error) {
	return models.Order{
		OrderUID:    OrderUID,
		TrackNumber: "CA123456789RU",
		Entry:       "Entry",
		Delivery: models.Delivery{
			Name:    "Костя",
			Phone:   "+7...",
			Zip:     "zip",
			City:    "Пенза",
			Address: "проспект Победы, 144",
			Region:  "Пензенская область",
			Email:   "spms@pnz.ru",
		},
		Payment: models.Payment{
			Transaction:  OrderUID,
			RequestID:    "",
			Currency:     "RUB",
			Provider:     "Google Pay",
			Amount:       300,
			PaymentDT:    time.Now(),
			Bank:         "WB-Bank",
			DeliveryCost: 100,
			GoodsTotal:   200,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				OrderUID:    OrderUID,
				ChrtID:      0,
				TrackNumber: "",
				Price:       200,
				RID:         "",
				Name:        "Просто яблоко",
				Sale:        25,
				Size:        "",
				TotalPrice:  150,
				NMID:        0,
				Brand:       "apple",
				Status:      0,
			},
			{
				OrderUID:    OrderUID,
				ChrtID:      0,
				TrackNumber: "",
				Price:       100,
				RID:         "",
				Name:        "cerf",
				Sale:        50,
				Size:        "",
				TotalPrice:  50,
				NMID:        0,
				Brand:       "samsung",
				Status:      0,
			},
			{
				OrderUID:    OrderUID,
				ChrtID:      0,
				TrackNumber: "",
				Price:       50,
				RID:         "",
				Name:        "padlo",
				Sale:        0,
				Size:        "",
				TotalPrice:  50,
				NMID:        0,
				Brand:       "dgfdrfdh",
				Status:      0,
			},
		},
		Locale:            "",
		InternalSignature: "",
		CustomerID:        "",
		DeliveryService:   "Деловые линии",
		ShardKey:          "",
		SMID:              0,
		DateCreated:       time.Now(),
		OofShard:          "",
	}, nil
}

func main() {
	t, err := templates.GetTemplates()
	if err != nil {
		log.Println(err)
		return
	}

	mux := httpcontroller.GetRouter(MockService{}, t)

	http.ListenAndServe(":8000", mux)
}
