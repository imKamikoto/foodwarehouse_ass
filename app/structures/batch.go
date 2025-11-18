package structures

import "time"

type Batch struct {
	ID          string
	Name        string
	Client      string
	ArrivalDate time.Time // дата прибытия
	ExpiryDate  time.Time // дата когда умирает(
}

func (b Batch) IsExpired() (bool, time.Time) {
	return !b.ExpiryDate.After(time.Now()), b.ExpiryDate
}
