package strdist

import (
	"fmt"
)

// NGramConfig holds information about how to construct the NGramSet from a
// given string
type NGramConfig struct {
	// Length is the target length of an NGram. If MinLength is set to a
	// value greater than zero and lower than this then NGrams of length
	// MinLength..Length will be generated.
	Length int
	// MinLength gives the minimum length of an NGram. If it is greater than
	// zero the generated NGrams will have a range of lengths from MinLength
	// to Length. If it is equal to 0 then all the NGrams will be of length
	// Length.
	MinLength int
	// OverflowTheSource is set to indicate that the collection of NGrams
	// generated will include ones starting before the source string and
	// finishing after the source string. Runes outside the source string
	// will be set to nil. This allows substrings at the start and end of the
	// string to be captured explicitly.
	OverFlowTheSource bool
}

// DfltNGramConfig is a suggested value for the NGram config
var DfltNGramConfig = NGramConfig{
	Length:            3,
	MinLength:         2,
	OverFlowTheSource: true,
}

// Check returns a non-nil error if the NGramConfig has invalid entries
func (ngc NGramConfig) Check() error {
	if ngc.Length <= 0 {
		return fmt.Errorf("the N-Gram Length (%d) must be > 0", ngc.Length)
	}
	if ngc.MinLength < 0 {
		return fmt.Errorf("the N-Gram MinLength (%d) must be >= 0",
			ngc.MinLength)
	}
	if ngc.MinLength > ngc.Length {
		return fmt.Errorf(
			"the N-Gram MinLength (%d) must be <= the Length (%d)",
			ngc.MinLength, ngc.Length)
	}
	return nil
}

// String returns a string describing the NGramConfig
func (ngc NGramConfig) String() string {
	s := fmt.Sprintf("MinLength: %2d", ngc.MinLength)
	s += fmt.Sprintf(", Length: %2d", ngc.Length)
	s += fmt.Sprintf(", Overflow: %-5.5v", ngc.OverFlowTheSource)
	return s
}

// String returns a string describing the NGramConfig
func (ngc NGramConfig) Desc(prefix string) string {
	s := fmt.Sprintf("%s Min: %2d", prefix, ngc.MinLength)
	s += fmt.Sprintf(" Len: %2d", ngc.Length)
	s += fmt.Sprintf(" O'flow: %5t", ngc.OverFlowTheSource)
	return s
}

// calcStartIdx calculates the start index and returns true if it is less
// than or equal to the maxIdx, false otherwise
func (ngc NGramConfig) calcStartIdx(endIdx, maxIdx, length int) (int, bool) {
	startIdx := endIdx - (length - 1)
	if startIdx < 0 {
		startIdx += ngc.Length
		if startIdx > maxIdx {
			if !ngc.OverFlowTheSource {
				return startIdx, false
			}
			startIdx = maxIdx
		}
	}
	return startIdx, true
}

// overflowStrings returns the strings produced by overflowing the source.
func (ngc NGramConfig) overflowStrings(idx, maxIdx int, ngStrings *[]string) {
	if idx != maxIdx {
		return
	}
	if !ngc.OverFlowTheSource {
		return
	}

	overflow := make([]rune, ngc.Length-1)

	for _, str := range *ngStrings {
		for l := len(str) + 1; l <= ngc.Length; l++ {
			diff := l - len(str)
			str += string(overflow[:diff])
			*ngStrings = append(*ngStrings, str)
		}
	}
}

// subGrams returns all the NGrams from the passed sub-string, the passed
// sub-string will have at most ngc.Length runes. It is called once for every
// non-stripped rune in the string.
func (ngc NGramConfig) subGrams(idx, maxIdx int, ss []rune) []string {
	ngStrings := []string{}
	overflow := make([]rune, ngc.Length-1)

	maxSSIdx := min(idx, ngc.Length-1)
	minLength := ngc.MinLength
	if minLength == 0 {
		minLength = ngc.Length
	}

	for l := minLength; l <= ngc.Length; l++ {
		endIdx := idx % ngc.Length
		startIdx, ok := ngc.calcStartIdx(endIdx, maxSSIdx, l)
		if !ok {
			break
		}

		var str string

		if startIdx <= endIdx {
			str = string(ss[startIdx : endIdx+1])
		} else {
			str = string(ss[startIdx:maxSSIdx+1]) +
				string(ss[:endIdx+1])
		}

		remainder := l - len(str)
		if remainder > 0 {
			str = string(overflow[:remainder]) + str
		}

		ngStrings = append(ngStrings, str)
	}

	ngc.overflowStrings(idx, maxIdx, &ngStrings)

	return ngStrings
}

// NGrams constructs an NGram set according to the given configuration
func (ngc NGramConfig) NGrams(s string) NGramSet {
	ngs := NGramSet{}

	subStr := make([]rune, ngc.Length)
	maxIdx := len(s) - 1
	for i, r := range s {
		subStr[i%ngc.Length] = r
		strs := ngc.subGrams(i, maxIdx, subStr)
		for _, str := range strs {
			ngs[str]++
		}
	}

	return ngs
}
