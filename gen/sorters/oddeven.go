package sorters

import (
//	"log"
	"math"
)

var newId int
var pos int

func compareAndSwap(array []int, comparators []Comparator, i int, j int) {
	newId += 2
	comparators[pos] = Comparator{array[i], array[j], newId - 2, newId - 1}
	pos++
	array[i] = newId - 2
	array[j] = newId - 1
	//	log.Println(array)
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

// createOddEvenEncoding creates sorting network. 
// size has to be power of 2. 
// See sorter.Sorter for the data structure. 
func createOddEvenEncoding(n int) Sorter {

	// conservatively estimating the size of the comparator to allocate memory
	nn := float64(n)
	size := int(nn * math.Log(nn) * math.Log(nn))

	comparators := make([]Comparator, size)
	output := make([]int, n)
	newId = n+2
	pos = 0

	for i, _ := range output {
		output[i] = i+2
	}

	oddevenMergeRange(output, comparators, 0, n-1)

	var last int

	for i, comp := range comparators {
		if comp.A == 0 && comp.B == 0 {
			last = i
			break
		}
	}

	comparators = comparators[:last]
	//log.Println("Created Sorter of size: ", n, "with comparators:", last)

	input := make([]int, n)

	for i, _ := range input {
		input[i] = i+2
	}

	return Sorter{comparators, input, output}
}
