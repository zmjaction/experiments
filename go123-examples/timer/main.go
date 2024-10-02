package main

import (
	"fmt"
	"time"
)

func main() {
	timeout := 50 * time.Millisecond
	t := time.NewTimer(timeout)

	time.Sleep(100 * time.Millisecond)

	start := time.Now()
	t.Reset(timeout)
	<-t.C

	fmt.Printf("煎鱼已经消失：%dms\n", time.Since(start).Milliseconds())
}
