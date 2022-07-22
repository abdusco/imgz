package files

import (
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
	Root      string
	Separator string
}

func (f FlatteningPathNormalizer) Normalize(path string) (string, error) {
	rel, err := filepath.Rel(f.Root, path)
	if err != nil {
		return "", err
	}
	flat := strings.Replace(rel, string(filepath.Separator), f.Separator, -1)
	return flat, nil
}

func FindImages(root string) ([]string, error) {
	var images []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
