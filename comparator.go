package fixedmap

import "bytes"

// KeyComparator is a function type that is used to compare keys.
type KeyComparator[K KeyType] func(K, K) bool

// KeyCompareComparable is a compare function for comparable type.
func KeyCompareComparable[K comparable](left, right K) bool {
	return left == right
}

// KeyCompareBytes is a compare function for comparable type.
func KeyCompareBytes(left, right []byte) bool {
	return bytes.Compare(left, right) == 0
}
