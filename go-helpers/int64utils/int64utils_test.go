package int64utils

import (
	"testing"
	"sort"
)

func BenchmarkInt64ArrayOrderedInsert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var i64s Int64Array
		i64s = []int64{1,3,5,7,9,11,13,15,17}

		i64s = i64s.Insert(0)
		i64s = i64s.Insert(2)
		i64s = i64s.Insert(12)

	}
}

func BenchmarkInt64ArrayUnOrderedSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i64s := Int64Array{1,3,5,7,9,11,0,2,12}
		sort.Sort(i64s)


	}
}
