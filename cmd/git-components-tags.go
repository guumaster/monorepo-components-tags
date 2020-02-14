package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"

	"gitlab-components-tags/internal/components"
	"gitlab-components-tags/internal/git"
)

var wg sync.WaitGroup

func ComponentsTags(c *cli.Context) error {
	gitlabToken := c.String("token")
	projectID := c.String("project")
	gitlabBaseUrl := c.String("base-url")
	fromCommit := c.String("commit")
	varPrefix := c.String("prefix")
	exportVars := c.Bool("export-shell")

	var allCommits []*git.Commit
	var tags *git.TagMap
	var err error

	g := git.NewGitlab(gitlabToken, gitlabBaseUrl)
	do(&wg, func() {
		tags, err = g.GetAllTags(projectID)
	})
	do(&wg, func() {
		allCommits, err = g.GetAllCommits(projectID)
	})
	wg.Wait()

	if err != nil {
		return err
	}

	list := components.MakeList(allCommits, *tags, fromCommit)

	if exportVars {
		if varPrefix != "" {
			varPrefix = strings.ToUpper(varPrefix) + "_"
		}
		for _, c := range list {
			fmt.Printf("export %s%s_VERSION=\"%s\"\n", varPrefix, toEnvVar(c.Name), c.Version)
		}
	} else {
		fmt.Println(list)
	}
	return nil
}

func toEnvVar(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}

func do(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
