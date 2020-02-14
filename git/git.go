package git

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/gosuri/uitable"
)

type Git interface {
	GetAllCommits(pid string) []*Commit
	GetAllTags(pid string) TagMap
}

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("UTC")
}

type Commit struct {
	ShortID       string
	CommittedDate *time.Time
	Tags          []Tag
}

func (c Commit) String() string {
	return fmt.Sprintf("[%s] %s", c.ShortID, c.CommittedDate.In(loc).Format(time.RFC3339))
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
			table.AddRow(t.Name, t.Version.String(), t.CommittedDate.In(loc).Format(time.RFC3339), t.ShortID)
		}
	}
	return table.String()
}
