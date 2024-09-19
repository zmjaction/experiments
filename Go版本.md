# Go版本

## go 1.23

https://go.dev/blog/go1.23

https://go.dev/doc/go1.23 Go 1.23 的发布说明文档

## go 1.22

https://go.dev/doc/go1.22 Go 1.22 的发布说明文档

## go 1.21

https://go.dev/doc/go1.21 Go 1.21 的发布说明文档

## go 1.20

https://go.dev/doc/go1.20 Go 1.20 的发布说明文档

## go 1.19

https://go.dev/doc/go1.19 Go 1.19 的发布说明文档

## go 1.18

https://go.dev/doc/go1.18 Go 1.18 的发布说明文档

## go 1.17

### 语法特性

#### 支持将切片转换为数组指针

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

### Go Module 构建模式的变化

#### 修剪的 module 依赖图

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

### go get 不再被用来安装命令可执行文件

新版本中，我们需要使用 go install 来安装，并且使用 go install 安装时还要用 @vx.y.z 明确要安装的命令的二进制文件的版本，或者是使用 @latest 来安装最新版本

### //go:build 形式的构建约束指示符

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

### 基于寄存器的调用惯例

Go 1.17 版本中，Go 编译器最大的变化是在 AMD64 架构下率先实现了**从基于堆栈的调用惯例到基于寄存器的调用惯例的切换**

所谓“调用惯例（calling convention）”，是指调用方和被调用方对于函数调用的一个明确的约定，包括函数参数与返回值的传递方式、传递顺序。只有双方都遵守同样的约定，函数才能被正确地调用和执行。如果不遵守这个约定，函数将无法正确执行。

Go 1.17 版本之前，Go 采用基于栈的调用约定，也就是说函数的参数与返回值都通过栈来传递，这种方式的优点是实现简单，不用担心底层 CPU 架构寄存器的差异，适合跨平台，但缺点就是牺牲了一些性能。

寄存器的访问速度是要远高于内存的，go1.17 在 AMD64 架构下率先实现基于寄存器的调用惯例

性能提升，**Go 程序的运行性能提高了约 5%，二进制文件大小减少约 2%**

```go
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

Go 1.16.5:

```
$go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/demo
cpu: Intel(R) Core(TM) i5-8257U CPU @ 1.40GHz
BenchmarkConcatStringByOperator-8                       12132355          91.51 ns/op
BenchmarkConcatStringBySprintf-8                         2707862         445.1 ns/op
BenchmarkConcatStringByJoin-8                           24101215          50.84 ns/op
BenchmarkConcatStringByStringsBuilder-8                 11104750         124.4 ns/op
BenchmarkConcatStringByStringsBuilderWithInitSize-8     24542085          48.24 ns/op
BenchmarkConcatStringByBytesBuffer-8                    14425054          77.73 ns/op
BenchmarkConcatStringByBytesBufferWithInitSize-8        20863174          49.07 ns/op
PASS
ok    github.com/demo  9.166s
```

Go 1.17:

```go
$go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/demo
cpu: Intel(R) Core(TM) i5-8257U CPU @ 1.40GHz
BenchmarkConcatStringByOperator-8                       13058850          89.47 ns/op
BenchmarkConcatStringBySprintf-8                         2889898         410.1 ns/op
BenchmarkConcatStringByJoin-8                           25469310          47.15 ns/op
BenchmarkConcatStringByStringsBuilder-8                 13064298          92.33 ns/op
BenchmarkConcatStringByStringsBuilderWithInitSize-8     29780911          41.14 ns/op
BenchmarkConcatStringByBytesBuffer-8                    16900072          70.28 ns/op
BenchmarkConcatStringByBytesBufferWithInitSize-8        27310650          43.96 ns/op
PASS
ok    github.com/demo  9.198s
```

相对于 Go 1.16.5 跑出的结果，Go 1.17 在每一个测试项上都有小幅的性能提升，有些性能提升甚至达到 10% 左右（以 BenchmarkConcatStringBySprintf 为例，它的性能提升为 (445.1-410.1)/445.1=7.8%）。也就是说你的 Go 源码使用 Go 1.17 版本重新编译一下，就能获得大约 5% 的性能提升

Go 1.17 发布说明也提到了：**改为基于寄存器的调用惯例后，绝大多数程序不会受到影响**。只有那些之前就已经违反 unsafe.Pointer 的使用规则的代码可能会受到影响



https://go.dev/blog/go1.17

https://go.dev/doc/go1.17  Go 1.17 的发布说明文档

https://github.com/golang/go/tree/release-branch.go1.17

https://github.com/golang/go/issues Go 语言项目的官方 issue 列表

https://go-review.googlesource.com/q/status:open+-is:wip Go 项目的代码 review 站点

## go 1.16

https://go.dev/doc/go1.16  Go 1.16 的发布说明文档