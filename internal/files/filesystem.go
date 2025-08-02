package files

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

func isImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
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
		if isImage(path) {
			images = append(images, path)
		}
		return nil
	})

	return images, err
}
