package source

import (
	"os"
)

type Source struct {
	FullPath string
	Contents []rune
}

func (s *Source) LineCol(offset uint) (line uint, col uint) {
	line = 1
	col = 1
	for i := uint(0); i < offset && i < uint(len(s.Contents)); i++ {
		if s.Contents[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
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
