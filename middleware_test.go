package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestWrapper(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	expected := []string{
		"First",
		"Second",
		"Third",
		"Handler",
		"Third-to-Last",
		"Second-to-Last",
		"Last",
	}

	t.Run("Append Test", func(t *testing.T) {
		var actual []string
		w := &Wrapper{}

		w.Append(getFakeMiddleware(&actual, "First", "Last"))
		w.Append(getFakeMiddleware(&actual, "Second", "Second-to-Last"))
		w.Append(getFakeMiddleware(&actual, "Third", "Third-to-Last"))
		w.WrapHandler(getFakeHandler(&actual))(res, req)

		if !reflect.DeepEqual(expected, actual) {
			t.Error("Expected:", stringify(expected), "| Actual:", stringify(actual))
		}
	})

	t.Run("Prepend Test", func(t *testing.T) {
		var actual []string
		w := &Wrapper{}

		w.Prepend(getFakeMiddleware(&actual, "Third", "Third-to-Last"))
		w.Prepend(getFakeMiddleware(&actual, "Second", "Second-to-Last"))
		w.Prepend(getFakeMiddleware(&actual, "First", "Last"))
		w.WrapHandler(getFakeHandler(&actual))(res, req)

		if !reflect.DeepEqual(expected, actual) {
			t.Error("Expected:", stringify(expected), "| Actual:", stringify(actual))
		}
	})

	t.Run("Mixed With Chaining Syntax", func(t *testing.T) {
		var actual []string
		w := &Wrapper{}

		w.
			Prepend(getFakeMiddleware(&actual, "Second", "Second-to-Last")).
			Append(getFakeMiddleware(&actual, "Third", "Third-to-Last")).
			Prepend(getFakeMiddleware(&actual, "First", "Last")).
			WrapHandler(getFakeHandler(&actual))(res, req)

		if !reflect.DeepEqual(expected, actual) {
			t.Error("Expected:", stringify(expected), "| Actual:", stringify(actual))
		}
	})
}

// ==================== Test Fixtures/Props =====================

func stringify(s []string) string {
	return strings.Join(s, ",")
}

func getFakeHandler(actual *[]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*actual = append(*actual, "Handler")
	}
}

func getFakeMiddleware(actual *[]string, before, after string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*actual = append(*actual, before)
			next(w, r)
			*actual = append(*actual, after)
		}
	}
}
