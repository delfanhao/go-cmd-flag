package flag

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	ctx = CommandLineContext{
		appName:    "",
		modules:    make([]string, 0),
		startIndex: 1,
		preParams:  make(map[string]string, 0),
		target:     "",
		sufParams:  make(map[string]string, 0),
	}
}

// Parse 解析用户定义的参数项
func Parse(define interface{}) {
	if len(ctx.modules) > 0 {
		ctx.startIndex = 2
	}

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

// SetModules 设置模块名
func SetModules(mod ...string) {
	ctx.modules = mod
}

func newDefineStruct() DefineStruct {
	return DefineStruct{
		flag:           "",
		full:           "",
		description:    "",
		defaultVal:     "",
		required:       false,
		arg:            false,
		valDescription: "",
	}
}

//
func processFlagAndFull(field *reflect.StructField, helpMsg *DefineStruct) (string, string, bool) {
	tag, full, ok, fullOk := "", "", false, false
	tag, ok = field.Tag.Lookup(TagFlag)
	if ok {
		helpMsg.flag = tag
		if full, fullOk = field.Tag.Lookup(TagFullName); fullOk {
			helpMsg.full = full
		}
	}

	return tag, full, ok
}

// setValueForType 根据用户定义的类型进行设置值
func setValueForType(defineValue *reflect.Value, field *reflect.StructField, value string) bool {
	ok := true
	setter := defineValue.FieldByName(field.Name)
	switch setter.Kind() {
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			setter.SetFloat(v)
		} else {
			ok = false
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			setter.SetInt(v)
		} else {
			ok = false
		}

	case reflect.String:
		setter.SetString(value)

	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err == nil {
			setter.SetBool(v)
		} else {
			ok = false
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			setter.SetUint(v)
		} else {
			ok = false
		}

	default:
		ok = false
	}

	return ok
}

//
func setValue(tag, full string, field *reflect.StructField, defineValue *reflect.Value, helpMsg *DefineStruct) {
	value, ok := findOnCmdLine(tag, full, helpMsg.arg)
	if ok {
		if !setValueForType(defineValue, field, value) {
			fmt.Println(fmt.Sprintf("Invalid value of param and value : %s=%s", tag, value))
			os.Exit(-1)
		}
	} else {
		if defaultValueTag, ok := field.Tag.Lookup(TagDefault); ok {
			if !setValueForType(defineValue, field, defaultValueTag) {
				panic("Struct definition error.")
			}
			helpMsg.defaultVal = defaultValueTag
		}
	}
}

// checkRequired 检查是否必须的参数
func checkRequired(helpMsg *DefineStruct, field *reflect.StructField) {
	req, ok := field.Tag.Lookup(TagRequired)
	helpMsg.required = ok &&
		strings.ToUpper(req) == ValueTrue ||
		strings.ToUpper(req) == ValueYes ||
		req == Value1
}

// isArg 判断是否后置参数
func isArg(field *reflect.StructField) bool {
	req, ok := field.Tag.Lookup(TagArg)
	return ok &&
		strings.ToUpper(req) == ValueTrue ||
		strings.ToUpper(req) == ValueYes ||
		req == Value1
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
			setDesc(&field, &helpMsg)
			helpMsg.arg = isArg(&field)
			setValue(tag, full, &field, &defineValue, &helpMsg)
			if !helpMsg.arg {
				preHelpInfo = append(preHelpInfo, helpMsg)
			} else {
				sufHelpInfo = append(sufHelpInfo, helpMsg)
			}
		}
	}

	// 定义扫描完成，如果lossReqCount为真，说明命令行中没有输入参数
	if hasRequired && len(os.Args) == 1 {
		Usage()
		os.Exit(-1)
	}

	// 处理结束后，如果ctx.preParams 中还有元素， 说明是无法解析的参数，报错
	if len(ctx.preParams) > 0 || len(ctx.sufParams) > 0 {
		showUnknownParams()
	}
}

// showUnknownParams 显示无法解析的命令行参数
func showUnknownParams() {
	showHelp()
	for p, _ := range ctx.preParams {
		fmt.Println("Unknown option :", p)
	}

	for s, _ := range ctx.sufParams {
		fmt.Println("Unknown arg : ", s)
	}

	os.Exit(-1)
}
