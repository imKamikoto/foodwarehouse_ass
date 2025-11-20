package structures

import (
	"fmt"
	"log/slog"
	"time"
)

type Loader struct {
	ID           string
	IsBusy       bool
	CurrentOrder *Order
}

func CreateLoader(id string) *Loader {
	return &Loader{
		ID:           id,
		IsBusy:       false,
		CurrentOrder: nil,
	}
}

func (l *Loader) ServeClient(w *Warehouse, order Order, store *Store) {
	slog.Info(fmt.Sprintf("[Loader: %s]: Начинает отгрузку заказа <%s>, для <%s>\n", l.ID, order.ID, order.Client))
	l.IsBusy = true
	defer func() { l.IsBusy = false }()

	productName := order.ProductName

	for {

		time.Sleep(time.Second * 1)

		batch, err := w.FetchBatchForClient(order.Client, productName)
		if err != nil {
			break
		}

		slog.Info(fmt.Sprintf("[Loader: %s]: Взял товар <%s>-<%s> для <%s>-<%s>\n", l.ID, batch.ID, batch.Name, store.ID, store.Name))
		w.Metrics.LogDelivery(*batch)

		// nyyyyyy
		store.AddAssortment(batch)

		time.Sleep(time.Second * 1)
	}
}
