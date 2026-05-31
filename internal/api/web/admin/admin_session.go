package admin

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

var sessionSecret []byte

func init() {
	sessionSecret = make([]byte, 32)
	rand.Read(sessionSecret)
}

func SetSessionCookie(w http.ResponseWriter, username string) {
	value := username + ":" + signValue(username)
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    value,
		Path:     "/admin",
		HttpOnly: true,
		Expires:  time.Now().Add(8 * time.Hour),
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/admin",
		MaxAge: -1,
	})
}

func ValidateSession(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return "", false
	}
	for i := len(cookie.Value) - 1; i >= 0; i-- {
		if cookie.Value[i] == ':' {
			username := cookie.Value[:i]
			sig := cookie.Value[i+1:]
			if sig == signValue(username) {
				return username, true
			}
			return "", false
		}
	}
	return "", false
}

func signValue(value string) string {
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

func RequireAdminSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := ValidateSession(r); !ok {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
