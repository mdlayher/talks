package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// START TEST OMIT
func TestClientStatusOK(t *testing.T) {
	const expect = "ok"
	c, done := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if want, got := "/status", r.URL.Path; want != got {
			t.Fatalf("unexpected URL path:\n- want: %q\n-  got: %q", want, got)
		}

		_, _ = io.WriteString(w, expect)
	})
	defer done()

	status, err := c.Status()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if want, got := expect, status; want != got {
		t.Fatalf("unexpected status:\n- want: %q\n-  got: %q", want, got)
	}
}

// END TEST OMIT

// START TESTCLIENT OMIT
func testClient(t *testing.T, fn func(w http.ResponseWriter, r *http.Request)) (*Client, func()) {
	s := httptest.NewServer(http.HandlerFunc(fn))

	c := NewClient(s.URL)
	done := func() { s.Close() }

	return c, done
}

// END TESTCLIENT OMIT
