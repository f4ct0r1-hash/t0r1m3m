package storage

import (
	"sync"
)

// 메모리 스토리지 구조 선언
type MemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// 메모리 스토리지 생성
func NewMemStore() *MemStore {
	return &MemStore{
		data: make(map[string]string),
	}
}

// SET 명령어(데이터 저장)
func (s *MemStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// GET 명령어(데이터 조회)
func (s *MemStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.data[key]
	return val, exists
}
