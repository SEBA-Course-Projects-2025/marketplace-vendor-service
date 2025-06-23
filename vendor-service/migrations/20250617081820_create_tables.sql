-- +goose Up
-- +goose StatementBegin
CREATE TABLE vendor_accounts
(
    id            uuid primary key,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    name          varchar(255) not null,
    description   text         not null,
    logo          varchar(255) not null,
    address       varchar(255) not null,
    website       varchar(255) not null,
    created_at    timestamp default now(),
    updated_at    timestamp default now()
);

CREATE TABLE reviews
(
    id            uuid primary key,
    vendor_id     uuid references vendor_accounts (id) on delete cascade,
    product_id    uuid         not null,
    reviewer_id   uuid         not null,
    reviewer_name varchar(255) not null,
    rating        float8       not null,
    comment       text         not null,
    created_at    timestamp default now(),
    updated_at    timestamp default now()
);

CREATE TABLE replies
(
    id           uuid primary key,
    review_id    uuid references reviews (id) on delete cascade,
    replier_id   uuid         not null,
    replier_name varchar(255) not null,
    comment      text         not null,
    created_at   timestamp default now(),
    updated_at   timestamp default now()
);

CREATE TABLE orders
(
    id          uuid primary key,
    vendor_id   uuid references vendor_accounts (id) on delete cascade,
    customer_id uuid           not null,
    total_price numeric(12, 2) not null,
    status      varchar(40)    not null,
    created_at  timestamp default now(),
    updated_at  timestamp default now()
);

CREATE TABLE order_items
(
    id           uuid primary key not null,
    product_id   uuid             not null,
    order_id     uuid references orders (id),
    product_name varchar(255)     not null,
    quantity     int              not null,
    unit_price   numeric(12, 2)   not null,
    image_url    text             not null,
    created_at   timestamp default now(),
    updated_at   timestamp default now()
);

CREATE INDEX idx_orders_created_at ON orders (created_at);
CREATE INDEX idx_orders_total_price ON orders (total_price);

CREATE INDEX idx_reviews_created_at ON reviews (created_at);
CREATE INDEX idx_reviews_rating ON reviews (rating);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_total_price;
DROP INDEX IF EXISTS idx_reviews_created_at;
DROP INDEX IF EXISTS idx_reviews_rating;

DROP TABLE IF EXISTS replies;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS vendor_accounts;
-- +goose StatementEnd
