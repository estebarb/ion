// Auxiliar functions to help in the construction of
// templates and handlers that use them
package templates

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/estebarb/ion/components/reqctx"
	"html/template"
	"net/http"
)

var ErrUnableToFindLayout = errors.New("Unable to find requested layout")

type Templates struct {
	debug         bool
	functionsMap  template.FuncMap
	templates     *template.Template
	states        *reqctx.State
	defaultLayout string
}

type ITemplateCtx interface {
	SetLayout(name string)
	Layout() string
	SetView(name string)
	View() string
	TemplateValues() map[string]interface{}
	SetTemplateValue(key string, value interface{})
	TemplateValue(key string) (interface{}, bool)
}

type TemplateCtx struct {
	LayoutName string
	ViewName   string
	Values     map[string]interface{}
}

func (t *TemplateCtx) SetLayout(name string) {
	t.LayoutName = name
}
func (t *TemplateCtx) Layout() string {
	return t.LayoutName
}
func (t *TemplateCtx) SetView(name string) {
	t.ViewName = name
}
func (t *TemplateCtx) View() string {
	return t.ViewName
}
func (t *TemplateCtx) TemplateValues() map[string]interface{} {
	return t.Values
}
func (t *TemplateCtx) SetTemplateValue(key string, value interface{}) {
	t.Values[key] = value
}
func (t *TemplateCtx) TemplateValue(key string) (interface{}, bool) {
	v, ok := t.Values[key]
	return v, ok
}

func ContextFactory() interface{} {
	return &TemplateCtx{
		Values: make(map[string]interface{}),
	}
}

// Creates a new template that has a default layout ("layout")
// and some additional functions:
// render <name> <ctx>  -> Renders other template
func New(state *reqctx.State) *Templates {
	t := &Templates{
		functionsMap:  make(template.FuncMap),
		states:        state,
		defaultLayout: "layout",
	}
	t.AddFunc("render", t.dispatchTemplate)
	return t
}

func (t *Templates) Debug(debug bool) {
	t.debug = debug
}

func (t *Templates) SetDefaultLayout(name string) {
	t.defaultLayout = name
}

func (t *Templates) AddFunc(name string, fun interface{}) {
	t.functionsMap[name] = fun
}

func (t *Templates) LoadPattern(pattern string) error {
	tl := template.New("").Funcs(t.functionsMap)
	tl, err := tl.ParseGlob(pattern)
	if err != nil {
		return err
	}
	t.templates = tl
	return nil
}

func (t *Templates) Lookup(name string) *template.Template {
	if t.templates != nil {
		return t.templates.Lookup(name)
	}
	return nil
}

func (t *Templates) execute(name string, env interface{}) string {
	var buf bytes.Buffer
	tmpl := t.Lookup(name)
	if tmpl != nil {
		tmpl.Execute(&buf, env)
		return buf.String()
	}
	if t.debug {
		return fmt.Sprintf("ERROR: Can't find template '%s'", name)
	} else {
		return ""
	}
}

func (t *Templates) BuildContext(r *http.Request) interface{} {
	return t.states.Context(r)
}

func (t *Templates) RenderTemplate(tmpl string) http.Handler {
	return t.RenderWithLayout(tmpl, tmpl)
}

func (t *Templates) RenderWithLayout(layout, view string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := t.BuildContext(r).(ITemplateCtx)
		ctx.SetView(view)
		view := t.Lookup(layout)
		if view != nil {
			view.Execute(w, ctx)
		} else {
			panic(ErrUnableToFindLayout)
		}
	})
}

func (t *Templates) RenderView(view string) http.Handler {
	return t.RenderWithLayout(t.defaultLayout, view)
}

func (t *Templates) dispatchTemplate(name string, env interface{}) template.HTML {
	return template.HTML(t.execute(name, env))
}
