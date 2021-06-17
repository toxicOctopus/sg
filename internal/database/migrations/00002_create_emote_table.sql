-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table emotes
(
    id serial
        constraint emotes_pk
            primary key,
    name varchar,
    image_path varchar
);

create unique index emotes_name_uindex
    on emotes (name);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists emotes
-- +goose StatementEnd
