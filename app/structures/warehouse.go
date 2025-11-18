package structures

import (
	"fmt"
	"time"
)

type Warehouse struct {
	Cameras    []*ColdStorage
	Loaders    []*Loader
	Dispatcher *Dispatcher
	Metrics    *Metrics
}

func NewWarehouse(storages []*ColdStorage, loaders []*Loader, metrics *Metrics) *Warehouse {
	w := &Warehouse{
		Cameras: storages,
		Loaders: loaders,
		Metrics: metrics,
	}
	dispatcher := &Dispatcher{
		Metrics:   metrics,
		Warehouse: w,
	}
	w.Dispatcher = dispatcher
	return w
}

// AcceptBatch:
// - возвращает вытесненную партию (если была) и ошибку (если совсем не смогли принять).
func (w *Warehouse) AcceptBatch(b Batch) (*Batch, error) {
	fmt.Printf("\n[Warehouse]: Поступление партии: %s, %s\n", b.Name, b.ID)

	// 1. Ищем незаполненную камеру
	for _, cs := range w.Cameras {
		if !cs.IsFull() {
			if err := cs.AddBatch(b); err != nil {
				return nil, fmt.Errorf("не удалось добавить партию в незаполненную камеру: %w", err)
			}
			return nil, nil // никто не вытеснен
		}
	}

	// 2. Все камеры заполнены — ищем камеру с самой "старой" партией
	cameraWithOldestBatch := w.getColdStorageWithOldestBatch()
	if cameraWithOldestBatch == nil {
		return nil, fmt.Errorf("нет доступных камер для размещения партии")
	}

	discarded, err := cameraWithOldestBatch.RemoveOldest()
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить старую партию: %w", err)
	}

	if err := cameraWithOldestBatch.AddBatch(b); err != nil {
		return discarded, fmt.Errorf("не удалось добавить новую партию после вытеснения: %w", err)
	}

	return discarded, nil
}

func (w *Warehouse) getColdStorageWithOldestBatch() *ColdStorage {
	var (
		oldestStorage *ColdStorage
		oldestTime    time.Time
		initialized   bool
	)

	for _, cs := range w.Cameras {
		b, _ := cs.GetOldestBatch()
		if b == nil {
			continue
		}

		if !initialized || b.ExpiryDate.Before(oldestTime) {
			oldestTime = b.ExpiryDate
			oldestStorage = cs
			initialized = true
		}
	}
	return oldestStorage
}

// поиск партии по названию продукта
func (w *Warehouse) FetchBatchForClient(client string, name string) (*Batch, error) {
	for _, cs := range w.Cameras {
		batch, err := cs.TakeBatchByClientAndName(client, name)
		if err == nil {
			fmt.Printf("[Warehouser]: Взят товар для клиента %s, с именем %s\n", client, name)
			return batch, nil
		}
	}
	return nil, fmt.Errorf("Необходимая партия не найдена в холодильных камерах")
}
