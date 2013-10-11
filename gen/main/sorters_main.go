package main

import (
	"../sorters"
	"flag"
	"fmt"
	"os"
)

var debug = flag.Bool("debug", false, "Adds debug information.")
var size = flag.Int("s", 8, "Size of the array to sort.")
var cut = flag.Int("cut", -1, "Cut marks the position which devides the array in two sorted [:cut],[cut:]. -1 defines no cut.")
var dot = flag.String("dot", "", "Create dot compatible output of graph")

func main() {

	flag.Parse()

	s := sorters.CreateSortingNetwork(*size, *cut, sorters.OddEven)

	//fmt.Println()
	if *dot != "" {
		printSorter(s, *dot)
	}

	if *debug {
		fmt.Println(s)
	}
}

func printSorter(sorter sorters.Sorter, filename string) {
	file, ok := os.Create(filename)
	if ok != nil {
		panic("Can open file to write.")
	}
	file.Write([]byte(fmt.Sprintln("digraph {")))

	rank := "{rank=same; "
	for i := 0; i < len(sorter.Out); i++ {
		rank += fmt.Sprintf(" t%v ", sorter.Out[i])
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
		file.Write([]byte(fmt.Sprintf("t%v -> t%v [dir=none]\n", comp.A, comp.B)))
		file.Write([]byte(fmt.Sprintf("t%v -> t%v \n", comp.A, comp.C)))
		file.Write([]byte(fmt.Sprintf("t%v -> t%v \n", comp.B, comp.D)))
	}
	file.Write([]byte(fmt.Sprintln("}")))
}
