package structures

type Order struct {
	ID          string
	Client      string // владелец заказа (магазин, который делает заказ партии у склада)
	ProductName string
}
