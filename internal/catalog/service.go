package catalog

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

type CreateVenueInput struct {
	Name    string
	Address string
	City    string
}

func (s *Service) CreateVenue(ctx context.Context, venue CreateVenueInput) (Venue, error) {
	name := strings.TrimSpace(venue.Name)
	address := strings.TrimSpace(venue.Address)
	city := strings.TrimSpace(venue.City)

	if name == "" || address == "" || city == "" {
		return Venue{}, fmt.Errorf("name, address, and city required")
	}

	newVenue, err := s.repository.CreateVenue(ctx, CreateVenueInput{
		Name:    name,
		Address: address,
		City:    city,
	})
	if err != nil {
		return Venue{}, err
	}

	return newVenue, nil
}

func (s *Service) GetVenues(ctx context.Context) ([]Venue, error) {
	if venues, err := s.repository.GetVenues(ctx); err != nil {
		return []Venue{}, err
	} else {
		return venues, nil
	}
}
