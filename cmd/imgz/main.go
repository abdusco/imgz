package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"

	"github.com/abdusco/imgz/internal/version"
	"github.com/abdusco/imgz/pkg/resizer"
)

type cliArgs struct {
	Version   kong.VersionFlag `help:"Show version and exit"`
	ImagePath string           `arg:"" help:"Path to image" type:"existingfile"`
	MaxSize   int              `default:"5000"`
	Quality   int              `default:"75"`
}

func (a cliArgs) Run() error {
	res := resizer.New(resizer.Options{
		MaxSize: 7000,
		Quality: 75,
	})

	f, err := os.Open(a.ImagePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	if err := res.Resize(context.Background(), f, os.Stdout); err != nil {
		return fmt.Errorf("failed to resize: %w", err)
	}
	return nil
}

func main() {
	var args cliArgs
	cli := kong.Parse(&args, kong.Name("imgz"), kong.Vars{"version": version.String()})
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
