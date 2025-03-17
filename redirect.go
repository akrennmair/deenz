package deenz

import "net/http"

func Redirect(f func(http.ResponseWriter, *http.Request) string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					http.Error(w, e.Error(), http.StatusInternalServerError)
				} else {
					panic(err)
				}
			}
		}()

		location := f(w, r)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusFound)
	})
}
