package strdist_test

import (
	"testing"

	"github.com/nickwells/strdist.mod/strdist"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// finderChecker is a test Helper function that calls the finder funtion and
// checks the results
func finderChecker(t *testing.T, tID, subID, target string, pop []string, f *strdist.Finder, expect []string) {
	t.Helper()
	results := f.FindStrLike(target, pop...)
	if testhelper.StringSliceDiff(expect, results) {
		t.Log(tID)
		t.Logf("\t: comparing: %s", target)
		t.Logf("\t:   against: %v", pop)
		t.Logf("\t:  expected: %v", expect)
		t.Logf("\t:       got: %v", results)
		t.Errorf("\t: %s : results are unexpected\n", subID)
	}
}

// finderCheckerMaxN is a test Helper function that calls the limited results
// finder funtion and checks the results
func finderCheckerMaxN(t *testing.T, tID, subID, target string, pop []string, n int, f *strdist.Finder, expect []string) {
	t.Helper()
	results := f.FindNStrLike(n, target, pop...)
	if testhelper.StringSliceDiff(expect, results) {
		t.Log(tID)
		t.Logf("\t: comparing: %s", target)
		t.Logf("\t:   against: %v", pop)
		t.Logf("\t:  expected: %v", expect)
		t.Logf("\t:       got: %v", results)
		t.Errorf("\t: %s (max %d values): results are unexpected\n", subID, n)
	}
}
