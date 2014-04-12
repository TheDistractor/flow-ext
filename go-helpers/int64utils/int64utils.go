//Package to supply some []int64 functionality, such as Sort, Inserts.
//TODO: Add Missing functional such as Deletes etc.
package int64utils

import "sort"


type Int64Array []int64


//Sort Interface
func (s Int64Array) Len() int {
return len(s)
}

func (s Int64Array) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Int64Array) Less(i, j int) bool {
	return s[i] < s[j]
}


//Int64Array Type helpers
func (s Int64Array) Search(v int64) (int) {
	return Int64ArraySearch(s,v)
}

func (s Int64Array) Insert(v int64) ([]int64) {
	return Int64OrderedArrayInsert(s,v)
}

func (s Int64Array) InsertAt(v int64, i int) ([]int64) {
	return Int64ArrayInsertAt(s,v,i)
}

//Int64Array Direct Primitives
func Int64ArraySearch(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func Int64OrderedArrayInsert(s []int64, v int64) ([]int64) {
	return Int64ArrayInsertAt(s, v, Int64ArraySearch(s,v))
}

func Int64ArrayInsertAt(s []int64, v int64, i int) ([]int64) {
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
