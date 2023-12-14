package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type TestAlgo struct{}

func (TestAlgo) Dist(_, _ string) float64 { return 0.0 }
func (TestAlgo) Name() string             { return "TestAlgo" }
func (TestAlgo) Desc() string             { return "" }

func TestCommonFinder(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fc strdist.FinderConfig
	}{
		{
			ID: testhelper.MkID("good"),
			fc: strdist.FinderConfig{Threshold: 1.2, MinStrLength: 4},
		},
		{
			ID: testhelper.MkID("bad MinStrLen"),
			fc: strdist.FinderConfig{Threshold: 1.2, MinStrLength: -4},
			ExpErr: testhelper.MkExpErr(
				"FinderConfig: the minimum string length (-4) must be >= 0"),
		},
		{
			ID: testhelper.MkID("bad threshold"),
			fc: strdist.FinderConfig{Threshold: -1.0, MinStrLength: 4},
			ExpErr: testhelper.MkExpErr(
				"FinderConfig: the Threshold (-1.000000) must be >= 0"),
		},
	}

	var a TestAlgo

	for _, tc := range testCases {
		f, err := strdist.NewFinder(tc.fc, a)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			if f == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: a nil pointer was returned but no error\n")
			}
		}
	}
}
