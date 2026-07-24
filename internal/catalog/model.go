package catalog

import (
	"time"

	"github.com/google/uuid"
)

type Venue struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID            uuid.UUID     `json:"id"`
	VenueID       uuid.UUID     `json:"venue_id"`
	Title         string        `json:"title"`
	StartsAt      time.Time     `json:"starts_at"`
	SalesOpensAt  time.Time     `json:"sales_opens_at"`
	SalesClosesAt *time.Time    `json:"sales_closes_at"`
	Status        SessionStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Seat struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	RowLabel   string     `json:"row_label"`
	SeatNumber int32      `json:"seat_number"`
	PriceCents int32      `json:"price_cents"`
	Status     SeatStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
}
