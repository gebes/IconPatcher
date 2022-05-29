package file

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

func Find(root, pattern string) ([]string, error) {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var a []string
	err = filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if compiledPattern.MatchString(d.Name()) {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}
