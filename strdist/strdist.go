package strdist

import "fmt"

// StrDist records a string and its associated distance
type StrDist struct {
	Str  string
	Dist float64
}

// String returns a string form of the StrDist
func (sd StrDist) String() string {
	return fmt.Sprintf("Str: %q, Dist: %.5f", sd.Str, sd.Dist)
}

// SDSlice holds a list of strings and associated distances
type SDSlice []StrDist

// Cmp tests whether the i'th element in the SDSlice is less than the j'th
// element. If the elements each have the same distance then the comparison
// is on the basis of the lexicographical order of the strings, otherwise it
// is based on the ordering of the distances
func (sd SDSlice) Cmp(i, j int) bool {
	if sd[i].Dist == sd[j].Dist {
		return sd[i].Str < sd[j].Str
	}
	return sd[i].Dist < sd[j].Dist
}

// lessThanFunc returns a function that will compare the two StrDist values
func lessThanFunc(strLen int) func(sd1, sd2 StrDist) bool {
	return func(sd1, sd2 StrDist) bool {
		// firstly compare by the distance metric
		if sd1.Dist != sd2.Dist {
			return sd1.Dist < sd2.Dist
		}

		// then compare by closeness in length to the target string
		lenDiff1, lenDiff2 := len(sd1.Str)-strLen, len(sd2.Str)-strLen
		sqLenDiff1, sqLenDiff2 := lenDiff1*lenDiff1, lenDiff2*lenDiff2
		if sqLenDiff1 != sqLenDiff2 {
			return sqLenDiff1 < sqLenDiff2
		}

		// finally compare by lexical ordering
		return sd1.Str < sd2.Str
	}
}
