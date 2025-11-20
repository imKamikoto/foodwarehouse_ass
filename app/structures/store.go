package structures

import (
	"fmt"
	"log/slog"
)

type Store struct {
	ID         string
	Name       string
	Orders     []Order
	Assortment []*Batch
}

func NewStore(id string, name string) *Store {
	return &Store{
		ID:         id,
		Name:       name,
		Orders:     make([]Order, 0),
		Assortment: make([]*Batch, 0),
	}
}

func (s *Store) CreateOrder(orderId string, productName string) Order {
	newOrder := Order{ID: orderId, Client: s.Name, ProductName: productName}
	s.Orders = append(s.Orders, newOrder)
	return newOrder
}

func (s *Store) AddAssortment(batch *Batch) {
	slog.Info(fmt.Sprintf("[Store: %s]: Заказанный товар добавлен в ассортимент магазина\n", s.ID))
	s.Assortment = append(s.Assortment, batch)
}
