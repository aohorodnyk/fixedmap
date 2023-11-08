package fixedmap

import (
	"sync/atomic"
)

func NewSyncMapString[V any](capacity int) *SyncMap[string, V] {
	return NewSyncMap[string, V](capacity, HashString, KeyCompareComparable[string], IndexUint64)
}

func NewSyncMapFlatByte[K KeyFlatByte, V any](capacity int) *SyncMap[K, V] {
	return NewSyncMap[K, V](capacity, HashFlatByte[K], KeyCompareComparable[K], IndexUint64)
}

func NewSyncMapBytes[V any](capacity int) *SyncMap[[]byte, V] {
	return NewSyncMap[[]byte, V](capacity, HashBytes, KeyCompareBytes, IndexUint64)
}

func NewSyncMap[K KeyType, V any](capacity int,
	keyHasher KeyHasher[K],
	keyComparator KeyComparator[K],
	indexCalculator IndexCalculator,
) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		table:           make([]atomic.Pointer[syncNode[K, V]], capacity),
		keyHasher:       keyHasher,
		keyComparator:   keyComparator,
		indexCalculator: indexCalculator,
	}
}

type SyncMap[K KeyType, V any] struct {
	table           []atomic.Pointer[syncNode[K, V]]
	keyHasher       KeyHasher[K]
	keyComparator   KeyComparator[K]
	indexCalculator IndexCalculator

	length atomic.Uint64
}

func (m *SyncMap[K, V]) Len() uint64 {
	if m == nil {
		return 0
	}

	return m.length.Load()
}

func (m *SyncMap[K, V]) Cap() int {
	if m == nil {
		return 0
	}

	return len(m.table)
}

func (m *SyncMap[K, V]) Get(key K) (value V, ok bool) {
	if m == nil {
		return value, false
	}

	index := m.index(key)
	current := m.table[index].Load()

	for current != nil && !m.keyComparator(current.key, key) {
		current = current.Next()
	}

	if current == nil {
		return value, false
	}

	// If at this point node is marked as deletedAt, we are ok, because we accessed it before it was deletedAt.
	// valueRef can never be nil, because we set it when we create a node.
	valueRef := current.value.Load()
	if valueRef.deleted {
		// Since we set new values (after delete) to the front, it makes no sense to search any further.
		return value, false
	}

	return valueRef.value, true
}

// Set sets a value by given key.
// It returns previous value and true if key already exists in FixedMap, otherwise it returns false.
func (m *SyncMap[K, V]) Set(key K, value V) (old V, ok bool) {
	if m == nil {
		return old, false
	}

	var (
		valueDefault V
		created      bool
		added        bool
		index        = m.index(key)
	)

	// Try to set new value until we succeed.
	for !added {
		// Reset old value, because we might try set node repeatedly.
		old = valueDefault
		// Reset created, because we might try set node repeatedly.
		created = false
		// Get current head of the linked list.
		head := m.table[index].Load()
		current := head

		// Find the first node with the same key.
		// Since we add new nodes to the front, we never need to search further than the first node with the same key.
		for current != nil && !m.keyComparator(current.key, key) {
			current = current.Next()
		}

		// If we did not find the node with the same key, let's try to create a new node.
		if current == nil {
			// Create new node and try to set it as a new head.
			added = m.insertNodeFront(index, head, key, value)
			created = true

			continue
		}

		// Get old node value.
		oldNodeValueRef := current.value.Load()
		if oldNodeValueRef.deleted {
			// If node is marked as deleted, someone deleted it.
			// Let's try to create a new node and set it as a new head.
			added = m.insertNodeFront(index, head, key, value)
			created = true

			continue
		}

		// Old value is not deleted. Let's try to set new value.
		nodeValue := newSyncNodeValue(value, false)
		added = current.value.CompareAndSwap(oldNodeValueRef, nodeValue)
		old = oldNodeValueRef.value
	}

	// If node was created, we need to increase length of the map.
	if created {
		m.length.Add(1)
	}

	return old, created
}

func (m *SyncMap[K, V]) Delete(key K) (old V, ok bool) {
	if m == nil {
		return old, false
	}

	var (
		valueDefault V
		deleted      bool
		index        = m.index(key)
	)

	// Try to delete node until we succeed.
	for !deleted {
		var prev *syncNode[K, V]

		old = valueDefault

		head := m.table[index].Load()
		current := head

		// Find the first node with the same key.
		for current != nil && !m.keyComparator(current.key, key) {
			current = current.Next()
		}

		// If we did not find the node with the same key, we are done.
		if current == nil {
			return old, false
		}

		// Get old node value.
		oldNodeValueRef := current.value.Load()
		if oldNodeValueRef.deleted {
			// If node is marked as deleted, someone deleted it.
			// We are done.
			return old, false
		}

		// Create new node value and use CAS to set it as a new value.
		// If CAS fails, it means that someone else already changed the value.
		nodeValue := newSyncNodeValue(oldNodeValueRef.value, true)
		if !current.value.CompareAndSwap(oldNodeValueRef, nodeValue) {
			continue
		}

		old = oldNodeValueRef.value
		// After set deleted flag, we need to delete the node from the linked list.
		// If our current node is head, we need to set next node as a new head.
		if head == current {
			// Try to set next node as a new head.
			deleted = m.table[index].CompareAndSwap(head, current.Next())

			continue
		}

		deleted = prev.next.CompareAndSwap(current, current.Next())
	}

	m.length.Add(^uint64(0))

	return old, true
}

func (m *SyncMap[K, V]) Range(callback Ranger[K, V]) {
	if m == nil {
		return
	}

	for index := range m.table {
		current := m.table[index].Load()

		for current != nil {
			valueRef := current.value.Load()
			if !valueRef.deleted && !callback(current.key, valueRef.value) {
				return
			}

			current = current.Next()
		}
	}
}

func (m *SyncMap[K, V]) insertNodeFront(index int, head *syncNode[K, V], key K, value V) bool {
	// We need to create a new node and add it to the front of the LinkedList.
	newNode := newSyncNode(key, value)
	// Set current head to the next of the new node.
	newNode.next.Store(head)
	// Try to set new node as head.
	return m.table[index].CompareAndSwap(head, newNode)
}

// index calculates index in table by given key.
func (m *SyncMap[K, V]) index(key K) int {
	return m.indexCalculator(m.keyHasher(key), len(m.table))
}
