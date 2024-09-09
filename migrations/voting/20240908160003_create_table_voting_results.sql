-- +goose Up
-- +goose StatementBegin
CREATE TABLE voting_results
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    voting_id    UUID NOT NULL REFERENCES voting(id) ON DELETE CASCADE,
    invariant_id UUID NOT NULL REFERENCES voting_invariance(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS voting_results;
-- +goose StatementEnd