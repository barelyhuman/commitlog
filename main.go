package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/barelyhuman/commitlog/commands"
	"github.com/urfave/cli/v2"
)

//go:embed .commitlog.release
var version string

func main() {
	app := &cli.App{
		Name:            "commitlog",
		Usage:           "commits to changelogs",
		CommandNotFound: cli.ShowCommandCompletions,
		Action: func(c *cli.Context) error {
			fmt.Println(
				"[commitlog] we no longer support direct invocation, please use the subcommand `generate` to generate a log or `release` to manage your .commitlog.release",
			)

			return nil
		},
		Version:     version,
		Compiled:    time.Now(),
		HideVersion: false,
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "commits to changelogs",
				Action: func(c *cli.Context) error {
					return commands.Commitlog(c)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "path",
						Value:   ".",
						Aliases: []string{"p"},
						Usage:   "root with the '.git' folder `PATH`",
					},
					&cli.BoolFlag{
						Name:  "promo",
						Usage: "add promo text to the end of output",
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "path to the output `FILE`",
					},
					&cli.BoolFlag{
						Name:  "stdio",
						Value: true,
						Usage: "print to the stdout",
					},
					&cli.StringFlag{
						Name:  "categories",
						Value: "",
						Usage: "categories to use, includes all commits by default. any text you add here will be used to create categories out of the commits",
					},
					&cli.StringFlag{
						Name:    "start",
						Aliases: []string{"s"},
						Usage: "`START` reference for the commit to include commits from," +
							"This is inclusive of the given commit reference",
					},
					&cli.StringFlag{
						Name:    "end",
						Aliases: []string{"e"},
						Usage: "`END` reference for the commit to stop including commits at." +
							"This is exclusive of the given commit reference",
					},
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "manage .commitlog.release version",
				Action: func(c *cli.Context) error {
					return commands.Release(c)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "init",
						Value: false,
						Usage: "initialise commitlog release",
					},
					&cli.StringFlag{
						Name:    "path",
						Value:   ".",
						Aliases: []string{"p"},
						Usage: "`PATH` to a folder where .commitlog.release exists or is to be created." +
							"(note: do not use `--commit` or `--push` if the file isn't in the root)",
					},
					&cli.BoolFlag{
						Name:  "pre",
						Value: false,
						Usage: "create a pre-release version. will default to patch increment unless" +
							"specified and not already a pre-release",
					},
					&cli.StringFlag{
						Name:  "pre-tag",
						Value: "beta",
						Usage: "create a pre-release version",
					},
					&cli.BoolFlag{
						Name:  "major",
						Value: false,
						Usage: "create a major version",
					},
					&cli.BoolFlag{
						Name:  "minor",
						Value: false,
						Usage: "create a minor version",
					},
					&cli.BoolFlag{
						Name:  "patch",
						Value: false,
						Usage: "create a patch version",
					},
					&cli.BoolFlag{
						Name:  "commit",
						Value: false,
						Usage: "if true will create a commit and tag, of the changed version",
					},
					&cli.BoolFlag{
						Name:  "push",
						Value: false,
						Usage: "if true will create push the created release commit + tag on origin",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[commitlog] %v", err)
	}
}
