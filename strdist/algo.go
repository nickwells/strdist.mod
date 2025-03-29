package strdist

// Algo describes the interface that a distance algorithm used by the
// Finder must satisfy
type Algo interface {
	Dist(s1, s2 string) float64
	Name() string
	Desc() string
}

// AlgoName consts are names of various distance algorithms
const (
	AlgoNameLevenshtein       = "Levenshtein"
	AlgoNameScaledLevenshtein = "scaled Levenshtein"
	AlgoNameCosine            = "cosine"
	AlgoNameHamming           = "Hamming"
	AlgoNameJaccard           = "Jaccard"
	AlgoNameWeightedJaccard   = "weighted Jaccard"

	caseBlind = "case-blind "

	CaseBlindAlgoNameLevenshtein       = caseBlind + AlgoNameLevenshtein
	CaseBlindAlgoNameScaledLevenshtein = caseBlind + AlgoNameScaledLevenshtein
	CaseBlindAlgoNameCosine            = caseBlind + AlgoNameCosine
	CaseBlindAlgoNameHamming           = caseBlind + AlgoNameHamming
	CaseBlindAlgoNameJaccard           = caseBlind + AlgoNameJaccard
	CaseBlindAlgoNameWeightedJaccard   = caseBlind + AlgoNameWeightedJaccard
)

// DfltThreshold consts are suggested default similarity thresholds for the
// Finder when using the corresponding distance algorithms
const (
	DfltThresholdLevenshtein     = 5.0
	DfltThresholdScaledLev       = 0.33
	DfltThresholdCosine          = 0.4
	DfltThresholdHamming         = 5.0
	DfltThresholdJaccard         = 0.5
	DfltThresholdWeightedJaccard = 0.7
)

// DefaultThresholds associates the default similarity thresholds with the
// algorithm names
var DefaultThresholds = map[string]float64{
	AlgoNameLevenshtein:       DfltThresholdLevenshtein,
	AlgoNameScaledLevenshtein: DfltThresholdScaledLev,
	AlgoNameCosine:            DfltThresholdCosine,
	AlgoNameHamming:           DfltThresholdHamming,
	AlgoNameJaccard:           DfltThresholdJaccard,
	AlgoNameWeightedJaccard:   DfltThresholdWeightedJaccard,
}

// DefaultFinders associates the Finders with the algorithm name
var DefaultFinders = map[string]*Finder{
	AlgoNameLevenshtein: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdLevenshtein,
			MinStrLength: DfltMinStrLength,
		},
		LevenshteinAlgo{}),
	AlgoNameScaledLevenshtein: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdScaledLev,
			MinStrLength: DfltMinStrLength,
		},
		ScaledLevAlgo{}),
	AlgoNameCosine: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdCosine,
			MinStrLength: DfltMinStrLength,
		},
		NewCosineAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),
	AlgoNameHamming: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdHamming,
			MinStrLength: DfltMinStrLength,
		},
		HammingAlgo{}),
	AlgoNameJaccard: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdJaccard,
			MinStrLength: DfltMinStrLength,
		},
		NewJaccardAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),
	AlgoNameWeightedJaccard: NewFinderOrPanic(
		FinderConfig{
			Threshold:    DfltThresholdJaccard,
			MinStrLength: DfltMinStrLength,
		},
		NewWeightedJaccardAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),

	CaseBlindAlgoNameLevenshtein: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdLevenshtein,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		LevenshteinAlgo{}),
	CaseBlindAlgoNameScaledLevenshtein: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdScaledLev,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		ScaledLevAlgo{}),
	CaseBlindAlgoNameCosine: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdCosine,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		NewCosineAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),
	CaseBlindAlgoNameHamming: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdHamming,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		HammingAlgo{}),
	CaseBlindAlgoNameJaccard: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdJaccard,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		NewJaccardAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),
	CaseBlindAlgoNameWeightedJaccard: NewFinderOrPanic(
		FinderConfig{
			Threshold:      DfltThresholdJaccard,
			MapToLowerCase: true,
			MinStrLength:   DfltMinStrLength,
		},
		NewWeightedJaccardAlgoOrPanic(DfltNGramConfig, DfltMaxCacheSize)),
}
