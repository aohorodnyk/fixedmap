package fixedmap

import (
	"hash/maphash"
	"unsafe"
)

// KeyType is a type that can be used as a key in maps.
type KeyType interface {
	KeyFlatByte | ~string | []byte
}

// KeyFlatByte is a type that can be used as a key in maps.
// It can be converted to a byte slice with unsafe.Slice with the same code.
type KeyFlatByte interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~uint | ~float32 | ~float64 | ~complex64 | ~complex128
}

// KeyHasher is a function type that is used to hash keys.
type KeyHasher[K KeyType] func(key K) uint64

// maphash uses hash/maphash algorithm to calculate seed.
// maphash ~23% faster than xxhash and ~35% faster than crc32.

//nolint:gochecknoglobals // This seed will be used for getting hash in FixedMap.
var mapHashSeed = maphash.MakeSeed()

// MapHashString calculates hash based on string input.
func HashString(key string) uint64 {
	return maphash.String(mapHashSeed, key)
}

// MapHashBytes calculates hash based on slice of bytes input.
func HashBytes(key []byte) uint64 {
	return maphash.Bytes(mapHashSeed, key)
}

// MapHashFlatByte calculates hash based on flat types input.
func HashFlatByte[K KeyFlatByte](key K) uint64 {
	return maphash.Bytes(mapHashSeed, unsafe.Slice((*byte)(unsafe.Pointer(&key)), unsafe.Sizeof(key)))
}
