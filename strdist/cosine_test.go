package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestNewCosineAlgo(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		ngc       strdist.NGramConfig
		cacheSize int
	}{
		{
			ID:  testhelper.MkID("good"),
			ngc: strdist.NGramConfig{Length: 2},
		},
		{
			ID:     testhelper.MkID("bad ngLen (==0)"),
			ngc:    strdist.NGramConfig{Length: 0},
			ExpErr: testhelper.MkExpErr("the N-Gram Length (0) must be > 0"),
		},
		{
			ID:     testhelper.MkID("bad ngLen (<0)"),
			ngc:    strdist.NGramConfig{Length: -1},
			ExpErr: testhelper.MkExpErr("the N-Gram Length (-1) must be > 0"),
		},
		{
			ID:        testhelper.MkID("bad cacheSize (<0)"),
			ngc:       strdist.NGramConfig{Length: 2},
			cacheSize: -1,
			ExpErr:    testhelper.MkExpErr("the maxCacheSize (-1) must be >= 0"),
		},
	}

	for _, tc := range testCases {
		a, err := strdist.NewCosineAlgo(tc.ngc, tc.cacheSize)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			if a == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: a nil pointer was returned but no error\n")

				continue
			}

			if a.Name() != strdist.AlgoNameCosine {
				t.Log(tc.IDStr())
				t.Errorf("\t: the Algo should be a Cosine algo\n")
			}
		}
	}
}

func TestNewCosineFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		ngc strdist.NGramConfig
		fc  strdist.FinderConfig
	}{
		{
			ID:  testhelper.MkID("good"),
			ngc: strdist.NGramConfig{Length: 2},
			fc:  strdist.FinderConfig{Threshold: 1.2, MinStrLength: 4},
		},
		{
			ID:  testhelper.MkID("bad MinStrLen"),
			ngc: strdist.NGramConfig{Length: 2},
			fc:  strdist.FinderConfig{Threshold: 1.2, MinStrLength: -4},
			ExpErr: testhelper.MkExpErr(
				"FinderConfig: the minimum string length (-4) must be >= 0"),
		},
		{
			ID:  testhelper.MkID("bad threshold"),
			ngc: strdist.NGramConfig{Length: 2},
			fc:  strdist.FinderConfig{Threshold: -1.0, MinStrLength: 4},
			ExpErr: testhelper.MkExpErr(
				"FinderConfig: the Threshold (-1.000000) must be >= 0"),
		},
	}

	for _, tc := range testCases {
		a, err := strdist.NewCosineAlgo(tc.ngc, 0)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: an unexpected error making the cosine algo: %s\n",
				err)

			continue
		}

		f, err := strdist.NewFinder(tc.fc, a)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			if f == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: a nil pointer was returned but no error\n")

				continue
			}
		}
	}
}

func TestCosine(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s1, s2  string
		ngc     strdist.NGramConfig
		expDist float64
	}{
		{
			ID:      testhelper.MkID("identical"),
			s1:      "abab",
			s2:      "abab",
			ngc:     strdist.NGramConfig{Length: 2},
			expDist: 0.0,
		},
		{
			ID:      testhelper.MkID("both empty"),
			s1:      "",
			s2:      "",
			ngc:     strdist.NGramConfig{Length: 2},
			expDist: 0.0,
		},
		{
			ID:      testhelper.MkID("first empty, second not"),
			s1:      "",
			s2:      "abab",
			ngc:     strdist.NGramConfig{Length: 2},
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("second empty, first not"),
			s1:      "abab",
			s2:      "",
			ngc:     strdist.NGramConfig{Length: 2},
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("no common n-grams"),
			s1:      "abab",
			s2:      "cdcd",
			ngc:     strdist.NGramConfig{Length: 2},
			expDist: 1.0,
		},
		{
			ID:      testhelper.MkID("bad n-gram length (== 0)"),
			s1:      "abab",
			s2:      "abab",
			ngc:     strdist.NGramConfig{Length: 0},
			expDist: 1.0,
			ExpErr:  testhelper.MkExpErr("the N-Gram Length (0) must be > 0"),
		},
		{
			ID:      testhelper.MkID("bad n-gram length (< 0)"),
			s1:      "abab",
			s2:      "abab",
			ngc:     strdist.NGramConfig{Length: -1},
			expDist: 1.0,
			ExpErr:  testhelper.MkExpErr("the N-Gram Length (-1) must be > 0"),
		},
	}

	for _, tc := range testCases {
		const epsilon = 0.00001

		a, err := strdist.NewCosineAlgo(tc.ngc, 0)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			dist := a.Dist(tc.s1, tc.s2)
			testhelper.DiffFloat(t, tc.IDStr(), "distance",
				dist, tc.expDist, epsilon)
		}
	}
}

func TestCosineFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		ngc                 strdist.NGramConfig
		cacheSize           int
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
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:        testhelper.MkID("short population entry, cacheSize: 1"),
			ngc:       strdist.NGramConfig{Length: 2},
			cacheSize: 1,
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 4,
			},
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:        testhelper.MkID("short population entry, cacheSize: 2"),
			ngc:       strdist.NGramConfig{Length: 2},
			cacheSize: 2,
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 4,
			},
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:        testhelper.MkID("short population entry, cacheSize: 3"),
			ngc:       strdist.NGramConfig{Length: 2},
			cacheSize: 3,
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 4,
			},
			maxResults:          1,
			target:              "hell",
			pop:                 []string{"HELLO", "hellos", "hel", "world"},
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			ID:  testhelper.MkID("empty target"),
			ngc: strdist.NGramConfig{Length: 2},
			fc: strdist.FinderConfig{
				Threshold:    0.3,
				MinStrLength: 0,
			},
			maxResults:          1,
			target:              "",
			pop:                 []string{"", "HELLO", "hello", "hel", "world"},
			expStringsNoChange:  []string{""},
			expStringsFlatCase:  []string{""},
			expNStringsFlatCase: []string{""},
		},
	}

	for _, tc := range testCases {
		a, err := strdist.NewCosineAlgo(tc.ngc, tc.cacheSize)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the Cosine Algo: %s", err)

			continue
		}

		fc := tc.fc
		fc.MapToLowerCase = true

		f2lower, err := strdist.NewFinder(fc, a)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the map2Lower CosineFinder: %s", err)

			continue
		}

		fc = tc.fc

		f, err := strdist.NewFinder(fc, a)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: Couldn't create the CosineFinder: %s", err)

			continue
		}

		finderChecker(t, tc.IDStr(), "no case change",
			tc.target, tc.pop, f, tc.expStringsNoChange)
		finderChecker(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, f2lower, tc.expStringsFlatCase)
		finderCheckerMaxN(t, tc.IDStr(), "flattened case",
			tc.target, tc.pop, tc.maxResults,
			f2lower, tc.expNStringsFlatCase)
	}
}
