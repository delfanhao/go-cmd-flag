package flag

import (
	"fmt"
	"os"
	"path"
)

const (
	TagFlag             = "flag"    // 参数名, 参数短名，要求只能有一个字母, 大小写敏感
	TagFullName         = "full"    // 长参数, 用两个连续减号作为前导
	TagDefault          = "default" // 缺省值，如果从命令行获取不到对应的值，则使用此值
	TagDescription      = "desc"    // 参数描述信息
	TagRequired         = "req"     // 是否必须的参数
	TagValueDescription = "valDesc" // 参数值的描述内容

	ValueTrue  = "T" // true
	ValueFalse = "F" // false
	ValueYes   = "Y" // yes
	ValueNo    = "N" // no

	InvalidParam = "Invalid param -%s" // 无效参数
)

// CommandLineContext 命令行参数结构, 对于不带-的
type CommandLineContext struct {
	appName    string
	startIndex int
	preParams  map[string]string
	target     string
	sufParams  map[string]string
}

// DefineStruct 定义tag结构
type DefineStruct struct {
	flag           string
	full           string
	description    string
	defaultVal     string
	required       bool
	valDescription string
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

	// 添加help帮助信息
	help := DefineStruct{
		flag:           "h",
		full:           "help",
		description:    "Show this help information.",
		defaultVal:     "",
		required:       false,
		valDescription: "",
	}
	helpInfo = append(helpInfo, help)

	fmt.Println("Usage: " + ctx.appName + " [OPTIONS] TARGET [ARGS...]")

	fmt.Println("Options:")
	for i := 0; i < len(helpInfo); i++ {
		line := fmt.Sprintf("-%s", helpInfo[i].flag)
		if len(helpInfo[i].full) > 0 {
			line += fmt.Sprintf(",--%s", helpInfo[i].full)
		}
		if len(helpInfo[i].valDescription) > 0 {
			line += fmt.Sprintf(" %s", helpInfo[i].valDescription)
		}

		req := ""
		if helpInfo[i].required {
			req = "Required. "
		}

		if len(helpInfo[i].description) > 0 {
			line += fmt.Sprintf("\n    %s%s", req, helpInfo[i].description)
		}
		if len(helpInfo[i].defaultVal) > 0 {
			line += fmt.Sprintf("\n    default value is %s", helpInfo[i].defaultVal)
		}
		fmt.Println(line)
	}
}
