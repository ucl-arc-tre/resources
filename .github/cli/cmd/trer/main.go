package main

import (
	"errors"
	"log"
	"os"

	"github.com/ucl-arc-tre/global/internal/check"
	"github.com/ucl-arc-tre/global/internal/matrix"
	"github.com/ucl-arc-tre/global/internal/script"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "trer"
	app.Usage = "UCL ARC TRE global packages CLI"

	app.Commands = []cli.Command{
		{
			Name:  "matrix",
			Usage: "Print a GitHub workflow matrix given the changes in this repo",
			Action: func(c *cli.Context) error {
				return matrix.Print(c.String("repo-root"))
			},
			Flags: []cli.Flag{cli.StringFlag{
				Name:  "repo-root, r",
				Value: ".",
				Usage: "Relative path to the git repository root",
			}},
		},
		{
			Name:  "script",
			Usage: "Print a shell script for building a package. Is context aware",
			Action: func(c *cli.Context) error {
				if dir := c.Args().Get(0); dir != "" {
					return script.Print(dir)
				} else {
					return errors.New("Missing path argument")
				}
			},
		},
		{
			Name:  "check",
			Usage: "Check that the version file has been modified if any files in the directory have",
			Action: func(c *cli.Context) error {
				return check.CheckVersionBumps(c.String("repo-root"))
			},
			Flags: []cli.Flag{cli.StringFlag{
				Name:  "repo-root, r",
				Value: ".",
				Usage: "Relative path to the git repository root",
			}},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
