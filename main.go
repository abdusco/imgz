package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/alitto/pond"
	"github.com/mholt/archiver/v4"
	"github.com/schollz/progressbar/v3"

	"imgz/cli"
	"imgz/files"
	"imgz/images"
)

func main() {
	app := cli.New(func(ctx context.Context, args cli.CliArgs) error {
		log.Printf("processing %q", args.Root)

		found, err := files.FindImages(args.Root)
		if err != nil {
			return fmt.Errorf("cannot find images in %q: %v", args.Root, err)
		}
		log.Printf("found %d images", len(found))

		fa, err := os.Create(args.Output)
		if err != nil {
			return fmt.Errorf("cannot create %q: %v", args.Output, err)
		}
		defer fa.Close()

		r := resizer{
			Normalizer: files.FlatteningPathNormalizer{Root: args.Root, Separator: ";"},
			Resizer: images.NewResizer(images.ResizeOptions{
				MaxSize: args.MaxSize,
				Quality: args.Quality,
			}),
		}
		err = r.Resize(ctx, found, fa)
		if err != nil {
			return fmt.Errorf("cannot resize images in %q: %v", args.Root, err)
		}

		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigCh
		log.Printf("received signal: %v", s)
		cancel()
	}()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

type resizer struct {
	Resizer    images.Resizer
	Normalizer files.Normalizer
}

func (d resizer) Resize(ctx context.Context, imagePaths []string, w io.Writer) error {
	z := archiver.Zip{ContinueOnError: true}

	fileCh := make(chan archiver.File)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := z.ArchiveAsync(context.Background(), w, fileCh)
		if err != nil {
			log.Printf("cannot archive images: %v", err)
		}
	}()

	pool := pond.New(runtime.NumCPU(), 0, pond.Context(ctx))
	bar := progressbar.Default(int64(len(imagePaths)))
	defer bar.Close()

	for _, imgPath := range imagePaths {
		imgPath := imgPath
		pool.Submit(func() {
			if err := ctx.Err(); err != nil {
				return
			}
			savePath, err := d.Normalizer.Normalize(imgPath)
			filename := filepath.Base(imgPath)

			if err != nil {
				//log.Printf("cannot normalize %q: %v", filename, err)
				return
			}
			f, err := os.Open(imgPath)
			if err != nil {
				//log.Printf("cannot open %q: %v", filename, err)
				return
			}
			//stat, err := f.Stat()
			//if err != nil || errors.Is(err, os.ErrNotExist) {
			//	//log.Printf("cannot stat %q: %v", filename, err)
			//	return
			//}
			//beforeSize := stat.Size()

			var buf bytes.Buffer
			err = d.Resizer.Resize(f, &buf)
			if err != nil {
				//log.Printf("cannot resize %q: %v", filename, err)
				return
			}
			//afterSize := buf.Len()
			//log.Printf("resized %q: %d -> %d", filename, beforeSize, afterSize)

			bf := bufFile{buf, filename}
			select {
			case fileCh <- archiver.File{
				FileInfo:      bf,
				NameInArchive: savePath,
				Open:          func() (io.ReadCloser, error) { return &bf, nil },
			}:
			default:
			}

			bar.Add(1)
		})
	}

	//log.Printf("waiting")
	pool.StopAndWait()
	close(fileCh)
	wg.Wait()
	return nil
}

type bufFile struct {
	bytes.Buffer
	name string
}

func (i bufFile) Close() error {
	return nil
}

func (i bufFile) Name() string {
	return i.name
}

func (i bufFile) Size() int64 {
	return int64(i.Len())
}

func (i bufFile) Mode() fs.FileMode {
	return 0644
}

func (i bufFile) ModTime() time.Time {
	return time.Now()
}

func (i bufFile) IsDir() bool {
	return false
}

func (i bufFile) Sys() any {
	return nil
}
