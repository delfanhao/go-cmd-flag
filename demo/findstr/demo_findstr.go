package main

import (
	"fmt"
	flag "github.com/delfanhao/go-cmd-flag"
)

/*
   程序运行时请加入参数:
     -o out.log test.txt --str '123 45' -r 10
   进行测试
*/

type Params struct {
	Output string `flag:"o" full:"output" desc:"Specified output filename."`
	Repeat int64  `flag:"r" full:"repeat" desc:"Repeat times." arg:"y"`
	Str    string `flag:"s" full:"str" desc:"Find string." arg:"y"`
}

func showValue(k string, v interface{}) {
	println(fmt.Sprintf("%s=%v", k, v))
}

func main() {
	p := Params{}
	flag.SetAppName("findstr")
	flag.Parse(&p)
	showValue("Output", p.Output)
	showValue("Target", flag.GetTarget())
	showValue("Str", p.Str)
	showValue("Repeat", p.Repeat)
}
