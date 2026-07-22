package catalog

import (
	"context"
)

type Repository interface {
	CreateVenue(ctx context.Context, venue CreateVenueInput) (Venue, error)

	GetVenues(ctx context.Context) ([]Venue, error)

	// CreateSession(ctx context.Context, session Session) (Session, error)

	// FindSessionById(ctx context.Context, id string) (Session, error)

	// FindSessionsByVenueId(ctx context.Context, venueID string) ([]Session, error)

	// CreateSeats(ctx context.Context, seats []Seat) ([]Seat, error)

	// FindSeatsBySessionId(ctx context.Context, sessionID string) ([]Seat, error)
}
