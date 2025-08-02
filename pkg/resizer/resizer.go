package resizer

import (
	"context"
	"io"
	"log/slog"
)

func init() {
	register(defaultStrategy, func(o Options) Resizer {
		return imagingResizer{o}
	})
}

type Options struct {
	MaxSize uint
	Quality uint
}

type Resizer interface {
	Resize(ctx context.Context, r io.Reader, w io.Writer) error
}

var strategies = make(map[string]func(o Options) Resizer)
var defaultStrategy = "imaging"

func register(name string, factory func(o Options) Resizer) {
	strategies[name] = factory
}

func New(o Options) Resizer {
	slog.Debug("resize strategy", "strategy", defaultStrategy)
	return strategies[defaultStrategy](o)
}
