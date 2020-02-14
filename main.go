package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"gitlab-components-tags/components"
	"gitlab-components-tags/git"
)

var projectID string
var fromCommit string
var varPrefix string
var exportVars bool
var wg sync.WaitGroup

func main() {
	gitlabBaseUrl := ""
	fromCommit = ""
	exportVars = false

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		fmt.Println("Missing GITLAB_TOKEN environment variable")
		os.Exit(-1)
	}

	flag.StringVar(&projectID, "project", "", "Gitlab Project Id")
	flag.StringVar(&gitlabBaseUrl, "base-url", "", "Gitlab base URL")
	flag.StringVar(&fromCommit, "commit", "", "Pick your commit")
	flag.BoolVar(&exportVars, "export-shell", false, "Export as shell environment vars")
	flag.StringVar(&varPrefix, "prefix", "", "Env vars prefix")
	flag.Parse()

	var allCommits []*git.Commit
	var tags git.TagMap

	g := git.NewGitlab(gitlabToken, gitlabBaseUrl)
	do(&wg, func() { tags = g.GetAllTags(projectID) })
	do(&wg, func() { allCommits = g.GetAllCommits(projectID) })
	wg.Wait()

	list := components.MakeList(allCommits, tags, fromCommit)

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
