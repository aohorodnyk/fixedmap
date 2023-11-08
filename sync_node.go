package fixedmap

import (
	"sync/atomic"
)

func newSyncNode[K KeyType, V any](key K, value V) *syncNode[K, V] {
	newNode := &syncNode[K, V]{
		key: key,
	}

	newNode.value.Store(newSyncNodeValue(value, false))

	return newNode
}

func newSyncNodeValue[V any](value V, deleted bool) *syncNodeValue[V] {
	return &syncNodeValue[V]{
		value:   value,
		deleted: deleted,
	}
}

// syncNode is a node in a linked list.
// It's used to resolve collisions in hash table.
// It's safe for concurrent use.
// It uses atomic operations to make it lock-free.
type syncNode[K KeyType, V any] struct {
	key   K
	value atomic.Pointer[syncNodeValue[V]]

	next atomic.Pointer[syncNode[K, V]]
}

// NodeValue is a value of the node.
// This struct contains different states of value, like deleted.
// It uses atomic operations to make it lock-free.
type syncNodeValue[V any] struct {
	value   V
	deleted bool
}

// Next function is a helper to check if the node is deletedAt.
// It's safe for concurrent use.
// It's created to make access to the next field with worry-free about nil value.
func (n *syncNode[K, V]) Next() *syncNode[K, V] {
	if n == nil {
		return nil
	}

	return n.next.Load()
}
