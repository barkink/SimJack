package engine

import (
	"fmt"
	"os"
	"simjack/config"
	"strings" 
)

type Engine struct {
	Deck                *Deck
	Dealer              *Dealer
	Boxes               []*Box
	Players             []*Player
	RoundCount          int
	CurrentRound        int
	CurrentShoeNumber   int
	HitOnSoft17         bool
	AllowDAS            bool
	DealerTakesHoleCard bool
	Logger              *Logger
	MaxSplits           int
	ShowProgress bool
	lastPercent  int
	MinBet float64
	MaxBet float64
	Debug bool
}

func NewEngine(cfg config.SimulationConfig, logger *Logger, showProgress bool, debug bool) *Engine {
	players := []*Player{}
	boxes := make([]*Box, 7)
	deck := NewDeck(cfg.NumDecks, ParseForcedCards(cfg.ForcedCards))

	// Oyuncuları oluştur
	for _, pc := range cfg.Players {
		strategy, err := LoadStrategyFromFile(pc.Strategy)
		if err != nil {
			fmt.Printf("Failed to load strategy: %v\n", err)
			os.Exit(1)
		}

		if cs, ok := strategy.(*CountingStrategy); ok && cs.CountingEnabled {
			cs.Deck = deck
		}

		p := NewPlayer(pc, strategy)

		p.Boxes = []*Box{}
		players = append(players, p)

		for _, b := range pc.Boxes {
			if b.Index < 1 || b.Index > 7 {
				continue // geçersiz box
			}
			idx := b.Index - 1
			if boxes[idx] != nil {
				continue // aynı box'a iki kişi oturamaz
			}
			box := NewBoxWithConfig(b, p)
			boxes[idx] = box
			p.Boxes = append(p.Boxes, box)
		}
	}

	return &Engine{
		Deck:                deck,
		Dealer:              NewDealer(),
		Players:             players,
		Boxes:               boxes,
		RoundCount:          cfg.RoundCount,
		HitOnSoft17:         cfg.HitOnSoft17,
		AllowDAS:            cfg.AllowDoubleAfterSplit,
		DealerTakesHoleCard: cfg.DealerTakesHoleCard,
		CurrentRound:        1,
		CurrentShoeNumber:   1,
		Logger:              logger,
		MaxSplits:           cfg.MaxSplits,
		ShowProgress:        showProgress,
		lastPercent:         -1,
		MinBet: 			 cfg.MinBet,
		MaxBet: 			 cfg.MaxBet,
		Debug: 				 debug,
	}
}

func (e *Engine) Run() {
	for i := 0; i < e.RoundCount; i++ {
		active := false
		for _, p := range e.Players {
			if !p.IsBusted && !p.IsRetired {
				active = true
				break
			}
		}
		if !active {
			break
		}

		e.PlayRound()
		e.CurrentRound++

		if e.Deck.ShuffleIfNeeded() {
			e.CurrentShoeNumber++
		}

		if e.ShowProgress && e.RoundCount > 0 {
			percent := e.CurrentRound * 100 / e.RoundCount
			if percent > 100 {
				percent = 100
			}
			if percent != e.lastPercent {
				e.lastPercent = percent
				barLength := 20
				filled := percent * barLength / 100
				if filled > barLength {
					filled = barLength
				}
				if filled < 0 {
					filled = 0
				}
				empty := barLength - filled
				if empty < 0 {
					empty = 0
				}
				bar := "[" + strings.Repeat("█", filled) + strings.Repeat(" ", empty) + "]"
				fmt.Printf("\r%s %3d%% \t", bar, percent)
				if percent == 100 {
					fmt.Println()
				}
			}
		}

	}
}

