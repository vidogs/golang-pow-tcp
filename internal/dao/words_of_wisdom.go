package dao

import (
	"github.com/sirupsen/logrus"
	"math/rand"
)

type WordsOfWisdom struct {
	log    *logrus.Entry
	quotes []string
}

func NewWordsOfWisdom(log *logrus.Entry, quotes []string) *WordsOfWisdom {
	return &WordsOfWisdom{
		log:    log,
		quotes: quotes,
	}
}

func (d *WordsOfWisdom) GetRandomQuote() (string, error) {
	if len(d.quotes) == 0 {
		return "", ErrNoRows
	}

	return d.quotes[rand.Intn(len(d.quotes))], nil
}
