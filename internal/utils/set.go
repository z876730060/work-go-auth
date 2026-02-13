package utils

import (
	"slices"
	"sync"
)

type Set[T comparable] struct {
	mux   sync.RWMutex
	datas []T
	safe  bool
}

func NewSet[T comparable](safe bool) *Set[T] {
	return &Set[T]{
		datas: make([]T, 0, 16),
		safe:  safe,
	}
}

func (s *Set[T]) Add(value T) {
	if s.safe {
		s.mux.Lock()
		defer s.mux.Unlock()
	}
	if slices.Contains(s.datas, value) {
		return
	}
	s.datas = append(s.datas, value)
}

func (s *Set[T]) Range(f func(value T) bool) {
	if s.safe {
		s.mux.RLock()
		defer s.mux.RUnlock()
	}
	for _, v := range s.datas {
		if !f(v) {
			break
		}
	}
}
