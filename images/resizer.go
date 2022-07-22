package images

import (
	"image"
	"io"
)

type ResizeOptions struct {
	MaxSize uint64
	Quality uint64
}

type Resizer interface {
	Resize(img image.Image, w io.Writer) error
}

func NewResizer(o ResizeOptions) Resizer {
	return &imagingResizer{o}
}