func (e *Engine) PlayRound() {
	e.Deck.ResetRoundCounter()
	e.Dealer.ResetHand()

	for _, p := range e.Players {
		if !p.IsBusted && !p.IsRetired {
			p.ResetRound()  // ✔️ yalnızca bir kez! Oyunucun round bazlı başında ve sonundaki kasa bilgisi her hand için aynı olsun diye. 
		}
	}

	// Box içeriğini sıfırla
	for _, box := range e.Boxes {
		if box == nil || box.Player == nil || box.Player.IsBusted || box.Player.IsRetired {
			continue
		}
		p := box.Player
		box.Reset()

		if cs, ok := p.Strategy.(*CountingStrategy); ok {
			box.MainBet = cs.GetBetUnit(box.OriginalMainBet)
		}

		// Bahisleri planla
		mainBet := box.MainBet
		ppBet := box.PerfectPairBet
		p21Bet := box.P21Bet
		total := mainBet + ppBet + p21Bet

		// Oyuncunun toplam parası tümünü karşılıyorsa, hepsini yatır
		if p.CanBet(total) {
			p.PlaceBet(mainBet)
			p.PlaceBet(ppBet)
			p.PlaceBet(p21Bet)
		} else if p.CanBet(mainBet) {
			// Sadece main bet yatırılabiliyor, sidebet'ler yapılamaz
			p.PlaceBet(mainBet)
			box.PerfectPairBet = 0
			box.P21Bet = 0
		} else if p.CanBet(e.MinBet + ppBet + p21Bet) {
			// Main bet düşürülür, yan bahislerle beraber
			mainBet = e.MinBet + ppBet + p21Bet
			p.PlaceBet(mainBet)
			p.PlaceBet(ppBet)
			p.PlaceBet(p21Bet)
		}  else if p.CanBet(e.MinBet) {
			// Main bet düşürülür, yan bahislerle beraber 
			mainBet = e.MinBet
			p.PlaceBet(e.MinBet)
			box.PerfectPairBet = 0
			box.P21Bet = 0
		}else {
			// Oyuncunun hiçbir bahis yapacak kadar parası yok
			continue
		}

		// Box için güncellenen bet miktarını ayarla
		box.MainBet = mainBet

		hand := NewHand(box.MainBet, box.ID, box.nextHandID)
		box.AddHand(hand)
	}

	// İlk kart dağıtımı (her box'a)
	for _, box := range e.Boxes {
		if box == nil || len(box.Hands) == 0 {
			continue
		}
		card, _ := e.Deck.DealCard()
		box.Hands[0].AddCard(card)
	}

	// Dealer ilk kart
	dc1, _ := e.Deck.DealCard()
	e.Dealer.Hand.AddCard(dc1)

	// İkinci kart dağıtımı (oyunculara)
	for _, box := range e.Boxes {
		if box == nil || len(box.Hands) == 0 {
			continue
		}
		card, _ := e.Deck.DealCard()
		box.Hands[0].AddCard(card)
	}

	// Dealer ikinci kart opsiyonel
	if e.DealerTakesHoleCard {
		dc2, _ := e.Deck.DealCard()
		e.Dealer.Hand.AddCard(dc2)
	}

	// Yan bahisleri değerlendir
	for _, box := range e.Boxes {
		if box == nil || len(box.Hands) == 0 {
			continue
		}
		//p := box.Player
		c1 := box.Hands[0].Cards[0]
		c2 := box.Hands[0].Cards[1]

		if box.PerfectPairBet > 0 {
			win, kind := GetPerfectPairPayout(c1, c2)
			if win > 0 {
				box.PerfectPairWin = float64(win + 1) *  box.PerfectPairBet
			} else {
				box.PerfectPairWin = 0
			}
			box.PerfectPairType = kind
		}

		if box.P21Bet > 0 && len(e.Dealer.Hand.Cards) > 0 {
			c3 := e.Dealer.Hand.Cards[0]
			win, kind := Get21Plus3Payout([]Card{c1, c2, c3})
			if win > 0 {
				box.P21Win = float64(win + 1) *  box.P21Bet
			} else {
				box.P21Win = 0
			}
			box.P21Type = kind
		}
	}

	// Sigorta: Dealer açık kartı A ise sor
	dealerHasAce := e.Dealer.Hand.Cards[0].Rank == "A"
	if dealerHasAce {
		for _, box := range e.Boxes {
			if box == nil || box.Player == nil || box.Player.IsBusted || box.Player.IsRetired {
				continue
			}
			p := box.Player
			if p.Strategy.DecideInsurance() {
				amount := box.MainBet / 2
				if p.PlaceBet(amount) {
					// sigorta başarıyla alındı
					box.InsuranceTaken = true
					box.InsuranceBet = amount
				}
			}
		}
	}

	// Dealer hole card aldıysa ve blackjack yaptıysa kontrol et
	if dealerHasAce && e.DealerTakesHoleCard && e.Dealer.Hand.IsBlackjack() {
		for _, box := range e.Boxes {
			if box == nil || box.Player == nil {
				continue
			}

			if box.InsuranceTaken {
				box.InsuranceResult = "win"
				box.InsurancePayout = box.InsuranceBet * 2
			} else {
				box.InsuranceResult = "lose"
				box.InsurancePayout = 0
			}

			for _, hand := range box.Hands {
				if hand.IsBlackjack() {
					hand.Result = "push"
					box.TotalPayout += hand.BetAmount
				} else {
					hand.Result = "lose"
				}
			}
		}
		e.handleRoundEnd() // round sonunda tüm box'ları işleyelim
		return             // round burada biter
	}

	// Dealer hole card aldıysa ve blackjack yapmadıysa, sigorta yapanlar kaybeder
	if dealerHasAce && e.DealerTakesHoleCard && !e.Dealer.Hand.IsBlackjack() {
		for _, box := range e.Boxes {
			if box == nil || box.Player == nil {
				continue
			}

			if box.InsuranceTaken {
				box.InsuranceResult = "lose"
				box.InsurancePayout = 0
			}
		}
	}

	// Buradan sonra oyuncu aksiyonları başlar
	e.executeBoxActions()

	// Dealer oynayacak mı kontrol et
	dealerShouldPlay := false
	for _, box := range e.Boxes {
		if box == nil || box.Player == nil {
			continue
		}
		for _, hand := range box.Hands {
			if !hand.IsBust() && !hand.IsBlackjack() {
				dealerShouldPlay = true
				break
			}
		}
		if dealerShouldPlay {
			break
		}
	}

	if dealerShouldPlay {
		if len(e.Dealer.Hand.Cards) == 1 {
			// Dealer ikinci kartı almamışsa şimdi alır
			dc2, _ := e.Deck.DealCard()
			e.Dealer.Hand.AddCard(dc2)

			// Şimdi eline bakalım, blackjack mi?
			if e.Dealer.Hand.IsBlackjack() {
				for _, box := range e.Boxes {
					if box == nil || box.Player == nil {
						continue
					}

					if box.InsuranceTaken {
						box.InsuranceResult = "win"
						box.InsurancePayout = box.InsuranceBet * 2
					} else {
						box.InsuranceResult = "lose"
						box.InsurancePayout = 0
					}

					for _, hand := range box.Hands {
						if hand.IsBlackjack() {
							hand.Result = "push"
							box.TotalPayout += hand.BetAmount
						} else {
							hand.Result = "lose"
						}
					}
				}
				e.handleRoundEnd() // tüm box'ları topluca işleyelim
				return             // round burada biter
			} else {
				// Blackjack değilse sigorta yapanlar kaybeder
				for _, box := range e.Boxes {
					if box == nil || box.Player == nil {
						continue
					}

					if box.InsuranceTaken {
						box.InsuranceResult = "lose"
						box.InsurancePayout = 0
					}
				}
				// Kurala göre devam et (soft 17 vs.)
				if e.HitOnSoft17 {
					e.Dealer.Play(e.Deck, true)
				} else {
					e.Dealer.Play(e.Deck, false)
				}
			}
		} else {

			// Kurala göre devam et (soft 17 vs.)
			if e.HitOnSoft17 {
				e.Dealer.Play(e.Deck, true)
			} else {
				e.Dealer.Play(e.Deck, false)
			}
		}
	}

	e.handleRoundEnd() // round sonunda tüm box'ları işleyelim
}

