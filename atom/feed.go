package atom

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"git.icyphox.sh/vite/config"
	"git.icyphox.sh/vite/markdown"
)

type AtomLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr,omitempty"`
}

type AtomSummary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type AtomAuthor struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email"`
}

type AtomTitle struct {
	XMLName xml.Name `xml:"title"`
	Title   string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type AtomEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   *AtomTitle
	Updated string `xml:"updated"`
	ID      string `xml:"id"`
	Link    *AtomLink
	Summary *AtomSummary
}

type AtomFeed struct {
	XMLName  xml.Name `xml:"feed"`
	Xmlns    string   `xml:"xmlns,attr"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle,omitempty"`
	ID       string   `xml:"id"`
	Updated  string   `xml:"updated"`
	Link     *AtomLink
	Author   *AtomAuthor `xml:"author"`
	Entries  []AtomEntry
}

// Creates a new Atom feed.
func NewAtomFeed(srcDir string, posts []markdown.Output) ([]byte, error) {
	entries := []AtomEntry{}

	for _, p := range posts {
		dateStr := p.Meta["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		rfc3339 := date.Format(time.RFC3339)

		entry := AtomEntry{
			Title:   &AtomTitle{Title: p.Meta["title"], Type: "html"},
			Updated: rfc3339,
			// tag:icyphox.sh,2019-10-23:blog/some-post/
			ID: fmt.Sprintf(
				"tag:%s,%s:%s",
				strings.TrimSuffix(
					config.Config.URL[8:], // strip https://
					"/",
				),
				dateStr,
				filepath.Join(srcDir, p.Meta["slug"]),
			),
			// filepath.Join strips the second / in http://
			Link: &AtomLink{Href: config.Config.URL + filepath.Join(srcDir, p.Meta["slug"])},
			Summary: &AtomSummary{
				Content: fmt.Sprintf("<h2>%s</h2>\n%s",
					p.Meta["subtitle"],
					string(p.HTML)),
				Type: "html",
			},
		}
		entries = append(entries, entry)
	}

	// 2021-07-14T00:00:00Z
	now := time.Now().Format(time.RFC3339)
	feed := &AtomFeed{
		Xmlns:    "http://www.w3.org/2005/Atom",
		Title:    config.Config.Title,
		ID:       config.Config.URL,
		Subtitle: config.Config.Desc,
		Link:     &AtomLink{Href: config.Config.URL},
		Author: &AtomAuthor{
			Name:  config.Config.Author.Name,
			Email: config.Config.Author.Email,
		},
		Updated: now,
		Entries: entries,
	}

	feedXML, err := xml.MarshalIndent(feed, " ", " ")
	if err != nil {
		return nil, err
	}
	// Add the <?xml...> header.
	return []byte(xml.Header + string(feedXML)), nil
}
