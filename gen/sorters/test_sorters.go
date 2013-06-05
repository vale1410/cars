package main

import (
    "./sorters"
    "flag"
)

var debug = flag.Bool("debug", false, "Adds debug information.")
var size =  flag.Int("s", 8, "Size of the array to sort.")

func main() {
	flag.Parse()
    oddeven.CreateOddEvenEncoding(*size)
}
