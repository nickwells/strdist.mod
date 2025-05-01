package strdist

import "fmt"

// DfltMinStrLength is the default value for the minimum string length
const DfltMinStrLength = 3

// FinderConfig contains details needed to configure an algorithm and to
// constrain the subsequent Finder.
type FinderConfig struct {
	// Threshold limits the distance that a string must have for the Finder
	// to recognise it as a match. The distance must be <= the threshold
	Threshold float64
	// MinStrLength is the shortest string that will be compared against
	// other strings. The problem with trying to find similar strings to very
	// short targets is that they can match with a lot of not obviously
	// similar alternatives. For instance a match for a single character
	// string might be every other single character string in the
	// population. For a number of use cases this is not particularly
	// helpful.
	MinStrLength int
	// MapToLowerCase is set to indicate that the string should be mapped to
	// a lower-case equivalent before calculating the distance.
	MapToLowerCase bool
	// StripRunes is a set of runes (unicode characters in a string) to be
	// removed from the string before calculating the distance.
	StripRunes string
}

// Check checks that the FinderConfig has valid values and returns an error
// if not
func (fc FinderConfig) Check() error {
	if fc.Threshold < 0 {
		return fmt.Errorf("FinderConfig: the Threshold (%f) must be >= 0",
			fc.Threshold)
	}

	if fc.MinStrLength < 0 {
		return fmt.Errorf(
			"FinderConfig: the minimum string length (%d) must be >= 0",
			fc.MinStrLength)
	}

	return nil
}

// Desc returns a string describing the finder configuration
func (fc FinderConfig) Desc() string {
	s := fmt.Sprintf("Threshold: %7.4f", fc.Threshold)
	s += fmt.Sprintf(", MinStrLength: %2d", fc.MinStrLength)
	s += fmt.Sprintf(", MapToLowerCase: %-5.5v", fc.MapToLowerCase)
	s += fmt.Sprintf(", StripRunes: %-9q", fc.StripRunes)

	return s
}
