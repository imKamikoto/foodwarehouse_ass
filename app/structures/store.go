package structures

type Store struct {
	ID     string
	Name   string
	Orders []Order
}

func (s *Store) CreateOrder(orderId string) Order {
	newOrder := Order{ID: orderId, Client: s.Name}
	s.Orders = append(s.Orders, newOrder)
	return newOrder
}
