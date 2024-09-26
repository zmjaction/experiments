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

### 1、泛型

**泛型是Go诞生以来最复杂、最难读和理解的语法特性**

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

https://tonybai.com/2022/04/20/some-changes-in-go-1-18/ go1.18变化

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

## go 1.16

https://go.dev/doc/go1.16  Go 1.16 的发布说明文档