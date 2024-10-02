package main

import "fmt"

// type MySlice[T int | float32] []T
//
//	func (s MySlice[T]) Sum() T {
//		var sum T
//		for _, value := range s {
//			sum += value
//		}
//		return sum
//	}
//
// type MySlice2 []int
//
//	func (s MySlice2) Sum2() int {
//		var sum int
//		for _, value := range s {
//			sum += value
//		}
//		return sum
//	}

type MyInt2 int

type MyInt interface {
	int | int8 | int16 | int32 | int64
}

func GetMaxNum[T MyInt](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(GetMaxNum[int](10, 20))
	//fmt.Println(GetMaxNum[MyInt2](10, 20))
}
