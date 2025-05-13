package engine

import (
	"fmt"
	"sort"
	"strings"
)

func GetPerfectPairPayout(c1, c2 Card) (float64, string) {
	r1 := strings.TrimSpace(strings.ToUpper(c1.Rank))
	r2 := strings.TrimSpace(strings.ToUpper(c2.Rank))
	s1 := strings.TrimSpace(strings.ToLower(c1.Suit))
	s2 := strings.TrimSpace(strings.ToLower(c2.Suit))

	if r1 != r2 {
		return 0, "none"
	}
	if s1 == s2 {
		return 25, "Perfect Pair"
	}
	if (isRed(s1) && isRed(s2)) || (isBlack(s1) && isBlack(s2)) {
		return 12, "Colored Pair"
	}
	return 6, "Mixed Pair"
}

func Get21Plus3Payout(cards []Card) (float64, string) {
	if len(cards) != 3 {
		return 0, "none"
	}
	ranks := []int{cardNumeric(cards[0].Rank), cardNumeric(cards[1].Rank), cardNumeric(cards[2].Rank)}
	suits := []string{cards[0].Suit, cards[1].Suit, cards[2].Suit}
	uniqueRanks := map[string]int{}
	for _, c := range cards {
		uniqueRanks[c.Rank]++
	}
	sameSuit := suits[0] == suits[1] && suits[1] == suits[2]
	consecutive := isStraight(ranks)

	if sameSuit && len(uniqueRanks) == 1 {
		return 100, "Suited Trips"
	}
	if sameSuit && consecutive {
		return 40, "Straight Flush"
	}
	for _, v := range uniqueRanks {
		if v == 3 {
			return 30, "Three of a Kind"
		}
	}
	if consecutive {
		return 10, "Straight"
	}
	if sameSuit {
		return 5, "Flush"
	}
	return 0, "none"
}

func isRed(suit string) bool {
	return suit == "hearts" || suit == "diamonds"
}

func isBlack(suit string) bool {
	return suit == "spades" || suit == "clubs"
}

func cardNumeric(r string) int {
	switch r {
	case "A":
		return 1
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	default:
		var v int
		fmt.Sscanf(r, "%d", &v)
		return v
	}
}

func isStraight(vals []int) bool {
	if len(vals) != 3 {
		return false
	}
	sort.Ints(vals)
	// Normal sıra kontrolü
	if vals[0]+1 == vals[1] && vals[1]+1 == vals[2] {
		return true
	}
	// A, Q, K
	if vals[0] == 1 && vals[1] == 12 && vals[2] == 13 {
		return true
	}
	// A, 2, 3
	if vals[0] == 1 && vals[1] == 2 && vals[2] == 3 {
		return true
	}
	return false
}
