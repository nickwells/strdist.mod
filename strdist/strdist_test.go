package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
)

func TestStrDistToString(t *testing.T) {
	testCases := []struct {
		name   string
		dist   strdist.StrDist
		expStr string
	}{
		{
			name:   "dflt",
			expStr: "Str: '', Dist: 0.00000",
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		s := tc.dist.String()
		if s != tc.expStr {
			t.Log(tcID)
			t.Logf("\t:      got: %s\n", s)
			t.Logf("\t: expected: %s\n", tc.expStr)
			t.Errorf("\t: bad string conversion\n")
		}
	}

}

func TestStrDistCmp(t *testing.T) {
	dists := []strdist.StrDist{
		{Str: "a", Dist: 1.1},
		{Str: "b", Dist: 1.1},
		{Str: "c", Dist: 1.2},
		{Str: "d", Dist: 1.2},
		{Str: "e", Dist: 1.3},
	}
	testCases := []struct {
		name   string
		i      int
		j      int
		expVal bool
	}{
		{name: "same", i: 0, j: 0, expVal: false},
		{name: "same dist, str[i] < str[j]", i: 0, j: 1, expVal: true},
		{name: "same dist, str[i] > str[j]", i: 1, j: 0, expVal: false},
		{name: "dist[i] < dist[j]", i: 2, j: 4, expVal: true},
		{name: "dist[i] > dist[j]", i: 4, j: 2, expVal: false},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		val := strdist.SDSlice(dists).Cmp(tc.i, tc.j)
		if val != tc.expVal {
			t.Log(tcID)
			t.Logf("\t: Comparing (%2d) %s\n", tc.i, dists[tc.i])
			t.Logf("\t:      With (%2d) %s\n", tc.j, dists[tc.j])
			t.Errorf("\t: bad comparison, expected %t, got %t\n",
				tc.expVal, val)
		}
	}

}
