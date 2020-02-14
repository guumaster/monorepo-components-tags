package monorepo

import (
	"fmt"
	"strings"
	"sync"

	"monorepo-components-tags/pkg/components"
	"monorepo-components-tags/pkg/git"
)

var wg sync.WaitGroup

// App contains the methods to scan and clean paths.
type App struct {
	provider git.RepoProvider
}

type Options struct {
	ProjectID  string
	FromCommit string
	VarPrefix  string
	VarSuffix  string
	ExportVars bool
}

func New(p git.RepoProvider) *App {
	return &App{
		provider: p,
	}
}

func (a *App) SetProvider(provider git.RepoProvider) {
	a.provider = provider
}

func (a *App) CollectTags(opts *Options) error {

	var allCommits []*git.Commit
	var tags *git.TagMap
	var err error

	do(&wg, func() {
		tags, err = a.provider.GetAllTags(opts.ProjectID)
	})
	do(&wg, func() {
		allCommits, err = a.provider.GetAllCommits(opts.ProjectID)
	})
	wg.Wait()

	if err != nil {
		return err
	}

	list := components.MakeList(allCommits, *tags, opts.FromCommit)

	if opts.ExportVars {
		if opts.VarPrefix != "" {
			opts.VarPrefix = strings.ToUpper(opts.VarPrefix) + "_"
		}
		opts.VarSuffix = strings.ToUpper(opts.VarSuffix)
		for _, c := range list {
			fmt.Printf("export %s%s_%s=\"%s\"\n", opts.VarPrefix, toEnvVar(c.Name), opts.VarSuffix, c.Version)
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
