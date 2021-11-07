package ui

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	fsx "github.com/bmatcuk/doublestar/v4"
)

//go:embed html
var htmlFS embed.FS

type Renderer interface {
	Render(w io.Writer, name string, data interface{}) error
}

type tmpl struct {
	templates     map[string]*template.Template
	baseTemplates []string
	partials      []string
	pages         []string
}

func (t *tmpl) Render(w io.Writer, name string, data interface{}) error {
	//t.parseTemplates()
	buf := &bytes.Buffer{}
	err := t.templates[name].ExecuteTemplate(buf, "base", data)

	if err != nil {
		fmt.Fprintf(w, "An error occured")
		return err
	}
	_, err = buf.WriteTo(w)

	return err
}

func NewRenderer() Renderer {
	t := tmpl{
		templates:     make(map[string]*template.Template),
		pages:         make([]string, 0),
		baseTemplates: make([]string, 0),
		partials:      make([]string, 0),
	}

	if err := t.parseTemplates(); err != nil {
		panic(err)
	}

	return &t
}

func (t *tmpl) parseTemplates() error {
	if err := t.extractPages(htmlFS); err != nil {
		return err
	}

	for _, page := range t.pages {
		name := filepath.Base(page)
		files := t.addBaseTemplatesAndPartialsToPages(page)
		temp := template.Must(template.ParseFS(htmlFS, files...))

		t.templates[name] = temp
	}

	return nil
}

func (t *tmpl) addBaseTemplatesAndPartialsToPages(page string) []string {
	files := make([]string, len(t.baseTemplates))
	copy(files, t.baseTemplates)
	files = append(files, page)
	files = append(files, t.partials...)
	return files
}

func (t *tmpl) extractPages(fsys fs.FS) error {
	matches, err := fsx.Glob(fsys, "html/**/*.html")

	if err != nil {
		return err
	}

	for _, match := range matches {
		if strings.HasSuffix(match, "base.html") {
			// add to base files
			t.baseTemplates = append(t.baseTemplates, match)
			continue
		}

		if strings.Contains(match, "partials") {
			// add to partials
			t.partials = append(t.partials, match)
			continue
		}

		t.pages = append(t.pages, match)
	}

	return nil
}
