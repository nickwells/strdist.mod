package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestNewNGramsFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		ngLen     int
		minStrLen int
		threshold float64
		caseMod   strdist.CaseMod
	}{
		{
			ID:        testhelper.MkID("good"),
			ngLen:     2,
			minStrLen: 4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
		},
		{
			ID:        testhelper.MkID("bad ngLen (==0)"),
			ngLen:     0,
			minStrLen: -4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad N-Gram length",
				"- it should be > 0"),
		},
		{
			ID:        testhelper.MkID("bad ngLen (<0)"),
			ngLen:     -1,
			minStrLen: -4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad N-Gram length",
				"- it should be > 0"),
		},
		{
			ID:        testhelper.MkID("bad MinStrLen"),
			ngLen:     2,
			minStrLen: -4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad minimum string length",
				"- it should be >= 0"),
		},
		{
			ID:        testhelper.MkID("bad threshold"),
			ngLen:     2,
			minStrLen: 4,
			threshold: -1.0,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad threshold",
				"- it should be >= 0.0"),
		},
	}

	for _, tc := range testCases {
		f, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, tc.caseMod)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			if f == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: a nil pointer was returned but no error\n")
				continue
			}

			id := tc.IDStr()
			ca, ok := f.Algo.(*strdist.CosineAlgo)
			if !ok {
				t.Log(tc.IDStr())
				t.Errorf("\t: the Algo should be a *CosineAlgo\n")
			} else {
				testhelper.DiffInt(t, id, "N-Gram length", ca.N, tc.ngLen)
			}
			testhelper.DiffInt(t, id, "MinStrLen", f.MinStrLen, tc.minStrLen)
			testhelper.DiffInt(t, id, "caseMod", int(f.CM), int(tc.caseMod))
			testhelper.DiffFloat64(t, id, "threshold", f.T, tc.threshold, 0.0)
		}
	}
}

func TestCosine(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s1, s2  string
		ngLen   int
		expDist float64
	}{
		{
			ID:      testhelper.MkID("identical"),
			s1:      "abab",
			s2:      "abab",
			ngLen:   2,
			expDist: 0.0,
		},
		{
			ID:      testhelper.MkID("both empty"),
			s1:      "",
			s2:      "",
			ngLen:   2,
			expDist: 0.0,
		},
		{
			ID:      testhelper.MkID("first empty, second not"),
			s1:      "",
			s2:      "abab",
			ngLen:   2,
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("second empty, first not"),
			s1:      "abab",
			s2:      "",
			ngLen:   2,
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("no common n-grams"),
			s1:      "abab",
			s2:      "cdcd",
			ngLen:   2,
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("bad n-gram length (== 0)"),
			s1:      "abab",
			s2:      "abab",
			ngLen:   0,
			expDist: 1.0,
			ExpErr:  testhelper.MkExpErr("invalid length of the n-gram:"),
		},
		{
			ID:      testhelper.MkID("bad n-gram length (< 0)"),
			s1:      "abab",
			s2:      "abab",
			ngLen:   -1,
			expDist: 1.0,
			ExpErr:  testhelper.MkExpErr("invalid length of the n-gram:"),
		},
	}

	for _, tc := range testCases {
		dist, err := strdist.CosineDistance(tc.s1, tc.s2, tc.ngLen)

		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffFloat64(t, tc.IDStr(), "distance",
				dist, tc.expDist, 0.00001)
		}
	}
}

func TestCosineFinder(t *testing.T) {
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
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:                  testhelper.MkID("empty target"),
			ngLen:               2,
			minStrLen:           0,
			threshold:           0.3,
			maxResults:          1,
			target:              "",
			pop:                 []string{"", "HELLO", "hello", "hel", "world"},
			expStringsNoChange:  []string{""},
			expStringsFlatCase:  []string{""},
			expNStringsFlatCase: []string{""},
		},
	}

	for _, tc := range testCases {
		noChangeFinder, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the NoCaseChange CosineFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the ForceToLower CosineFinder: %s",
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
