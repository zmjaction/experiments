# Go版本

Go并不是一成不变的编程语言。最初的Go1.0发布以来，Go语言习惯的模式已经发生了重大的变化

## go 1.23

*在 `Go 1.22` 版本发布 **6** 个月之后，`Go 1.23` ，2024年8月13日* GO1.23发布，于北京时间 **2024** 年 **8** 月 **14** 日凌晨 **1:03** 正式发布

快速安装：

```go
$ go install golang.org/dl/go1.23.0@latest
$ go1.23.0 download  
Downloaded   0.0% (   16384 / 71607288 bytes) ...
Downloaded  11.0% ( 7880688 / 71607288 bytes) ...
Downloaded  41.7% (29835136 / 71607288 bytes) ...
Downloaded  61.6% (44121808 / 71607288 bytes) ...
Downloaded  77.7% (55639936 / 71607288 bytes) ...
Downloaded  95.1% (68107968 / 71607288 bytes) ...
Downloaded 100.0% (71607288 / 71607288 bytes)
Unpacking /Users/chenmingyong/sdk/go1.23.0/go1.23.0.darwin-arm64.tar.gz ...
Success. You may now run 'go1.23.0'
$ go1.23.0 version
go version go1.23.0 darwin/arm64
```

### 一、语言特性更新

- 新的迭代器语法：在 "for-range" 循环中，现在可以使用迭代器函数作为 range 表达式，如 func (func (K) bool)。这[支持用户自定义任意序列的迭代器](https://www.oschina.net/action/GoToLink?url=http%3A%2F%2Fmp.weixin.qq.com%2Fs%3F__biz%3DMzIyNzM0MDk0Mg%3D%3D%26mid%3D2247497326%26idx%3D1%26sn%3D65618d9554bde8f43a19ca4c0be25492%26chksm%3De860118fdf17989908ac10d01debc852c9356c3034a9e2068fcfe9c88999f7f921c4ab4930c4%26scene%3D142%23wechat_redirect)。标准库的 slices 和 maps 包也添加了支持迭代器的新功能。
- 泛型类型别名预览: Go 1.23 包含了对泛型类型别名的预览支持。

#### 1、函数迭代器

在 `Go 1.23` 中，**迭代器** 实际上是指符合以下三种函数签名之一的函数：

```go
func(yield func() bool)

func(yield func(V) bool)

func(yield func(K, V) bool)
```

如果一个函数或方法返回的值符合上述形式之一，那么该返回值就可以被称为 **迭代器**。

代码示例：

```go
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
```

#### Range Over Function Types （对函数类型遍历）

在 `Go 1.23` 版本中，`for-range` 循环中的范围表达式得到了改进。此前，范围表达式仅支持 `array`（数组）、`slice`（切片） 和 `map`（映射） 等类型，而从 `Go 1.23` 开始，新增了对函数类型的支持。不过，函数类型必须是前面所提到的三种类型之一，也就是函数需要实现迭代器。

```go
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
  s := string{"a", "b"}
  for i, v := range Backward(s) {
		fmt.Println(i, v)

	}
}

```

①返回函数接收yield函数作为参数，接收一个int类型和一个E类型

②它使用一个for循环从切片的最后一个元素开始向前遍历切片。对于每个元素，它调用yield函数，并将当前元素的索引和元素本身作为参数传递给yield函数

③如果`yield`函数返回`false`，则内部函数会立即返回，停止遍历切片。如果循环正常结束（即`yield`函数始终返回`true`），则内部函数执行完毕。

#### iter包

为了简化迭代器的使用，`Go 1.23` 版本新增了一个 `iter` 包，该包定义了两种迭代器类型，分别是 `Seq` 和 `Seq2`，用于处理不同的迭代场景。

```go
package iter

type Seq[V any] func(yield func(V) bool)

type Seq2[K, V any] func(yield func(K, V) bool)
```

**`Seq` 和 `Seq2` 的区别：**

- **`Seq[V any]`**
  `Seq` 是一个泛型类型的函数，接收一个 `yield` 函数作为参数。它推出单个元素，例如切片的索引或映射中的键。`yield` 函数返回 `bool`，决定是否继续迭代。

  使用场景：可以用于返回一个单值的迭代，比如切片中的索引或值，或映射中的键或值。

- **`Seq2[K, V any]`**
  `Seq2` 是一个泛型类型的函数，接收一个 `yield` 函数，推送一对元素，例如切片中的索引和值，或者映射中的键值对。`yield` 函数同样返回 `bool`，以决定是否继续迭代。

  使用场景：当需要同时返回两个值（如键和值）时使用

Seq代码示例：

```go
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
for line := range Lines("a.txt") {
		fmt.Println(line)
}
```

Seq2代码示例：

```go
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

// 可以通过 for-range 循环直接接收被推送的值
for k, v := range Entries(`{"name":"go", "version":"1.23.0"}`) {
		fmt.Printf("%v=%v\n", k, v)
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
```



### 2、标准库的 slices 和 maps 包也添加了支持迭代器的新功能

前面说过，函数迭代器转正了。标准库中有一些包立即就提供了一些便利的、可以与函数迭代器一起使用的函数，以slices、maps两个后加入Go标准库的泛型容器包为主。

`slices` 包新增的函数：

- `All([]E) iter.Seq2[int, E]`
- `Values([]E) iter.Seq[E]`
- `Collect(iter.Seq[E]) []E`
- `AppendSeq([]E, iter.Seq[E]) []E`
- `Backward([]E) iter.Seq2[int, E]`
- `Sorted(iter.Seq[E]) []E`
- `SortedFunc(iter.Seq[E], func(E, E) int) []E`
- `SortedStableFunc(iter.Seq[E], func(E, E) int) []E`
- `Repeat([]E, int) []E`
- `Chunk([]E, int) iter.Seq([]E)`s

1. **`All`**：`func All(t any) (s()t) iter.iterator(t)`，该函数返回一个迭代器，用于遍历切片的索引和值。这对于需要同时访问切片元素及其索引的操作非常有用。
2. **`Values`**：`func Values(t any) (s()t) iter.iterator(t)`，此函数返回一个仅遍历切片元素的迭代器，忽略索引。当只关心切片中的元素值时，可以使用这个函数。
3. **`Backward`**：`func Backward(t any) (s()t) iter.iterator(t)`，该函数返回一个反向遍历切片的迭代器，按照从后向前的顺序访问切片元素，适用于需要反向处理切片数据的场景。
4. **`Collect`**：`func Collect(t any) (it iter.iterator(t)) ()t`，用于将迭代器中的值收集到一个新的切片中，方便将迭代操作的结果转换为切片形式。
5. **`AppendSeq`**：`func AppendSeq(t any) (dst()t, it iter.iterator(t)) ()t`，可以将一个迭代器中的值追加到一个已有的切片中，实现切片的扩展。
6. **`Sorted`**：`func Sorted(t constraints.Ordered) (it iter.iterator(t)) ()t`，该函数从迭代器中收集值到一个新切片，并对其进行排序。它要求切片元素类型必须满足`constraints.Ordered`约束，即元素类型必须支持比较操作。
7. **`SortedFunc`**：`func SortedFunc(t any) (it iter.iterator(t), less func(a, b t) bool) ()t`，与`Sorted`类似，但允许用户自定义比较函数`less`来对元素进行排序，提供了更灵活的排序方式。

`maps` 包新增的函数：

- `All(map[K]V) iter.Seq2[K, V]`
- `Keys(map[K]V) iter.Seq[K]`
- `Values(map[K]V) iter.Seq[V]`
- `Collect(iter.Seq2[K, V]) map[K, V]`
- `Insert(map[K, V], iter.Seq2[K, V])`

### 二、工具链改进

- `Go telemetry` 遥测系统：允许 `Go` 的工具链（编译器、调试器等工具）收集使用和故障统计数据。这些数据的收集是为了帮助开发团队了解 `Go` 工具链的使用情况和运行状态，从而对工具链进行改进和优化。
- `Go` 命令：新增了一些便利的功能。例如运行 `go env -changed` 可以更容易地查看哪些设置的有效值与默认值不同，而 `go mod tidy -diff` 可以在不修改 `go.mod` 和 `go.sum` 文件的情况下，帮助你查看需要进行的更改。
- `Go vet 子命令`：现在能够检测代码中使用的某些特性或函数是否对于目标 `Go` 版本来说太新了。

### 三、标准库更新

- 优化了 `time.Timer` 与 `time.Ticker` 两个定时器的实现。

- 标准库中新增了总共三个包：`iter`、`structs` 和 `unique`。

   `iter`：提供了与序列上的迭代器相关的基本定义和操作。

  `structs`：定义了标记类型，用于修改结构体的属性。

  `unique`：提供了规范化（“interning”）可比较值的工具。

  https://before80.github.io/go_docs/goBlog/2024/NewUniquePackage/

- GODEBUG 设置：支持在 go.mod 和 go.work 文件中使用新的 godebug 指令

### 四、参考

Go 1.23 https://go.dev/blog/go1.23

Go 1.23 的发布说明文档 https://go.dev/doc/go1.23 

Go 1.23版本里程碑 https://github.com/golang/go/milestone/212 

Next Release Notes Deaft https://tip.golang.org/doc/next 

Go Release Dashboard https://dev.golang.org/release

 Go 1.23新特性前瞻 https://mp.weixin.qq.com/s/c7UuQetStkA7Tw2DLfMjvA

Go 1.23 unique库 https://mp.weixin.qq.com/s/NDqeknAm7q77siHm0Jbcxg

unique 背景 https://github.com/golang/go/issues/62483

*https://github.com/go4org/intern*

time.Reset 过期时间问题 https://mp.weixin.qq.com/s/NijdOmdfKGLJowhbe9yqPg

time.After 泄露问题 https://mp.weixin.qq.com/s/Qcpj7TqMeOwCs--59kD3Kw

//go:linkname特性 https://segmentfault.com/a/1190000045164130

slice https://pkg.go.dev/slices@master

go1.23迭代器 https://chenmingyong.cn/posts/go1.23-iterator

## go 1.22

2024年 2月6日，Go 官方发布了最新的 Go1.22 版本

快速安装：

```go
go install golang.org/dl/go1.23.0@latest
```

最新的Go版本1.22比Go 1.21晚了6个月。 它的大部分变化都在工具链、运行时和库的实现中。 与往常一样，该版本保持了Go 1对兼容性

#### 一、语言特性更新

Go 1.22对“for”循环做了两个更改。

##### 1、for循环声明的变量只创建一次，并在每次迭代中更新。在Go 1.22中，循环的每次迭代都会创建新变量，以避免意外的共享错误

先来看一段代码：

```go
package main

import (
 "fmt"
)

func main() {
    done := make(chan bool)

    // 遍历 slice 中的所有元素，分别开协程对其进行一段逻辑操作
    // 这里，用打印元素来代表一段逻辑
    values := []string{"a", "b", "c"}
    for i, v := range values {
        go func() {
            fmt.Printf("&p\n", &v) // 在Go122之前
            fmt.Println(i, v)
            done <- true
        }()
    }

    // 等所有协程执行完
    for _ = range values {
        <-done
    }
}
```

![image-20241004100822628](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20241004100822628.png)

![image-20241004101152693](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20241004101152693.png)

数据竞争问题：

我们可以将上面的 for range 语句做一个等价转换，这样可以帮助你理解 for range 的工作原理。等价转换后的结果是这样的：

```go
func main() {
    done := make(chan bool)
  
    values := []string{"a", "b", "c"}
  
    {
      i, v := 0, 0
      for i, v := range values {
        go func() {
          fmt.Printf("&p\n", &v)
          fmt.Println(i, v)
          done <- true
        }()
      }
    }


    // 等所有协程执行完
    for _ = range values {
        <-done
    }
}
```

通过等价转换后的代码，我们可以清晰地看到循环变量 i 和 v 在每次迭代时的重用。而 Goroutine 执行的闭包函数引用了它的外层包裹函数中的变量 i、v，这样，变量 i、v 在主 Goroutine 和新启动的 Goroutine 之间实现了共享，而 i, v 值在整个循环过程中是重用的，仅有一份。

##### 2.range 关键字支持整型数据

在Go 1.22版本中，for range后面的range表达式除了支持传统的像数组、切片、map、channel等表达式外，**还支持放置整型表达式**

示例代码：

```go
package main

import "fmt"

func main() {
	for i := range 10 {
		fmt.Println(i)
	}
}
// Notice 下面的代码只能在go1.22版本运行，以下的版本会报错
0
1
2
3
4
5
6
7
8
9
```

如果n <= 0，则循环不运行任何迭代。

这个新语法特性，可以理解为是一种“语法糖”，是下面等价代码的“语法糖”：

```go
for i:=0;i<5;i++ {
  
}
```

迭代总是从0开始，似乎限制了该语法糖的使用范围

##### 3.函数迭代器的

Go增加了函数迭代器（iterator），函数迭代器以试验特性提供，通过GOEXPERIMENT=rangefunc可以体验该功能特性

示例代码：

```go
func Backward(s []E) func(func(int, E) bool) {
    return func(yield func(int, E) bool) {
        for i := len(s)-1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
        return
    }
}
```

### 二、编译器、工具链

1、Go 1.22版本在编译上优化PGO(profile-guided optimization)

2、在工具链方面：

> go work支持vendor

在Go 1.22版本中，通过go work vendor可以将workspace中的依赖放到vendor⽬录下，同时在构建时，如果workspace下有vendor⽬录，那么默认的构建是go build -mod=vendor，即基于vendor的构建。

> 改进go test -cover的输出

对于没有自己的测试文件的包，go test -cover在go 1.22版本之前会输出：

```
? mymod/mypack [no test files]
```

但在Go 1.22版本之后，会报告覆盖率为0.0%：

```
mymod/mypack coverage: 0.0% of statements
```

### 三、标准库

#### 1.math/rand/v2：标准库的第一个v2版本包

Go 1.22中新增了math/rand/v2包，这里之所以将它列为Go 1.22版本标准库的⼀次重要变化，是因为这是标准库第一次为某个包建⽴v2版本，按照Russ Cox的说法，这次math/rand/v2包的创建，算是为标准库中的其他可能的v2包“探探路”，找找落地路径。关于math/rand/v2包相对于原math/rand包的变化有很多，具体可以参考[issue 61716](https://go.dev/issue/61716)中的设计与讨论。

#### 2.slice

##### Concat：高效拼接切片

Concat函数接受一个不定参数slices，参数类型为切片，该函数用于将多个切片拼接到一个新的切片里并返回新切片。

在以前的 Go 版本中，有一个很常见的使用场景，如果我们想要拼接两个切片。必须要手写类似如下的代码：

```go
func main() {
 s1 := []string{"a", "b", "c"}
 s2 := []string{"1", "2", "3"}

 s3 := append(s1, s2...)
 fmt.Println(s3)
}
```

Concat 函数签名如下：

```go
func Concat[S ~[]E, E any](slices ...S) S
```

使用Concat函数，示例：

```go
import (
 "fmt"
 "slices"
)

func main() {
 s1 := []string{"hello"}
 s2 := []string{"a", "b"}
 s3 := []string{"1", "2"}
 resp := slices.Concat(s1, s2, s3)
 fmt.Println(resp)
 fmt.Printf("cap: %d, len: %d\n", cap(resp), len(resp))
}
```

其内部函数实现也比较简单。如下代码：

```go
// Concat returns a new slice concatenating the passed in slices.
func Concat[S ~[]E, E any](slices ...S) S {
 size := 0
 for _, s := range slices {
  size += len(s)
  if size < 0 {
   panic("len out of range")
  }
 }
 newslice := Grow[S](nil, size)
 for _, s := range slices {
  newslice = append(newslice, s...)
 }
 return newslice
}
```

`Concat` 函数的源码实现非常简洁，它在拼接切片之前先计算了新切片所需的长度，然后利用 `Grow` 函数初始化新切片。这样做的好处是避免了后续 `append` 操作中因为切片扩容而导致的内存重新分配和复制问题，使得函数更加高效。

##### Delete`、`DeleteFunc`、`Compact`、`CompactFunc` 和 `Replace 函数，零化处理

在 Go 1.22 版本中，对 Delete、DeleteFunc、Compact、CompactFunc 和 Replace 函数进行了更新。这些函数的共同点是接受一个给定的切片参数，记为 s1，并返回一个新切片，记为 s2。被移除的元素会在 s1 中被置为零值（被移除的元素 是指从 s1 中移除的指定元素，在s2 中不存在）。

**Delete 函数**
通过不同 Go 版本的代码示例来感受 Delete 函数 零化处理 的更新。
Go1.21版本的代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 3, 4, 5}
        s2 := slices.Delete(s1, 3, 5)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 3 4 5]
