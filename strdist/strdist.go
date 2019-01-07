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

type SDSlice []StrDist

func (sd SDSlice) Cmp(i, j int) bool {
	if sd[i].Dist == sd[j].Dist {
		return sd[i].Str < sd[j].Str
	}
	return sd[i].Dist < sd[j].Dist
}
