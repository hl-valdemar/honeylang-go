package source

import (
	"os"
)

type Source struct {
	FullPath string
	Contents []rune
}

func Load(fullPath string) (*Source, error) {
	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	runes := []rune(string(bytes))
	src := Source{fullPath, runes}
	return &src, nil
}
