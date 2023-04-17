package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"golang.org/x/exp/slog"

	"github.com/abdusco/imgz/internal/files"
	"github.com/abdusco/imgz/internal/version"
	"github.com/abdusco/imgz/pkg/resizer"
)

func main() {
	var cliApp struct {
		Debug     bool             `help:"Enable debug logging"`
		Version   kong.VersionFlag `help:"Show version and exit"`
		Resize    resizeCmd        `cmd:"" help:"Resize an image"`
		ResizeDir resizeDirCmd     `cmd:"" help:"Resize a folder of images"`
	}
	cli := kong.Parse(
		&cliApp,
		kong.Name("imgz"),
		kong.Vars{"version": version.String()},
	)
	if cliApp.Debug {
		slog.SetDefault(slog.New(slog.HandlerOptions{Level: slog.LevelDebug}.NewTextHandler(os.Stderr)))
	}
	if err := cli.Run(); err != nil {
		slog.Error("exit with error", "error", err)
		os.Exit(1)
	}
}

type resizeDirCmd struct {
	SourceDirs []string `arg:"" help:"List of paths to image folders" type:"existingdir"`
	OutputDir  string   `short:"o" help:"Dir to save zip files" required:""`
	MaxSize    uint     `default:"5000" help:"Max side length of resized images"`
	Quality    uint     `default:"75" help:"JPEG quality"`
	Clean      bool     `help:"Delete source dirs after resizing"`
	SkipDone   bool     `help:"Skip processing a dir if there's a matching ZIP file in output dir'"`
}

func (c resizeDirCmd) Run() error {
	ctx := context.Background()
	res := resizer.New(resizer.Options{
		MaxSize: c.MaxSize,
		Quality: c.Quality,
	})

	var dirsToRemove []string
	for _, dir := range c.SourceDirs {
		norm := files.NewFlattener(dir)

		images, err := files.FindImages(ctx, dir)
		if err != nil {
			return fmt.Errorf("failed to find images in %q: %w", dir, err)
		}

		if len(images) == 0 {
			slog.Info("dir has no images", "dir", dir)
			continue
		}

		zipPath := filepath.Join(c.OutputDir, fmt.Sprintf("%s.zip", filepath.Base(dir)))
		if c.SkipDone {
			if _, err := os.Stat(zipPath); err == nil {
				slog.Info("dir has already been processed, skipping", "dir", dir, "zip_path", zipPath)
				continue
			}
		}

		slog.Debug("processing dir", "dir", dir, "total_images", len(images))

		zf, err := os.Create(zipPath)
		if err != nil {
			return fmt.Errorf("failed to create zip file: %w", err)
		}
		zw := zip.NewWriter(zf)

		for _, imagePath := range images {
			slog.Debug("resizing image", "filename", filepath.Base(imagePath))
			f, err := os.Open(imagePath)
			if err != nil {
				slog.Error("failed to open image", "path", imagePath, "error", err)
				continue
			}

			archivePath, err := norm.Normalize(imagePath)
			if err != nil {
				return fmt.Errorf("failed to normalize archive path: %w", err)
			}

			w, err := zw.Create(archivePath)
			if err != nil {
				return fmt.Errorf("failed to create zip file entry: %w", err)
			}

			if err := res.Resize(ctx, f, w); err != nil {
				return fmt.Errorf("failed to resize image: %w", err)
			}

			f.Close()
		}
		zw.Close()
		dirsToRemove = append(dirsToRemove, dir)
	}

	if c.Clean {
		for _, dir := range dirsToRemove {
			slog.Debug("removing dir", "dir", dir)
			if err := os.RemoveAll(dir); err != nil {
				slog.Error("failed to remove dir", "dir", dir, "error", err)
				continue
			}
		}
	}

	return nil
}

type resizeCmd struct {
	ImagePath string `arg:"" help:"Path to image. Use \"-\" for stdin"`
	Output    string `short:"o" help:"Path to output file. Use \"-\" for stdout. Defaults to $sourceDir/resized/$source.jpg or stdout if input is stdin" `
	MaxSize   uint   `default:"5000"`
	Quality   uint   `default:"75"`
}

func (c resizeCmd) Run() error {
	res := resizer.New(resizer.Options{
		MaxSize: c.MaxSize,
		Quality: c.Quality,
	})

	var r io.Reader
	if c.ImagePath == "-" {
		r = os.Stdin
		slog.Debug("input is set to stdin")
	} else {
		var err error
		c.ImagePath, err = filepath.Abs(c.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to resolve image path: %w", err)
		}

		f, err := os.Open(c.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		slog.Debug("input is set to file", "path", c.ImagePath)
		r = f
	}

	if c.Output == "" && c.ImagePath == "-" {
		c.Output = "-"
	}

	var w io.Writer
	if c.Output == "-" {
		w = os.Stdout
		slog.Info("saving to stdout")
	} else {
		ext := ".jpg"
		basename := strings.TrimSuffix(filepath.Base(c.ImagePath), filepath.Ext(c.ImagePath))
		outputPath := filepath.Join(filepath.Dir(c.ImagePath), "resized", basename+ext)
		dir := filepath.Dir(outputPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory at %q: %w", dir, err)
		}
		outputPath, err := filepath.Abs(outputPath)
		if err != nil {
			return fmt.Errorf("failed resolve output path: %w", err)
		}

		if outputPath == c.ImagePath {
			return fmt.Errorf("output path cannot be the same as input path")
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		slog.Debug("will save to file", "path", outputPath)
		w = f
	}

	if err := res.Resize(context.Background(), r, w); err != nil {
		return fmt.Errorf("failed to resize: %w", err)
	}
	return nil
}
