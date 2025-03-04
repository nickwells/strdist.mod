package strdist

import (
	"sort"
	"strings"
)

// Finder records the parameters of the finding algorithm
type Finder struct {
	// FinderConfig holds the Finder configuration - various details about the
	// configuration of the underlying algorithm and constraints on the
	// behaviour of the Finder itself.
	FinderConfig
	// Algo is the algorithm with which to calculate the distance between two
	// strings
	Algo Algo
}

// NewFinder checks that the parameters are valid and creates a new
// Finder if they are. The minLen and threshold limit must each be >=
// 0. A zero threshold wil require an exact match.
func NewFinder(fc FinderConfig, algo Algo) (*Finder, error) {
	if err := fc.Check(); err != nil {
		return nil, err
	}

	return &Finder{
		FinderConfig: fc,
		Algo:         algo,
	}, nil
}

// NewFinderOrPanic returns a new Finder. It will panic if there is anything
// wrong with the config.
func NewFinderOrPanic(fc FinderConfig, algo Algo) *Finder {
	f, err := NewFinder(fc, algo)
	if err != nil {
		panic(err)
	}

	return f
}

// prepStr converts the string according to the FinderConfig
func (fc FinderConfig) prepStr(s string) string {
	if fc.MapToLowerCase {
		s = strings.ToLower(s)
	}

	if fc.StripRunes != "" {
		stripped := make([]rune, 0, len(s))

		for _, r := range s {
			if strings.ContainsRune(fc.StripRunes, r) {
				continue
			}

			stripped = append(stripped, r)
		}

		s = string(stripped)
	}

	return s
}

// FindLike returns StrDists for those strings in the population (pop) which
// are similar to the string (s). A string is similar if it has a common
// difference calculated from the n-grams which is less than or equal to the
// NGram finder's threshold value.
func (f *Finder) FindLike(s string, pop ...string) []StrDist {
	if len(pop) == 0 {
		return nil
	}

	s = f.FinderConfig.prepStr(s)
	if len(s) < f.FinderConfig.MinStrLength {
		return nil
	}

	dists := make([]StrDist, 0, len(pop))

	for _, pOrig := range pop {
		p := f.FinderConfig.prepStr(pOrig)

		if len(p) < f.FinderConfig.MinStrLength {
			continue
		}

		d := f.Algo.Dist(s, p)
		if d > f.FinderConfig.Threshold {
			continue
		}

		dists = append(dists, StrDist{
			Str:  pOrig,
			Dist: d,
		})
	}

	lt := lessThanFunc(len(s))

	sort.Slice(dists, func(i, j int) bool { return lt(dists[i], dists[j]) })

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
