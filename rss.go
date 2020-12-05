package main

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	. "github.com/gorilla/feeds"
)

func generateRSS(posts []Post, cfg Config) {
	now := time.Now()
	feed := &Feed{
		Title:       cfg.Title,
		Link:        &Link{Href: cfg.SiteURL},
		Description: cfg.Description,
		Author: &Author{
			Name:  cfg.Author["name"],
			Email: cfg.Author["email"],
		},
		Created: now,
	}

	// Sort posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[j].Fm.Date.Time.Before(posts[i].Fm.Date.Time)
	})

	atomfile, err := os.Create(filepath.Join("build", "blog", "feed.xml"))
	if err != nil {
		printErr(err)
	}
	for _, p := range posts {
		feed.Items = append(feed.Items, &Item{
			Title:       p.Fm.Title,
			Link:        &Link{Href: cfg.RSSPrefixURL + p.Fm.URL},
			Description: string(p.Fm.Body),
			Created:     p.Fm.Date.Time,
		})
	}

	err = feed.WriteAtom(atomfile)
	if err != nil {
		printErr(err)
	}
}
