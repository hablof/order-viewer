-- +goose Up
-- +goose StatementBegin

CREATE TABLE order (
    order_uid    VARCHAR(32) PRIMARY KEY, -- payment_id = order_uid
    track_number VARCHAR(32) NOT NULL UNIQUE, 
    delivery_id  INTEGER     NOT NULL, 
    "entry": "WBIL",
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test",
    "delivery_service": "meest",
    "shardkey": "9",
    "sm_id": 99,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
);

CREATE TABLE delivery (
    delivery_id SERIAL       PRIMARY KEY,
    name        VARCHAR(32)  NOT NULL,
    phone       VARCHAR(18)  NOT NULL,
    zip         VARCHAR(12)  NOT NULL,
    city        VARCHAR(32)  NOT NULL,
    address     VARCHAR(64)  NOT NULL,
    email       VARCHAR(256) NOT NULL
);

CREATE TABLE payment (
    payment_id    VARCHAR(32)  PRIMARY KEY,
    request_id    VARCHAR(256) NOT NULL DEFAULT '',
    currency      VARCHAR(16)  NOT NULL, -- Вынести в отдельную таблицу?
    provider      VARCHAR(32)  NOT NULL, -- Вынести в отдельную таблицу?
    amount        INTEGER      NOT NULL, -- констреинт на сумму?
    "payment_dt": 1637907727,
    bank          VARCHAR(256) NOT NULL, -- Вынести в отдельную таблицу?
    delivery_cost INTEGER      NOT NULL, -- констреинт на сумму?
    goods_total:  INTEGER      NOT NULL, -- констреинт на сумму?
    custom_fee    INTEGER      NOT NULL  -- default 0 ?
);

CREATE TABLE item (
    rid     VARCHAR(32) PRIMARY KEY,
    nm_id   INTEGER NOT NULL,
    chrt_id INTEGER NOT NULL,
    track_number": "WBILMTESTTRACK",
    price": 453,
    name": "Mascaras",
    sale": 30,
    size": "0",
    total_price": 317,
    brand": "Vivienne Sabo",
    status": 202
);

CREATE TABLE order_item (
    order_uid VARCHAR(19) FOREIGN KEY,
    item_rid  VARCHAR(19) FOREIGN KEY,
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
