package repository

import (
	"time"

	"github.com/hablof/order-viewer/internal/models"
)

var (
	testObjs = []models.Order{
		{
			OrderUID:    "first_uid",
			TrackNumber: "first_track",
			Entry:       "first_entry",
			Delivery: models.Delivery{
				Name:    "first_delivery",
				Phone:   "+7...",
				Zip:     "first_zip",
				City:    "first_city",
				Address: "first_addres",
				Region:  "first_region",
				Email:   "first_email",
			},
			Payment: models.Payment{
				Transaction:  "first_uid",
				RequestID:    "first_request",
				Currency:     "first_currency",
				Provider:     "first_provider",
				Amount:       111,
				PaymentDT:    time.Date(2023, time.September, 16, 23, 9, 15, 0, time.UTC),
				Bank:         "first_bank",
				DeliveryCost: 100,
				GoodsTotal:   10,
				CustomFee:    1,
			},
			Items: []models.Item{
				{
					OrderUID:    "first_uid",
					ChrtID:      1,
					TrackNumber: "first_track",
					Price:       2,
					RID:         "first_first_rid",
					Name:        "first_first_name",
					Sale:        50,
					Size:        "first_first_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "first_first_brand",
					Status:      1,
				},
				{
					OrderUID:    "first_uid",
					ChrtID:      1,
					TrackNumber: "first_track",
					Price:       2,
					RID:         "first_second_rid",
					Name:        "first_second_name",
					Sale:        50,
					Size:        "first_second_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "first_second_brand",
					Status:      1,
				},
				{
					OrderUID:    "first_uid",
					ChrtID:      1,
					TrackNumber: "first_track",
					Price:       8,
					RID:         "first_third_rid",
					Name:        "first_third_name",
					Sale:        0,
					Size:        "first_third_size",
					TotalPrice:  8,
					NMID:        1,
					Brand:       "first_third_brand",
					Status:      1,
				},
			},
			Locale:            "first_locale",
			InternalSignature: "first_signature",
			CustomerID:        "first_id",
			DeliveryService:   "first_delivery_service",
			ShardKey:          "first_shard",
			SMID:              1,
			DateCreated:       time.Date(2023, time.September, 16, 22, 9, 15, 0, time.UTC),
			OofShard:          "first_oof_shard",
		},
		{
			OrderUID:    "second_uid",
			TrackNumber: "second_track",
			Entry:       "second_entry",
			Delivery: models.Delivery{
				Name:    "second_delivery",
				Phone:   "+7...",
				Zip:     "second_zip",
				City:    "second_city",
				Address: "second_addres",
				Region:  "second_region",
				Email:   "second_email",
			},
			Payment: models.Payment{
				Transaction:  "second_uid",
				RequestID:    "second_request",
				Currency:     "second_currency",
				Provider:     "second_provider",
				Amount:       111,
				PaymentDT:    time.Date(2023, time.September, 16, 23, 9, 15, 0, time.UTC),
				Bank:         "second_bank",
				DeliveryCost: 100,
				GoodsTotal:   10,
				CustomFee:    1,
			},
			Items: []models.Item{
				{
					OrderUID:    "second_uid",
					ChrtID:      1,
					TrackNumber: "second_track",
					Price:       2,
					RID:         "second_first_rid",
					Name:        "second_first_name",
					Sale:        50,
					Size:        "second_first_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "second_first_brand",
					Status:      1,
				},
				{
					OrderUID:    "second_uid",
					ChrtID:      1,
					TrackNumber: "second_track",
					Price:       2,
					RID:         "second_second_rid",
					Name:        "second_second_name",
					Sale:        50,
					Size:        "second_second_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "second_second_brand",
					Status:      1,
				},
				{
					OrderUID:    "second_uid",
					ChrtID:      1,
					TrackNumber: "second_track",
					Price:       8,
					RID:         "second_third_rid",
					Name:        "second_third_name",
					Sale:        0,
					Size:        "second_third_size",
					TotalPrice:  8,
					NMID:        1,
					Brand:       "second_third_brand",
					Status:      1,
				},
			},
			Locale:            "second_locale",
			InternalSignature: "second_signature",
			CustomerID:        "second_id",
			DeliveryService:   "second_delivery_service",
			ShardKey:          "second_shard",
			SMID:              1,
			DateCreated:       time.Date(2023, time.September, 16, 22, 9, 15, 0, time.UTC),
			OofShard:          "second_oof_shard",
		},
		{
			OrderUID:    "third_uid",
			TrackNumber: "third_track",
			Entry:       "third_entry",
			Delivery: models.Delivery{
				Name:    "third_delivery",
				Phone:   "+7...",
				Zip:     "third_zip",
				City:    "third_city",
				Address: "third_addres",
				Region:  "third_region",
				Email:   "third_email",
			},
			Payment: models.Payment{
				Transaction:  "third_uid",
				RequestID:    "third_request",
				Currency:     "third_currency",
				Provider:     "third_provider",
				Amount:       111,
				PaymentDT:    time.Date(2023, time.September, 16, 23, 9, 15, 0, time.UTC),
				Bank:         "third_bank",
				DeliveryCost: 100,
				GoodsTotal:   10,
				CustomFee:    1,
			},
			Items: []models.Item{
				{
					OrderUID:    "third_uid",
					ChrtID:      1,
					TrackNumber: "third_track",
					Price:       2,
					RID:         "third_first_rid",
					Name:        "third_first_name",
					Sale:        50,
					Size:        "third_first_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "third_first_brand",
					Status:      1,
				},
				{
					OrderUID:    "third_uid",
					ChrtID:      1,
					TrackNumber: "third_track",
					Price:       2,
					RID:         "third_second_rid",
					Name:        "third_second_name",
					Sale:        50,
					Size:        "third_second_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "third_second_brand",
					Status:      1,
				},
				{
					OrderUID:    "third_uid",
					ChrtID:      1,
					TrackNumber: "third_track",
					Price:       8,
					RID:         "third_third_rid",
					Name:        "third_third_name",
					Sale:        0,
					Size:        "third_third_size",
					TotalPrice:  8,
					NMID:        1,
					Brand:       "third_third_brand",
					Status:      1,
				},
			},
			Locale:            "third_locale",
			InternalSignature: "third_signature",
			CustomerID:        "third_id",
			DeliveryService:   "third_delivery_service",
			ShardKey:          "third_shard",
			SMID:              1,
			DateCreated:       time.Date(2023, time.September, 16, 22, 9, 15, 0, time.UTC),
			OofShard:          "third_oof_shard",
		},
		{
			OrderUID:    "fourth_uid",
			TrackNumber: "fourth_track",
			Entry:       "fourth_entry",
			Delivery: models.Delivery{
				Name:    "fourth_delivery",
				Phone:   "+7...",
				Zip:     "fourth_zip",
				City:    "fourth_city",
				Address: "fourth_addres",
				Region:  "fourth_region",
				Email:   "fourth_email",
			},
			Payment: models.Payment{
				Transaction:  "fourth_uid",
				RequestID:    "fourth_request",
				Currency:     "fourth_currency",
				Provider:     "fourth_provider",
				Amount:       111,
				PaymentDT:    time.Date(2023, time.September, 16, 23, 9, 15, 0, time.UTC),
				Bank:         "fourth_bank",
				DeliveryCost: 100,
				GoodsTotal:   10,
				CustomFee:    1,
			},
			Items: []models.Item{
				{
					OrderUID:    "fourth_uid",
					ChrtID:      1,
					TrackNumber: "fourth_track",
					Price:       2,
					RID:         "fourth_first_rid",
					Name:        "fourth_first_name",
					Sale:        50,
					Size:        "fourth_first_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "fourth_first_brand",
					Status:      1,
				},
				{
					OrderUID:    "fourth_uid",
					ChrtID:      1,
					TrackNumber: "fourth_track",
					Price:       2,
					RID:         "fourth_second_rid",
					Name:        "fourth_second_name",
					Sale:        50,
					Size:        "fourth_second_size",
					TotalPrice:  1,
					NMID:        1,
					Brand:       "fourth_second_brand",
					Status:      1,
				},
				{
					OrderUID:    "fourth_uid",
					ChrtID:      1,
					TrackNumber: "fourth_track",
					Price:       8,
					RID:         "fourth_third_rid",
					Name:        "fourth_third_name",
					Sale:        0,
					Size:        "fourth_third_size",
					TotalPrice:  8,
					NMID:        1,
					Brand:       "fourth_third_brand",
					Status:      1,
				},
			},
			Locale:            "fourth_locale",
			InternalSignature: "fourth_signature",
			CustomerID:        "fourth_id",
			DeliveryService:   "fourth_delivery_service",
			ShardKey:          "fourth_shard",
			SMID:              1,
			DateCreated:       time.Date(2023, time.September, 16, 22, 9, 15, 0, time.UTC),
			OofShard:          "fourth_oof_shard",
		},
	}
)
