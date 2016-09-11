package templates

import (
	"github.com/estebarb/ion/components/reqctx"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var layout_wrong = `begin_layout

end_layout`

var layout_debug = `begin_layout
ERROR: Can't find template 'wrong'
end_layout`

var layout_hello = `begin_layout
hello
end_layout`

var layout2_hello = `begin2_layout
hello
end2_layout`

var layout_world = `begin_layout
world
end_layout`

var layout2_world = `begin2_layout
world
end2_layout`

func compare(t *testing.T, ts *httptest.Server, expected string) {
	req, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Error(err)
	}
	req.Body.Close()
	str := string(content)
	if str != expected {
		t.Errorf("Expecting '%s', got '%s'",
			expected, str)
	}
}

func doTest(t *testing.T, handler http.Handler, expected string) {
	ts := httptest.NewServer(handler)

	compare(t, ts, expected)
	ts.Close()
}

func TestTemplates_RenderWithLayout(t *testing.T) {
	tmpl := New()
	err := tmpl.LoadPattern("*.html")
	if err != nil {
		t.Error(err)
	}
	doTest(t,
		tmpl.RenderWithLayout("layout", "hello"),
		layout_hello)

	doTest(t,
		tmpl.RenderWithLayout("layout2", "hello"),
		layout2_hello)

	doTest(t,
		tmpl.RenderWithLayout("layout", "world"),
		layout_world)

	doTest(t,
		tmpl.RenderWithLayout("layout2", "world"),
		layout2_world)

	doTest(t,
		tmpl.RenderWithLayout("layout", "wrong"),
		layout_wrong)

	tmpl.Debug(true)
	doTest(t,
		tmpl.RenderWithLayout("layout", "wrong"),
		layout_debug)

	tmpl.Debug(false)
	doTest(t,
		tmpl.RenderWithLayout("layout", "wrong"),
		layout_wrong)
}

func TestTemplates_RenderTemplate(t *testing.T) {
	tmpl := New()
	err := tmpl.LoadPattern("*.html")
	if err != nil {
		t.Error(err)
	}
	doTest(t,
		tmpl.RenderTemplate("hello"),
		"hello")
}

func TestTemplates_RenderView(t *testing.T) {
	tmpl := New()
	err := tmpl.LoadPattern("*.html")
	if err != nil {
		t.Error(err)
	}
	doTest(t,
		tmpl.RenderView("hello"),
		layout_hello)

	tmpl.SetDefaultLayout("layout2")
	doTest(t,
		tmpl.RenderView("hello"),
		layout2_hello)
}

func TestTemplates_LoadPattern(t *testing.T) {
	tmpl := New()
	err := tmpl.LoadPattern("nomatch")
	if err == nil {
		t.Error("Expected error, but didn't happen")
	}
	tm := tmpl.Lookup("anything")
	if tm != nil {
		t.Error("Expected t=nil, got", tm)
	}
}

func TestTemplates_AddStateContainer(t *testing.T) {
	tmpl := New()

	container := reqctx.NewStateContainer()

	tmpl.AddStateContainer("middleware", container)

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := container.GetState(r)
			state.Set("msg", "Hello World")
			container.Middleware(next).ServeHTTP(w, r)
		})
	}

	err := tmpl.LoadPattern("*.html")
	if err != nil {
		t.Error(err)
	}

	handler := middleware(tmpl.RenderTemplate("middleware"))

	doTest(t, handler, "Hello World")
}
