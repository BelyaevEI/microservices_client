-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
INSERT INTO roles (name) VALUES ('user') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('admin') ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name IN ('user', 'admin');
drop table roles;
-- +goose StatementEnd