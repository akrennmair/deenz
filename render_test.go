package deenz_test

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrennmair/deenz"
)

type Params struct {
	Name string
}

func TestRender(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`Hello, {{ .Values.Name }}!`))

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		return &Params{Name: "world"}
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := `Hello, world!`, w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}

func TestRenderNilValues(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`Hello, {{ .Values.Name }}!`))

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		return nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := ``, w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}

func TestRenderNilTemplate(t *testing.T) {
	handler := deenz.Render(nil, func(w http.ResponseWriter, r *http.Request) *Params {
		return &Params{Name: "world"}
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := "provided template is nil\n", w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}

func TestRenderError(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`{{ if .Error }}{{ .Error }}{{ else }}Hello, {{ .Values.Name }}!{{ end }}`))

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		w.WriteHeader(http.StatusInternalServerError)
		deenz.HandleError(errors.New("test error"))
		return &Params{Name: "world"}
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := http.StatusInternalServerError, w.Code; expected != got {
		t.Fatalf("Got wrong HTTP status code. Expected %d, got %d instead.", expected, got)
	}

	if expected, got := `test error`, w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}

func TestRenderErrorThatIsWrongType(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`{{ if .Error }}{{ .Error }}{{ else }}Hello, {{ .Values.Name }}!{{ end }}`))

	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("Expected to recover something, got %v instead", err)
		} else if expected, got := "string", fmt.Sprintf("%T", err); expected != got {
			t.Fatalf("Types don't match. Expected %q, got %q instead.", expected, got)
		}
	}()

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		w.WriteHeader(http.StatusInternalServerError)
		panic("not an actual error")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)
}

func TestRenderMust(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`{{ if .Error }}{{ .Error }}{{ else }}Hello, {{ .Values.Name }}!{{ end }}`))

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		return deenz.Must(func() (*Params, error) {
			return &Params{Name: "world"}, nil
		}())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := `Hello, world!`, w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}

func TestRenderMustError(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`{{ if .Error }}{{ .Error }}{{ else }}Hello, {{ .Values.Name }}!{{ end }}`))

	handler := deenz.Render(tmpl, func(w http.ResponseWriter, r *http.Request) *Params {
		return deenz.Must(func() (*Params, error) {
			return nil, errors.New("test error")
		}())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)

	if expected, got := `test error`, w.Body.String(); expected != got {
		t.Fatalf("Got wrong render result. Expected %q, got %q instead.", expected, got)
	}
}
