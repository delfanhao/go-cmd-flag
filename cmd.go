package flag

import (
	"fmt"
	"os"
	"strings"
)

// findOnCmdLine 在命令行中获取对应的参数,如果存在则返回对应的值及设置ok为true
func findOnCmdLine(tag string, full string) (string, bool) {
	if v, ok := ctx.preParams[tag]; ok {
		delete(ctx.preParams, tag)
		delete(ctx.preParams, full)
		return v, true
	} else if v, ok := ctx.preParams[full]; ok {
		delete(ctx.preParams, full)
		return v, true
	}

	return "", false
}

// split 根据指定分隔符，返回 k,v对，如果不存在分隔符， 则v为空字符串
func split(token string, sep string) (k, v string) {
	if strings.Index(token, sep) >= 1 {
		arr := strings.Split(token, "=")
		if len(arr) < 2 {
			panic("error")
		}

		return arr[0], arr[1]
	}

	return token, ""
}

// processToken 处理当前位置的token，如果是参数名，则向后读取一个，如果是value，则放入map, 否则认为就是个标记
func processToken(token string, pos int) (int, string, string) {
	isSingleFlag := true
	if token[0] == '-' {
		token = token[1:]
		isSingleFlag = false
	}
	step, k, v := 1, "", ""
	if strings.Index(token, "=") >= 1 {
		k, v = split(token, "=")
		step = 1
	} else if strings.Index(token, ":") >= 1 {
		k, v = split(token, ":")
		step = 1
	} else {
		if len(os.Args) > 1 && pos <= len(os.Args)-1 && len(os.Args) > pos+1 && os.Args[pos+1][0] != '-' {
			step, k, v = 2, token, os.Args[pos+1]
		} else {
			step, k, v = 1, token, ""
		}
	}

	if isSingleFlag && len(k) != 1 {
		panic(fmt.Sprintf(InvalidParam, k))
	}

	return step, k, v
}

// parseCmdLine 分解当前命令行，解析出参数, 放入context的map中
func parseCmdLine() {
	pos := ctx.startIndex

	state := -1 // 设置为前导参数解析阶段， 遇到第一个不在token后的且不带 - 的token，则认为是target
	for pos < len(os.Args) {
		token := os.Args[pos]
		if token[0] == '-' {
			token = token[1:]
			step, k, v := processToken(token, pos)
			if state == -1 {
				ctx.preParams[k] = v
			} else {
				ctx.sufParams[k] = v
			}
			pos += step
		} else {
			ctx.target = token
			state = 1
			pos++
		}
	}
}

// showHelp 显示帮助信息
func showHelp() {
	_, hasTagHelp := ctx.preParams["h"]
	_, hasTagHelpFull := ctx.preParams["help"]
	if hasTagHelp || hasTagHelpFull {
		Usage()
		os.Exit(0)
	}
}
