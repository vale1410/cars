package sorters

import (
//	"fmt"
	"math"
)

var bitonicNewId int
var bitonicPos int

func bitonicCompareAndSwap(up bool, array []int, comparators []Comparator, i int, j int) {
	//fmt.Println("compare", up, i, j)
	bitonicNewId += 2
	var C, D int

	if !up {
		C = bitonicNewId - 2
		D = bitonicNewId - 1
	} else {
		C = bitonicNewId - 1
		D = bitonicNewId - 2
	}

	comparators[bitonicPos] = Comparator{array[i], array[j], C, D}

	if !up {
		array[i] = D
		array[j] = C
	} else {
		array[i] = C
		array[j] = D
	}

	bitonicPos++
}

func bitonicMerge(up bool, array []int, comparators []Comparator, lo int, hi int) {
	if (hi - lo) >= 1 {
	    //fmt.Println("merge", up, lo, hi)
		mid := lo + ((hi - lo) / 2)
		for i := 0; i <= mid-lo; i++ {
			bitonicCompareAndSwap(up, array, comparators, lo+i, mid+i+1)
		}
		bitonicMerge(up, array, comparators, lo, mid)
		bitonicMerge(up, array, comparators, mid+1, hi)
	}
}

func bitonicSort(up bool, array []int, comparators []Comparator, lo int, hi int) {
	if (hi - lo) >= 1 {
	    //fmt.Println("sort", up, lo, hi)
		mid := lo + ((hi - lo) / 2)
		bitonicSort(true, array, comparators, lo, mid)
		bitonicSort(false, array, comparators, mid+1, hi)
		bitonicMerge(up, array, comparators, lo, hi)
	}
}

// createOddEvenEncoding creates sorting network.
// size has to be power of 2.
// See sorter.Sorter for the data structure.
func createBitonicEncoding(n int) Sorter {

	// conservatively estimating the size of the comparator to allocate memory
	nn := float64(n)
	size := int(nn * math.Log(nn) * math.Log(nn))

	comparators := make([]Comparator, size)
	output := make([]int, n)
	bitonicNewId = n + 2
	bitonicPos = 0

	offset := 2

	for i, _ := range output {
		output[i] = i + offset
	}

	bitonicSort(true, output, comparators, 0, n-1)

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
		input[i] = i + offset
	}

	return Sorter{comparators, input, output}
}

//func triangle(n int)
//    triangle
//
//
//
//func createBitonicEncoding(n int) Sorter {
//
//    triangle(n)
//
//
//    //dummy stuff
//	comparators := make([]Comparator, n)
//	output := make([]int, n)
//	input := make([]int, n)
//	return Sorter{comparators, input, output}
//}
