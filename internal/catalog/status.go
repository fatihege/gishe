package catalog

// session status

type SessionStatus string

const (
	SessionStatusDraft     SessionStatus = "draft"
	SessionStatusPublished SessionStatus = "published"
	SessionStatusCancelled SessionStatus = "cancelled"
	SessionStatusCompleted SessionStatus = "completed"
)

func ParseSessionStatus(value string) (SessionStatus, error) {
	status := SessionStatus(value)

	if !status.IsValid() {
		return "", ErrInvalidStatus
	}

	return status, nil
}

func (s SessionStatus) IsValid() bool {
	switch s {
	case SessionStatusDraft, SessionStatusPublished, SessionStatusCancelled, SessionStatusCompleted:
		return true
	default:
		return false
	}
}

// seat status

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusHeld      SeatStatus = "held"
	SeatStatusSold      SeatStatus = "sold"
)

func ParseSeatStatus(value string) (SeatStatus, error) {
	status := SeatStatus(value)

	if !status.IsValid() {
		return "", ErrInvalidStatus
	}

	return status, nil
}

func (s SeatStatus) IsValid() bool {
	switch s {
	case SeatStatusAvailable, SeatStatusHeld, SeatStatusSold:
		return true
	default:
		return false
	}
}

func (s SeatStatus) CanTransitionTo(next SeatStatus) bool {
	switch s {
	case SeatStatusAvailable:
		return next == SeatStatusHeld
	case SeatStatusHeld:
		return next == SeatStatusAvailable || next == SeatStatusSold
	case SeatStatusSold:
		return next == SeatStatusAvailable
	default:
		return false
	}
}
