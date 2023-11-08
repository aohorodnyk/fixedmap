package fixedmap

// NewMapString creates a new map with string keys.
// The map is backed by a hash table with the given capacity.
func NewMapString[V any](cap int) *Map[string, V] {
	return NewMap[string, V](cap, HashString, KeyCompareComparable[string], IndexUint64)
}

// NewMapFlatByte creates a new map with keys that are generic type from KeyFlatByte.
// The map is backed by a hash table with the given capacity.
func NewMapFlatByte[K KeyFlatByte, V any](cap int) *Map[K, V] {
	return NewMap[K, V](cap, HashFlatByte[K], KeyCompareComparable[K], IndexUint64)
}

// NewMapBytes creates a new map with []byte keys.
// The map is backed by a hash table with the given capacity.
func NewMapBytes[V any](cap int) *Map[[]byte, V] {
	return NewMap[[]byte, V](cap, HashBytes, KeyCompareBytes, IndexUint64)
}

// NewMap creates a new map with all custom parameters.
func NewMap[K KeyType, V any](cap int, keyHasher KeyHasher[K],
	keyComparator KeyComparator[K],
	indexCalculator IndexCalculator) *Map[K, V] {
	return &Map[K, V]{
		table:           make([]*node[K, V], cap),
		keyHasher:       keyHasher,
		keyComparator:   keyComparator,
		indexCalculator: indexCalculator,
	}
}

// Map is a hash table with fixed capacity.
// The map is backed by a hash table with the given capacity and uses LinkedList to resolve collisions.
// The map uses the given KeyHasher to hash keys and the given KeyComparator to compare keys.
// The map is NOT safe for concurrent use.
type Map[K KeyType, V any] struct {
	table           []*node[K, V]
	keyHasher       KeyHasher[K]
	keyComparator   KeyComparator[K]
	indexCalculator IndexCalculator

	length int
}

// node is a node of LinkedList.
// It's used to resolve collisions in hash table.
type node[K KeyType, V any] struct {
	key   K
	value V
	next  *node[K, V]
}

// Len returns a number of elements stored in FixedMap at given time.
func (m *Map[K, V]) Len() int {
	return m.length
}

// Cap returns size of pre-allocated table to store data.
func (m *Map[K, V]) Cap() int {
	return len(m.table)
}

// Get returns a value by given key.
// The second parameter returns true if value found in FixedMap, otherwise it returns false.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	if m == nil {
		return value, false
	}

	index := m.index(key)
	current := m.table[index]

	for current != nil && !m.keyComparator(current.key, key) {
		current = current.next
	}

	if current == nil {
		return value, false
	}

	return current.value, true
}

// Set sets a value by given key.
// If key already exists in FixedMap, it will be overwritten.
// If key does not exist in FixedMap, it will be added to a front of a LinkedList.
func (m *Map[K, V]) Set(key K, value V) {
	if m == nil {
		return
	}

	index := m.index(key)
	current := m.table[index]

	for current != nil && !m.keyComparator(current.key, key) {
		current = current.next
	}

	if current == nil {
		m.table[index] = &node[K, V]{key, value, m.table[index]}
	} else {
		current.value = value
	}

	m.length++
}

// Delete deletes a value by given key.
func (m *Map[K, V]) Delete(key K) {
	if m == nil {
		return
	}

	index := m.index(key)
	current := m.table[index]
	var previous *node[K, V]

	for current != nil && !m.keyComparator(current.key, key) {
		previous = current
		current = current.next
	}

	if current == nil {
		return
	}

	if previous == nil {
		m.table[index] = current.next
	} else {
		previous.next = current.next
	}

	m.length--
}

// Range iterates over all elements in FixedMap.
// If the given function returns false, the iteration will be stopped.
func (m *Map[K, V]) Range(callback Ranger[K, V]) {
	if m == nil {
		return
	}

	for _, current := range m.table {
		for current != nil {
			if !callback(current.key, current.value) {
				return
			}

			current = current.next
		}
	}
}

// index calculates index in table by given key.
func (m *Map[K, V]) index(key K) int {
	return m.indexCalculator(m.keyHasher(key), len(m.table))
}
