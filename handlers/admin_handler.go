package handlers

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"

	repo "sos/internal/repo/clickhouse"

	"golang.org/x/crypto/bcrypt"
)

func AdminUsersHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("admin users")
		name := r.URL.Query().Get("name")
		login := r.URL.Query().Get("login")

		users, err := repo.GetUsers(name, login)
		if err != nil {
			http.Error(w, "failed to GetUsers", http.StatusInternalServerError)
			return
		}

		session, _ := store.Get(r, "auth-session")
		log.Println("session", session)
		role, ok := session.Values["role"].(int)
		if !ok {
			http.Error(w, "session role not found", http.StatusInternalServerError)
			//http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		log.Println("AdminUsersHandler", role)

		tmpl, err := template.ParseFiles("templates/admin.html")
		if err != nil {
			http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, map[string]interface{}{
			"Users": users,
			"Role":  role,
		})
		if err != nil {
			http.Error(w, "Execute error", http.StatusInternalServerError)
		}
	}
}

func AdminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("add")
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	err := repo.CreateUser(username, string(hashedPassword), role)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
