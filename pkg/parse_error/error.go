package parseerror

import "fmt"

func New(err string, fileName string) error {
	return fmt.Errorf("%s: %s", fileName, err)
}

func FromError(err error, fileName string) error {
	return fmt.Errorf("%s: %s", fileName, err.Error())
}
