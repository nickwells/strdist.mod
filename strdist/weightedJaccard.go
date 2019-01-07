package strdist

import (
	"fmt"
	"strings"
)

const DfltWeightedJaccardThreshold = 0.33

// DfltWeightedJaccardFinder is a Finder with some default values suitable
// for a WeightedJaccard algorithm already set.
var DfltWeightedJaccardFinder *Finder

// CaseBlindWeightedJaccardFinder is a Finder with some default values suitable
// for a WeightedJaccard algorithm already set. CaseMod is set to ForceToLower.
var CaseBlindWeightedJaccardFinder *Finder

func init() {
	var err error
	DfltWeightedJaccardFinder, err =
		NewWeightedJaccardFinder(2,
			DfltMinStrLen, DfltWeightedJaccardThreshold, NoCaseChange)
	if err != nil {
		panic("Cannot construct the default WeightedJaccardFinder: " +
			err.Error())
	}
	CaseBlindWeightedJaccardFinder, err =
		NewWeightedJaccardFinder(2, DfltMinStrLen,
			DfltWeightedJaccardThreshold, ForceToLower)
	if err != nil {
		panic("Cannot construct the case-blind WeightedJaccardFinder: " +
			err.Error())
	}
}

// WeightedJaccardAlgo encapsulates the details needed to provide the
// WeightedJaccard distance.
type WeightedJaccardAlgo struct {
	N         int
	ngsTarget NGramSet
}

// NewWeightedJaccardFinder returns a new Finder having a WeightedJaccard
// algo and an error which will be non-nil if the parameters are invalid. The
// n-gram length must be > 0; for other invalid parameters see the NewFinder
// func.
func NewWeightedJaccardFinder(ngLen, minStrLen int, threshold float64, cm CaseMod) (*Finder, error) {
	if ngLen <= 0 {
		return nil,
			fmt.Errorf("bad N-Gram length (%d) - it should be > 0", ngLen)
	}
	algo := &WeightedJaccardAlgo{
		N: ngLen,
	}

	return NewFinder(minStrLen, threshold, cm, algo)
}

// Prep for a WeightedJaccardAlgo will pre-calculate the n-gram set for the
// target string
func (a *WeightedJaccardAlgo) Prep(s string, cm CaseMod) {
	switch cm {
	case ForceToLower:
		a.ngsTarget, _ = NGrams(strings.ToLower(s), a.N)
	default:
		a.ngsTarget, _ = NGrams(s, a.N)
	}
}

// Dist for a WeightedJaccardAlgo will calculate the distance from the target
// string
func (a *WeightedJaccardAlgo) Dist(_, s string, cm CaseMod) float64 {
	var ngs NGramSet
	switch cm {
	case ForceToLower:
		ngs, _ = NGrams(strings.ToLower(s), a.N)
	default:
		ngs, _ = NGrams(s, a.N)
	}
	return 1.0 - WeightedJaccardIndex(a.ngsTarget, ngs)
}

// WeightedJaccardIndex returns the Weighted Jaccard index of the two n-gram
// sets. It uses the NGramWeightedLen... functions to calculate the length
func WeightedJaccardIndex(ngs1, ngs2 NGramSet) float64 {
	if len(ngs1) == 0 && len(ngs2) == 0 {
		return 1.0
	}

	uLen := NGramWeightedLenUnion(ngs1, ngs2)
	iLen := NGramWeightedLenIntersection(ngs1, ngs2)

	return float64(iLen) / float64(uLen)
}

// WeightedJaccardDistance returns the Weighted Jaccard distance of the two
// strings. This is 1 minus the WeightedJaccardIndex
func WeightedJaccardDistance(s1, s2 string, n int) (float64, error) {
	ngs1, err := NGrams(s1, n)
	if err != nil {
		return 1.0, err
	}
	ngs2, err := NGrams(s2, n)
	if err != nil {
		return 1.0, err
	}

	return 1.0 - WeightedJaccardIndex(ngs1, ngs2), nil
}
