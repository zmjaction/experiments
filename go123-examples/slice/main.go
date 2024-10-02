package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

func slice() {
	arr := []int{1, 2, 3}
	brr := slices.Repeat(arr, 3) // [1 2 3 1 2 3 1 2 3]
	fmt.Println(brr)
	slices.Sort(brr)
	fmt.Println(brr) // [1 1 1 2 2 2 3 3 3]
	slices.Reverse(brr)
	fmt.Println(brr) // [3 3 3 2 2 2 1 1 1]
}

func All() {
	names := []string{"Alice", "Bob", "Vera"}
	for i, v := range slices.All(names) {
		fmt.Println(i, ":", v)
	}

}

// map按key排序
func sortMap() {
	m := map[string]struct{}{"赵六": {}, "张三": {}, "王五": {}, "李四": {}}
	for _, key := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%s\t", key)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
}

func seq() {
	m := map[string]struct{}{"赵六": {}, "张三": {}, "王五": {}, "李四": {}}
	s1 := maps.Keys(m) // 不保证每次的顺序都一样
	// Seq函数: type Seq[V any] func(yield func(V) bool)
	// 可以通过 for range 直接遍历这种函数
	for key := range s1 {
		fmt.Println(key)
	}
	// 也可以借助于Pull() 和next() 遍历seq
	next, stop := iter.Pull(s1)
	defer stop()
	for {
		key, valid := next()
		if valid {
			fmt.Println(key)
		} else {
			break
		}
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
}

func seq2() {
	m := map[string]string{"赵六": "1", "张三": "2", "王五": "3", "李四": "4"}
	s2 := maps.All(m)
	next, stop := iter.Pull2(s2)
	defer stop()
	for {
		key, value, valid := next()
		if valid {
			fmt.Println(key, value)
		} else {
			break
		}
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
}

func main() {
	slice()
	sortMap()
	seq()
	seq2()
}
