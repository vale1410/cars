package main

import (
	"../base"
	"../pbo"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var name = flag.String("f", "", "Path of the file specifying the car sequencing according to the CSPlib.")
var ver = flag.Bool("ver", false, "Show version info.")
var e1 = flag.Bool("e1", false, "Collection of flags: -ex1 -cnt -ca    -id7 -id8 -id9.")
var e2 = flag.Bool("e2", false, "Collection of flags: -ex1 -cnt -id6 3 -id7 -id8 -id9.")
var e3 = flag.Bool("e3", false, "Collection of flags: -ex1 -cnt -ca -id6 3 -id7 -id8 -id9.")
var ex1 = flag.Bool("ex1", false, "Adds clauses to state that in each position there is exactly one car.")
var cnt = flag.Bool("cnt", false, "Meta flag: sets id1, id2, id3, id4, id5.")
var id1 = flag.Bool("id1", false, "Sequential Counter for cardinality, clauses 1 (see paper).")
var id2 = flag.Bool("id2", false, "Sequential Counter for cardinality, clauses 2 (see paper).")
var id3 = flag.Bool("id3", false, "Sequential Counter for cardinality, clauses 3 (see paper).")
var id4 = flag.Bool("id4", false, "Sequential Counter for cardinality, clauses 4 (see paper).")
var id5 = flag.Bool("id5", false, "Sequential Counter for cardinality, clauses 4 (see paper).")
var ca = flag.Bool("ca", false, "Meta flag: sets ca1, ca2, ca3, ca4, ca5.")
var ca1 = flag.Bool("ca1", false, "Sequential Counter for Capacity constraints, clauses 1 (see paper).")
var ca2 = flag.Bool("ca2", false, "Sequential Counter for Capacity constraints, clauses 2 (see paper).")
var ca3 = flag.Bool("ca3", false, "Sequential Counter for Capacity constraints, clauses 3 (see paper).")
var ca4 = flag.Bool("ca4", false, "Sequential Counter for Capacity constraints, clauses 4 (see paper).")
var ca5 = flag.Bool("ca5", false, "Initializes the counter to be a AtMost constraint.")
var id6 = flag.Int("id6", 0, "AtMostSeqCard reusing the aux variables of cardinality constraints on the demand. 0: none, 1: just on options; 2: just on classes; 3: on options and classes")
var id7 = flag.Bool("id7", false, "Implications from Classes to Options.")
var id8 = flag.Bool("id8", false, "Class to Option relations, alternative 1: Completion Clause. (alternative to id9)")
var id9 = flag.Bool("id9", false, "Class to Option relations, alternative 2: class implies neg options. Adds binary"+
	" clauses (alternative to id8).")
var re1 = flag.Bool("re1", false, "Implications of options that are of the form 1/q. Adds binary clauses.")
var re2 = flag.Bool("re2", false, "Implications of options that are of the form 2/q. Adds binary clauses.")
var sym = flag.Bool("sym", false, "Order the sequence in one direction (first car has lower or equal class id than last).")
var ian = flag.Bool("ian", false, "A redundant constraint that precomputes sets of classes that need to be neighbours.")
var sbd = flag.Bool("sbd", false, "For initial grounding use simple bounds to generate counters.  (context optimization).")
var opt = flag.Int("opt", -1, "Adds optimization statement with value given. Should be used with -sbd and without -re1 -re2.")
var add = flag.Int("add", 0, "Add n dummy cars without any option. (simulates optimization).")
var debug = flag.Bool("debug", false, "Adds debug information to the cnf (symbol table and textual clauses).")
var symbolsFile = flag.String("symbols", "", "Outputs the symbol table; meaning of the variables.")
var pb = flag.Bool("pbo", false, "Create Pseudo Boolean Model; simple version")

var digitRegexp = regexp.MustCompile("([0-9]+ )*[0-9]+")

var size, class_count, option_count int
var gen IdGen

func main() {
	flag.Parse()
	if *ver {
		fmt.Println(`CNF generator for car sequencing problem from CSPLib 
Version tag: 1.3a
For infos about flags use -help
Copyright (C) NICTA and Valentin Mayer-Eichberger
License GPLv2+: GNU GPL version 2 or later <http://gnu.org/licenses/gpl.html>
There is NO WARRANTY, to the extent permitted by law.`)
		return
	}
	setFlags()
	if *name == "" {
		*name = flag.Arg(0)
	}
	options, classes, class2option := parse(*name)

	if *add > 0 {
		cId := base.CountableId{base.ClassType, class_count}
		dummy := base.Countable{CId: cId, Window: 1, Capacity: 1, Demand: *add}
		classes = append(classes, dummy)
		class2option = append(class2option, make([]bool, option_count))
		class_count++
		size += *add
	}

	if *pb {
		pbo.CreatePBOModel(size, options, classes, class2option)
	} else {
		createSATModel(options, classes, class2option)
	}
}

func setFlags() {
	t := true
	n3 := 3

	if *e1 {
		ex1 = &t
		cnt = &t
		ca = &t
		id7 = &t
		id8 = &t
		id9 = &t
	}

	if *e2 {
		ex1 = &t
		cnt = &t
		id6 = &n3
		id7 = &t
		id8 = &t
		id9 = &t
	}

	if *e3 {
		ex1 = &t
		cnt = &t
		ca = &t
		id6 = &n3
		id7 = &t
		id8 = &t
		id9 = &t
	}

	if *ca {
		ca1 = &t
		ca2 = &t
		ca3 = &t
		ca4 = &t
		ca5 = &t
	}

	if *cnt {
		id1 = &t
		id2 = &t
		id3 = &t
		id4 = &t
		id5 = &t
	}
}

type clause struct {
	desc     string
	literals []int
}

type IdGen struct {
	Id           int
	CountVarMap  map[base.CountVar]int
	PosVarMap    map[base.PosVar]int
	AtMostVarMap map[base.AtMostVar]int
}

func NewIdGen() {
	gen.Id = 0
	gen.PosVarMap = make(map[base.PosVar]int, size*(class_count+option_count)) //just an approximation of size of map
	gen.CountVarMap = make(map[base.CountVar]int, size*class_count^2)          //just an approximation of size of map
	gen.AtMostVarMap = make(map[base.AtMostVar]int, size*class_count^2)        //just an approximation of size of map
	return
}

func printClausesDIMACS(clauses []clause) {

	fmt.Printf("p cnf %v %v\n", len(gen.PosVarMap)+len(gen.CountVarMap)+len(gen.AtMostVarMap), len(clauses))

	for _, c := range clauses {
		for _, l := range c.literals {
			fmt.Printf("%v ", l)
		}
		fmt.Printf("0\n")
	}
}

func generateSymbolTable() []string {

	symbolTable := make([]string, len(gen.CountVarMap)+len(gen.PosVarMap)+len(gen.AtMostVarMap)+1)

	for key, valueInt := range gen.PosVarMap {
		s := ""
		switch key.CId.Typ {
		case base.OptionType:
			s = "pos(option,"
		case base.ClassType:
			s = "pos(class,"
		case base.ExactlyOne:
			s = "pos(aux,"
		case base.OptimizationType:
			s = "pos(opti,"
		}
		s += strconv.Itoa(key.CId.Index)
		s += ","
		s += strconv.Itoa(key.Pos)
		s += ")"
		symbolTable[valueInt] = s
	}

	for key, valueInt := range gen.CountVarMap {
		s := ""
		switch key.CId.Typ {
		case base.OptionType:
			s = "count(option,"
		case base.ClassType:
			s = "count(class,"
		case base.OptimizationType:
			s = "count(opti,"
		}
		s += strconv.Itoa(key.CId.Index)
		s += ","
		s += strconv.Itoa(key.Pos)
		s += ","
		s += strconv.Itoa(key.Count)
		s += ")"
		symbolTable[valueInt] = s
	}

	for key, valueInt := range gen.AtMostVarMap {
		s := ""
		switch key.CId.Typ {
		case base.OptionType:
			s = "atMost(option,"
		case base.ClassType:
			s = "atMost(class,"
		case base.OptimizationType:
			s = "atMost(opti,"
		}
		s += strconv.Itoa(key.CId.Index)
		s += ","
		s += strconv.Itoa(key.First)
		s += ","
		s += strconv.Itoa(key.Pos)
		s += ","
		s += strconv.Itoa(key.Count)
		s += ")"
		symbolTable[valueInt] = s
	}

	return symbolTable
}

func printSymbolTable(symbolTable []string, filename string) {

	symbolFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := symbolFile.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(symbolFile)

	for i, s := range symbolTable {
		// write a chunk
		if _, err := w.Write([]byte(fmt.Sprintln(i, "\t:", s))); err != nil {
			panic(err)
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}

}

func printDebug(symbolTable []string, clauses []clause) {

	// first print symbol table into file
	fmt.Println("c pos(Type,Id,Position).")
	fmt.Println("c count(Type,Id,Position,Count).")
	fmt.Println("c atMost(Type,Id,First,Position,Count).")
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
		"ca1",
		"ca2",
		"ca3",
		"ca4",
		"ca5",
		"lt1",
		"gt1",
		"sym",
		"re1",
		"re2",
		"op0",
		"op1",
		"op2",
		"op3",
		"op4"}

	for _, key := range all {
		fmt.Printf("c %v\t: %v\t%.1f \n", key, stat[key], 100*float64(stat[key])/float64(len(clauses)))
	}
	fmt.Printf("c %v\t: %v\t%.1f \n", "tot", len(clauses), 100.0)
}

func getCountId(v base.CountVar) (id int) {
	id, b := gen.CountVarMap[v]

	if !b {
		gen.Id++
		id = gen.Id
		gen.CountVarMap[v] = id
	}
	return id
}

func getPosId(v base.PosVar) (id int) {
	id, b := gen.PosVarMap[v]

	if !b {
		gen.Id++
		id = gen.Id
		gen.PosVarMap[v] = id
	}
	return id
}

func getAtMostId(v base.AtMostVar) (id int) {
	id, b := gen.AtMostVarMap[v]

	if !b {
		gen.Id++
		id = gen.Id
		gen.AtMostVarMap[v] = id
	}
	return id
}

func parse(filename string) (options, classes []base.Countable, class2option [][]bool) {
	input, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Please specifiy correct path to instance. File does not exist: ", filename)
		return
	}

	b := bytes.NewBuffer(input)

	lines := strings.Split(b.String(), "\n")

	state := 0

	for _, l := range lines {
		numbers := strings.Split(strings.TrimSpace(l), " ")
		if digitRegexp.MatchString(numbers[0]) {
			switch state {
			case 0:
				{
					size, _ = strconv.Atoi(numbers[0])
					option_count, _ = strconv.Atoi(numbers[1])
					options = make([]base.Countable, option_count)
					class_count, _ = strconv.Atoi(numbers[2])
					classes = make([]base.Countable, class_count)
					class2option = make([][]bool, class_count)
				}
			case 1:
				{
					for i, v := range numbers {
						capacity, _ := strconv.Atoi(v)
						options[i].CId = base.CountableId{base.OptionType, i}
						options[i].Capacity = capacity
					}
				}
			case 2:
				{
					for i, v := range numbers {
						window, _ := strconv.Atoi(v)
						options[i].Window = window
					}
				}
			default:
				{
					num, _ := strconv.Atoi(numbers[0])
					classes[num].CId = base.CountableId{base.ClassType, num}
					class2option[num] = make([]bool, option_count)

					// find option with lowest slope
					// to determine capacity and windows

					classes[num].Capacity = 1
					classes[num].Window = 1
					slope := 1.0

					for i, v := range numbers {
						if i == 1 {
							demand, _ := strconv.Atoi(v)
							classes[num].Demand = demand
						} else if i > 1 {
							value, _ := strconv.Atoi(v)
							has_option := value == 1
							class2option[num][i-2] = has_option
							if has_option {
								options[i-2].Demand += classes[num].Demand
								slope2 := float64(options[i-2].Capacity) / float64(options[i-2].Window)
								if slope2 < slope {
									slope = slope2
									classes[num].Capacity = options[i-2].Capacity
									classes[num].Window = options[i-2].Window
								}
							}
						}
					}
				}
			}
			state++
		} else {
			//fmt.Println("c ", l)
		}
	}
	return

}

