-- +goose Up
-- +goose StatementBegin

CREATE TABLE order (
    order_uid    VARCHAR(32) PRIMARY KEY,     -- FOREIGN KEY to payment payment_id = order_uid id заказа
    track_number VARCHAR(32) NOT NULL UNIQUE, -- FOREIGN KEY to item??
    delivery_id  INTEGER     NOT NULL,        -- FOREIGN KEY to delivery

    entry              VARCHAR(32)  NOT NULL, -- ??
    locale             VARCHAR(32)  NOT NULL,
    internal_signature VARCHAR(256) NOT NULL DEFAULT '',
    customer_id        VARCHAR(32)  NOT NULL,
    delivery_service   VARCHAR(32)  NOT NULL, -- Вынести в отдельную таблицу?
    shardkey           VARCHAR(16)  NOT NULL,
    sm_id              INTEGER      NOT NULL, -- непонятно что такое
    date_created       TIMESTAMP    NOT NULL,
    oof_shard          VARCHAR(16)  NOT NULL  -- непонятно что такое
);

CREATE TABLE delivery (
    delivery_id SERIAL       PRIMARY KEY,
    name        VARCHAR(32)  NOT NULL,
    phone       VARCHAR(18)  NOT NULL,
    zip         VARCHAR(12)  NOT NULL,
    city        VARCHAR(32)  NOT NULL, -- Вынести в отдельную таблицу?
    address     VARCHAR(64)  NOT NULL,
    email       VARCHAR(256) NOT NULL
);

CREATE TABLE payment (
    payment_id    VARCHAR(32)  PRIMARY KEY,
    request_id    VARCHAR(256) NOT NULL DEFAULT '',
    currency      VARCHAR(16)  NOT NULL, -- Вынести в отдельную таблицу?
    provider      VARCHAR(32)  NOT NULL, -- Вынести в отдельную таблицу?
    amount        INTEGER      NOT NULL, -- констреинт на сумму?
    payment_dt    TIMESTAMP    NOT NULL, 
    bank          VARCHAR(256) NOT NULL, -- Вынести в отдельную таблицу?
    delivery_cost INTEGER      NOT NULL, -- констреинт на сумму?
    goods_total   INTEGER      NOT NULL, -- констреинт на сумму?
    custom_fee    INTEGER      NOT NULL  -- default 0 ? констреинт на сумму?
);

CREATE TABLE item (
    rid          VARCHAR(32) PRIMARY KEY, -- уникальный на доставку+артикул
    track_number VARCHAR(32) NOT NULL,    -- not unique: по одному трекномеру несколько товаров

    nm_id       INTEGER      NOT NULL, -- видимо, артикул
    chrt_id     INTEGER      NOT NULL, -- непонятно что такое
    price       INTEGER      NOT NULL,
    name        VARCHAR(256) NOT NULL,
    sale        INTEGER      NOT NULL DEFAULT 0,
    size        VARCHAR(32)  NOT NULL,
    total_price INTEGER      NOT NULL,
    brand       VARCHAR(256) NOT NULL,
    status      INTEGER      NOT NULL  -- непонятно что за числа
);

-- Удобно доставать товары, относящиеся к искомому заказу, исползуя track_number как внешний ключ
CREATE INDEX track_number_idx ON item (track_number);

-- rid уникальный на доставку+артикул, нет связи "многие ко многим"
-- CREATE TABLE order_item (
--     order_uid VARCHAR(32) FOREIGN KEY,
--     item_rid  VARCHAR(32) FOREIGN KEY,
-- );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS track_number_idx;

DROP TABLE IF EXISTS item;

DROP TABLE IF EXISTS payment;

DROP TABLE IF EXISTS delivery;

DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
