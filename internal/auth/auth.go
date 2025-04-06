package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var sessionName = "auth-session"

func NewSessionStore() *sessions.CookieStore {
	// В продакшене используйте секреты из переменных окружения
	authKey := []byte("your-auth-key-32-byte")
	encryptionKey := []byte("your-encryption-key-16-byte")

	return sessions.NewCookieStore(authKey, encryptionKey)
}

func AuthMiddleware(store *sessions.CookieStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, sessionName)

			// Проверка аутентификации
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Redirect(w, r, "/main", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
