package structures

import (
	"fmt"
	"log/slog"
	"strings"
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
	// сейчас первый всегда занят, для демонстрации, что выбор между двумя погрузчиками
	// TODO что-то придумать с этими погрузчиками, пока не попадают под план того, как работает система
	loaders := []*Loader{
		NewLoader(idGen["loader"](), true),
		NewLoader(idGen["loader"](), false),
	}

	// склад + диспетчер
	warehouse := NewWarehouse([]*ColdStorage{cs}, loaders, metrics)

	// поставщики
	suppliers := []*Supplier{
		NewSupplier(idGen["supplier"](), "Молочная ферма", "Молоко"),
		NewSupplier(idGen["supplier"](), "Сырзавод", "Сыр"),
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
	if s.rng.Float64() < 0.8 {
		s.createRandomOrder(productName)
	}
	// s.createRandomOrder(i, productName)
}

// отдельный метод для формирования заказа
func (s *Simulation) createRandomOrder(productName string) {
	orderStore := s.Stores[s.rng.Intn(len(s.Stores))]
	order := orderStore.CreateOrder(s.IDGen["order"](), productName)
	// order.ProductName = productName

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
}

// RunUntil — автоматический режим: крутить симуляцию, пока не пройдёт duration модельного времени.
func (s *Simulation) RunUntil() {}

// logFinalStats — просто лог итоговых метрик
func (s *Simulation) LogFinalStats() {
	var b strings.Builder

	fmt.Fprintf(&b, "\n=== Итоговая статистика ===")
	fmt.Fprintf(&b, "\nreceived = %d", s.Metrics.Received)
	fmt.Fprintf(&b, "\ndiscarded = %d", s.Metrics.Discarded)
	fmt.Fprintf(&b, "\ndelivered = %d", s.Metrics.Delivered)

	slog.Info("warehouse cameras",
		"count", len(s.Warehouse.Cameras),
	)
	fmt.Fprintf(&b, "\n\n=== Информация о складе после выполнения программы ===\n")

	if len(s.Warehouse.Cameras) > 0 {
		fmt.Fprintf(&b, "\nКамера хранения : %s\n", s.Warehouse.Cameras[0].ID)
		fmt.Fprintf(&b, "Количество товаров внутри : %d\n\n", len(s.Warehouse.Cameras[0].Batches))
	}
	slog.Info(b.String())
}

func (s *Simulation) LogStoreStats() {
	var b strings.Builder

	fmt.Fprintf(&b, "\n==================== 🏪 Stores summary ====================\n\n")

	for _, store := range s.Stores {
		fmt.Fprintf(&b, "🏪 Магазин: %s (%s)\n", store.Name, store.ID)

		// Заказы
		fmt.Fprintf(&b, "  📑 Заказы (%d):\n", len(store.Orders))
		if len(store.Orders) == 0 {
			fmt.Fprintf(&b, "    — нет оформленных заказов\n")
		} else {
			for _, o := range store.Orders {
				product := o.ProductName
				if product == "" {
					product = "<не указан>"
				}
				fmt.Fprintf(&b, "    • [%s] %s (товар: %s)\n", o.ID, o.Client, product)
			}
		}

		// Ассортимент
		fmt.Fprintf(&b, "  📦 Ассортимент (%d):\n", len(store.Assortment))
		if len(store.Assortment) == 0 {
			fmt.Fprintf(&b, "    — полки пусты :(\n")
		} else {
			for _, bch := range store.Assortment {
				fmt.Fprintf(
					&b,
					"    • [%s] %s (от поставки для %s)\n",
					bch.ID, bch.Name, bch.Client,
				)
			}
		}

		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "===========================================================\n")

	slog.Info(b.String())
}
