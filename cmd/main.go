package main

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"

	"sos/handlers"
	//"sos/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	store := sessions.NewCookieStore([]byte("your-secret-key"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // Время жизни сессии в секундах
		HttpOnly: true,
		Secure:   false, // Установите true для HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler(store)).Methods("POST")

	// r.HandleFunc("/admin", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/admin", handlers.AdminHandler(store)).Methods("GET")
	r.HandleFunc("/admin/users", handlers.AdminUsersHandler(store)).Methods("GET", "POST")
	r.HandleFunc("/admin/users/add", handlers.AdminAddUserHandler).Methods("POST")

	r.HandleFunc("/logout", handlers.LogoutHandler(store)).Methods("POST")
	r.HandleFunc("/main", handlers.MainHandler(store)).Methods("GET")
	r.HandleFunc("/report", handlers.ReportHandler).Methods("GET")
	r.HandleFunc("/getStops", handlers.GetStopsHandler).Methods("POST")
	r.HandleFunc("/getRouteNumbers", handlers.GetRouteNumbers).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Задайте дату и время
//dateStr := "2025-02-14 23:59:00"                    // Укажите вашу дату
//location, err := time.LoadLocation("Europe/Moscow") // Укажите нужный часовой пояс
//if err != nil {
//panic(err)
//}
//
//// Парсинг строки даты в объект time.Time
//dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, location)
//if err != nil {
//panic(err)
//}
//
//// Получение Unix timestamp
//unixTimestamp := dateTime.Unix()
//fmt.Println("Unix Timestamp:", unixTimestamp) // Вывод: Unix Timestamp: 1710055200
