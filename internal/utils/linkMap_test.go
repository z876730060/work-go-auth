package utils

import "testing"

func TestLinkMap_Set(t *testing.T) {
	l := NewLinkMap[int](false, false, true)
	l.Set("a", 1)
	l.Set("b", 2)
	l.Set("c", 3)
	if v, ok := l.Get("a"); !ok || v != 1 {
		t.Errorf("a = %v, want 1", v)
	}
	if v, ok := l.Get("b"); !ok || v != 2 {
		t.Errorf("b = %v, want 2", v)
	}
	if v, ok := l.Get("c"); !ok || v != 3 {
		t.Errorf("c = %v, want 3", v)
	}
	l.Set("a", 4)
	if v, ok := l.Get("a"); !ok || v != 4 {
		t.Errorf("a = %v, want 4", v)
	}
	t.Log("range:")
	l.Range(func(key string, value int) bool {
		t.Log(key, value)
		return true
	})

	t.Log(len(l.sort), cap(l.sort))
	l.Delete("a")
	if v, ok := l.Get("a"); ok || v != 0 {
		t.Errorf("a = %v, want 0", v)
	}
	t.Log(len(l.sort), cap(l.sort))
}
