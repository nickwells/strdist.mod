package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// TestJaccard tests the Jaccard functions
func TestJaccard(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		s1, s2         string
		expVal         float64
		expWeightedVal float64
	}{
		{
			ID:             testhelper.MkID("abc and abcd"),
			s1:             "abc",
			s2:             "abcd",
			expVal:         0.666666667,
			expWeightedVal: 0.4,
		},
		{
			ID:             testhelper.MkID("both empty"),
			s1:             "",
			s2:             "",
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
	}

	for _, tc := range testCases {
		ngc := strdist.NGramConfig{Length: 2}
		ngs1 := ngc.NGrams(tc.s1)
		ngs2 := ngc.NGrams(tc.s2)

		const epsilon = 0.00001

		ji := strdist.JaccardIndex(ngs1, ngs2)
		testhelper.DiffFloat(t, tc.IDStr(), "Jaccard index",
			ji, tc.expVal, epsilon)

		wji := strdist.WeightedJaccardIndex(ngs1, ngs2)
		testhelper.DiffFloat(t, tc.IDStr(), "WeightedJaccard index",
			wji, tc.expWeightedVal, epsilon)

		j, err := strdist.NewJaccardAlgo(ngc, 0)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the JaccardAlgo: %s", err)
		}
		d := j.Dist(tc.s1, tc.s2)
		testhelper.DiffFloat(t, tc.IDStr(), "Jaccard distance",
			d, 1.0-tc.expVal, epsilon)

		wj, err := strdist.NewWeightedJaccardAlgo(ngc, 0)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the WeightedJaccardAlgo: %s", err)
		}
		d = wj.Dist(tc.s1, tc.s2)
		testhelper.DiffFloat(t, tc.IDStr(), "Weighted Jaccard distance",
			d, 1.0-tc.expWeightedVal, epsilon)
	}
}

func TestJaccardFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ngc                 strdist.NGramConfig
		fc                  strdist.FinderConfig
		maxResults          int
		target              string
		pop                 []string
		expStringsNoChange  []string
		expStringsFlatCase  []string
		expNStringsFlatCase []string
	}{
		{
			ID:  testhelper.MkID("std"),
			ngc: strdist.NGramConfig{Length: 2},
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 4,
			},
			maxResults:          0,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELL"},
			expNStringsFlatCase: []string{},
		},
		{
			ID:  testhelper.MkID("short target"),
			ngc: strdist.NGramConfig{Length: 2},
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 6,
			},
			maxResults:          99,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{},
			expNStringsFlatCase: []string{},
		},
		{
			ID:  testhelper.MkID("short population entry"),
			ngc: strdist.NGramConfig{Length: 2},
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 4,
			},
			maxResults: 1,
			target:     "hell",
			pop:        []string{"HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELLO"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:  testhelper.MkID("empty target"),
			ngc: strdist.NGramConfig{Length: 2},
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 0,
			},
			maxResults: 1,
			target:     "",
			pop:        []string{"", "HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{""},
			expStringsFlatCase:  []string{""},
			expNStringsFlatCase: []string{""},
		},
	}

	for _, tc := range testCases {
		ja, err := strdist.NewJaccardAlgo(tc.ngc, 0)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the Jaccard Algo: %s", err)
			continue
		}

		noChangeFinder, err := strdist.NewFinder(tc.fc, ja)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the standard JaccardFinder: %s", err)
			continue
		}

		fc := tc.fc
		fc.MapToLowerCase = true
		flatCaseFinder, err := strdist.NewFinder(fc, ja)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower JaccardFinder: %s", err)
			continue
		}

		finderChecker(t, tc.IDStr(), "no case change",
			tc.target, tc.pop, noChangeFinder, tc.expStringsNoChange)
		finderChecker(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, flatCaseFinder, tc.expStringsFlatCase)
		finderCheckerMaxN(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, tc.maxResults,
			flatCaseFinder, tc.expNStringsFlatCase)
	}
}