func (e *Engine) executeBoxActions() {
	for _, box := range e.Boxes {
		if box == nil || box.Player == nil || len(box.Hands) == 0 {
			continue
		}
		p := box.Player
		i := 0
		for i < len(box.Hands) {
			hand := box.Hands[i]
			if hand.IsSplitChild && hand.Cards[0].Rank == "A" && len(hand.Cards) == 2 && hand.Cards[1].Rank != "A" {
				i++
				continue
			}

			action := p.Strategy.GetAction(hand, e.Dealer.Hand.Cards[0])

			if action == "split" && hand.CanSplit() && len(box.Hands) < e.MaxSplits+1 {
				if p.PlaceBet(hand.BetAmount) {
					c1 := hand.Cards[0]
					c2 := hand.Cards[1]

					h1 := NewSplitHand(hand, box.nextHandID)
					h1.AddCard(c1)
					card1, _ := e.Deck.DealCard()
					h1.AddCard(card1)
					box.nextHandID++

					h2 := NewSplitHand(hand, box.nextHandID)
					h2.AddCard(c2)
					card2, _ := e.Deck.DealCard()
					h2.AddCard(card2)
					box.nextHandID++

					h1.IsSplitChild = true
					h2.IsSplitChild = true
					box.SplitCount++

					box.Hands[i] = h1
					box.Hands = append(box.Hands[:i+1], append([]*Hand{h2}, box.Hands[i+1:]...)...)
					continue
				} else {
					action = p.Strategy.FallbackAction("split")
				}
			}

			if action == "double" {
				if hand.IsSplitChild && !e.AllowDAS {
					action = p.Strategy.FallbackAction("double")
				} else if p.PlaceBet(hand.BetAmount) {
					hand.MarkAsDoubled()
					card, _ := e.Deck.DealCard()
					hand.AddCard(card)
					i++
					continue
				} else {
					action = p.Strategy.FallbackAction("double")
				}
			}

			if action == "hit" {
				for hand.CalculateValue() < 21 {
					card, _ := e.Deck.DealCard()
					hand.AddCard(card)
					if hand.CalculateValue() >= 17 {
						break
					}
				}
			}

			i++
		}
	}
}

