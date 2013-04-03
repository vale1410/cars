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

var name = flag.String("file", "test.txt", "Name of the file specifing the car sequencing according to the CSPlib.")
var sym = flag.Bool("sym", false, "Order the sequence in one direction")
var id8 = flag.Bool("id8", false, "Class to Option relations, alternative 1: Completion Clause. (alternative to id9)")
var id9 = flag.Bool("id9", false, "Class to Option relations, alternative 2: class implies neg options. Adds binary clauses (alternative to id8) ")
var re1 = flag.Bool("re1", false, "Implications of options that are of the form 1/q. Adds binary clauses")
var re2 = flag.Bool("re2", false, "Implications of options that are of the form 2/q. Adds binary clauses")
var sbd = flag.Bool("sbd", false, "For initial grounding use simple bounds to generate counters.  (this is needed for correct optimization version of the sequencing problem")
var debug = flag.Bool("debug", false, "Adds debug information to the cnf (symbol table and textual clauses) ")

var digitRegexp = regexp.MustCompile("([0-9]+ )*[0-9]+")

var size, class_count, option_count int
var gen IdGen

func main() {
	flag.Parse()
	parse(*name)
}

const (
	optionType countType = iota
	classType
	exactlyOne
	optType
)

type countType int
type clause struct {
	desc     string
	literals []int
}

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
	gen.posVarMap = make(map[PosVar]int, size*(class_count+option_count))
	gen.countVarMap = make(map[CountVar]int, size*class_count^2)
	return
}

func printClausesDIMACS(clauses []clause) {

	fmt.Printf("p cnf %v %v\n", len(gen.posVarMap)+len(gen.countVarMap), len(clauses))

	for _, c := range clauses {
		for _, l := range c.literals {
			fmt.Printf("%v ", l)
		}
		fmt.Printf("0\n")
	}
}

