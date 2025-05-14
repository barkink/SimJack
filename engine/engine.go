package engine

import (
	"fmt"
	"os"
	"simjack/config"
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
}

func NewEngine(cfg config.SimulationConfig, logger *Logger) *Engine {
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

		for _, boxIndex := range pc.BoxIndexes {
			if boxIndex < 1 || boxIndex > 7 {
				continue // Geçersiz box
			}
			idx := boxIndex - 1
			if boxes[idx] == nil {
				boxes[idx] = NewBox(boxIndex, p)
			} else {
				// aynı box'a birden fazla kişi oturamaz
				continue
			}
			p.Boxes = append(p.Boxes, boxes[idx])
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
	}
}

func (e *Engine) PlayRound() {
	e.Deck.ResetRoundCounter()
	e.Dealer.ResetHand()

	// Box içeriğini sıfırla
	for _, box := range e.Boxes {
		if box == nil || box.Player == nil || box.Player.IsBusted || box.Player.IsRetired {
			continue
		}
		p := box.Player
		p.ResetRound()
		box.Reset()

		if cs, ok := p.Strategy.(*CountingStrategy); ok {
			p.BetUnitUsed = cs.GetBetUnit(p.BetUnit)
		}
		if !p.PlaceMainBet() {
			//if !p.PlaceBet(p.BetUnit) {
			p.IsBusted = true
			continue
		}
		if pp := p.SidebetAmount("perfect_pair"); pp > 0 {
			if !p.PlaceBet(pp) {
				p.IsBusted = true
				continue
			}
		}
		if p21 := p.SidebetAmount("21+3"); p21 > 0 {
			if !p.PlaceBet(p21) {
				p.IsBusted = true
				continue
			}
		}

		hand := NewHand(p.BetUnit, box.ID, box.nextHandID)
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
				box.PerfectPairWin = win + box.PerfectPairBet
			} else {
				box.PerfectPairWin = 0
			}
			box.PerfectPairType = kind
		}

		if box.P21Bet > 0 && len(e.Dealer.Hand.Cards) > 0 {
			c3 := e.Dealer.Hand.Cards[0]
			win, kind := Get21Plus3Payout([]Card{c1, c2, c3})
			if win > 0 {
				box.P21Win = win + box.P21Bet
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
				amount := p.BetUnit / 2
				if p.PlaceInsurance(amount) {
					// sigorta başarıyla alındı
				}
			}
		}
	}

	// Dealer hole card aldıysa ve blackjack yaptıysa kontrol et
	if e.DealerTakesHoleCard && e.Dealer.Hand.IsBlackjack() {
		for _, box := range e.Boxes {
			if box == nil || box.Player == nil {
				continue
			}
			p := box.Player

			if p.InsuranceTaken {
				p.WinInsurance()
			} else {
				p.LoseInsurance()
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

	// Dealer hole card aldıysa ve blackjack yaptıysa kontrol et
	if e.DealerTakesHoleCard && e.Dealer.Hand.IsBlackjack() {
		for _, box := range e.Boxes {
			if box == nil || box.Player == nil {
				continue
			}
			p := box.Player

			if p.InsuranceTaken {
				p.WinInsurance()
			} else {
				p.LoseInsurance()
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
			p := box.Player
			if p.InsuranceTaken {
				p.LoseInsurance()
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
			c2, _ := e.Deck.DealCard()
			e.Dealer.Hand.AddCard(c2)

			// Şimdi eline bakalım, blackjack mi?
			if e.Dealer.Hand.IsBlackjack() {
				for _, box := range e.Boxes {
					if box == nil || box.Player == nil {
						continue
					}
					p := box.Player

					if p.InsuranceTaken {
						p.WinInsurance()
					} else {
						p.LoseInsurance()
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
					p := box.Player
					if p.InsuranceTaken {
						p.LoseInsurance()
					}
				}
			}
		}

		// Kurala göre devam et (soft 17 vs.)
		if e.HitOnSoft17 {
			e.Dealer.Play(e.Deck, true)
		} else {
			e.Dealer.Play(e.Deck, false)
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
				if p.PlaceBet(p.BetUnit) {
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
		box.TotalPayout += box.Player.InsurancePayout

		// Ödemeyi yap
		box.Player.ReceivePayout(box.TotalPayout)
		box.Player.CheckStatus()

		// ileride loglama vs. de buraya gelir
		if e.Logger != nil {
			e.Logger.LogRound(e.CurrentRound, e.CurrentShoeNumber, box, e.Deck, e.Dealer)
		}
	}
}
