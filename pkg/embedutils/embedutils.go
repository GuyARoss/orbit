package embedutils

import "io/fs"

type FileReader interface {
	Read() (fs.File, error)
}
