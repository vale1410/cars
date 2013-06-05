package oddeven

import (
	"fmt"
	"math"
)

var newId int
var pos int

// c = a && b
// d = a || b
type Comparator struct {
	a, b, c, d int
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

// it generates the id set for one comparator run of size n
func CreateOddEvenEncoding(s int) []Comparator {
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

	//fmt.Println(array)

	oddevenMergeRange(array, comparators, 0, len(array)-1)

    var last int
	for i, comp := range comparators {
		if comp.a == 0 && comp.b == 0 {
            last = i
            break 
        } 
    } 
    comparators = comparators[:last]
    fmt.Println("Number of comparators:",last)

	//fmt.Println(comparators)
    //printGraph(comparators)

	//do the magic

	// shrink again
	return comparators
}

func printGraph(comparators []Comparator) {
	fmt.Println("digraph {")
	for _, comp := range comparators {
		if comp.a != 0 || comp.b != 0 {
			fmt.Printf("n%v -> t%v \n", comp.a, comp.a)
			fmt.Printf("n%v -> t%v \n", comp.b, comp.b)
			fmt.Printf("t%v -> t%v \n", comp.a, comp.b)
			fmt.Printf("t%v -> n%v \n", comp.a, comp.c)
			fmt.Printf("t%v -> n%v \n", comp.b, comp.d)
		}
	}
	fmt.Println("}")
}
