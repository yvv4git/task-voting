-- +goose Up
-- +goose StatementBegin
CREATE TABLE voting_invariance
(
    id          SERIAL PRIMARY KEY,
    voting_id   INT NOT NULL REFERENCES voting(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS voting_invariance;
-- +goose StatementEnd