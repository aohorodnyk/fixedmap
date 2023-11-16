package fixedmap_test

import (
	"math"
	"testing"

	"github.com/aohorodnyk/fixedmap"
)

func BenchmarkIndexFactory2n(b *testing.B) {
	const size = 536_870_912

	tests := []struct {
		hash uint64
		exp  int
	}{
		{
			hash: 0,
			exp:  0,
		},
		{
			hash: 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: size + 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: size*5 + 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: math.MaxUint64,
			exp:  536870911,
		},
	}

	indexCalculator := fixedmap.IndexFactory(size)

	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			index := indexCalculator(test.hash)
			if index != test.exp {
				b.Errorf("Unexpectedly index.\n\tExpected: %d\n\tActual: %d", test.exp, index)
			}
		}
	}
}

func BenchmarkIndexFactory(b *testing.B) {
	const size = 536_870_913

	tests := []struct {
		hash uint64
		exp  int
	}{
		{
			hash: 0,
			exp:  0,
		},
		{
			hash: 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: size + 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: size*5 + 2_363_253,
			exp:  2_363_253,
		},
		{
			hash: math.MaxUint64,
			exp:  63,
		},
	}

	indexCalculator := fixedmap.IndexFactory(size)

	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			index := indexCalculator(test.hash)
			if index != test.exp {
				b.Errorf("Unexpectedly index.\n\tExpected: %d\n\tActual: %d", test.exp, index)
			}
		}
	}
}
