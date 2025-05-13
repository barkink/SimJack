package config

type SimulationConfig struct {
	NumDecks              int            `json:"num_decks"`
	RoundCount            int            `json:"round_count"`
	HitOnSoft17           bool           `json:"hit_on_soft_17"`
	AllowDoubleAfterSplit bool           `json:"allow_double_after_split"`
	AcceptInsurance       bool           `json:"accept_insurance"`
	ForcedCards           []string       `json:"forced_cards"`
	StrategyDirectory     string         `json:"strategy_directory"`
	Players               []PlayerConfig `json:"players"`
	TrackScenario         string         `json:"track_scenario"`
	DealerTakesHoleCard   bool           `json:"dealer_takes_hole_card"`
}

type PlayerConfig struct {
	PlayerID       int                `json:"player_id"`
	InitialBalance float64            `json:"initial_balance"`
	BetUnit        float64            `json:"bet_unit"`
	MaxSplits      int                `json:"max_splits"`
	Strategy       string             `json:"strategy"`
	Owner          string             `json:"owner"`
	TargetBalance  float64            `json:"target_balance"`
	BoxIndexes     []int              `json:"box_indexes"`
	Sidebets       map[string]float64 `json:"sidebets"`
}
