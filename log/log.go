// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log 日志初始化
package log

import (
	"errors"
	"fmt"
	"os"

	beegolog "github.com/liangdas/mqant/log/beego"
	"github.com/liangdas/mqant/logv2"
	mqanttools "github.com/liangdas/mqant/utils"
)

var hostname, _ = os.Hostname()
var DefaultLogger = newDefaultBeegoLogger(true)

var beeLogger logv2.Logger

// newDefaultBeegoLogger mqant 默认日志实例
func newDefaultBeegoLogger(debug bool, o ...Option) logv2.Logger {
	cc := &Options{}
	for _, opt := range o {
		opt(cc)
	}
	if len(o) == 0 {
		cc = NewOptions()
	}
	defaultLogger := &beegoLogger{
		NewBeegoLoggerV2(Console, cc),
	}
	return logv2.With(defaultLogger)
}

// 自动测试日志
func NewLogger(debug bool, o ...Option) logv2.Logger {
	cc := &Options{}
	for _, opt := range o {
		opt(cc)
	}
	if len(o) == 0 {
		cc = NewOptions()
	}
	defaultLogger := &beegoLogger{
		NewBeegoLoggerV2(Console, cc),
	}
	return logv2.With(defaultLogger)
}

// NewLastVersionLogger 初始化前一个版本的日志
func NewLastVersionLogger(debug bool, ProcessID string, Logdir string, settings map[string]interface{}) {
	defaultLogger := &beegoLogger{
		NewBeegoLogger(debug, hostname, Logdir, settings),
	}
	beeLogger = logv2.With(defaultLogger)
}

// RegisterMqantLogger 注册业务日志
func RegisterMqantLogger(other logv2.Logger) {
	beeLogger = other
}

// LogBeego LogBeego
func LogBeego() logv2.Logger {
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	return beeLogger
}

// CreateTrace CreateTrace
func CreateTrace(trace, span string) TraceSpan {
	return &TraceSpanImp{
		Trace: trace,
		Span:  span,
	}
}

// Info Info
func Info(format string, a ...interface{}) {
	// LogBeego().Log(logv2.LevelInfo.String(), b...)
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelInfo, "", "", "mqant", x)
}

// Error Error
func Error(format string, a ...interface{}) {
	//gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
	x := fmt.Sprintf(format, a...)
	// LogBeego().Log(logv2.LevelError.String(), b...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelError, "", "", "mqant", x)
}

// Warning Warning
func Warning(format string, a ...interface{}) {
	//gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
	// LogBeego().Log(logv2.LevelError.String(), b...)
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelWarn, "", "", "mqant", x)
}

// CreateRootTrace CreateRootTrace
func CreateRootTrace() TraceSpan {
	return &TraceSpanImp{
		Trace: mqanttools.GenerateID().String(),
		Span:  mqanttools.GenerateID().String(),
	}
}

// BiReport BiReport
func BiReport(msg string) {
	// TODO
}

// TDebug TDebug
func TDebug(span TraceSpan, format string, a ...interface{}) {
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelWarn, "", "", "mqant", x)
}

// TInfo TInfo
func TInfo(span TraceSpan, format string, a ...interface{}) {
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelInfo, span.SpanId(), span.TraceId(), "mqant", x)
}

// TError TError
func TError(span TraceSpan, format string, a ...interface{}) {
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelError, span.SpanId(), span.TraceId(), "mqant", x)
}

// TWarning TWarning
func TWarning(span TraceSpan, format string, a ...interface{}) {
	x := fmt.Sprintf(format, a...)
	if beeLogger == nil {
		beeLogger = newDefaultBeegoLogger(false)
	}
	beeLogger.Log(logv2.LevelWarn, span.SpanId(), span.TraceId(), "mqant", x)
}

// Flush() 刷新日志
func Flush() {
	if beeLogger != nil {
		beeLogger.Flush()
	}
}

// Close Close
func Close() {
	if beeLogger != nil {
		beeLogger.Close()
	}
}

type beegoLogger struct {
	*beegolog.BeeLogger
}

func (bg *beegoLogger) Log(level logv2.Level, keyvals ...interface{}) error {
	if len(keyvals) <= 3 {
		return errors.New("keyvals is less than 3")
	}
	// 前四位是Span
	span := &beegolog.BeegoTraceSpan{}
	span.Span = keyvals[0].(string)
	span.Trace = keyvals[1].(string)
	switch level {
	case logv2.LevelWarn:
		bg.Warning(span, keyvals[2].(string), keyvals[3:]...)
	case logv2.LevelError:
		bg.Error(span, keyvals[2].(string), keyvals[3:]...)
	case logv2.LevelInfo:
		bg.Info(span, keyvals[2].(string), keyvals[3:]...)
	case logv2.LevelDebug:
		bg.Debug(span, "xxxx%s", "xxxxxx")
	}

	return nil
}

// Flush() 刷新日志
func (bg *beegoLogger) Flush() {
	bg.BeeLogger.Flush()
}

// Close() 关闭日志写入
func (bg *beegoLogger) Close() {
	bg.BeeLogger.Close()
}
