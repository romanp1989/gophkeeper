create table if not exists "users"
(
    id serial primary key,
    login varchar(100) not null,
    password varchar(100) not null
);

create unique index if not exists users_login_udx
    on "users" (login);

create type secret_type as enum ('credential', 'text', 'blob', 'card');

create table if not exists "secrets"
(
    id serial primary key,
    user_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    title varchar(255) not null,
    metadata jsonb,
    payload bytea not null,
    type secret_type not null,
    text varchar(2000)
);