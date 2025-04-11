package utils

import (
	"math"

	"github.com/ztrue/tracerr"
)

type generator struct {
	alphabet  string
	base      uint
	index     uint64
	length    uint
	maxLength uint
}

// NewGenerator - публичная функция для создания нового генератора
func NewGenerator(alphabet string, startIndex uint64, maxLength uint) *generator {
	base := len([]rune(alphabet))
	length := 1

	// Определяем минимальную длину, которая соответствует startIndex
	for startIndex >= uint64(math.Pow(float64(base), float64(length))) {
		startIndex -= uint64(math.Pow(float64(base), float64(length)))
		length++
	}

	return &generator{
		alphabet:  alphabet,
		base:      uint(base),
		index:     startIndex,
		length:    uint(length),
		maxLength: maxLength,
	}
}

// indexToWord - преобразует индекс в слово на основе алфавита и длины
func (g *generator) indexToWord(index uint64, length uint) string {
	alphabetRunes := []rune(g.alphabet)
	word := make([]rune, length)

	for i := int(length - 1); i >= 0; i-- {
		word[i] = alphabetRunes[index%uint64(g.base)]
		index /= uint64(g.base)
	}

	return string(word)
}

// Next - публичный метод для получения следующего слова
func (g *generator) Next() (string, error) {
	if g.length > g.maxLength {
		return "", tracerr.New("max length")
	}

	word := g.indexToWord(g.index, g.length)
	g.index++

	// Если достигли предела текущей длины, увеличиваем длину
	if g.index >= uint64(math.Pow(float64(g.base), float64(g.length))) {
		g.index = 0 // Начинаем сначала, но с увеличенной длиной
		g.length++
	}

	return word, nil
}
