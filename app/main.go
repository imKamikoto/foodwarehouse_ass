package main

import (
	"fmt"
	"time"
	. "warehouse_app/structures"
	. "warehouse_app/utils"

	"golang.org/x/exp/rand"
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

	// создание объектов
	metrics := &Metrics{}
	coldStorage := NewColdStorage(idGenerator["coldStorage"](), 5)
	loaders := []*Loader{
		CreateLoader(idGenerator["loader"]()),
		CreateLoader(idGenerator["loader"]()),
	}
	warehouse := NewWarehouse([]*ColdStorage{coldStorage}, loaders, metrics)

	suppliers := []*Supplier{
		{ID: idGenerator["supplier"](), Name: "Молочная ферма", ProductType: "Молоко"},
		{ID: idGenerator["supplier"](), Name: "Сырзавод", ProductType: "Сыр"},
	}

	stores := []*Store{
		{ID: idGenerator["store"](), Name: "Перекресток"},
		{ID: idGenerator["store"](), Name: "Магнит"},
	}

	now := time.Now()

	for step := 0; step < 1; step++ {
		fmt.Printf("\n--- Шаг %d ---\n\n", step+1)

		// 1. Приход новой партии от случайного поставщика для случайного магазина
		sup := suppliers[0]
		store := stores[0]
		batch := sup.GenerateBatch(idGenerator["batch"](), "Молоко блин", store.Name, now)

		fmt.Printf("Поставщик <%s> привёз партию <%s> для клиента <%s>, товар: <%s>\n",
			sup.Name, batch.ID, batch.Client, batch.Name)

		status, err := warehouse.AcceptBatch(batch)
		if !status || err != nil {
			fmt.Println(err)
		}

		// 2. Иногда приходит заказ от магазина
		if step%2 == 1 {
			orderStore := stores[rand.Intn(len(stores))]
			order := orderStore.CreateOrder(idGenerator["order"]())
			warehouse.Dispatcher.ProcessOrder(order)
		}

		time.Sleep(300 * time.Millisecond)
	}
}
