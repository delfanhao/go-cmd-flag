# GO-CMD-FLAG

> go-cmd-flag 是一个让您能够在Golang应用中便捷添加命令行参数(Options)、程序参数(Args)的工具包，相比Golang本身提供的flag包更加简单便捷。

## 为什么要开发GO-CMD-FLAG工具包

我们在使用Golang开发命令行工具的时候，通常需要设计一些参数来达到灵活使用的目的。这部分功能在Golang的标准库flag中，已经有了基本的功能，但使用起来略烦琐：

- 不支持子命令，诸如``docker run``这中形式的命令；
- 每一个参数都需要通过语句进行定义，参数多的时候，需要写很多重复代码；
- 不区分``-``及``--``，弱化了参数缩写及全名的概念，不太符合习惯；  
- 解析命令行的时候， 遇到第一个不带``-``或``--``的参数则解析结束，不处理程序参数;

``GO-CMD-FLAG``则是为了解决这些问题而产生，着力于解决以上问题。

## GO-CMD-FLAG的能力

首先总结一下常见的命令行组成，常用的命令行都会有以下几部分：

``命令 [子命令] [前置参数...] [目标] [后置参数...]``

这几个部分的解释如下：

- 命令：   
  程序本身；
  
- 子命令：   
  对应程序内部命令模块，例如开发人员熟悉的``docker run``, 其中``run``就是这种形式的子命令；
  
- 命令行参数(OPTIONS)：
  前置参数，这部分参数列表通常都出现在 ``目标`` 前面，用来设置程序初始的值等相关内容；
  
- 目标：   
  程序需要执行的对象，例如一个对文件操作的命令，目标就是指定的文件名；
  
- 程序参数(ARGS)：   
  后置参数，放在目标后的参数列表，通常用于指明程序运行时的设置；
  
例如，我们编写了一个统计文本文件中出现指定字符串次数的程序， 需要把找到的位置记录在日志中，且当找到30个的时候就结束，则使用时可能是如下命令：

``findstr --output result.log abc.txt --str=tom --repeat 30``

其中， ``findstr``为命令， ``--output result.log``为命令行参数， ``abc.txt`` 为目标， ``--str=tom --repeat 30``为程序参数。

``GO-CMD-FLAG``能够很方便的处理以上这些情况。

## 安装/更新
```
go get -u github.com/delfanhao/go-cmd-flag
```

## 使用

例如需要开发上述 findstr 的功能，参数定义如下：

```
type Params struct {
   Output string `flag:"o" full:"output" desc:"Specified output filename."`
   Repeat int64 `flag:"r" full:"repeat" desc:"Repeat times." arg:"y"`
   Str string `flag:"s" full:"str" desc:"Find string." arg:"y"`
}
```

> 注意：由于Golang的特性，需要导出的参数均需要大写，因此定义struct及参数变量时，需符合这一个语法要求；

引入本package并调用解析参数的方法，即可在用户定义的struct中填充解析后的值：
```
import (
    flag "github.com/delfanhao/go-cmd-flag"
)

type Params struct { ... } // 定义略，参见上面Params定义

func main() {
	p := Params{}
	flag.SetAppName("findstr")
	flag.ParseFromPosition(&p, 1)
	target := flag.GetTarget()   
	println(p.Output)
	println(p.Str)
	println(p.Repeat)
}    
```

### 定义参数时使用的结构体标签
用户定义参数结构体时，根据需要定义出参数变量及类型，并在结构体标签中定义参数行为，可用的标签如下：
> 注: bool类型可用t/y/1表示true, f/n/0表示false,不区分大小写

| 标签名 | 描述  | 缺省值 |
|----|----|-----|
|flag| 必要标签，用于定义单字母的参数名，用于匹配命令行中``-``符号前导的参数||
|full| 可选标签，用于定义参数全名，用于匹配命令行中``--``符号前导的参数||
|default |可选标签，设置当前参数的缺省值||
|desc|可选标签，用于帮助信息中说明此参数的作用||
|req|可选标签，声明当前参数是否用户必须指定|n|
|valDesc|可选标签，参数值描述||
|arg|可选标签，声明当前参数是命令行参数还是程序参数|n|

### API

- SetAppName(string)  
  设置应用名称，用于帮助信息中的显示，如未设置则使用程序名本身；

- SetModules(...string)  
  设置子命令，参数为可变长参数；

- Parse(struct)   
  指定自定义的参数结构体并开始解析命令行参数，这里需要使用指针，即``&结构体变量``形式；

- ParseFromPosition(struct, int)  
  作用同Parse函数，不同的是可指定从命令行的第几项开始解析；
  
- GetModule() string  
  获取命令行中指定的子命令;
  
- GetTarget() string  
  获取命令行中的目标参数;
  

## 内置参数

我们在使用命令行的时候，都会通过使用 ``-h`` 或者 ``--help`` 来获取应用本身的帮助信息， go-cmd-flag 内置了这部分内容的处理工作，通过您自定义的struct中进行组织并输出。

同样与标准flag包类似，go-cmd-flag 中的Usage函数也被定义为一个闭包， 如果您对自动生成的帮助信息不满意， 可以通过修改Usage的定义来定制您需要的内容。

例如:``findstr -h -o test.txt``，程序不会有其它动作，只是显示帮助信息然后退出。




