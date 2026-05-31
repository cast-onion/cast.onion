package web

import (
	"net/http"
	"os"
	"path/filepath"

	webadmin "github.com/cast-onion/internal/api/web/admin"
	"github.com/cast-onion/internal/db/queries"
)

type Handler struct {
	distDir string
}

func NewHandler(distDir string) *Handler {
	return &Handler{distDir: distDir}
}

func NewHandlerWithAPI(distDir, _ string, _ *queries.Queries) (*Handler, error) {
	return NewHandler(distDir), nil
}

func (h *Handler) Routes(adminH *webadmin.Handler) http.Handler {
	mux := http.NewServeMux()

	adminRoutes := adminH.Routes()
	mux.Handle("/admin", adminRoutes)
	mux.Handle("/admin/", adminRoutes)

	fs := http.FileServer(http.Dir(h.distDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(h.distDir, filepath.Clean(r.URL.Path))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(h.distDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	return mux
}
