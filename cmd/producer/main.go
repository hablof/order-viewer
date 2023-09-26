// Вспомогательная программа-producer
// Эта программа отправит в канал NATS-Streaming несколько json-посылок
// Тем самым позволит продемонстрировать работу основного сервиса
//

package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	invalidJSON = `{
	"name": "Test Testov",
	"zip": "2639809",{}
	"city": "Kiryat Mozkin",
	"address": "Ploshad Mira 15"
	"region": "Kraiot",
	"email": "smps@pnz.com"
}`

	invalidOrder_WrongPaymentID = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "invalid",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": [
		{
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`

	invalidOrder_NoItems = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": null,
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`

	invalidOrder_WrongTrackNumber = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": [
		{
			"chrt_id": 9934930,
			"track_number": "WrongTrackNumber",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`

	validOrder = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	},
	"payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	},
	"items": [
		{
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
}`
	sameIDviolatesUniqueConstraint = `{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "track",
		"entry": "WBIL",
		"delivery": {
			"name": "Bob",
			"phone": "+9725550000",
			"zip": "2639809",
			"city": "Kiryat Mozkin",
			"address": "Ploshad Mira 16",
			"region": "Kraiot",
			"email": "bob@gmail.com"
		},
		"payment": {
			"transaction": "b563feb7b2b84b6test",
			"request_id": "",
			"currency": "RUB",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "track",
				"price": 453,
				"rid": "ab4219087a764ae0btest",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "Vivienne Sabo",
				"status": 222
			}
		],
		"locale": "ru",
		"internal_signature": "",
		"customer_id": "test2",
		"delivery_service": "meest2",
		"shardkey": "92",
		"sm_id": 992,
		"date_created": "2021-10-26T06:22:19Z",
		"oof_shard": "2"
	}`

	orderWith15Items = `{
		"order_uid": "15_items",
		"track_number": "CA123456780100RU",
		"entry": "WBIL",
		"delivery": {
			"name": "Константин",
			"phone": "+7(800)555-35-35",
			"zip": "440062",
			"city": "Пенза",
			"address": "Проспект Победы, 73",
			"region": "Пензенская облать",
			"email": "smps@pnz.ru"
		},
		"payment": {
			"transaction": "15_items",
			"request_id": "",
			"currency": "RUB",
			"provider": "wbpay",
			"amount": 2000,
			"payment_dt": 1695300556,
			"bank": "the Big Green",
			"delivery_cost": 500,
			"goods_total": 1500,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 1,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test1",
				"name": "Чудо-ложка",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Magazin na diwane",
				"status": 200
			},
			{
				"chrt_id": 2,
				"track_number": "CA123456780100RU",
				"price": 200,
				"rid": "test1",
				"name": "Тряпка (микрофибра)",
				"sale": 50,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Magazin na diwane",
				"status": 200
			},
			{
				"chrt_id": 3,
				"track_number": "CA123456780100RU",
				"price": 150,
				"rid": "test3",
				"name": "The Real Edge of Glory",
				"sale": 33,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Magazin na diwane",
				"status": 200
			},
			{
				"chrt_id": 4,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test1",
				"name": "Гейнер Mr.Big",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "ООО Спорт-Пит",
				"status": 200
			},
			{
				"chrt_id": 5,
				"track_number": "CA123456780100RU",
				"price": 10000,
				"rid": "test5",
				"name": "Iphone 3G",
				"sale": 99,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Apple",
				"status": 200
			},
			{
				"chrt_id": 6,
				"track_number": "CA123456780100RU",
				"price": 125,
				"rid": "test6",
				"name": "ПКЦТ.685282.001 Веб-камера",
				"sale": 20,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "НИИ ПКБ ЛЦТ",
				"status": 200
			},
			{
				"chrt_id": 7,
				"track_number": "CA123456780100RU",
				"price": 800,
				"rid": "test7",
				"name": "Средство гель чистка для мытья посуды антижир эко жидкость",
				"sale": 87,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Китай",
				"status": 200
			},
			{
				"chrt_id": 8,
				"track_number": "CA123456780100RU",
				"price": 101,
				"rid": "test8",
				"name": "Швабра с распылителем",
				"sale": 1,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Magazin na diwane",
				"status": 200
			},
			{
				"chrt_id": 9,
				"track_number": "CA123456780100RU",
				"price": 150,
				"rid": "test9",
				"name": "Пастила натуральная ассорти 2 кг",
				"sale": 33,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Сладости Западной Сибири",
				"status": 200
			},
			{
				"chrt_id": 10,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test10",
				"name": "Сухой экстракт Родиола розовая, капсулы 60шт",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "ВсемЧай",
				"status": 200
			},
			{
				"chrt_id": 11,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test11",
				"name": "Опрыскиватель садовый 2л для растений",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "Садовик",
				"status": 200
			},
			{
				"chrt_id": 12,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test12",
				"name": "Стол книжка раскладной большой",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "OMI",
				"status": 200
			},
			{
				"chrt_id": 13,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test13",
				"name": "солнцезащитные очки",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "patolli",
				"status": 200
			},
			{
				"chrt_id": 14,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test14",
				"name": "Набор магнитных отверток",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "STYLE CITY",
				"status": 200
			},
			{
				"chrt_id": 15,
				"track_number": "CA123456780100RU",
				"price": 100,
				"rid": "test15",
				"name": "Сварочный аппарат инверторный САИ 190К",
				"sale": 0,
				"size": "0",
				"total_price": 100,
				"nm_id": 100,
				"brand": "РЕСАНТА",
				"status": 200
			}
		],
		"locale": "ru",
		"internal_signature": "yes",
		"customer_id": "yoHNBHifjE",
		"delivery_service": "Деловые конверты",
		"shardkey": "92",
		"sm_id": 1,
		"date_created": "2023-09-14T17:43:01Z",
		"oof_shard": "1"
	}`
)

var (
	Msgs = []string{
		invalidJSON,                    // error
		invalidOrder_WrongPaymentID,    // error
		invalidOrder_NoItems,           // error
		invalidOrder_WrongTrackNumber,  // error
		validOrder,                     // success
		sameIDviolatesUniqueConstraint, // error
		orderWith15Items,               // success
	} // Total: 5 errs, 2 ok
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	clusterID := "my_cluster"
	clientID := "test_pub"
	URL := "localhost:4222"

	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}

	// Connect to NATS
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect")
	}
	defer sc.Close()

	for _, msg := range Msgs {

		err = sc.Publish("orders", []byte(msg))
		if err != nil {
			log.Fatal().Err(err).Msg("Error during publish")
		}
		log.Info().Msg("Published msg")
	}

}
