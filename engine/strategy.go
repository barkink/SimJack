package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var strategyDirectory string = "strategies"

func SetStrategyDirectory(path string) error {
	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		return fmt.Errorf("invalid strategy directory: %s", path)
	}
	strategyDirectory = path
	return nil
}

// Strategy interface - oyuncuya atanacak stratejiler bunu implement etmeli
type Strategy interface {
	GetAction(hand *Hand, dealerUp Card) string
	FallbackAction(failed string) string
	DecideInsurance() bool
	String() string
}

// Deviation: count'a göre farklı aksiyon
type DeviationRule struct {
	AtCount int    `json:"at_count"`
	Action  string `json:"action"`
}

// Bahis rampası: count >= MinCount ise BetUnit kullanılır
type BetRampTier struct {
	MinCount int     `json:"min_count"`
	BetUnit  float64 `json:"bet_unit"`
}

// Ana strateji tipi: BaseStrategy (dynamic), Deviations ve BetRamp içerir
type CountingStrategy struct {
	BaseStrategy    *DynamicStrategy
	Deviations      map[string]DeviationRule `json:"deviations"`
	BetRamp         []BetRampTier            `json:"bet_ramp"`
	AcceptInsurance bool                     `json:"decide_insurance"`
	Deck            *Deck                    `json:"-"` // runtime'da atanır
	CountingEnabled bool                     // 💡 yeni alan
	Name            string 
}

func (s *CountingStrategy) GetAction(hand *Hand, dealerUp Card) string {
	val := hand.CalculateValue()
	dealerKey := getDealerRankKey(dealerUp)
	key := ""
	if hand.CanSplit() {
		key = fmt.Sprintf("pair_%s_vs_%s", hand.Cards[0].Rank, dealerKey)
	} else if hasAceButNotPair(hand) {
		key = fmt.Sprintf("soft_%d_vs_%s", val, dealerKey)
	} else {
		key = fmt.Sprintf("hard_%d_vs_%s", val, dealerKey)
	}

	if dev, ok := s.Deviations[key]; ok && s.Deck != nil && s.Deck.GetRunningCount() >= dev.AtCount {
		hand.DecisionTrace = append(hand.DecisionTrace, DecisionLogEntry{
			Key:    key,
			Action: dev.Action + " (deviation)",
		})
		return dev.Action
	}
	action := s.BaseStrategy.GetAction(hand, dealerUp)
	// base strategy key mevcut mi kontrol edelim
	_, has := s.BaseStrategy.Actions[key]
	suffix := " (main)"
	if !has {
		suffix = " (fallback)"
	}

	hand.DecisionTrace = append(hand.DecisionTrace, DecisionLogEntry{
		Key:    key,
		Action: action + suffix,
	})
	return action
}

func (s *CountingStrategy) FallbackAction(failed string) string {
	return s.BaseStrategy.FallbackAction(failed)
}

func (s *CountingStrategy) DecideInsurance() bool {
	if s.Deck != nil {
		trueCount := int(float64(s.Deck.GetRunningCount()) / (float64(len(s.Deck.Cards)) / 52.0))
		if trueCount >= 3 {
			return true // kart sayan oyuncu için TC >= 3'te insurance alınır
		}
	}
	return s.AcceptInsurance // kart saymıyorsa ya da deck atanmadıysa config'teki davranışı uygula
}

func (s *CountingStrategy) GetBetUnit(base float64) float64 {
	if s.Deck == nil {
		return base
	}
	count := s.Deck.GetRunningCount()
	for i := len(s.BetRamp) - 1; i >= 0; i-- {
		if count >= s.BetRamp[i].MinCount {
			return base * s.BetRamp[i].BetUnit
		}
	}
	return base
}

// JSON formatına uygun geçici yapı
type CountingStrategyFile struct {
	Fallback        string                   `json:"fallback"`
	DecideInsurance bool                     `json:"decide_insurance"`
	Actions         map[string][]string      `json:"actions"`
	Deviations      map[string]DeviationRule `json:"deviations"`
	BetRamp         []BetRampTier            `json:"bet_ramp"`
	CountingEnabled bool                     `json:"counting_enabled"`
}

// JSON dosyasından CountingStrategy yükler
func LoadCountingStrategyFromFile(name string) (*CountingStrategy, error) {
	path := filepath.Join(strategyDirectory, fmt.Sprintf("%s.json", name))
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data CountingStrategyFile
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode strategy: %w", err)
	}

	base := &DynamicStrategy{
		Fallback:        data.Fallback,
		AcceptInsurance: data.DecideInsurance,
		Actions:         data.Actions,
	}

	return &CountingStrategy{
		BaseStrategy:    base,
		Deviations:      data.Deviations,
		BetRamp:         data.BetRamp,
		AcceptInsurance: data.DecideInsurance,
		Deck:            nil, // deck dışarıdan yalnızca aktifse atanır
		CountingEnabled: data.CountingEnabled,
		Name:            name,
	}, nil
}

// Uyum için eski fonksiyon ismi korunur
type DynamicStrategy struct {
	Fallback        string              `json:"fallback"`
	AcceptInsurance bool                `json:"decide_insurance"`
	Actions         map[string][]string `json:"actions"`
}

func (s *DynamicStrategy) GetAction(hand *Hand, dealerUp Card) string {
	val := hand.CalculateValue()
	dealerKey := getDealerRankKey(dealerUp)
	key := ""
	if hand.CanSplit() {
		key = fmt.Sprintf("pair_%s_vs_%s", hand.Cards[0].Rank, dealerKey)
	} else if hasAceButNotPair(hand) {
		key = fmt.Sprintf("soft_%d_vs_%s", val, dealerKey)
	} else {
		key = fmt.Sprintf("hard_%d_vs_%s", val, dealerKey)
	}
	if actions, ok := s.Actions[key]; ok && len(actions) > 0 {
		return actions[0]
	}
	return s.Fallback
}

func (s *DynamicStrategy) FallbackAction(failed string) string {
	for _, actions := range s.Actions {
		if len(actions) > 1 && actions[0] == failed {
			fallback := strings.TrimSpace(actions[1])
			if fallback != "" {
				return fallback
			}
		}
	}
	if failed == "double" || failed == "split" {
		return "hit"
	}
	return s.Fallback
}

func (s *DynamicStrategy) DecideInsurance() bool {
	return s.AcceptInsurance
}

func LoadStrategyFromFile(name string) (Strategy, error) {
	return LoadCountingStrategyFromFile(name)
}

func hasAceButNotPair(h *Hand) bool {
	if len(h.Cards) != 2 {
		return false
	}
	hasAce := false
	for _, c := range h.Cards {
		if c.Rank == "A" {
			hasAce = true
		} else if c.Rank == h.Cards[0].Rank && h.Cards[0].Rank != "A" {
			return false
		}
	}
	return hasAce
}

func getDealerRankKey(card Card) string {
	switch card.Rank {
	case "J", "Q", "K":
		return "10"
	default:
		return strings.ToUpper(card.Rank)
	}
}

func (s *CountingStrategy) String() string {
	return s.Name
}