// TODO: this is just copy paste, needs added clauses from alternative stuff
// also this will also be the first time a sorting network will be used.
func createAltSATModel(options, classes []base.Countable, class2option [][]bool) bool {

	NewIdGen()

	clauses := make([]clause, 0)

	//clauses 7 and 9
	for i := 0; i < class_count; i++ {
		for j := 0; j < option_count; j++ {
			if class2option[i][j] {
				if *id7 {
					clauses = append(clauses, createAtMostSeq7(classes[i].CId, options[j].CId)...)
				}
			} else {
				if *id9 {
					clauses = append(clauses, createAtMostSeq9(classes[i].CId, options[j].CId)...)
				}
			}
		}
	}

	//clauses 8
	if *id8 {
		for j := 0; j < option_count; j++ {

			ops := make([]base.CountableId, 0, class_count)

			for i := 0; i < class_count; i++ {
				if class2option[i][j] {
					k := len(ops)
					ops = ops[:k+1]
					ops[k] = classes[i].CId
				}
			}
			clauses = append(clauses, createAtMostSeq8(options[j].CId, ops)...)
		}
	}

	//clauses exactly one class per position
	if *ex1 {
		clauses = append(clauses, createExactlyOne()...)
	}

	//symmetry breaking
	if *sym {
		clauses = append(clauses, createSymmetry()...)
	}

	//fmt.Println("number of clauses: ", len(clauses))
	//fmt.Println("number of pos variables: ", len(gen.PosVarMap))
	//fmt.Println("number of count variables: ", len(gen.CountVarMap))

	if *ian {
		createIanConstraints(options, classes, class2option)
	}

	if len(clauses) > 0 {
		printClausesDIMACS(clauses)
	}

	if *debug || *symbolsFile != "" {

		symbolTable := generateSymbolTable()

		if *debug {
			fmt.Println("c options: ", options)
			fmt.Println("c classes: ", classes)
			fmt.Println("c class2option: ", class2option)
			printDebug(symbolTable, clauses)
		}

		if *symbolsFile != "" {
			printSymbolTable(symbolTable, *symbolsFile)
		}
	}

	return true
}

