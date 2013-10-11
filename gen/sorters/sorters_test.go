package sorters

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
)

// TestCardinality check constraint sum n <= k
func TestCardinality(t *testing.T) {

	sizes := []int{3, 4, 6, 9, 9, 9, 68, 123, 250, 543, 1024, 1025}
	ks := []int{2, 2, 3, 2, 6, 7, 8, 8, 100, 200, 256, 800}
	//sizes := []int{123}
	//ks := []int{100}

	for i, size := range sizes {
		cardinalitySorting(size, ks[i], t)
	}
}

func cardinalitySorting(size int, k int, t *testing.T) {

	array1 := rand.Perm(size)
	array2 := make([]int, size)

	copy(array2, array1)
	sort.Ints(array2)

	mapping := make(map[int]int, size)

	sorter := CreateSortingNetwork(size, -1, OddEven)

	printSorter(sorter, "test1.dot")

	for i := size - k; i < size; i++ {
		mapping[sorter.Out[i]] = 0
		sorter.Out[i] = 0
		array2[i] = 0
	}

	sorter.PropagateBackwards(mapping)
	printSorter(sorter, "test2.dot")

	sortAndCompareArrays(sorter, array1, array2, t)
}

func TestSorting(t *testing.T) {

	sizes := []int{3, 4, 6, 8, 9, 31, 32, 33, 63, 65, 123, 234, 256, 1024, 1025}

	for _, size := range sizes {
		normalSorting(size, t)
		normalSorting(size, t)
	}

}

func TestCut(t *testing.T) {

	sizes := []int{3, 4, 6, 9, 68, 123, 250, 543, 1024, 1025}
	cuts := []int{2, 2, 3, 7, 2, 100, 100, 234, 256, 800}

	for i, size := range sizes {
		cutSorting(size, cuts[i], t)
	}
}

func cutSorting(size int, cut int, t *testing.T) {

	array1 := rand.Perm(size)
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
	sortAndCompareArrays(sorter, array1, array2, t)
}

func normalSorting(size int, t *testing.T) {

	array1 := rand.Perm(size)
	array2 := make([]int, size)
	copy(array2, array1)
	sort.Ints(array2)
	sorter := CreateSortingNetwork(len(array1), -1, OddEven)
	sortAndCompareArrays(sorter, array1, array2, t)
}

func sortAndCompareArrays(sorter Sorter, array1, array2 []int, t *testing.T) {

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

		if comp.D > 1 { // 0,1, specific meaning
			mapping[comp.D] = max(a, b)
		}
		if comp.C > 1 { // 0,1, specific meaning
			mapping[comp.C] = min(a, b)
		}

	}

	output := make([]int, len(array1))

	e := false

	for i, x := range sorter.Out {
		output[i] = mapping[x]
		if output[i] != array2[i] {
			t.Error("Output array does not coincide in position", i)
			e = true
		}
	}
	if e {
		t.Error("ideal", array2)
		t.Error("output", output)
		//t.Error(mapping)
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

func printSorter(sorter Sorter, filename string) {
	file, ok := os.Create(filename)
	if ok != nil {
		panic("Can open file to write.")
	}
	file.Write([]byte(fmt.Sprintln("digraph {")))
	file.Write([]byte(fmt.Sprintln("  graph [rankdir = LR, splines=ortho];")))

	rank := "{rank=same; "
	for i := 0; i < len(sorter.Out); i++ {
		if sorter.Out[i] > 1 {
			rank += fmt.Sprintf(" t%v ", sorter.Out[i])
		}
	}
	rank += "}; "

	for i := 0; i < len(sorter.Out); i++ {
		file.Write([]byte(fmt.Sprintf("n%v -> t%v\n", sorter.In[i], sorter.In[i])))
	}

	file.Write([]byte(rank))
	rank = "{rank=same; "
	for i := 0; i < len(sorter.Out); i++ {
		rank += fmt.Sprintf(" t%v ", sorter.In[i])
	}
	rank += "}; "
	file.Write([]byte(rank))

	//var rank string
	for _, comp := range sorter.Comparators {
		rank = "{rank=same; "
		rank += fmt.Sprintf(" t%v t%v ", comp.A, comp.B)
		rank += "}; "
		file.Write([]byte(rank))
	}

	for _, comp := range sorter.Comparators {
		if comp.A > 1 && comp.B > 1 {
			//file.Write([]byte(fmt.Sprintf("t%v -> t%v [dir=none]\n", comp.A, comp.B)))
			file.Write([]byte(fmt.Sprintf("t%v -> t%v \n", comp.B, comp.A)))
		}
		if comp.C > 1 {
			file.Write([]byte(fmt.Sprintf("t%v -> t%v \n", comp.A, comp.C)))
		}
		if comp.D > 1 {

			file.Write([]byte(fmt.Sprintf("t%v -> t%v \n", comp.B, comp.D)))
		}
	}
	file.Write([]byte(fmt.Sprintln("}")))
}
