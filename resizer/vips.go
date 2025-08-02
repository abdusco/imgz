package resizer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func init() {
	if _, err := exec.LookPath("vips"); err != nil {
		return
	}
	register("vips", func(o Options) Resizer {
		return vipsResizer{o}
	})
	defaultStrategy = "vips"
}

type vipsResizer struct {
	Options Options
}

func (v vipsResizer) Resize(ctx context.Context, r io.Reader, w io.Writer) error {
	args := []string{"thumbnail", "/dev/stdin"}

	jpegSaveOpts := []string{"optimize_coding", "strip=true"}
	if v.Options.Quality > 0 {
		jpegSaveOpts = append(jpegSaveOpts, fmt.Sprintf("Q=%d", v.Options.Quality))
	}
	args = append(args, fmt.Sprintf(".jpg[%s]", strings.Join(jpegSaveOpts, ",")))

	if v.Options.MaxSize > 0 {
		args = append(args, fmt.Sprintf("%dx%d", v.Options.MaxSize, v.Options.MaxSize), "--size", "down")
	}

	cmd := exec.CommandContext(ctx, "vips", args...)
	cmd.Stdin = r
	cmd.Stdout = w
	var b bytes.Buffer
	cmd.Stderr = &b
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if err.ExitCode() != 0 {
				return fmt.Errorf("failed to run vips: %w: stderr=%s", err, b.String())
			}
		}
	}
	return nil
}
