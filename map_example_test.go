package fixedmap_test

import (
	"fmt"

	"github.com/aohorodnyk/fixedmap"
)

func ExampleNewMapBytes() {
	mapBytes := fixedmap.NewMapBytes[string](10)
	mapBytes.Set([]byte{123, 35, 12}, "bytes: [123, 35, 12]")
	mapBytes.Set([]byte{}, "bytes2: []")
	mapBytes.Set([]byte{0}, "bytes3: [0]")

	fmt.Println(mapBytes.Get([]byte{123, 35, 12}))
	fmt.Println(mapBytes.Get([]byte{}))
	fmt.Println(mapBytes.Get([]byte{0}))

	mapBytes.Delete([]byte{123, 35, 12})

	mapBytes.Set([]byte{0}, "bytes3: {0}")

	fmt.Println(mapBytes.Get([]byte{123, 35, 12}))
	fmt.Println(mapBytes.Get([]byte{}))
	fmt.Println(mapBytes.Get([]byte{0}))

	// Output:
	// bytes: [123, 35, 12] true
	// bytes2: [] true
	// bytes3: [0] true
	//  false
	// bytes2: [] true
	// bytes3: {0} true
}

func ExampleNewMapFlatByte() {
	mapInt := fixedmap.NewMapFlatByte[int, string](10)
	mapInt.Set(213, "two hundred three")
	mapInt.Set(267342463255, "Huge number")
	mapInt.Set(1, "one")
	mapInt.Set(0, "zero")
	mapInt.Set(5, "five")

	fmt.Println(mapInt.Get(267342463255))
	fmt.Println(mapInt.Get(1))

	mapInt.Delete(1)

	fmt.Println(mapInt.Get(267342463255))
	fmt.Println(mapInt.Get(1))

	mapInt.Set(267342463255, "267342463255")
	fmt.Println(mapInt.Get(267342463255))

	// Output:
	// Huge number true
	// one true
	// Huge number true
	//  false
	// 267342463255 true
}

func ExampleNewMapString() {
	mapStr := fixedmap.NewMapString[int](10)
	mapStr.Set("key", 1)
	mapStr.Set("key2", 2)
	mapStr.Set("key3", 3)

	value, ok := mapStr.Get("key")
	fmt.Println(value, ok)

	value, ok = mapStr.Get("key3")
	fmt.Println(value, ok)

	value, ok = mapStr.Get("key4")
	fmt.Println(value, ok)

	mapStr.Delete("key")

	value, ok = mapStr.Get("key")
	fmt.Println(value, ok)

	value, ok = mapStr.Get("key2")
	fmt.Println(value, ok)

	value, ok = mapStr.Get("key3")
	fmt.Println(value, ok)

	value, ok = mapStr.Get("key4")
	fmt.Println(value, ok)

	// Output:
	// 1 true
	// 3 true
	// 0 false
	// 0 false
	// 2 true
	// 3 true
	// 0 false
}
