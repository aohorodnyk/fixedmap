package fixedmap

import "testing"

func TestMap_Len(t *testing.T) {
	m := NewMapString[int](1)
	m.Set("key", 1)
	m.Set("key2", 2)
	m.Set("key3", 3)

	if m.Len() != 3 {
		t.Errorf("Expected map size is 3, got %d", m.Len())
	}

	m.Delete("key")

	if m.Len() != 2 {
		t.Errorf("Expected map size is 2, got %d", m.Len())
	}

	m.Delete("key2")
	m.Delete("key3")
	m.Delete("key3")
	m.Delete("key")
	m.Delete("key2")
	m.Delete("key5")

	if m.Len() != 0 {
		t.Errorf("Expected map size is 0, got %d", m.Len())
	}

	m.Set("key", 1)
	m.Set("key2", 2)
	m.Set("key3", 3)

	if m.Len() != 3 {
		t.Errorf("Expected map size is 3, got %d", m.Len())
	}
}
