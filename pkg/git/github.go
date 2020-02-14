package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type Github struct {
	client *github.Client
}

func NewGithub(token string) *Github {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Github{
		client,
	}
}

func (g *Github) GetAllCommits(pid string) ([]*Commit, error) {
	var list []*Commit
	ctx := context.TODO()
	page := 1
	repo := strings.Split(pid, "/")

	for {
		cms, resp, err := g.client.Repositories.ListCommits(ctx, repo[0], repo[1], &github.CommitsListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 500,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error getting all commits: %w", err)
		}

		for _, c := range cms {

			list = append(list, &Commit{
				ShortID:       c.GetSHA(),
				CommittedDate: c.Commit.Author.Date,
				Tags:          []Tag{},
			})
		}

		if resp.NextPage == 0 {
			break
		}
		page += 1
	}

	return list, nil
}

func (g *Github) GetAllTags(pid string) (*TagMap, error) {
	sort := new(string)
	*sort = "desc"

	tagMap := TagMap{}
	page := 1

	ctx := context.TODO()
	repo := strings.Split(pid, "/")
	for {
		ps, resp, err := g.client.Repositories.ListTags(ctx, repo[0], repo[1], &github.ListOptions{
			Page:    page,
			PerPage: 300,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting all tags: %w", err)
		}

		for _, p := range ps {
			name, version := parseVersion(p.Name)
			// Skip weird tags
			if version == nil {
				continue
			}

			if _, ok := tagMap[p.Commit.GetSHA()]; !ok {
				tagMap[p.Commit.GetSHA()] = []Tag{}
			}

			var t *time.Time
			if p.Commit.Author != nil && p.Commit.Author.Date != nil {
				t = p.Commit.Author.Date
			}

			tagMap[p.Commit.GetSHA()] = append(tagMap[p.Commit.GetSHA()], Tag{
				ShortID:       p.Commit.GetSHA(),
				CommittedDate: t,
				Name:          *name,
				Version:       *version,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		page += 1
	}
	return &tagMap, nil
}

// utility to return a pointer to a bool for optional parameters
func (g *Github) bool(f bool) *bool {
	return &f
}
