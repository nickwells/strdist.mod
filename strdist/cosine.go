package strdist

import (
	"fmt"
	"math"
	"strings"
)

const DfltCosineThreshold = 0.33

// DfltCosineFinder is a Finder with some default values suitable
// for a Cosine algorithm already set.
var DfltCosineFinder *Finder

// CaseBlindCosineFinder is a Finder with some default values suitable
// for a Cosine algorithm already set. CaseMod is set to ForceToLower.
var CaseBlindCosineFinder *Finder

func init() {
	var err error
	DfltCosineFinder, err =
		NewCosineFinder(2, DfltMinStrLen, DfltCosineThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default CosineFinder: " + err.Error())
	}
	CaseBlindCosineFinder, err =
		NewCosineFinder(2, DfltMinStrLen, DfltCosineThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind CosineFinder: " + err.Error())
	}
}

// CosineAlgo encapsulates the details needed to provide the cosine distance.
type CosineAlgo struct {
	N         int
	ngsTarget NGramSet
	lenTarget float64
}

// NewCosineFinder returns a new Finder having a cosine algo and an error
// which will be non-nil if the parameters are invalid. The n-gram length
// must be > 0; for other invalid parameters see the NewFinder func.
func NewCosineFinder(ngLen, minStrLen int, threshold float64, cm CaseMod) (*Finder, error) {
	if ngLen <= 0 {
		return nil,
			fmt.Errorf("bad N-Gram length (%d) - it should be > 0", ngLen)
	}
	algo := &CosineAlgo{
		N: ngLen,
	}

	return NewFinder(minStrLen, threshold, cm, algo)
}

// Prep for a CosineAlgo will pre-calculate the n-gram set for the target string
func (a *CosineAlgo) Prep(s string, cm CaseMod) {
	switch cm {
	case ForceToLower:
		a.ngsTarget, _ = NGrams(strings.ToLower(s), a.N)
	default:
		a.ngsTarget, _ = NGrams(s, a.N)
	}
	a.lenTarget = a.ngsTarget.Length()
}

// Dist for a CosineAlgo will calculate the distance from the target string
func (a *CosineAlgo) Dist(_, s string, cm CaseMod) float64 {
	var ngs NGramSet
	switch cm {
	case ForceToLower:
		ngs, _ = NGrams(strings.ToLower(s), a.N)
	default:
		ngs, _ = NGrams(s, a.N)
	}
	return a.ngsTarget.cosineDistance(a.lenTarget, ngs)
}

// (ngs NGramSet)cosineDistance works out the cosine distance for a string
// having already worked out the N-Gram set for the other string and its
// distance
func (ngs NGramSet) cosineDistance(len float64, sNgs NGramSet) float64 {
	len2s := sNgs.lengthSquared()

	if len2s == 0.0 && len == 0.0 {
		return 0.0
	}
	if len2s == 0.0 || len == 0.0 {
		return 1.0
	}

	dot := Dot(ngs, sNgs)
	if dot == 0 {
		return 1.0
	}
	lenS := math.Sqrt(len2s)
	return 1.0 - (float64(dot) / (lenS * len))
}

// CosineSimilarity returns the cosine similarity between two NGramSets ngs1
// and ngs2. This is the dot-product of the two n-gram sets divided by
// the product of the lengths. Note that if both sets are of length 0 then
// the similarity is set to 1 (meaning identical) but if either one but not
// both is of length 0 (or the dot-product is zero) then the similarity is
// set to 0 (meaning completely different).
func CosineSimilarity(ngs1, ngs2 NGramSet) float64 {
	len2s1 := ngs1.lengthSquared()
	len2s2 := ngs2.lengthSquared()

	if len2s1 == 0.0 && len2s2 == 0.0 {
		return 1.0
	}
	if len2s1 == 0.0 || len2s2 == 0.0 {
		return 0.0
	}

	d := Dot(ngs1, ngs2)
	if d == 0.0 {
		return 0.0
	}

	lenS1 := math.Sqrt(len2s1)
	lenS2 := math.Sqrt(len2s2)

	return float64(d) / (lenS1 * lenS2)
}

// CosineDistance measures the cosine distance between two strings. This is 1
// minus the cosine similarity
func CosineDistance(s1, s2 string, n int) (float64, error) {
	ngs1, err := NGrams(s1, n)
	if err != nil {
		return 1.0, err
	}
	ngs2, err := NGrams(s2, n)
	if err != nil {
		return 1.0, err
	}

	return 1.0 - CosineSimilarity(ngs1, ngs2), nil
}
