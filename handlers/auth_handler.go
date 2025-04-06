package handlers

import (
	"log"
	"net/http"

	"sos/internal/model"

	"github.com/gorilla/sessions"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func LoginHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := model.GetUser(username)
		if err != nil || !user.CheckPassword(password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "auth-session")
		session.Values["authenticated"] = true
		session.Values["username"] = username

		var roleInt int
		switch user.Role {
		case "admin":
			roleInt = 0
		case "moderator":
			roleInt = 2
		case "user":
			roleInt = 3
		}
		session.Values["role"] = roleInt
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}
		log.Println("LoginHandler", user.Role)
		log.Println("LoginHandler", session)

		switch user.Role {
		case "user":
			log.Println("login to main")
			http.Redirect(w, r, "/main", http.StatusSeeOther)
		case "moderator":
			http.Redirect(w, r, "/report", http.StatusSeeOther)
		case "admin":
			log.Println("login to admin")
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		default:
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
}

func LogoutHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		session.Values["authenticated"] = false
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func InternalHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/main.html")
}
