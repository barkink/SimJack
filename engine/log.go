package engine

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Logger struct {
	file        *os.File
	writer      *csv.Writer
	gzipWriter  *gzip.Writer
	finalPath   string
	tempPath    string
	flushEvery  int
	counter     int
	gzipEnabled bool
	FinalPath   string
	headerWritten bool
}

func NewLogger(path string, gzipEnabled bool) (*Logger, error) {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	tempPath := base + "_0" + ext
	finalPath := base + "_1" + ext

	if gzipEnabled {
		tempPath += ".gz"
		finalPath += ".gz"
	}

	f, err := os.Create(tempPath)
	if err != nil {
		return nil, err
	}

	var w *csv.Writer
	var gw *gzip.Writer

	if gzipEnabled {
		gw = gzip.NewWriter(f)
		w = csv.NewWriter(gw)
	} else {
		w = csv.NewWriter(f)
	}

	return &Logger{
		file:        f,
		writer:      w,
		gzipWriter:  gw,
		tempPath:    tempPath,
		finalPath:   finalPath,
		flushEvery:  10000,
		counter:     0,
		gzipEnabled: gzipEnabled,
		FinalPath:   finalPath,
		headerWritten: false,
	}, nil
}

func (l *Logger) writeHeader() {
	l.writer.Write([]string{
		"round", "shoe", "deck_running_count", "true_count", "real_count_till_cut_card", "box_id", "player_id", "hand_id", "owner", "strategy", 
		"bet_from_config", "bet_unit_used", "hand_payout", "main_payout", "pp_bet", "pp_win", "pp_type",
		"p21_bet", "p21_win", "p21_type", "insurance_taken", "insurance_bet", "insurance_payout", "insurance_result",
		"initial_balance", "round_start_balance", "player_balance",
		"hand", "result",
		"is_blackjack", "is_doubled", "is_split_child", "split_count",
		"dealer_upcard", "dealer_final_hand", "dealer_blackjack", "dealer_bust",
		"player_bust", "player_draws", "player_is_bankrupt", "player_is_retired",
		"num_decks", "cut_card_position", "cards_drawn_total", "cards_drawn_round", "cards_left_after_round", 
		"decision_trace",
		"box_total_invested","box_total_earned",
	})
	l.writer.Flush()
}

func (l *Logger) LogRound(round int, shoe int, box *Box, deck *Deck, dealer *Dealer) {
	if !l.headerWritten {
		l.writeHeader()
		l.headerWritten = true
	}
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

	for i, hand := range box.Hands {
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
		traceBytes, err := json.Marshal(hand.DecisionTrace)
		traceStr := "[]" // Varsayılan olarak boş JSON dizisi
		if err == nil {
			traceStr = string(traceBytes)
		}
		record := []string{
			strconv.Itoa(round),
			strconv.Itoa(shoe),
			strconv.Itoa(int(float64(deck.RunningCount))),
			strconv.Itoa(trueCount),
			strconv.Itoa(deck.RealCountTillCutCard),
			box.ID,
			strconv.Itoa(p.ID),
			hand.ID,
			p.Owner,
			p.Strategy.String(),
			fmt.Sprintf("%.2f", box.OriginalMainBet),
			fmt.Sprintf("%.2f", hand.BetAmount),
			fmt.Sprintf("%.2f", hand.Payout),
			fmt.Sprintf("%.2f", box.TotalPayout),
			fmt.Sprintf("%.2f", box.PerfectPairBet),
			fmt.Sprintf("%.2f", box.PerfectPairWin),
			box.PerfectPairType,
			fmt.Sprintf("%.2f", box.P21Bet),
			fmt.Sprintf("%.2f", box.P21Win),
			box.P21Type,
			boolToStr(box.InsuranceTaken),
			fmt.Sprintf("%.2f", box.InsuranceBet),
			fmt.Sprintf("%.2f", box.InsurancePayout),
			fmt.Sprintf("%s", box.InsuranceResult),
			fmt.Sprintf("%.2f", p.InitialBalance),
			fmt.Sprintf("%.2f", p.RoundStartBal),
			fmt.Sprintf("%.2f", p.Balance),
			hand.String(),
			hand.Result,
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
			boolToStr(p.IsBusted),
			boolToStr(p.IsRetired),
			strconv.Itoa(deck.NumDecks),
			strconv.Itoa(deck.CutCardPosition),
			strconv.Itoa(deck.DrawnThisShoe),
			strconv.Itoa(deck.DrawnThisRound),
			strconv.Itoa(len(deck.Cards)),
			traceStr,
		}

		if i == len(box.Hands)-1 {
			// Bu box için son eldeyiz

			// Total yatırım = ana bahislerin toplamı + sidebet
			totalBet := box.PerfectPairBet + box.P21Bet
			for _, h := range box.Hands {
				totalBet += h.BetAmount
			}

			// Total kazanç = yan bahis kazançları + ellerin sonucuna göre payout
			totalWin := box.PerfectPairWin + box.P21Win
			for _, h := range box.Hands {
				switch h.Result {
				case "win":
					totalWin += h.BetAmount * 2
				case "push":
					totalWin += h.BetAmount
				}
			}

			record = append(record, fmt.Sprintf("%.2f", totalBet)) // box_total_invested
			record = append(record, fmt.Sprintf("%.2f", totalWin)) // box_total_earned
		} else {
			record = append(record, "")
			record = append(record, "")
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
	if l.gzipEnabled && l.gzipWriter != nil {
		l.gzipWriter.Close()
	}
	l.file.Close()
	os.Rename(l.tempPath, l.finalPath)
}

func boolToStr(b bool) string {
	if b {
		return "True"
	}
	return "False"
}
