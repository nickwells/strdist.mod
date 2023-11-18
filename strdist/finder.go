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

// DistAlgo describes the algorithm which the Finder will use to calculate
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
	// CM, if set to ForceToLower, will convert all strings to lower case
	// before generating the distance
	CM CaseMod
	// Algo is the algorithm with which to calculate the distance between two
	// strings
	Algo DistAlgo

	// pop holds the default population of strings for the Find... methods to
	// search if no strings are provided.
	pop []string
}

// NewFinder checks that the parameters are valid and creates a new
// Finder if they are. The minLen and threshold limit must each be >=
// 0. A zero threshold wil require an exact match.
func NewFinder(minLen int, limit float64, cm CaseMod, a DistAlgo) (*Finder, error) {
	if minLen < 0 {
		return nil,
			fmt.Errorf("bad minimum string length (%d) - it should be >= 0",
				minLen)
	}
	if limit < 0.0 {
		return nil,
			fmt.Errorf("bad threshold (%f) - it should be >= 0.0", limit)
	}
	f := &Finder{
		MinStrLen: minLen,
		T:         limit,
		CM:        cm,
		Algo:      a,
	}
	return f, nil
}

// SetPop will set the population of strings to be searched by the
// Find... methods
func (f *Finder) SetPop(pop []string) {
	f.pop = pop
}

// FindLike returns StrDists for those strings in the population (pop) which
// are similar to the string (s). A string is similar if it has a common
// difference calculated from the n-grams which is less than or equal to the
// NGram finder's threshold value. If the list of strings to search is empty
// then the default population from the Finder will be used. This should be
// set in advance using the SetPop method.
func (f *Finder) FindLike(s string, pop ...string) []StrDist {
	if len(pop) == 0 {
		pop = f.pop
	}
	if len(pop) == 0 || len(s) < f.MinStrLen {
		return []StrDist{}
	}

	dists := make([]StrDist, 0, len(pop))

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

	sort.Slice(dists, func(i, j int) bool {
		if dists[i].Dist != dists[j].Dist {
			return dists[i].Dist < dists[j].Dist
		}

		lenDiffI, lenDiffJ := len(dists[i].Str)-len(s), len(dists[j].Str)-len(s)
		sqLenDiffI, sqLenDiffJ := lenDiffI*lenDiffI, lenDiffJ*lenDiffJ
		if sqLenDiffI != sqLenDiffJ {
			return sqLenDiffI < sqLenDiffJ
		}

		return dists[i].Str < dists[j].Str
	})
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
