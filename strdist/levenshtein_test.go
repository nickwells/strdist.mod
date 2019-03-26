package strdist_test

import (
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
		for _, order := range []string{"a,b", "b,a"} {
			a, b := tc.a, tc.b
			if order == "b,a" {
				a, b = b, a
			}
			dist := strdist.LevenshteinDistance(a, b)
			if dist != tc.expDist {
				t.Logf(tc.IDStr())
				t.Errorf("\t: LevenshteinDistance('%s', '%s')"+
					" expected distance: %d got: %d",
					a, b, tc.expDist, dist)
			}
		}
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
			ID:                  testhelper.MkID("std"),
			minStrLen:           4,
			threshold:           2,
			maxResults:          0,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELL"},
			expNStringsFlatCase: []string{},
		},
		{
			ID:                  testhelper.MkID("short target"),
			minStrLen:           6,
			threshold:           2,
			maxResults:          99,
			target:              "hello",
			pop:                 []string{"HELL", "world"},
			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{},
			expNStringsFlatCase: []string{},
		},
		{
			ID:                  testhelper.MkID("short population entry"),
			minStrLen:           4,
			threshold:           2,
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
	}

	for _, tc := range testCases {
		noChangeLF, err := strdist.NewLevenshteinFinder(
			tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the NoCaseChange LevenshteinFinder: %s",
				err)
			continue
		}
		flatCaseLF, err := strdist.NewLevenshteinFinder(
			tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower LevenshteinFinder: %s",
				err)
			continue
		}

		noChangeSlice := noChangeLF.FindStrLike(tc.target, tc.pop...)
		flatCaseSlice := flatCaseLF.FindStrLike(tc.target, tc.pop...)
		flatCaseSliceShort := flatCaseLF.FindNStrLike(
			tc.maxResults, tc.target, tc.pop...)

		if testhelper.StringSliceDiff(noChangeSlice, tc.expStringsNoChange) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %v", tc.expStringsNoChange)
			t.Logf("\t:      got: %v", noChangeSlice)
			t.Errorf("\t: results are unexpected - no case change\n")
		}
		if testhelper.StringSliceDiff(flatCaseSlice, tc.expStringsFlatCase) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %v", tc.expStringsFlatCase)
			t.Logf("\t:      got: %v", flatCaseSlice)
			t.Errorf("\t: results are unexpected - flattened case\n")
		}
		if testhelper.StringSliceDiff(flatCaseSliceShort, tc.expNStringsFlatCase) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %v", tc.expNStringsFlatCase)
			t.Logf("\t:      got: %v", flatCaseSliceShort)
			t.Errorf(
				"\t: results are unexpected - flattened case, max %d values\n",
				tc.maxResults)
		}
	}
}
