package engine

import (
	"fmt"
	"simjack/config"
)

type Box struct {
	ID              string
	Player          *Player
	Hands           []*Hand
	MainBet         float64
	PerfectPairBet  float64
	PerfectPairWin  float64
	PerfectPairType string
	P21Bet          float64
	P21Win          float64
	P21Type         string
	TotalPayout     float64
	SplitCount      int
	nextHandID      int
	OriginalMainBet        float64
	OriginalPerfectPairBet float64
	OriginalP21Bet         float64
	InsuranceTaken  bool
	InsuranceBet    float64
	InsuranceResult string
	InsurancePayout float64
}

func (b *Box) AddHand(h *Hand) {
	b.Hands = append(b.Hands, h)
}

func (b *Box) Reset() {
	if b.Player != nil && b.Player.IsBusted {
		b.Player = nil
		b.Hands = nil
		return // Elenmiş oyuncunun Box'ı tamamen sıfırlanır
	}
	b.Hands = []*Hand{}
	b.PerfectPairWin = 0
	b.PerfectPairType = "none"
	b.P21Win = 0
	b.P21Type = "none"
	b.TotalPayout = 0
	b.SplitCount = 0
	b.nextHandID = 1
	b.MainBet = b.OriginalMainBet
	b.PerfectPairBet = b.OriginalPerfectPairBet
	b.P21Bet = b.OriginalP21Bet
	b.InsuranceTaken = false
	b.InsuranceBet = 0
	b.InsuranceResult = "none"
	b.InsurancePayout = 0
}

func NewBoxWithConfig(cfg config.BoxAssignment, player *Player) *Box {
	return &Box{
		ID:              fmt.Sprintf("B%d", cfg.Index),
		Player:          player,
		MainBet:         cfg.MainBet,
		PerfectPairBet:  cfg.Sidebets["perfect_pair"],
		P21Bet:          cfg.Sidebets["21+3"],
		OriginalMainBet:      cfg.MainBet,
		OriginalPerfectPairBet: cfg.Sidebets["perfect_pair"],
		OriginalP21Bet:         cfg.Sidebets["21+3"],
		PerfectPairType: "none",
		P21Type:         "none",
		Hands:           []*Hand{},
		nextHandID:      1,
		InsuranceTaken:  false,
		InsuranceBet:    0,
		InsuranceResult: "none",
		InsurancePayout: 0,
	}
}
