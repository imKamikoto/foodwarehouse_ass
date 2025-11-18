package structures

import (
	"fmt"
)

type ColdStorage struct {
	ID       string
	Capacity int
	Batches  []Batch
}

func NewColdStorage(id string, capacity int) *ColdStorage {
	return &ColdStorage{
		ID:       id,
		Capacity: capacity,
		Batches:  make([]Batch, 0, capacity),
	}
}

// - если есть место — просто добавляем;
// - если места нет — удаляем самую старую партию и добавляем новую;
// - если по какой-то причине удалить не удалось — возвращаем false.
func (cs *ColdStorage) AddBatch(newBatch Batch) bool {
	if cs.IsFull() {
		ok, err := cs.removeOldest()
		if !ok || err != nil {
			return false
		}
	}
	fmt.Printf(
		"[ColdStorage: %s]: Принимает новую поставку <%s>, <%s>, для <%s>",
		cs.ID, newBatch.ID, newBatch.Name, newBatch.Client,
	)
	cs.Batches = append(cs.Batches, newBatch)
	return true
}

func (cs *ColdStorage) removeOldest() (bool, error) {
	if len(cs.Batches) == 0 {
		return false, fmt.Errorf("cold storage empty (Batches size = 0)")
	}

	oldestBatch, oldestIdx := cs.GetOldestBatch()
	if oldestBatch == nil || oldestIdx < 0 {
		return false, fmt.Errorf("failed to find oldest batch")
	}

	fmt.Printf(
		"[ColdStorage: %s]: Утилизирует старый товар <%s>, <%s>, для <%s>\n",
		cs.ID, oldestBatch.ID, oldestBatch.Name, oldestBatch.Client,
	)

	cs.Batches = append(cs.Batches[:oldestIdx], cs.Batches[oldestIdx+1:]...)

	return true, nil
}

func (cs *ColdStorage) IsFull() bool {
	return len(cs.Batches) >= cs.Capacity
}

func (cs *ColdStorage) GetOldestBatch() (*Batch, int) {
	if len(cs.Batches) == 0 {
		return nil, -1
	}

	oldestIdx := 0
	oldestExpiry := cs.Batches[0].ExpiryDate

	for i := 1; i < len(cs.Batches); i++ {
		if cs.Batches[i].ExpiryDate.Before(oldestExpiry) {
			oldestExpiry = cs.Batches[i].ExpiryDate
			oldestIdx = i
		}
	}

	return &cs.Batches[oldestIdx], oldestIdx
}

func (cs *ColdStorage) TakeBatchByClientAndName(clientName string, name string) (*Batch, error) {
	for i, b := range cs.Batches {
		if b.Client == clientName && b.Name == name {
			batch := b

			cs.Batches = append(cs.Batches[:i], cs.Batches[i+1:]...)

			return &batch, nil
		}
	}
	return nil, fmt.Errorf("партия для клиента %q не найдена", clientName)
}
