package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/mathutil.mod/mathutil"
	"github.com/nickwells/strdist.mod/strdist"
)

// TestJaccard tests the Jaccard functions
func TestJaccard(t *testing.T) {
	testCases := []struct {
		name           string
		s1, s2         string
		expVal         float64
		expWeightedVal float64
	}{
		{
			name:           "abc and abcd",
			s1:             "abc",
			s2:             "abcd",
			expVal:         0.666666667,
			expWeightedVal: 0.4,
		},
		{
			name:           "both empty",
			s1:             "",
			s2:             "",
			expVal:         1.0,
			expWeightedVal: 1.0,
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.name)
		ngs1, err := strdist.NGrams(tc.s1, 2)
		if err != nil {
			t.Errorf("%s : Couldn't construct the ngrams for %s: %s",
				testID, tc.s1, err)
		}
		ngs2, err := strdist.NGrams(tc.s2, 2)
		if err != nil {
			t.Errorf("%s : Couldn't construct the ngrams for %s: %s",
				testID, tc.s2, err)
		}
		ji := strdist.JaccardIndex(ngs1, ngs2)

		const epsilon = 0.00001
		if !mathutil.AlmostEqual(ji, tc.expVal, epsilon) {
			t.Errorf("%s : the returned index should have been"+
				" within %f of %9.7f but was %9.7f",
				testID, epsilon, tc.expVal, ji)
		}

		ji, err = strdist.JaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Errorf(
				"%s : Unexpected error constructing the JaccardDistance: %s",
				testID, err)
		}
		if !mathutil.AlmostEqual(ji, 1.0-tc.expVal, epsilon) {
			t.Errorf("%s : the returned distance should have been"+
				" within %f of %9.7f but was %9.7f",
				testID, epsilon, 1.0-tc.expVal, ji)
		}

		wji := strdist.WeightedJaccardIndex(ngs1, ngs2)
		if !mathutil.AlmostEqual(wji, tc.expWeightedVal, epsilon) {
			t.Errorf("%s (weighted) : the returned index should have been"+
				" within %f of %9.7f but was %9.7f",
				testID, epsilon, tc.expWeightedVal, wji)
		}

		wji, err = strdist.WeightedJaccardDistance(tc.s1, tc.s2, 2)
		if err != nil {
			t.Errorf(
				"%s : Unexpected error constructing the WeightedJaccardDistance: %s",
				testID, err)
		}
		if !mathutil.AlmostEqual(wji, 1.0-tc.expWeightedVal, epsilon) {
			t.Errorf("%s (weighted) : the returned distance should have been"+
				" within %f of %9.7f but was %9.7f",
				testID, epsilon, 1.0-tc.expWeightedVal, wji)
		}
	}
}

func TestJaccardFinder(t *testing.T) {
	testCases := []struct {
		name                string
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
			name:                "std",
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
			name:                "short target",
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
			name:                "short population entry",
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
			name:                "empty target",
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

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		noChangeFinder, err := strdist.NewJaccardFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tcID)
			t.Errorf("Couldn't create the NoCaseChange JaccardFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewJaccardFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tcID)
			t.Errorf("Couldn't create the ForceToLower JaccardFinder: %s",
				err)
			continue
		}

		finderChecker(t, tcID, "no case change",
			tc.target, tc.pop, noChangeFinder, tc.expStringsNoChange)
		finderChecker(t, tcID, "flattened case",
			tc.target, tc.pop, flatCaseFinder, tc.expStringsFlatCase)
		finderCheckerMaxN(t, tcID, "flattened case",
			tc.target, tc.pop, tc.maxResults,
			flatCaseFinder, tc.expNStringsFlatCase)
	}
}
