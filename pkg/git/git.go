package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/gosuri/uitable"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("UTC")
}

type RepoProvider interface {
	GetAllCommits(pid string) ([]*Commit, error)
	GetAllTags(pid string) (*TagMap, error)
}

type Commit struct {
	ShortID       string
	CommittedDate *time.Time
	Tags          []Tag
}

func (c Commit) String() string {
	if c.CommittedDate != nil {
		return fmt.Sprintf("[%s] %s", c.ShortID, c.CommittedDate.In(loc).Format(time.RFC3339))

	}
	return fmt.Sprintf("[%s]", c.ShortID)
}

type Tag struct {
	ShortID       string
	CommittedDate *time.Time
	Name          string
	Version       semver.Version
}

type TagMap map[string][]Tag

func (m TagMap) String() string {
	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("Component", "Version", "Date", "Commit")
	for _, ts := range m {
		for _, t := range ts {
			date := ""
			if t.CommittedDate != nil {
				date = t.CommittedDate.In(loc).Format(time.RFC3339)
			}
			table.AddRow(t.Name, t.Version.String(), date, t.ShortID)
		}
	}
	return table.String()
}

func parseVersion(name *string) (*string, *semver.Version) {
	ss := strings.Split(*name, "-v")
	if len(ss) < 2 {
		return nil, nil
	}
	v := ss[len(ss)-1]

	version, err := semver.Make(v)
	if err != nil {
		return nil, nil
	}

	return &ss[0], &version
}