// TODO: this is just copy paste, needs added clauses from alternative stuff
// also this will also be the first time a sorting network will be used.
func createAlternative(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	pV := base.PosVar{c.CId, 0}
	cV2 := base.CountVar{c.CId, 0, 1}

	if *id3 {
		cn := clause{"id3", []int{getPosId(pV), -getCountId(cV2)}}
		clauses = append(clauses, cn)
	}

	for i := 0; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j <= c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j
			if *id1 {
				c1 := clause{"id1", []int{-1 * getCountId(cV1), getCountId(cV2)}}
				clauses = append(clauses, c1)
			}
			if *id3 {
				c3 := clause{"id3", []int{getPosId(pV), getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, c3)
			}
		}
	}

	pV = base.PosVar{c.CId, 0}
	cV2 = base.CountVar{c.CId, 0, 1}

	if *id4 {
		cn := clause{"id4", []int{-getPosId(pV), getCountId(cV2)}}
		clauses = append(clauses, cn)
	}

	for i := 0; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j < c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j + 1
			if *id2 {
				c2 := clause{"id2", []int{getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, c2)
			}
			if *id4 {
				c4 := clause{"id4", []int{-getPosId(pV), -getCountId(cV1), getCountId(cV2)}}
				clauses = append(clauses, c4)
			}
		}
	}

	return
}

