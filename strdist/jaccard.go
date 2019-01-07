package strdist

import (
	"fmt"
	"strings"
)

const DfltJaccardThreshold = 0.33

// DfltJaccardFinder is a Finder with some default values suitable
// for a Jaccard algorithm already set.
var DfltJaccardFinder *Finder

// CaseBlindJaccardFinder is a Finder with some default values suitable
// for a Jaccard algorithm already set. CaseMod is set to ForceToLower.
var CaseBlindJaccardFinder *Finder

func init() {
	var err error
	DfltJaccardFinder, err =
		NewJaccardFinder(2, DfltMinStrLen, DfltJaccardThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default JaccardFinder: " + err.Error())
	}
	CaseBlindJaccardFinder, err =
		NewJaccardFinder(2, DfltMinStrLen, DfltJaccardThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind JaccardFinder: " + err.Error())
	}
}

// JaccardAlgo encapsulates the details needed to provide the cosine distance.
type JaccardAlgo struct {
	N         int
	ngsTarget NGramSet
}

// NewJaccardFinder returns a new Finder having a Jaccard algo and an error
// which will be non-nil if the parameters are invalid. The n-gram length
// must be > 0; for other invalid parameters see the NewFinder func.
func NewJaccardFinder(ngLen, minStrLen int, threshold float64, cm CaseMod) (*Finder, error) {
	if ngLen <= 0 {
		return nil,
			fmt.Errorf("bad N-Gram length (%d) - it should be > 0", ngLen)
	}
	algo := &JaccardAlgo{
		N: ngLen,
	}

	return NewFinder(minStrLen, threshold, cm, algo)
}

// Prep for a JaccardAlgo will pre-calculate the n-gram set for the target
// string
func (a *JaccardAlgo) Prep(s string, cm CaseMod) {
	switch cm {
	case ForceToLower:
		a.ngsTarget, _ = NGrams(strings.ToLower(s), a.N)
	default:
		a.ngsTarget, _ = NGrams(s, a.N)
	}
}

// Dist for a JaccardAlgo will calculate the distance from the target string
func (a *JaccardAlgo) Dist(_, s string, cm CaseMod) float64 {
	var ngs NGramSet
	switch cm {
	case ForceToLower:
		ngs, _ = NGrams(strings.ToLower(s), a.N)
	default:
		ngs, _ = NGrams(s, a.N)
	}
	return 1.0 - JaccardIndex(a.ngsTarget, ngs)
}

// JaccardIndex returns the Jaccard index of the two n-gram sets
func JaccardIndex(ngs1, ngs2 NGramSet) float64 {
	if len(ngs1) == 0 && len(ngs2) == 0 {
		return 1.0
	}

	uLen := NGramLenUnion(ngs1, ngs2)
	iLen := NGramLenIntersection(ngs1, ngs2)

	return float64(iLen) / float64(uLen)
}

// JaccardDistance returns the Jaccard distance of the two strings. This is 1
// minus the JaccardIndex
func JaccardDistance(s1, s2 string, n int) (float64, error) {
	ngs1, err := NGrams(s1, n)
	if err != nil {
		return 1.0, err
	}
	ngs2, err := NGrams(s2, n)
	if err != nil {
		return 1.0, err
	}

	return 1.0 - JaccardIndex(ngs1, ngs2), nil
}
