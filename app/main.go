package main

import (
	// "fmt"

	"log/slog"
	"time"

	// "time"
	. "warehouse_app/structures"
	// . "warehouse_app/utils"
)

func main() {
	slog.Info(" Начало работы программы ")

	start := time.Now()
	step := 500 * time.Millisecond

	sim := NewSimulation(start, step)

	// Пошаговый режим:
	sim.RunSteps(5)

	// // Автоматический режим (например, 10 минут модельного времени):
	// sim.RunUntil(10 * time.Minute)
}
