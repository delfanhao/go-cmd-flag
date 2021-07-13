package flag

import (
	"fmt"
	"os"
	"path"
)

const (
	TagFlag      = "flag"     // 参数名, 参数短名，要求只能有一个字母, 大小写敏感
	TagFullName  = "full"     // 长参数, 用两个连续减号作为前导
	TagKind      = "kind"     // 类别，可设置为 string/int/float/bool, 缺省为string
	TagDefault   = "default"  // 缺省值，如果从命令行获取不到对应的值，则使用此值
	TagModule    = "module"   // 命令
	TagRequired  = "required" // 是否必须的参数
	DefaultStart = 1          // 参数开始位置

	ValueTrue  = "T" // bool常量的True
	ValueFalse = "F" // bool 常量的false

	InvalidParam = "Invalid param -%s" // 无效参数
)

// CommandLineContext 命令行参数结构, 对于不带-的
type CommandLineContext struct {
	appName    string
	startIndex int
	tokens     []string
	preParams  map[string]string
	target     string
	sufParams  map[string]string
}

// DefineStruct 定义
type DefineStruct struct {
	flag        string
	full        string
	description string
	defaultVal  string
	required    string
}

var (
	ctx      CommandLineContext
	helpInfo []DefineStruct
)

func init() {
	helpInfo = make([]DefineStruct, 0)
}

//Usage 显示帮助信息
var Usage = func() {
	if len(ctx.appName) == 0 {
		_, file := path.Split(os.Args[0])
		ctx.appName = file
	}

	fmt.Println("Usage: " + ctx.appName + " [OPTIONS] TARGET [ARG...]\n")

	fmt.Println("Options:")
	for i := 0; i < len(helpInfo); i++ {
		line := fmt.Sprintf("-%s", helpInfo[i].flag)
		if len(helpInfo[i].full) > 0 {
			line += fmt.Sprintf(",--%s", helpInfo[i].full)
		}
		if len(helpInfo[i].description) > 0 {
			line += fmt.Sprintf("\n    %s", helpInfo[i].description)
		}
		if len(helpInfo[i].defaultVal) > 0 {
			line += fmt.Sprintf("\n    default value is %s", helpInfo[i].defaultVal)
		}
		println(line)
	}
}
