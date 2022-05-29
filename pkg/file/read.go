package file

import (
	"io/ioutil"
	"os"
)

func Read(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}
