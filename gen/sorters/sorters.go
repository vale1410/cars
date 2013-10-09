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

func CreateSortingNetwork(s int, cut int, typ SortingNetworkType) (sorter Sorter) {

	//grow to be 2^n
	n := 1
	for n < s {
		n *= 2
	}

	log.Println("Input: ", s, "Power of 2: ", n)

	switch typ {
	case OddEven:
		sorter = createOddEvenEncoding(n)
	default:
		log.Panic("Type of sorting network not implemented")
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

		mapping := make(map[int]int, 0)
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

			if a < s && b < s && (a >= cut || b < cut) {
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

		log.Println("Propagate Ordering with cut", cut, ", removed ", len(comparators)-l, "sorters")

		sorter.Comparators = out

		return
	}
}

// PropagteZeros propagates from left to right a given set of zeros in mapping
// mapping contains the set of zeros in the input vector of sorter (-1 stands for zero)
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
	sorter.Comparators = out

	log.Println("Propagate Zeros; removed ", len(comparators)-l, "sorters")

	return
}

// ChangeSize shrinkes the sorter to a size s
func (sorter *Sorter) changeSize(s int) {

	n := len(sorter.In)

	mapping := make(map[int]int, n-s)

	for i := s; i < n; i++ {
		mapping[i] = -1
	}

	sorter.propagateZeros(mapping)

	//potential check for s..n being -1

	for i, x := range sorter.Out {
		if r, ok := mapping[x]; ok {
			sorter.Out[i] = r
		}
	}

	sorter.In = sorter.In[:s]
	sorter.Out = sorter.Out[:s]

	return
}

func (sorter *Sorter) propagateForward(zeros map[int]bool, ones map[int]bool) {

	mapping := make(map[int]int, 0) next step
	l := 0
	comparators := sorter.Comparators
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

		if zeros[comp.A]  {
			zeros[comp.D] = true
            // C is like B
            mapping[comp.C] = comp.B
			removed = true
		}

		if zeros[comp.B]  {
			zeros[comp.D] = true
            //C is like A
            mapping[comp.C] = comp.A
			removed = true
		}

		if ones[comp.A]  {
			ones[comp.C] = true
            // D is like B
            mapping[comp.D] = comp.B
			removed = true
		}

		if ones[comp.B]  {
			ones[comp.C] = true
            //D is like A
            mapping[comp.D] = comp.A
			removed = true
		}

		if zeros[comp.A] && zeros[comp.B] {
			zeros[comp.C] = true
			removed = true
		}

		if ones[comp.A] && ones[comp.B] {
			ones[comp.D] = true
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

func propagateBackwards(sorter Sorter, zeros map[int]bool, ones map[int]bool) {

	l := 0
	comparators := sorter.Comparators
	remove := Comparator{0, 0, 0, 0}

	for i := len(comparators) - 1; i >= 0; i-- {

		comp := comparators[i]

		removed := false

		if zeros[comp.C] {
			zeros[comp.A] = true
			zeros[comp.B] = true
			if ones[comp.C] {
				log.Panic("Sorting network has problems in propagateBackwards", comp)
			}
			zeros[comp.D] = true
			removed = true
		}

		if ones[comp.D] {
			ones[comp.A] = true
			ones[comp.B] = true
			if zeros[comp.C] {
				log.Panic("Sorting network has problems in propagateBackwards", comp)
			}
			ones[comp.C] = true
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
