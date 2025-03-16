package deenz

import (
	"html/template"
	"net/http"
)

type tmplParams[T any] struct {
	Error  string
	Values *T
}

func Render[T any](tmpl *template.Template, f func(http.ResponseWriter, *http.Request) *T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tmpl == nil {
			http.Error(w, "provided template is nil", http.StatusInternalServerError)
			return
		}

		params := &tmplParams[T]{}

		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					params.Error = e.Error()
				}
			}

			if params.Values != nil || params.Error != "" {
				_ = tmpl.Execute(w, params)
			}
		}()

		params.Values = f(w, r)
	})
}
