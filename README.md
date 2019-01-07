# strdist
String distance metrics

## How to use the strdist package

There are two ways to use this package, you can either call the named
distance functions which have been defined for the various algorithms or else
you can use a `Finder`.

## Distance functions

There are `...Distance` functions provided for a number of algorithms and
these all take a pair of strings and return a measure of distance between
them. Note that the distances returned by different functions are not
comparible so, for instance, the CosineDistance and the Levenshtein distance
are on completely different scales.

```go
var s1 = "hello"
var s2 = "world"

fmt.Println("the Levenshtein distance between ",
            s1, " and ", s2, " is ",
			strdist.LevenshteinDistance(s1, s2))
```

## Finder

The `Finder` allows you to search a slice of strings for the closest match to
a target string. You can specify the shortest string to consider, the
threshold for similarity, whether or not to flatten case so that you can
treat "hello" and "HELLO" as the same string and finally you specify the
particular algorithm to be used. There are various default Finders
constructed for the supplied algorithms.
