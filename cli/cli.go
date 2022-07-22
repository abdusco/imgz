package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

type CliArgs struct {
	Root    string
	Output  string
	MaxSize uint64
	Quality uint64
}

type RunnerFunc func(ctx context.Context, args CliArgs) error

func New(runner RunnerFunc) *cli.App {
	return &cli.App{
		Name:           "imgz",
		DefaultCommand: "resize",
		Commands: []*cli.Command{
			{
				Name:      "resize",
				ArgsUsage: "source_dir",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output ZIP file",
					},
					&cli.Uint64Flag{
						Name:  "max-size",
						Usage: "Length of longest side of the image in px",
						Value: 5000,
					},
					&cli.Uint64Flag{
						Name:  "quality",
						Usage: "JPEG quality between 0-100",
						Value: 75,
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 1 {
						cli.ShowAppHelpAndExit(ctx, 1)
						return nil
					}
					root := ctx.Args().First()
					root = expandHome(root)
					output := ctx.String("output")
					if output == "" {
						output = fmt.Sprintf("%s.zip", filepath.Base(root))
					}
					args := CliArgs{
						Root:    root,
						Output:  output,
						MaxSize: ctx.Uint64("max-size"),
						Quality: ctx.Uint64("quality"),
					}

					return runner(ctx.Context, args)
				},
			},
		},
	}
}

func expandHome(root string) string {
	if strings.HasPrefix(root, "~") {
		home, _ := os.UserHomeDir()
		root = filepath.Join(home, root[2:])
	}
	return root
}
