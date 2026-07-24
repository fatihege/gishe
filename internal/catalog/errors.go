package catalog

import "errors"

var (
	ErrInvalidStatus          = errors.New("invalid status")
	ErrVenueFieldsRequired    = errors.New("name, address, and city are required")
	ErrVenueNotFound          = errors.New("venue not found")
	ErrSessionNotFound        = errors.New("session not found")
	ErrInvalidVenue           = errors.New("invalid venue")
	ErrSessionTitleRequired   = errors.New("session title is required")
	ErrSessionTimesRequired   = errors.New("session times are required")
	ErrInvalidSessionSchedule = errors.New("invalid session schedule")
)
