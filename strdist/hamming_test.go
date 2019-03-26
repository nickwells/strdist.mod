package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// TestHamming ...
func TestHamming(t *testing.T) {
	testCases := []struct {
		a, b    string
		expDist float64
	}{
		{
			a:       "a",
			b:       "b",
			expDist: 1,
		},
		{
			a:       "ab",
			b:       "b",
			expDist: 2,
		},
		{
			a:       "aaa",
			b:       "aba",
			expDist: 1,
		},
		{
			a:       "aaa",
			b:       "a¶a",
			expDist: 1,
		},
		{
			a:       "a§a",
			b:       "a¶a",
			expDist: 1,
		},
		{
			a:       "a§abc",
			b:       "a¶a",
			expDist: 3,
		},
		{
			a:       "a",
			b:       "a",
			expDist: 0,
		},
		{
			a:       "",
			b:       "",
			expDist: 0,
		},
		{
			a:       "abc",
			b:       "abc",
			expDist: 0,
		},
	}

	for i, tc := range testCases {
		for _, order := range []string{"a,b", "b,a"} {
			a, b := tc.a, tc.b
			if order == "b,a" {
				a, b = b, a
			}

			if dist := strdist.HammingDistance(a, b); dist != tc.expDist {
				t.Errorf("test %d (%s): HammingDistance('%s', '%s')"+
					" should have been %f but was %f",
					i, order, a, b, tc.expDist, dist)
			}
		}
	}
}

func TestHammingFinder(t *testing.T) {
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
			threshold:           2.0,
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
			threshold:           2.0,
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
			threshold:           2.0,
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:                  testhelper.MkID("empty target"),
			minStrLen:           0,
			threshold:           2.0,
			maxResults:          1,
			target:              "",
			pop:                 []string{"", "HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{""},
			expStringsFlatCase:  []string{""},
			expNStringsFlatCase: []string{""},
		},
	}

	for _, tc := range testCases {
		noChangeFinder, err := strdist.NewHammingFinder(
			tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the NoCaseChange HammingFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewHammingFinder(
			tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower HammingFinder: %s",
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
