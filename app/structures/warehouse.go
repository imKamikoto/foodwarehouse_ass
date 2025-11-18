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

// Приём партии с вытеснением самой старой, если всё занято
func (w *Warehouse) AcceptBatch(b Batch) (bool, error) {
	fmt.Printf("\n[Warehouse]: Поступление партии: %s, %s\n", b.Name, b.ID)

	// поиск незаполненных камер
	for _, cs := range w.Cameras {
		if !cs.IsFull() {
			if cs.AddBatch(b) {
				return true, nil
			}
			return false, fmt.Errorf("неудалось добавить партию в незаполненную камеру")
		}
	}

	cameraWithOldestBatch := w.getColdStorageWithOldestBatch()
	if cameraWithOldestBatch == nil {
		return false, fmt.Errorf("Нет доступных камер для размещения партии")
	}

	res := cameraWithOldestBatch.AddBatch(b)
	if res {
		return true, nil
	}
	return false, fmt.Errorf("Неудалось принять партию")
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
			return batch, nil
		}
	}
	return nil, fmt.Errorf("Необходимая партия не найдена в холодильных камерах")
}
