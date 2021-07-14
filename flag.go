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

//
func newDefineStruct() DefineStruct {
	return DefineStruct{
		flag:           "",
		full:           "",
		description:    "",
		defaultVal:     "",
		required:       false,
		valDescription: "",
	}
}

//
func processFlagAndFull(field *reflect.StructField, helpMsg *DefineStruct) (string, string, bool) {
	tag, full, ok := "", "", false
	tag, ok = field.Tag.Lookup(TagFlag)
	if ok {
		helpMsg.flag = tag
		if full, ok = field.Tag.Lookup(TagFullName); ok {
			helpMsg.full = full
		}
	}

	return tag, full, ok
}

//
func setValue(tag, full string, field *reflect.StructField, defineValue *reflect.Value, helpMsg *DefineStruct) {
	value, ok := findOnCmdLine(tag, full)
	if ok {
		setter := defineValue.FieldByName(field.Name)
		setter.SetString(value)
	} else {
		if defaultValueTag, ok := field.Tag.Lookup(TagDefault); ok {
			setter := defineValue.FieldByName(field.Name)
			setter.SetString(defaultValueTag)
			helpMsg.defaultVal = defaultValueTag
		}
	}
}

// checkRequired 检查是否必须的参数
func checkRequired(helpMsg *DefineStruct, field *reflect.StructField) {
	req, ok := field.Tag.Lookup(TagRequired)
	if ok && strings.ToUpper(req) == ValueTrue || strings.ToUpper(req) == ValueYes {
		helpMsg.required = true
	}
}

func setDesc(field *reflect.StructField, helpMsg *DefineStruct) {
	if desc, ok := field.Tag.Lookup(TagDescription); ok {
		helpMsg.description = desc
	}

	if valDesc, ok := field.Tag.Lookup(TagValueDescription); ok {
		helpMsg.valDescription = valDesc
	}
}

// parseDefine 解析用户定义的结构，获取需要参数定义的flag，缺省值，类型等信息
func parseDefine(define interface{}) {
	defineType := reflect.TypeOf(define)
	if defineType.Kind() != reflect.Ptr && defineType.Kind() != reflect.Struct {
		panic("Need struct type, Check your code please.")
	}

	defineValue := reflect.ValueOf(define).Elem()
	numEle := defineType.Elem().NumField()
	hasRequired := false
	for idx := 0; idx < numEle; idx++ {
		field := defineType.Elem().Field(idx)
		helpMsg := newDefineStruct()

		if tag, full, ok := processFlagAndFull(&field, &helpMsg); ok {
			helpMsg.flag = tag
			checkRequired(&helpMsg, &field)
			hasRequired = hasRequired || helpMsg.required
			setValue(tag, full, &field, &defineValue, &helpMsg)
			setDesc(&field, &helpMsg)
			helpInfo = append(helpInfo, helpMsg)
		}
	}

	// 定义扫描完成，如果lossReqCount为真，说明命令行中没有输入参数
	if hasRequired && len(os.Args) == 1 {
		Usage()
		os.Exit(-1)
	}

	// 处理结束后，如果ctx.preParams 中还有元素， 说明是无法解析的参数，报错
	if len(ctx.preParams) > 0 {
		showUnknownParams()
	}
}

// showUnknownParams 显示无法解析的命令行参数
func showUnknownParams() {
	showHelp()
	for k, _ := range ctx.preParams {
		fmt.Println("Unknown param :", k)
	}
	os.Exit(-1)
}
