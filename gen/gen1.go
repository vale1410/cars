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

type countType int

const (
	optionType countType = iota
	classType
)

type Countable struct {
	index    int
	typ      countType
	window   int
	capacity int
	demand   int
    lower   []int
    upper   []int
}

type CountVar struct {
	item  Countable
	pos   int
	count int
}

type PosVar struct {
	item Countable
	pos  int
}

type IdGen struct {
    id int
    countVarMap map[CountVar]int
    posVarMap map[PosVar]int
}

func NewIdGen(size int,class_count int, option_count int) (gen IdGen) {
	gen.id = 1
    gen.countVarMap = make([CountVar]int, size*class_count^2)
	gen.posVarMap = make([PosVar]int, size*(class_count+option_count))
    return
}

func (gen *IdGen) getIdCountVar(v *CountVar) (id int) {
    id, b := gen.countVarMap[v]
    
    if !b {
        gen.id++
        id = gen.id
        gen.countVarMap[v] = id
    } 
    return id
}


type clause []int

func parse(filename string) bool {
	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("problem in newInput:\n ")
		return false
	}

	b := bytes.NewBuffer(input)

	lines := strings.Split(b.String(), "\n")

	state := 0

	var size int
	var option_count, class_count int
	var options []Countable
	var classes []Countable
	var class2option [][]bool

	// parse stuff
	for _, l := range lines {
		numbers := strings.Split(l, " ")
		if digitRegexp.MatchString(numbers[0]) {
			switch state {
			case 0:
				{
					size, _ = strconv.Atoi(numbers[0])
					option_count, _ = strconv.Atoi(numbers[1])
					options = make([]Countable, option_count)
					class_count, _ = strconv.Atoi(numbers[2])
					classes = make([]Countable, class_count)
					class2option = make([][]bool, class_count)
				}
			case 1:
				{
					for i, v := range numbers {
						capacity, _ := strconv.Atoi(v)
						options[i].typ = optionType
						options[i].index = i
						options[i].capacity = capacity
					}
				}
			case 2:
				{
					for i, v := range numbers {
						window, _ := strconv.Atoi(v)
						options[i].window = window
					}
				}
			default:
				{
					num, _ := strconv.Atoi(numbers[0])
					classes[num].index = num
					classes[num].typ = classType
					class2option[num] = make([]bool, option_count)
					slope := 1.0
					for i, v := range numbers {
						if i == 1 {
							demand, _ := strconv.Atoi(v)
							classes[num].demand = demand
						} else if i > 1 {
							value, _ := strconv.Atoi(v)
							has_option := value == 1
							class2option[num][i-2] = has_option
							if has_option {
								options[i-2].demand += classes[num].demand
								slope2 := float64(options[i-2].capacity) / float64(options[i-2].window)
								if slope2 < slope {
									slope = slope2
									classes[num].capacity = options[i-2].capacity
									classes[num].window = options[i-2].window
								}
							}
						}
					}
				}
			}
			state++
		} else {
			fmt.Println(l)
		}
	}
	
    // compute lower and upper bounds for each Countable

    for i := range options {
        options[i].createBounds(size)
    } 

    for i := range classes {
        classes[i].createBounds(size)
    } 


	fmt.Println("options: ", options)
	fmt.Println("classes: ", classes)
	fmt.Println("class2option: ", class2option)

    //
    //clauses := make([]clause,0)
    //

	return true
}

func createAtMostSeq(c Countable, size int) (clauses []clause) {
    

}


func (c *Countable) createBounds(size int) {
	c.lower = make([]int, size)
	c.upper = make([]int, size)

	h := c.demand

	for i := size - 1; i >= 0; i-- {
		q := c.window
		u := c.capacity

		for i >= 0 {

			c.lower[i] = h

			if u > 0 {
				u--
				if h > 0 {
					h--
				}
			}
			q--
			if q <= 0 {
				break
			}
			i--
		}
	}

	h = 1

	for i := 0; i < size; i++ {
		q := c.window
		u := c.capacity

		for i < size {

			c.upper[i] = h + 1

			if q <= u && h < c.demand {
				h++
			}
			q--
			if q <= 0 {
				break
			}
			i++
		}

	}
}
