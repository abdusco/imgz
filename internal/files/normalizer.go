package files

import (
	"path/filepath"
	"strings"
)

type Normalizer interface {
	Normalize(path string) (string, error)
}

type NoopNormalizer struct {
}

func (n NoopNormalizer) Normalize(path string) (string, error) {
	return path, nil
}

type FlatteningPathNormalizer struct {
	root      string
	separator string
}

func NewFlattener(root string) FlatteningPathNormalizer {
	return FlatteningPathNormalizer{root: root, separator: ";"}
}

func (f FlatteningPathNormalizer) Normalize(path string) (string, error) {
	rel, err := filepath.Rel(f.root, path)
	if err != nil {
		return "", err
	}
	flat := strings.ReplaceAll(rel, string(filepath.Separator), f.separator)
	return flat, nil
}
