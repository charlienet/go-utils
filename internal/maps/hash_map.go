package maps

import (
	"maps"
	"sync"
)

type hashmap[M ~map[K]V, K comparable, V any] struct {
	m    M
	mu   sync.RWMutex
	sync bool
}

func NewHashMap[M ~map[K]V, K comparable, V any](mm ...M) *hashmap[M, K, V] {
	m := make(M)
	for _, v := range mm {
		maps.Copy(m, v)
	}

	return &hashmap[M, K, V]{
		m: m,
	}
}

func (h *hashmap[M, K, V]) Synchronize() *hashmap[M, K, V] {
	h.sync = true

	return h
}

func (h *hashmap[M, K, V]) Set(k K, v V) {
	h.mu.Lock()
	h.m[k] = v
	h.mu.Unlock()
}

func (h *hashmap[M, K, V]) Get(key K) (V, bool) {
	h.mu.RLock()

	v, ok := h.m[key]

	h.mu.RUnlock()
	return v, ok
}
