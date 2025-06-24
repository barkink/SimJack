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
	AllowSurrender      bool
	SurrenderAgainstAce bool
	DealerTakesHoleCard bool
	Logger              *Logger
	MaxSplits           int
	ShowProgress bool
	lastPercent  int
	MinBet float64
	MaxBet float64
	MinSideBet float64
	MaxSideBet float64
	Debug bool
}

func NewEngine(cfg config.SimulationConfig, logger *Logger, showProgress bool, debug bool, stdinStrategies map[string]CountingStrategyFile) *Engine {
	players := []*Player{}
	boxes := make([]*Box, 7)
	deck := NewDeck(cfg.NumDecks, ParseForcedCards(cfg.ForcedCards))

	// Oyuncuları oluştur
	for _, pc := range cfg.Players {
		var strategy Strategy
		var err error

		if stdinStrategies != nil {
			data, ok := stdinStrategies[pc.Strategy]
			if !ok {
				fmt.Printf("Strategy %s not found in stdin input\n", pc.Strategy)
				os.Exit(1)
			}
			strategy, err = LoadCountingStrategyFromData(pc.Strategy, data)
			if err != nil {
				fmt.Printf("Failed to load strategy from stdin data: %v\n", err)
				os.Exit(1)
			}
		} else {
			strategy, err = LoadStrategyFromFile(pc.Strategy)
			if err != nil {
				fmt.Printf("Failed to load strategy from file: %v\n", err)
				os.Exit(1)
			}
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
		AllowSurrender:      cfg.AllowSurrender,
		SurrenderAgainstAce: cfg.SurrenderAgainstAce, 
		DealerTakesHoleCard: cfg.DealerTakesHoleCard,
		CurrentRound:        1,
		CurrentShoeNumber:   1,
		Logger:              logger,
		MaxSplits:           cfg.MaxSplits,
		ShowProgress:        showProgress,
		lastPercent:         -1,
		MinBet: 			 cfg.MinBet,
		MaxBet: 			 cfg.MaxBet,
		MinSideBet:          cfg.MinBet / 5,
		MaxSideBet:          cfg.MaxBet / 5,
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

		// Stratejinin önerdiği bahsi masanın limitleri (MinBet/MaxBet) içinde kalacak şekilde ayarla.
		if box.MainBet > e.MaxBet {
			box.MainBet = e.MaxBet
		}
		if box.MainBet < e.MinBet {
			box.MainBet = e.MinBet
		}

		// Yan bahisleri masanın yan bahis limitleri içinde kalacak şekilde ayarla.
		if box.PerfectPairBet > 0 { // Sadece pozitif bir bahis varsa kontrol et
			if box.PerfectPairBet < e.MinSideBet {
				box.PerfectPairBet = e.MinSideBet
			}
			if box.PerfectPairBet > e.MaxSideBet {
				box.PerfectPairBet = e.MaxSideBet
			}
		}
		if box.P21Bet > 0 { // Sadece pozitif bir bahis varsa kontrol et
			if box.P21Bet < e.MinSideBet {
				box.P21Bet = e.MinSideBet
			}
			if box.P21Bet > e.MaxSideBet {
				box.P21Bet = e.MaxSideBet
			}
		}

		if !e.determineAndPlaceBets(box) {
			continue // Oyuncu bahis yapamadı, sıradaki box'a geç.
		}

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
	anyInsuranceTaken := false
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
					anyInsuranceTaken = true
				}
			}
		}
	}

	// Dealer hole card aldıysa ve blackjack yaptıysa kontrol et
	if dealerHasAce && e.DealerTakesHoleCard && e.Dealer.Hand.IsBlackjack() {
		e.handleDealerBlackjackInsurance()
		e.handleRoundEnd() // round sonunda tüm box'ları işleyelim
		return             // round burada biter
	}

	// Dealer hole card aldıysa ve blackjack yapmadıysa, sigorta yapanlar kaybeder
	if dealerHasAce && e.DealerTakesHoleCard && !e.Dealer.Hand.IsBlackjack() {
		e.handleNoDealerBlackjackInsurance()
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
		if !e.DealerTakesHoleCard {
			// Dealer ikinci kartı almamışsa şimdi alır
			dc2, _ := e.Deck.DealCard()
			e.Dealer.Hand.AddCard(dc2)

			// Şimdi eline bakalım, blackjack mi?
			if e.Dealer.Hand.IsBlackjack() {
				e.handleDealerBlackjackInsurance()
				e.handleRoundEnd() // tüm box'ları topluca işleyelim
				return             // round burada biter
			} else {
				// Blackjack değilse sigorta yapanlar kaybeder
				e.handleNoDealerBlackjackInsurance()
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
	} else if anyInsuranceTaken && !e.DealerTakesHoleCard {  // Eğer sigorta alındı ve dağıtıcının hole card almadıysa, dealer'ın hole card'ını al
		dc2, _ := e.Deck.DealCard()
		e.Dealer.Hand.AddCard(dc2)
		if e.Dealer.Hand.IsBlackjack() {
			e.handleDealerBlackjackInsurance()
		} else {
			e.handleNoDealerBlackjackInsurance()
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

	handLoop:
		for i < len(box.Hands) {
			hand := box.Hands[i]

			if hand.IsSplitChild && hand.Cards[0].Rank == "A" {
				i++
				continue handLoop
			}

			if hand.CalculateValue() >= 21 {
				i++
				continue handLoop
			}

			actions, isFallback, isDeviation, key := p.Strategy.GetAction(hand, e.Dealer.Hand.Cards[0])
			hand.SetDecisionTrace(actions) // Önerilen tüm eylemleri geçici olarak sakla

			actionLoop:
				for _, action := range actions {
					
					finalizeAndLog := func(finalAction string) {
						hand.FinalizeDecision(key, finalAction, isDeviation, isFallback)
					}

					switch action {
					case "surrender":
						// YENİ KONTROL: Dağıtıcının açık kartı As ise, kuralın izin verip vermediğine bak.
						dealerHasAce := e.Dealer.Hand.Cards[0].Rank == "A"
						if dealerHasAce && !e.SurrenderAgainstAce {
							continue actionLoop // İzin yoksa, bu eylemi atla ve bir sonrakini dene.
						}
						// Surrender sadece ilk iki kartla ve kural izin veriyorsa mümkündür.
						if e.AllowSurrender && len(hand.Cards) == 2  && !hand.IsSplitChild {
							finalizeAndLog("surrender")
							hand.Result = "surrender" // Elin sonucunu ayarla
							i++                       // Sıradaki ele geç
							continue handLoop
						}
						// Koşullar sağlanmıyorsa, bir sonraki strateji eylemini dene.
						continue actionLoop

					case "split":
						if hand.CanSplit() && len(box.Hands) < e.MaxSplits+1 && p.PlaceBet(hand.BetAmount) {
							finalizeAndLog("split")
							c1 := hand.Cards[0]
							c2 := hand.Cards[1]

							h1 := NewSplitHand(hand, box.nextHandID)
							h1.AddCard(c1)
							card1, _ := e.Deck.DealCard()
							h1.AddCard(card1)
							box.nextHandID++
							h1.IsSplitChild = true

							h2 := NewSplitHand(hand, box.nextHandID)
							h2.AddCard(c2)
							card2, _ := e.Deck.DealCard()
							h2.AddCard(card2)
							box.nextHandID++
							h2.IsSplitChild = true
							
							box.SplitCount++
							box.Hands[i] = h1
							box.Hands = append(box.Hands[:i+1], append([]*Hand{h2}, box.Hands[i+1:]...)...)

							continue handLoop
						}
						continue actionLoop 

					case "double":
						if !(hand.IsSplitChild && !e.AllowDAS) && p.PlaceBet(hand.BetAmount) {
							finalizeAndLog("double")
							hand.MarkAsDoubled()
							card, _ := e.Deck.DealCard()
							hand.AddCard(card)
							i++
							continue handLoop
						}
						continue actionLoop

					case "hit":
						finalizeAndLog("hit")
						card, _ := e.Deck.DealCard()
						hand.AddCard(card)
						continue handLoop

					case "stand":
						finalizeAndLog("stand")
						i++
						continue handLoop
					}
				}
			// Eğer actionLoop'tan çıkıldıysa (hiçbir eylem başarılı olmadı), bu stand anlamına gelir.
			hand.FinalizeDecision(key, "stand", isDeviation, isFallback)
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
			if hand.Result == "surrender" {
				hand.Payout = hand.BetAmount / 2
			} else {
				result := e.Dealer.Evaluate(hand)
				if hand.Result == "" {
					hand.Result = result 
				}

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

// Bu fonksiyon, dağıtıcının Blackjack yaptığı durumu ele alır.
// Sigorta bahislerini "kazandı" olarak sonuçlandırır ve oyuncu ellerini
// dağıtıcının Blackjack'ine göre (push veya lose) ayarlar.
func (e *Engine) handleDealerBlackjackInsurance() {
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
				hand.FinalizeDecision("No Decision", "Player Blackjack", false, false)
				hand.Result = "push"
			} else {
				hand.FinalizeDecision("No Decision", "Dealer Blackjack", false, false)
				hand.Result = "lose"
			}
		}
	}
}

// Bu fonksiyon, dağıtıcının Blackjack YAPMADIĞI durumu ele alır.
// Sadece sigorta bahislerini "kaybetti" olarak sonuçlandırır.
func (e *Engine) handleNoDealerBlackjackInsurance() {
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

// determineAndPlaceBets, bir box için tüm bahis mantığını yönetir.
// Oyuncunun bakiyesine göre mümkün olan en iyi bahis kombinasyonunu belirler,
// bahisleri oyuncudan tahsil eder ve box'ı günceller.
// Bahis başarılıysa true, değilse false döner.
func (e *Engine) determineAndPlaceBets(box *Box) bool {
	p := box.Player
	mainBet := box.MainBet
	ppBet := box.PerfectPairBet
	p21Bet := box.P21Bet

	finalMainBet := 0.0
	finalPPBet := 0.0
	finalP21Bet := 0.0

	// --- 1. Aşama: Stratejinin önerdiği ana bahisle dene ---
	if ppBet > 0 && p21Bet > 0 && p.CanBet(mainBet+ppBet+p21Bet) {
		finalMainBet = mainBet
		finalPPBet = ppBet
		finalP21Bet = p21Bet
	} else if p21Bet > 0 && p.CanBet(mainBet+p21Bet) {
		finalMainBet = mainBet
		finalPPBet = ppBet
	} else if ppBet > 0 && p.CanBet(mainBet+ppBet) {
		finalMainBet = mainBet
		finalP21Bet = p21Bet
	} else if p.CanBet(mainBet) {
		finalMainBet = mainBet
		// --- 2. Aşama: Strateji bahsi yetmiyorsa, Minimum bahisle dene ---
	} else if ppBet > 0 && p21Bet > 0 && p.CanBet(e.MinBet+ppBet+p21Bet) {
		finalMainBet = e.MinBet
		finalPPBet = ppBet
		finalP21Bet = p21Bet
	} else if p21Bet > 0 && p.CanBet(e.MinBet+p21Bet) {
		finalMainBet = e.MinBet
		finalPPBet = ppBet
	} else if ppBet > 0 && p.CanBet(e.MinBet+ppBet) {
		finalMainBet = e.MinBet
		finalP21Bet = p21Bet
	} else if p.CanBet(e.MinBet) {
		finalMainBet = e.MinBet
	}

	if finalMainBet <= 0 {
		return false // Bahis yapılamadı.
	}

	// Belirlenen nihai bahisleri oyuncunun bakiyesinden düş.
	p.PlaceBet(finalMainBet)
	p.PlaceBet(finalPPBet)
	p.PlaceBet(finalP21Bet)

	// Box'ın durumunu nihai bahislerle güncelle.
	box.MainBet = finalMainBet
	box.PerfectPairBet = finalPPBet
	box.P21Bet = finalP21Bet

	return true // Bahis başarıyla yapıldı.
}