package utils

import "testing"

func TestAll(t *testing.T) {
	t.Run("Add", TestSet_Add)
	t.Run("Range", TestSet_Range)
}

func TestSet_Add(t *testing.T) {
	s := NewSet[int](false)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	if len(s.datas) != 3 {
		t.Errorf("len(s.datas) = %v, want 3", len(s.datas))
	}
}

func TestSet_Range(t *testing.T) {
	s := NewSet[int](false)
	s.Add(1)
	s.Add(2)
	s.Add(3)
	var sum int
	s.Range(func(value int) bool {
		sum += value
		return true
	})
	if sum != 6 {
		t.Errorf("sum = %v, want 6", sum)
	}
}
