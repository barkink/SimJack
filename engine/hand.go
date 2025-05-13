package engine

import (
	"strings"
)

// external link to Box registry for global access (registered by engine)
var boxRegistry = map[int]*Box{}

func registerBox(box *Box) {
	boxRegistry[box.ID] = box
}

type Hand struct {
	Cards        []Card
	IsDoubled    bool
	IsSplitChild bool
	BetAmount    float64
	Result       string
	BoxID        int
	Payout       float64
}

func NewHand(bet float64, boxID int) *Hand {
	return &Hand{
		Cards:     []Card{},
		BetAmount: bet,
		BoxID:     boxID,
		Payout:    0,
	}
}

func NewSplitHand(from *Hand) *Hand {
	return &Hand{
		Cards:        []Card{},
		BetAmount:    from.BetAmount,
		IsSplitChild: true,
		BoxID:        from.BoxID,
		Payout:       0,
	}
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

func (h *Hand) CalculateValue() int {
	total := 0
	aces := 0
	for _, c := range h.Cards {
		v := c.Value()
		total += v
		if c.Rank == "A" {
			aces++
		}
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.CalculateValue() == 21 && !h.IsSplitChild
}

func (h *Hand) IsBust() bool {
	return h.CalculateValue() > 21
}

func (h *Hand) CanSplit() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == h.Cards[1].Rank
}

func (h *Hand) MarkAsDoubled() {
	h.IsDoubled = true
	h.BetAmount *= 2
}

func (h *Hand) String() string {
	var parts []string
	for _, c := range h.Cards {
		parts = append(parts, c.String())
	}
	return strings.Join(parts, ";")
}
