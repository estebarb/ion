// Auxiliar functions to help in the construction of
// templates and handlers that use them
package templates

import (
	"bytes"
	"fmt"
	"github.com/estebarb/ion/components/reqctx"
	"html/template"
	"net/http"
)

type Templates struct {
	debug         bool
	functionsMap  template.FuncMap
	templates     *template.Template
	states        map[string]*reqctx.StateContainer
	defaultLayout string
}

// Creates a new template that has a default layout ("layout")
// and some additional functions:
// render <name> <ctx>  -> Renders other template
func New() *Templates {
	t := &Templates{
		functionsMap:  make(template.FuncMap),
		states:        make(map[string]*reqctx.StateContainer),
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

func (t *Templates) AddStateContainer(namespace string, state *reqctx.StateContainer) {
	t.states[namespace] = state
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

func (t *Templates) BuildContext(r *http.Request) map[string]interface{} {
	ctx := make(map[string]interface{})
	for namespace, stateContainer := range t.states {
		ctx[namespace] = stateContainer.GetState(r).GetAll()
	}
	return ctx
}

func (t *Templates) RenderTemplate(tmpl string) http.Handler {
	return t.RenderWithLayout(tmpl, tmpl)
}

func (t *Templates) RenderWithLayout(layout, view string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := t.BuildContext(r)
		ctx["__view__"] = view
		t.Lookup(layout).Execute(w, ctx)
	})
}

func (t *Templates) RenderView(view string) http.Handler {
	return t.RenderWithLayout(t.defaultLayout, view)
}

func (t *Templates) dispatchTemplate(name string, env interface{}) template.HTML {
	return template.HTML(t.execute(name, env))
}
