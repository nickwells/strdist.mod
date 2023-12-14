package strdist

import (
	"math"
	"unicode/utf8"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// ScaledLevAlgo encapsulates the details needed to provide the ScaledLev
// distance.
type ScaledLevAlgo struct{}

// Name returns the algorithm name
func (ScaledLevAlgo) Name() string {
	return AlgoNameScaledLevenshtein
}

// Desc returns a string describing the algorithm configuration
func (ScaledLevAlgo) Desc() string {
	return ""
}

// Dist for a ScaledLevAlgo will calculate the ScaledLev distance between
// the two strings
func (ScaledLevAlgo) Dist(s1, s2 string) float64 {
	return ScaledLevDistance(s1, s2)
}

// ScaledLevDistance calculates the Scaled Levenshtein distance between
// strings a and b. This is the Levenshtein distance divided by the max of
// the lengths of the two strings. Two zero-length strings are taken as
// identical (with a zero distance between them)
func ScaledLevDistance(a, b string) float64 {
	aLen := utf8.RuneCountInString(a)
	bLen := utf8.RuneCountInString(b)

	if aLen == 0 && bLen == 0 {
		return 0.0
	}

	d := make([][]int, aLen+1)
	for i := range d {
		d[i] = make([]int, bLen+1)
		d[i][0] = i
	}

	for i := 1; i <= bLen; i++ {
		d[0][i] = i
	}

	for j, bRune := range b {
		for i, aRune := range a {
			var subsCost int
			if aRune != bRune {
				subsCost = 1
			}

			del := d[i][j+1] + 1
			ins := d[i+1][j] + 1
			sub := d[i][j] + subsCost

			d[i+1][j+1] = mathutil.MinOfInt(del, ins, sub)
		}
	}

	return float64(d[aLen][bLen]) / math.Max(float64(aLen), float64(bLen))
}
