package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
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

	for _, tc := range testCases {
		h := strdist.HammingAlgo{}
		dist := h.Dist(tc.a, tc.b)
		testhelper.DiffFloat(t,
			fmt.Sprintf("HammingDistance(%q, %q)", tc.a, tc.b), "distance",
			dist, tc.expDist, 0)
		dist = h.Dist(tc.b, tc.a)
		testhelper.DiffFloat(t,
			fmt.Sprintf("HammingDistance(%q, %q)", tc.b, tc.a), "distance",
			dist, tc.expDist, 0)
	}
}

func TestHammingFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		fc                  strdist.FinderConfig
		maxResults          int
		target              string
		pop                 []string
		expStringsNoChange  []string
		expStringsFlatCase  []string
		expNStringsFlatCase []string
	}{
		{
			ID: testhelper.MkID("std"),
			fc: strdist.FinderConfig{
				Threshold:    2.0,
				MinStrLength: 4,
			},
			maxResults: 0,
			target:     "hello",
			pop:        []string{"HELL", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{"HELL"},
			expNStringsFlatCase: []string{},
		},
		{
			ID: testhelper.MkID("short target"),
			fc: strdist.FinderConfig{
				Threshold:    2.0,
				MinStrLength: 6,
			},
			maxResults: 99,
			target:     "hello",
			pop:        []string{"HELL", "world"},

			expStringsNoChange:  []string{},
			expStringsFlatCase:  []string{},
			expNStringsFlatCase: []string{},
		},
		{
			ID: testhelper.MkID("short population entry"),
			fc: strdist.FinderConfig{
				Threshold:    2.0,
				MinStrLength: 4,
			},
			maxResults: 1,
			target:     "hell",
			pop:        []string{"HELLO", "hellos", "hel", "world"},

			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID: testhelper.MkID("empty target"),
			fc: strdist.FinderConfig{
				Threshold:    2.0,
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
		h := strdist.HammingAlgo{}
		f, err := strdist.NewFinder(tc.fc, h)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the NoCaseChange HammingFinder: %s", err)
			continue
		}
		fc := tc.fc
		fc.MapToLowerCase = true
		flatCaseFinder, err := strdist.NewFinder(fc, h)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("Couldn't create the ForceToLower HammingFinder: %s", err)
			continue
		}
		finderChecker(t, tc.IDStr(), "no case change",
			tc.target, tc.pop, f, tc.expStringsNoChange)
		finderChecker(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, flatCaseFinder, tc.expStringsFlatCase)
		finderCheckerMaxN(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, tc.maxResults,
			flatCaseFinder, tc.expNStringsFlatCase)
	}
}
