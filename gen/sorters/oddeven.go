package sorters

import (
	"fmt"
	"math"
)

var newId int
var pos int

type Sorter struct {
	Comparators []Comparator
	In          []int
	Out         []int
}

// B --|-- D = A && B
//     |
// A --|-- C = A || B
type Comparator struct {
	A, B, C, D int
}

func compareAndSwap(array []int, comparators []Comparator, i int, j int) {
	newId += 2
	comparators[pos] = Comparator{array[i], array[j], newId - 2, newId - 1}
	pos++
	array[i] = newId - 2
	array[j] = newId - 1
	//	fmt.Println(array)
}

func oddevenMerge(array []int, comparators []Comparator, lo int, hi int, r int) {
	step := r * 2
	if step < hi-lo {
		oddevenMerge(array, comparators, lo, hi, step)
		oddevenMerge(array, comparators, lo+r, hi-r, step)
		for i := lo + r; i <= hi-r; i += step {
			compareAndSwap(array, comparators, i, i+r)
		}
	} else {
		compareAndSwap(array, comparators, lo, lo+r)
	}
}

func oddevenMergeRange(array []int, comparators []Comparator, lo int, hi int) {
	if (hi - lo) >= 1 {
		mid := lo + ((hi - lo) / 2)
		oddevenMergeRange(array, comparators, lo, mid)
		oddevenMergeRange(array, comparators, mid+1, hi)
		oddevenMerge(array, comparators, lo, hi, 1)
	}
}

// it generates the id set for one comparator run of size s
func CreateOddEvenEncoding(s int, cut int) Sorter {
	//grow to be 2^n
	n := 1
	for n < s {
		n *= 2
	}

	nn := float64(n)
	size := int(nn * math.Log(nn) * math.Log(nn))
	fmt.Println("Input: ", s, "Power of 2: ", n, "Estimated number of comparators: ", size)

	comparators := make([]Comparator, size)
	array := make([]int, n)
	newId = n
	pos = 0

	for i, _ := range array {
		array[i] = i
	}

	oddevenMergeRange(array, comparators, 0, len(array)-1)

	var last int

	for i, comp := range comparators {
		if comp.A == 0 && comp.B == 0 {
			last = i
			break
		}
	}

	comparators = comparators[:last]
	fmt.Println("Number of comparators before shrinking:", last)

	// shrink the comparator to size s by setting the last n-s to 0
	// and propagate through
	mapping := make(map[int]int, n-s)

	for i := s; i < n; i++ {
		mapping[i] = -1
	}

	//fmt.Println(comparators)
	comparators = propagateZeros(comparators, mapping)
	fmt.Println("Number of comparators after propagating zeros:", len(comparators))

	for i, x := range array {
		if r, ok := mapping[x]; ok {
			array[i] = r
		}
	}

	mapping = make(map[int]int, n-s)

	if cut >= 0 {
		comparators = propagateOrdering(comparators, mapping, s, cut)
		fmt.Println("Number of comparators after propagating ordering:", len(comparators))
	}

	for i, x := range array {
		if r, ok := mapping[x]; ok {
			array[i] = r
		}
	}

	input := make([]int, s)

	for i, _ := range input {
		input[i] = i
	}

	return Sorter{comparators, input, array[:s]}
}

// from 0..cut-1 sorted and from cut .. length-1 sorted
// propagated and remove comparators
func propagateOrdering(comparators []Comparator, mapping map[int]int, s int, cut int) []Comparator {

	l := 0
	zero := Comparator{0, 0, 0, 0}

	for i, comp := range comparators {
		a, aok := mapping[comp.A]
		b, bok := mapping[comp.B]

		if aok {
			comparators[i].A = a
		} else {
			a = comp.A
		}
		if bok {
			comparators[i].B = b
		} else {
			b = comp.B
		}

		if a < s && b < s && (a >= cut || b < cut) {
			// we have an already sorted input
			//fmt.Println("Sorted",a,b)
			mapping[comp.C] = a
			mapping[comp.D] = b
			comparators[i] = zero
			l++
		}
	}

	//remove zeros and then return comparators
	out := make([]Comparator, 0, l)

	for _, comp := range comparators {
		if comp != zero {
			out = append(out, comp)
		}
	}

	return out
}

//func propagateForward(sorter Sorter, zeros map[int]bool, ones map[int]bool) {
//
//	mapping := make(map[int]int,len(sorter.comparators))
//
//	for i, comp := range comparators {
//
//}
//
//func propagateBackwards(sorter Sorter, zeros map[int]bool, ones map[int]bool) {
//
//	mapping := make(map[int]int,len(sorter.comparators))
//
//}

func propagateZeros(comparators []Comparator, mapping map[int]int) []Comparator {

	l := 0

	zero := Comparator{0, 0, 0, 0}

	for i, comp := range comparators {
		a, aok := mapping[comp.A]
		b, bok := mapping[comp.B]

		if aok {
			comparators[i].A = a
		} else {
			a = comp.A
		}

		if bok {
			comparators[i].B = b
		} else {
			b = comp.B
		}

		if aok && a == -1 {
			comparators[i] = zero
			mapping[comp.D] = -1
			mapping[comp.C] = b
		}

		if bok && b == -1 {
			comparators[i] = zero
			mapping[comp.D] = -1
			mapping[comp.C] = a
		}

		if comparators[i] == zero {
			l++
		}
	}

	//remove zeros and then return comparators
	out := make([]Comparator, 0, l)

	for _, comp := range comparators {
		if comp != zero {
			out = append(out, comp)
		}
	}
	return out
}
