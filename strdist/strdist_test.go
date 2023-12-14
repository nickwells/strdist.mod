package strdist_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestStrDistToString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		dist   strdist.StrDist
		expStr string
	}{
		{
			ID:     testhelper.MkID("dflt"),
			expStr: `Str: "", Dist: 0.00000`,
		},
		{
			ID:     testhelper.MkID("with vals"),
			dist:   strdist.StrDist{Str: "Hello", Dist: 1.23456},
			expStr: `Str: "Hello", Dist: 1.23456`,
		},
	}

	for _, tc := range testCases {
		s := tc.dist.String()
		testhelper.DiffString(t, tc.IDStr(), "StrDist.String()", s, tc.expStr)
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
		testhelper.ID
		i      int
		j      int
		expVal bool
	}{
		{ID: testhelper.MkID("same"), i: 0, j: 0, expVal: false},
		{ID: testhelper.MkID("same dist, str[i] < str[j]"), i: 0, j: 1, expVal: true},
		{ID: testhelper.MkID("same dist, str[i] > str[j]"), i: 1, j: 0, expVal: false},
		{ID: testhelper.MkID("dist[i] < dist[j]"), i: 2, j: 4, expVal: true},
		{ID: testhelper.MkID("dist[i] > dist[j]"), i: 4, j: 2, expVal: false},
	}

	for _, tc := range testCases {
		id := tc.IDStr() +
			fmt.Sprintf(" - Cmp(%d (%q), %d (%q))",
				tc.i, dists[tc.i], tc.j, dists[tc.j])

		val := strdist.SDSlice(dists).Cmp(tc.i, tc.j)

		testhelper.DiffBool(t, id, "comparison", val, tc.expVal)
	}
}
