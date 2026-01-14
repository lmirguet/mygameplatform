package httpx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSPA_ServesIndexOnRoot(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	s := NewSPA(http.Dir(dir), "index.html")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)
	s.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "INDEX" {
		t.Fatalf("body: got %q want %q", body, "INDEX")
	}
}

func TestSPA_ServesFileWhenExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "asset.txt"), []byte("ASSET"), 0o644); err != nil {
		t.Fatalf("write asset.txt: %v", err)
	}

	s := NewSPA(http.Dir(dir), "index.html")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/asset.txt", nil)
	s.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "ASSET" {
		t.Fatalf("body: got %q want %q", body, "ASSET")
	}
}

func TestSPA_FallsBackToIndexOnUnknownPath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	s := NewSPA(http.Dir(dir), "index.html")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/some/deep/link", nil)
	s.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "INDEX" {
		t.Fatalf("body: got %q want %q", body, "INDEX")
	}
}

func TestSPA_NonGETIsNotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	s := NewSPA(http.Dir(dir), "index.html")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "http://example.test/", strings.NewReader("x"))
	s.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status: got %d want %d", rr.Code, http.StatusNotFound)
	}
	b, _ := io.ReadAll(rr.Result().Body)
	if len(b) == 0 {
		t.Fatalf("expected non-empty 404 body")
	}
}

