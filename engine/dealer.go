package engine

type Dealer struct {
	Hand *Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: NewHand(0, "dealer", -1), // -1: dealer box yok
	}
}

func (d *Dealer) ResetHand() {
	d.Hand = NewHand(0, "dealer", -1)
}

func (d *Dealer) Play(deck *Deck, hitOnSoft17 bool) {
	for {
		value := d.Hand.CalculateValue()
		if value > 17 {
			break
		}
		if value == 17 {
			if hitOnSoft17 && hasAce(d.Hand.Cards) {
				card, _ := deck.DealCard()
				d.Hand.AddCard(card)
				continue
			}
			break
		}
		card, _ := deck.DealCard()
		d.Hand.AddCard(card)
	}
}

func (d *Dealer) Evaluate(playerHand *Hand) string {
	dealerValue := d.Hand.CalculateValue()
	playerValue := playerHand.CalculateValue()

	if playerHand.IsBlackjack() && !d.Hand.IsBlackjack() {
		return "blackjack"
	}
	if playerHand.IsBust() {
		return "lose"
	}
	if d.Hand.IsBust() {
		return "win"
	}
	if playerValue > dealerValue {
		return "win"
	}
	if playerValue < dealerValue {
		return "lose"
	}
	return "push"
}

func hasAce(cards []Card) bool {
	for _, c := range cards {
		if c.Rank == "A" {
			return true
		}
	}
	return false
}
