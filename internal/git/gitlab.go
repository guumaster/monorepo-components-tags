package git

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/xanzy/go-gitlab"
)

type Gitlab struct {
	client *gitlab.Client
}

func NewGitlab(token, baseUrl string) *Gitlab {
	client := gitlab.NewClient(nil, token)

	if baseUrl != "" {
		_ = client.SetBaseURL(baseUrl)
	}
	return &Gitlab{
		client,
	}
}

func (g *Gitlab) GetAllCommits(pid string) ([]*Commit, error) {
	var list []*Commit
	page := 1
	for {
		cms, resp, err := g.client.Commits.ListCommits(pid, &gitlab.ListCommitsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 500,
			},
			All:       g.bool(false),
			WithStats: g.bool(false),
		})
		if err != nil {
			return nil, fmt.Errorf("error getting all commits: %w", err)
		}

		for _, c := range cms {
			list = append(list, &Commit{
				ShortID:       c.ShortID,
				CommittedDate: c.CommittedDate,
				Tags:          []Tag{},
			})
		}

		if resp.CurrentPage >= resp.TotalPages || len(list) >= resp.TotalItems {
			break
		}
		page += 1
	}

	return list, nil
}

func (g *Gitlab) GetAllTags(pid string) (*TagMap, error) {
	sort := new(string)
	*sort = "desc"

	tagMap := TagMap{}
	page := 1

	for {
		ps, resp, err := g.client.Tags.ListTags(pid, &gitlab.ListTagsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
			Sort: sort,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting all tags: %w", err)
		}

		for _, p := range ps {
			name, version := g.parseVersion(p.Name)
			// Skip weird tags
			if version == nil {
				continue
			}

			if _, ok := tagMap[p.Commit.ShortID]; !ok {
				tagMap[p.Commit.ShortID] = []Tag{}
			}

			tagMap[p.Commit.ShortID] = append(tagMap[p.Commit.ShortID], Tag{
				ShortID:       p.Commit.ShortID,
				CommittedDate: p.Commit.CommittedDate,
				Name:          *name,
				Version:       *version,
			})
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		page += 1
	}
	return &tagMap, nil
}

func (g *Gitlab) parseVersion(name string) (*string, *semver.Version) {
	ss := strings.Split(name, "-v")
	if len(ss) != 2 {
		return nil, nil
	}

	version, err := semver.Make(ss[1])
	if err != nil {
		return nil, nil
	}

	return &ss[0], &version
}

// utility to return a pointer to a bool for optional parameters
func (g *Gitlab) bool(f bool) *bool {
	return &f
}
