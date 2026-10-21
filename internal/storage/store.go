package storage

import (
	"sync"
	"time"
)

// 데이터를 감싸는 내부 구조체 선언
type item struct {
	value string
	ttl   time.Time
}

// 메모리 스토리지 구조 선언
type MemStore struct {
	mu   sync.RWMutex
	data map[string]item
}

// 메모리 스토리지 생성
func NewMemStore() *MemStore {
	store := &MemStore{
		data: make(map[string]item),
	}
	go store.startExpirationCleaner()
	return store
}

// SET 명령어(데이터 저장) => TTL 설정 X
func (s *MemStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item{
		value: value,
		ttl:   time.Time{},
	}
}

// SETEX 명령어(데이터 저장) => TTL 설정 O
func (s *MemStore) SetEx(key, value string, ttlTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item{
		value: value,
		ttl:   time.Now().Add(ttlTime * time.Second),
	}
}

// GET 명령어(데이터 조회)
func (s *MemStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, exists := s.data[key]
	if !exists {
		return "", false
	}
	if !it.ttl.IsZero() && time.Now().After(it.ttl) {
		delete(s.data, key)
	}
	return it.value, exists
}

// DEL 명령어(데이터 삭제)
func (s *MemStore) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, exists := s.data[key]
	if !exists {
		return false
	}
	if !it.ttl.IsZero() && time.Now().After(it.ttl) {
		delete(s.data, key)
	}
	delete(s.data, key)
	return true
}

// EXISTS 명령어(데이터 존재 여부 확인)
func (s *MemStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, exists := s.data[key]
	if !exists {
		return false
	}
	if !it.ttl.IsZero() && time.Now().After(it.ttl) {
		delete(s.data, key)
	}
	return true
}

// KEYS 명령어 (데이터 전체 조회)
func (s *MemStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	keys := make([]string, 0, len(s.data))

	for k, it := range s.data {
		if !it.ttl.IsZero() && now.After(it.ttl) {
			delete(s.data, k)
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

// 1초마다 주기적으로 검사하여 TTL 시간 지난 키 지우기
func (s *MemStore) startExpirationCleaner() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for k, it := range s.data {
			if !it.ttl.IsZero() && now.After(it.ttl) {
				delete(s.data, k)
			}
		}
	}
}
