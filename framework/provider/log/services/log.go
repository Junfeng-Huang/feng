package services

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"feng/framework/container"
	"feng/framework/contract"
	"feng/framework/provider/log/formatter"
)

// FengLog 的通用实例
type FengLog struct {
	// 五个必要参数
	level      contract.LogLevel   // 日志级别
	formatter  contract.Formatter  // 日志格式化方法
	ctxFielder contract.CtxFielder // ctx获取上下文字段
	output     io.Writer           // 输出
	c          container.Container // 容器
}

// IsLevelEnable 判断这个级别是否可以打印
func (log *FengLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

// logf 为打印日志的核心函数
func (log *FengLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	// 先判断日志级别
	if !log.IsLevelEnable(level) {
		return nil
	}

	// 使用ctxFielder 获取context中的信息
	fs := fields
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	// // 如果绑定了trace服务，获取trace信息
	// if log.c.IsBind(contract.TraceKey) {
	// 	tracer := log.c.MustMake(contract.TraceKey).(contract.Trace)
	// 	tc := tracer.GetTrace(ctx)
	// 	if tc != nil {
	// 		maps := tracer.ToMap(tc)
	// 		for k, v := range maps {
	// 			fs[k] = v
	// 		}
	// 	}
	// }

	// 将日志信息按照formatter序列化为字符串
	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}
	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用log进行panic
	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	// 通过output进行输出
	log.output.Write(ct)
	// 因为每个ct都是单独的对象，无需加锁，且在写入日志文件时文件模式为O_APPEND是原子操作。
	// 标准库的log写入时之所以加锁是因为其写日志的对象是共享的。
	log.output.Write([]byte("\r\n"))
	return nil
}

// SetOutput 设置output
func (log *FengLog) SetOutput(output io.Writer) {
	log.output = output
}

// Panic 输出panic的日志信息
func (log *FengLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

// Fatal will add fatal record which contains msg and fields
func (log *FengLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

// Error will add error record which contains msg and fields
func (log *FengLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

// Warn will add warn record which contains msg and fields
func (log *FengLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

// Info 会打印出普通的日志信息
func (log *FengLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

// Debug will add debug record which contains msg and fields
func (log *FengLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

// Trace will add trace info which contains msg and fields
func (log *FengLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

// SetLevel set log level, and higher level will be recorded
func (log *FengLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

// SetCxtFielder will get fields from context
func (log *FengLog) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

// SetFormatter will set formatter handler will covert data to string for recording
func (log *FengLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}
