package structures

import (
	"math/rand"
	"time"
)

type Supplier struct {
	ID          string // айдишник поставщика
	Name        string // имя поставщика
	ProductType string // имя поставляемого продукта (молочко, творог, сыр тд и тп)
}

func (s *Supplier) GenerateBatch(batchId string, name string, clientName string, nowTime time.Time) Batch {
	expiry := nowTime.Add(time.Duration(24+rand.Intn(4)*24) * time.Hour)
	return Batch{
		ID:          batchId,
		Name:        name,
		Client:      clientName,
		ArrivalDate: nowTime,
		ExpiryDate:  expiry,
	}
}
