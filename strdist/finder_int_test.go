package strdist

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestConvertStrDist(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		dists           []StrDist
		n               int
		expLen          int
		expShortLen     int
		expFirstVal     string // first val is the same for long or short
		expLastVal      string
		expLastShortVal string
	}{
		{
			ID: testhelper.MkID("long dist slice"),
			dists: []StrDist{
				{Str: "str0", Dist: 1.2},
				{Str: "str1", Dist: 1.2},
				{Str: "str2", Dist: 1.2},
				{Str: "str3", Dist: 1.2},
				{Str: "str4", Dist: 1.2},
				{Str: "str5", Dist: 1.2},
				{Str: "str6", Dist: 1.2},
			},
			n:               3,
			expLen:          7,
			expShortLen:     3,
			expFirstVal:     "str0",
			expLastVal:      "str6",
			expLastShortVal: "str2",
		},
		{
			ID: testhelper.MkID("short dist slice"),
			dists: []StrDist{
				{Str: "str0", Dist: 1.2},
				{Str: "str1", Dist: 1.2},
			},
			n:               3,
			expLen:          2,
			expShortLen:     2,
			expFirstVal:     "str0",
			expLastVal:      "str1",
			expLastShortVal: "str1",
		},
	}

	for _, tc := range testCases {
		strsAll := convertStrDist(tc.dists)
		strsShort := convertStrDistN(tc.n, tc.dists)
		testhelper.DiffInt(t, tc.IDStr(), "all", len(strsAll), tc.expLen)
		testhelper.DiffInt(t, tc.IDStr(), "short", len(strsShort), tc.expShortLen)
		if len(strsAll) > 0 {
			checkVal(t, tc.IDStr(), "all", "first", strsAll[0], tc.expFirstVal)
			checkVal(t, tc.IDStr(), "all", "last",
				strsAll[len(strsAll)-1], tc.expLastVal)
		}
		if len(strsShort) > 0 {
			checkVal(t, tc.IDStr(), "short", "first", strsShort[0], tc.expFirstVal)
			checkVal(t, tc.IDStr(), "short", "last",
				strsShort[len(strsShort)-1], tc.expLastShortVal)
		}
		if len(strsShort) > 0 &&
			strsShort[0] != tc.expFirstVal {
			t.Log(tc.IDStr())
			t.Errorf("\t: bad first val (short): should be: %q, was: %q\n",
				tc.expFirstVal, strsShort[0])
		}
		if len(strsShort) > 0 &&
			strsShort[len(strsShort)-1] != tc.expLastShortVal {
			t.Log(tc.IDStr())
			t.Errorf("\t: bad last val (short): should be: %q, was: %q\n",
				tc.expLastShortVal, strsShort[len(strsShort)-1])
		}
	}
}

// checkVal reports values that differ from expectation
func checkVal(t *testing.T, tcID, name, vName string, val, expVal string) {
	t.Helper()

	if val != expVal {
		t.Log(tcID)
		t.Errorf("\t: bad %s val (%s): should be: %q, was: %q\n",
			vName, name, expVal, val)
	}
}