[1 2 3]
```

Go 1.22版本代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 3, 4, 5}
        s2 := slices.Delete(s1, 3, 5)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 3 0 0]
[1 2 3]
```

通过对比不同版本的代码运行结果可知，**被移除的元素** 在原切片里被置为了 **零值**。

DeleteFunc 函数

通过不同 `Go` 版本的代码示例来感受 `DeleteFunc` 函数 **零化处理** 的更新。

在Go 1.21版本的代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 3, 4, 5}
        s2 := slices.DeleteFunc(s1, func(e int) bool {
                return e%2 == 0
        })
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 3 5 4 5]
[1 3 5]
```

①在传入的函数`func(e int) bool`中，逻辑为`return e%2 == 0`，即当元素为偶数时，该函数返回`true`。`slices.DeleteFunc`函数会遍历原始切片`s1`，对于每个元素，都会调用这个条件函数。如果条件函数返回`true`，则该元素会被删除。

在Go1.22版本代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 3, 4, 5}
        s2 := slices.DeleteFunc(s1, func(e int) bool {
                return e%2 == 0
        })
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 3 5 0 0]
[1 3 5]
```

通过对比不同版本的代码运行结果可知，**被移除的元素** 在原切片里被置为了 **零值**。

Compact函数

通过不同 `Go` 版本的代码示例来感受 `Compact` 函数 **零化处理** 的更新。

在Go1.21版本示例代码：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 2, 3, 3, 4, 5}
        s2 := slices.Compact(s1)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 3 4 5 4 5]
[1 2 3 4 5]
```

在Go1.22版本示例代码：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 2, 2, 3, 3, 4, 5}
        s2 := slices.Compact(s1)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 3 4 5 0 0]
[1 2 3 4 5]
```

通过对比不同版本的代码运行结果可知，**被移除的元素** 在原切片里被置为了 **零值**。

CompactFunc 函数

通过不同 `Go` 版本的代码示例来感受 `CompactFunc` 函数 **零化处理** 的更新。

在Go1.21版本示例代码：

```go
package main

import (
        "fmt"
        "slices"
        "strings"
)

func main() {
        s1 := []string{"hello world", "Hello world", "hello World"}
        s2 := slices.CompactFunc(s1, func(a, b string) bool {
                return strings.ToLower(a) == strings.ToLower(b)
        })
        fmt.Printf("%#v\n", s1)
        fmt.Printf("%#v\n", s2)
}
// 代码运行结果
[]string{"hello world", "Hello world", "hello World"}
[]string{"hello world"}
```

在Go 1.22版本示例代码：

```go
package main

import (
        "fmt"
        "slices"
        "strings"
)

func main() {
        s1 := []string{"hello world", "Hello world", "hello World"}
        s2 := slices.CompactFunc(s1, func(a, b string) bool {
                return strings.ToLower(a) == strings.ToLower(b)
        })
        fmt.Printf("%#v\n", s1)
        fmt.Printf("%#v\n", s2)
}
// 代码运行结果
[]string{"hello world", "", ""}
[]string{"hello world"}
```

通过对比不同版本的代码运行结果可知，**被移除的元素** 在原切片里被置为了 **零值**。

Replace 函数

通过不同 `Go` 版本的代码示例来感受 `Replace` 函数 **零化处理** 的更新。

在Go1.21版本示例代码：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 6, 7, 4, 5}
        s2 := slices.Replace(s1, 1, 3, 2)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 4 5 5]
[1 2 4 5]
```

在Go1.22 版本代码中示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []int{1, 6, 7, 4, 5}
        s2 := slices.Replace(s1, 1, 3, 2)
        fmt.Println(s1)
        fmt.Println(s2)
}
// 代码运行结果
[1 2 4 5 0]
[1 2 4 5]
```

