package fixedmap

import (
	"math"
)

type IndexCalculator func(hash uint64) int

// IndexFactory returns an index calculator using different algorithms.
// In case of size that is equal to 2^n, it will use an algorithm: checksum & (size - 1).
// In case of size that is not 2^n, it will use an algorithms: checksum % size.
// Our benchmark tests show that in Xeon processors in a cloud bitwise operator slightly faster.
// If 2^n size number is preferable, Closest2n function can be used to get the closest 2^n number to input size.
func IndexFactory(size int) (callback IndexCalculator) {
	sizeU := uint64(size)
	{
		nf := math.Log2(float64(size))
		if nf != math.Trunc(nf) {
			callback = func(hash uint64) int {
				return int(hash % sizeU)
			}

			return callback
		}
	}

	sizeU--

	callback = func(hash uint64) int {
		return int(hash & sizeU)
	}

	return callback
}

// Closest2n returns the next closest number to input size from 2^n.
// If number is 2^n the same number will be returned.
// Otherwise, this function will return bigger number than input.
func Closest2n(size int) int {
	if size <= 0 {
		return 1
	}

	nf := math.Log2(float64(size))
	if nf == math.Trunc(nf) {
		return size
	}

	return 1 << (int(nf) + 1)
}
