package middleware

import (
	"log"
	"net/http"
	"website/internal/config"
)

var store = config.GetSessionStore()

func IsAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Erreur lors de la récupération de la session: %v", err)
		return false
	}
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func GetUserInfo(r *http.Request) (string, bool) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Erreur lors de la récupération de la session: %v", err)
		return "", false
	}
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
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Erreur lors de la récupération de la session: %v", err)
		return err
	}
	session.Values["authenticated"] = true
	session.Values["username"] = username
	if err := session.Save(r, w); err != nil {
		log.Printf("Erreur lors de la sauvegarde de la session: %v", err)
		return err
	}
	return nil
}

func ClearAuthSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Erreur lors de la récupération de la session: %v", err)
		return err
	}
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	if err := session.Save(r, w); err != nil {
		log.Printf("Erreur lors de la sauvegarde de la session: %v", err)
		return err
	}
	return nil
}
