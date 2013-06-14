package pbo

import (
	"../base"
	"fmt"
)

//func CreatePBOModelNative(N int, options, classes []base.Countable, class2option [][]bool) bool {
//} 

func CreatePBOModel(N int, options, classes []base.Countable, class2option [][]bool) bool {

	O := len(options)
	M := len(classes)

	ncons := 0

	// exactly one car per slot
	ncons += N
	// channel cars and options
	ncons += N * M * O
	// demand on options as boolean sums too
	ncons += O
	// capacity constraint as boolean sums

	for opt := 0; opt < O; opt++ {
		ncons += N - options[opt].Window + 1
	}
	// demand on classes
	ncons += M

	fmt.Printf("* #variable= %v #constraint= %v\n", (N * (M + O)), ncons)

	// exactly one car per slot
	for pos := 0; pos < N; pos++ {
		for cla := 0; cla < M; cla++ {
			fmt.Printf(" +1 x%v", 1+pos*M+cla)
		}
		fmt.Println(" = 1 ;")
	}

	// channel cars and options
	for car := 0; car < N; car++ {
		for cla := 0; cla < M; cla++ {
			for opt := 0; opt < O; opt++ {
				if class2option[cla][opt] {
					fmt.Printf(" -1 x%v +1 x%v >= 0 ;\n", (1 + car*M + cla), (N*M + 1 + car*O + opt))
				} else {
					fmt.Printf(" +1 x%v +1 x%v <= 1 ;\n", (1 + car*M + cla), (N*M + 1 + car*O + opt))
				}
			}
		}
	}

	for opt := 0; opt < O; opt++ {
		// demand on options as boolean sums too
		for car := 0; car < N; car++ {
			fmt.Printf(" +1 x%v", (N*M + 1 + car*O + opt))
		}
		fmt.Printf(" = %v ;\n", options[opt].Demand)

		// capacity constraint as boolean sums
		p := options[opt].Capacity
		q := options[opt].Window
		n := N - q + 1

		for i := 0; i < n; i++ {
			for j := 0; j < q; j++ {
				fmt.Printf(" +1 x%v", (N*M + 1 + (i+j)*O + opt))
			}
			fmt.Printf(" <= %v ;\n", p)
		}
	}

	for cla := 0; cla < M; cla++ {
		// demand on classes
		for car := 0; car < N; car++ {
			fmt.Printf(" +1 x%v", (1 + car*M + cla))
		}
		fmt.Printf(" = %v ;\n", classes[cla].Demand)
	}
	return true
}
