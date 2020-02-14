package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"gitlab-components-tags/cmd"
)

var (
	ErrMissingGitlabToken = errors.New("missing GITLAB_TOKEN environment variable")
)

func main() {
	app := &cli.App{
		Name:     "gitlab-components-tags",
		Usage:    "get components tags from gitlab",
		Version:  "v1.0.0",
		HideHelp: true,
		Before: func(c *cli.Context) error {
			if c.Bool("help") {
				cli.ShowAppHelpAndExit(c, 0)
				return nil
			}
			token := os.Getenv("GITLAB_TOKEN")
			if token == "" {
				return ErrMissingGitlabToken
			}

			_ = c.Set("token", token)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "project", Required: true, Aliases: []string{"p"}, Usage: "Gitlab ProjectID (required)"},
			&cli.StringFlag{Name: "commit", Usage: "Short ID commit from where to start the search"},
			&cli.StringFlag{Name: "base-url", Usage: "Gitlab base url", DefaultText: "https://gitlab.com/api/v4"},
			&cli.StringFlag{Name: "prefix", Usage: "Add prefix to exported names"},
			&cli.BoolFlag{Name: "export-shell", Aliases: []string{"e"}, Usage: "format output as shell variables"},
			&cli.BoolFlag{Name: "help", Aliases: []string{"h"}},
		},
		Action: cmd.ComponentsTags,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
}
