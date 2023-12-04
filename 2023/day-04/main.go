package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
)

type Ticket struct {
	Id            int
	Numbers       []int
	WiningNumbers []int
	Cnt           int
}

func (t Ticket) Matches() int {
	matchCnt := 0
	for _, n := range t.Numbers {
		for _, wn := range t.WiningNumbers {
			if n == wn {
				matchCnt++
			}
		}
	}
	return matchCnt
}

func ParseNumbers(input string) ([]int, error) {
	input = strings.TrimSpace(input)
	numbersStr := strings.Split(input, " ")
	out := make([]int, 0, len(numbersStr))
	for _, numStr := range numbersStr {
		if numStr == "" {
			continue
		}
		n, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func ParseTickets(input []byte) ([]Ticket, error) {
	var (
		lines   = strings.Split(string(input), "\n")
		tickets = make([]Ticket, 0, len(lines))
	)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		ticket := Ticket{Cnt: 1}
		if ticketRecord := strings.Split(line, ":"); len(ticketRecord) > 1 {
			var (
				ticketHeader = ticketRecord[0]
				recordStr    = ticketRecord[1]
				err          error
			)
			if ticketId := strings.Split(ticketHeader, " "); len(ticketId) > 1 {
				ticket.Id, err = strconv.Atoi(ticketId[len(ticketId)-1])
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("invalid ticket id string: %s", ticketHeader)
			}
			if numbers := strings.Split(recordStr, "|"); len(numbers) != 2 {
				return nil, fmt.Errorf("invalid ticket record: %s", recordStr)
			} else {
				ticket.Numbers, err = ParseNumbers(numbers[0])
				if err != nil {
					return nil, err
				}
				ticket.WiningNumbers, err = ParseNumbers(numbers[1])
				if err != nil {
					return nil, err
				}
			}
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func main() {
	var (
		matchPowSum  int64
		totalTickets int
	)
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read input", "err", err)
		os.Exit(-1)
	}
	tickets, err := ParseTickets(input)
	if err != nil {
		slog.Error("could not parse tickets", "err", err)
		os.Exit(-1)
	}
	for i := 0; i < len(tickets); i++ {
		matches := tickets[i].Matches()
		if matches > 0 {
			matchPowSum += int64(math.Pow(2, float64(matches-1)))
		}
		for m := 1; m <= matches; m++ {
			if i+m > len(tickets)-1 {
				break
			}
			tickets[i+m].Cnt += tickets[i].Cnt
		}
		totalTickets += tickets[i].Cnt

	}
	slog.Info("great success", "p1_sum", matchPowSum, "p2_sum", totalTickets)
}
