package main

import (
	"errors"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	SeedToSoil            = "seed-to-soil"
	SoilToFertilizer      = "soil-to-fertilizer"
	FertilizerToWater     = "fertilizer-to-water"
	WaterToLight          = "water-to-light"
	LightToTemperature    = "light-to-temperature"
	TemperatureToHumidity = "temperature-to-humidity"
	HumidityToLocation    = "humidity-to-location"
)

type Instruction struct {
	Src int64
	Dst int64
	Len int64
}

type Step struct {
	Instructions []Instruction
}

func (s Step) MaxKey() int64 {
	var max int64
	for _, ins := range s.Instructions {
		lastKey := ins.Src + ins.Len - 1
		if lastKey > max {
			max = lastKey
		}
	}
	return max
}

func (s Step) GetDst(src int64) int64 {
	for _, ins := range s.Instructions {
		if src >= ins.Src && src < ins.Src+ins.Len {
			offset := src - ins.Src
			return ins.Dst + offset
		}
	}
	return src
}

func InstructioFromString(input string) (Instruction, error) {
	intStr := strings.Split(input, " ")
	ints, err := StringsToInts(intStr)
	if err != nil {
		return Instruction{}, err
	}
	if len(ints) != 3 {
		return Instruction{}, errors.New("wrong Instruction format")
	}
	return Instruction{
		Src: ints[1],
		Dst: ints[0],
		Len: ints[2],
	}, nil
}

func StringsToInts(input []string) ([]int64, error) {
	out := make([]int64, 0, len(input))
	for _, numberStr := range input {
		number, err := strconv.ParseInt(numberStr, 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, number)
	}
	return out, nil
}

func FindMin(start, l int64, steps []Step) int64 {
	var min int64 = int64(math.MaxInt64)
	startSeed := start
	// TODO some recursion magic to speed things up
	for j := int64(0); j < l; j++ {
		seed := startSeed + j
		key := seed
		for _, step := range steps {
			key = step.GetDst(key)
		}
		if key < min {
			min = key
		}
	}
	return min
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read input", "err", err)
	}
	paragraphs := strings.Split(string(input), "\n\n")
	header := strings.Split(paragraphs[0], ":")
	paragraphs = paragraphs[1:]

	seedsStr := strings.Split(strings.TrimSpace(header[1]), " ")
	seeds, err := StringsToInts(seedsStr)
	if err != nil {
		slog.Error("seeds not an int", "err", err)
	}

	almanac := []Step{}
	for _, p := range paragraphs {
		step := Step{}
		if lines := strings.Split(p, "\n"); len(lines) > 0 {
			for i := 1; i < len(lines); i++ {
				if lines[i] == "" {
					continue
				}
				ins, err := InstructioFromString(lines[i])
				if err != nil {
					slog.Error("could not read instruction", "err", err)
					os.Exit(-1)
				}
				step.Instructions = append(step.Instructions, ins)
			}
		}
		almanac = append(almanac, step)
	}
	var (
		min          = int64(math.MaxInt64)
		minWithRange = int64(math.MaxInt64)
	)
	for _, seed := range seeds {
		key := seed
		for _, step := range almanac {
			key = step.GetDst(key)
		}
		if key < min {
			min = key
		}
	}
	for i := 0; i < len(seeds)-2; i = i + 2 {
		min := FindMin(seeds[i], seeds[i+1], almanac)
		if min < minWithRange {
			minWithRange = min
		}
	}
	slog.Info("great success", "min_location", min, "min_with_range", minWithRange)
}
