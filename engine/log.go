package engine

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Logger struct {
	file       *os.File
	writer     *csv.Writer
	finalPath  string
	tempPath   string
	flushEvery int
	counter    int
}

func NewLogger(path string) (*Logger, error) {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	tempPath := base + "_0" + ext
	finalPath := base + "_1" + ext

	f, err := os.Create(tempPath)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	w.Write([]string{
		"round", "shoe", "deck_running_count", "true_count", "real_count_till_cut_card", "box_id", "player_id", "hand_id", "owner", "hand", "result",
		"bet", "hand_payout", "main_payout", "pp_bet", "pp_win", "pp_type",
		"p21_bet", "p21_win", "p21_type", "insurance_bet", "insurance_payout",
		"initial_balance", "round_start_balance", "player_balance",
		"is_blackjack", "is_doubled", "is_split_child", "split_count",
		"dealer_upcard", "dealer_final_hand", "dealer_blackjack", "dealer_bust",
		"player_bust", "player_draws", "player_is_bankrupt", "player_is_retired",
		"num_decks", "cut_card_position", "cards_drawn_total", "cards_drawn_round", "cards_left_after_round", "strategy_key",
	})
	w.Flush()
	return &Logger{
		file: f, writer: w,
		tempPath:   tempPath,
		finalPath:  finalPath,
		flushEvery: 100000, // her 1000 kayıtta flush
		counter:    0,
	}, nil
}

func (l *Logger) LogRound(round int, shoe int, box *Box, deck *Deck, dealer *Dealer) {
	p := box.Player
	trueCount := 0
	if deck != nil && deck.NumDecks > 0 {
		remainingDecks := float64(len(deck.Cards)) / 52.0
		if remainingDecks > 0 {
			trueCount = int(float64(deck.RunningCount) / remainingDecks)
		}
	}
	dealerUp := "?"
	if len(dealer.Hand.Cards) > 0 {
		dealerUp = dealer.Hand.Cards[0].String()
	}
	dealerFinal := ""
	for i, c := range dealer.Hand.Cards {
		if i > 0 {
			dealerFinal += ";"
		}
		dealerFinal += c.String()
	}
	dealerBJ := boolToStr(dealer.Hand.IsBlackjack())
	dealerBust := boolToStr(dealer.Hand.CalculateValue() > 21)

	for _, hand := range box.Hands {
		draws := ""
		if len(hand.Cards) > 2 {
			for i, c := range hand.Cards[2:] {
				if i > 0 {
					draws += ";"
				}
				draws += c.String()
			}
		}

		//action Key hesaplama
		traceStr := "no_decision_needed"
		if len(hand.DecisionTrace) > 0 {
			traceStr = ""
			for i, entry := range hand.DecisionTrace {
				if i > 0 {
					traceStr += "; "
				}
				traceStr += fmt.Sprintf("%s:%s", entry.Key, entry.Action)
			}
		}
		record := []string{
			strconv.Itoa(round),
			strconv.Itoa(shoe),
			strconv.Itoa(int(float64(deck.RunningCount))),
			strconv.Itoa(trueCount),
			strconv.Itoa(deck.RealCountTillCutCard),
			strconv.Itoa(box.ID),
			strconv.Itoa(p.ID),
			strconv.Itoa(hand.ID),
			p.Owner,
			hand.String(),
			hand.Result,
			fmt.Sprintf("%.2f", hand.BetAmount),
			fmt.Sprintf("%.2f", hand.Payout),
			fmt.Sprintf("%.2f", box.TotalPayout),
			fmt.Sprintf("%.2f", box.PerfectPairBet),
			fmt.Sprintf("%.2f", box.PerfectPairWin),
			box.PerfectPairType,
			fmt.Sprintf("%.2f", box.P21Bet),
			fmt.Sprintf("%.2f", box.P21Win),
			box.P21Type,
			fmt.Sprintf("%.2f", p.InsuranceBet),
			fmt.Sprintf("%.2f", p.InsurancePayout),
			fmt.Sprintf("%.2f", p.InitialBalance),
			fmt.Sprintf("%.2f", p.RoundStartBal),
			fmt.Sprintf("%.2f", p.Balance),
			boolToStr(hand.IsBlackjack()),
			boolToStr(hand.IsDoubled),
			boolToStr(hand.IsSplitChild),
			strconv.Itoa(box.SplitCount),
			dealerUp,
			dealerFinal,
			dealerBJ,
			dealerBust,
			boolToStr(hand.CalculateValue() > 21),
			draws,
			boolToStr(p.Balance < p.BetUnit),
			boolToStr(p.TargetBalance > 0 && p.Balance >= p.TargetBalance),
			strconv.Itoa(deck.NumDecks),
			strconv.Itoa(deck.CutCardPosition),
			strconv.Itoa(deck.DrawnThisShoe),
			strconv.Itoa(deck.DrawnThisRound),
			strconv.Itoa(len(deck.Cards)),
			traceStr,
		}
		l.writer.Write(record)
		l.counter++
		if l.counter%l.flushEvery == 0 {
			l.writer.Flush()
		}
	}
}

func (l *Logger) Close() {
	l.writer.Flush()
	l.file.Close()
	os.Rename(l.tempPath, l.finalPath)
}

func boolToStr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}
