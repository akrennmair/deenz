package deenz_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrennmair/deenz"
)

func TestRedirect(t *testing.T) {
	handler := deenz.Redirect(func(w http.ResponseWriter, r *http.Request) string {
		return "/foobar"
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))

	handler.ServeHTTP(w, r)

	if expected, got := 302, w.Code; expected != got {
		t.Fatalf("Got wrong status code. Expected %d, got %d.", expected, got)
	}

	if expected, got := "/foobar", w.Header().Get("Location"); expected != got {
		t.Fatalf("Got wrong Location header. Expected %q, got %q.", expected, got)
	}
}

func TestRedirectError(t *testing.T) {
	handler := deenz.Redirect(func(w http.ResponseWriter, r *http.Request) string {
		deenz.HandleError(errors.New("test error"))
		return "/foobar"
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))

	handler.ServeHTTP(w, r)

	if expected, got := http.StatusInternalServerError, w.Code; expected != got {
		t.Fatalf("Got wrong status code. Expected %d, got %d.", expected, got)
	}

	if expected, got := "", w.Header().Get("Location"); expected != got {
		t.Fatalf("Got wrong Location header. Expected %q, got %q.", expected, got)
	}
}
