package deenz_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrennmair/deenz"
)

func TestRedirectWithCatchError(t *testing.T) {
	var caughtErr error

	handler := deenz.Redirect(func(w http.ResponseWriter, r *http.Request) (loc string) {
		defer deenz.CatchError(func(err error) {
			caughtErr = err
		})

		loc = "/quux"

		deenz.HandleError(errors.New("test error"))
		return
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))

	handler.ServeHTTP(w, r)

	if expected, got := 302, w.Code; expected != got {
		t.Fatalf("Got wrong status code. Expected %d, got %d.", expected, got)
	}

	if caughtErr == nil {
		t.Fatalf("Expected to catch an error, got none.")
	}

	if expected, got := "test error", caughtErr.Error(); expected != got {
		t.Fatalf("Got wrong error message. Expected %q, got %q.", expected, got)
	}

	if expected, got := "/quux", w.Header().Get("Location"); expected != got {
		t.Fatalf("Got wrong Location header. Expected %q, got %q.", expected, got)
	}
}

func TestRedirectWithHandleErrorNil(t *testing.T) {
	var caughtErr error

	handler := deenz.Redirect(func(w http.ResponseWriter, r *http.Request) (loc string) {
		defer deenz.CatchError(func(err error) {
			caughtErr = err
		})

		loc = "/quux"

		deenz.HandleError(nil)
		return
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))

	handler.ServeHTTP(w, r)

	if caughtErr != nil {
		t.Fatalf("Expected not to catch an error, got an error: %v", caughtErr)
	}
}

func TestRedirectWithPanicThatIsNotError(t *testing.T) {
	var caughtErr error

	defer func() {
		if caughtErr != nil {
			t.Fatalf("Caught error %v where we didn't expect that.", caughtErr)
		}

		if err := recover(); err == nil {
			t.Fatalf("Expected to recover something, got %v instead", err)
		} else if expected, got := "string", fmt.Sprintf("%T", err); expected != got {
			t.Fatalf("Types don't match. Expected %q, got %q instead.", expected, got)
		}
	}()

	handler := deenz.Redirect(func(w http.ResponseWriter, r *http.Request) (loc string) {
		defer deenz.CatchError(func(err error) {
			caughtErr = err
		})

		loc = "/quux"

		panic("hello world")

		return
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))

	handler.ServeHTTP(w, r)
}