func printDebug(clauses []clause) {

	symbolTable := make([]string, len(gen.countVarMap)+len(gen.posVarMap)+1)

	for key, valueInt := range gen.posVarMap {
		s := ""
		switch key.cId.typ {
		case optionType:
			s = "pos(option,"
		case classType:
			s = "pos(class,"
		case exactlyOne:
			s = "pos(aux,"
		}
		s += strconv.Itoa(key.cId.index)
		s += ","
		s += strconv.Itoa(key.pos)
		s += ")"
		symbolTable[valueInt] = s
	}

	for key, valueInt := range gen.countVarMap {
		s := ""
		switch key.cId.typ {
		case optionType:
			s = "count(option,"
		case classType:
			s = "count(class,"
		}
		s += strconv.Itoa(key.cId.index)
		s += ","
		s += strconv.Itoa(key.pos)
		s += ","
		s += strconv.Itoa(key.count)
		s += ")"
		symbolTable[valueInt] = s
	}

	fmt.Println("c pos(Type,Id,Position).")
	fmt.Println("c count(Type,Id,Position,Count).")
	for i, s := range symbolTable {
		fmt.Println("c", i, "\t:", s)
	}

	stat := make(map[string]int, 10)

	for _, c := range clauses {

		count, ok := stat[c.desc]
		if !ok {
			stat[c.desc] = 1
		} else {
			stat[c.desc] = count + 1
		}

		fmt.Printf("c %s ", c.desc)
		first := true
		for _, l := range c.literals {
			if !first {
				fmt.Printf(",")
			}
			first = false
			if l < 0 {
				fmt.Printf("-%s", symbolTable[-l])
			} else {
				fmt.Printf("+%s", symbolTable[l])
			}
		}
		fmt.Println(".")
	}

	all := []string{"id1",
		"id2",
		"id3",
		"id4",
		"id5",
		"id6",
		"id7",
		"id8",
		"id9",
		"lt1",
		"gt1",
		"sym",
		"re1",
		"re2"}

	for _, key := range all {
		fmt.Printf("c %v\t: %v\t%.1f \n", key, stat[key], 100*float64(stat[key])/float64(len(clauses)))
	}
	fmt.Printf("c %v\t: %v\t%.1f \n", "tot", len(clauses), 100.0)
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
		fmt.Println("Could not read file: ", filename)
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
		numbers := strings.Split(strings.TrimSpace(l), " ")
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
				}
			}
			state++
		} else {
			fmt.Println("c ", l)
		}
	}

	for i := range options {
		options[i].createBounds()
	}

	fmt.Println("c options: ", options)
	fmt.Println("c classes: ", classes)
	fmt.Println("c class2option: ", class2option)

	NewIdGen()

	clauses := make([]clause, 0)

	//clauses 1-6 for classes
	for _, c := range classes {
		clauses = append(clauses, createAtMostSeq13(c)...)
		clauses = append(clauses, createAtMostSeq24(c)...)
		clauses = append(clauses, createAtMostSeq5(c)...)
		clauses = append(clauses, createAtMostSeq6(c)...)
	}

	//clauses 1-6 for options
	for _, o := range options {
		clauses = append(clauses, createAtMostSeq13(o)...)
		clauses = append(clauses, createAtMostSeq24(o)...)
		clauses = append(clauses, createAtMostSeq5(o)...)
		clauses = append(clauses, createAtMostSeq6(o)...)
		if *re1 {
			clauses = append(clauses, createRedundant1(o)...)
		}
		if *re2 {
			clauses = append(clauses, createRedundant2(o)...)
		}
	}

	//clauses 7 and 9
	for i := 0; i < class_count; i++ {
		for j := 0; j < option_count; j++ {
			if class2option[i][j] {
				clauses = append(clauses, createAtMostSeq7(classes[i].cId, options[j].cId)...)
			} else {
				if *id9 {
					clauses = append(clauses, createAtMostSeq9(classes[i].cId, options[j].cId)...)
				}
			}
		}
	}

	//clauses 8
	if *id8 {
		for j := 0; j < option_count; j++ {

			ops := make([]CountableId, 0, class_count)

			for i := 0; i < class_count; i++ {
				if class2option[i][j] {
					k := len(ops)
					ops = ops[:k+1]
					ops[k] = classes[i].cId
				}
			}
			clauses = append(clauses, createAtMostSeq8(options[j].cId, ops)...)
		}
	}

	//clauses exaclty one class per position
	clauses = append(clauses, createExactlyOne()...)

	//symmetry breaking
	if *sym {
		clauses = append(clauses, createSymmetry()...)
	}

	//fmt.Println("number of clauses: ", len(clauses))
	//fmt.Println("number of pos variables: ", len(gen.posVarMap))
	//fmt.Println("number of count variables: ", len(gen.countVarMap))

	printClausesDIMACS(clauses)
	if *debug {
		printDebug(clauses)
	}

	return true
}

func createAtMostSeq13(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	pV := PosVar{c.cId, 0}
	cV2 := CountVar{c.cId, 0, 1}

	cn := clause{"id3", []int{getPosId(pV), -getCountId(cV2)}}
	clauses = append(clauses, cn)

	for i := 0; i < size-1; i++ {

		cV1 := CountVar{c.cId, i, -1}
		cV2.pos = i + 1
		pV.pos = i + 1

		for j := c.lower[i]; j <= c.upper[i]; j++ {
			cV1.count = j
			cV2.count = j
			c1 := clause{"id1", []int{-1 * getCountId(cV1), getCountId(cV2)}}
			c3 := clause{"id3", []int{getPosId(pV), getCountId(cV1), -getCountId(cV2)}}
			clauses = append(clauses, c1, c3)
		}
	}

	return
}

