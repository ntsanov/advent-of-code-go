package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	dot        = rune('.')
	gearSymbol = rune('*')
)

type Gear struct {
	Parts []Number
	Raw   int
	Col   int
}

func (n Number) isAdjacent(raw, col int) bool {
	if raw > n.Raw+1 || raw < n.Raw-1 {
		return false
	}
	if col < n.Start-1 || col > n.End {
		return false
	}
	return true
}

type Number struct {
	rune  []rune
	Start int
	End   int
	Raw   int
}

func (n Number) Int() int {
	number, err := strconv.Atoi(string(n.rune))
	if err != nil {
		slog.Error("not a valid number", "err", string(n.rune))
		os.Exit(-1)
	}
	return number
}

type Symbols struct {
	Matrix [][]rune
}

func (s Symbols) HasSymbol(startRaw, startCol, endRaw, endCol int) bool {
	for _, raw := range s.Matrix[startRaw:endRaw] {
		for _, s := range raw[startCol:endCol] {
			if s != dot {
				return true
			}
		}
	}
	return false
}

func main() {
	var (
		partNumbers  = [][]Number{}
		symbolMatrix = [][]rune{}
		gears        = []Gear{}
		input, err   = os.ReadFile("input.txt")
	)
	if err != nil {
		slog.Error("could not read input", "err", err)
	}
	lines := strings.Split(string(input), "\n")
	for rawIdx, line := range lines {
		var (
			raw              = []Number{}
			number           = []rune{}
			numStart         int
			symbolMatrixLine = []rune{}
		)
		if len(line) == 0 {
			continue
		}
		for colIdx, r := range line {
			// Dot
			if r == dot {
				// Not currently in a number
				if len(number) == 0 {
					symbolMatrixLine = append(symbolMatrixLine, dot)
					continue
					// Reached end of number
				} else {
					raw = append(raw, Number{
						rune:  number,
						Start: numStart,
						End:   numStart + len(number),
						Raw:   rawIdx,
					})
					number = []rune{}
				}
				symbolMatrixLine = append(symbolMatrixLine, dot)
				// Digit
			} else if unicode.IsDigit(r) {
				if len(number) == 0 {
					numStart = colIdx
				}
				number = append(number, r)
				symbolMatrixLine = append(symbolMatrixLine, dot)
				// Symbol
			} else {
				if len(number) != 0 {
					raw = append(raw, Number{
						rune:  number,
						Start: numStart,
						End:   numStart + len(number),
						Raw:   rawIdx,
					})
					number = []rune{}
				}

				if r == gearSymbol {
					gears = append(gears, Gear{Raw: rawIdx, Col: colIdx})
				}
				symbolMatrixLine = append(symbolMatrixLine, r)
			}
		}
		if len(number) != 0 {
			raw = append(raw, Number{
				rune:  number,
				Start: numStart,
				End:   numStart + len(number),
				Raw:   rawIdx,
			})
		}
		partNumbers = append(partNumbers, raw)
		symbolMatrix = append(symbolMatrix, symbolMatrixLine)
	}

	var (
		partNumSum int64
		ratioSum   int64
	)
	symbols := Symbols{symbolMatrix}
	for rawIdx, raw := range partNumbers {
		for _, n := range raw {
			var (
				endIdx     = n.Start + len(n.rune)
				minCol int = n.Start
				maxCol int = endIdx
				minRaw int = rawIdx
				maxRaw int = rawIdx
			)
			if n.Start > 0 {
				minCol--
			}
			if endIdx < len(symbolMatrix[rawIdx])-1 {
				maxCol++
			}
			if rawIdx > 0 {
				minRaw--
			}
			if rawIdx < len(symbolMatrix)-1 {
				maxRaw++
			}
			if symbols.HasSymbol(minRaw, minCol, maxRaw+1, maxCol) {
				partNumSum += int64(n.Int())
			}

		}
	}
	for _, gear := range gears {
		var (
			minRaw       = gear.Raw
			maxRaw       = gear.Raw
			validNumbers = []int64{}
		)
		if minRaw > 0 {
			minRaw = minRaw - 1
		}
		if maxRaw < len(symbolMatrix)-1 {
			maxRaw = maxRaw + 1
		}
		for _, raw := range partNumbers[minRaw : maxRaw+1] {
			for _, partNumber := range raw {
				if partNumber.isAdjacent(gear.Raw, gear.Col) {
					number, _ := strconv.ParseInt(string(partNumber.rune), 10, 64)
					validNumbers = append(validNumbers, number)
				}
			}
		}
		if len(validNumbers) > 1 {
			ratio := int64(1)
			for _, n := range validNumbers {
				ratio *= n
			}
			slog.Info("valid gear", "ratio", ratio, "col", gear.Col, "raw", gear.Raw)
			ratioSum += ratio
		}
	}
	slog.Info("great success", "sum_p1", partNumSum, "sum_p2", ratioSum)
}
