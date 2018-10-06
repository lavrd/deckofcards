package deckofcards

import (
	"math/rand"
	"strings"
	"time"
)

const (
	DefaultRemaining = 52

	Discard = "discard"
)

var (
	Codes = []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "0S", "JS", "QS", "KS",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "0D", "JD", "QD", "KD",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "0H", "JH", "QH", "KH",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "0C", "JC", "QC", "KC",
	}

	Suits = map[string]string{
		"S": "SPADES",
		"D": "DIAMONDS",
		"H": "HEARTS",
		"C": "CLUBS",
	}

	Values = map[string]string{
		"A": "ACE",
		"K": "KING",
		"Q": "QUEEN",
		"J": "JACK",
		"0": "10",
		"9": "9",
		"8": "8",
		"7": "7",
		"6": "6",
		"5": "5",
		"4": "4",
		"3": "3",
		"2": "2",
	}
)

type Pile struct {
	Cards     Cards `json:"cards"`
	Remaining int   `json:"remaining"`
}
type Piles map[string]*Pile

type Card struct {
	SVG   string `json:"svg"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}
type Cards []*Card

type Deck struct {
	Shuffled  bool  `json:"shuffled"`
	Remaining int   `json:"remaining"`
	Cards     Cards `json:"cards"`
	Piles     Piles `json:"piles"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (d *Deck) Default() *Deck {
	var (
		cards Cards
	)

	for _, c := range Codes {
		code := c
		split := strings.Split(code, "")
		value := Values[split[0]]
		suit := Suits[split[1]]

		card := &Card{
			Value: value,
			Code:  code,
			Suit:  suit,
			SVG:   "",
		}

		cards = append(cards, card)
	}

	var (
		deck = &Deck{
			Cards:     cards,
			Remaining: DefaultRemaining,
			Piles: Piles{
				Discard: &Pile{},
			},
		}
	)

	d = deck
	return deck
}

func (d *Deck) Shuffle() *Deck {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})

	d.Shuffled = true
	return d
}

func (d *Deck) Partial(partials ...string) *Deck {
	var (
		cards Cards
	)

	for _, p := range partials {
		for _, c := range d.Cards {
			if p == c.Code {
				cards = append(cards, c)
			}
		}
	}

	d.Cards = cards
	d.Remaining = len(cards)
	return d
}

func (d *Deck) Pile(name string, cards Cards) *Pile {
	var (
		pile = &Pile{
			Cards:     cards,
			Remaining: len(cards),
		}
	)

	d.Piles[name] = pile
	return pile
}

func (p *Pile) Draw(count int) Cards {
	var (
		cards Cards
	)

	for range toRange(count) {
		i := rand.Intn(p.Remaining)
		cards = append(cards, p.Cards[i])
		p.Delete(i)
	}

	return cards
}

func (p *Pile) Delete(i int) *Pile {
	p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
	p.Remaining--
	return p
}

func New(shuffle bool) *Deck {
	var (
		deck = (&Deck{}).Default()
	)

	if shuffle {
		deck.Shuffle()
	}

	return deck
}

func toRange(n int) []struct{} {
	return make([]struct{}, n)
}

func (d *Deck) Delete(i int) *Deck {
	d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
	d.Remaining--
	return d
}

func (d *Deck) Draw(count int) Cards {
	var (
		cards Cards
	)

	if count > d.Remaining {
		count = d.Remaining
	}

	for range toRange(count) {
		i := rand.Intn(d.Remaining)
		cards = append(cards, d.Cards[i])
		d.Delete(i)
	}

	return cards
}
