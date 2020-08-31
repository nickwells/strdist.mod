package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// TestLevenshtein ...
func TestLevenshtein(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		a, b    string
		expDist int
	}{
		{
			ID:      testhelper.MkID("zero char same"),
			a:       "",
			b:       "",
			expDist: 0,
		},
		{
			ID:      testhelper.MkID("single char same"),
			a:       "a",
			b:       "a",
			expDist: 0,
		},
		{
			ID:      testhelper.MkID("single char differ"),
			a:       "a",
			b:       "b",
			expDist: 1,
		},
		{
			ID:      testhelper.MkID("differ 2"),
			a:       "aa",
			b:       "ab",
			expDist: 1,
		},
		{
			ID:      testhelper.MkID("Kitten/Sitting"),
			a:       "Kitten",
			b:       "Sitting",
			expDist: 3,
		},
		{
			ID:      testhelper.MkID("Saturday/Sunday"),
			a:       "Saturday",
			b:       "Sunday",
			expDist: 3,
		},
	}

	for _, tc := range testCases {
		dist := strdist.LevenshteinDistance(tc.a, tc.b)
		testhelper.CmpValInt(t,
			tc.IDStr(), fmt.Sprintf("LevenshteinDistance(%q, %q)", tc.a, tc.b),
			dist, tc.expDist)

		dist = strdist.LevenshteinDistance(tc.b, tc.a)
		testhelper.CmpValInt(t,
			tc.IDStr(), fmt.Sprintf("LevenshteinDistance(%q, %q)", tc.b, tc.a),
			dist, tc.expDist)
	}
}

func TestLevenshteinFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
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
			ID:         testhelper.MkID("std"),
			minStrLen:  4,
			threshold:  2,
			maxResults: 0,
			target:     "hello",
			pop:        []string{"HELL", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELL"},
			expNStringsFlatCase: []string{},
		},
		{
			ID:         testhelper.MkID("short target"),
			minStrLen:  6,
			threshold:  2,
			maxResults: 99,
			target:     "hello",
			pop:        []string{"HELL", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{},
			expNStringsFlatCase: []string{},
		},
		{
			ID:         testhelper.MkID("short population entry"),
			minStrLen:  4,
			threshold:  2,
			maxResults: 1,
			target:     "hell",
			pop:        []string{"HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
	}

	for _, tc := range testCases {
		noChangeFinder, err := strdist.NewLevenshteinFinder(
			tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the NoCaseChange LevenshteinFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewLevenshteinFinder(
			tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower LevenshteinFinder: %s",
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
