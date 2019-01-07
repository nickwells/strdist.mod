package strdist

import (
	"fmt"
	"testing"
)

func TestConvertStrDist(t *testing.T) {
	testCases := []struct {
		name            string
		dists           []StrDist
		n               int
		expLen          int
		expShortLen     int
		expFirstVal     string
		expLastVal      string
		expLastShortVal string
	}{
		{
			name: "long dist slice",
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
			name: "short dist slice",
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

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		strsAll := convertStrDist(tc.dists)
		strsShort := convertStrDistN(tc.n, tc.dists)
		if len(strsAll) != tc.expLen {
			t.Log(tcID)
			t.Errorf("\t: bad conversion (all): len should be: %d, was: %d\n",
				tc.expLen, len(strsAll))
		}
		if len(strsShort) != tc.expShortLen {
			t.Log(tcID)
			t.Errorf("\t: bad conversion (short): len should be: %d, was: %d\n",
				tc.expShortLen, len(strsShort))
		}
		if len(strsAll) > 0 &&
			strsAll[0] != tc.expFirstVal {
			t.Log(tcID)
			t.Errorf("\t: bad first val (all): should be: '%s', was: '%s'\n",
				tc.expFirstVal, strsAll[0])
		}
		if len(strsShort) > 0 &&
			strsShort[0] != tc.expFirstVal {
			t.Log(tcID)
			t.Errorf("\t: bad first val (short): should be: '%s', was: '%s'\n",
				tc.expFirstVal, strsShort[0])
		}
		if len(strsAll) > 0 &&
			strsAll[len(strsAll)-1] != tc.expLastVal {
			t.Log(tcID)
			t.Errorf("\t: bad last val (all): should be: '%s', was: '%s'\n",
				tc.expLastVal, strsAll[len(strsAll)-1])
		}
		if len(strsShort) > 0 &&
			strsShort[len(strsShort)-1] != tc.expLastShortVal {
			t.Log(tcID)
			t.Errorf("\t: bad last val (short): should be: '%s', was: '%s'\n",
				tc.expLastShortVal, strsShort[len(strsShort)-1])
		}
	}

}
