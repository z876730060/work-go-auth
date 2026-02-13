package utils

import (
	"slices"
	"sync"
)

type LinkMap[T any] struct {
	mux      sync.RWMutex
	maps     map[string]T
	sort     []string
	override bool
	reSort   bool
	safe     bool
}

func NewLinkMap[T any](override, reSort, safe bool) *LinkMap[T] {
	return &LinkMap[T]{
		maps:     make(map[string]T),
		sort:     make([]string, 0, 16),
		override: override,
		reSort:   reSort,
		safe:     safe,
	}
}

func (l *LinkMap[T]) Get(key string) (T, bool) {
	if l.safe {
		l.mux.RLock()
		defer l.mux.RUnlock()
	}
	v, ok := l.maps[key]
	return v, ok
}

func (l *LinkMap[T]) Set(key string, value T) {
	if l.safe {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	if _, ok := l.maps[key]; ok {
		if l.override {
			l.maps[key] = value
			if l.reSort {
				index := slices.Index(l.sort, key)
				if index == -1 {
					l.sort = append(l.sort, key)
				} else {
					l.sort = append(l.sort[:index], l.sort[index+1:]...)
					l.sort = append(l.sort, key)
				}
			}
		}
		return
	}
	l.maps[key] = value
	l.sort = append(l.sort, key)
}

func (l *LinkMap[T]) Range(f func(key string, value T) bool) {
	if l.safe {
		l.mux.RLock()
		defer l.mux.RUnlock()
	}
	for _, key := range l.sort {
		if !f(key, l.maps[key]) {
			break
		}
	}
}

func (l *LinkMap[T]) Delete(key string) {
	if l.safe {
		l.mux.Lock()
		defer l.mux.Unlock()
	}
	delete(l.maps, key)
	index := slices.Index(l.sort, key)
	if index == -1 {
		return
	}
	l.sort = append(l.sort[:index], l.sort[index+1:]...)
}
