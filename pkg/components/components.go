package components

import (
	"fmt"
	"sort"
	"time"

	"github.com/blang/semver"
	"github.com/gosuri/uitable"

	"monorepo-components-tags/pkg/git"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("UTC")
}

type Component struct {
	Name          string
	ShortID       string
	CommittedDate *time.Time
	Version       semver.Version
}

type ComponentMap map[string]Component

func (c *Component) String() string {
	if c.CommittedDate != nil {
		return fmt.Sprintf("%s-%s (%s)", c.Name, c.Version.String(), c.CommittedDate.In(loc).Format(time.RFC3339))

	}
	return fmt.Sprintf("%s-%s (-)", c.Name, c.Version.String())
}

func (m ComponentMap) String() string {
	table := uitable.New()
	table.MaxColWidth = 50

	var list []Component

	for _, c := range m {
		list = append(list, c)
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].CommittedDate == nil || list[j].CommittedDate == nil {
			return true
		}
		if list[i].CommittedDate.Unix() == list[j].CommittedDate.Unix() {
			return list[i].Version.GTE(list[j].Version)
		}
		return list[i].CommittedDate.UnixNano() > list[j].CommittedDate.UnixNano()
	})

	table.AddRow("Component", "Version", "Date", "Commit")
	for _, c := range list {
		date := ""
		if c.CommittedDate != nil {
			date = c.CommittedDate.In(loc).Format(time.RFC3339)
		}
		table.AddRow(c.Name, c.Version.String(), date, c.ShortID)
	}
	return table.String()
}

func MakeList(commits []*git.Commit, tagMap git.TagMap, fromCommit string) ComponentMap {
	// Add tags info into commit list
	for _, c := range commits {
		tag, ok := tagMap[c.ShortID]
		if ok {
			c.Tags = append(c.Tags, tag...)
		}
	}

	found := fromCommit == ""
	cMap := ComponentMap{}
	for _, c := range commits {
		if !found && c.ShortID == fromCommit {
			found = true
		}

		if len(c.Tags) == 0 || found == false {
			continue
		}

		for _, t := range c.Tags {
			last, ok := cMap[t.Name]

			if !ok || t.Version.GT(last.Version) {
				cMap[t.Name] = Component{
					Name:          t.Name,
					ShortID:       t.ShortID,
					CommittedDate: t.CommittedDate,
					Version:       t.Version,
				}
			}
		}
	}

	return cMap
}
