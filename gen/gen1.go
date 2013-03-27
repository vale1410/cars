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
var size, class_count, option_count int
var gen IdGen

func main() {
	flag.Parse()
	filename := *name
	parse(filename)
}

const (
	optionType countType = iota
	classType
)

type countType int
type clause []int

type CountableId struct {
	typ   countType
	index int
}

type Countable struct {
	cId      CountableId
	window   int
	capacity int
	demand   int
	lower    []int
	upper    []int
}

type CountVar struct {
	cId   CountableId
	pos   int
	count int
}

type PosVar struct {
	cId CountableId
	pos int
}

type IdGen struct {
	id          int
	countVarMap map[CountVar]int
	posVarMap   map[PosVar]int
}

func NewIdGen() {
	gen.id = 0
	gen.countVarMap = make(map[CountVar]int, size*class_count^2)
	gen.posVarMap = make(map[PosVar]int, size*(class_count+option_count))
	return
}

func getCountId(v CountVar) (id int) {
	id, b := gen.countVarMap[v]

	if !b {
		gen.id++
		id = gen.id
		gen.countVarMap[v] = id
	}
	return id
}

func getPosId(v PosVar) (id int) {
	id, b := gen.posVarMap[v]

	if !b {
		gen.id++
		id = gen.id
		gen.posVarMap[v] = id
	}
	return id
}

func parse(filename string) bool {
	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("problem in newInput:\n ")
		return false
	}

	b := bytes.NewBuffer(input)

	lines := strings.Split(b.String(), "\n")

	state := 0

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
						options[i].cId = CountableId{optionType, i}
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
					classes[num].cId = CountableId{classType, num}
					class2option[num] = make([]bool, option_count)

					// find option with lowest slope
					// to determine capacity and windows

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
					classes[num].createBounds()

					//// backprop the bounds back to the options
					//for p := 0; p < size; p++ {       
					//	for j := 0; j < option_count; j++ {
					//        options[j].upper[p] = options[j].demand+1
					//		if class2option[num][j] {
					//			options[j].lower[p] += classes[num].lower[p]
					//			options[j].upper[p] -= (classes[num].demand + 1 - classes[num].upper[p])
					//		}
					//	}
					//}
				}
			}
			state++
		} else {
			//fmt.Println(l)
		}
	}

	for i := range options {
		options[i].createBounds()
	}

	//fmt.Println("options: ", options)
	//fmt.Println("classes: ", classes)
	//fmt.Println("class2option: ", class2option)

	NewIdGen()

	clauses := make([]clause, 0)

	for _, c := range classes {
		clauses = append(clauses, createAtMostSeq13(c)...)
		clauses = append(clauses, createAtMostSeq24(c)...)
		clauses = append(clauses, createAtMostSeq5(c)...)
		clauses = append(clauses, createAtMostSeq6(c)...)
	}

	for _, o := range options {
		clauses = append(clauses, createAtMostSeq13(o)...)
		clauses = append(clauses, createAtMostSeq24(o)...)
		clauses = append(clauses, createAtMostSeq5(o)...)
		clauses = append(clauses, createAtMostSeq6(o)...)
	}
	
    for i := 0; i < class_count; i++ {
	    for j := 0; j < option_count; j++ {
            if class2option[i][j] {
		        clauses = append(clauses, createAtMostSeq7(classes[i].cId,options[j].cId)...)
            }
        }
    }
    
	for j := 0; j < option_count; j++ {

        ops := make([]CountableId,0,class_count)

        for i := 0; i < class_count; i++ {
            if class2option[i][j] {
                k := len(ops)
                ops = ops[:k+1]
                ops[k] = classes[i].cId
            }
        }
		clauses = append(clauses, createAtMostSeq8(options[j].cId, ops)...)
    }

	fmt.Println("number of clauses: ", len(clauses))
	fmt.Println("number of pos variables: ", len(gen.posVarMap))
	fmt.Println("number of count variables: ", len(gen.countVarMap))

	return true
}

func createAtMostSeq13(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 CountVar
    var pV PosVar
    
    pV.cId = c.cId
	cV1.cId = c.cId
	cV2.cId = c.cId

	for i := 0; i < size-1; i++ {

		cV1.pos = i
		cV2.pos = i + 1
        pV.pos = i+1

		for j := c.lower[i]; j <= c.upper[i]; j++ {
			cV1.count = j
			cV2.count = j
			c1 := clause{-1 * getCountId(cV1), getCountId(cV2)}
			c3 := clause{-getPosId(pV), getCountId(cV1), -getCountId(cV2)}
			clauses = append(clauses, c1, c3)
		}
	}

	//fmt.Printf("13 for %v added clauses %v\n", c.cId, len(clauses))

	return
}

func createAtMostSeq24(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 CountVar
    var pV PosVar
    
    pV.cId = c.cId
	cV1.cId = c.cId
	cV2.cId = c.cId

	for i := 0; i < size-1; i++ {

		cV1.pos = i
		cV2.pos = i + 1
        pV.pos = i+1

		for j := c.lower[i]; j < c.upper[i]; j++ {
			cV1.count = j
			cV2.count = j + 1
			c2 := clause{getCountId(cV1), -getCountId(cV2)}
			c4 := clause{getPosId(pV), getCountId(cV1), -getCountId(cV2)}
			clauses = append(clauses, c2, c4)
		}
	}

	//fmt.Printf("24 for %v added clauses %v\n", c.cId, len(clauses))

	return
}

func createAtMostSeq5(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV CountVar
	cV.cId = c.cId

	for i := 0; i < size; i++ {
		cV.pos = i

		cV.count = c.lower[i]
		cn := clause{getCountId(cV)}
		clauses = append(clauses, cn)

		cV.count = c.upper[i]
		cn = clause{-getCountId(cV)}
		clauses = append(clauses, cn)
	}

	//fmt.Printf("5 for %v added clauses %v\n", c.cId, len(clauses))

	return
}

func createAtMostSeq6(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 CountVar

	cV1.cId = c.cId
	cV2.cId = c.cId
    q := c.window
    u := c.capacity

	for i := q; i < size; i++ {

		cV1.pos = i-q
		cV2.pos = i

		for j := c.lower[i-q]; j < c.upper[i-q]; j++ {
			cV1.count = j
			cV2.count = j + u
            if j+u < c.upper[i] {
			    cn := clause{getCountId(cV1), -getCountId(cV2)}
			    clauses = append(clauses, cn)
            }
		}
	}

	//fmt.Printf("6 for %v added clauses %v\n", c.cId, len(clauses))

	return
}

func createAtMostSeq7(cId1 CountableId, cId2 CountableId) (clauses []clause) {
    
    var pV1, pV2  PosVar
    
    pV1.cId = cId1
    pV2.cId = cId2

	for i := 0; i < size; i++ {
		pV1.pos = i
		pV2.pos = i
		clauses = append(clauses, clause{-getPosId(pV1), getPosId(pV2)})
    }

	//fmt.Printf("7 added clauses %v\n", len(clauses))

    return
}

func createAtMostSeq8(cId1 CountableId, cId2s []CountableId) (clauses []clause) {
    
    var pV1  PosVar
    
    pV1.cId = cId1

	for i := 0; i < size; i++ {
		pV1.pos = i

        c := make(clause,len(cId2s)+1)
		c[0] = -getPosId(pV1)

        for j,id := range cId2s {
		    c[j+1] = getPosId(PosVar{id,i})
        }

        clauses = append(clauses, c)
    }

	//fmt.Printf("8 added clauses %v\n", len(clauses))

    return
}

func (c *Countable) createBounds() {
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
