package main

import (
	"html/template"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/tools/present"
)

//go:generate bash -c "go run main.go > README.md"

// prefix is the base URL for the go-talks website with this repo's
// name prefixed.
const prefix = "http://go-talks.appspot.com/github.com/mdlayher/talks/"

func main() {
	base, err := url.Parse(prefix)
	if err != nil {
		log.Fatalf("failed to parse base URL: %v", err)
	}

	// Look for all presentations in '.slide' format.
	var ps []presentation
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".slide" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Parse a presentation to retrieve its metadata
		doc, err := present.Parse(f, path, 0)
		if err != nil {
			return err
		}

		u, err := url.Parse(path)
		if err != nil {
			return err
		}

		// Build a presentation to be output as part of the template
		ps = append(ps, presentation{
			Title:    doc.Title,
			Subtitle: doc.Subtitle,
			Date:     doc.Time.Format("2006 Jan 02"),
			Time:     doc.Time,
			URL:      base.ResolveReference(u).String(),
		})
		return nil
	})
	if err != nil {
		log.Fatalf("unexpected error during filesystem walk: %v", err)
	}

	// Order all presentations by their date and time, and then output the template
	sort.Sort(byTime(ps))
	_ = markdown.Execute(os.Stdout, ps)
}

// byTime sorts a slice of presentations by their Time field.
type byTime []presentation

func (t byTime) Len() int           { return len(t) }
func (t byTime) Less(i, j int) bool { return t[i].Time.Before(t[j].Time) }
func (t byTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// A presentation is a present-format presentation.
type presentation struct {
	Title    string
	Subtitle string
	Time     time.Time
	Date     string
	URL      string
}

// markdown is the markdown template for README.md
var markdown = template.Must(template.New("README.md").Parse(strings.TrimSpace(`
talks [![Build Status](https://travis-ci.org/mdlayher/talks.svg?branch=master)](https://travis-ci.org/mdlayher/talks)
=====

Talks by Matt Layher using the Go present tool.  MIT Licensed.

Talks
-----
{{range .}}- [{{.Date}} - {{.Title}}]({{.URL}}){{if .Subtitle}}
  - {{.Subtitle}}{{end}}
{{end}}
`)))
