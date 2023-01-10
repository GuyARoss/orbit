// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package log

import (
	"fmt"
	"runtime"
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"

	blockUnderline = "\033[4m"
	clear          = "\x1B[2J\x1B[3J\x1B[H"
	clearWin       = "\x1B[2J\x1B[0f"
)

type Logger interface {
	Error(text string) (int, error)
	Success(text string) (int, error)
	Warn(text string) (int, error)
	Info(text string) (int, error)
	Clear()
}

type DefaultLogger struct {
	outf func(string, ...interface{}) (int, error)
}

func NewEmptyLogger() *DefaultLogger {
	return &DefaultLogger{
		outf: func(s string, i ...interface{}) (int, error) {
			return 0, nil
		},
	}
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		outf: fmt.Printf,
	}
}

func (l *DefaultLogger) Clear() {
	if runtime.GOOS == "windows" {
		l.outf(clearWin)
		return
	}

	l.outf(clear)
}

func (l *DefaultLogger) Error(text string) (int, error) {
	return l.outf("%s%s%s\n", string(colorRed), text, string(colorReset))
}

func (l *DefaultLogger) Success(text string) (int, error) {
	return l.outf("%s%s%s\n", string(colorGreen), text, string(colorReset))
}

func (l *DefaultLogger) Warn(text string) (int, error) {
	return l.outf("(!) %s%s%s\n", string(colorYellow), text, string(colorReset))
}

func (l *DefaultLogger) Info(text string) (int, error) {
	return l.outf("(i) %s%s%s\n", string(colorBlue), text, string(colorReset))
}

func (l *DefaultLogger) Title(text string) (int, error) {
	return l.outf("%s%s%s\n", string(blockUnderline), text, string(colorReset))
}
