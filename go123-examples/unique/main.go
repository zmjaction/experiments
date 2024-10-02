package main

import (
	"fmt"
	"structs"
	"unique"
)

type Studuent struct {
	No   int
	Name string
	_    structs.HostLayout
}

func main() {
	h1 := unique.Make("a")
	h2 := unique.Make("b")
	if h1 == h2 {
		fmt.Println("a==b ?", true)
	}
	fmt.Println("h1 value", h1.Value())

	sa := unique.Make(Studuent{No: 1, Name: "a"})
	sb := unique.Make(Studuent{No: 2, Name: "b"})
	if sa == sb {
		fmt.Println("sa == sb")
	}
}
