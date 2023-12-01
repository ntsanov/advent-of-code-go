package main

import (
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

var stringNumbers = [9]string{
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

type numberMatch struct {
	idx int
	num int
}

func findStringNumbers(src string) []numberMatch {
	// iterate in reverse becuase there is string "eightwo" in input and it expects 8
	matches := []numberMatch{}
	for i, numStr := range stringNumbers {
		if idx := strings.Index(src, numStr); idx != -1 {
			matches = append(matches, numberMatch{idx, i + 1})
			lastIdx := strings.LastIndex(src, numStr)
			if lastIdx != idx {
				matches = append(matches, numberMatch{lastIdx, i + 1})
			}
		}
	}
	return matches
}

func main() {
	var sumP1, sumP2 int64
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read input", err, err.Error())
	}
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		stringNumMatches := findStringNumbers(line)
		matches := []numberMatch{}
		for idx, r := range line {
			if num, err := strconv.Atoi(string(r)); err == nil {
				matches = append(matches, numberMatch{idx, num})
			}
		}
		if len(matches) == 0 {
			continue
		}
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].idx < matches[j].idx
		})
		coord := 10*matches[0].num + matches[len(matches)-1].num
		sumP1 += int64(coord)
		matches = append(matches, stringNumMatches...)
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].idx < matches[j].idx
		})
		coord = 10*matches[0].num + matches[len(matches)-1].num
		sumP2 += int64(coord)
	}

	slog.Info("great success", "sum_p1", sumP1, "part2_sum", sumP2)
}
