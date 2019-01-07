package strdist

import (
	"fmt"
	"sort"
)

// CaseMod represents the different behaviours with regards to case
// handling when measuring distances
type CaseMod int

const (
	// NoCaseChange indicates that the case should not be changed
	NoCaseChange CaseMod = iota
	// ForceToLower indicates that the case should be forced to lower case
	// when calculating distances
	ForceToLower
)

// DistAlgo describes the algorithm which the Finder will use to
// distance. There is a Prep func provided which will allow some common tasks
// to be performed before the distance is calculated - some algorithms can
// cache some intermediate results to save time when calculating the
// string-to-string distance.
type DistAlgo interface {
	Prep(s string, cm CaseMod)
	Dist(s1, s2 string, cm CaseMod) float64
}

// DfltMinStrLen is a suggested minimum length of string to be matched. The
// problem with trying to find similar strings to very short targets is that
// they can match with a lot of not obviously similar alternatives. For
// instance a match for a single character string might be every other single
// character string in the population. For a number of use cases this is not
// particularly helpful.
const DfltMinStrLen = 4

// Finder records the parameters of the finding algorithm
type Finder struct {
	// MinStrLen records the minimum length of string to be matched
	MinStrLen int
	// T is the threshold for similarity for this finder
	T float64
	// CH, if set to ForceToLower, will convert all strings to lower case
	// before generating the distance
	CM CaseMod
	// Algo is the algorithm with which to calculate the distance between two
	// strings
	Algo DistAlgo
}

// NewFinder checks that the parameters are valid and creates a new
// Finder if they are. The minStrLen and threshold must each be >=
// 0. A zero threshold wil require an exact match.
func NewFinder(minStrLen int, threshold float64, cm CaseMod, a DistAlgo) (*Finder, error) {
	if minStrLen < 0 {
		return nil,
			fmt.Errorf("bad minimum string length (%d) - it should be >= 0",
				minStrLen)
	}
	if threshold < 0.0 {
		return nil,
			fmt.Errorf("bad threshold (%f) - it should be >= 0.0", threshold)
	}
	f := &Finder{
		MinStrLen: minStrLen,
		T:         threshold,
		CM:        cm,
		Algo:      a,
	}
	return f, nil
}

// FindLike returns StrDists for those strings in the population (pop) which
// are similar to the string (s). A string is similar if it has a common
// difference calculated from the n-grams which is less than or equal to the
// NGram finder's threshold value
func (f *Finder) FindLike(s string, pop ...string) []StrDist {
	lp := len(pop)
	if lp == 0 || len(s) < f.MinStrLen {
		return []StrDist{}
	}

	dists := make([]StrDist, 0, lp)

	f.Algo.Prep(s, f.CM)

	for _, p := range pop {
		if len(p) < f.MinStrLen {
			continue
		}

		d := f.Algo.Dist(s, p, f.CM)
		if d > f.T {
			continue
		}

		dists = append(dists, StrDist{
			Str:  p,
			Dist: d,
		})
	}

	sort.Slice(dists, func(i, j int) bool { return SDSlice(dists).Cmp(i, j) })
	return dists
}

// FindStrLike returns those strings in the population (pop) which are
// similar to the string (s). Similarity is as for the Find func.
func (f *Finder) FindStrLike(s string, pop ...string) []string {
	return convertStrDist(f.FindLike(s, pop...))
}

// FindNStrLike returns the first n strings in the population (pop) which are
// similar to the string (s). Similarity is as for the Find func.
func (f *Finder) FindNStrLike(n int, s string, pop ...string) []string {
	return convertStrDistN(n, f.FindLike(s, pop...))
}

// convertStrDist returns the strings from a slice of StrDists
func convertStrDist(dists []StrDist) []string {
	rval := make([]string, 0, len(dists))
	for _, d := range dists {
		rval = append(rval, d.Str)
	}

	return rval
}

// convertStrDistN returns the first n strings from a slice of StrDists
func convertStrDistN(n int, dists []StrDist) []string {
	if len(dists) < n {
		n = len(dists)
	}
	rval := make([]string, 0, n)
	for i, d := range dists {
		if i >= n {
			break
		}
		rval = append(rval, d.Str)
	}

	return rval
}
