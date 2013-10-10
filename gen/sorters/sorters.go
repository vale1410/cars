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

type SortingNetworkType int

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

func CreateSortingNetwork(s int, cut int, typ SortingNetworkType) (sorter Sorter) {

	//grow to be 2^n
	n := 1
	for n < s {
		n *= 2
	}

	//log.Println("Input: ", s, "Power of 2: ", n)

	switch typ {
	case OddEven:
		sorter = createOddEvenEncoding(n)
	default:
		log.Panic("Type of sorting network not implemented yet")
	}

	sorter.changeSize(s)
	sorter.propagateOrdering(cut)

	return
}

// from 0..cut-1 sorted and from cut .. length-1 sorted
// propagated and remove comparators
func (sorter *Sorter) propagateOrdering(cut int) {

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

// PropagteZeros propagates from left to right a given set of zeros in mapping
// mapping contains the set of zeros in the input vector of sorter (0 stands for zero)
func (sorter *Sorter) propagateZeros(mapping map[int]int) {

	l := 0
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

		if aok && a == 0 {
			comparators[i] = zero
			mapping[comp.D] = 0
			mapping[comp.C] = b
		}

		if bok && b == 0 {
			comparators[i] = zero
			mapping[comp.D] = 0
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
	sorter.Comparators = out

	//log.Println("Propagate Zeros; new size ", l, "sorters")

	return
}

// ChangeSize shrinks the sorter to a size s
func (sorter *Sorter) changeSize(s int) {

	n := len(sorter.In)

	mapping := make(map[int]int, n-s)

	for i := s; i < n; i++ {
		//setting the top n-s elements to zero
		mapping[sorter.In[i]] = 0
	}

	sorter.propagateZeros(mapping)

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

		if a == 0 && b == 0 {
			mapping[comp.D] = 1
			removed = true
		}

		if removed {
			l++
			comparators[i] = remove
			removed = false
		}
	}

	//remove zeros and then return comparators
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

	for i := len(comparators) - 1; i >= 0; i-- {

		comp := comparators[i]

		removed := false

		if mapping[comp.C] == 0 {
			mapping[comp.A] = 0
			mapping[comp.B] = 0
			if mapping[comp.D] == 1 {
				log.Panic("Sorting network has problems in propagateBackwards", comp)
			}
			mapping[comp.D] = 0
			removed = true
		}

		if mapping[comp.D] == 1 {
			mapping[comp.A] = 1
			mapping[comp.B] = 1
			if mapping[comp.C] == 0 {
				log.Panic("Sorting network has problems in propagateBackwards", comp)
			}
			mapping[comp.C] = 1
			removed = true
		}

		if removed {
			l++
			comparators[i] = remove
		}

	}

	//remove zeros and then return comparators
	out := make([]Comparator, 0, l)
	for _, comp := range comparators {
		if comp != remove {
			out = append(out, comp)
		}
	}
	sorter.Comparators = out
}