func (e *Engine) handleRoundEnd() {
	for _, box := range e.Boxes {
		if box == nil || box.Player == nil {
			continue
		}

		// Her hand için dealer'a karşı sonucu değerlendir
		for _, hand := range box.Hands {
			result := e.Dealer.Evaluate(hand)
			hand.Result = result

			switch result {
			case "win":
				hand.Payout += hand.BetAmount * 2
			case "blackjack":
				hand.Payout += hand.BetAmount * 2.5
			case "push":
				hand.Payout += hand.BetAmount
			case "lose":
				// kayıp, ödeme yok
			}
			box.TotalPayout += hand.Payout
		}

		// Yan bahisleri ekle
		box.TotalPayout += box.PerfectPairWin + box.P21Win
		// Sigorta ödemesi
		box.TotalPayout += box.InsurancePayout

		// Ödemeyi yap
		box.Player.ReceivePayout(box.TotalPayout)
	}

	for _, p := range e.Players {
		if !p.IsBusted && !p.IsRetired {
			p.CheckStatus(e.MinBet, e.CurrentRound)
		}
	}

	for _, box := range e.Boxes {
		if box == nil || box.Player == nil {
			continue
		}
		if e.Logger != nil {
			e.Logger.LogRound(e.CurrentRound, e.CurrentShoeNumber, box, e.Deck, e.Dealer)
		}
		box.Reset()
	}

	if e.Debug {
		for _, p := range e.Players {
			if p.BustedAtRound == e.CurrentRound {
				fmt.Printf("❌ Player %d (%s) BUSTED at round %d | Balance: %.2f\n",
					p.ID, p.Owner, p.BustedAtRound, p.Balance)
			}
			if p.RetiredAtRound == e.CurrentRound {
				fmt.Printf("🏁 Player %d (%s) RETIRED at round %d | Balance: %.2f\n",
					p.ID, p.Owner, p.RetiredAtRound, p.Balance)
			}
		}
	}

}
