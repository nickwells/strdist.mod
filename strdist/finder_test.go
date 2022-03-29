package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type TestAlgo struct{}

func (ta TestAlgo) Prep(_ string, _ strdist.CaseMod)            {}
func (ta TestAlgo) Dist(_, _ string, _ strdist.CaseMod) float64 { return 0.0 }

func TestCommonFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		minStrLen int
		threshold float64
		caseMod   strdist.CaseMod
	}{
		{
			ID:        testhelper.MkID("good"),
			minStrLen: 4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
		},
		{
			ID:        testhelper.MkID("bad MinStrLen"),
			minStrLen: -4,
			threshold: 1.2,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad minimum string length",
				"- it should be >= 0"),
		},
		{
			ID:        testhelper.MkID("bad threshold"),
			minStrLen: 4,
			threshold: -1.0,
			caseMod:   strdist.NoCaseChange,
			ExpErr: testhelper.MkExpErr(
				"bad threshold",
				"- it should be >= 0.0"),
		},
	}

	var a TestAlgo

	for _, tc := range testCases {
		cfi, err := strdist.NewFinder(
			tc.minStrLen, tc.threshold, tc.caseMod, a)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			if cfi == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: a nil pointer was returned but no error\n")
			} else {
				if cfi.MinStrLen != tc.minStrLen {
					t.Log(tc.IDStr())
					t.Errorf("\t: minStrLen should be: %d, was: %d\n",
						tc.minStrLen, cfi.MinStrLen)
				}
				if cfi.T != tc.threshold {
					t.Log(tc.IDStr())
					t.Errorf("\t: threshold should be: %f, was: %f\n",
						tc.threshold, cfi.T)
				}
				if cfi.CM != tc.caseMod {
					t.Log(tc.IDStr())
					t.Errorf("\t: caseMod should be: %d, was: %d\n",
						tc.caseMod, cfi.CM)
				}
			}
		}
	}
}
