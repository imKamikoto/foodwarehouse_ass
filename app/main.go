package main

import (
	"fmt"
	"time"

	. "warehouse_app/structures"
	. "warehouse_app/utils"
)

func main() {
	fmt.Println("Начало работы программы")

	idGenerator := map[string]func() string{
		"loader":      NewGenerator("loader", 1511).Next,
		"supplier":    NewGenerator("supplier", 1250).Next,
		"store":       NewGenerator("store", 2672).Next,
		"batch":       NewGenerator("batch", 3982).Next,
		"order":       NewGenerator("order", 4092).Next,
		"coldStorage": NewGenerator("coldStorage", 5927).Next,
	}

	metrics := &Metrics{}
	coldStorage := NewColdStorage(idGenerator["coldStorage"](), 5)
	loaders := []*Loader{
		CreateLoader(idGenerator["loader"]()),
		CreateLoader(idGenerator["loader"]()),
	}

	warehouse := NewWarehouse([]*ColdStorage{coldStorage}, loaders, metrics)
	dispatcher := warehouse.Dispatcher

	suppliers := []*Supplier{
		{ID: idGenerator["supplier"](), Name: "Молочная ферма", ProductType: "Молоко"},
		{ID: idGenerator["supplier"](), Name: "Сырзавод", ProductType: "Сыр"},
	}

	stores := []*Store{
		{ID: idGenerator["store"](), Name: "Перекресток"},
		{ID: idGenerator["store"](), Name: "Магнит"},
	}

	now := time.Now()

	for step := 0; step < 2; step++ {
		fmt.Printf("\n--- Шаг %d ---\n\n", step+1)

		sup := suppliers[0]
		store := stores[0]

		productName := sup.ProductType
		batch := sup.GenerateBatch(idGenerator["batch"](), productName, store.Name, now)

		fmt.Printf(
			"Поставщик <%s> привёз партию <%s> для клиента <%s>, товар: <%s>\n",
			sup.Name, batch.ID, batch.Client, batch.Name,
		)

		if err := dispatcher.ReceiveBatch(batch); err != nil {
			fmt.Println("Ошибка при приёме партии:", err)
		}

		fmt.Printf("\n\nwarehouse cameras capasity %d\n", len(warehouse.Cameras))
		fmt.Printf("camera capasity %d\n\n", len(warehouse.Cameras[0].Batches))

		if step%2 == 1 {
			orderStore := stores[0]
			order := orderStore.CreateOrder(idGenerator["order"]())

			order.ProductName = productName

			fmt.Printf(
				"\nМагазин <%s> оформляет заказ <%s> на товар <%s>\n\n",
				orderStore.Name, order.ID, order.ProductName,
			)

			dispatcher.ProcessOrder(order, orderStore)
		}

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Printf(
		"\n=== Итоговая статистика ===\nПринято партий:   %d\nСписано партий:   %d\nОтгружено партий: %d\n",
		metrics.Received, metrics.Discarded, metrics.Delivered,
	)

	fmt.Printf("warehouse cameras capasity %d\n", len(warehouse.Cameras))
	fmt.Printf("camera capasity %d", len(warehouse.Cameras[0].Batches))
}
