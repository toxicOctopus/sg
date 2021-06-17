-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO action_types (name, action_source) VALUES ('damage', 1),('dodge', 1),('block', 1),('overpower', 1);
INSERT INTO action_types (name, action_source) VALUES ('sweep', 0),('block', 0),('target_execute', 0),('fake', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
delete from action_types;
-- +goose StatementEnd
