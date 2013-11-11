package sorters

import (
	"log"
)

const (
	OddEven SortingNetworkType = iota
	Pairwise
	Bitonic
	ShellSort
	Insertion
	Bubble
)

const (
	AtMost CardinalityType = iota
	AtLeast
	Equal
)

type SortingNetworkType int
type CardinalityType int

// The slice of comparators must be in correct order,
// meaning that the comparator with input A and B must
// occur after the comparator with output A and B.
// expection are 0 and 1, that define true and false. They may only occur as
// output of a comparator (C or D), but not as input (as will be propagated
// and the comparator removed).
type Sorter struct {
	Comparators []Comparator
	In          []int
	Out         []int
}

// Ids for the connections (A,B,C,D) start at 2 and are incremented.
// Id 0 and 1 are reserved for true and false respectively
// B --|-- D = A && B
//     |
// A --|-- C = A || B
type Comparator struct {
	A, B, C, D int
}

func CreateCardinalityNetwork(size int, k int, cType CardinalityType, sType SortingNetworkType) (sorter Sorter) {

	mapping := make(map[int]int, size)

	sorter = CreateSortingNetwork(size, -1, sType)

	switch cType {
	case AtMost:
		for i := size - k; i < size; i++ {
			mapping[sorter.Out[i]] = 0
			sorter.Out[i] = 0
		}
		sorter.PropagateBackwards(mapping)
		sorter.Out = sorter.Out[:k]
	case AtLeast:
		for i := 0; i < k; i++ {
			mapping[sorter.Out[i]] = 1
			sorter.Out[i] = 1
		}
		sorter.PropagateBackwards(mapping)
		sorter.Out = sorter.Out[k:]
	case Equal:
		for i := size - k; i < size; i++ {
			mapping[sorter.Out[i]] = 0
			sorter.Out[i] = 0
		}
		for i := 0; i < k; i++ {
			mapping[sorter.Out[i]] = 1
			sorter.Out[i] = 1
		}
		sorter.PropagateBackwards(mapping)
		sorter.Out = nil
	default:
		log.Panic("CardnalityNot implemented yet")
	}
	return
}

// CreateSorting Networks creates a sorting network of arbitrary size and cut
// and of type.
func CreateSortingNetwork(s int, cut int, typ SortingNetworkType) (sorter Sorter) {

	//grow to be 2^n
	n := 1
	for n < s {
		n *= 2
	}

	comparators := make([]Comparator, 0)
	output := make([]int, n)

	offset := 2

	for i, _ := range output {
		output[i] = i + offset
	}
	input := make([]int, n)
	copy(input, output)

	newId := n + offset

	switch typ {
	case OddEven:
		oddevenSort(&newId, output, &comparators, 0, n-1)
	case Bitonic:
		triangleBitonic(&newId, output, &comparators, 0, n-1)
	default:
		log.Panic("Type of sorting network not implemented yet")
	}

	sorter = Sorter{comparators, input, output}
	sorter.changeSize(s)
	sorter.PropagateOrdering(cut)

	return
}

// from 0..cut-1 sorted and from cut .. length-1 sorted
// propagated and remove comparators
func (sorter *Sorter) PropagateOrdering(cut int) {

	if cut <= 0 || cut >= len(sorter.In) {
		return
	} else {

		mapping := make(map[int]int, len(sorter.Comparators))
		l := 0
		s := len(sorter.In)
		comparators := sorter.Comparators

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

			if a < s && b < s && (a >= sorter.In[cut] || b < sorter.In[cut]) {
				// we have an already sorted input
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

		//log.Println("Propagate Ordering with cut", cut, ", after size ", len(comparators)-l, "sorters")

		sorter.Comparators = out

		return
	}
}

// ChangeSize shrinks the sorter to a size s
func (sorter *Sorter) changeSize(s int) {

	n := len(sorter.In)

	mapping := make(map[int]int, n-s)

	for i := s; i < n; i++ {
		//setting the top n-s elements to zero
		mapping[sorter.In[i]] = 0
	}

	//sorter.PropagateZeros(mapping)
	sorter.PropagateForward(mapping)

	//potential check for s..n being 0

	for i, x := range sorter.Out {
		if r, ok := mapping[x]; ok {
			sorter.Out[i] = r
		}
	}

	sorter.In = sorter.In[:s]
	sorter.Out = sorter.Out[:s]

	return
}

func (sorter *Sorter) PropagateForward(mapping map[int]int) {

	l := 0
	comparators := sorter.Comparators
	// remove is a comparator with no functionality
	remove := Comparator{0, 0, 0, 0}

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

		removed := false

		if a == 0 {
			mapping[comp.D] = 0
			mapping[comp.C] = b
			removed = true
		}

		if b == 0 {
			mapping[comp.D] = 0
			mapping[comp.C] = a
			removed = true
		}

		if a == 1 {
			mapping[comp.C] = 1
			mapping[comp.D] = b
			removed = true
		}

		if b == 1 {
			mapping[comp.C] = 1
			mapping[comp.D] = a
			removed = true
		}

		if a == 0 && b == 0 {
			mapping[comp.C] = 0
			removed = true
		}

		if a == 1 && b == 1 {
			mapping[comp.D] = 1
			removed = true
		}

		if removed {
			l++
			comparators[i] = remove
			removed = false
		}
	}

	//remove the unused comparators
	out := make([]Comparator, 0, l)
	for _, comp := range comparators {
		if comp != remove {
			out = append(out, comp)
		}
	}

	sorter.Comparators = out
}

func (sorter *Sorter) PropagateBackwards(mapping map[int]int) {

	l := 0
	comparators := sorter.Comparators
	remove := Comparator{0, 0, 0, 0}

	cleanMapping := make(map[int]int, 0)

	for i := len(comparators) - 1; i >= 0; i-- {

		comp := comparators[i]

		removed := false

		if value, ok := mapping[comp.C]; ok && value == 0 {
			mapping[comp.A] = 0
			mapping[comp.B] = 0
			cleanMapping[comp.D] = 0
			removed = true
		}

		if value, ok := mapping[comp.D]; ok && value == 1 {
			mapping[comp.A] = 1
			mapping[comp.B] = 1
			cleanMapping[comp.C] = 1
			removed = true
		}

		if removed {

			l++
			comparators[i] = remove
		}
	}

	//remove the unused comparators
	out := make([]Comparator, 0, l)
	for _, comp := range comparators {
		if comp != remove {
			if value, ok := mapping[comp.C]; ok {
				comp.C = value
			}

			if value, ok := mapping[comp.D]; ok {
				comp.D = value
			}
			out = append(out, comp)
		}
	}

	sorter.Comparators = out

	if len(cleanMapping) > 0 {
		sorter.PropagateForward(cleanMapping)
	}

}

// Functions for creating sorters
func compareAndSwap(newId *int, array []int, comparators *[]Comparator, i int, j int) {
	*newId += 2
	*comparators = append(*comparators, Comparator{array[i], array[j], *newId - 2, *newId - 1})

	array[i] = *newId - 2
	array[j] = *newId - 1
}
