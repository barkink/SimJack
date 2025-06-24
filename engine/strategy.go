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
	// GetAction, eylem listesini, fallback olup olmadığını, deviation olup olmadığını ve strateji anahtarını döndürür.
	GetAction(hand *Hand, dealerUp Card) (actions []string, isFallback bool, isDeviation bool, key string)
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

func (s *CountingStrategy) GetAction(hand *Hand, dealerUp Card) ([]string, bool, bool, string) {
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

	if dev, ok := s.Deviations[key]; ok && s.Deck != nil && s.getTrueCount() >= float64(dev.AtCount) {
		// Deviation (sapma) varsa, bu tek ve öncelikli eylemdir.
		return []string{dev.Action}, false, true, key
	}
	
	// Sapma yoksa, temel stratejiyi çağır.
	actions, isFallback := s.BaseStrategy.GetAction(hand, dealerUp)
	return actions, isFallback, false, key
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
	trueCount := s.getTrueCount()
	for i := len(s.BetRamp) - 1; i >= 0; i-- {
		if trueCount >= float64(s.BetRamp[i].MinCount) {
			return base * s.BetRamp[i].BetUnit
		}
	}
	return base
}

// JSON formatına uygun geçici yapı
type CountingStrategyFile struct {
	Fallback        string                   `json:"fallback"`
	Actions         map[string][]string      `json:"actions"`
	Deviations      map[string]DeviationRule `json:"deviations"`
	BetRamp         []BetRampTier            `json:"bet_ramp"`
	CountingEnabled bool                     `json:"counting_enabled"`
	AcceptInsurance  bool                     `json:"decide_insurance"`
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
		Actions:         data.Actions,
	}

	return &CountingStrategy{
		BaseStrategy:    base,
		Deviations:      data.Deviations,
		BetRamp:         data.BetRamp,
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

func (s *DynamicStrategy) GetAction(hand *Hand, dealerUp Card) ([]string, bool) {
	val := hand.CalculateValue()
	dealerKey := getDealerRankKey(dealerUp)
	// ... (key oluşturma mantığı aynı)
	key := ""
	if hand.CanSplit() {
		key = fmt.Sprintf("pair_%s_vs_%s", hand.Cards[0].Rank, dealerKey)
	} else if hasAceButNotPair(hand) {
		key = fmt.Sprintf("soft_%d_vs_%s", val, dealerKey)
	} else {
		key = fmt.Sprintf("hard_%d_vs_%s", val, dealerKey)
	}

	if actions, ok := s.Actions[key]; ok && len(actions) > 0 {
		return actions, false // Ana strateji, fallback değil
	}
	return []string{s.Fallback}, true // Bu bir fallback
}

func (s *DynamicStrategy) DecideInsurance() bool {
	return s.AcceptInsurance
}

func LoadStrategyFromFile(name string) (Strategy, error) {
	return LoadCountingStrategyFromFile(name)
}

func hasAceButNotPair(h *Hand) bool {
	// Sadece iki kartlı eller soft olabilir (split sonrası gelenler dahil)
	if len(h.Cards) != 2 {
		return false
	}
	// Eğer kartlar aynıysa (çift ise), soft değildir. AA çifti split'e girer.
	if h.Cards[0].Rank == h.Cards[1].Rank {
		return false
	}
	// Kartlardan biri As ise, el soft'tur.
	return h.Cards[0].Rank == "A" || h.Cards[1].Rank == "A"
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

func (s *CountingStrategy) getTrueCount() float64 {
	if s.Deck == nil || len(s.Deck.Cards) == 0 {
		return 0
	}
	remainingDecks := float64(len(s.Deck.Cards)) / 52.0
	if remainingDecks == 0 {
		return 0
	}
	return float64(s.Deck.RunningCount) / remainingDecks
}

func LoadCountingStrategyFromData(name string, data CountingStrategyFile) (*CountingStrategy, error) {
	base := &DynamicStrategy{
		Fallback:        data.Fallback,
		Actions:         data.Actions,
		AcceptInsurance: data.AcceptInsurance, // bu alan optional ama sorun yaratmaz
	}

	return &CountingStrategy{
		BaseStrategy:    base,
		Deviations:      data.Deviations,
		BetRamp:         data.BetRamp,
		Deck:            nil,
		CountingEnabled: data.CountingEnabled,
		Name:            name,
	}, nil
}
