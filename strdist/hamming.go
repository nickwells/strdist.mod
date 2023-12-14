package strdist

import (
	"unicode/utf8"
)

// HammingAlgo encapsulates the details needed to provide the Hamming distance.
type HammingAlgo struct{}

// Name returns the algorithm name
func (HammingAlgo) Name() string {
	return AlgoNameHamming
}

// Desc returns a string describing the algorithm configuration
func (HammingAlgo) Desc() string {
	return ""
}

// HammingDistance returns the Hamming distance of the two strings. if the
// two strings are of different length then the Hamming distance is increased
// by the difference in lengths. Note that it compares runes rather than
// characters or chars
func (HammingAlgo) Dist(s1, s2 string) float64 {
	d := utf8.RuneCountInString(s2) - utf8.RuneCountInString(s1)
	if d < 0 {
		d *= -1
		s1, s2 = s2, s1 // s1 is longer than s2 so swap
	}

	r1, r2 := []rune(s1), []rune(s2)
	for i, s1r := range r1 {
		s2r := r2[i]
		if s1r != s2r {
			d++
		}
	}

	return float64(d)
}
