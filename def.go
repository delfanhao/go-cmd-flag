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
	TagArg              = "arg"     // 是否后置参数arg

	ValueTrue  = "T" // true
	ValueFalse = "F" // false
	ValueYes   = "Y" // yes
	ValueNo    = "N" // no
	Value1     = "1" // 1, true
	Value0     = "0" // 0, false

	InvalidParam = "Invalid param -%s" // 无效参数
)

// CommandLineContext 上下文结构
type CommandLineContext struct {
	appName    string
	startIndex int
	modules    []string
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
	arg            bool
	valDescription string
}

var (
	ctx         CommandLineContext
	preHelpInfo []DefineStruct
	sufHelpInfo []DefineStruct
)

func init() {
	preHelpInfo = make([]DefineStruct, 0)
	sufHelpInfo = make([]DefineStruct, 0)
}

func formatInfo(s string) string {
	if l := len(s); l > 0 {
		s = " " + s
	}
	return s
}

func showParams(title string, helpInfo []DefineStruct) {
	if len(helpInfo) == 0 {
		return
	}
	fmt.Println(title)
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
	preHelpInfo = append(preHelpInfo, help)
	mods := ""
	for i := range ctx.modules {
		if len(mods) > 0 {
			mods += "|"
		}
		mods += ctx.modules[i]
	}

	mods = formatInfo(mods)
	target := formatInfo(ctx.target)
	sufParams := ""
	if len(sufHelpInfo) > 0 {
		sufParams = " TARGET [ARGS...]"
	}

	usageInfo := fmt.Sprintf("Usage: %s%s [OPTIONS]%s%s", ctx.appName, mods, target, sufParams)
	fmt.Println(usageInfo)

	showParams("Options:", preHelpInfo)
	showParams("Args:", sufHelpInfo)
}
