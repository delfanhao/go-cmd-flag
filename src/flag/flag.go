package flag

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func init() {
	ctx = CommandLineContext{
		startIndex: 1,
		preParams:  make(map[string]string, 0),
		target:     "",
		sufParams:  make(map[string]string, 0),
	}
}

// Parse 解析用户定义的参数项
func Parse(define interface{}) {
	ParseFromPosition(define, ctx.startIndex)
}

// ParseFromPosition 从指定位置开始解析对应的命令行参数
// define 用户定义的参数, 要求是 struct 结构
// startIndex 参数开始的位置
func ParseFromPosition(define interface{}, startIndex int) {
	ctx.startIndex = startIndex
	parseCmdLine()
	parseDefine(define)
}

// SetAppName 设置应用名称
func SetAppName(name string) {
	ctx.appName = name
}

// GetModule 获取模块名
func GetModule() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	return ""
}

// GetTarget 获取执行对象
func GetTarget() string {
	return ctx.target
}

// parseDefine 解析用户定义的结构，获取需要参数定义的flag，缺省值，类型等信息
func parseDefine(define interface{}) {
	defineType := reflect.TypeOf(define)
	if defineType.Kind() != reflect.Ptr && defineType.Kind() != reflect.Struct {
		panic("Need struct type, Check your code please.")
	}

	defineValue := reflect.ValueOf(define).Elem()
	numEle := defineType.Elem().NumField()
	for idx := 0; idx < numEle; idx++ {
		field := defineType.Elem().Field(idx)
		helpMsg := DefineStruct{
			flag:        "",
			full:        "",
			description: "",
			defaultVal:  "",
		}
		if tag, ok := field.Tag.Lookup(TagFlag); ok {
			helpMsg.flag = tag
			full := ""
			if full, ok = field.Tag.Lookup(TagFullName); ok {
				helpMsg.full = full
			}

			if value, ok := findOnCmdLine(tag, full); ok {
				setter := defineValue.FieldByName(field.Name)
				setter.SetString(value)
			} else {
				if req, ok := field.Tag.Lookup(TagRequired); ok {
					if strings.ToUpper(req) == ValueTrue {
						panic(fmt.Sprintf("The param %s must specified.", helpMsg.flag ))
					}
				}
				// 如果用户定义了缺省值，则设置为缺省值
				if defaultValueTag, ok := field.Tag.Lookup(TagDefault); ok {
					setter := defineValue.FieldByName(field.Name)
					setter.SetString(defaultValueTag)
					helpMsg.defaultVal = defaultValueTag
				}
			}

			if desc, ok := field.Tag.Lookup("desc"); ok {
				helpMsg.description = desc
			}

			helpInfo = append(helpInfo, helpMsg)
		}
	}

	// 处理结束后，如果ctx.preParams 中还有元素， 说明是无法解析的参数，报错
	if len(ctx.preParams) > 0 {
		showUnknownParams()
	}
}

// showUnknownParams 显示无法解析的命令行参数
func showUnknownParams() {
	for k, _ := range ctx.preParams {
		if k == "h" || k == "help" {
			showHelp()
		}
		fmt.Println("Unknown param :", k)
	}
	os.Exit(-1)
}
