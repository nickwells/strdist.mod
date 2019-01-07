package strdist

import (
	"fmt"
	"math"
)

// NGramSet represents a set of n-grams. Each n-gram has an associated weight
// which is the number of occurences in the original string from which it is
// derived
type NGramSet map[string]int

// Dot returns the dot product of the two n-gram sets
func Dot(n1, n2 NGramSet) int64 {
	var d int64
	for k, v := range n1 {
		d += int64(v * n2[k])
	}
	return d
}

// lengthSquared returns the square of the Length of the n-gram set
func (ngs NGramSet) lengthSquared() float64 {
	var l float64
	for _, v := range ngs {
		l += float64(v * v)
	}
	return l
}

// Length returns the Length (sometimes called the Magnitude) of the n-gram set
func (ngs NGramSet) Length() float64 {
	return math.Sqrt(ngs.lengthSquared())
}

// NGrams transforms the string s into a map of n-grams (substrings of s each
// of length n). The key is the n-gram and the value is the number of
// occurrences. The length of the n-grams (n) must be greater than zero or an
// error will be returned
func NGrams(s string, n int) (NGramSet, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid length of the n-gram: %d", n)
	}

	ngrams := make(NGramSet)
	if len(s) < n {
		return ngrams, nil
	}

	chars := make([]rune, n)
	for i, r := range s {
		chars[i%n] = r
		if i >= n-1 {
			offset := (i + 1) % n
			str := string(chars[offset:])
			if offset != 0 {
				str += string(chars[0:offset])
			}
			ngrams[str]++
		}
	}

	return ngrams, nil
}

// NGramUnion returns a map which contains the union of the two sets of strings
// the associated counts are added together
func NGramUnion(ngs1, ngs2 NGramSet) NGramSet {
	union := make(NGramSet)

	for k, v := range ngs1 {
		union[k] = v
	}

	for k, v := range ngs2 {
		union[k] += v
	}

	return union
}

// NGramLenUnion returns the length of the union of the two sets of strings
// without having to construct the union
func NGramLenUnion(ngs1, ngs2 NGramSet) int {
	lenNGS1 := len(ngs1)
	lenNGS2 := len(ngs2)

	unionLen := lenNGS1 + lenNGS2

	rangeMap, otherMap := ngs1, ngs2
	if lenNGS1 > lenNGS2 {
		rangeMap, otherMap = otherMap, rangeMap
	}

	for k := range rangeMap {
		_, ok := otherMap[k]
		if ok {
			unionLen--
		}
	}

	return unionLen
}

// NGramWeightedLenUnion returns the weighted length of the union of the two
// sets of strings.  The weights are the map values (the number of instances
// of the key) and for the union we take the sum of the two values.
func NGramWeightedLenUnion(ngs1, ngs2 NGramSet) int {
	unionLen := 0

	for _, v := range ngs1 {
		unionLen += v
	}

	for _, v := range ngs2 {
		unionLen += v
	}

	return unionLen
}

// NGramIntersection returns a map which contains the intersection of the two
// sets of strings.  We take the minimum of the two associated counts
func NGramIntersection(ngs1, ngs2 NGramSet) NGramSet {
	intersection := make(NGramSet)

	for k, ngs1v := range ngs1 {
		ngs2v, ok := ngs2[k]
		if ok {
			v := ngs1v
			if ngs2v < v {
				v = ngs2v
			}
			intersection[k] = v
		}
	}

	return intersection
}

// NGramLenIntersection returns the length of the intersection of the two
// sets of strings without having to construct the intersection
func NGramLenIntersection(ngs1, ngs2 NGramSet) int {
	intersectionLen := 0

	rangeMap, otherMap := ngs1, ngs2
	if len(ngs1) > len(ngs2) {
		rangeMap, otherMap = otherMap, rangeMap
	}

	for k := range rangeMap {
		_, ok := otherMap[k]
		if ok {
			intersectionLen++
		}
	}

	return intersectionLen
}

// NGramWeightedLenIntersection returns the weighted length of the
// intersection of the two sets of strings. The weights are the map values
// (the number of instances of the key) and for the intersection we take the
// minimum of the two values.
func NGramWeightedLenIntersection(ngs1, ngs2 NGramSet) int {
	intersectionLen := 0

	rangeMap, otherMap := ngs1, ngs2
	if len(ngs1) > len(ngs2) {
		rangeMap, otherMap = otherMap, rangeMap
	}

	for k, rmV := range rangeMap {
		omV, ok := otherMap[k]
		if ok {
			v := rmV
			if omV < v {
				v = omV
			}
			intersectionLen += v
		}
	}

	return intersectionLen
}

// WeightedLen returns the weighted length of the set of strings. The weights
// are the map values (the number of instances of the key)
func (ngs NGramSet) WeightedLen() int {
	wLen := 0

	for _, v := range ngs {
		wLen += v
	}

	return wLen
}

// NGramsEqual compares the two sets and returns true if they are equal,
// false otherwise
func NGramsEqual(ngs1, ngs2 NGramSet) bool {
	if len(ngs1) != len(ngs2) {
		return false
	}

	for k, ngs1v := range ngs1 {
		ngs2v, ok := ngs2[k]
		if !ok {
			return false
		}
		if ngs1v != ngs2v {
			return false
		}
	}

	return true
}

// OverlapCoefficient constructs the overlap coefficient (sometimes known as the
// Szymkiewicz-Simpson coefficient) of the two n-gram sets
func OverlapCoefficient(ngs1, ngs2 NGramSet) float64 {
	minLen := len(ngs1)
	if len(ngs2) < minLen {
		minLen = len(ngs2)
	}
	if minLen == 0 {
		return 1.0
	}

	iLen := NGramLenIntersection(ngs1, ngs2)
	return float64(iLen) / float64(minLen)
}

// WeightedOverlapCoefficient constructs the weighted overlap coefficient
// (sometimes known as the Szymkiewicz-Simpson coefficient) of the two n-gram
// sets
func WeightedOverlapCoefficient(ngs1, ngs2 NGramSet) float64 {
	lenNGS1 := ngs1.WeightedLen()
	lenNGS2 := ngs2.WeightedLen()
	minLen := lenNGS1

	if lenNGS2 < minLen {
		minLen = lenNGS2
	}

	if minLen == 0 {
		return 1.0
	}

	iLen := NGramWeightedLenIntersection(ngs1, ngs2)
	return float64(iLen) / float64(minLen)
}
