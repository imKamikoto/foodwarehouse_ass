package structures

import (
	"fmt"
	"log/slog"
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

// AddBatch:
// - если есть место — кладём;
// - если нет — возвращаем ошибку (вытеснение НЕ здесь).
func (cs *ColdStorage) AddBatch(newBatch Batch) error {
	if cs.IsFull() {
		return fmt.Errorf("cold storage %s is full", cs.ID)
	}

	slog.Info(fmt.Sprintf(
		"🧊 [ColdStorage: %s]: Принимает новую поставку <%s>, <%s>, для <%s>\n",
		cs.ID, newBatch.ID, newBatch.Name, newBatch.Client,
	))
	cs.Batches = append(cs.Batches, newBatch)
	return nil
}

// RemoveOldest — удаляет и возвращает самую "старую" партию по ExpiryDate.
func (cs *ColdStorage) RemoveOldest() (*Batch, error) {
	if len(cs.Batches) == 0 {
		return nil, fmt.Errorf("cold storage %s empty (Batches size = 0)", cs.ID)
	}

	oldestBatch, oldestIdx := cs.GetOldestBatch()
	if oldestBatch == nil || oldestIdx < 0 {
		return nil, fmt.Errorf("failed to find oldest batch in %s", cs.ID)
	}

	slog.Info(fmt.Sprintf(
		"[ColdStorage: %s]: Утилизирует старый товар <%s>, <%s>, для <%s>\n",
		cs.ID, oldestBatch.ID, oldestBatch.Name, oldestBatch.Client,
	))

	cs.Batches = append(cs.Batches[:oldestIdx], cs.Batches[oldestIdx+1:]...)

	return oldestBatch, nil
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
