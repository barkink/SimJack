package engine

import (
	"simjack/config"
)

type Player struct {
	ID              int
	Owner           string
	Balance         float64
	InitialBalance  float64
	TargetBalance   float64
	IsBusted        bool
	IsRetired       bool
	Boxes           []*Box
	Strategy        Strategy
	RoundStartBal   float64
	TotalSpent      float64
	TotalEarned     float64
	BustedAtRound  int
	RetiredAtRound int
}

func NewPlayer(cfg config.PlayerConfig, strategy Strategy) *Player {
	// Strategy'ye sigorta davranışını aktar
	if cs, ok := strategy.(*CountingStrategy); ok {
		cs.AcceptInsurance = cfg.AcceptInsurance
	}

	return &Player{
		ID:              cfg.PlayerID,
		Owner:           cfg.Owner,
		Balance:         cfg.InitialBalance,
		InitialBalance:  cfg.InitialBalance,
		TargetBalance:   cfg.TargetBalance,
		Strategy:        strategy,
	}
}

func (p *Player) CanBet(amount float64) bool {
	return p.Balance >= amount
}

func (p *Player) PlaceBet(amount float64) bool {
	if p.CanBet(amount) {
		p.Balance -= amount
		p.TotalSpent += amount
		return true
	}
	return false
}

func (p *Player) ReceivePayout(amount float64) {
	p.Balance += amount
	p.TotalEarned += amount
}

func (p *Player) ResetRound() {
	p.RoundStartBal = p.Balance
	p.TotalSpent = 0
	p.TotalEarned = 0
}

func (p *Player) CheckStatus(minBet float64, round int) {
	if p.TargetBalance > 0 && p.Balance >= p.TargetBalance && !p.IsRetired {
		p.IsRetired = true
		p.RetiredAtRound = round
	}

	if p.Balance < minBet && !p.IsBusted {
		p.IsBusted = true
		p.BustedAtRound = round
	}
}

