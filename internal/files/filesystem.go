package files

import (
	"context"
	"os"
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
	flat := strings.Replace(rel, string(filepath.Separator), f.separator, -1)
	return flat, nil
}

func FindImages(ctx context.Context, root string) ([]string, error) {
	var images []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" {
			images = append(images, path)
		}
		return nil
	})

	return images, err
}
