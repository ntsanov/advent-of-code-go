package main

import (
	"log/slog"
	"os"
)

func main() {
	limit := Combo{
		Red:   12,
		Green: 13,
		Blue:  14,
	}
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read input", "err", err)
		os.Exit(-1)
	}
	games, err := ParseGames(input)
	if err != nil {
		slog.Error("could not parse input", "err", err)
		os.Exit(-1)
	}
	var (
		idSum    int
		powerSum int64
	)
	for _, game := range games {
		rq := game.MinReq()
		powerSum += int64(rq.Blue) * int64(rq.Green) * int64(rq.Red)
		isPossible := true
		for _, grab := range game.Grabs {
			if !grab.FitsIn(limit) {
				isPossible = false
				break
			}
		}
		if isPossible {
			idSum += game.Id
		}

	}
	slog.Info("great success", "id_sum", idSum, "power_sum", powerSum)
}
