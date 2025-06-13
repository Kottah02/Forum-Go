package middleware

import (
	"net/http"
	"website/internal/config"
)

var store = config.GetSessionStore()

func IsAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func GetUserInfo(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	return username, ok
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func SetAuthSession(w http.ResponseWriter, r *http.Request, username string) error {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	return session.Save(r, w)
}

func ClearAuthSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	return session.Save(r, w)
}
