package file

import (
	"os"
	"time"
)

func Touch(path string) error {
	currentTime := time.Now().Local()
	return os.Chtimes(path, currentTime, currentTime)
}
