package storage

import "sync"

type HashStorage struct {
	mu     sync.Mutex
	Hashes map[string]string
}

func (h *HashStorage) AddHash(resHash string, filePath string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Hashes[filePath] = resHash
}
