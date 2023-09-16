
CREATE TABLE orders (
    order_uid    VARCHAR(32) PRIMARY KEY,     -- FOREIGN KEY to payment transaction = order_uid id заказа, FOREIGN KEY to item
    track_number VARCHAR(32) NOT NULL UNIQUE, -- FOREIGN KEY to item??
    delivery_id  INTEGER     NOT NULL UNIQUE, -- FOREIGN KEY to delivery

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
    region      VARCHAR(64)  NOT NULL
    email       VARCHAR(256) NOT NULL
);

CREATE TABLE payment (
    transaction   VARCHAR(32)  PRIMARY KEY, -- transaction = order_uid id заказа
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

    -- поскольку в таблице orders и track_number и order_uid обладают свойством UNIQUE, в качестве внешного ключа можно импользовать и то и другое.
    -- Будем использовать именно order_uid.
    order_uid    VARCHAR(32)  NOT NULL, 
    chrt_id      INTEGER      NOT NULL, -- непонятно что такое
    track_number VARCHAR(32)  NOT NULL, -- not unique: по одному трекномеру несколько товаров
    price        INTEGER      NOT NULL,
    rid          VARCHAR(32)  NOT NULL, -- возможно уникальный на доставку+артикул, но непонятно, как хранятся заказы с несколькими одинаковами товарами
    name         VARCHAR(256) NOT NULL,
    sale         INTEGER      NOT NULL DEFAULT 0,
    size         VARCHAR(32)  NOT NULL,
    total_price  INTEGER      NOT NULL,
    nm_id        INTEGER      NOT NULL, -- видимо, артикул
    brand        VARCHAR(256) NOT NULL,
    status       INTEGER      NOT NULL  -- непонятно что за числа
);

-- Удобно доставать товары, относящиеся к искомому заказу, исползуя track_number как внешний ключ
-- CREATE INDEX track_number_idx ON item (track_number);

-- А ещё удобнее доставать товары, относящиеся к искомому заказу, исползуя ID заказа как внешний ключ
CREATE INDEX order_uid_idx ON item (order_uid);

-- rid уникальный на доставку+артикул, нет связи "многие ко многим"
-- CREATE TABLE order_item (
--     order_uid VARCHAR(32) FOREIGN KEY,
--     item_rid  VARCHAR(32) FOREIGN KEY,
-- );
