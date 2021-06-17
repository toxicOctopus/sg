-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table channel_emotes
(
    id serial
        constraint channel_emotes_pk
            primary key,
    channel_id int not null,
    emote_id int not null,
    action_type_id int not null
);

create unique index channel_emotes_channel_id_emote_uindex
    on channel_emotes (channel_id, emote_id);

alter table channel_emotes
    add constraint channel_emotes_channels_id_fk
        foreign key (channel_id) references channels
            on update cascade on delete cascade;

alter table channel_emotes
    add constraint channel_emotes_emotes_id_fk
        foreign key (emote_id) references emotes
            on update cascade on delete cascade;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop index if exists channel_emotes_channel_id_emote_uindex;
drop table if exists channel_emotes;
-- +goose StatementEnd
