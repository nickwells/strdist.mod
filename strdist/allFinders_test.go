package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
)

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

	cosineFinder, err :=
		strdist.NewCosineFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Cosine Finder: ", err)
	}
	jaccardFinder, err :=
		strdist.NewJaccardFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Jaccard Finder: ", err)
	}
	weightedJaccardFinder, err :=
		strdist.NewWeightedJaccardFinder(2, 0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the WeightedJaccard Finder: ", err)
	}
	levenshteinFinder, err :=
		strdist.NewLevenshteinFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the Levenshtein Finder: ", err)
	}
	scaledLevFinder, err :=
		strdist.NewScaledLevFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the ScaledLev Finder: ", err)
	}
	hammingFinder, err :=
		strdist.NewHammingFinder(0, 1.0, strdist.NoCaseChange)
	if err != nil {
		t.Fatal("couldn't create the ScaledLev Finder: ", err)
	}

	testCases := []struct {
		name     string
		finder   *strdist.Finder
		distFunc func(string, string) float64
	}{
		{
			name:   "cosine",
			finder: cosineFinder,
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
			name:   "jaccard",
			finder: jaccardFinder,
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
			name:   "weightedJaccard",
			finder: weightedJaccardFinder,
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
			name:   "levenshtein",
			finder: levenshteinFinder,
			distFunc: func(s1, s2 string) float64 {
				return float64(strdist.LevenshteinDistance(s1, s2))
			},
		},
		{
			name:   "scaledLev",
			finder: scaledLevFinder,
			distFunc: func(s1, s2 string) float64 {
				return float64(strdist.ScaledLevDistance(s1, s2))
			},
		},
		{
			name:   "hamming",
			finder: hammingFinder,
			distFunc: func(s1, s2 string) float64 {
				return float64(strdist.HammingDistance(s1, s2))
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)
		sdslice := tc.finder.FindLike(target, pop...)
		for _, sd := range sdslice {
			d := tc.distFunc(target, sd.Str)
			if d != sd.Dist {
				t.Log(tcID)
				t.Errorf("\t:   distance from: %s\n", target)
				t.Errorf("\t:              to: %s\n", sd.Str)
				t.Errorf("\t: finder distance: %.5f\n", sd.Dist)
				t.Errorf("\t: dist func gives: %.5f\n", d)
				t.Errorf("\t: the calculated distances differ\n")
			}
		}
	}
}
