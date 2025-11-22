package structures

import (
	"fmt"
	"log/slog"
	"time"
)

type Dispatcher struct {
	Metrics   *Metrics
	Warehouse *Warehouse
	lastIndex int
}

func (d *Dispatcher) AssignLoader() *Loader {
	if len(d.Warehouse.Loaders) == 0 {
		return nil
	}

	for {
		loaders := d.Warehouse.Loaders
		for i, loader := range loaders {
			if !loader.IsBusy {
				slog.Info(fmt.Sprintf(
					"🧑‍💼 [Dispatcher]: Назначаем погрузчика: %s\n",
					loader.ID,
				))
				return loaders[i]
			}
		}

		// сюда попадаем, если все заняты
		slog.Info("🧑‍💼 [Dispatcher]: все погрузчики заняты, ждём...")
		time.Sleep(100 * time.Millisecond)
	}
}

// Приём новой партии от поставщика
func (d *Dispatcher) ReceiveBatch(b Batch) error {
	discarded, err := d.Warehouse.AcceptBatch(b)
	if err != nil {
		return err
	}

	if discarded != nil {
		d.Metrics.LogDiscard(*discarded)
		// сделать что-то типо "Waste" для утилизированного товара?
		slog.Info(fmt.Sprintf("🧑‍💼 [Dispatcher]: Партия <%s> утилизирована\n", discarded.ID))
	}

	d.Metrics.LogArrival(b)
	return nil
}

func (d *Dispatcher) ProcessOrder(order Order, store *Store) {
	slog.Info(fmt.Sprintf("🧑‍💼 [Dispatcher]: Обработка заказа %s - %s от магазина %s\n", order.ID, order.ProductName, store.Name))
	loader := d.AssignLoader()
	if loader == nil {
		slog.Info("🧑‍💼 [Dispatcher]: Диспетчер: Нет доступных погрузчиков!")
		return
	}
	loader.ServeClient(d.Warehouse, order, store)
}
