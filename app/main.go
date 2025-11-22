package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	. "warehouse_app/structures"
)

func main() {
	modeFlag := flag.String("mode", "step", "simulation mode: 'step' or 'auto'")
	stepFlag := flag.String("step", "500ms", "simulation step duration (e.g. '200ms', '1s')")
	stepsFlag := flag.Int("steps", 5, "number of steps in 'step' mode")
	// durationFlag := flag.String("duration", "10s", "simulation duration in 'auto' mode (e.g. '10s', '2m')")

	flag.Parse()

	step, err := time.ParseDuration(*stepFlag)
	if err != nil {
		slog.Error("failed to parse step duration", "value", *stepFlag, "error", err)
		os.Exit(1)
	}

	start := time.Now()
	slog.Info("Начало работы программы",
		"mode", *modeFlag,
		"step", step,
	)

	sim := NewSimulation(start, step)

	switch *modeFlag {
	case "step":
		if *stepsFlag <= 0 {
			slog.Error("steps must be > 0", "steps", *stepsFlag)
			os.Exit(1)
		}
		slog.Info("Запуск в пошаговом режиме",
			"steps", *stepsFlag,
		)
		sim.RunSteps(*stepsFlag)

	// case "auto":
	// 	duration, err := time.ParseDuration(*durationFlag)
	// 	if err != nil {
	// 		slog.Error("failed to parse duration", "value", *durationFlag, "error", err)
	// 		os.Exit(1)
	// 	}
	// 	if duration <= 0 {
	// 		slog.Error("duration must be > 0", "duration", duration)
	// 		os.Exit(1)
	// 	}
	// 	slog.Info("Запуск в автоматическом режиме",
	// 		"duration", duration,
	// 	)
	// 	sim.RunUntil(duration)

	default:
		slog.Error("unknown mode", "mode", *modeFlag)
		os.Exit(1)
	}

	sim.LogFinalStats()
	sim.LogStoreStats()
}
