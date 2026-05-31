package cdn

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const rootDir = "cdn-files"

var inlineTypes = map[string]bool{
	"text/plain":       true,
	"text/html":        true,
	"text/css":         true,
	"text/javascript":  true,
	"application/json": true,
	"image/png":        true,
	"image/jpeg":       true,
	"image/gif":        true,
	"image/webp":       true,
	"image/svg+xml":    true,
	"audio/mpeg":       true,
	"audio/ogg":        true,
	"video/mp4":        true,
	"video/webm":       true,
	"application/pdf":  true,
}

func Handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "application/json")

		ip := r.Header.Get("CF-Connecting-IP")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		if strings.Contains(ip, ",") {
			ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		}
		if ip == "" {
			ip = r.RemoteAddr
		}
		if strings.Contains(ip, ":") && !strings.Contains(ip, "[") {
			ip = strings.Split(ip, ":")[0]
		}

		userAgent := r.Header.Get("User-Agent")
		method := r.Method

		resp := map[string]string{
			"message":     "Welcome to the cast.onion CDN!",
			"version":     "0.1.0",
			"commit_hash": "78b790e48d2ba3bb73df7a6ecf3dd6dca1c07973",
			"ip":          ip,
			"country":     "US",
			"method":      method,
			"user_agent":  userAgent,
		}

		json.NewEncoder(w).Encode(resp)
		return
	}

	rel := strings.TrimPrefix(r.URL.Path, "/")
	if rel == "" || strings.Contains(rel, "..") {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	fullPath := filepath.Join(rootDir, filepath.FromSlash(rel))

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			serveCorruptedError(w, rel)
		}
		return
	}

	if info.IsDir() {
		serveDirectory(w, r, rel, fullPath)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		serveCorruptedError(w, rel)
		return
	}

	ext := strings.ToLower(filepath.Ext(fullPath))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	base := strings.Split(mimeType, ";")[0]
	if inlineTypes[base] {
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write(data)
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(fullPath)))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	}
}

func serveDirectory(w http.ResponseWriter, _ *http.Request, rel, fullPath string) {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		serveCorruptedError(w, rel)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta charset="utf-8">
<title>%s</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0a0a0a;color:#e0e0e0;font-family:'Courier New',monospace;font-size:14px;padding:40px}
h1{color:#c8ff00;font-size:16px;margin-bottom:24px}
a{color:#e0e0e0;text-decoration:none;display:block;padding:8px 0;border-bottom:1px solid #1a1a1a}
a:hover{color:#c8ff00}
.dir{color:#c8ff00}
.back{margin-bottom:16px}
</style></head><body>
<h1>/%s</h1>`, rel, rel)

	if rel != "" {
		parent := filepath.Dir(rel)
		if parent == "." {
			parent = ""
		}
		fmt.Fprintf(w, `<a class="back" href="/%s">../</a>`, parent)
	}

	for _, e := range entries {
		name := e.Name()
		path := rel + "/" + name
		if rel == "" {
			path = name
		}
		class := ""
		suffix := ""
		if e.IsDir() {
			class = ` class="dir"`
			suffix = "/"
		}
		fmt.Fprintf(w, `<a%s href="/%s">%s%s</a>`, class, path, name, suffix)
	}

	fmt.Fprintf(w, `</body></html>`)
}

func serveCorruptedError(w http.ResponseWriter, path string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta charset="utf-8">
<title>error — cast.onion cdn</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0a0a0a;color:#e0e0e0;font-family:'Courier New',monospace;font-size:14px;display:flex;align-items:center;justify-content:center;min-height:100vh}
.box{max-width:480px;text-align:center}
.title{color:#ff4444;font-size:18px;margin-bottom:12px}
.path{color:#555;font-size:12px;margin-bottom:24px}
.msg{color:#888;font-size:13px;line-height:1.7}
</style></head><body>
<div class="box">
<div class="title">we had a problem loading this file</div>
<div class="path">%s</div>
<div class="msg">the file or directory may be corrupted or unreadable.</div>
</div></body></html>`, path)
}
