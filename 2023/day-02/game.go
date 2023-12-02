package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	red   = "red"
	green = "green"
	blue  = "blue"
)

type Game struct {
	Id    int
	Grabs []Combo
}

type Combo struct {
	Green int
	Blue  int
	Red   int
}

func (c Combo) FitsIn(target Combo) bool {
	if target.Green >= c.Green &&
		target.Blue >= c.Blue &&
		target.Red >= c.Red {
		return true
	} else {
		return false
	}
}

func (c Combo) Add(n Combo) Combo {
	return Combo{
		Green: c.Green + n.Green,
		Red:   c.Red + n.Red,
		Blue:  c.Blue + n.Blue,
	}
}

func (g Game) Sum() Combo {
	sum := Combo{}
	for _, grab := range g.Grabs {
		sum = sum.Add(grab)
	}
	return sum
}

func (g Game) MinReq() Combo {
	req := Combo{}
	for _, grab := range g.Grabs {
		if req.Red < grab.Red {
			req.Red = grab.Red
		}
		if req.Green < grab.Green {
			req.Green = grab.Green
		}
		if req.Blue < grab.Blue {
			req.Blue = grab.Blue
		}
	}
	return req
}

func ParseGames(input []byte) ([]Game, error) {
	var (
		lines = strings.Split(string(input), "\n")
		games = make([]Game, 0, len(lines))
	)
	for _, line := range lines {
		game := Game{
			Grabs: make([]Combo, 0),
		}
		if gameRecord := strings.Split(line, ":"); len(gameRecord) > 1 {
			var (
				gameStr   = gameRecord[0]
				recordStr = gameRecord[1]
				err       error
			)
			if gameId := strings.Split(gameStr, " "); len(gameId) > 1 {
				game.Id, err = strconv.Atoi(gameId[1])
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("invalid game id string: %s", gameStr)
			}
			grabs := strings.Split(recordStr, ";")
			for _, grabStr := range grabs {
				grab := Combo{}
				colors := strings.Split(grabStr, ",")
				for _, grabColors := range colors {
					cnt := 0
					grabColors = strings.TrimSpace(grabColors)
					colorCnt := strings.Split(grabColors, " ")
					if len(colorCnt) != 2 {
						return nil, fmt.Errorf("invalid color grab: %s", grabColors)
					}
					if cnt, err = strconv.Atoi(colorCnt[0]); err != nil {
						return nil, err
					}
					switch colorCnt[1] {
					case red:
						grab.Red += cnt
					case green:
						grab.Green += cnt
					case blue:
						grab.Blue += cnt
					default:
						return nil, fmt.Errorf("unknown color: %s", colorCnt[1])
					}

				}
				game.Grabs = append(game.Grabs, grab)
			}
		}
		games = append(games, game)
	}
	return games, nil
}
