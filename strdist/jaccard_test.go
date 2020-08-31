package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
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
		ngs1, err := strdist.NGrams(tc.s1, 2)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't construct the ngrams for %q: %s", tc.s1, err)
		}
		ngs2, err := strdist.NGrams(tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't construct the ngrams for %q: %s", tc.s2, err)
		}

		const epsilon = 0.00001
		ji := strdist.JaccardIndex(ngs1, ngs2)
		testhelper.CmpValFloat64(t, tc.IDStr(), "Jaccard index",
			ji, tc.expVal, epsilon)

		ji, err = strdist.JaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't calculate the JaccardDistance: %s", err)
		}
		testhelper.CmpValFloat64(t, tc.IDStr(), "Jaccard distance",
			ji, 1.0-tc.expVal, epsilon)

		wji := strdist.WeightedJaccardIndex(ngs1, ngs2)
		testhelper.CmpValFloat64(t, tc.IDStr(), "weighted Jaccard index",
			wji, tc.expWeightedVal, epsilon)

		wji, err = strdist.WeightedJaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr() + " (weighted)")
			t.Errorf("\t: Couldn't calculate the WeightedJaccardDistance: %s",
				err)
		}
		testhelper.CmpValFloat64(t, tc.IDStr(), "weighted Jaccard distance",
			wji, 1.0-tc.expWeightedVal, epsilon)
	}
}

func TestJaccardFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ngLen               int
		minStrLen           int
		threshold           float64
		maxResults          int
		target              string
		pop                 []string
		expStringsNoChange  []string
		expStringsFlatCase  []string
		expNStringsFlatCase []string
	}{
		{
			ID:                  testhelper.MkID("std"),
			ngLen:               2,
			minStrLen:           4,
			threshold:           0.3,
			maxResults:          0,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELL"},
			expNStringsFlatCase: []string{},
		},
		{
			ID:                  testhelper.MkID("short target"),
			ngLen:               2,
			minStrLen:           6,
			threshold:           0.3,
			maxResults:          99,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{},
			expNStringsFlatCase: []string{},
		},
		{
			ID:         testhelper.MkID("short population entry"),
			ngLen:      2,
			minStrLen:  4,
			threshold:  0.3,
			maxResults: 1,
			target:     "hell",
			pop:        []string{"HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELLO"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:         testhelper.MkID("empty target"),
			ngLen:      2,
			minStrLen:  0,
			threshold:  0.3,
			maxResults: 1,
			target:     "",
			pop:        []string{"", "HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{""},
			expStringsFlatCase:  []string{""},
			expNStringsFlatCase: []string{""},
		},
	}

	for _, tc := range testCases {
		noChangeFinder, err := strdist.NewJaccardFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the NoCaseChange JaccardFinder: %s", err)
			continue
		}
		flatCaseFinder, err := strdist.NewJaccardFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.ForceToLower)
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
