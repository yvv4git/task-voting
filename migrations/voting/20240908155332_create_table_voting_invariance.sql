-- +goose Up
-- +goose StatementBegin
CREATE TABLE voting_invariance
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voting_id   UUID NOT NULL REFERENCES voting(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS voting_invariance;
-- +goose StatementEnd