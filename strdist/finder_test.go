package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

type TestAlgo struct{}

func (ta TestAlgo) Prep(_ string, _ strdist.CaseMod)            {}
func (ta TestAlgo) Dist(_, _ string, _ strdist.CaseMod) float64 { return 0.0 }

func TestCommonFinder(t *testing.T) {
	testCases := []struct {
		name        string
		minStrLen   int
		threshold   float64
		caseMod     strdist.CaseMod
		errExpected bool
		errContains []string
	}{
		{
			name:      "good",
			minStrLen: 4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
		},
		{
			name:        "bad MinStrLen",
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

	var a TestAlgo

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		cfi, err := strdist.NewFinder(
			tc.minStrLen, tc.threshold, tc.caseMod, a)
		if err == nil {
			if tc.errExpected {
				t.Log(tcID)
				t.Errorf("\t: an error was expected but none was found\n")
			} else {
				if cfi == nil {
					t.Log(tcID)
					t.Errorf("\t: a nil pointer was returned but no error\n")
				} else {
					if cfi.MinStrLen != tc.minStrLen {
						t.Log(tcID)
						t.Errorf("\t: minStrLen should be: %d, was: %d\n",
							tc.minStrLen, cfi.MinStrLen)
					}
					if cfi.T != tc.threshold {
						t.Log(tcID)
						t.Errorf("\t: threshold should be: %f, was: %f\n",
							tc.threshold, cfi.T)
					}
					if cfi.CM != tc.caseMod {
						t.Log(tcID)
						t.Errorf("\t: caseMod should be: %d, was: %d\n",
							tc.caseMod, cfi.CM)
					}
				}
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
	}
}
