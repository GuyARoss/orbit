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

type Log interface {
	Error(text string) (int, error)
	Success(text string) (int, error)
	Warn(text string) (int, error)
	Info(text string) (int, error)
}

func Error(text string) (int, error) {
	return fmt.Printf("%s%s%s\n", string(colorRed), text, string(colorReset))
}

func Success(text string) (int, error) {
	return fmt.Printf("%s%s%s\n", string(colorGreen), text, string(colorReset))
}

func Warn(text string) (int, error) {
	return fmt.Printf("(!) %s%s%s\n", string(colorYellow), text, string(colorReset))
}

func Info(text string) (int, error) {
	return fmt.Printf("%s%s%s\n", string(colorBlue), text, string(colorReset))
}

func Title(text string) (int, error) {
	return fmt.Printf("%s%s%s\n", string(blockUnderline), text, string(colorReset))
}
