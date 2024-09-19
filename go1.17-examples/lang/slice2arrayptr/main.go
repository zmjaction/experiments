//go:build linux

package main

import "fmt"

func array2slice() {
	a := [3]int{11, 12, 13}
	b := a[:]
	b[1] += 10
	fmt.Printf("%v\n", a)
	fmt.Printf("%T\n", a)
}

func slice2arrayptr() {
	var b = []int{11, 12, 13}
	var p = (*[3]int)(b)
	p[1] = p[1] + 10
	fmt.Printf("%v\n", b)
	fmt.Printf("%T\n", p)
}

func main() {
	array2slice()
	slice2arrayptr()
}
