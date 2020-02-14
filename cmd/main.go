package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"monorepo-components-tags/pkg/git"
	"monorepo-components-tags/pkg/monorepo"
)

var (
	// These are build-time variables that get set by goreleaser.
	version = "dev"

	ErrMissingToken  = errors.New("missing GITLAB_TOKEN or GITHUB_TOKEN  environment variable")
	ErrMultipleToken = errors.New("setup only one of GITLAB_TOKEN or GITHUB_TOKEN variable")
)

func main() {
	app := monorepo.New(nil)
	cmd := buildCLI(app)

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func buildCLI(app *monorepo.App) *cli.App {
	return &cli.App{
		Name:     "monorepo-components-tags",
		Usage:    "get components tags from a git monorepo",
		Version:  "v1.0.0",
		HideHelp: true,
		Before: func(c *cli.Context) error {
			if c.Bool("help") {
				cli.ShowAppHelpAndExit(c, 0)
				return nil
			}

			p, err := getProvider(c)
			if err != nil {
				return err
			}

			app.SetProvider(p)

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token", Hidden: true},
			&cli.StringFlag{
				Name: "project", Required: true, Aliases: []string{"p"},
				Usage: "Your Gitlab Project ID or Github project in the form `user/repo` (required)",
			},
			&cli.StringFlag{
				Name:  "commit",
				Usage: "Find tags including only this commit and any older ones. Any short or long SHA are valid",
			},
			&cli.StringFlag{
				Name: "base-url", DefaultText: "https://gitlab.com/api/v4",
				Usage: "Gitlab base url",
			},
			&cli.StringFlag{
				Name:  "provider",
				Usage: "Set this if you want to force an specific provider. Possible values are `GITLAB|GITHUB`",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "Add prefix to exported names",
			},
			&cli.StringFlag{
				Name:  "suffix",
				Usage: "Change variable names suffix to different value", Value: "VERSION",
			},
			&cli.BoolFlag{
				Name: "export-shell", Aliases: []string{"e"},
				Usage: "format output as shell variables",
			},
			&cli.BoolFlag{Name: "help", Aliases: []string{"h"}},
		},
		Action: func(c *cli.Context) error {
			return app.CollectTags(&monorepo.Options{
				ProjectID:  c.String("project"),
				FromCommit: c.String("commit"),
				VarPrefix:  c.String("prefix"),
				VarSuffix:  c.String("suffix"),
				ExportVars: c.Bool("export-shell"),
			})
		},
	}
}

func getProvider(c *cli.Context) (git.RepoProvider, error) {
	tokGithub := os.Getenv("GITHUB_TOKEN")
	tokGitlab := os.Getenv("GITLAB_TOKEN")
	provider := strings.ToUpper(c.String("provider"))

	// forced provider no token present
	if (provider == "GITLAB" && tokGitlab == "") || (provider == "GITHUB" && tokGithub == "") {
		return nil, ErrMissingToken
	}

	// No provider both tokens exists
	if provider == "" && tokGithub != "" && tokGitlab != "" {
		return nil, ErrMultipleToken
	}

	var p git.RepoProvider

	if provider == "GITHUB" || (tokGithub != "" && tokGitlab == "") {
		p = git.NewGithub(tokGithub)
	}

	if provider == "GITLAB" || (tokGitlab != "" && tokGithub == "") {
		gitlabBaseUrl := c.String("base-url")
		p = git.NewGitlab(tokGitlab, gitlabBaseUrl)
	}

	// Some weird combination
	if p == nil {
		return nil, ErrMissingToken
	}

	return p, nil
}
