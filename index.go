package fixedmap

type IndexCalculator func(hash uint64, len int) int

// IndexUint64 is an index calculator that assumes that we use checksum calculator x64.
func IndexUint64(sum uint64, length int) int {
	return int(sum % uint64(length))
}

// IndexUint32 is an index calculator that assumes that we use checksum calculator x32.
// This calculator can be ONLY used with map cap less than MaxUint32 and checksum calculator x32.
// It can be much faster in some use cases.
func IndexUint32(sum uint64, length int) int {
	return int(uint32(sum) % uint32(length))
}