func createSATModel(options, classes []base.Countable, class2option [][]bool) bool {

	for i := range options {
		if *sbd {
			options[i].ComputeSimpleBounds(size)
		} else {
			options[i].ComputeImprovedBounds(size)
		}
	}

	for i := range classes {
		if *sbd {
			classes[i].ComputeSimpleBounds(size)
		} else {
			classes[i].ComputeImprovedBounds(size)
		}
	}

	NewIdGen()

	clauses := make([]clause, 0)

	//clauses 1-6 for classes
	for _, c := range classes {
		if *id1 || *id2 || *id3 || *id4 {
			clauses = append(clauses, createCounter(c)...)
		}
		if *id5 {
			clauses = append(clauses, createAtMostSeq5(c)...)
		}
		if *id6 == 2 || *id6 == 3 {
			clauses = append(clauses, createAtMostSeq6(c)...)
		}
	}

	//clauses 1-6 for options
	for _, o := range options {
		if *ca1 || *ca2 || *ca3 || *ca4 || *ca5 {
			clauses = append(clauses, createCapacityConstraints(o)...)
		}
		if *id1 || *id2 || *id3 || *id4 {
			clauses = append(clauses, createCounter(o)...)
		}
		if *id5 {
			clauses = append(clauses, createAtMostSeq5(o)...)
		}
		if *id6 == 1 || *id6 == 3 {
			clauses = append(clauses, createAtMostSeq6(o)...)
		}
		if *re1 {
			clauses = append(clauses, createRedundant1(o)...)
		}
		if *re2 {
			clauses = append(clauses, createRedundant2(o)...)
		}
		if *opt > 0 {
			clauses = append(clauses, createOptPositions(o)...)
			clauses = append(clauses, createOptCounter(o)...)
		}
	}

	//clauses 7 and 9
	for i := 0; i < class_count; i++ {
		for j := 0; j < option_count; j++ {
			if class2option[i][j] {
				if *id7 {
					clauses = append(clauses, createAtMostSeq7(classes[i].CId, options[j].CId)...)
				}
			} else {
				if *id9 {
					clauses = append(clauses, createAtMostSeq9(classes[i].CId, options[j].CId)...)
				}
			}
		}
	}

	//clauses 8
	if *id8 {
		for j := 0; j < option_count; j++ {

			ops := make([]base.CountableId, 0, class_count)

			for i := 0; i < class_count; i++ {
				if class2option[i][j] {
					k := len(ops)
					ops = ops[:k+1]
					ops[k] = classes[i].CId
				}
			}
			clauses = append(clauses, createAtMostSeq8(options[j].CId, ops)...)
		}
	}

	//clauses exactly one class per position
	if *ex1 {
		clauses = append(clauses, createExactlyOne()...)
	}

	//symmetry breaking
	if *sym {
		clauses = append(clauses, createSymmetry()...)
	}

	//fmt.Println("number of clauses: ", len(clauses))
	//fmt.Println("number of pos variables: ", len(gen.PosVarMap))
	//fmt.Println("number of count variables: ", len(gen.CountVarMap))

	if *ian {
		createIanConstraints(options, classes, class2option)
	}

	if len(clauses) > 0 {
		printClausesDIMACS(clauses)
	}

	if *debug || *symbolsFile != "" {

		symbolTable := generateSymbolTable()

		if *debug {
			fmt.Println("c options: ", options)
			fmt.Println("c classes: ", classes)
			fmt.Println("c class2option: ", class2option)
			printDebug(symbolTable, clauses)
		}

		if *symbolsFile != "" {
			printSymbolTable(symbolTable, *symbolsFile)
		}
	}

	return true
}

