package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"os"
	"slices"
	"strings"
)

//func Backward[E any](s []E) func(func(int, E) bool) {
//	return func(yield func(int, E) bool) {
//		for i := len(s) - 1; i >= 0; i-- {
//			if !yield(i, s[i]) {
//				return
//			}
//		}
//		return
//	}
//}
//
////func main() {
////	sl := []string{"hello", "world", "golang"}
////	Backward(sl)(func(i int, s string) bool {
////		fmt.Printf("%d : %s\n", i, s)
////		return true
////	})
////}
//
//func main() {
//	sl := []string{"hello", "world", "golang"}
//	for i, s := range Backward(sl) {
//		fmt.Printf("%d : %s\n", i, s)
//	}
//}

// Lines 对文件里的文本行进行遍历
func Lines(path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
}

// Entries 对json对象的属性进行遍历
func Entries(object string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		dec := json.NewDecoder(strings.NewReader(object))
		var kvs map[string]string
		if err := dec.Decode(&kvs); err != nil {
			panic(err)
		}
		for k, v := range kvs {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Head 对HTTP head 请求回应的header进行遍历
func Head(url string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		r, err := http.Head(url)
		if err != nil {
			panic(err)
		}
		r.Body.Close()
		if r.StatusCode != http.StatusOK {
			panic(fmt.Errorf("invalid status %v", r.StatusCode))
		}
		for k, v := range r.Header {
			if !yield(k, strings.Join(v, ",")) {
				return
			}
		}
	}
}

/*Backward[E any] Backward 是一个泛型函数，反向遍历切片
 * @Description: 接受一个切片（s），切片的元素类型为泛型类型E（可以是任何类型）s
 * @param s
 * @return func(yield func(int, E) bool)
 */
func Backward[E any](s []E) func(yield func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

func main() {
	s := []string{"a", "b", "c"}
	for i := len(s) - 1; i >= 0; i-- {
		fmt.Println(i, s[i])

	}
	//iterFunc := slices.Backward(s)
	//callBack := func(i, int, str string) bool {
	//	fmt.Println(i, str)
	//	return true
	//}
	//iterFunc(callBack)
	for i, v := range Backward(s) {
		fmt.Println(i, v)

	}
	fmt.Println(strings.Repeat("-", 50))
	for k, v := range slices.Backward(s) { // 使用迭代器和手写循环本质上是一样的
		fmt.Println(k, v)
	}
	for line := range Lines("a.txt") {
		fmt.Println(line)
	}
	for k, v := range Entries(`{"name":"go", "version":"1.23.0"}`) {
		fmt.Printf("%v=%v\n", k, v)
	}

	for k, v := range Head("https://golang.google.cn/") {
		if len(v) < 100 {
			fmt.Printf("%v=%v\n", k, v)
		}
	}
}
