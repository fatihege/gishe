package catalogpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatihege/gishe/internal/catalog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateVenue(ctx context.Context, venue catalog.CreateVenueInput) (catalog.Venue, error) {
	query := `
		INSERT INTO venues (name, address, city)
		VALUES ($1, $2, $3)
		RETURNING id, name, address, city, created_at
	`

	var newVenue catalog.Venue

	err := r.db.QueryRow(
		ctx, query, venue.Name, venue.Address, venue.City,
	).Scan(
		&newVenue.ID, &newVenue.Name, &newVenue.Address, &newVenue.City, &newVenue.CreatedAt,
	)
	if err != nil {
		return catalog.Venue{}, fmt.Errorf("create new venue: %v", err)
	}

	return newVenue, nil
}

func (r *Repository) GetVenues(ctx context.Context) ([]catalog.Venue, error) {
	query := `
		SELECT * FROM venues
	`

	var venues []catalog.Venue

	rows, err := r.db.Query(ctx, query)
	if errors.Is(err, pgx.ErrNoRows) {
		return []catalog.Venue{}, catalog.ErrNoVenuesFound
	}
	if err != nil {
		return []catalog.Venue{}, fmt.Errorf("get venues: %v", err)
	}

	for i := 0; rows.Next(); i++ {
		venues = append(venues, catalog.Venue{})
		if err := rows.Scan(
			&venues[i].ID, &venues[i].Name, &venues[i].Address, &venues[i].City, &venues[i].CreatedAt,
		); err != nil {
			return []catalog.Venue{}, fmt.Errorf("iterate venues: %v", err)
		}
	}

	return venues, nil
}
