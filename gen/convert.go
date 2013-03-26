package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var digitRegexp = regexp.MustCompile("([0-9]+ )*[0-9]+")

var name = flag.String("file", "test.txt", "Name of the txt file")

func main() {
	flag.Parse()
	filename := *name

	parse(filename)
}

func parse(filename string) bool {
	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("problem in newInput:\n ")
		return false
	}

	b := bytes.NewBuffer(input)

	lines := strings.Split(b.String(), "\n")

	sw := 0
	//size := 0
	//option_count := 0
	//class_count := 0
	for _, l := range lines {
		numbers := strings.Split(l, " ")
		if digitRegexp.MatchString(numbers[0]) {
			switch sw {
			case 0:
				{
					printAtom("def", numbers)
					//size,_ = strconv.Atoi(numbers[0])
					//option_count,_ = strconv.Atoi(numbers[1])
					//class_count,_ = strconv.Atoi(numbers[2])
				}
			case 1:
				{
					p := make([]string, 2, 2)
					for i, v := range numbers {
						p[0] = strconv.Itoa(i + 1)
						p[1] = v
						printAtom("option_max", p)
					}
				}
			case 2:
				{
					p := make([]string, 2, 2)
					for i, v := range numbers {
						p[0] = strconv.Itoa(i + 1)
						p[1] = v
						printAtom("option_window", p)
					}
				}
			default:
				{

					printAtom("class_count", numbers[:2])
					for i, v := range numbers {
						if i > 1 {
							value, _ := strconv.Atoi(v)
							if value == 1 {
								p := make([]string, 2, 2)
								p[0] = numbers[0]
								p[1] = strconv.Itoa(i - 1)
								printAtom("class_option", p)
							}
						}
					}
				}
			}
			sw++
		} else {
			fmt.Println(l)
		}
	}
	return true
}

func printAtom(name string, values []string) {
	// for each element in the array, print stuff
	fmt.Printf("%v(", name)
	for n, s := range values {
		if n < len(values)-1 {
			fmt.Printf("%v,", s)
		} else {
			fmt.Printf("%v).", s)
		}
	}
	fmt.Println()
}
