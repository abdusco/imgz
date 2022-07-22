package images

import (
	"io"
	"log"
)

func init() {
	register(defaultStrategy, func(o ResizeOptions) Resizer {
		return imagingResizer{o}
	})
}

type ResizeOptions struct {
	MaxSize uint64
	Quality uint64
}

type Resizer interface {
	Resize(r io.Reader, w io.Writer) error
}

var strategies = make(map[string]func(o ResizeOptions) Resizer)
var defaultStrategy = "imaging"

func register(name string, factory func(o ResizeOptions) Resizer) {
	strategies[name] = factory
}

func NewResizer(o ResizeOptions) Resizer {
	log.Printf("Using %s strategy", defaultStrategy)
	return strategies[defaultStrategy](o)
}
