-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM users
        GROUP BY lower(email)
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION
            'cannot normalize users.email: duplicate emails exist ignoring case';
    END IF;
END
$$;
-- +goose StatementEnd

ALTER TABLE users
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING created_at AT TIME ZONE 'UTC';

ALTER TABLE users
    DROP CONSTRAINT users_email_key;

UPDATE users
SET email = lower(email);

ALTER TABLE users
    ADD CONSTRAINT users_email_lowercase_check
        CHECK (email = lower(email)),
    ADD CONSTRAINT users_email_unique
        UNIQUE (email);

-- +goose Down
ALTER TABLE users
    DROP CONSTRAINT users_email_unique,
    DROP CONSTRAINT users_email_lowercase_check;

ALTER TABLE users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE users
    ALTER COLUMN created_at TYPE TIMESTAMP
    USING created_at AT TIME ZONE 'UTC';

ALTER TABLE users
    ALTER COLUMN id DROP DEFAULT;
