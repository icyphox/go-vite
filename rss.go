package main

import (
	. "github.com/gorilla/feeds"
	"sort"
	"time"
	"os"
	"path/filepath"
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
		return posts[j].fm.Date.Time.Before(posts[i].fm.Date.Time)
	})

	atomfile, err := os.Create(filepath.Join("build", "blog", "feed.xml"))
	if err != nil {
		printErr(err)
	}
	for _, p := range posts {
		feed.Items = append(feed.Items, &Item{
			Title: p.fm.Title,
			Link:  &Link{Href: cfg.RSSPrefixURL + p.fm.URL},
			Description: string(p.fm.Body),
			Created: p.fm.Date.Time,
		})
	}

	err = feed.WriteAtom(atomfile)
	if err != nil {
		printErr(err)
	}
}
