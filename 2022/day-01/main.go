package main

import (
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var (
		elfFood       uint64
		threeFatElves = make([]uint64, 3)
		total         uint64
	)

	inputData, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read file", "error", err.Error())
		os.Exit(-1)
	}
	lines := strings.Split(string(inputData), "\n")
	for _, line := range lines {
		sort.Slice(threeFatElves, func(i, j int) bool {
			return threeFatElves[i] < threeFatElves[j]
		})
		if line == "" {
			if elfFood > threeFatElves[0] {
				threeFatElves[0] = elfFood
			}
			elfFood = 0
		} else {
			food, err := strconv.ParseUint(line, 10, 64)
			if err != nil {
				slog.Error("food not real", "error", err.Error())
				os.Exit(-1)
			}
			elfFood += food
		}
	}
	for _, elfFood := range threeFatElves {
		total += elfFood
	}
	slog.Info("found three fatest elves", "elves", threeFatElves, "total", total)
}
