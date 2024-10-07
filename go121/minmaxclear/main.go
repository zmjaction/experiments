package main

import "fmt"

func main() {
	var x, y = 5, 6
	i := min(x, y)
	fmt.Println(i)
	fmt.Println(max("abc", "hello", "golang")) // hello
	//var f float64 = 5.6
	//fmt.Printf("%T\n", max(x, y, f))    // invalid argument: mismatched types int (previous argument) and float64 (type of f)
	//fmt.Printf("%T\n", max(x, y, 10.1)) // (untyped float constant) truncated to int
	var sl = []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("before clear, sl=%v, len(sl)=%d, cap(sl)=%d\n", sl, len(sl), cap(sl))
	clear(sl)
	fmt.Printf("after clear, sl=%v, len(sl)=%d, cap(sl)=%d\n", sl, len(sl), cap(sl))

	var m = map[string]int{
		"li":   13,
		"zhao": 14,
		"wang": 15,
	}
	fmt.Printf("before clear, m=%v, len(m)=%d\n", m, len(m))
	clear(m)
	fmt.Printf("after clear, m=%v, len(m)=%d\n", m, len(m))
}
