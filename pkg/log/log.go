package log

import "fmt"

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"

	blockUnderline = "\033[4m"
)

type Logger interface {
	Error(text string) (int, error)
	Success(text string) (int, error)
	Warn(text string) (int, error)
	Info(text string) (int, error)
}

type DefaultLogger struct {
	outf func(string, ...interface{}) (int, error)
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		outf: fmt.Printf,
	}
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
	return l.outf("%s%s%s\n", string(colorBlue), text, string(colorReset))
}

func (l *DefaultLogger) Title(text string) (int, error) {
	return l.outf("%s%s%s\n", string(blockUnderline), text, string(colorReset))
}
