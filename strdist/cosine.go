package strdist

import (
	"math"
)

// cosineStrDetails holds details about a previously scanned string
type cosineStrDetails struct {
	ngs    NGramSet
	strLen float64
}

// CosineAlgo encapsulates the details needed to provide the cosine
// distance. Note that the cosine distance is not a true distance metric as
// it does not exhibit the triangle inequality property.
type CosineAlgo struct {
	ngc   NGramConfig
	cache *cache[cosineStrDetails]
}

// NewCosineAlgo returns a new CosineAlgo with the config and cache size set
func NewCosineAlgo(ngc NGramConfig, maxCacheSize int) (*CosineAlgo, error) {
	if err := ngc.Check(); err != nil {
		return nil, err
	}

	cache, err := newCache[cosineStrDetails](maxCacheSize)
	if err != nil {
		return nil, err
	}

	return &CosineAlgo{
		ngc:   ngc,
		cache: cache,
	}, nil
}

// NewCosineAlgoOrPanic returns a new CosineAlgo. it will panic if the algo
// cannot be created withour errors.
func NewCosineAlgoOrPanic(ngc NGramConfig, maxCacheSize int) *CosineAlgo {
	a, err := NewCosineAlgo(ngc, maxCacheSize)
	if err != nil {
		panic(err)
	}
	return a
}

// Dist for a CosineAlgo will calculate the distance from the target string
func (a *CosineAlgo) Dist(s1, s2 string) float64 {
	sd1 := a.getStrDetails(s1)
	sd2 := a.getStrDetails(s2)
	return a.cosineDistance(sd1, sd2)
}

// Name returns the algorithm Name
func (a CosineAlgo) Name() string {
	return AlgoNameCosine
}

// Desc returns a string describing the algorithm configuration
func (a CosineAlgo) Desc() string {
	return a.cache.Desc() + " " + a.ngc.Desc("N-Gram:")
}

// getStrDetails returns the strDetails associated with the given string and
// caaches the results if the cache is of non-zero size.
func (a *CosineAlgo) getStrDetails(s string) cosineStrDetails {
	if csd, ok := a.cache.getCachedEntry(s); ok {
		return csd
	}
	var csd cosineStrDetails
	csd.ngs = a.ngc.NGrams(s)
	csd.strLen = csd.ngs.Length()

	a.cache.setCachedEntry(s, csd)

	return csd
}

// cosineDistance returns the cosine distance for the two strings given their
// NGrams and lengths
func (a CosineAlgo) cosineDistance(sd1, sd2 cosineStrDetails) float64 {
	if sd1.strLen == 0 && sd2.strLen == 0 {
		return 0
	}
	if sd1.strLen == 0 || sd2.strLen == 0 {
		return 1
	}
	dot := Dot(sd1.ngs, sd2.ngs)
	if dot == 0 {
		return 1
	}
	return 1 - (float64(dot) / (sd1.strLen * sd2.strLen))
}

// CosineSimilarity returns the cosine similarity between two NGramSets ngs1
// and ngs2. This is the dot-product of the two n-gram sets divided by
// the product of the lengths. Note that if both sets are of length 0 then
// the similarity is set to 1 (meaning identical) but if either one but not
// both is of length 0 (or the dot-product is zero) then the similarity is
// set to 0 (meaning completely different).
func CosineSimilarity(ngs1, ngs2 NGramSet) float64 {
	len2s1 := ngs1.lengthSquared()
	len2s2 := ngs2.lengthSquared()

	if len2s1 == 0.0 && len2s2 == 0.0 {
		return 1.0
	}
	if len2s1 == 0.0 || len2s2 == 0.0 {
		return 0.0
	}

	d := Dot(ngs1, ngs2)
	if d == 0.0 {
		return 0.0
	}

	lenS1 := math.Sqrt(len2s1)
	lenS2 := math.Sqrt(len2s2)

	return float64(d) / (lenS1 * lenS2)
}
