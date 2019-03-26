package strdist_test

import (
	"testing"

	"github.com/nickwells/mathutil.mod/mathutil"
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
			t.Errorf("\t: Couldn't construct the ngrams for %s: %s",
				tc.s1, err)
		}
		ngs2, err := strdist.NGrams(tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't construct the ngrams for %s: %s", tc.s2, err)
		}
		ji := strdist.JaccardIndex(ngs1, ngs2)

		const epsilon = 0.00001
		if !mathutil.AlmostEqual(ji, tc.expVal, epsilon) {
			t.Log(tc.IDStr())
			t.Errorf("\t: the returned index should have been"+
				" within %f of %9.7f but was %9.7f",
				epsilon, tc.expVal, ji)
		}

		ji, err = strdist.JaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Unexpected error calculating the JaccardDistance: %s",
				err)
		}
		if !mathutil.AlmostEqual(ji, 1.0-tc.expVal, epsilon) {
			t.Log(tc.IDStr())
			t.Errorf("\t: the returned distance should have been"+
				" within %f of %9.7f but was %9.7f",
				epsilon, 1.0-tc.expVal, ji)
		}

		wji := strdist.WeightedJaccardIndex(ngs1, ngs2)
		if !mathutil.AlmostEqual(wji, tc.expWeightedVal, epsilon) {
			t.Log(tc.IDStr() + " (weighted)")
			t.Errorf("\t: the returned index should have been"+
				" within %f of %9.7f but was %9.7f",
				epsilon, tc.expWeightedVal, wji)
		}

		wji, err = strdist.WeightedJaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Log(tc.IDStr() + " (weighted)")
			t.Errorf("\t: Unexpected error calculating the"+
				" WeightedJaccardDistance: %s",
				err)
		}
		if !mathutil.AlmostEqual(wji, 1.0-tc.expWeightedVal, epsilon) {
			t.Log(tc.IDStr() + " (weighted)")
			t.Errorf("\t: the returned distance should have been"+
				" within %f of %9.7f but was %9.7f",
				epsilon, 1.0-tc.expWeightedVal, wji)
		}
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
			ID:                  testhelper.MkID("short population entry"),
			ngLen:               2,
			minStrLen:           4,
			threshold:           0.3,
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELLO"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:                  testhelper.MkID("empty target"),
			ngLen:               2,
			minStrLen:           0,
			threshold:           0.3,
			maxResults:          1,
			target:              "",
			pop:                 []string{"", "HELLO", "hellos", "hel", "world"},
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
			t.Errorf("Couldn't create the NoCaseChange JaccardFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewJaccardFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower JaccardFinder: %s",
				err)
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
