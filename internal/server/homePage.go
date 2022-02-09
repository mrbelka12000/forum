package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//HomePage ..
func (handle *Handle) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return 
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	t, err := template.ParseFiles("./templates/html/homePage.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	posts, err := database.GetPosts(handle.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	HomeInfo := forum.Homepage{
		Posts:      posts,
		Categories: categories,
	}
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		if err := t.Execute(w, HomeInfo); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		return
	}

	HomeInfo.InSession = true
	if err := t.Execute(w, HomeInfo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}
