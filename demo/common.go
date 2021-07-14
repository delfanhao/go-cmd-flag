package main

// CmdLineDefine 定义命令行样例
type CmdLineDefine struct {
	Output  string `flag:"o" full:"output" desc:"Specified output filename" req:"y" valDesc:"filename"`
	Mode    string `flag:"m" desc:"Specified display mode." default:"text"`
	LineNum string `flag:"l" desc:"Show number of line."`
}
