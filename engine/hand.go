package engine

import (
	"fmt"
	"strings"
)

// DecisionLogEntry, bir el için alınan stratejik kararın tüm adımlarını kaydeder.
type DecisionLogEntry struct {
	Key          string   `json:"key"`
	Actions      []string `json:"actions"`
	FinalAction  string   `json:"final_action"`
	IsDeviation  bool     `json:"is_deviation"`
	IsFallback   bool     `json:"is_fallback"`
}

type Hand struct {
	ID            string             `json:"id"`
	BoxID         string             `json:"box_id"`
	Cards         []Card             `json:"cards"`
	BetAmount     float64            `json:"bet_amount"`
	Payout        float64            `json:"payout"`
	Result        string             `json:"result"`
	IsSplitChild  bool               `json:"is_split_child"`
	IsDoubled     bool               `json:"is_doubled"`
	DecisionTrace []DecisionLogEntry `json:"decision_trace"`
	FinalAction   string             `json:"-"` // Bu loglama için geçici bir alandır
	StrategyActions []string         `json:"-"` // Bu loglama için geçici bir alandır
}

func NewHand(bet float64, boxID string, handID int) *Hand {
	return &Hand{
		ID:            fmt.Sprintf("%s-%d", boxID, handID),
		BoxID:         boxID,
		Cards:         []Card{},
		BetAmount:     bet,
		DecisionTrace: []DecisionLogEntry{},
	}
}

func NewSplitHand(from *Hand, handID int) *Hand {
	return &Hand{
		ID:            fmt.Sprintf("%s-%d", from.BoxID, handID),
		BoxID:         from.BoxID,
		Cards:         []Card{},
		BetAmount:     from.BetAmount,
		IsSplitChild:  true,
		DecisionTrace: append([]DecisionLogEntry{}, from.DecisionTrace...),
	}
}

func (h *Hand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

// SetDecisionTrace, bir el için önerilen stratejiyi ve nihai kararı kaydeder.
// Bu fonksiyon, executeBoxActions içinde çağrılır.
func (h *Hand) SetDecisionTrace(actions []string) {
	h.StrategyActions = actions
}

// FinalizeDecision, el için nihai kararı kaydeder.
func (h *Hand) FinalizeDecision(key, finalAction string, isDeviation, isFallback bool) {
	logEntry := DecisionLogEntry{
		Key:         key,
		Actions:     h.StrategyActions,
		FinalAction: finalAction,
		IsDeviation: isDeviation,
		IsFallback:  isFallback,
	}
	h.DecisionTrace = append(h.DecisionTrace, logEntry)
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
