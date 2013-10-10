package sorters

import (
	"log"
	"math/rand"
	"sort"
	"testing"
)

// TestCardinality check constraint sum n <= k
func TestCardinality(t *testing.T) {

	size := 8
	k := 3

	array1 := make([]int, size)
	array2 := make([]int, size)

	copy(array2, array1)
	sort.Ints(array2)

	mapping := make(map[int]int, size)

	sorter := CreateSortingNetwork(size, -1, OddEven)

	for i := size - k; i < size; i++ {
		mapping[sorter.Out[i]] = 0
		sorter.Out[i] = 0
	}

	sorter.PropagateBackwards(mapping)

	log.Println(sorter)

}

func TestSorting(t *testing.T) {

	sizes := []int{3, 4, 6, 8, 123, 234, 256, 1024, 1025}

	for _, size := range sizes {
		normalSorting(size, t)
	}

}

func TestCut(t *testing.T) {

	sizes := []int{3, 4, 6, 8, 123, 234, 256, 1024, 1025}
	cuts := []int{2, 2, 3, 2, 68, 123, 250, 543, 800}

	for i, size := range sizes {
		cutSorting(size, cuts[i], t)
	}
}

func cutSorting(size int, cut int, t *testing.T) {
	array1 := make([]int, size)
	array2 := make([]int, size)

	element := 0

	for i, _ := range array1 {
		if i == cut {
			element = 0
		}
		array1[i] = element
		element++
	}

	copy(array2, array1)
	sort.Ints(array2)
	sorter := CreateSortingNetwork(len(array1), cut, OddEven)
	compareArrays(sorter, array1, array2, t)
}

func normalSorting(size int, t *testing.T) {

	array1 := rand.Perm(size)
	array2 := make([]int, size)
	copy(array2, array1)
	sort.Ints(array2)
	sorter := CreateSortingNetwork(len(array1), -1, OddEven)
	compareArrays(sorter, array1, array2, t)
}

func compareArrays(sorter Sorter, array1, array2 []int, t *testing.T) {

	mapping := make(map[int]int, len(sorter.Comparators))

	for i, x := range sorter.In {
		mapping[x] = array1[i]
	}

	for _, comp := range sorter.Comparators {

		b, bok := mapping[comp.B]
		a, aok := mapping[comp.A]

		if !aok {
			t.Error("not in mapping", comp.A)
		}

		if !bok {
			t.Error("not in mapping", comp.B)
		}

		mapping[comp.D] = max(a, b)
		mapping[comp.C] = min(a, b)

	}

	output := make([]int, len(array1))

	for i, x := range sorter.Out {
		output[i] = mapping[x]
		if output[i] != array2[i] {
			t.Error(output[i], array2[i], "Output array does not coincide")
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func TestGenerateSAT(t *testing.T) {
}
