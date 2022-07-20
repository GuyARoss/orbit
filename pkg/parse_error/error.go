// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package parseerror

import "fmt"

func New(err string, fileName string) error {
	return fmt.Errorf("%s: %s", fileName, err)
}

func FromError(err error, fileName string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %s", fileName, err.Error())
}
