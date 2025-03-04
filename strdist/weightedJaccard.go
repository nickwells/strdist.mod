package strdist

// WeightedJaccardAlgo encapsulates the details needed to provide the
// WeightedJaccard distance.
type WeightedJaccardAlgo struct {
	*JaccardAlgo
}

// NewWeightedJaccardAlgo returns a new WeightedJaccardAlgo with the config
// and cache size set
func NewWeightedJaccardAlgo(ngc NGramConfig, maxCacheSize int) (
	*WeightedJaccardAlgo, error,
) {
	ja, err := NewJaccardAlgo(ngc, maxCacheSize)
	if err != nil {
		return nil, err
	}

	return &WeightedJaccardAlgo{JaccardAlgo: ja}, nil
}

// NewWeightedJaccardAlgoOrPanic returns a new WeightedJaccardAlgo. it will
// panic if the algo cannot be created withour errors.
func NewWeightedJaccardAlgoOrPanic(ngc NGramConfig, maxCacheSize int,
) *WeightedJaccardAlgo {
	a, err := NewWeightedJaccardAlgo(ngc, maxCacheSize)
	if err != nil {
		panic(err)
	}

	return a
}

// Name returns the algorithm name
func (WeightedJaccardAlgo) Name() string {
	return AlgoNameWeightedJaccard
}

// Desc returns a string describing the algorithm configuration
func (a WeightedJaccardAlgo) Desc() string {
	return a.JaccardAlgo.Desc()
}

// Dist for a WeightedJaccardAlgo will calculate the distance from the target
// string
func (a *WeightedJaccardAlgo) Dist(s1, s2 string) float64 {
	ngs1 := a.getNGramSet(s1)
	ngs2 := a.getNGramSet(s2)

	return 1.0 - WeightedJaccardIndex(ngs1, ngs2)
}

// WeightedJaccardIndex returns the Weighted Jaccard index of the two n-gram
// sets. It uses the NGramWeightedLen... functions to calculate the length
func WeightedJaccardIndex(ngs1, ngs2 NGramSet) float64 {
	if len(ngs1) == 0 && len(ngs2) == 0 {
		return 1.0
	}

	uLen := NGramWeightedLenUnion(ngs1, ngs2)
	iLen := NGramWeightedLenIntersection(ngs1, ngs2)

	return float64(iLen) / float64(uLen)
}
