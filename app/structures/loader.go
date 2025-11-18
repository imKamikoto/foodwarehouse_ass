package structures

type Loader struct {
	ID           string // id погрузчика
	IsBusy       bool   // занятость
	CurrentOrder *Order // текущий заказ
}

func CreateLoader(id string) *Loader {
	return &Loader{
		ID:           id,
		IsBusy:       false,
		CurrentOrder: nil,
	}
}

func (l *Loader) ServeClient(w *Warehouse, order Order) {
	// l.IsBusy = true
	// defer func() { l.IsBusy = false }()

	// for {
	// 	batch, ok := w.FetchBatchForClient(order.Client)
	// 	if !ok {
	// 		// fmt.Printf("Погрузчик %d: Больше нет партий для клиента %s\n", l.ID, order.Client)
	// 		break
	// 	}
	// 	w.Metrics.LogDelivery(batch)
	// 	// fmt.Printf("Погрузчик %d: Отгружена партия %s для клиента %s\n",
	// 	// l.ID, batch.ID, batch.Client)
	// 	// для наглядности чуть "замедлим" отгрузку
	// 	time.Sleep(200 * time.Millisecond)
	// }
	// // fmt.Printf("Погрузчик %d: Завершил обслуживание заказа %d клиента %s\n",
	// // l.ID, order.ID, order.Client)
}
