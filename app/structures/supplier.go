package structures

import (
	"math/rand"
	"time"
)

type Supplier struct {
	ID          string // айдишник поставщика
	Name        string // имя поставщика
	ProductType string // имя поставляемого продукта (молочко, творог, сыр тд и тп)
	Batches     []*Batch
}

func NewSupplier(id string, name string, prodType string) *Supplier {
	return &Supplier{
		ID:          id,
		Name:        name,
		ProductType: prodType,
		Batches:     make([]*Batch, 0),
	}
}

func (s *Supplier) GenerateBatch(batchId string, name string, clientName string, nowTime time.Time) Batch {
	expiry := nowTime.Add(time.Duration(24+rand.Intn(4)*24) * time.Hour)
	batch := Batch{
		ID:          batchId,
		Name:        name,
		Client:      clientName,
		ArrivalDate: nowTime,
		ExpiryDate:  expiry,
	}
	s.Batches = append(s.Batches, &batch)
	return batch
}
