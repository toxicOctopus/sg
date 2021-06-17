-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table channels
(
	id serial
		constraint channels_pk
			primary key,
	twitch_name varchar not null,
	action_cooldown int
);

create unique index channels_twitch_name_uindex
	on channels (twitch_name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists channels CASCADE;
-- +goose StatementEnd
