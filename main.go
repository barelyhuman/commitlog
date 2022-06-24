package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/barelyhuman/commitlog/commands"
	"github.com/urfave/cli/v2"
)

// go:embed .commitlog.release
var version string

func main() {
	app := &cli.App{
		Name:  "commitlog",
		Usage: "commits to changelogs",
		Action: func(c *cli.Context) error {
			fmt.Println(
				"[commitlog] we no longer support direct invocation, please use the subcommand `generate` to generate a log or `release` to manage your .commitlog.release",
			)

			return nil
		},
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
						Usage: "initialise commitlog release",
					},
					&cli.BoolFlag{
						Name: "pre-release",
						Usage: "create a pre-release version. will default to patch increment unless" +
							"specified and not already a pre-release",
					},
					&cli.StringFlag{
						Name:  "pre-release-tag",
						Value: "beta",
						Usage: "create a pre-release version",
					},
					&cli.BoolFlag{
						Name:  "major",
						Usage: "create a major version",
					},
					&cli.BoolFlag{
						Name:  "minor",
						Usage: "create a minor version",
					},
					&cli.BoolFlag{
						Name:  "patch",
						Usage: "create a patch version",
					},
					&cli.BoolFlag{
						Name:  "commit",
						Value: false,
						Usage: "if true will create a commit, of the changed version",
					},
					&cli.BoolFlag{
						Name:  "tag",
						Value: false,
						Usage: "if true will create a tag, with the given version",
					},
				},
			},
		},
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}

	app.Version = version

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
