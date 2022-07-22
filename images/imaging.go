package images

import (
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

type imagingResizer struct {
	options ResizeOptions
}

func (i imagingResizer) Resize(img image.Image, w io.Writer) error {
	resized := imaging.Fit(img, int(i.options.MaxSize), int(i.options.MaxSize), imaging.Lanczos)

	err := imaging.Encode(w, resized, imaging.JPEG, imaging.JPEGQuality(int(i.options.Quality)))
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}

	return nil
}