func createAtMostSeq24(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	pV := PosVar{c.cId, 0}
	cV2 := CountVar{c.cId, 0, 1}

	cn := clause{"id4", []int{-getPosId(pV), getCountId(cV2)}}
	clauses = append(clauses, cn)

	for i := 0; i < size-1; i++ {

		cV1 := CountVar{c.cId, i, -1}
		cV2.pos = i + 1
		pV.pos = i + 1

		for j := c.lower[i]; j < c.upper[i]; j++ {
			cV1.count = j
			cV2.count = j + 1
			c2 := clause{"id2", []int{getCountId(cV1), -getCountId(cV2)}}
			c4 := clause{"id4", []int{-getPosId(pV), -getCountId(cV1), getCountId(cV2)}}
			clauses = append(clauses, c2, c4)
		}
	}

	return
}

func createAtMostSeq5(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV CountVar
	cV.cId = c.cId

	for i := 0; i < size; i++ {
		cV.pos = i

		cV.count = c.lower[i]
		cn := clause{"id5", []int{getCountId(cV)}}
		clauses = append(clauses, cn)

		cV.count = c.upper[i]
		cn = clause{"id5", []int{-getCountId(cV)}}
		clauses = append(clauses, cn)
	}

	return
}

func createAtMostSeq6(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 CountVar

	cV1.cId = c.cId
	cV2.cId = c.cId
	q := c.window
	u := c.capacity

	if *sbd {
		// needed because I tried to avoid the first column, now extra work for sbd
		cV1.pos = q - 1
		cV1.count = u + 1
		cn := clause{"id6", []int{-getCountId(cV1)}}
		clauses = append(clauses, cn)

	}

	for i := q; i < size; i++ {

		cV1.pos = i - q
		cV2.pos = i

		for j := c.lower[i-q]; j < c.upper[i-q]; j++ {
			cV1.count = j
			cV2.count = j + u
			if c.lower[i] <= j+u && j+u < c.upper[i] {
				cn := clause{"id6", []int{getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, cn)
			}
		}
	}

	return
}

func createAtMostSeq7(cId1 CountableId, cId2 CountableId) (clauses []clause) {

	var pV1, pV2 PosVar

	pV1.cId = cId1
	pV2.cId = cId2

	for i := 0; i < size; i++ {
		pV1.pos = i
		pV2.pos = i
		clauses = append(clauses, clause{"id7", []int{-getPosId(pV1), getPosId(pV2)}})
	}

	return
}

func createAtMostSeq8(cId1 CountableId, cId2s []CountableId) (clauses []clause) {

	var pV1 PosVar

	pV1.cId = cId1

	for i := 0; i < size; i++ {
		pV1.pos = i

		c := make([]int, len(cId2s)+1)
		c[0] = -getPosId(pV1)

		for j, id := range cId2s {
			c[j+1] = getPosId(PosVar{id, i})
		}

		clauses = append(clauses, clause{"id8", c})
	}

	return
}

func createAtMostSeq9(cId1 CountableId, cId2 CountableId) (clauses []clause) {

	var pV1, pV2 PosVar

	pV1.cId = cId1
	pV2.cId = cId2

	for i := 0; i < size; i++ {
		pV1.pos = i
		pV2.pos = i
		clauses = append(clauses, clause{"id9", []int{-getPosId(pV1), -getPosId(pV2)}})
	}

	return
}

func createExactlyOne() (clauses []clause) {

	clauses = make([]clause, 0)

	var posV1, posV2, auxV1, auxV2 PosVar

	for i := 0; i < size; i++ {

		posV1.pos = i
		posV2.pos = i
		auxV1.pos = i
		auxV2.pos = i

		atLeastOne := make([]int, class_count)

		for j := 0; j < class_count-1; j++ {

			posV1.cId = CountableId{classType, j}
			posV2.cId = CountableId{classType, j + 1}
			atLeastOne[j] = getPosId(posV1)

			auxV1.cId = CountableId{exactlyOne, j}
			auxV2.cId = CountableId{exactlyOne, j + 1}

			c1 := clause{"lt1", []int{-getPosId(posV1), getPosId(auxV1)}}
			c2 := clause{"lt1", []int{-getPosId(posV2), -getPosId(auxV1)}}
			clauses = append(clauses, c1, c2)
			if j < class_count-2 {
				c3 := clause{"lt1", []int{-getPosId(auxV1), getPosId(auxV2)}}
				clauses = append(clauses, c3)
			}

		}

		atLeastOne[class_count-1] = getPosId(posV2)
		clauses = append(clauses, clause{"gt1", atLeastOne})

	}

	return
}

func createSymmetry() (clauses []clause) {

	var pV1, pVn PosVar

	pV1.pos = 0
	pVn.pos = size - 1

	for i := 0; i < class_count-1; i++ {

		pV1.cId = CountableId{exactlyOne, i}
		pVn.cId = CountableId{exactlyOne, i}

		clauses = append(clauses, clause{"sym", []int{getPosId(pV1), -getPosId(pVn)}})
	}

	pV1.cId = CountableId{classType, class_count - 1}
	pVn.cId = CountableId{classType, class_count - 1}
	pVn2 := PosVar{CountableId{exactlyOne, class_count - 2}, size - 1}

	clauses = append(clauses, clause{"sym", []int{getPosId(pV1), -getPosId(pVn), -getPosId(pVn2)}})

	return
}

func createRedundant1(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var pV1, pV2 PosVar

	pV1.cId = c.cId
	pV2.cId = c.cId

	q := c.window
	u := c.capacity

	if u == 1 {
		for i := 0; i < size; i++ {

			pV1.pos = i

			for j := i + 1; j < i+q && j < size; j++ {
				pV2.pos = j
				cn := clause{"re1", []int{-getPosId(pV1), -getPosId(pV2)}}
				clauses = append(clauses, cn)
			}
		}
	}

	return
}

func createRedundant2(c Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	q := c.window
	u := c.capacity

	if u == 2 {

		var pV1, pV2, pV3 PosVar

		pV1.cId = c.cId
		pV2.cId = c.cId
		pV3.cId = c.cId

		for i := 0; i < size; i++ {

			pV1.pos = i

			for j := i + 1; j < i+q && j < size; j++ {

				pV2.pos = j

				for k := j + 1; k < i+q && k < size; k++ {

					pV3.pos = k

					cn := clause{"re2", []int{-getPosId(pV1), -getPosId(pV2), -getPosId(pV3)}}
					clauses = append(clauses, cn)
				}
			}
		}
	}

	return
}

//func createOpt1(c Countable) (clauses []clause) {
//
//	clauses = make([]clause, 0)
//
//	var cV1, cV2 CountVar
//	var optV PosVar
//
//	cV1.cId = c.cId
//	cV2.cId = c.cId
//	optV.cId = c.cId
//
//	q := c.window
//	u := c.capacity
//
//	for i := q; i < size; i++ {
//
//		cV1.pos = i - q
//		optV.pos = i - q
//		cV2.pos = i
//
//		for j := c.lower[i-q]; j < c.upper[i-q]; j++ {
//			cV1.count = j
//			cV2.count = j + u
//			if j+u < c.upper[i] {
//				cn := clause{"id6", []int{-getPosId(optV), getCountId(cV1), -getCountId(cV2)}}
//				clauses = append(clauses, cn)
//			}
//		}
//	}
//
//	return
//}

func (c *Countable) createBounds() {
	if *sbd {
		c.createSimpleBounds()
	} else {
		c.createImprovedBounds()
	}
}

func (c *Countable) createSimpleBounds() {

	c.lower = make([]int, size)
	c.upper = make([]int, size)

	h := c.demand

	for i := size - 1; i >= 0; i-- {
		c.lower[i] = h
		if h > 0 {
			h--
		}
	}

	h = 2

	for i := 0; i < size; i++ {
		c.upper[i] = h
		if h <= c.demand {
			h++
		}
	}
}

func (c *Countable) createImprovedBounds() {
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
	q := c.window - 1
	u := c.capacity - 1

	for i := 0; i < size; i++ {

		for i < size {

			c.upper[i] = h + 1

			if u > 0 && h < c.demand {
				u--
				h++
			}
			q--
			if q <= 0 {
				break
			}
			i++
		}

		q = c.window
		u = c.capacity

	}
}
