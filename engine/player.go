package engine

import "simjack/config"

type Player struct {
	ID              int
	Owner           string
	Balance         float64
	InitialBalance  float64
	BetUnit         float64
	BetUnitUsed     float64
	MaxSplits       int
	TargetBalance   float64
	IsBusted        bool
	IsRetired       bool
	Boxes           []*Box
	Strategy        Strategy
	Sidebets        map[string]float64
	InsuranceTaken  bool
	InsuranceBet    float64
	InsuranceResult string
	InsurancePayout float64
	RoundStartBal   float64
	TotalSpent      float64
	TotalEarned     float64
}

func NewPlayer(cfg config.PlayerConfig, strategy Strategy) *Player {
	return &Player{
		ID:              cfg.PlayerID,
		Owner:           cfg.Owner,
		Balance:         cfg.InitialBalance,
		InitialBalance:  cfg.InitialBalance,
		BetUnit:         cfg.BetUnit,
		TargetBalance:   cfg.TargetBalance,
		Strategy:        strategy,
		Sidebets:        cfg.Sidebets,
		InsuranceResult: "none",
	}
}

func (p *Player) CanBet(amount float64) bool {
	return p.Balance >= amount
}

func (p *Player) PlaceMainBet() bool {
	return p.PlaceBet(p.BetUnitUsed)
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
	p.InsuranceTaken = false
	p.InsuranceBet = 0
	p.InsuranceResult = "none"
	p.InsurancePayout = 0
	p.BetUnitUsed = p.BetUnit
}

func (p *Player) CheckStatus() {
	if p.Balance < p.BetUnit {
		p.IsBusted = true
	}
	if p.TargetBalance > 0 && p.Balance >= p.TargetBalance {
		p.IsRetired = true
	}
}

func (p *Player) PlaceInsurance(amount float64) bool {
	if p.PlaceBet(amount) {
		p.InsuranceTaken = true
		p.InsuranceBet = amount
		return true
	}
	return false
}

func (p *Player) WinInsurance() {
	p.InsuranceResult = "win"
	p.InsurancePayout = p.InsuranceBet * 2
}

func (p *Player) LoseInsurance() {
	p.InsuranceResult = "lose"
	p.InsurancePayout = 0
}

func (p *Player) StrategyName() string {
	return "basic"
}

func (p *Player) SidebetAmount(name string) float64 {
	if val, ok := p.Sidebets[name]; ok {
		return val
	}
	return 0
}
