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
	loader := d.Warehouse.Loaders[d.lastIndex]
	return loader
}

func (d *Dispatcher) ProcessOrder(order Order) {
	loader := d.AssignLoader()
	if loader == nil {
		fmt.Println("Диспетчер: Нет доступных погрузчиков!")
		return
	}

	loader.ServeClient(d.Warehouse, order)
}
