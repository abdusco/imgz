//go:build vips

package images

import (
	"bytes"
	"fmt"
	"io"

	"github.com/h2non/bimg"
)

func init() {
	register("vips", func(o ResizeOptions) Resizer {
		return vipsResizer{o}
	})
	defaultStrategy = "vips"
}

type vipsResizer struct {
	Options ResizeOptions
}

func (v vipsResizer) Resize(r io.Reader, w io.Writer) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}
	img := bimg.NewImage(buf.Bytes())

	out, err := img.Resize(int(v.Options.MaxSize), int(v.Options.MaxSize))
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}
	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}
	return nil
}