##### 越界插入优化

`Go 1.22` 版本对 `slices` 库的 `Insert` 函数进行了优化。在使用 `Insert` 函数时，若参数 `i` 超出切片的范围，总会触发 `panic`。而在 `Go 1.22` 版本之前，即使 `i` 越界了，在没有指定插入元素的情况下，该行为不会触发 `panic`。

在Go1.21版本代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []string{"hello", "world"}
        s2 := slices.Insert(s1, 3)
        fmt.Println(s2)
}
// 代码运行结果
[hello world]
```

在Go1.22版本代码示例：

```go
package main

import (
        "fmt"
        "slices"
)

func main() {
        s1 := []string{"hello", "world"}
        s2 := slices.Insert(s1, 3)
        fmt.Println(s2)
}
// 代码运行结果
panic: runtime error: slice bounds out of range [3:2]

goroutine 1 [running]:
slices.Insert[...]({0x14000092020?, 0x1400004c738?, 0x0?}, 0x60?, {0x0?, 0x1400009ef38?, 0x1003c4ad0?})
        /Users/zmj/go/go1.22/src/slices/slices.go:133 +0x434
main.main()
        /Users/zmj/workspace/experiments/go122/slices-example/main.go:38 +0x70
```

#### 3.增强http.ServeMux表达能力

```go
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
```

关于新版http.ServeMux的具体使用方法，其作者Jonathan Amsterdam（也是[log/slog](https://tonybai.com/2022/10/30/first-exploration-of-slog)的作者）在官博发表了一篇名为“[Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements)”的文章，大家可以详细参考。

### 四、参考

https://go.dev/doc/go1.22 Go 1.22 的发布说明文档

https://juejin.cn/post/7351336619594743848 go1.22版本

https://tonybai.com/2023/12/25/go-1-22-foresight/ go 1.22 前瞻性

https://colobu.com/2023/12/24/new-math-rand-in-Go/ Go标准库新的math/rand

https://mp.weixin.qq.com/s?__biz=MzUxMDI4MDc1NA==&mid=2247500406&idx=1&sn=d91ed868609a8a6217867bd2ba15cc36&scene=21#wechat_redirect Go1.22 新特性 math/rand/v2库

https://mp.weixin.qq.com/s?__biz=MzUxMDI4MDc1NA==&mid=2247500437&idx=1&sn=9d21f73b106e0a6a9c466f6db98c7163&scene=21#wechat_redirect Go 1.22 Slices 变更Concat、Delete、Insert函数

https://blog.csdn.net/weixin_44604586/article/details/136227324 Go 1.22slices库

https://tonybai.com/2024/02/18/some-changes-in-go-1-22/ Go1.22值得关注的变化 

## go 1.21

2023年 8月8日，Go 官方发布了最新的 Go1.22 版本

最新的Go版本1.21比Go 1.20晚了6个月。 它的大部分变化都在工具链、运行时和库的实现中。 与往常一样，该版本保持了Go 1对兼容性的承诺; 事实上，Go 1.21改进了这个承诺。 我们期望几乎所有的Go程序都能像以前一样编译和运行。

### 1、语言特性更新

#### min、max、clear

在Go 1.21版本中，builtin增加了三个预定义函数：min、max和clear。

顾名思义，min和max函数分别返回参数列表中的最小值和最大值，它们都是泛型函数，原型如下：

```
func min[T cmp.Ordered](x T, y ...T) T
func max[T cmp.Ordered](x T, y ...T) T
```

示例：

```go
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
}
```

我们看到：Go 1.21编译器报错，即便是untyped constant，如果类型不同，也会提醒你可能存在值精度的truncated。

max和min支持哪些类型呢？通过min和max原型中的类型参数(type parameter)可以看到，其约束类型(constraint)为cmp.Ordered，我们看一下该约束类型的定义：

```
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
        ~float32 | ~float64 |
        ~string
}
```

符合Ordered约束的上述这些类型以及衍生类型都可以使用min、max获取最小值和最大值。

新增一个clear函数:

```go
func clear[T ~[]Type | ~map[Type]Type1](t T)
```

示例：

```go
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
// 结果
before clear, sl=[1 2 3 4 5 6], len(sl)=6, cap(sl)=6
after clear, sl=[0 0 0 0 0 0], len(sl)=6, cap(sl)=6
before clear, m=map[li:13 wang:15 zhao:14], len(m)=3
after clear, m=map[], len(m)=0

```

- 针对slice，clear保持slice的长度和容量，但将所有slice内已存在的元素(len个)都置为元素类型的零值；
- 针对map，clear则是清空所有map的键值对，clear后，我们将得到一个empty map。

> 注：clear函数在清空map中的键值对时，并未释放掉这些键值所占用的内存。



https://go.dev/doc/go1.21 Go 1.21 的发布说明文档

https://tonybai.com/2023/08/20/some-changes-in-go-1-21/ Go1.21 值得关注的变化

## go 1.20



https://go.dev/doc/go1.20 Go 1.20 的发布说明文档

## go 1.19

Go官方团队在2022.06.11发布了Go 1.19 Beta 1版本，Go 1.19的正式release版本预计会在今年8月份发布。

### 一、语言特性更新

### 二、可移植

正式在linux上支持龙芯架构(GOOS=linux, GOARCH=loong64)

这一点不得不提，因为这一变化都是国内龙芯团队贡献的。不过目前龙芯支持的linux kernel版本最低也是5.19，意味着龙芯在老版本linux上还无法使用Go。

三、新的编译约束 `unix`

Go语言支持使用编译约束(build constraint)进行条件编译。Go 1.19版本新增了编译约束 `unix` ，可以在`//go:build`后面使用`unix`。

```go
//go:build unix
```

`unix`表示编译的目标操作系统是Unix或者类Unix系统。对于Go 1.19版本而言，如果`GOOS`是 `aix`, `android`, `darwin`, `dragonfly`, `freebsd`, `hurd`, `illumos`, `ios`, `linux`, `netbsd`, `openbsd`, 或 `solaris`中的某一个，那就满足`unix`这个编译约束。

未来`unix`约束还会匹配一些新的类Unix操作系统。



https://go.dev/doc/go1.19 Go 1.19 的发布说明文档

https://tonybai.com/2022/08/22/some-changes-in-go-1-19/ Go1.19值得关注的变化

https://segmentfault.com/a/1190000042005487 Go1.19 变化

## go 1.18

### 1、泛型引入

**泛型是Go诞生以来最复杂、最难读和理解的语法特性**

假设我们有个计算两数之和的函数：

```go
func Add(a int, b int) int {
    return a + b
}
```

这个函数很简单，但是它有个问题——无法计算int类型之外的和。如果我们想计算浮点或者字符串的和该怎么办？解决办法之一就是像下面这样为不同类型定义不同的函数

```go
func AddFloat32(a float32, b float32) float32 {
    return a + b
}

func AddString(a string, b string) string {
    return a + b
}
```

可是除此之外还有没有更好的方法？答案是有的，我们可以来回顾下函数的 **形参(parameter)** 和 **实参(argument)** 这一基本概念：

```go
func Add(a int, b int) int {  
    // 变量a,b是函数的形式参数   "a int, b int" 这一串被称为形参列表
    return a + b
}

Add(100,200) // 调用函数时，传入的100和200是实际参数
```

我们知道，函数的 **形参(parameter)** 只是类似占位符的东西并没有具体的值，只有我们调用函数传入**实参(argument)** 之后才有具体的值。

那么，如果我们将 **形参 实参** 这个概念推广一下，给变量的类型也引入和类似形参实参的概念的话，问题就迎刃而解：在这里我们将其称之为 **类型形参(type parameter)** 和 **类型实参(type argument)**，如下：

```go
// 假设 T 是类型形参，在定义函数时它的类型是不确定的，类似占位符
func Add(a T, b T) T {  
    return a + b
}
```

在上面这段伪代码中， T 被称为 **类型形参(type parameter)**， 它不是具体的类型，在定义函数时类型并不确定。因为 T 的类型并不确定，所以我们需要像函数的形参那样，在调用函数的时候再传入具体的类型。这样我们不就能一个函数同时支持多个不同的类型了吗？在这里被传入的具体类型被称为 **类型实参(type argument)**:

下面一段伪代码展示了调用函数时传入类型实参的方式：

```go
// [T=int]中的 int 是类型实参，代表着函数Add()定义中的类型形参 T 全都被 int 替换
Add[T=int](100, 200)  
// 传入类型实参int后，Add()函数的定义可近似看成下面这样：
func Add( a int, b int) int {
    return a + b
}

// 另一个例子：当我们想要计算两个字符串之和的时候，就传入string类型实参
Add[T=string]("Hello", "World") 
// 类型实参string传入后，Add()函数的定义可近似视为如下
func Add( a string, b string) string {
    return a + b
}
```

通过引入 **类型形参** 和 **类型实参** 这两个概念，我们让一个函数获得了处理多种不同类型数据的能力，这种编程方式被称为 **泛型编程**。

### 2、类型形参、类型实参

观察下面这个简单的例子：

```go
type IntSlice []int

var a IntSlice = []int{1, 2, 3} // 正确
var b IntSlice = []float32{1.0, 2.0, 3.0} // ✗ 错误，因为IntSlice的底层类型是[]int，浮点类型的切片无法赋值
```

这里定义了一个新的类型 `IntSlice` ，它的底层类型是 `[]int` ，理所当然只有int类型的切片能赋值给 `IntSlice` 类型的变量。

接下来如果我们想要定义一个可以容纳 `float32` 或 `string` 等其他类型的切片的话该怎么办？很简单，给每种类型都定义个新类型：

```go
type StringSlice []string
type Float32Slie []float32
type Float64Slice []float64
```

但是这样做的问题显而易见，它们结构都是一样的只是成员类型不同就需要重新定义这么多新类型。那么有没有一个办法能只定义一个类型就能代表上面这所有的类型呢？答案是可以的，这时候就需要用到泛型了：

```go
type Slice[T int|float32|float64 ] []T
```

不同于一般的类型定义，这里类型名称 `Slice` 后带了中括号，对各个部分做一个解说就是：

