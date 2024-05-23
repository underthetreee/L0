-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders (
    uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(255),
    internal_signature VARCHAR(255),
    custom_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shard_key VARCHAR(255),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_uid VARCHAR(255),
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255),

    CONSTRAINT order_uid_fk FOREIGN KEY (order_uid) REFERENCES orders(uid)
);

CREATE TABLE IF NOT EXISTS payments (
    order_uid VARCHAR(255),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INT,
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,

    CONSTRAINT order_uid_fk FOREIGN KEY (order_uid) REFERENCES orders(uid)
);

CREATE TABLE IF NOT EXISTS items (
    order_uid VARCHAR(255),
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(255),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT,

    CONSTRAINT order_uid_fk FOREIGN KEY (order_uid) REFERENCES orders(uid)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd
