package main

import (
	"fmt"
)

func main() {
	done := make(chan bool)

	// 遍历 slice 中的所有元素，分别开协程对其进行一段逻辑操作
	// 这里，用打印元素来代表一段逻辑
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}

	// 等所有协程执行完
	for _ = range values {
		<-done
	}
}
