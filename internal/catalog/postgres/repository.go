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

func (r *Repository) CreateVenue(ctx context.Context, input catalog.CreateVenueInput) (catalog.Venue, error) {
	query := `
		INSERT INTO venues (name, address, city)
		VALUES ($1, $2, $3)
		RETURNING id, name, address, city, created_at
	`

	var venue catalog.Venue

	err := r.db.QueryRow(
		ctx,
		query,
		input.Name,
		input.Address,
		input.City,
	).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Address,
		&venue.City,
		&venue.CreatedAt,
	)
	if err != nil {
		return catalog.Venue{}, fmt.Errorf("create new venue: %w", err)
	}

	return venue, nil
}

func (r *Repository) FindVenueByID(ctx context.Context, id string) (catalog.Venue, error) {
	query := `
		SELECT id, name, address, city, created_at
		FROM venues
		WHERE id = $1
	`

	var venue catalog.Venue

	err := r.db.QueryRow(
		ctx, query, id,
	).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Address,
		&venue.City,
		&venue.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return catalog.Venue{}, catalog.ErrVenueNotFound
	}
	if err != nil {
		return catalog.Venue{}, fmt.Errorf("find venue by id: %w", err)
	}

	return venue, nil
}

func (r *Repository) GetVenues(ctx context.Context) ([]catalog.Venue, error) {
	query := `
		SELECT id, name, address, city, created_at
		FROM venues
	`

	var venues = make([]catalog.Venue, 0)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get venues: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var venue catalog.Venue

		if err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.Address,
			&venue.City,
			&venue.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("iterate venues: %w", err)
		}

		venues = append(venues, venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate venues: %w", err)
	}

	return venues, nil
}

func (r *Repository) CreateSession(ctx context.Context, input catalog.CreateSessionInput) (catalog.Session, error) {
	query := `
		INSERT INTO sessions (venue_id, title, starts_at, sales_opens_at, sales_closes_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, venue_id, title, starts_at, sales_opens_at, sales_closes_at, status, created_at, updated_at
	`

	var session catalog.Session

	err := r.db.QueryRow(
		ctx,
		query,
		input.VenueID,
		input.Title,
		input.StartsAt,
		input.SalesOpensAt,
		input.SalesClosesAt,
		catalog.SessionStatusDraft,
	).Scan(
		&session.ID,
		&session.VenueID,
		&session.Title,
		&session.StartsAt,
		&session.SalesOpensAt,
		&session.SalesClosesAt,
		&session.Status,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return catalog.Session{}, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

func (r *Repository) FindSessionByID(ctx context.Context, id string) (catalog.Session, error) {
	query := `
		SELECT id, venue_id, title, starts_at, sales_opens_at, sales_closes_at, status, created_at, updated_at
		FROM sessions
		WHERE id = $1
	`

	var session catalog.Session

	err := r.db.QueryRow(
		ctx, query, id,
	).Scan(
		&session.ID,
		&session.VenueID,
		&session.Title,
		&session.StartsAt,
		&session.SalesOpensAt,
		&session.SalesClosesAt,
		&session.Status,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return catalog.Session{}, catalog.ErrSessionNotFound
	}
	if err != nil {
		return catalog.Session{}, fmt.Errorf("find session by id: %w", err)
	}

	return session, nil
}

func (r *Repository) FindSessionsByVenueID(ctx context.Context, venueID string) ([]catalog.Session, error) {
	query := `
		SELECT id, venue_id, title, starts_at, sales_opens_at, sales_closes_at, status, created_at, updated_at
		FROM sessions
		WHERE venue_id = $1
	`

	var sessions = make([]catalog.Session, 0)

	rows, err := r.db.Query(ctx, query, venueID)
	if err != nil {
		return nil, fmt.Errorf("find sessions by venue id: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var session catalog.Session

		if err := rows.Scan(
			&session.ID,
			&session.VenueID,
			&session.Title,
			&session.StartsAt,
			&session.SalesOpensAt,
			&session.SalesClosesAt,
			&session.Status,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("iterate sessions: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}

	return sessions, nil
}

func (r *Repository) GetSessions(ctx context.Context) ([]catalog.Session, error) {
	query := `
		SELECT id, venue_id, title, starts_at, sales_opens_at, sales_closes_at, status, created_at, updated_at
		FROM sessions
	`

	var sessions = make([]catalog.Session, 0)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get sessions: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var session catalog.Session

		if err := rows.Scan(
			&session.ID,
			&session.VenueID,
			&session.Title,
			&session.StartsAt,
			&session.SalesOpensAt,
			&session.SalesClosesAt,
			&session.Status,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("iterate sessions: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}

	return sessions, nil
}
