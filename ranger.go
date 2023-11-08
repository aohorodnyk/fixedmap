package fixedmap

// Ranger is a function that is called for each element in FixedMap in Range function.
type Ranger[K any, V any] func(key K, value V) (next bool)
