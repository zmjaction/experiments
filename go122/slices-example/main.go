package main

import (
	"fmt"
	"net/http"
	"slices"
)

func main() {
	//s1 := []int{11, 12, 13, 14}
	//s2 := slices.Delete(s1, 1, 3)
	//fmt.Println("s1:", s1)
	//fmt.Println("s2:", s2)
	//s1 := []string{"Go slices", "Go maps"}
	//s2 := []string{"Go strings", "Go strconv"}
	//s4 := slices.Concat(s1, s2)
	//fmt.Printf("cap: %d, len: %d\n", cap(s4), len(s4))
	//fmt.Println(s4)
	//s1DeleteFunc := []int{1, 2, 3, 4, 5}
	//s2DeleteFunc := slices.DeleteFunc(s1DeleteFunc, func(e int) bool {
	//	return e%2 == 0
	//})
	//fmt.Printf("cap: %d, len: %d\n", cap(s1DeleteFunc), len(s1DeleteFunc))
	//fmt.Println(s1DeleteFunc)
	//fmt.Println(s2DeleteFunc)
	//s1 := []string{"hello world", "Hello world", "hello World"}
	////s1 := []string{"Gopher", "MingYong Chen", "mingyong chen"}
	//
	//s2 := slices.CompactFunc(s1, func(a, b string) bool {
	//	return strings.ToLower(a) == strings.ToLower(b)
	//})
	//fmt.Printf("%#v\n", s1)
	//fmt.Printf("%#v\n", s2)
	//s1 := []int{1, 6, 7, 4, 5}
	//s2 := slices.Replace(s1, 1, 3, 2)
	//fmt.Println(s1)
	//fmt.Println(s2)
	s1 := []string{"hello", "world"}
	s2 := slices.Insert(s1, 3)
	fmt.Println(s2)
	//Hello()
}

func Hello() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != "GET" {
	//		fmt.Fprintf(w, "warn: 只支持GET方法")
	//	} else {
	//		fmt.Fprintf(w, "你好 "+r.PathValue("name"))
	//	}
	//})
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "你好 "+r.PathValue("name"))
	})
}
