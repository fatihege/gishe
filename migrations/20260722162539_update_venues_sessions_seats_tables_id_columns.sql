-- +goose Up
ALTER TABLE venues
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE sessions
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE seats
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE venues
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE sessions
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE seats
    ALTER COLUMN id DROP DEFAULT;
