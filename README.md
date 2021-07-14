# GO-CMD-FLAG

> go-cmd-flag 是一个让您能够在golang应用中便捷添加命令行参数、程序参数的工具包，相比 golang 本身提供的 flag 包更加简单便捷，并且增加了程序参数的功能。

#### 概念

我们在开发命令行工具的时候， 通常需要设计一些参数来达到灵活使用的目的。这部分功能在Golang的标准库flag中，已经有了基本的功能， 但使用起来略烦琐，且flag包解析命令行的时候， 遇到第一个不带"-"的参数则解析结束。
go-cmd-flag则是为了解决这些问题而开发。

在go-cmd-flag里，把命令行所有内容分为如下几个段落：

``命令 [模块名] [前置参数...] [目标] [后置参数...]``

- 命令： 程序本身；
- 模块名： 对应程序内部命令模块，例如我们熟悉的Docker, 通常后面后跟着一个命令用来指明要使用什么功能，例如 ``docker run``；
- 前置参数： 这部分参数列表通常都出现在 ``目标`` 前，用来设置程序初始的值；
- 目标： 程序需要执行的对象，例如一个对文件操作的命令，目标就是指定的文件名；
- 后置参数： 放在目标后的参数列表，通常用于指明程序运行时的设置；

例如，我们编写了一个统计文本文件中出现指定字符串次数的程序， 需要把找到的文职记录下，且当找到30个的时候就结束，则使用时可能是如下命令：

``findstr --output result.log --str=tom abc.txt --repeat 30``

其中， ``findstr``为命令， ``--oupput``,``--str`` 为前置参数， ``abc.txt`` 为目标， ``--repeat``为后置参数。

## 安装/更新

```
go get -u github.com/delfan/go-cmd-flag
```

## 使用

```
package main



```




例, 业务代码中有如下定义：

```
type struct FlagSample {
   Output string `flag:"o" full:"output" desc:"Specified output filename."`
}
```
引入本package并调用解析参数的方法，即可在用户定义的struct中填充解析后的值：
```
import (
    flag "github.com/delfanhao/go-cmd-flag"
)

func main() {
	p := CmdLineParam{}
	flag.SetAppName("Go-cmd-flag")
	flag.ParseFromPosition(&p, 1)
	target := flag.GetTarget()   
	println(flag.Output)
}    
```
编译后测试:
```
./sample -o /golang/test.txt
```

## 内置前置参数

我们在使用命令行的时候，都会通过使用 -h 或者 --help 来获取应用本身的帮助信息， go-cmd-flag 内置了这部分内容的处理工作，通过您自定义的struct中进行组织并输出。

同样与标准flag包类似，go-cmd-flag 中的Usage函数也被定义为一个闭包， 如果您对自动生成的帮助信息不满意， 可以通过修改Usage的定义来定制话您需要的内容。

例如:``sample -o test.txt -h``，程序不会有其它动作，只是显示帮助信息然后退出。