- `T` 就是上面介绍过的**类型形参(Type parameter)**，在定义Slice类型的时候 T 代表的具体类型并不确定，类似一个占位符
- `int|float32|float64` 这部分被称为**类型约束(Type constraint)**，中间的 `|` 的意思是告诉编译器，类型形参 T 只可以接收 int 或 float32 或 float64 这三种类型的实参
- 中括号里的 `T int|float32|float64` 这一整串因为定义了所有的类型形参(在这个例子里只有一个类型形参T），所以我们称其为 **类型形参列表(type parameter list)**
- 这里新定义的类型名称叫 `Slice[T]`

这种类型定义的方式中带了类型形参，很明显和普通的类型定义非常不一样，所以我们将这种

> 类型定义中带 **类型形参** **的类型，称之为** 泛型类型(Generic type)**

泛型类型不能直接拿来使用，必须传入**类型实参(Type argument)** 将其确定为具体的类型之后才可使用。而传入类型实参确定具体类型的操作被称为 **实例化(Instantiations)** ：

```go
// 这里传入了类型实参int，泛型类型Slice[T]被实例化为具体的类型 Slice[int]
var a Slice[int] = []int{1, 2, 3}  
fmt.Printf("Type Name: %T",a)  //输出：Type Name: Slice[int]

// 传入类型实参float32, 将泛型类型Slice[T]实例化为具体的类型 Slice[string]
var b Slice[float32] = []float32{1.0, 2.0, 3.0} 
fmt.Printf("Type Name: %T",b)  //输出：Type Name: Slice[float32]

// ✗ 错误。因为变量a的类型为Slice[int]，b的类型为Slice[float32]，两者类型不同
a = b  

// ✗ 错误。string不在类型约束 int|float32|float64 中，不能用来实例化泛型类型
var c Slice[string] = []string{"Hello", "World"} 

// ✗ 错误。Slice[T]是泛型类型，不可直接使用必须实例化为具体的类型
var x Slice[T] = []int{1, 2, 3} 
```

对于上面的例子，我们先给泛型类型 `Slice[T]` 传入了类型实参 `int` ，这样泛型类型就被实例化为了具体类型 `Slice[int]` ，被实例化之后的类型定义可近似视为如下：

```go
type Slice[int] []int     // 定义了一个普通的类型 Slice[int] ，它的底层类型是 []int
```

我们用实例化后的类型 `Slice[int]` 定义了一个新的变量 `a` ，这个变量可以存储int类型的切片。之后我们还用同样的方法实例化出了另一个类型 `Slice[float32]` ，并创建了变量 `b` 。

因为变量 a 和 b 就是具体的不同类型了(一个 Slice[int] ，一个 Slice[float32]），所以 `a = b` 这样不同类型之间的变量赋值是不允许的。

同时，因为 Slice[T] 的类型约束限定了只能使用 int 或 float32 或 float64 来实例化自己，所以 `Slice[string]` 这样使用 string 类型来实例化是错误的。

map简单示例：

```go
// MyMap类型定义了两个类型形参 KEY 和 VALUE。分别为两个形参指定了不同的类型约束
// 这个泛型类型的名字叫： MyMap[KEY, VALUE]
type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE  

// 用类型实参 string 和 flaot64 替换了类型形参 KEY 、 VALUE，泛型类型被实例化为具体的类型：MyMap[string, float64]
var a MyMap[string, float64] = map[string]float64{
    "jack_score": 9.6,
    "bob_score":  8.4,
}
```

用上面的例子重新复习下各种概念的话：

- KEY和VALUE是**类型形参**
- `int|string` 是KEY的**类型约束**， `float32|float64` 是VALUE的**类型约束**
- `KEY int|string, VALUE float32|float64` 整个一串文本因为定义了所有形参所以被称为**类型形参列表**
- Map[KEY, VALUE] 是**泛型类型**，类型的名字就叫 Map[KEY, VALUE]
- `var a MyMap[string, float64] = xx` 中的string和float64是**类型实参**，用于分别替换KEY和VALUE，**实例化**出了具体的类型 `MyMap[string, float64]`

![image-20240928213036621](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240928213036621.png)



### 3、泛型方法和泛型函数

#### 3.1 泛型方法

看了上的例子，你一定会说，介绍了这么多复杂的概念，但好像泛型类型根本没什么用处啊？

是的，单纯的泛型类型实际上对开发来说用处并不大。但是如果将泛型类型和接下来要介绍的泛型receiver相结合的话，泛型就有了非常大的实用性了

我们知道，定义了新的普通类型之后可以给类型添加方法。那么可以给泛型类型添加方法吗？答案自然是可以的，如下：

```go
package main

// 泛型
type MySlice[T int | float32] []T

func (s MySlice[T]) Sum() T {
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

// 普通
type MySlice2 []int // 实例化后的类型名叫 MySlice2[int]
// 方法中所有类型形参 T 都被替换为类型实参 int
func (s MySlice2) Sum2() int {
	var sum int
	for _, value := range s {
		sum += value
	}
	return sum
}
```

这个例子为泛型类型 `MySlice[T]` 添加了一个计算成员总和的方法 `Sum()` 。注意观察这个方法的定义：

- 首先看receiver `(s MySlice[T])` ，所以我们直接把类型名称 `MySlice[T]` 写入了receiver中
- 然后方法的返回参数我们使用了类型形参 T ****(实际上如果有需要的话，方法的接收参数也可以实用类型形参)
- 在方法的定义中，我们也可以使用类型形参 T （在这个例子里，我们通过 `var sum T` 定义了一个新的变量 `sum` )

对于这个泛型类型 `MySlice[T]` 我们该如何使用？泛型类型无论如何都需要先用类型实参实例化，所以用法如下：

```go
var s MySlice[int] = []int{1, 2, 3, 4}
fmt.Println(s.Sum()) // 输出：10

var s2 MySlice[float32] = []float32{1.0, 2.0, 3.0, 4.0}
fmt.Println(s2.Sum()) // 输出：10.0
```

#### 3.2 泛型函数

假设我们想要写一个计算两个数之和的函数：

```go
func Add(a int, b int) int {
    return a + b
}
```

这个函数理所当然只能计算int的和，而浮点的计算是不支持的。这时候我们可以像下面这样定义一个泛型函数：

```go
func Add[T int | float32 | float64](a T, b T) T {
    return a + b
}
```

上面就是泛型函数的定义。

> 这种带类型形参的函数被称为**泛型函数**

它和普通函数的点不同在于函数名之后带了类型形参。这里的类型形参的意义、写法和用法因为与泛型类型是一模一样的，就不再赘述了。

和泛型类型一样，泛型函数也是不能直接调用的，要使用泛型函数的话必须传入类型实参之后才能调用。

```go
Add[int](1,2) // 传入类型实参int，计算结果为 3
Add[float32](1.0, 2.0) // 传入类型实参float32, 计算结果为 3.0

Add[string]("hello", "world") // 错误。因为泛型函数Add的类型约束中并不包含string
```

或许你会觉得这样每次都要手动指定类型实参太不方便了。所以Go还支持类型实参的自动推导：

```go
Add(1, 2)  // 1，2是int类型，编译请自动推导出类型实参T是int
Add(1.0, 2.0) // 1.0, 2.0 是浮点，编译请自动推导出类型实参T是float32
```

自动推导的写法就好像免去了传入实参的步骤一样，但请记住这仅仅只是编译器帮我们推导出了类型实参，实际上传入实参步骤还是发生了的。

### 4、自定义泛型约束

有时候使用泛型编程时，我们会书写长长的类型约束，如下：

```go
// 一个可以容纳所有int,uint以及浮点类型的泛型切片
type Slice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] []T
```

理所当然，这种写法是我们无法忍受也难以维护的，而Go支持将类型约束单独拿出来定义到接口中，从而让代码更容易维护：

```go
type IntUintFloat interface {
    int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}
type Slice[T IntUintFloat] []T
```

这段代码把类型约束给单独拿出来，写入了接口类型 `IntUintFloat` 当中。需要指定类型约束的时候直接使用接口 `IntUintFloat` 即可。

不过这样的代码依旧不好维护，而接口和接口、接口和普通类型之间也是可以通过 `|` 进行组合：

```go
type Int interface {
    int | int8 | int16 | int32 | int64
}

type Uint interface {
    uint | uint8 | uint16 | uint32
}

type Float interface {
    float32 | float64
}

type Slice[T Int | Uint | Float] []T  // 使用 '|' 将多个接口类型组合
```

上面的代码中，我们分别定义了 Int, Uint, Float 三个接口类型，并最终在 Slice[T] 的类型约束中通过使用 `|` 将它们组合到一起。

同时，在接口里也能直接组合其他接口，所以还可以像下面这样：

```go
type SliceElement interface {
    Int | Uint | Float | string // 组合了三个接口类型并额外增加了一个 string 类型
}

type Slice[T SliceElement] []T 
```

上面定义的 Slie[T] 虽然可以达到目的，但是有一个缺点：

```go
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
```

![image-20240928215640666](../../Library/Application Support/typora-user-images/image-20240928215640666.png)

这里发生错误的原因是，泛型类型 Slice[T] 允许的是 int 作为类型实参，而不是 MyInt2 （虽然 MyInt2 类型底层类型是 int ，但它依旧不是 int 类型）。

为了从根本上解决这个问题，Go新增了一个符号 `~` ，在类型约束中使用类似 `~int` 这种写法的话，就代表着不光是 int ，所有以 int 为底层类型的类型也都可用于实例化。

使用 ~ 对代码进行改写之后如下：

```go
type Int interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32
}
type Float interface {
    ~float32 | ~float64
}

type Slice[T Int | Uint | Float] []T 

var s Slice[int] // 正确

type MyInt int
var s2 Slice[MyInt]  // MyInt底层类型是int，所以可以用于实例化

type MyMyInt MyInt
var s3 Slice[MyMyInt]  // 正确。MyMyInt 虽然基于 MyInt ，但底层类型也是int，所以也能用于实例化

type MyFloat32 float32  // 正确
var s4 Slice[MyFloat32]
```

**限制**：使用 `~` 时有一定的限制：

1. ~后面的类型不能为接口
2. ~后面的类型必须为基本类型

```go
type MyInt int

type _ interface {
    ~[]byte  // 正确
    ~MyInt   // 错误，~后的类型必须为基本类型
    ~error   // 错误，~后的类型不能为接口
}
```

https://segmentfault.com/a/1190000041634906  Go1.18泛型全面讲清楚

### 2、**模糊测试Fuzzing**

#### 2.1 Fuzzing介绍及特点：

Fuzzing中文含义是[模糊测试](https://zhida.zhihu.com/search?content_id=197551581&content_type=Article&match_order=2&q=模糊测试&zhida_source=entity)，是一种自动化测试技术，可以随机生成[测试数据](https://zhida.zhihu.com/search?content_id=197551581&content_type=Article&match_order=1&q=测试数据&zhida_source=entity)集，然后调用要测试的功能代码来检查功能是否符合预期。

-  模糊测试(fuzz test)是对[单元测试](https://zhida.zhihu.com/search?content_id=197551581&content_type=Article&match_order=1&q=单元测试&zhida_source=entity)(unit test)的补充，并不是要替代单元测试。

- 单元测试是检查指定的输入得到的结果是否和预期的输出结果一致，测试数据集比较有限。

- 模糊测试可以生成随机测试数据，找出单元测试覆盖不到的场景，进而发现程序的潜在bug和安全漏洞。

#### 2.2 Fuzzing 使用

Go Fuzzing模糊测试函数的语法如下所示：

- 模糊测试函数定义在`xxx_test.go`文件里，这点和Go已有的单元测试(unit test)和性能测试(benchmark test)一样。
- 函数名以`Fuzz`开头，参数是`* testing.F`类型，`testing.F`类型有2个重要方法`Add`和`Fuzz`。
- `Add`方法是用于添加种子语料(seed corpus)数据，Fuzzing底层可以根据种子语料数据自动生成随机测试数据。
- `Fuzz`方法接收一个函数类型的变量作为参数，该函数类型的第一个参数必须是`*testing.T`类型，其余的参数类型和`Add`方法里传入的实参类型保持一致

<img src="https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240926234858885.png" alt="image-20240926234858885" style="zoom:50%;" />

我们声明如下翻转函数：

```go
func Reverse(s string) string {
    b := []byte(s)
    for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}
func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := utils.Reverse(input)
	doubleRev := utils.Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
输出：
$ original: "The quick brown fox jumped over the lazy dog"
$ reversed: "god yzal eht revo depmuj xof nworb kciuq ehT"
$ reversed again: "The quick brown fox jumped over the lazy dog"

```

添加单元测试并验证：

```go
func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, tc := range testcases {
		rev := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}
输出：
--- PASS: TestReverse (0.00s)
PASS
```

将增加模糊测试：

```go
func FuzzReverse(f *testing.F) {
    testcases := []string{"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc)  // Use f.Add to provide a seed corpus
    }
    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)
        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
# 执行
go test -v -fuzz=FuzzReverse
```

![image-20240927000833414](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240927000833414.png)

这个例子里，随机生成了一个字符串，这是由2个字节组成的一个UTF-8字符串，按照`Reverse`函数进行反转后，得到了一个非UTF-8的字符串

所以我们之前实现的按照字节进行字符串反转的函数`Reverse`是有bug的，该函数对于ASCII码里的字符组成的字符串是可以正确反转的，但是对于非ASCII码里的字符，如果简单按照字节进行反转，得到的可能是一个非法的字符串。

修改后的反转函数：

```go
func Reverse(s string) (string, error) {
    if !utf8.ValidString(s) { // 价差字符串是否有效的UTF-8编码字符串，如果不是则返回空字符串
        return s, errors.New("input is not valid UTF-8")
    }
    r := []rune(s) // 将字符串转换为rune切片，确保可以处理多字节字符 如中文
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r), nil
}
```

![image-20240927000415512](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240927000415512.png)

https://tonybai.com/2021/12/01/first-class-fuzzing-in-go-1-18 

https://pkg.go.dev/testing@master#hdr-Fuzzing

https://go.googlesource.com/proposal/+/master/design/draft-fuzzing.md

### 3、**Workspaces**

解决go mod遗留下来的**本地多模块开发依赖问题**

在没有 go1.18 之前，只能使用 replace，如下：

```go
replace (
 common =>  "xx/xx/common"
)
```

当前我们的项目目录如下:

```go
.
├── common // 子模块
│   ├── go.mod
│   └── util.go
└── hello // 子模块
    ├── go.mod
    └── main.go
```

进入项目目录, 我们使用`go work init ./hello `来初始化一个新的工作区, 同时加入需要的的子模块

```go
.
├── common // 子模块
│   ├── go.mod
│   └── util.go
├── go.work //工作区
└── hello // 子模块
    ├── go.mod
    └── main.go

```

把 common项目移动到 workspace 目录下，然后执行：

```go
go work use ./common
```

work工作区:

```go
go 1.19

use (
	./common
	./hello
)
```

main函数代码：执行命令：

```go
package main

import (
	"example.com/common"
	"fmt"
)

func main() {
	common.PrintUtil()
	fmt.Printf("this is hello")
}

```

执行命令：

```go
thie is common
this is hello%  
```

参考：

https://go.dev/doc/go1.18 Go 1.18 的发布说明文档

https://tonybai.com/2022/04/20/some-changes-in-go-1-18/ go1.18值得关注的几个变化

https://go.dev/doc/tutorial/workspaces 

https://go.googlesource.com/proposal/+/master/design/45713-workspace.md

## go 1.17

### 语法特性

#### 1.1、支持将切片转换为数组指针

基础功能，通过数组切片化，我们可以将一个数组转换为切片，转换后，数组将成为转换后的切片的底层数组，通过切片，我们可以直接改变数组中的元素

```go
func array2slice() {
	a := [3]int{11, 12, 13}
	b := a[:]
	b[1] += 10
	fmt.Printf("%v\n", a) // [11 22 13]
}
```

在 Go 1.17 版本中，我们可以像下面代码这样将一个切片转换为数组类型指针

```go
b := []int{11, 12, 13}
p := (*[3]int)(b) // 将切片转换为数组类型指针
p[1] = p[1] + 10
fmt.Printf("%v\n", b) // [11 22 13]
```

不过，这里你要注意的是，Go 会通过运行时而不是编译器去对这类切片到数组指针的转换代码做检查，如果发现越界行为，就会触发运行时 panic。Go 运行时实施检查的一条原则就是“转换后的数组长度不能大于原切片的长度”，注意这里是切片的长度（len），而不是切片的容量（cap）。于是你会看到，下面的转换有些合法，有些非法：

```go
var b = []int{11, 12, 13}
var p = (*[4]int)(b) // cannot convert slice with length 3 to pointer to array with length 4
var p = (*[0]int)(b) // ok，*p = []
var p = (*[1]int)(b) // ok，*p = [11]
var p = (*[2]int)(b) // ok，*p = [11, 12]
var p = (*[3]int)(b) // ok，*p = [11, 12, 13]
var p = (*[3]int)(b[:1]) // cannot convert slice with length 1 to pointer to array with length 3 
```

另外，nil 切片或 cap 为 0 的 empty 切片都可以被转换为一个长度为 0 的数组指针，比如：

```go
var b1 []int // nil切片
p1 := (*[0]int)(b1)
var b2 = []int{} // empty切片
p2 := (*[0]int)(b2)
```

### 2、Go Module 构建模式的变化

#### 2.1、修剪的 module 依赖图

Go 1.17 版本中，Go Module 最重要的一个变化就是 pruned module graph，即**修剪的 module 依赖图**

用下图中的例子来详细解释一下 module 依赖图修剪的原理

![image-20240919230623850](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240919230623850.png)

在这个示例中，main module 中的 lazy.go 导入了 module a 的 package x，后者则导入了 module b 中的 package b。并且，module a 还有一个 package y，这个包导入了 module c 的 package c。通过 go mod graph 命令，我们可以得到 main module 的完整 module 依赖图，也就是上图的右上角的那张。

现在问题来了！package y 是因为自身是 module a 的一部分而被 main module 依赖的，它自己没有为 main module 的构建做出任何“代码级贡献”，同理，package y 所依赖的 module c 亦是如此。但是在 Go 1.17 之前的版本中，如果 Go 编译器找不到 module c，那么 main module 的构建也会失败，这会让开发者们觉得不够合理！

在 Go 1.16.5 下，这个示例的 go.mod：

```go
module example.com/lazy

go 1.15

// require 块中保留的都是 main module 的直接依赖
require example.com/a v0.1.0

// replace 块主要是为了能找到各种依赖 module 而设置的
replace (
  example.com/a v0.1.0 => ./a
  example.com/b v0.1.0 => ./b
  example.com/c v0.1.0 => ./c1
  example.com/c v0.2.0 => ./c2
)
```

此时，如果我们将 replace 中的第三行（example.com/c v0.1.0 => ./c1 这一行）删除，也就是让 Go 编译器找不到 module c@v0.1.0，那么我们构建 main modue 时就会得到下面的错误提示：

```go
$go build
go: example.com/a@v0.1.0 requires
  example.com/c@v0.1.0: missing go.sum entry; to add it:
  go mod download example.com/c
```

在 Go 1.17 下，示例的 go.mod：

```go
$go mod tidy
$cat go.mod

module example.com/lazy

go 1.17

require example.com/a v0.1.0

require example.com/b v0.1.0 // indirect

replace (
  example.com/a v0.1.0 => ./a
  example.com/b v0.1.0 => ./b
  example.com/c v0.1.0 => ./c1
  example.com/c v0.2.0 => ./c2
)
```

将 go.mod replace 块中的第三行（example.com/c v0.1.0 => ./c1 这一行）删除，再来用 go 1.17 构建一次 main module

这一次我们没有看到 Go 编译器的错误提示。也就是说在构建过程中，Go 编译器看到的 main module 依赖图中并没有 module c@v0.1.0 ，这种将那些对构建完全没有“贡献”的间接依赖 module 从构建时使用的依赖图中修剪掉的过程，就被称为 **module 依赖图修剪（pruned module graph）**

### 、go get 不再被用来安装命令可执行文件

新版本中，我们需要使用 go install 来安装，并且使用 go install 安装时还要用 @vx.y.z 明确要安装的命令的二进制文件的版本，或者是使用 @latest 来安装最新版本

### 4、//go:build 形式的构建约束指示符

Go 编译器还在 Go 1.17 中引入了 //go:build 形式的构建约束指示符，以替代原先易错的 // +build 形式

在 Go 1.17 之前，我们可以通过在源码文件头部放置 // +build 构建约束指示符来实现构建约束，但这种形式十分易错，并且它并不支持 && 和||这样的直观的逻辑操作符，而是用逗号、空格替代，这里你可以看下原 // +build 形式构建约束指示符的用法及含义：

![image-20240919232304938](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240919232304938.png)

go1.17 新形式将支持 && 和||逻辑操作符

```go
//go:build linux && (386 || amd64 || arm || arm64 || mips64 || mips64le || ppc64 || ppc64le)
//go:build linux && (mips64 || mips64le)
//go:build linux && (ppc64 || ppc64le)
//go:build linux && !386 && !arm
```

### 5、基于寄存器的调用惯例

所谓“调用惯例（calling convention）”，是指调用方和被调用方对于函数调用的一个明确的约定，包括函数参数与返回值的传递方式、传递顺序。只有双方都遵守同样的约定，函数才能被正确地调用和执行。如果不遵守这个约定，函数将无法正确执行。

先把结论提前摆出来：

1）Go1.17.1之前的函数调用，参数都在栈上传递；Go1.17.1以后，9个以内的参数在寄存器传递，9个以外的在栈上传递；

2）Go1.17.1之前版本，callee(调用者)函数返回值通过caller(被调用者)栈传递；Go1.17.1以后，函数调用的返回值，9个以内通过寄存器传递回caller，9个以外在栈上传递；

3）由于CPU访问寄存器的速度要远高于栈内存，各业务团队将自己程序的Go版本升级到Go1.17以上，能够提高程序性能。

4）在Go 1.17的版本发布说明文档中有提到：切换到基于寄存器的调用惯例后，一组有代表性的Go包和程序的基准测试显示，Go程序的运行性能提高了约5%，二进制文件大小减少约2%。

#### go 1.16函数版本分析：

![image-20240922210128007](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240922210128007.png)

```go
package main

func main() {
    var r1, r2, r3, r4, r5, r6, r7 int64 = 1, 2, 3, 4, 5, 6, 7
    A(r1, r2, r3, r4, r5, r6, r7)
}

func Add(p1, p2, p3, p4, p5, p6, p7 int64) int64 {
    return p1 + p2 + p3 + p4 + p5 + p6 + p7
}
```

分析汇编

```go
# GOOS=linux GOARCH=amd64 go tool compile -S -N -l main.go

# "".main STEXT size=189 args=0x0 locals=0x80 funcid=0x0
        0x0000 00000 (main.go:3)        TEXT    "".main(SB), ABIInternal, $128-0 #main函数定义, $128-0：128表示将分配的main函数的栈帧大小；0指定了调用方传入的参数，由于main是最上层函数，这里没有入参
        0x0000 00000 (main.go:3)        MOVQ    (TLS), CX # 将本地线程存储信息保存到CX寄存器中
        0x0009 00009 (main.go:3)        CMPQ    SP, 16(CX) # 栈溢出检测：比较当前栈顶地址(SP寄存器存放的)与本地线程存储的栈顶地址
        0x000d 00013 (main.go:3)        PCDATA  $0, $-2 # PCDATA，FUNCDATA用于Go汇编额外信息，不必关注
        0x000d 00013 (main.go:3)        PCDATA  $0, $-2            
        0x000d 00013 (main.go:3)        JLS     179 # 如果当前栈顶地址(SP寄存器存放的)小于本地线程存储的栈顶地址，则跳到180处代码处进行栈分裂扩容操作
        0x0013 00019 (main.go:3)        PCDATA  $0, $-1
        0x0013 00019 (main.go:3)        ADDQ    $-128, SP # 为main函数栈帧分配了128字节的空间，注意此时的SP寄存器指向，会往下移动128个字节
        0x0017 00023 (main.go:3)        MOVQ    BP, 120(SP) # BP寄存器存放的是main函数caller的基址，movq这条指令是将main函数caller的基址入栈。
        0x001c 00028 (main.go:3)        LEAQ    120(SP), BP # 将main函数的基址存放到到BP寄存器
        0x0021 00033 (main.go:3)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0021 00033 (main.go:3)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0021 00033 (main.go:4)        MOVQ    $1, "".r1+112(SP) # main函数局部变量r1入栈
        0x002a 00042 (main.go:4)        MOVQ    $2, "".r2+104(SP) # main函数局部变量r2入栈
        0x0033 00051 (main.go:4)        MOVQ    $3, "".r3+96(SP) # main函数局部变量r3入栈
        0x003c 00060 (main.go:4)        MOVQ    $4, "".r4+88(SP) # main函数局部变量r4入栈
        0x0045 00069 (main.go:4)        MOVQ    $5, "".r5+80(SP) # main函数局部变量r5入栈
        0x004e 00078 (main.go:4)        MOVQ    $6, "".r6+72(SP) # main函数局部变量r6入栈
        0x0057 00087 (main.go:4)        MOVQ    $7, "".r7+64(SP) # main函数局部变量r7入栈
        0x0060 00096 (main.go:5)        MOVQ    "".r2+104(SP), AX # 将局部变量r2传给寄存器AX
        0x0065 00101 (main.go:5)        MOVQ    "".r3+96(SP), CX # 将局部变量r3传给寄存器CX
        0x006a 00106 (main.go:5)        MOVQ    "".r4+88(SP), DX # 将局部变量r4传给寄存器DX
        0x006f 00111 (main.go:5)        MOVQ    "".r5+80(SP), BX # 将局部变量r5传给寄存器BX
        0x0074 00116 (main.go:5)        MOVQ    "".r6+72(SP), SI # 将局部变量r6传给寄存器SI
        0x0079 00121 (main.go:5)        MOVQ    "".r1+112(SP), DI # 将局部变量r1传给寄存器DI
        0x007e 00126 (main.go:5)        MOVQ    DI, (SP) # 寄存器DI将局部变量r1加入栈头SP指向的位置
        0x0082 00130 (main.go:5)        MOVQ    AX, 8(SP) # 寄存器AX将局部变量r2加入栈头SP+8指向的位置
        0x0087 00135 (main.go:5)        MOVQ    CX, 16(SP) # 寄存器CX将局部变量r3加入栈头SP+16指向的位置 
        0x008c 00140 (main.go:5)        MOVQ    DX, 24(SP)  # 寄存器DX将局部变量r4加入栈头SP+24指向的位置 
        0x0091 00145 (main.go:5)        MOVQ    BX, 32(SP)  # 寄存器BX将局部变量r3加入栈头SP+32指向的位置 
        0x0096 00150 (main.go:5)        MOVQ    SI, 40(SP)  # 寄存器SI将局部变量r3加入栈头SP+40指向的位置 
        0x009b 00155 (main.go:5)        MOVQ    $7, 48(SP)  # 将局部变量r7加入栈头SP+48指向的位置 
        0x00a4 00164 (main.go:5)        PCDATA  $1, $0
        0x00a4 00164 (main.go:5)        CALL    "".Add(SB) # 调用 A函数
        0x00a9 00169 (main.go:6)        MOVQ    120(SP), BP # 将栈上存储的main函数的调用方的基地址恢复到BP
        0x00ae 00174 (main.go:6)        SUBQ    $-128, SP # 增加SP的值，栈收缩，收回分配给main函数栈帧的128字节空间
        0x00b2 00178 (main.go:6)        RET
        0x00b3 00179 (main.go:6)        NOP
        0x00b3 00179 (main.go:3)        PCDATA  $1, $-1
        0x00b3 00179 (main.go:3)        PCDATA  $0, $-2
        0x00b3 00179 (main.go:3)        CALL    runtime.morestack_noctxt(SB)
        0x00b8 00184 (main.go:3)        PCDATA  $0, $-1
        0x00b8 00184 (main.go:3)        JMP     0
        0x0000 64 48 8b 0c 25 00 00 00 00 48 3b 61 10 0f 86 a0  dH..%....H;a....
        0x0010 00 00 00 48 83 c4 80 48 89 6c 24 78 48 8d 6c 24  ...H...H.l$xH.l$
        0x0020 78 48 c7 44 24 70 01 00 00 00 48 c7 44 24 68 02  xH.D$p....H.D$h.
        0x0030 00 00 00 48 c7 44 24 60 03 00 00 00 48 c7 44 24  ...H.D$`....H.D$
        0x0040 58 04 00 00 00 48 c7 44 24 50 05 00 00 00 48 c7  X....H.D$P....H.
        0x0050 44 24 48 06 00 00 00 48 c7 44 24 40 07 00 00 00  D$H....H.D$@....
        0x0060 48 8b 44 24 68 48 8b 4c 24 60 48 8b 54 24 58 48  H.D$hH.L$`H.T$XH
        0x0070 8b 5c 24 50 48 8b 74 24 48 48 8b 7c 24 70 48 89  .\$PH.t$HH.|$pH.
        0x0080 3c 24 48 89 44 24 08 48 89 4c 24 10 48 89 54 24  <$H.D$.H.L$.H.T$
        0x0090 18 48 89 5c 24 20 48 89 74 24 28 48 c7 44 24 30  .H.\$ H.t$(H.D$0
        0x00a0 07 00 00 00 e8 00 00 00 00 48 8b 6c 24 78 48 83  .........H.l$xH.
        0x00b0 ec 80 c3 e8 00 00 00 00 e9 43 ff ff ff           .........C...
        rel 5+4 t=17 TLS+0
        rel 165+4 t=8 "".Add+0
        rel 180+4 t=8 runtime.morestack_noctxt+0
"".Add STEXT nosplit size=50 args=0x40 locals=0x0 funcid=0x0
        0x0000 00000 (main.go:8)        TEXT    "".Add(SB), NOSPLIT|ABIInternal, $0-64 #A函数定义, $0-64：0表示将分配的A函数的栈帧大小；64指定了调用方传入的参数和函数的返回值的大小，入参7个，返回值1个，都是8字节，共64字节
        0x0000 00000 (main.go:8)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:8)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:8)        MOVQ    $0, "".~r7+64(SP) # 这里 SP+64就是上面main栈空间中用来接收返回值的地址
        0x0009 00009 (main.go:9)        MOVQ    "".p1+8(SP), AX # A返回值和r1参数求和后，放入AX寄存器
        0x000e 00014 (main.go:9)        ADDQ    "".p2+16(SP), AX # AX寄存器的值再和r2参数求和，结果放入AX
        0x0013 00019 (main.go:9)        ADDQ    "".p3+24(SP), AX # AX寄存器的值再和r3参数求和，结果放入AX
        0x0018 00024 (main.go:9)        ADDQ    "".p4+32(SP), AX # AX寄存器的值再和r4参数求和，结果放入AX
        0x001d 00029 (main.go:9)        ADDQ    "".p5+40(SP), AX # AX寄存器的值再和r5参数求和，结果放入AX
        0x0022 00034 (main.go:9)        ADDQ    "".p6+48(SP), AX # AX寄存器的值再和r6参数求和，结果放入AX
        0x0027 00039 (main.go:9)        ADDQ    "".p7+56(SP), AX # AX寄存器的值再和r7参数求和，结果放入AX
        0x002c 00044 (main.go:9)        MOVQ    AX, "".~r7+64(SP) # AX寄存器的值 写回main栈空间中用来接收返回值的地址SP+64中
        0x0031 00049 (main.go:9)        RET
        0x0000 48 c7 44 24 40 00 00 00 00 48 8b 44 24 08 48 03  H.D$@....H.D$.H.
        0x0010 44 24 10 48 03 44 24 18 48 03 44 24 20 48 03 44  D$.H.D$.H.D$ H.D
        0x0020 24 28 48 03 44 24 30 48 03 44 24 38 48 89 44 24  $(H.D$0H.D$8H.D$
        0x0030 40 c3                                            @.
go.cuinfo.packagename. SDWARFCUINFO dupok size=0
        0x0000 6d 61 69 6e                                      main
""..inittask SNOPTRDATA size=24
        0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 00 00 00 00 00 00 00 00                          ........
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
```

#### go 1.17函数版本分析：

```go
package main

func main() {
    var r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11 int64 = 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
    a, b := A(r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11)
    c := a + b
    print(c)
}

func A(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11 int64) (int64, int64) {
    return p1 + p2 + p3 + p4 + p5 + p6 + p7, p2 + p4 + p6 + p7 + p8 + p9 + p10 + p11  
}
```

分析main函数的汇编代码：

```go
# GOOS=linux GOARCH=amd64 go tool compile -S -N -l main.go
"".main STEXT size=362 args=0x0 locals=0xe0 funcid=0x0
        0x0000 00000 (main.go:3)        TEXT    "".main(SB), ABIInternal, $224-0    #main函数定义, $224-0：224表示将分配的main函数的栈帧大小；0指定了调用方传入的参数，由于main是最上层函数，这里没有入参
        0x0000 00000 (main.go:3)        LEAQ    -96(SP), R12
        0x0005 00005 (main.go:3)        CMPQ    R12, 16(R14)
        0x0009 00009 (main.go:3)        PCDATA  $0, $-2
        0x0009 00009 (main.go:3)        JLS     349
        0x000f 00015 (main.go:3)        PCDATA  $0, $-1
        0x000f 00015 (main.go:3)        SUBQ    $224, SP                     # 为main函数栈帧分配了224字节的空间，注意此时的SP寄存器指向，会往下移动224个字节
        0x0016 00022 (main.go:3)        MOVQ    BP, 216(SP)                  # BP寄存器存放的是main函数caller的基址，movq这条指令是将main函数caller的基址入栈
        0x001e 00030 (main.go:3)        LEAQ    216(SP), BP                  # 将main函数的基址存放到到BP寄存器
        0x0026 00038 (main.go:3)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0026 00038 (main.go:3)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0026 00038 (main.go:4)        MOVQ    $1, "".r1+168(SP)            # main函数局部变量r1入栈    
        0x0032 00050 (main.go:4)        MOVQ    $2, "".r2+144(SP)            # main函数局部变量r2入栈
        0x003e 00062 (main.go:4)        MOVQ    $3, "".r3+136(SP)            # main函数局部变量r3入栈
        0x004a 00074 (main.go:4)        MOVQ    $4, "".r4+128(SP)            # main函数局部变量r4入栈
        0x0056 00086 (main.go:4)        MOVQ    $5, "".r5+120(SP)            # main函数局部变量r5入栈
        0x005f 00095 (main.go:4)        MOVQ    $6, "".r6+112(SP)            # main函数局部变量r6入栈
        0x0068 00104 (main.go:4)        MOVQ    $7, "".r7+104(SP)            # main函数局部变量r7入栈
        0x0071 00113 (main.go:4)        MOVQ    $8, "".r8+96(SP)             # main函数局部变量r8入栈 
        0x007a 00122 (main.go:4)        MOVQ    $9, "".r9+88(SP)             # main函数局部变量r9入栈
        0x0083 00131 (main.go:4)        MOVQ    $10, "".r10+160(SP)          # main函数局部变量r10入栈
        0x008f 00143 (main.go:4)        MOVQ    $11, "".r11+152(SP)          # main函数局部变量r11入栈
        0x009b 00155 (main.go:5)        MOVQ    "".r2+144(SP), BX            # 将局部变量r2传给寄存器BX
        0x00a3 00163 (main.go:5)        MOVQ    "".r3+136(SP), CX            # 将局部变量r3传给寄存器CX
        0x00ab 00171 (main.go:5)        MOVQ    "".r4+128(SP), DI            # 将局部变量r4传给寄存器DI
        0x00b3 00179 (main.go:5)        MOVQ    "".r5+120(SP), SI            # 将局部变量r5传给寄存器SI
        0x00b8 00184 (main.go:5)        MOVQ    "".r6+112(SP), R8            # 将局部变量r6传给寄存器R8
        0x00bd 00189 (main.go:5)        MOVQ    "".r7+104(SP), R9            # 将局部变量r7传给寄存器R9
        0x00c2 00194 (main.go:5)        MOVQ    "".r8+96(SP), R10            # 将局部变量r8传给寄存器R10
        0x00c7 00199 (main.go:5)        MOVQ    "".r9+88(SP), R11            # 将局部变量r9传给寄存器R11
        0x00cc 00204 (main.go:5)        MOVQ    "".r10+160(SP), DX           # 将局部变量r10传给寄存器DX
        0x00d4 00212 (main.go:5)        MOVQ    "".r1+168(SP), AX            # 将局部变量r1传给寄存器DX
        0x00dc 00220 (main.go:5)        MOVQ    DX, (SP)                     # 将寄存器DX保存的r10传给SP指向的栈顶
        0x00e0 00224 (main.go:5)        MOVQ    $11, 8(SP)                   # 将变量r11传给SP+8
        0x00e9 00233 (main.go:5)        PCDATA  $1, $0
        0x00e9 00233 (main.go:5)        CALL    "".A(SB)                     # 调用 A 函数
        0x00ee 00238 (main.go:5)        MOVQ    AX, ""..autotmp_14+208(SP)   # 将寄存器AX存的函数A的第一个返回值a赋值给SP+208
        0x00f6 00246 (main.go:5)        MOVQ    BX, ""..autotmp_15+200(SP)   # 将寄存器BX存的函数A的第二个返回值b赋值给SP+200
        0x00fe 00254 (main.go:5)        MOVQ    ""..autotmp_14+208(SP), DX   # 将SP+208保存的A函数第一个返回值a传给寄存器DX
        0x0106 00262 (main.go:5)        MOVQ    DX, "".a+192(SP)             # 将A函数第一个返回值a通过寄存器DX入栈到SP+192
        0x010e 00270 (main.go:5)        MOVQ    ""..autotmp_15+200(SP), DX   # 将SP+200保存的A函数第二个返回值b传给寄存器DX
        0x0116 00278 (main.go:5)        MOVQ    DX, "".b+184(SP)             # 将第二个返回值b通过寄存器DX入栈到SP+184
        0x011e 00286 (main.go:6)        MOVQ    "".a+192(SP), DX             # 将返回值a传给DX寄存器
        0x0126 00294 (main.go:6)        ADDQ    "".b+184(SP), DX             # 将a+b赋值给DX寄存器
        0x012e 00302 (main.go:6)        MOVQ    DX, "".c+176(SP)             # 将DX寄存器的值入栈到SP+176
        0x0136 00310 (main.go:7)        CALL    runtime.printlock(SB)        
        0x013b 00315 (main.go:7)        MOVQ    "".c+176(SP), AX             # 将SP+176存储的入参c赋值给AX
        0x0143 00323 (main.go:7)        CALL    runtime.printint(SB)         # 调用打印函数打印c
        0x0148 00328 (main.go:7)        CALL    runtime.printunlock(SB)
        0x014d 00333 (main.go:8)        MOVQ    216(SP), BP
        0x0155 00341 (main.go:8)        ADDQ    $224, SP
        0x015c 00348 (main.go:8)        RET
```

通过上面汇编代码的注释，我们可以看到，main函数调用A函数的参数个数为11个，其中前 9 个参数分别是通过寄存器 AX，BX，CX，DI，SI，R8，R9, R10, R11传递，后面两个通过栈顶的SP，SP+8地址传递。

![img](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/umzs9oqs3g.png)

下面看 Add函数在Go1.17.1的汇编代码：

```go
"".A STEXT nosplit size=175 args=0x58 locals=0x18 funcid=0x0
        0x0000 00000 (main.go:10)       TEXT    "".A(SB), NOSPLIT|ABIInternal, $24-88
        0x0000 00000 (main.go:10)       SUBQ    $24, SP                        # 为A函数栈帧分配了24字节的空间
        0x0004 00004 (main.go:10)       MOVQ    BP, 16(SP)
        0x0009 00009 (main.go:10)       LEAQ    16(SP), BP
        0x000e 00014 (main.go:10)       FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (main.go:10)       FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (main.go:10)       FUNCDATA        $5, "".A.arginfo1(SB)
        0x000e 00014 (main.go:10)       MOVQ    AX, "".p1+48(SP)               # 寄存器AX存储的r1赋值给SP+48    
        0x0013 00019 (main.go:10)       MOVQ    BX, "".p2+56(SP)               # 寄存器BX存储的r2赋值给SP+56
        0x0018 00024 (main.go:10)       MOVQ    CX, "".p3+64(SP)               # 寄存器CX存储的r3赋值给SP+64
        0x001d 00029 (main.go:10)       MOVQ    DI, "".p4+72(SP)               # 寄存器DI存储的r4赋值给SP+72
        0x0022 00034 (main.go:10)       MOVQ    SI, "".p5+80(SP)               # 寄存器SI存储的r5赋值给SP+80
        0x0027 00039 (main.go:10)       MOVQ    R8, "".p6+88(SP)               # 寄存器R8存储的r6赋值给SP+88
        0x002c 00044 (main.go:10)       MOVQ    R9, "".p7+96(SP)               # 寄存器R9存储的r7赋值给SP+96
        0x0031 00049 (main.go:10)       MOVQ    R10, "".p8+104(SP)             # 寄存器R10存储的r8赋值给SP+104
        0x0036 00054 (main.go:10)       MOVQ    R11, "".p9+112(SP)             # 寄存器R11存储的r9赋值给SP+112
        0x003b 00059 (main.go:10)       MOVQ    $0, "".~r11+8(SP)              # 初始化第一个返回值a存放地址SP+8为0
        0x0044 00068 (main.go:10)       MOVQ    $0, "".~r12(SP)                # 初始化第二个返回值b存放地址SP为0
        0x004c 00076 (main.go:11)       MOVQ    "".p1+48(SP), CX               # SP+48存储的r1赋值给CX寄存器
        0x0051 00081 (main.go:11)       ADDQ    "".p2+56(SP), CX               # CX+r2赋值给CX寄存器
        0x0056 00086 (main.go:11)       ADDQ    "".p3+64(SP), CX               # CX+r3赋值给CX寄存器
        0x005b 00091 (main.go:11)       ADDQ    "".p4+72(SP), CX               # CX+r4赋值给CX寄存器
        0x0060 00096 (main.go:11)       ADDQ    "".p5+80(SP), CX               # CX+r5赋值给CX寄存器
        0x0065 00101 (main.go:11)       ADDQ    "".p6+88(SP), CX               # CX+r6赋值给CX寄存器
        0x006a 00106 (main.go:11)       ADDQ    "".p7+96(SP), CX               # CX+r7赋值给CX寄存器
        0x006f 00111 (main.go:11)       MOVQ    CX, "".~r11+8(SP)              # CX寄存器赋值给第一个返回值存放地址SP+8
        0x0074 00116 (main.go:11)       MOVQ    "".p2+56(SP), BX               # r2赋值给BX寄存器
        0x0079 00121 (main.go:11)       ADDQ    "".p4+72(SP), BX               # BX+r4赋值给BX寄存器
        0x007e 00126 (main.go:11)       ADDQ    "".p6+88(SP), BX               # BX+r6赋值给BX寄存器
        0x0083 00131 (main.go:11)       ADDQ    "".p7+96(SP), BX               # BX+r7赋值给BX寄存器
        0x0088 00136 (main.go:11)       ADDQ    "".p8+104(SP), BX              # BX+r8赋值给BX寄存器
        0x008d 00141 (main.go:11)       ADDQ    "".p9+112(SP), BX              # BX+r9赋值给BX寄存器
        0x0092 00146 (main.go:11)       ADDQ    "".p10+32(SP), BX              # BX+r11赋值给BX寄存器
        0x0097 00151 (main.go:11)       ADDQ    "".p11+40(SP), BX              # BX+r10赋值给BX寄存器
        0x009c 00156 (main.go:11)       MOVQ    BX, "".~r12(SP)                # BX寄存器赋值给第二个返回值存放地址SP
        0x00a0 00160 (main.go:11)       MOVQ    "".~r11+8(SP), AX              # 第一个返回值SP+8的值赋值给AX寄存器
        0x00a5 00165 (main.go:11)       MOVQ    16(SP), BP                     # main返回地址赋值给BP
        0x00aa 00170 (main.go:11)       ADDQ    $24, SP                        # 回收A函数栈帧空间
        0x00ae 00174 (main.go:11)       RET
```

在A函数栈中，我们可以看到，程序先把r1~r9参数分别从寄存器赋值到main栈帧的入参地址部分，即当前的SP+48~SP+112位，其实跟GO1.15.14的函数调用参数传递过程差不多，只不过一个是在caller中做参数从寄存器拷贝到栈上，一个是在callee中做参数从寄存器拷贝到栈上，而且前者只使用了AX一个寄存器，后者使用了9个不同的寄存器。

![img](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/sy1cud3w67.png)

禁用内联优化：

go1.15:

```go
GOOS=linux GOARCH=amd64 go tool compile -S main.go 
```

```go
"".A STEXT nosplit size=59 args=0x40 locals=0x0
        0x0000 00000 (main.go:8)        TEXT    "".A(SB), NOSPLIT|ABIInternal, $0-64
        0x0000 00000 (main.go:8)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:8)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:9)        MOVQ    "".p1+8(SP), AX     #参数从栈赋值到寄存器AX
        0x0005 00005 (main.go:9)        MOVQ    "".p2+16(SP), CX    #参数从栈赋值到寄存器CX
        0x000a 00010 (main.go:9)        ADDQ    CX, AX              
        0x000d 00013 (main.go:9)        MOVQ    "".p3+24(SP), CX    #参数从栈赋值到寄存器CX
        0x0012 00018 (main.go:9)        ADDQ    CX, AX
        0x0015 00021 (main.go:9)        MOVQ    "".p4+32(SP), CX
        0x001a 00026 (main.go:9)        ADDQ    CX, AX
        0x001d 00029 (main.go:9)        MOVQ    "".p5+40(SP), CX
        0x0022 00034 (main.go:9)        ADDQ    CX, AX
        0x0025 00037 (main.go:9)        MOVQ    "".p6+48(SP), CX
        0x002a 00042 (main.go:9)        ADDQ    CX, AX
        0x002d 00045 (main.go:9)        MOVQ    "".p7+56(SP), CX
        0x0032 00050 (main.go:9)        ADDQ    CX, AX
        0x0035 00053 (main.go:9)        MOVQ    AX, "".~r7+64(SP)
        0x003a 00058 (main.go:9)        RET
```

go1.17：

```go
"".A STEXT nosplit size=21 args=0x38 locals=0x0 funcid=0x0
        0x0000 00000 (main.go:8)        TEXT    "".A(SB), NOSPLIT|ABIInternal, $0-56
        0x0000 00000 (main.go:8)        FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:8)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0000 00000 (main.go:8)        FUNCDATA        $5, "".A.arginfo1(SB)
        0x0000 00000 (main.go:9)        LEAQ    (BX)(AX*1), DX
        0x0004 00004 (main.go:9)        ADDQ    DX, CX            #直接在寄存器之间做加法
        0x0007 00007 (main.go:9)        ADDQ    DI, CX            #直接在寄存器之间做加法
        0x000a 00010 (main.go:9)        ADDQ    SI, CX
        0x000d 00013 (main.go:9)        ADDQ    R8, CX
        0x0010 00016 (main.go:9)        LEAQ    (R9)(CX*1), AX
        0x0014 00020 (main.go:9)        RET
```



性能测试：

```go
package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var sl []string = []string{
	"Rob Pike ",
	"Robert Griesemer ",
	"Ken Thompson ",
}

func concatStringByOperator(sl []string) string {
	var s string
	for _, v := range sl {
		s += v
	}
	return s
}

func concatStringBySprintf(sl []string) string {
	var s string
	for _, v := range sl {
		s = fmt.Sprintf("%s%s", s, v)
	}
	return s
}

func concatStringByJoin(sl []string) string {
	return strings.Join(sl, "")
}

func concatStringByStringsBuilder(sl []string) string {
	var b strings.Builder
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByStringsBuilderWithInitSize(sl []string) string {
	var b strings.Builder
	b.Grow(64)
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByBytesBuffer(sl []string) string {
	var b bytes.Buffer
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByBytesBufferWithInitSize(sl []string) string {
	buf := make([]byte, 0, 64)
	b := bytes.NewBuffer(buf)
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func BenchmarkConcatStringByOperator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByOperator(sl)
	}
}

func BenchmarkConcatStringBySprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringBySprintf(sl)
	}
}

func BenchmarkConcatStringByJoin(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByJoin(sl)
	}
}

func BenchmarkConcatStringByStringsBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByStringsBuilder(sl)
	}
}

func BenchmarkConcatStringByStringsBuilderWithInitSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByStringsBuilderWithInitSize(sl)
	}
}

func BenchmarkConcatStringByBytesBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByBytesBuffer(sl)
	}
}

func BenchmarkConcatStringByBytesBufferWithInitSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByBytesBufferWithInitSize(sl)
	}
}

```

Go 1.16:

![image-20240922092129789](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240922092129789.png)

Go 1.17:

![image-20240922092206866](https://imghosting-1257040086.cos.ap-nanjing.myqcloud.com/img/image-20240922092206866.png)

相对于 Go 1.16 跑出的结果，Go 1.17 在每一个测试项上都有小幅的性能提升，有些性能提升甚至达到 10% 左右（以 BenchmarkConcatStringBySprintf 为例，它的性能提升为 (445.1-410.1)/445.1=7.8%）。也就是说你的 Go 源码使用 Go 1.17 版本重新编译一下，就能获得大约 5% 的性能提升

Go 1.17 发布说明也提到了：**改为基于寄存器的调用惯例后，绝大多数程序不会受到影响**。只有那些之前就已经违反 unsafe.Pointer 的使用规则的代码可能会受到影响



https://go.dev/blog/go1.17

https://go.dev/doc/go1.17  Go 1.17 的发布说明文档

https://github.com/golang/go/tree/release-branch.go1.17

https://github.com/golang/go/issues Go 语言项目的官方 issue 列表

https://go-review.googlesource.com/q/status:open+-is:wip Go 项目的代码 review 站点

https://www.bilibili.com/video/BV1iS4y1z7V5/?spm_id_from=333.337.search-card.all.click&vd_source=39d622d5b0f294cef4556e52ef149a30 从汇编角度理解函数调用过程

https://cloud.tencent.com/developer/article/2126557 深入分析go1.17函数调用栈参数传递

https://blog.csdn.net/weixin_52690231/article/details/125305807 调用约定修改
