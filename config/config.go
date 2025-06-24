package config

type SimulationConfig struct {
	NumDecks              int            `json:"num_decks"`
	RoundCount            int            `json:"round_count"`
	HitOnSoft17           bool           `json:"hit_on_soft_17"`
	AllowDoubleAfterSplit bool           `json:"allow_double_after_split"`
	DealerTakesHoleCard   bool           `json:"dealer_takes_hole_card"`
	ForcedCards           []string       `json:"forced_cards"`
	StrategyDirectory     string         `json:"strategy_directory"`
	GzipEnabled           bool           `json:"gzip_log"`
	MaxSplits             int            `json:"max_splits"`
	AllowSurrender        bool           `json:"allow_surrender"`
	SurrenderAgainstAce   bool           `json:"surrender_against_ace"` 
	MinBet                float64        `json:"min_bet"` 
	MaxBet                float64        `json:"max_bet"`
	Players               []PlayerConfig `json:"players"`
}

type PlayerConfig struct {
	PlayerID       int               `json:"player_id"`
	InitialBalance float64           `json:"initial_balance"`
	TargetBalance  float64           `json:"target_balance"`
	Strategy       string            `json:"strategy"`
	Owner          string            `json:"owner"`
	Boxes          []BoxAssignment   `json:"boxes"`
	AcceptInsurance  bool            `json:"accept_insurance"` 
}


// Her oyuncunun oturduğu box ve bahis ayarları
type BoxAssignment struct {
	Index    int               `json:"index"`
	MainBet  float64           `json:"main_bet"`
	Sidebets map[string]float64 `json:"sidebets"`
}
