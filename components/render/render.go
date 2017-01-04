// Package render provides a little sugar to the usage of standard Go templates.
// It allows to easily load templates and to setup functions and "inheritance".
// Based on ideas taken from https://elithrar.github.io/article/approximating-html-template-inheritance/
package render

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Templates facilitates usage of templates with inheritance
type Templates struct {
	stack     []*template.Template
	funcs     template.FuncMap
	templates map[string]*template.Template
}

// New creates a new template manager
func New() *Templates {
	return &Templates{
		stack:     make([]*template.Template, 0),
		funcs:     make(map[string]interface{}),
		templates: make(map[string]*template.Template),
	}
}

// Push adds a set of files (path is a regex) to the stack
// of "base templates". If err is non-nil then the stack is left
// untouched. The new template is based on the previous top of the stack.
func (t *Templates) Push(path string) (err error) {
	base, err := t.Peek().Clone()
	if err != nil {
		return err
	}
	newBase, err := base.ParseGlob(path)
	if err != nil {
		return err
	}
	t.stack = append(t.stack, newBase)
	return nil
}

// Peek returns the current top of the "base template" stack
func (t *Templates) Peek() *template.Template {
	if len(t.stack) > 0 {
		return t.stack[len(t.stack)-1]
	}
	return template.New("").Funcs(t.funcs)
}

// Pop removes a "base template" from the stack
func (t *Templates) Pop() {
	t.stack = t.stack[:len(t.stack)-1]
}

// AddTemplate creates a new template, with the same name as the
// file (without extension). The template is based on the current
// top of the stack.
func (t *Templates) AddTemplate(fileglob string) error {
	if t.templates == nil {
		t.templates = make(map[string]*template.Template)
	}

	layouts, err := filepath.Glob(fileglob)
	if err != nil {
		return err
	}

	for _, layout := range layouts {
		base, err := t.Peek().Clone()
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.Base(layout),
			filepath.Ext(layout))
		newTemplate, err := base.ParseFiles(layout)
		if err != nil {
			return err
		}
		t.templates[name] = newTemplate
	}
	return nil
}

// AddFunc adds a new function to the map of functions. All the functions
// must be added before adding templates or base layouts.
func (t *Templates) AddFunc(name string, fun interface{}) {
	t.funcs[name] = fun
}

// Execute renders the template with the given name, using the base layout, io.Writer
// and data given.
func (t *Templates) Execute(w io.Writer, base, name string, data interface{}) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	return tmpl.ExecuteTemplate(w, base, data)
}
