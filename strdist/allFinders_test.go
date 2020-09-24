package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// finders holds all the finders
type finders struct {
	cosineFinder          *strdist.Finder
	jaccardFinder         *strdist.Finder
	weightedJaccardFinder *strdist.Finder
	levenshteinFinder     *strdist.Finder
	scaledLevFinder       *strdist.Finder
	hammingFinder         *strdist.Finder
}

// makeFinders returns a finders struct
func makeFinders(t *testing.T) finders {
	t.Helper()

	f := finders{}
	ff, err := strdist.NewCosineFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Cosine Finder: ", err)
	}
	f.cosineFinder = ff
	ff, err = strdist.NewJaccardFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Jaccard Finder: ", err)
	}
	f.jaccardFinder = ff
	ff, err = strdist.NewWeightedJaccardFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the WeightedJaccard Finder: ", err)
	}
	f.weightedJaccardFinder = ff
	ff, err = strdist.NewLevenshteinFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Levenshtein Finder: ", err)
	}
	f.levenshteinFinder = ff
	ff, err = strdist.NewScaledLevFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the ScaledLev Finder: ", err)
	}
	f.scaledLevFinder = ff
	ff, err = strdist.NewHammingFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the ScaledLev Finder: ", err)
	}
	f.hammingFinder = ff

	return f
}

func TestAllFinders(t *testing.T) {
	target := "test"
	pop := []string{
		"",
		"a",
		"b",
		"short",
		"a longer string.....................",
		"t",
		"te",
		"tes",
		"test",
		"not test",
		"test XXX",
	}

	f := makeFinders(t)

	testCases := []struct {
		testhelper.ID
		finder   *strdist.Finder
		distFunc func(string, string) float64
	}{
		{
			ID:     testhelper.MkID("cosine"),
			finder: f.cosineFinder,
			distFunc: func(s1, s2 string) float64 {
				d, err := strdist.CosineDistance(s1, s2, 2)
				if err != nil {
					panic("CosineDistance returned a non-nil error: " +
						err.Error())
				}
				return d
			},
		},
		{
			ID:     testhelper.MkID("jaccard"),
			finder: f.jaccardFinder,
			distFunc: func(s1, s2 string) float64 {
				d, err := strdist.JaccardDistance(s1, s2, 2)
				if err != nil {
					panic("JaccardDistance returned a non-nil error: " +
						err.Error())
				}
				return d
			},
		},
		{
			ID:     testhelper.MkID("weightedJaccard"),
			finder: f.weightedJaccardFinder,
			distFunc: func(s1, s2 string) float64 {
				d, err := strdist.WeightedJaccardDistance(s1, s2, 2)
				if err != nil {
					panic("WeightedJaccardDistance returned a non-nil error: " +
						err.Error())
				}
				return d
			},
		},
		{
			ID:     testhelper.MkID("levenshtein"),
			finder: f.levenshteinFinder,
			distFunc: func(s1, s2 string) float64 {
				return float64(strdist.LevenshteinDistance(s1, s2))
			},
		},
		{
			ID:       testhelper.MkID("scaledLev"),
			finder:   f.scaledLevFinder,
			distFunc: strdist.ScaledLevDistance,
		},
		{
			ID:       testhelper.MkID("hamming"),
			finder:   f.hammingFinder,
			distFunc: strdist.HammingDistance,
		},
	}

	for _, tc := range testCases {
		sdslice := tc.finder.FindLike(target, pop...)
		for _, sd := range sdslice {
			d := tc.distFunc(target, sd.Str)
			if d != sd.Dist {
				t.Log(tc.IDStr())
				t.Logf("\t:   distance from: %s\n", target)
				t.Logf("\t:              to: %s\n", sd.Str)
				t.Logf("\t: finder distance: %.5f\n", sd.Dist)
				t.Logf("\t: dist func gives: %.5f\n", d)
				t.Errorf("\t: the calculated distances differ\n")
			}
		}
	}
}
