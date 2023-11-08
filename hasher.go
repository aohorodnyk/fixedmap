package fixedmap

import (
	"unsafe"

	"github.com/cespare/xxhash/v2"
)

// KeyType is a type that can be used as a key in maps.
type KeyType interface {
	KeyFlatByte | ~string | []byte
}

// KeyFlatByte is a type that can be used as a key in maps.
// It can be converted to a byte slice with unsafe.Slice with the same code.
type KeyFlatByte interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 |
		~int | ~uint | ~float32 | ~float64 | ~complex64 | ~complex128
}

// KeyHasher is a function type that is used to hash keys.
type KeyHasher[K KeyType] func(key K) uint64

// HashString is a hash function for string type.
func HashString(key string) uint64 {
	return xxhash.Sum64(unsafe.Slice(unsafe.StringData(key), len(key)))
}

// HashFlatByte is a hash function for KeyNumbers types.
// It converts the key to a byte slice with unsafe.Slice.
func HashFlatByte[K KeyFlatByte](key K) uint64 {
	return xxhash.Sum64(unsafe.Slice((*byte)(unsafe.Pointer(&key)), unsafe.Sizeof(key)))
}

// HashBytes is a hash function for byte slice type.
func HashBytes(key []byte) uint64 {
	return xxhash.Sum64(key)
}
