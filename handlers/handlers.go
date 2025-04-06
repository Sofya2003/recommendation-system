package handlers

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sos/internal/model"
	//"../internal"
)

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	//if r.Method == "POST" {
// 	//	// Вызов функции из internal для авторизации
// 	//	internal.HandleLogin(w, r)
// 	//	return
// 	//}
// 	tmpl := template.Must(template.ParseFiles("../templates/login.html"))
// 	tmpl.Execute(w, nil)
// }

//func MainHandler(w http.ResponseWriter, r *http.Request) {
//	tmpl := template.Must(template.ParseFiles("templates/main.html", "header.html"))
//	err := tmpl.Execute(w, nil)
//	if err != nil {
//		log.Println(err)
//	}
//}

func MainHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			"templates/main.html",
			"templates/header.html",
		)
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			log.Println("Error parsing templates:", err)
			return
		}

		session, _ := store.Get(r, "auth-session")
		log.Println("session", session)
		role, ok := session.Values["role"].(int)
		if !ok {
			log.Println("Role not found")
		}
		log.Println("MainHandler", role)

		err = tmpl.Execute(w, map[string]interface{}{
			"Role": role,
		})
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Execute error", http.StatusInternalServerError)
		}
		//tmpl := template.Must(template.ParseFiles("templates/main.html", "header.html"))
		//err := tmpl.Execute(w, nil)
		//if err != nil {
		//	log.Println(err)
		//}
	}
}

//func AdminHandler(w http.ResponseWriter, r *http.Request) {
//	tmpl := template.Must(template.ParseFiles("../templates/admin.html"))
//	tmpl.Execute(w, nil)
//}

func AdminHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		role, ok := session.Values["role"].(int)

		if !ok || role != 0 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		//tmpl := template.Must(template.ParseFiles("templates/admin.html"))
		//tmpl.Execute(w, nil)
	}
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/report.html"))
	tmpl.Execute(w, nil)
}

func GetStopsHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RouteNumber string `json:"routeNumber"`
		TimePeriod  string `json:"timePeriod"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Логика для получения данных о остановках
	stops := []model.Stop{
		{Workload: float64(rand.Intn(100))},
		{Workload: float64(rand.Intn(100))},
		{Workload: float64(rand.Intn(100))},
		{Workload: float64(rand.Intn(100))},
		{Workload: float64(rand.Intn(100))},
	}
	log.Println(stops)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]model.Stop{"stops": stops})
}

func GetRouteNumbers(w http.ResponseWriter, r *http.Request) {
	//var request struct {
	//	RouteNumber string `json:"routeNumber"`
	//	TimePeriod  string `json:"timePeriod"`
	//}
	//if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	//	http.Error(w, "Invalid request", http.StatusBadRequest)
	//	return
	//}

	// Логика для получения данных о остановках
	routeNumbers := []int{1, 2, 3}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]int{"routeNumbers": routeNumbers})
}
