-- +goose Up
create table users(
    id uuid primary key,
    email text not null unique,
    hashed_password text default 'unset' not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    is_admin boolean default false not null
);

-- +goose Down
drop table users;
