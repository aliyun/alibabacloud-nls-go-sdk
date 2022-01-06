/*
log.go

Copyright 1999-present Alibaba Group Holding Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nls

import (
	"io"
	"log"
	"os"
)

type NlsLogger struct {
	logger *log.Logger
	sil    bool
	debug  bool
}

var defaultLog *NlsLogger
var defaultStd = stdoutLogger(log.LstdFlags|log.Lmicroseconds, "NLS")

func DefaultNlsLog() *NlsLogger {
	return defaultStd
}

func stdoutLogger(flag int, tag string) *NlsLogger {
	logger := new(NlsLogger)
	logger.logger = log.New(os.Stderr, tag, flag)
	logger.sil = false
	logger.debug = false
	return logger
}

func NewNlsLogger(w io.Writer, tag string, flag int) *NlsLogger {
	logger := new(NlsLogger)
	logger.logger = log.New(w, tag, flag)
	logger.sil = false
	logger.debug = false
	return logger
}

func (l *NlsLogger) SetLogSil(sil bool) {
	l.sil = sil
}

func (l *NlsLogger) SetDebug(debug bool) {
	l.debug = debug
}

func (l *NlsLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *NlsLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *NlsLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func (l *NlsLogger) Fatalln(v ...interface{}) {
	l.logger.Fatalln(v...)
}

func (l *NlsLogger) Panic(v ...interface{}) {
	l.logger.Panic(v...)
}

func (l *NlsLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

func (l *NlsLogger) panicln(v ...interface{}) {
	l.logger.Panicln(v...)
}

func (l *NlsLogger) Print(v ...interface{}) {
	if l.sil {
		return
	}
	l.logger.Print(v...)
}

func (l *NlsLogger) Printf(format string, v ...interface{}) {
	if l.sil {
		return
	}
	l.logger.Printf(format, v...)
}

func (l *NlsLogger) Println(v ...interface{}) {
	if l.sil {
		return
	}
	l.logger.Println(v...)
}

func (l *NlsLogger) Debugln(v ...interface{}) {
	if l.debug {
		l.logger.Println(v...)
	}
}

func (l *NlsLogger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.logger.Printf(format, v...)
	}
}

func (l *NlsLogger) SetFlags(flags int) {
	l.logger.SetFlags(flags)
}

func (l *NlsLogger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}
