package dashboard

import (
	"html/template"
	"log"
	"net/http"
)

// GET request Webpage serving
func Web(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("static/html/dashboard.html")
		err := t.Execute(w, nil)
		if err != nil {
			log.Println("[ERROR] Template Execute : ", err)
		}
	}
}
