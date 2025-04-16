-- +goose Up
create table auth (
                        uuid text not null,
                         guid text primary key,
                         refresh_token text not null
);

-- +goose Down
drop table auth;