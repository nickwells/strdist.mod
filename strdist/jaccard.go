package strdist

// JaccardAlgo encapsulates the details needed to provide the jaccard distance.
type JaccardAlgo struct {
	ngc   NGramConfig
	cache *cache[NGramSet]
}

// NewJaccardAlgo returns a new JaccardAlgo with the config and cache size set
func NewJaccardAlgo(ngc NGramConfig, maxCacheSize int) (*JaccardAlgo, error) {
	if err := ngc.Check(); err != nil {
		return nil, err
	}

	cache, err := newCache[NGramSet](maxCacheSize)
	if err != nil {
		return nil, err
	}

	return &JaccardAlgo{
		ngc:   ngc,
		cache: cache,
	}, nil
}

// NewJaccardAlgoOrPanic returns a new JaccardAlgo. it will panic if the algo
// cannot be created withour errors.
func NewJaccardAlgoOrPanic(ngc NGramConfig, maxCacheSize int) *JaccardAlgo {
	a, err := NewJaccardAlgo(ngc, maxCacheSize)
	if err != nil {
		panic(err)
	}
	return a
}

// Name returns the algorithm name
func (JaccardAlgo) Name() string {
	return AlgoNameJaccard
}

// Desc returns a string describing the algorithm configuration
func (a JaccardAlgo) Desc() string {
	return a.cache.Desc() + " " + a.ngc.Desc("N-Gram:")
}

// getNGramSet returns the strDetails associated with the given string and
// caaches the results if the cache is of non-zero size.
func (a *JaccardAlgo) getNGramSet(s string) NGramSet {
	if ngs, ok := a.cache.getCachedEntry(s); ok {
		return ngs
	}

	ngs := a.ngc.NGrams(s)
	a.cache.setCachedEntry(s, ngs)

	return ngs
}

// Dist for a JaccardAlgo will calculate the distance from the target string
func (a *JaccardAlgo) Dist(s1, s2 string) float64 {
	ngs1 := a.getNGramSet(s1)
	ngs2 := a.getNGramSet(s2)

	return 1.0 - JaccardIndex(ngs1, ngs2)
}

// JaccardIndex returns the Jaccard index of the two n-gram sets
func JaccardIndex(ngs1, ngs2 NGramSet) float64 {
	if len(ngs1) == 0 && len(ngs2) == 0 {
		return 1.0
	}

	uLen := NGramLenUnion(ngs1, ngs2)
	iLen := NGramLenIntersection(ngs1, ngs2)

	return float64(iLen) / float64(uLen)
}
