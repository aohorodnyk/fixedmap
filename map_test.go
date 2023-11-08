package fixedmap_test

import (
	"testing"

	"github.com/aohorodnyk/fixedmap"
)

func TestMap_Len(t *testing.T) {
	t.Parallel()

	fmap := fixedmap.NewMapString[int](1)
	fmap.Set("key", 1)
	fmap.Set("key2", 2)
	fmap.Set("key3", 3)

	if fmap.Len() != 3 {
		t.Errorf("Expected map size is 3, got %d", fmap.Len())
	}

	fmap.Delete("key")

	if fmap.Len() != 2 {
		t.Errorf("Expected map size is 2, got %d", fmap.Len())
	}

	fmap.Delete("key2")
	fmap.Delete("key3")
	fmap.Delete("key3")
	fmap.Delete("key")
	fmap.Delete("key2")
	fmap.Delete("key5")

	if fmap.Len() != 0 {
		t.Errorf("Expected map size is 0, got %d", fmap.Len())
	}

	fmap.Set("key", 1)
	fmap.Set("key2", 2)
	fmap.Set("key3", 3)

	if fmap.Len() != 3 {
		t.Errorf("Expected map size is 3, got %d", fmap.Len())
	}
}
