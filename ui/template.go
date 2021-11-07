package ui

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	fsx "github.com/bmatcuk/doublestar/v4"
)

//go:embed html
var htmlFS embed.FS

type Renderer interface {
	Render(w io.Writer, name string, data interface{}) error
	RegisterExtraContext(funcs ...func() map[string]interface{})
}

type tmpl struct {
	templates     map[string]*template.Template
	baseTemplates []string
	partials      []string
	pages         []string
	context       map[string]interface{}
	template.FuncMap
}

func (t *tmpl) Render(w io.Writer, name string, data interface{}) error {
	// t.parseTemplates()
	buf := &bytes.Buffer{}
	var passedContext interface{}

	if data == nil {
		passedContext = t.context
	} else if dataMap, ok := data.(map[string]interface{}); ok {
		for k, v := range dataMap {
			t.context[k] = v
		}
		passedContext = t.context
	} else {
		passedContext = data
	}

	err := t.templates[name].ExecuteTemplate(buf, "base", passedContext)

	if err != nil {
		fmt.Fprintf(w, "An error occured")
		log.Println(err)
		return err
	}
	_, err = buf.WriteTo(w)

	return err
}

func (t *tmpl) RegisterExtraContext(funcs ...func() map[string]interface{}) {
	for _, f := range funcs {
		res := f()
		for key, val := range res {
			t.context[key] = val
		}
	}
}

func NewRenderer() Renderer {
	t := tmpl{
		templates:     make(map[string]*template.Template),
		pages:         make([]string, 0),
		baseTemplates: make([]string, 0),
		partials:      make([]string, 0),
		context:       make(map[string]interface{}),
		FuncMap: template.FuncMap{
			"isProd": isProd,
		},
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
		temp := &template.Template{}
		temp = temp.Funcs(t.FuncMap)
		temp = template.Must(temp.ParseFS(htmlFS, files...))

		//temp := template.Must(template.ParseFiles(files...)) // UNCOMMENT WHEN USING YOUR ACTUAL FILESYSTEM

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

func readDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var files []string
	for _, de := range entries {
		name := fmt.Sprintf("%s/%s", dir, de.Name())
		if !de.IsDir() {
			files = append(files, name)
		} else {
			fs, err := readDir(name)

			if err != nil {
				return nil, err
			}

			files = append(files, fs...)
		}
	}
	return files, nil
}

func (t *tmpl) extractPages(fsys fs.FS) error {
	matches, err := fsx.Glob(fsys, "html/**/*.html")

	//matches, err := readDir("./ui/html")

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

func isProd() bool {
	return os.Getenv("RELEASE") == "production"
}
