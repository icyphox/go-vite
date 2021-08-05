package atom

import (
	"encoding/xml"
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

type AtomEntry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	ID      string   `xml:"id"`
	Link    *AtomLink
	Summary *AtomSummary
}

type AtomFeed struct {
	XMLName  xml.Name `xml:"feed"`
	Xmlns    string   `xml:"xmlns,attr"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
	ID       string   `xml:"id"`
	Updated  string   `xml:"updated"`
	Link     *AtomLink
	Author   *AtomAuthor `xml:"author"`
	Entries  []AtomEntry
}

// Creates a new Atom feed.
func NewAtomFeed(srcDir string, posts []markdown.Output) ([]byte, error) {
	entries := []AtomEntry{}
	config := config.Config
	for _, p := range posts {
		dateStr := p.Meta["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		rfc3339 := date.Format(time.RFC3339)

		entry := AtomEntry{
			Title:   p.Meta["title"],
			Updated: rfc3339,
			ID:      NewUUID().String(),
			Link:    &AtomLink{Href: config.URL + srcDir + p.Meta["slug"]},
			Summary: &AtomSummary{Content: string(p.HTML), Type: "html"},
		}
		entries = append(entries, entry)
	}

	// 2021-07-14T00:00:00Z
	now := time.Now().Format(time.RFC3339)
	feed := &AtomFeed{
		Xmlns:    "http://www.w3.org/2005/Atom",
		Title:    config.Title,
		ID:       config.URL,
		Subtitle: config.Desc,
		Link:     &AtomLink{Href: config.URL},
		Author: &AtomAuthor{
			Name:  config.Author.Name,
			Email: config.Author.Email,
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
