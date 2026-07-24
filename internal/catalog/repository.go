package catalog

import (
	"context"
)

type Repository interface {
	CreateVenue(ctx context.Context, venue CreateVenueInput) (Venue, error)

	FindVenueByID(ctx context.Context, id string) (Venue, error)

	GetVenues(ctx context.Context) ([]Venue, error)

	CreateSession(ctx context.Context, session CreateSessionInput) (Session, error)

	FindSessionByID(ctx context.Context, id string) (Session, error)

	FindSessionsByVenueID(ctx context.Context, venueID string) ([]Session, error)

	GetSessions(ctx context.Context) ([]Session, error)

	// CreateSeats(ctx context.Context, seats []Seat) ([]Seat, error)

	// FindSeatsBySessionId(ctx context.Context, sessionID string) ([]Seat, error)
}
