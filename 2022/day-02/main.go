package main

import (
	"log/slog"
	"os"
	"strings"
)

const (
	X int = iota
	Y
	Z
)

// Truth table
//
//	X Y Z
//
// A
// B
// C
// var tt = [][]int{
// 	{0, 1, -1},
// 	{-1, 0, 1},
// 	{1, -1, 0},
// }

var tt1 = map[string]int{
	"AX": 4, "AY": 8, "AZ": 3,
	"BX": 1, "BY": 5, "BZ": 9,
	"CX": 7, "CY": 2, "CZ": 6,
}

var tt2 = map[string]int{
	"AX": 3, "AY": 4, "AZ": 8,
	"BX": 1, "BY": 5, "BZ": 9,
	"CX": 2, "CY": 6, "CZ": 7,
}

func points(key string) int {
	return tt1[key]
}

func main() {
	var totalScore1, totalScore2 int
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could no read input", "err", err.Error)
	}
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		rule := strings.ReplaceAll(line, " ", "")
		totalScore1 += tt1[rule]
		totalScore2 += tt2[rule]
		slog.Info("info", "line", line, "score", tt2[rule])
	}
	slog.Info("successfuly calculated score", "part1", totalScore1, "part2", totalScore2)
}
