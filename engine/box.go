package engine

type Box struct {
	ID              int
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
}

func NewBox(id int, player *Player) *Box {
	box := &Box{
		ID:              id,
		Player:          player,
		MainBet:         player.BetUnit,
		PerfectPairBet:  player.SidebetAmount("perfect_pair"),
		P21Bet:          player.SidebetAmount("21+3"),
		PerfectPairType: "none",
		P21Type:         "none",
		Hands:           []*Hand{},
		nextHandID:      1,
	}
	registerBox(box)
	return box
}

func (b *Box) AddHand(h *Hand) {
	b.Hands = append(b.Hands, h)
}

func (b *Box) Reset() {
	b.Hands = []*Hand{}
	b.PerfectPairWin = 0
	b.PerfectPairType = "none"
	b.P21Win = 0
	b.P21Type = "none"
	b.TotalPayout = 0
	b.SplitCount = 0
	b.nextHandID = 1
}
