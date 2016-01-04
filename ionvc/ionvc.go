// This packages provides the View and Controller helpers
// for implementing a MVC style programing.
// With IonVC you create controllers that put the data required
// by the views in the context, and the framework will inject this
// data on the views
// Also provides a centralized way to load templates, and to inject
// functions on them.
package ionvc
import (
	"html/template"
	"net/http"
	"fmt"
	"github.com/estebarb/ion"
)

var FunctionsMap template.FuncMap
var _templates *template.Template

func init(){
	FunctionsMap = make(template.FuncMap)
}

// Adds a function to the FunctionMap
func AddFunc(name string, fun interface{}){
	FunctionsMap[name] = fun
}

func LoadTemplates(pattern string){
	tl := template.New("").Funcs(FunctionsMap)
	tl, err := tl.ParseGlob(pattern)
	if err != nil{
		panic(fmt.Sprintf("Can't load templates: %v", err))
	}
	_templates = tl
}

// Wraps the handler given, and after it executes
// IonVC will render the corresponding template.
func ControllerFunc(handler http.HandlerFunc, templateName string) http.HandlerFunc{
	tmpl := _templates.Lookup(templateName)
	if tmpl == nil{
		panic(fmt.Sprintf("Template '%v' not found", templateName))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		ion.RenderTemplate(tmpl).ServeHTTP(w, r)
	}
}

func TemplateLookup(templateName string) *template.Template {
    return _templates.Lookup(templateName)
}
