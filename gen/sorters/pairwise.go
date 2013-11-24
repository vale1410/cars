package sorters

//func buildTriangleBitonic(newId *int, array []int, comparators *[]Comparator, lo int, hi int) {
//	if (hi - lo) >= 1 {
//		//fmt.Println("compare", lo, hi)
//		compareAndSwap(newId, array, comparators, lo, hi)
//		buildTriangleBitonic(newId, array, comparators, lo+1, hi-1)
//	}
//}
//
//func pairwiseMerge(newId *int, array []int, comparators *[]Comparator, lo int, hi int) {
//	if (hi - lo) >= 1 {
//		//fmt.Println("waterfall", lo, hi)
//		mid := lo + ((hi - lo) / 2)
//		for i := 0; i <= mid-lo; i++ {
//			//fmt.Println("compare", lo+i, mid+i+1)
//			compareAndSwap(newId, array, comparators, lo+i, mid+i+1)
//		}
//		waterfallBitonic(newId, array, comparators, lo, mid)
//		waterfallBitonic(newId, array, comparators, mid+1, hi)
//	}
//}
//
//func pairwiseSplit(newId *int, array []int, comparators *[]Comparator, lo int, hi int, spread) {
//    mid := lo + ((hi - lo) / 2)
//    for i := 0; i*spead <= mid-lo; i++ {
//    	//fmt.Println("compare", lo+i, mid+i+1)
//    	compareAndSwap(newId, array, comparators, lo+i*spread, mid+i*speard)
//    }
//
//func pairwiseSort(newId *int, array []int, comparators *[]Comparator, lo int, hi int, spread int) {
//	if (hi - lo) >= 1 {
//		//fmt.Println("triangle", lo, hi)
//		pairwiseSplit(newId, array, comparators, lo, mid,spread) // spread goes up by *2
//		pairwiseSort(newId, array, comparators, lo, mid-1,spread*2)
//		pairwiseSort(newId, array, comparators, lo+1, mid,spread*2)
//
//		pairwiseMerge(newId, array, comparators, lo, mid,spread)
//	}
//}
