package main

import (
	"encoding/json"
	"fmt"
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

//go:generate go run main.go

const (
	// talksPrefix is the base URL for the talks.godoc.org website with this
	// repository's name prefixed.
	talksPrefix = "https://talks.godoc.org/github.com/mdlayher/talks/"

	// repoPrefix is the base URL for resources whose links are hosted in
	// this repository.
	repoPrefix = "https://github.com/mdlayher/talks/blob/master/"

	// talksJSON is the name of the JSON metadata file produced by this script.
	talksJSON = "talks.json"
)

func main() {
	base, err := url.Parse(talksPrefix)
	if err != nil {
		log.Fatalf("failed to parse base URL: %v", err)
	}

	// Look for all presentations in '.slide' format.
	var ps []*presentation
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var fn func(path string, base *url.URL) (*presentation, error)
		switch filepath.Ext(path) {
		case ".json":
			// Skip special metadata file.
			if path == talksJSON {
				return nil
			}

			fn = parseJSON
		case ".slide":
			// Don't parse slides if a JSON file exists instead.
			noExt := strings.TrimSuffix(path, filepath.Ext(path))
			if _, err := os.Stat(noExt + ".json"); err == nil {
				return nil
			}

			fn = parsePresentation
		default:
			return nil
		}

		p, err := fn(path, base)
		if err != nil {
			return err
		}

		// Ensure valid resources.
		for i, r := range p.Resources {
			switch r.Kind {
			case audio, blog, slides:
			default:
				log.Fatalf("unexpected resource kind for %q: %q", p.Time, r.Kind)
			}

			if r.Link == "" {
				log.Fatalf("empty resource link for %q, kind %q", p.Title, r.Kind)
			}

			// Post-processing for URLs relative to this repository.
			link, err := resolveLink(repoPrefix, r.Link)
			if err != nil {
				log.Fatalf("failed to resolve resource link for %q: %v", p.Title, err)
			}

			p.Resources[i].Link = link
		}

		// Create a README for individual presentations, placing it at the same
		// directory level as the presentation metadata.
		readme, err := os.Create(filepath.Join(filepath.Dir(path), "README.md"))
		if err != nil {
			log.Fatalf("failed to create directory README: %v", err)
		}
		defer readme.Close()

		input := input{
			Title:         p.Title,
			Description:   p.Description,
			VideoLink:     p.VideoLink,
			ResourcesList: markdownList(p.Resources),
		}

		if err := directory.Execute(readme, input); err != nil {
			log.Fatalf("failed to execute directory template: %v", err)
		}

		ps = append(ps, p)
		return nil
	})
	if err != nil {
		log.Fatalf("unexpected error during filesystem walk: %v", err)
	}

	// Order all presentations by latest date and time, and then output the template.
	sort.Slice(ps, func(i int, j int) bool {
		return ps[i].Time.After(ps[j].Time)
	})

	// Generate top-level README.md.
	readme, err := os.Create("README.md")
	if err != nil {
		log.Fatalf("failed to create index README: %v", err)
	}
	defer readme.Close()

	// Render the presentations in an input format suitable for the README
	// markdown template.
	inputs := make([]input, 0, len(ps))
	for _, p := range ps {
		inputs = append(inputs, input{
			Title:         p.Title,
			Description:   p.Description,
			VideoLink:     p.VideoLink,
			ResourcesList: markdownList(p.Resources),
		})
	}

	if err := index.Execute(readme, inputs); err != nil {
		log.Fatalf("failed to execute index template: %v", err)
	}

	// Generate talks.json metadata.
	talks, err := os.Create(talksJSON)
	if err != nil {
		log.Fatalf("failed to create talks.json: %v", err)
	}
	defer talks.Close()

	if err := json.NewEncoder(talks).Encode(ps); err != nil {
		log.Fatalf("failed to encode JSON: %v", err)
	}
}

func parseJSON(path string, base *url.URL) (*presentation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var p presentation
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func parsePresentation(path string, base *url.URL) (*presentation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Parse a presentation to retrieve its metadata.
	doc, err := present.Parse(f, path, 0)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Build a presentation to be output as part of the template.
	return &presentation{
		Title:       doc.Title,
		Description: doc.Subtitle,
		Time:        doc.Time,
		Resources: []resource{{
			Kind: slides,
			Link: base.ResolveReference(u).String(),
		}},
	}, nil
}

// A presentation is a presentation's metadata.
type presentation struct {
	Title       string
	Description string
	Time        time.Time
	VideoLink   string
	Resources   []resource
}

// A resource is a type of external content resource.
type resource struct {
	Kind kind
	Link string
}

// A kind is a resource type.
//
// Video is explicitly not a kind because it is formatted differently in output.
type kind string

// Known kind types. Any unrecognized type will result in an error.
const (
	audio  kind = "audio"
	blog   kind = "blog"
	slides kind = "slides"
)

// An input is an input for the README templates.
type input struct {
	Title         string
	Description   string
	VideoLink     string
	ResourcesList string
}

// markdownList generates a markdown-formatted string of resource links.
func markdownList(resources []resource) string {
	ss := make([]string, 0, len(resources))
	for _, r := range resources {
		ss = append(ss, fmt.Sprintf("[%s](%s)", r.Kind, r.Link))
	}

	return strings.Join(ss, ", ")
}

func resolveLink(prefix, rel string) (string, error) {
	pu, err := url.Parse(prefix)
	if err != nil {
		return "", err
	}

	relu, err := url.Parse(rel)
	if err != nil {
		return "", err
	}

	return pu.ResolveReference(relu).String(), nil
}

// index is the markdown template for the top-level README.md.
var index = template.Must(template.New("index.md").Parse(strings.TrimSpace(`
# talks [![Build Status](https://travis-ci.org/mdlayher/talks.svg?branch=master)](https://travis-ci.org/mdlayher/talks)

Talks by Matt Layher. MIT Licensed.

## Talks
{{range .}}
- {{if .VideoLink}}[{{.Title}}]({{.VideoLink}}){{else}}{{.Title}}{{end}}{{if .Description}}
  - {{.Description}}{{end}}{{if .ResourcesList}}
  - {{.ResourcesList}}{{end}}{{end}}
`)))

// directory is the markdown template for individual directory README.md files.
var directory = template.Must(template.New("directory.md").Parse(strings.TrimSpace(`
# {{if .VideoLink}}[{{.Title}}]({{.VideoLink}}){{else}}{{.Title}}{{end}}

{{.Description}}
{{if .ResourcesList}}
- {{.ResourcesList}}
{{end}}
`)))
