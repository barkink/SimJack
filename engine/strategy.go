package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Strategy interface {
	GetAction(hand *Hand, dealerUp Card) string
	FallbackAction(failed string) string
	DecideInsurance() bool
}

type DynamicStrategy struct {
	AcceptInsurance bool                `json:"decide_insurance"`
	Fallback        string              `json:"fallback"`
	Actions         map[string][]string `json:"actions"`
}

var strategyDirectory string = "strategies"

func SetStrategyDirectory(path string) error {
	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		return fmt.Errorf("invalid strategy directory: %s", path)
	}
	strategyDirectory = path
	return nil
}

func LoadStrategyFromFile(name string) (Strategy, error) {
	path := filepath.Join(strategyDirectory, fmt.Sprintf("%s.json", name))
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load strategy file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var strat DynamicStrategy
	if err := decoder.Decode(&strat); err != nil {
		return nil, fmt.Errorf("failed to decode strategy: %w", err)
	}

	if strat.Fallback == "" {
		return nil, fmt.Errorf("strategy file missing 'fallback' value")
	}

	return &strat, nil
}

func (s *DynamicStrategy) GetAction(hand *Hand, dealerUp Card) string {
	val := hand.CalculateValue()
	key := ""

	if hand.CanSplit() {
		key = fmt.Sprintf("pair_%s_vs_%s", hand.Cards[0].Rank, strings.ToUpper(dealerUp.Rank))
	} else if hasAceButNotPair(hand) {
		key = fmt.Sprintf("soft_%d_vs_%s", val, strings.ToUpper(dealerUp.Rank))
	} else {
		key = fmt.Sprintf("hard_%d_vs_%s", val, strings.ToUpper(dealerUp.Rank))
	}

	if actions, ok := s.Actions[key]; ok && len(actions) > 0 {
		return actions[0] // ana aksiyon
	}
	return s.Fallback
}

func (s *DynamicStrategy) FallbackAction(failed string) string {
	for _, actions := range s.Actions {
		if len(actions) > 1 && actions[0] == failed {
			fallback := strings.TrimSpace(actions[1])
			if fallback != "" {
				return fallback // fallback aksiyon
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