func createIanConstraints(options []base.Countable, classes []base.Countable, class2option [][]bool) (clauses []clause) {

	first := make([]bool, option_count, option_count)

	sets := createSubSets(0, first)

	// how many 1/2
	cap12 := make([]int, len(sets))
	// how many 1/k, k > 2
	cap1k := make([]int, len(sets))
	// how many 2/k > 2
	cap2k := make([]int, len(sets))
	demands := make([]int, len(sets))
	supplies := make([]int, len(sets))
	rest := make([]int, len(sets))

	for s, set := range sets {
		// find max.Capacity among options
		for j, b := range set {
			if b && options[j].Window > 1 {
				if options[j].Window == 2 && options[j].Capacity == 1 {
					cap12[s]++
				} else if options[j].Window > 2 && options[j].Capacity == 1 {
					cap1k[s]++
				} else if options[j].Window > 2 && options[j].Capacity == 2 {
					cap2k[s]++
				}
			}
		}

		for i, class := range classes {
			alwaysfit := true
			neverfit := true
			for j, b := range set {
				if b {
					if class2option[i][j] {
						neverfit = false
					} else {
						alwaysfit = false
					}
				}
			}
			if alwaysfit {
				demands[s] += class.Demand
			}
			if neverfit {
				supplies[s] += class.Demand
			}
			if !alwaysfit && !neverfit {
				rest[s] += class.Demand
			}

		}
	}

	fmt.Println("-----\tcap12\tcap1k\tcap2k\tdemand\tsupply\trest")

	for s, set := range sets {

		restriction := false

		if cap12[s] > 0 && cap2k[s] == 0 && demands[s]-2 >= supplies[s] {
			restriction = true
		} else if cap1k[s] > 0 && cap2k[s] == 0 && 2*(demands[s]-1) >= supplies[s] {
			restriction = true
		} else if cap1k[s] > 0 && cap2k[s] == 1 && demands[s]-2 >= supplies[s] {
			restriction = true
		}

		if restriction {

			for i, b := range set {
				if b {
					fmt.Print(i)
				} else {
					fmt.Print("-")
				}

			}
			fmt.Println("\t", cap12[s], "\t", cap1k[s], "\t", cap2k[s], "\t", demands[s], "\t", supplies[s], "\t", rest[s])
		}

	}
	return
}

