package server

import (
	"forum/internal/database"
	"net/http"
	"strconv"
)

//CommentsHandler ..
func (handle *Handle) CommentsHandler(w http.ResponseWriter, r *http.Request) {
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
	comment := r.FormValue("comment")
	comment_btn := r.FormValue("comment_btn")
	for v := range r.Form {
		if v != "comment" && v != "comment_btn" {
			w.WriteHeader(400)
			return
		}
	}
	if comment_btn != "sent" {
		w.WriteHeader(400)
		return
	}
	if err := database.InsertCommentToDB(handle.DB, comment, postid, cookie); err != nil {
		handle.Post.ErrorVal.Err = true
		handle.Post.ErrorVal.MSG = err.Error()
	}
	http.Redirect(w, r, prevURL, 301)
}

//RateComment ..
func (handle *Handle) RateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	prevURL := r.Header.Get("Referer")

	var commentID int
	like := r.FormValue("likecom")
	dislike := r.FormValue("dislikecom")
	statement := "INSERT INTO CommentRaiting (Like, DisLike, CommentsId, UserId) VALUES (?,?,?,?)"
	checkSel := "SELECT Like , DisLike FROM CommentRaiting WHERE CommentsId = ? AND UserId = ?"
	checkDel := "DELETE FROM CommentRaiting WHERE CommentsId = ? AND UserId = ?"
	if len(like) != 0 {
		commentID, _ = strconv.Atoi(like)
	} else if len(dislike) != 0 {
		commentID, _ = strconv.Atoi(dislike)
	}
	if commentID != 0 {
		database.InsertLikeDislikeIntoDB(handle.DB, like, dislike, commentID, cookie, statement, checkSel, checkDel)
	}
	http.Redirect(w, r, prevURL, 301)
}
