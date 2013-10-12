package sorters

//import (
//
//)
//
//type clause struct {
//	desc     string
//	literals []int
//}
//
//type IdGen struct {
//	nextId        int
//	mapping map[Var2]int
//}
//
//type Var2 struct {
//	V1 int
//	V2 int
//}
//
//var clauseTypes []string
//
//var gen IdGen
//
//func TestGenerateSAT(t *testing.T) {
//    size := 8
//    k := 3
//
//    sorter := createCardinalityNetwork(size,k,AtMost,OddEven)
//}
//
//func getId(v Var2) (id int) {
//	id, b := gen.mapping[v]
//
//	if !b {
//		gen.nextId++
//		id = gen.Id
//		gen.mapping[v] = id
//	}
//
//	return id
//}
//
//
//func solve(clauses []clause) {
//
//    printClausesDIMACS(clauses)
//
//}
//
//func printClausesDIMACS(clauses []clause) {
//
//	fmt.Printf("p cnf %v %v\n", len(gen.mapping), len(clauses))
//
//	for _, c := range clauses {
//		for _, l := range c.literals {
//			fmt.Printf("%v ", l)
//		}
//		fmt.Printf("0\n")
//	}
//}
//
//func generateSymbolTable() []string {
//
//	symbolTable := make([]string, len(gen.mapping)+1)
//
//	for key, cnfId := range gen.mapping {
//		s := "var" +"("
//		s += strconv.Itoa(key.V1)
//		s += ","
//		s += strconv.Itoa(key.V2)
//		s += ")"
//		symbolTable[cnfId] = s
//	}
//
//	return symbolTable
//}
//
//func printSymbolTable(symbolTable []string, filename string) {
//
//	symbolFile, err := os.Create(filename)
//	if err != nil {
//		panic(err)
//	}
//	// close fo on exit and check for its returned error
//	defer func() {
//		if err := symbolFile.Close(); err != nil {
//			panic(err)
//		}
//	}()
//
//	// make a write buffer
//	w := bufio.NewWriter(symbolFile)
//
//	for i, s := range symbolTable {
//		// write a chunk
//		if _, err := w.Write([]byte(fmt.Sprintln(i, "\t:", s))); err != nil {
//			panic(err)
//		}
//	}
//
//	if err = w.Flush(); err != nil {
//		panic(err)
//	}
//
//}
//
//func printDebug(symbolTable []string, clauses []clause) {
//
//	// first print symbol table into file
//	fmt.Println("c pred2(V1,V2).")
//
//	for i, s := range symbolTable {
//		fmt.Println("c", i, "\t:", s)
//	}
//
//	stat := make(map[string]int, 10)
//
//	for _, c := range clauses {
//
//		count, ok := stat[c.desc]
//		if !ok {
//			stat[c.desc] = 1
//		} else {
//			stat[c.desc] = count + 1
//		}
//
//		fmt.Printf("c %s ", c.desc)
//		first := true
//		for _, l := range c.literals {
//			if !first {
//				fmt.Printf(",")
//			}
//			first = false
//			if l < 0 {
//				fmt.Printf("-%s", symbolTable[-l])
//
//			} else {
//				fmt.Printf("+%s", symbolTable[l])
//			}
//		}
//		fmt.Println(".")
//	}
//
//	for _, key := range clauseTypes {
//		fmt.Printf("c %v\t: %v\t%.1f \n", key, stat[key], 100*float64(stat[key])/float64(len(clauses)))
//	}
//	fmt.Printf("c %v\t: %v\t%.1f \n", "tot", len(clauses), 100.0)
//}
