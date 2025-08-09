-- +goose Up
-- +goose StatementBegin
create table outboxes
(
    id uuid primary key not null,
    exchange varchar(255) not null,
    event_type varchar(255) not null,
    payload jsonb not null,
    created_at timestamp default now(),
    processed bool not null default false,
    processed_at timestamp
);

create table processed_messages
(
    message_id uuid primary key not null
);

alter table order_items add column stock_id uuid;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table order_items drop column stock_id;

DROP TABLE IF EXISTS outboxes;
DROP TABLE IF EXISTS processed_messages;
-- +goose StatementEnd