func createSubSets(i int, set []bool) (sets [][]bool) {

	if i == option_count {
		sets = make([][]bool, 1)
		sets[0] = set
		return
	}

	pos := make([]bool, len(set))
	neg := make([]bool, len(set))
	copy(pos, set)
	copy(neg, set)

	pos[i] = true
	i++
	sets = createSubSets(i, pos)
	sets = append(sets, createSubSets(i, neg)...)
	return
}

func createCapacityConstraints(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	// first,position,count

	q := c.Window
	u := c.Capacity

	for first := 0; first < size-q+1; first++ {

		// first,position,count
		cV1 := base.AtMostVar{c.CId, first, first, 0}
		cV2 := base.AtMostVar{c.CId, first, first, 1}
		pV := base.PosVar{c.CId, first}

		if *ca3 {
			cn := clause{"ca3", []int{getPosId(pV), -getAtMostId(cV2)}}
			clauses = append(clauses, cn)
		}

		for i := first; i < first+q-1; i++ {

			cV1.Pos = i
			cV2.Pos = i + 1
			pV.Pos = i + 1

			for j := 0; j <= u+1; j++ {
				cV1.Count = j
				cV2.Count = j
				if *ca1 {
					c1 := clause{"ca1", []int{-1 * getAtMostId(cV1), getAtMostId(cV2)}}
					clauses = append(clauses, c1)
				}
				if *ca3 {
					c3 := clause{"ca3", []int{getPosId(pV), getAtMostId(cV1), -getAtMostId(cV2)}}
					clauses = append(clauses, c3)
				}
			}
		}

		cV1 = base.AtMostVar{c.CId, first, first, 0}
		cV2 = base.AtMostVar{c.CId, first, first, 1}
		pV = base.PosVar{c.CId, first}

		if *ca4 {
			cn := clause{"ca4", []int{-getPosId(pV), getAtMostId(cV2)}}
			clauses = append(clauses, cn)
		}

		for i := first; i < first+q-1; i++ {

			cV1.Pos = i
			cV2.Pos = i + 1
			pV.Pos = i + 1

			for j := 0; j <= u; j++ {
				cV1.Count = j
				cV2.Count = j + 1

				if *ca2 {
					c2 := clause{"ca2", []int{getAtMostId(cV1), -getAtMostId(cV2)}}
					clauses = append(clauses, c2)
				}
				if *ca4 {
					c4 := clause{"ca4", []int{-getPosId(pV), -getAtMostId(cV1), getAtMostId(cV2)}}
					clauses = append(clauses, c4)
				}
			}
		}

		if *ca5 { //initialize

			cV1 := base.AtMostVar{c.CId, first, first, 2}
			cV2 := base.AtMostVar{c.CId, first, first + q - 1, u + 1}
			cV3 := base.AtMostVar{c.CId, first, first, 0}

			clauses = append(clauses, clause{"ca5", []int{-getAtMostId(cV1)}})
			clauses = append(clauses, clause{"ca5", []int{-getAtMostId(cV2)}})
			clauses = append(clauses, clause{"ca5", []int{getAtMostId(cV3)}})

		}

	}

	return

}

func createCounter(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	pV := base.PosVar{c.CId, 0}
	cV2 := base.CountVar{c.CId, 0, 1}

	if *id3 {
		cn := clause{"id3", []int{getPosId(pV), -getCountId(cV2)}}
		clauses = append(clauses, cn)
	}

	for i := 0; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j <= c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j
			if *id1 {
				c1 := clause{"id1", []int{-1 * getCountId(cV1), getCountId(cV2)}}
				clauses = append(clauses, c1)
			}
			if *id3 {
				c3 := clause{"id3", []int{getPosId(pV), getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, c3)
			}
		}
	}

	pV = base.PosVar{c.CId, 0}
	cV2 = base.CountVar{c.CId, 0, 1}

	if *id4 {
		cn := clause{"id4", []int{-getPosId(pV), getCountId(cV2)}}
		clauses = append(clauses, cn)
	}

	for i := 0; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j < c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j + 1
			if *id2 {
				c2 := clause{"id2", []int{getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, c2)
			}
			if *id4 {
				c4 := clause{"id4", []int{-getPosId(pV), -getCountId(cV1), getCountId(cV2)}}
				clauses = append(clauses, c4)
			}
		}
	}

	return
}

