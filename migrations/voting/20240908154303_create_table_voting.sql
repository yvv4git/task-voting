-- +goose Up
-- +goose StatementBegin
CREATE TABLE voting
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP,
    started_at  TIMESTAMP,
    ended_at    TIMESTAMP,
    deleted_at  TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS voting;
-- +goose StatementEnd