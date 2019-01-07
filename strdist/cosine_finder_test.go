package strdist_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nickwells/mathutil.mod/mathutil"
	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestNewNGramsFinder(t *testing.T) {
	testCases := []struct {
		name        string
		ngLen       int
		minStrLen   int
		threshold   float64
		caseMod     strdist.CaseMod
		errExpected bool
		errContains []string
	}{
		{
			name:      "good",
			ngLen:     2,
			minStrLen: 4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
		},
		{
			name:        "bad ngLen (==0)",
			ngLen:       0,
			minStrLen:   -4,
			threshold:   1.2,
			caseMod:     strdist.NoCaseChange,
			errExpected: true,
			errContains: []string{
				"bad N-Gram length",
				"- it should be > 0",
			},
		},
		{
			name:        "bad ngLen (<0)",
			ngLen:       -1,
			minStrLen:   -4,
			threshold:   1.2,
			caseMod:     strdist.NoCaseChange,
			errExpected: true,
			errContains: []string{
				"bad N-Gram length",
				"- it should be > 0",
			},
		},
		{
			name:        "bad MinStrLen",
			ngLen:       2,
			minStrLen:   -4,
			threshold:   1.2,
			caseMod:     strdist.NoCaseChange,
			errExpected: true,
			errContains: []string{
				"bad minimum string length",
				"- it should be >= 0",
			},
		},
		{
			name:        "bad threshold",
			ngLen:       2,
			minStrLen:   4,
			threshold:   -1.0,
			caseMod:     strdist.NoCaseChange,
			errExpected: true,
			errContains: []string{
				"bad threshold",
				"- it should be >= 0.0",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		f, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, tc.caseMod)
		if err == nil {
			if tc.errExpected {
				t.Log(tcID)
				t.Errorf("\t: an error was expected but none was found\n")
			} else {
				if f == nil {
					t.Log(tcID)
					t.Errorf("\t: a nil pointer was returned but no error\n")
				} else {
					ca, ok := f.Algo.(*strdist.CosineAlgo)
					if !ok {
						t.Log(tcID)
						t.Errorf("\t: the Algo should be a *CosineAlgo\n")
					} else {
						if ca.N != tc.ngLen {
							t.Log(tcID)
							t.Errorf("\t: N-Gram Len should be: %d, was: %d\n",
								tc.ngLen, ca.N)
						}
					}
					if f.MinStrLen != tc.minStrLen {
						t.Log(tcID)
						t.Errorf("\t: minStrLen should be: %d, was: %d\n",
							tc.minStrLen, f.MinStrLen)
					}
					if f.T != tc.threshold {
						t.Log(tcID)
						t.Errorf("\t: threshold should be: %f, was: %f\n",
							tc.threshold, f.T)
					}
					if f.CM != tc.caseMod {
						t.Log(tcID)
						t.Errorf("\t: caseMod should be: %d, was: %d\n",
							tc.caseMod, f.CM)
					}
				}
			}
		} else {
			if !tc.errExpected {
				t.Log(tcID)
				t.Errorf("\t: an unexpected error was seen: %s\n", err)
			} else {
				var problemReported bool
				for _, errPart := range tc.errContains {
					if !strings.Contains(err.Error(), errPart) {
						if !problemReported {
							t.Log(tcID)
							t.Errorf(
								"\t: an unexpected error value was seen: %s\n",
								err)
						}
						t.Logf("\t: the error should contain : %s\n", errPart)
						problemReported = true
					}
				}
			}
		}
	}
}

func TestCosine(t *testing.T) {
	testCases := []struct {
		name        string
		s1, s2      string
		ngLen       int
		expDist     float64
		errExpected bool
		errContains []string
	}{
		{
			name:        "identical",
			s1:          "abab",
			s2:          "abab",
			ngLen:       2,
			expDist:     0.0,
			errExpected: false,
		},
		{
			name:        "both empty",
			s1:          "",
			s2:          "",
			ngLen:       2,
			expDist:     0.0,
			errExpected: false,
		},
		{
			name:        "first empty, second not",
			s1:          "",
			s2:          "abab",
			ngLen:       2,
			expDist:     1.0,
			errExpected: false,
		},
		{
			name:        "second empty, first not",
			s1:          "abab",
			s2:          "",
			ngLen:       2,
			expDist:     1.0,
			errExpected: false,
		},
		{
			name:        "no common n-grams",
			s1:          "abab",
			s2:          "cdcd",
			ngLen:       2,
			expDist:     1.0,
			errContains: []string{"invalid length of the n-gram:"},
		},
		{
			name:        "bad n-gram length (== 0)",
			s1:          "abab",
			s2:          "abab",
			ngLen:       0,
			expDist:     1.0,
			errExpected: true,
			errContains: []string{"invalid length of the n-gram:"},
		},
		{
			name:        "bad n-gram length (< 0)",
			s1:          "abab",
			s2:          "abab",
			ngLen:       -1,
			expDist:     1.0,
			errExpected: true,
			errContains: []string{"invalid length of the n-gram:"},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		dist, err := strdist.CosineDistance(tc.s1, tc.s2, tc.ngLen)

		if err == nil {
			if tc.errExpected {
				t.Log(tcID)
				t.Errorf("\t: an error was expected but not seen\n")
			}
		} else {
			if !tc.errExpected {
				t.Log(tcID)
				t.Errorf("\t: an unexpected error was seen: %s\n", err)
			} else {
				testhelper.ShouldContain(t, tcID, "error",
					err.Error(), tc.errContains)
			}
		}

		const epsilon = 0.00001
		if !mathutil.AlmostEqual(dist, tc.expDist, epsilon) {
			t.Log(tcID)
			t.Errorf("\t: the distance differs by more than %f"+
				" - expected: %.6f, got %.6f\n",
				epsilon, tc.expDist, dist)
		}
	}
}

func TestCosineFinder(t *testing.T) {
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
			expStringsNoChange:  []string{"hellos"},
			expStringsFlatCase:  []string{"HELLO", "hellos"},
			expNStringsFlatCase: []string{"HELLO"},
		},
		{
			name:                "empty target",
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

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		noChangeFinder, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.NoCaseChange)
		if err != nil {
			t.Log(tcID)
			t.Errorf("Couldn't create the NoCaseChange CosineFinder: %s",
				err)
			continue
		}
		flatCaseFinder, err := strdist.NewCosineFinder(
			tc.ngLen, tc.minStrLen, tc.threshold, strdist.ForceToLower)
		if err != nil {
			t.Log(tcID)
			t.Errorf("Couldn't create the ForceToLower CosineFinder: %s",
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
