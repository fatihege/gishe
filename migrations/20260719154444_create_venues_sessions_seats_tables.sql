-- +goose Up
CREATE TABLE venues (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    venue_id UUID NOT NULL REFERENCES venues(id) ON DELETE RESTRICT,
    title TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    sales_opens_at TIMESTAMPTZ NOT NULL,
    sales_closes_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'cancelled', 'completed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT sessions_sale_before_event CHECK (sales_opens_at < starts_at)
);

CREATE INDEX idx_sessions_venue_id ON sessions(venue_id);
CREATE INDEX idx_sessions_status_sale_opens ON sessions(status, sales_opens_at);

CREATE TABLE seats (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE RESTRICT,
    row_label TEXT NOT NULL,
    seat_number INT NOT NULL CHECK (seat_number > 0),
    price_cents INT NOT NULL CHECK (price_cents > 0),
    status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'held', 'sold')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (session_id, row_label, seat_number)
);

CREATE INDEX idx_seats_session_id ON seats(session_id);
CREATE INDEX idx_seats_session_available ON seats(session_id) WHERE status = 'available';

-- +goose Down
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS venues;