func createAtMostSeq5(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV base.CountVar
	cV.CId = c.CId

	for i := 0; i < size; i++ {
		cV.Pos = i

		cV.Count = c.Lower[i]
		cn := clause{"id5", []int{getCountId(cV)}}
		clauses = append(clauses, cn)

		cV.Count = c.Upper[i]
		cn = clause{"id5", []int{-getCountId(cV)}}
		clauses = append(clauses, cn)
	}

	return
}

func createAtMostSeq6(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 base.CountVar

	cV1.CId = c.CId
	cV2.CId = c.CId
	q := c.Window
	u := c.Capacity

	if *sbd {
		// needed because I tried to avoid the first column, now extra work for sbd
		cV1.Pos = q - 1
		cV1.Count = u + 1
		cn := clause{"id6", []int{-getCountId(cV1)}}
		clauses = append(clauses, cn)

	}

	for i := q; i < size; i++ {

		cV1.Pos = i - q
		cV2.Pos = i

		for j := c.Lower[i-q]; j < c.Upper[i-q]; j++ {
			cV1.Count = j
			cV2.Count = j + u
			if c.Lower[i] <= j+u && j+u < c.Upper[i] {
				cn := clause{"id6", []int{getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, cn)
			}
		}
	}

	return
}

func createAtMostSeq7(cId1 base.CountableId, cId2 base.CountableId) (clauses []clause) {

	var pV1, pV2 base.PosVar

	pV1.CId = cId1
	pV2.CId = cId2

	for i := 0; i < size; i++ {
		pV1.Pos = i
		pV2.Pos = i
		clauses = append(clauses, clause{"id7", []int{-getPosId(pV1), getPosId(pV2)}})
	}

	return
}

func createAtMostSeq8(cId1 base.CountableId, cId2s []base.CountableId) (clauses []clause) {

	var pV1 base.PosVar

	pV1.CId = cId1

	for i := 0; i < size; i++ {
		pV1.Pos = i

		c := make([]int, len(cId2s)+1)
		c[0] = -getPosId(pV1)

		for j, id := range cId2s {
			c[j+1] = getPosId(base.PosVar{id, i})
		}

		clauses = append(clauses, clause{"id8", c})
	}

	return
}

func createAtMostSeq9(cId1 base.CountableId, cId2 base.CountableId) (clauses []clause) {

	var pV1, pV2 base.PosVar

	pV1.CId = cId1
	pV2.CId = cId2

	for i := 0; i < size; i++ {
		pV1.Pos = i
		pV2.Pos = i
		clauses = append(clauses, clause{"id9", []int{-getPosId(pV1), -getPosId(pV2)}})
	}

	return
}

func createExactlyOne() (clauses []clause) {

	clauses = make([]clause, 0)

	var posV1, posV2, auxV1, auxV2 base.PosVar

	for i := 0; i < size; i++ {

		posV1.Pos = i
		posV2.Pos = i
		auxV1.Pos = i
		auxV2.Pos = i

		atLeastOne := make([]int, class_count)

		for j := 0; j < class_count-1; j++ {

			posV1.CId = base.CountableId{base.ClassType, j}
			posV2.CId = base.CountableId{base.ClassType, j + 1}
			atLeastOne[j] = getPosId(posV1)

			auxV1.CId = base.CountableId{base.ExactlyOne, j}
			auxV2.CId = base.CountableId{base.ExactlyOne, j + 1}

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

	var pV1, pVn base.PosVar

	pV1.Pos = 0
	pVn.Pos = size - 1

	for i := 0; i < class_count-1; i++ {

		pV1.CId = base.CountableId{base.ExactlyOne, i}
		pVn.CId = base.CountableId{base.ExactlyOne, i}

		clauses = append(clauses, clause{"sym", []int{getPosId(pV1), -getPosId(pVn)}})
	}

	pV1.CId = base.CountableId{base.ClassType, class_count - 1}
	pVn.CId = base.CountableId{base.ClassType, class_count - 1}
	pVn2 := base.PosVar{base.CountableId{base.ExactlyOne, class_count - 2}, size - 1}

	clauses = append(clauses, clause{"sym", []int{getPosId(pV1), -getPosId(pVn), -getPosId(pVn2)}})

	return
}

func createRedundant1(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var pV1, pV2 base.PosVar

	pV1.CId = c.CId
	pV2.CId = c.CId

	q := c.Window
	u := c.Capacity

	if u == 1 {
		for i := 0; i < size; i++ {

			pV1.Pos = i

			for j := i + 1; j < i+q && j < size; j++ {
				pV2.Pos = j
				cn := clause{"re1", []int{-getPosId(pV1), -getPosId(pV2)}}
				clauses = append(clauses, cn)
			}
		}
	}

	return
}

func createRedundant2(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	q := c.Window
	u := c.Capacity

	if u == 2 {

		var pV1, pV2, pV3 base.PosVar

		pV1.CId = c.CId
		pV2.CId = c.CId
		pV3.CId = c.CId

		for i := 0; i < size; i++ {

			pV1.Pos = i

			for j := i + 1; j < i+q && j < size; j++ {

				pV2.Pos = j

				for k := j + 1; k < i+q && k < size; k++ {

					pV3.Pos = k

					cn := clause{"re2", []int{-getPosId(pV1), -getPosId(pV2), -getPosId(pV3)}}
					clauses = append(clauses, cn)
				}
			}
		}
	}

	return
}

// only create these with options (1. definition of optimization statement)
func createOptPositions(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	var cV1, cV2 base.CountVar
	var optV base.PosVar

	cV1.CId = c.CId
	cV2.CId = c.CId
	optV.CId = c.CId
	optV.CId.Typ = base.OptimizationType

	q := c.Window
	u := c.Capacity

	if *sbd {
		// needed because avoid zero column, now extra work for sbd
		cV1.Pos = q - 1
		optV.Pos = q - 1
		cV1.Count = u + 1
		cn := clause{"op1", []int{getPosId(optV), -getCountId(cV1)}}
		clauses = append(clauses, cn)

	}

	for i := q; i < size; i++ {

		cV1.Pos = i - q
		optV.Pos = i
		cV2.Pos = i

		for j := c.Lower[i-q]; j < c.Upper[i-q]; j++ {
			cV1.Count = j
			cV2.Count = j + u
			if j+u < c.Upper[i] {
				cn := clause{"op0", []int{getPosId(optV), getCountId(cV1), -getCountId(cV2)}}
				clauses = append(clauses, cn)
			}
		}
	}

	return
}

func createOptCounter(c base.Countable) (clauses []clause) {

	clauses = make([]clause, 0)

	{ // set upper and lower bound for counters
		c.Lower = make([]int, size)
		c.Upper = make([]int, size)

		h := c.Demand

		for i := 0; i < size; i++ {
			c.Lower[i] = 0
			c.Upper[i] = h
			if h <= c.Demand {
				h++
			}
		}
	}

	pV := base.PosVar{c.CId, 0}
	cV2 := base.CountVar{c.CId, 0, 1}

	for i := *opt - 1; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j <= c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j
			c1 := clause{"op1", []int{-getCountId(cV1), getCountId(cV2)}}
			c3 := clause{"op3", []int{getPosId(pV), getCountId(cV1), -getCountId(cV2)}}
			clauses = append(clauses, c1, c3)
		}
	}

	for i := *opt - 1; i < size-1; i++ {

		cV1 := base.CountVar{c.CId, i, -1}
		cV2.Pos = i + 1
		pV.Pos = i + 1

		for j := c.Lower[i]; j < c.Upper[i]; j++ {
			cV1.Count = j
			cV2.Count = j + 1
			c2 := clause{"op2", []int{getCountId(cV1), -getCountId(cV2)}}
			c4 := clause{"op4", []int{-getPosId(pV), -getCountId(cV1), getCountId(cV2)}}
			clauses = append(clauses, c2, c4)
		}
	}

	return
}
