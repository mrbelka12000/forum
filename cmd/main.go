package main

import (
	"forum/internal/database"
	"forum/internal/server"
	"log"
	"net/http"
)

func main() {
	db, err := database.DbInit()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	create := make(chan int)

	go database.RemoveExpiredCookie(db, create)

	go database.CreateDb(db, create)

	handle := server.CreateHandle(db)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("templates/assets/"))))
	http.HandleFunc("/", handle.HomePage)

	http.HandleFunc("/login", handle.GetSignIn)
	http.HandleFunc("/postlogin", handle.PostSignIn)

	http.HandleFunc("/signout", handle.GetSignOut)
	http.HandleFunc("/postsignout", handle.PostSignOut)

	http.HandleFunc("/registration", handle.GetSignUp)
	http.HandleFunc("/postregistration", handle.PostSignUp)

	http.HandleFunc("/post/", handle.ChosenPost)
	http.HandleFunc("/createpost", handle.GetCreatePost)
	http.HandleFunc("/postcreate", handle.PostCreatePost)
	http.HandleFunc("/ratepost", handle.RatePost)
	http.HandleFunc("/ratecomment", handle.RateComment)
	http.HandleFunc("/placecomment", handle.CommentsHandler)

	http.HandleFunc("/profile", handle.Profile)
	http.HandleFunc("/filter", handle.Filter)

	log.Println("Server is Listening..." + "\n" + "http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
