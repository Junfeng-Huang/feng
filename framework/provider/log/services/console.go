package services

import (
	"os"

	"feng/framework/container"
	"feng/framework/contract"
)

// FengConsoleLog 代表控制台输出
type FengConsoleLog struct {
	FengLog
}

// NewFengConsoleLog 实例化FengConsoleLog
func NewFengConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(container.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &FengConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
