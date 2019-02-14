package strdist

import "fmt"

// StrDist records a string and its associated distance
type StrDist struct {
	Str  string
	Dist float64
}

// String returns a string form of the StrDist
func (sd StrDist) String() string {
	return fmt.Sprintf("Str: '%s', Dist: %.5f", sd.Str, sd.Dist)
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
