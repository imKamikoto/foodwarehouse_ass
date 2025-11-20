package structures

import (
	"log/slog"
	"time"

	"golang.org/x/exp/rand"

	"warehouse_app/utils"
)

type Simulation struct {
	Clock     time.Time     // модельное время
	Step      time.Duration // шаг симуляции
	Warehouse *Warehouse
	Metrics   *Metrics

	Suppliers []*Supplier
	Stores    []*Store
	Loaders   []*Loader

	IDGen map[string]func() string // генераторы ID
	rng   *rand.Rand               // рандом (кто делает заказ и т.п.)
}

func NewSimulation(start time.Time, step time.Duration) *Simulation {
	// генераторы ID
	idGen := map[string]func() string{
		"loader":      utils.NewGenerator("🦺 loader", 1511).Next,
		"supplier":    utils.NewGenerator("🏭 supplier", 1250).Next,
		"store":       utils.NewGenerator("🏪 store", 2672).Next,
		"batch":       utils.NewGenerator("📦 batch", 3982).Next,
		"order":       utils.NewGenerator("🧾 order", 4092).Next,
		"coldStorage": utils.NewGenerator("🧊 coldStorage", 5927).Next,
	}

	metrics := &Metrics{}

	// мороз камеры
	cs := NewColdStorage(idGen["coldStorage"](), 5)

	// погрузчики
	loaders := []*Loader{
		CreateLoader(idGen["loader"]()),
		CreateLoader(idGen["loader"]()),
	}

	// склад + диспетчер
	warehouse := NewWarehouse([]*ColdStorage{cs}, loaders, metrics)

	// поставщики
	suppliers := []*Supplier{
		{ID: idGen["supplier"](), Name: "Молочная ферма", ProductType: "Молоко"},
		{ID: idGen["supplier"](), Name: "Сырзавод", ProductType: "Сыр"},
	}

	// магазины
	stores := []*Store{
		NewStore(idGen["store"](), "Перекресток"),
		NewStore(idGen["store"](), "Магнит"),
	}

	src := rand.NewSource(uint64(time.Now().UnixNano()))
	rng := rand.New(src)

	return &Simulation{
		Clock:     start,
		Step:      step,
		Warehouse: warehouse,
		Metrics:   metrics,
		Suppliers: suppliers,
		Stores:    stores,
		Loaders:   loaders,
		IDGen:     idGen,
		rng:       rng,
	}
}

// Выполнение одного шага симуляции
func (s *Simulation) SimulationStep() {
	s.Clock = s.Clock.Add(s.Step)

	// 1. Поставки: для примера — всегда одна поставка от случайного поставщика в случайный магазин
	simSup := s.Suppliers[s.rng.Intn(len(s.Suppliers))]
	simStore := s.Stores[s.rng.Intn(len(s.Stores))]

	productName := simSup.ProductType
	batch := simSup.GenerateBatch(
		s.IDGen["batch"](),
		productName,
		simStore.Name,
		s.Clock,
	)

	slog.Info("Generated new batch",
		"supplier", simSup.Name,
		"batch_id", batch.ID,
		"client", batch.Client,
		"product", batch.Name,
	)

	// Запрос диспетчеру на передачу партии
	if err := s.Warehouse.Dispatcher.ReceiveBatch(batch); err != nil {
		slog.Error("failed to receive batch", "error", err)
	}

	// 2. Заказы: например, с вероятностью 0.5 на шаг
	if s.rng.Float64() < 0.5 {
		s.createRandomOrder(productName)
	}
}

// отдельный метод для формирования заказа
func (s *Simulation) createRandomOrder(productName string) {
	orderStore := s.Stores[s.rng.Intn(len(s.Stores))]
	order := orderStore.CreateOrder(s.IDGen["order"]())
	order.ProductName = productName

	slog.Info("new order",
		"store", orderStore.Name,
		"order_id", order.ID,
		"product", order.ProductName,
	)

	s.Warehouse.Dispatcher.ProcessOrder(order, orderStore)
}

// RunSteps — пошаговый режим
func (s *Simulation) RunSteps(steps int) {
	for i := 0; i < steps; i++ {
		slog.Info("\n\nsimulation step", "index", i+1, "time", s.Clock)
		s.SimulationStep()
	}
	s.logFinalStats()
}

// RunUntil — автоматический режим: крутить симуляцию, пока не пройдёт duration модельного времени.
func (s *Simulation) RunUntil() {}

// logFinalStats — просто лог итоговых метрик
func (s *Simulation) logFinalStats() {
	slog.Info("=== Итоговая статистика ===\n",
		"received", s.Metrics.Received,
		"discarded", s.Metrics.Discarded,
		"delivered", s.Metrics.Delivered,
	)

	slog.Info("warehouse cameras",
		"count", len(s.Warehouse.Cameras),
	)

	if len(s.Warehouse.Cameras) > 0 {
		slog.Info("camera state",
			"id", s.Warehouse.Cameras[0].ID,
			"batches", len(s.Warehouse.Cameras[0].Batches),
		)
	}
}
