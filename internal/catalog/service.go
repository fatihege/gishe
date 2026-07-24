package catalog

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

type CreateVenueInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

func (s *Service) CreateVenue(ctx context.Context, input CreateVenueInput) (Venue, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.City = strings.TrimSpace(input.City)

	if input.Name == "" || input.Address == "" || input.City == "" {
		return Venue{}, ErrVenueFieldsRequired
	}

	venue, err := s.repository.CreateVenue(ctx, input)
	if err != nil {
		return Venue{}, err
	}

	return venue, nil
}

func (s *Service) GetVenues(ctx context.Context) ([]Venue, error) {
	return s.repository.GetVenues(ctx)
}

type CreateSessionInput struct {
	VenueID       string     `json:"venue_id"`
	Title         string     `json:"title"`
	StartsAt      time.Time  `json:"starts_at"`
	SalesOpensAt  time.Time  `json:"sales_opens_at"`
	SalesClosesAt *time.Time `json:"sales_closes_at"`
}

func (s *Service) CreateSession(ctx context.Context, input CreateSessionInput) (Session, error) {
	input.VenueID = strings.TrimSpace(input.VenueID)
	input.Title = strings.TrimSpace(input.Title)

	if input.VenueID == "" {
		return Session{}, ErrInvalidVenue
	}

	if input.Title == "" {
		return Session{}, ErrSessionTitleRequired
	}

	if input.StartsAt.IsZero() || input.SalesOpensAt.IsZero() {
		return Session{}, ErrSessionTimesRequired
	}

	if !input.SalesOpensAt.Before(input.StartsAt) {
		return Session{}, ErrInvalidSessionSchedule
	}

	_, err := s.repository.FindVenueByID(ctx, input.VenueID)
	if errors.Is(err, ErrVenueNotFound) {
		return Session{}, ErrInvalidVenue
	}
	if err != nil {
		return Session{}, err
	}

	if input.SalesClosesAt != nil {
		salesClosesAt := *input.SalesClosesAt

		if salesClosesAt.IsZero() ||
			!salesClosesAt.After(input.SalesOpensAt) ||
			salesClosesAt.After(input.StartsAt) {
			return Session{}, ErrInvalidSessionSchedule
		}
	}

	session, err := s.repository.CreateSession(ctx, input)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s *Service) GetSessionById(ctx context.Context, id string) (Session, error) {
	return s.repository.FindSessionByID(ctx, id)
}

func (s *Service) GetSessions(ctx context.Context) ([]Session, error) {
	return s.repository.GetSessions(ctx)
}

func (s *Service) GetSessionsByVenueID(ctx context.Context, venueID string) ([]Session, error) {
	_, err := s.repository.FindVenueByID(ctx, venueID)
	if errors.Is(err, ErrVenueNotFound) {
		return nil, ErrInvalidVenue
	}
	if err != nil {
		return nil, err
	}

	return s.repository.FindSessionsByVenueID(ctx, venueID)
}
