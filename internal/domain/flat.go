package domain

type Flat struct {
	ID         int
	FlatNumber int // Есть условие "номер квартиры", но его почему-то нет в API.
	Price      int
	Rooms      int
	Status     Status
}

type Status string

const (
	StatusCreated      Status = "created"
	StatusApproved     Status = "approved"
	StatusDeclined     Status = "declined"
	StatusOnModeration Status = "moderating"
)

func (s Status) Validate() bool {
	switch s {
	case StatusCreated, StatusApproved, StatusDeclined, StatusOnModeration:
		return true
	}
	return false
}

func (s Status) String() string {
	return string(s)
}
