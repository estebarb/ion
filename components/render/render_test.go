package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tm := New()
	tm.AddFunc("upper", strings.ToUpper)

	tm.Push("wrapperA.html")
	tm.AddTemplate("contentA.html")
	tm.AddTemplate("contentB.html")

	var bufA bytes.Buffer
	var bufB bytes.Buffer
	var bufC bytes.Buffer
	var bufD bytes.Buffer
	var bufE bytes.Buffer

	t.Logf("Loaded %v templates:", len(tm.templates))
	for k := range tm.templates {
		t.Logf("Found template: %v", k)
	}

	err := tm.Execute(&bufA, "base", "contentA", "hello")
	if err != nil {
		t.Errorf("For base, contentA, hello: %v", err)
	}

	err = tm.Execute(&bufB, "base", "contentB", "hello")
	if err != nil {
		t.Errorf("For base2, contentB, hello: %v", err)
	}

	err = tm.Execute(&bufC, "base2", "contentA", "hello")
	if err != nil {
		t.Errorf("For base, contentA, hello: %v", err)
	}

	err = tm.Execute(&bufD, "base2", "contentB", "hello")
	if err != nil {
		t.Errorf("For base2, contentB, hello: %v", err)
	}

	err = tm.Execute(&bufE, "base2", "not here", "hello")
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}


	expectedA := "<wrapper><content>hello</content></wrapper>"
	expectedB := "<wrapper><contentB>HELLO</contentB></wrapper>"
	expectedC := "<wrapper2><content>hello</content></wrapper2>"
	expectedD := "<wrapper2><contentB>HELLO</contentB></wrapper2>"

	if bufA.String() != expectedA {
		t.Errorf("Expecting `%s`, got: `%s`", expectedA, bufA.String())
	}

	if bufB.String() != expectedB {
		t.Errorf("Expecting `%s`, got: `%s`", expectedB, bufB.String())
	}

	if bufC.String() != expectedC {
		t.Errorf("Expecting `%s`, got: `%s`", expectedC, bufC.String())
	}

	if bufD.String() != expectedD {
		t.Errorf("Expecting `%s`, got: `%s`", expectedD, bufD.String())
	}
}
