package server

import (
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//ChosenPost ..
func (handle *Handle) ChosenPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	t, err := template.ParseFiles("./templates/html/chosenpost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	id := r.URL.Path[6:]
	passed, postid := database.IsPostInDB(handle.DB, id)

	if !passed {
		w.WriteHeader(404)
		return
	}
	post, err := database.GetOnePost(handle.DB, postid)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	if handle.Post.ErrorVal.Err != false {
		post.ErrorVal = handle.Post.ErrorVal
	}
	if err := t.Execute(w, post); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	handle.Post.ErrorVal.Err = false
}

//RatePost ..
func (handle *Handle) RatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	if r.URL.Path != "/ratepost" {
		w.WriteHeader(404)
		return
	}

	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	prevURL := r.Header.Get("Referer")
	id := prevURL[27:]
	postid, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	statement := "INSERT INTO PostRaiting (Like, DisLike, PostId, LikerId) VALUES (?,?,?,?)"
	checkSel := "SELECT Like , DisLike FROM PostRaiting WHERE PostId = ? AND LikerId = ?"
	checkDel := "DELETE FROM PostRaiting WHERE PostId = ? AND LikerId = ?"
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

	database.InsertLikeDislikeIntoDB(handle.DB, like, dislike, postid, cookie, statement, checkSel, checkDel)
	http.Redirect(w, r, prevURL, 301)
}
