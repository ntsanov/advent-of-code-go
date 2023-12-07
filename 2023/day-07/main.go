package main

import (
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

type HandType int

type FindTypeFunc func([5]int) HandType

const (
	HighCard = HandType(iota)
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

var HandTypeString = map[HandType]string{
	HighCard:     "High Card",
	OnePair:      "One Pair",
	TwoPair:      "Two Piar",
	ThreeOfAKind: "Three of a kind",
	FullHouse:    "Full House",
	FourOfAKind:  "Four of a kind",
	FiveOfAKind:  "Five of a kind",
}

type Hand struct {
	Bid      int
	Cards    [5]int
	cardsStr string
	Type     HandType
	Jorkers  int
}

func IncType(h HandType) HandType {
	switch h {
	case HighCard:
		return OnePair
	case OnePair:
		return ThreeOfAKind
	case TwoPair:
		return FullHouse
	case ThreeOfAKind:
		return FourOfAKind
	// Can't have FullHouse and a Jocker
	case FullHouse:
		slog.Error("can't have Full house and a Joker")
		return FullHouse
	case FourOfAKind:
		return FiveOfAKind
	case FiveOfAKind:
		slog.Debug("5 Jokers, nice")
		return FiveOfAKind
	default:
		slog.Error("unknown HandType")
		return h
	}
}

func (h Hand) Stronger(other Hand) bool {
	if h.Type > other.Type {
		return true
	}

	if h.Type < other.Type {
		return false
	}
	// Equal, check cards
	for i, card := range h.Cards {
		if card > other.Cards[i] {
			return true
		}
		if card < other.Cards[i] {
			return false
		}
	}
	return true
}

func getHandType(uniqueCards, highestRep int) HandType {
	switch uniqueCards {
	case 1:
		return FiveOfAKind
	case 2:
		if highestRep == 4 {
			return FourOfAKind
		}
		return FullHouse
	case 3:
		if highestRep == 3 {
			return ThreeOfAKind
		}
		return TwoPair
	case 4:
		return OnePair
	default:
		return HighCard
	}
}

func findTypeJocker(cards [5]int) HandType {
	var (
		cardByType = map[int]int{}
		highestRep = 0
		jokers     = 0
	)
	for _, card := range cards {
		// Joker
		if card == 1 {
			jokers++
			continue
		}
		if match, found := cardByType[card]; found {
			match++
			if highestRep < match {
				highestRep = match
			}
			cardByType[card] = match
		} else {
			cardByType[card] = 1
		}
	}

	// We want to count jokers as distinct cards at this stage
	uniqueCards := len(cardByType) + jokers
	handType := getHandType(uniqueCards, highestRep)
	for i := 0; i < jokers; i++ {
		handType = IncType(handType)
	}
	return handType
}

func findType(cards [5]int) HandType {
	var (
		uniqueCards = map[int]int{}
		highestRep  = 0
	)
	for _, card := range cards {
		if match, found := uniqueCards[card]; found {
			match++
			if highestRep < match {
				highestRep = match
			}
			uniqueCards[card] = match
		} else {
			uniqueCards[card] = 1
		}
	}
	return getHandType(len(uniqueCards), highestRep)
}

func NewHand(cardsString string, bid int, CardsMap map[rune]int, findType FindTypeFunc) Hand {
	var (
		cards  = [5]int{}
		jokers int
	)
	for i, r := range cardsString {
		if v, found := CardsMap[r]; found {
			if r == 'J' {
				jokers++
			}
			cards[i] = v
		} else {
			slog.Error("unknown card", "card", string(r))
		}
	}
	return Hand{
		Bid:      bid,
		Cards:    cards,
		Type:     findType(cards),
		cardsStr: cardsString,
		Jorkers:  jokers,
	}
}

var cardsMap = map[rune]int{
	'2': 2, '3': 3, '4': 4, '5': 5,
	'6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
	'J': 11, 'Q': 12, 'K': 13, 'A': 14,
}

var cardsMapJoker = map[rune]int{
	'2': 2, '3': 3, '4': 4, '5': 5,
	'6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
	'J': 1, 'Q': 12, 'K': 13, 'A': 14,
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		slog.Error("could not read input", "err", err)
		os.Exit(-1)
	}
	var (
		lines        = strings.Split(string(input), "\n")
		hands        = []Hand{}
		handsWJokers = []Hand{}
	)
	for _, line := range lines {
		if values := strings.Split(line, " "); len(values) == 2 {
			cards := values[0]
			if bid, err := strconv.Atoi(values[1]); err == nil {
				hand := NewHand(cards, bid, cardsMap, findType)
				handWJoker := NewHand(cards, bid, cardsMapJoker, findTypeJocker)
				hands = append(hands, hand)
				handsWJokers = append(handsWJokers, handWJoker)
			} else {
				slog.Error("could not parse bid", "bid", bid, "err", err)
			}
		}
	}
	sort.Slice(hands, func(i, j int) bool {
		return hands[j].Stronger(hands[i])
	})

	var (
		sumP1 int64
	)
	for i, hand := range hands {
		rank := i + 1
		sumP1 += int64(rank * hand.Bid)
		slog.Debug("", "rank", strconv.Itoa(rank), "cards", hand.cardsStr, "type", HandTypeString[hand.Type], "bid", hand.Bid)
	}
	slog.Info("great success", "p1_sum", sumP1)

	var (
		sumP2 int64
	)
	sort.Slice(handsWJokers, func(i, j int) bool {
		return handsWJokers[j].Stronger(handsWJokers[i])
	})
	for i, hand := range handsWJokers {
		rank := i + 1
		sumP2 += int64(rank * hand.Bid)
		slog.Debug("", "rank", strconv.Itoa(rank), "cards", hand.cardsStr, "type", HandTypeString[hand.Type], "bid", hand.Bid)
	}
	slog.Info("great success", "p2_sum", sumP2)
}
