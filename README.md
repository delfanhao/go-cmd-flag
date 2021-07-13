GO-CMD-FLAG

> 本package是一个golang解析命令行参数的包，用户在业务代码中简单定义一个struct，通过元数据的描述即可通过本package完成参数、目标的解析工作。


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

