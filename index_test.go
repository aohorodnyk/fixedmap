package fixedmap_test

import (
	"fmt"
	"testing"

	"github.com/aohorodnyk/fixedmap"
)

func TestIndexFactory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		size int
		is2n bool
		hash uint64
		exp  int
	}{
		{
			size: 5,
			is2n: false,
			hash: 4,
			exp:  4,
		},
		{
			size: 5,
			is2n: false,
			hash: 9,
			exp:  4,
		},
		{
			size: 128,
			is2n: true,
			hash: 255,
			exp:  127,
		},
		{
			size: 1_024,
			is2n: true,
			hash: 42_342,
			exp:  358,
		},
		{
			size: 1_000,
			is2n: false,
			hash: 57_632,
			exp:  632,
		},
		{
			size: 536_870_912,
			is2n: true,
			hash: 57_632,
			exp:  57_632,
		},
		{
			size: 536_870_912,
			is2n: true,
			hash: 539_276_951,
			exp:  2_406_039,
		},
		{
			size: 500_000_000,
			is2n: false,
			hash: 539_276_951,
			exp:  39_276_951,
		},
		{
			size: 9_768_574_635_223_546,
			is2n: false,
			hash: 924_768_574_635_223_546,
			exp:  6_522_558_924_210_222,
		},
		{
			size: 9_768_574_635_223_546,
			is2n: false,
			hash: 8_768_574_635_223_546,
			exp:  8_768_574_635_223_546,
		},
		{
			size: 8_796_093_022_208,
			is2n: true,
			hash: 8_768_574_635_223_546,
			exp:  7_665_985_104_378,
		},
	}

	for idx, test := range tests {
		test := test

		t.Run(fmt.Sprintf("TestIndexFactory_%d", idx), func(t *testing.T) {
			t.Parallel()

			callback := fixedmap.IndexFactory(test.size)
			act := callback(test.hash)

			if act != test.exp {
				t.Errorf("Unexpected index.\n\vExpected: %d\n\vActual: %d\n\v", test.exp, act)
			}
		})
	}
}
