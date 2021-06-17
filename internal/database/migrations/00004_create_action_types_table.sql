-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table action_types
(
    id serial
        constraint action_types_pk
            primary key,
    name varchar,
    action_source int -- host(0) or viewer(1)
    );

alter table channel_emotes
    add constraint channel_emotes_action_types_id_fk
        foreign key (action_type_id) references action_types
            on update restrict on delete restrict;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table channel_emotes
    drop constraint if exists channel_emotes_action_types_id_fk;

drop table if exists action_types;
-- +goose StatementEnd
