// Copyright 2021 tobyzxj
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netproxy

import (
	"log"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var debugEnable bool

type TraceFunc func(format string, level int, v ...interface{})

var UserTrace TraceFunc = nil

func init() {
	debugEnable = false
	log.SetPrefix("[netproxy] TRACE: ")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// Debug Enable debug
func Debug(enable bool) {
	debugEnable = enable
}

// SetUserDebug 配置其他日志输出
func SetUserDebug(f TraceFunc) {
	UserTrace = f
}

// TraceInfo 调试信息日志
func TraceInfo(format string, v ...interface{}) {
	if debugEnable {
		if UserTrace != nil {
			UserTrace(format, LevelInformational, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}

// TraceError 错误日志
func TraceError(format string, v ...interface{}) {
	if debugEnable {
		if UserTrace != nil {
			UserTrace(format, LevelError, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}
