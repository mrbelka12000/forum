package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//Profile ..
func (handle *Handle) Profile(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./templates/html/profile.html")
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		user := database.GetUserFromDB(handle.DB, cookie)
		createdPosts, err := database.GetPostByUser(handle.DB, user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		likedPosts, err := database.GetLikedPostsByUser(handle.DB, user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		ditails := forum.UserProfile{
			User:         user,
			CreatedPosts: createdPosts,
			LikedPosts:   likedPosts,
		}
		if err := t.Execute(w, ditails); err != nil {
			w.WriteHeader(500)
			return
		}
	default:
		w.WriteHeader(405)
		return
	}
}
