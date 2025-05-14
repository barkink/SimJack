package engine

import (
	"fmt"
	"math/rand"
	"strings"
)

type Card struct {
	Rank string
	Suit string
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

func (c Card) Value() int {
	switch c.Rank {
	case "A":
		return 11
	case "K", "Q", "J":
		return 10
	default:
		var v int
		fmt.Sscanf(c.Rank, "%d", &v)
		return v
	}
}

var Suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Deck struct {
	Cards                []Card
	DrawnThisRound       int
	DrawnThisShoe        int
	CutCardPosition      int
	NeedsNewDeck         bool
	NumDecks             int
	ForcedCards          []Card
	RunningCount         int
	RealCountTillCutCard int
}

func NewDeck(numDecks int, forced []Card) *Deck {
	d := &Deck{
		NumDecks:    numDecks,
		ForcedCards: forced,
	}
	d.SetupShoe()
	return d
}

func (d *Deck) SetupShoe() {
	full := []Card{}
	for i := 0; i < d.NumDecks; i++ {
		for _, suit := range Suits {
			for _, rank := range Ranks {
				full = append(full, Card{Rank: rank, Suit: suit})
			}
		}
	}

	rand.Shuffle(len(full), func(i, j int) {
		full[i], full[j] = full[j], full[i]
	})

	if len(d.ForcedCards) > 0 {
		// Say: destede kaç tane her karttan var?
		available := make(map[string]int)
		for _, c := range full {
			key := c.String()
			available[key]++
		}

		// Say: kullanıcı kaç adet zorunlu kart istemiş
		forcedCounts := make(map[string]int)
		for _, fc := range d.ForcedCards {
			key := fc.String()
			forcedCounts[key]++
		}

		// Validasyon: fazla istek varsa hata ver
		for key, count := range forcedCounts {
			if count > available[key] {
				panic(fmt.Sprintf("forced card '%s' exceeds available copies in deck (%d > %d)", key, count, available[key]))
			}
		}

		// Forced'ları çıkararak kalan desteyi oluştur
		remaining := []Card{}
		for _, c := range full {
			key := c.String()
			if forcedCounts[key] > 0 {
				forcedCounts[key]--
			} else {
				remaining = append(remaining, c)
			}
		}

		rand.Shuffle(len(remaining), func(i, j int) {
			remaining[i], remaining[j] = remaining[j], remaining[i]
		})

		d.Cards = append([]Card{}, d.ForcedCards...)
		d.Cards = append(d.Cards, remaining...)
	} else {
		d.Cards = full
	}

	minCut := int(float64(len(d.Cards)) * 0.5)
	maxCut := int(float64(len(d.Cards)) * 0.6)
	d.CutCardPosition = rand.Intn(maxCut-minCut) + minCut
	d.NeedsNewDeck = false
	d.DrawnThisShoe = 0
	d.RunningCount = 0

	d.RealCountTillCutCard = 0
	for i := 0; i < d.CutCardPosition; i++ {
		d.adjustRealCountTillCutCard(d.Cards[i])
	}
}

func (d *Deck) DealCard() (Card, error) {
	if len(d.Cards) == 0 {
		return Card{}, fmt.Errorf("no cards left in deck")
	}
	d.DrawnThisRound++
	d.DrawnThisShoe++
	if len(d.Cards) <= d.CutCardPosition {
		d.NeedsNewDeck = true
	}
	c := d.Cards[0]
	d.Cards = d.Cards[1:]
	d.adjustRunningCount(c)
	d.adjustRealCountTillCutCard(c)
	return c, nil
}

func (d *Deck) ShuffleIfNeeded() bool {
	if d.NeedsNewDeck {
		d.SetupShoe()
		return true
	}
	return false
}

func (d *Deck) ResetRoundCounter() {
	d.DrawnThisRound = 0
}

func ParseForcedCards(raw []string) []Card {
	var cards []Card
	for _, s := range raw {
		var rank, suit string
		_, err := fmt.Sscanf(s, "%s of %s", &rank, &suit)
		if err == nil {
			cards = append(cards, Card{Rank: rank, Suit: suit})
		}
	}
	return cards
}

func getHiLoValue(c Card) int {
	switch strings.ToUpper(c.Rank) {
	case "2", "3", "4", "5", "6":
		return 1
	case "10", "J", "Q", "K", "A":
		return -1
	default:
		return 0
	}
}

func (d *Deck) adjustRunningCount(c Card) {
	d.RunningCount += getHiLoValue(c)
}

func (d *Deck) adjustRealCountTillCutCard(c Card) {
	d.RealCountTillCutCard += getHiLoValue(c)
}

func (d *Deck) GetRunningCount() int {
	return d.RunningCount
}
