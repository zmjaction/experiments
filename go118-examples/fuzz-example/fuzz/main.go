package main

import (
	"fmt"
	"fuzz/utils"
)

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := utils.Reverse(input)
	doubleRev := utils.Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
