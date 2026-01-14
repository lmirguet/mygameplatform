package httpx

import (
	"net/http"
	"path"
)

// SPA serves static assets from fs and falls back to serving indexFile for
// unknown paths. Intended for hosting a client-side routed single-page app.
type SPA struct {
	fs        http.FileSystem
	indexFile string
	fileSrv   http.Handler
}

func NewSPA(fs http.FileSystem, indexFile string) *SPA {
	if indexFile == "" {
		indexFile = "index.html"
	}
	return &SPA{
		fs:        fs,
		indexFile: indexFile,
		fileSrv:   http.FileServer(fs),
	}
}

func (s *SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.NotFound(w, r)
		return
	}

	requested := path.Clean("/" + r.URL.Path)
	if f, err := s.fs.Open(requested); err == nil {
		fi, statErr := f.Stat()
		_ = f.Close()
		if statErr == nil && !fi.IsDir() {
			s.fileSrv.ServeHTTP(w, r)
			return
		}
	}

	indexPath := path.Clean("/" + s.indexFile)
	f, err := s.fs.Open(indexPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() { _ = f.Close() }()

	fi, err := f.Stat()
	if err != nil || fi.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, s.indexFile, fi.ModTime(), f)
}
