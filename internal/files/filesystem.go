package files

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

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
