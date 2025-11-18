package structures

import "fmt"

type Dispatcher struct {
	Metrics   *Metrics
	Warehouse *Warehouse
	lastIndex int
}

func (d *Dispatcher) AssignLoader() *Loader {
	if len(d.Warehouse.Loaders) == 0 {
		return nil
	}
	d.lastIndex = (d.lastIndex + 1) % len(d.Warehouse.Loaders)
	fmt.Printf("[Dispatcher]: Назначаем погрузчика(ов): %s\n", d.Warehouse.Loaders[d.lastIndex].ID)
	return d.Warehouse.Loaders[d.lastIndex]
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
		fmt.Printf("[Dispatcher]: Партия <%s> утилизирована\n", discarded.ID)
	}

	d.Metrics.LogArrival(b)
	return nil
}

func (d *Dispatcher) ProcessOrder(order Order, store *Store) {
	loader := d.AssignLoader()
	if loader == nil {
		fmt.Println("Диспетчер: Нет доступных погрузчиков!")
		return
	}
	fmt.Printf("[Dispatcher]: Обработка заказа %s - %s от магазина %s\n", order.ID, order.ProductName, store.Name)
	loader.ServeClient(d.Warehouse, order, store)
}
