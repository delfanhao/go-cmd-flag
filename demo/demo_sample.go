package main

import flag "github.com/delfanhao/go-cmd-flag"

func main() {
	p := CmdLineDefine{}
	flag.Parse(&p)
	println("param output = ", p.Output)
	println("param intValue = ", p.IntValue)
}
