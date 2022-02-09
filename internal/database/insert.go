package database

import (
	"database/sql"
	"errors"
	"fmt"
	forum "forum/internal"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//InsertNewUserIntoDB ..
func InsertNewUserIntoDB(db *sql.DB, user forum.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("When generating hash from password InsertNewUserIntoDB() %w : ", err)
	}
	_, err = db.Exec("INSERT INTO SignInUser (Name, Login , Password) VALUES (?,?,?)", user.Name, user.Login, string(hash))
	if err != nil {
		return fmt.Errorf("Sorry name or login is already taken")
	}
	return nil
}

//InsertCookieIntoDB ..
func InsertCookieIntoDB(db *sql.DB, login string, cookie *http.Cookie) {
	var id int
	row := db.QueryRow("SELECT Id FROM SignInUser WHERE Login =?;", login)
	row.Scan(&id)

	_, err := db.Exec("DELETE FROM Cookie WHERE UserId =?", id)
	if err != nil {
		log.Println()
		return
	}

	_, err = db.Exec("INSERT INTO Cookie (Value, Expires,UserId)VALUES(?,?,?)", cookie.Value, cookie.Expires, id)
	if err != nil {
		log.Println(err)
		return
	}
}

//InsertPostIntoDB ..
func InsertPostIntoDB(db *sql.DB, post *forum.CreatePost, cookie *http.Cookie) error {
	user := GetUserFromDB(db, cookie)
	userID := user.ID
	title := strings.Trim(post.Title, "\r\n")
	title = strings.Trim(title, " ")

	text := strings.Trim(post.Text, "\r\n")
	text = strings.Trim(text, " ")
	if title == "" {
		return errors.New("Title must not be empty")
	}
	if text == "" {
		return errors.New("Text must not be empty")
	}

	if len(title) > 30 {
		return errors.New("Length of title must be less than 30 characters")
	}

	if len(post.Categories) == 0 {
		return errors.New("Please choose one category")
	}
	if title == "" {
		return errors.New("Invalid title")

	} else if text == "" {
		return errors.New("Invalid description")
	}
	res, err := db.Exec("INSERT INTO Post (Title , Post, AuthorID) VALUES(?,?,?)", post.Title, post.Text, userID)
	if err != nil {
		return err
	}

	postid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if err := InsertCategoriesToDB(db, post.Categories, postid); err != nil {
		_, err := db.Exec("DELETE FROM Post WHERE Id =?; ", postid)
		if err != nil {
			return err
		}
	}
	return nil
}

//InsertCategoriesToDB ..
func InsertCategoriesToDB(db *sql.DB, categories []string, postid int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	for _, category := range categories {
		_, err := tx.Exec("INSERT INTO Categories (Name,PostId) VALUES (?,?)", category, postid)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return nil
}

//InsertCommentToDB ..
func InsertCommentToDB(db *sql.DB, comment string, postID int, cookie *http.Cookie) error {
	test := strings.Replace(comment, "\r\n", " ", -1)
	test = strings.Trim(test, " ")
	if test == "" {
		return errors.New("Invalid comment")
	}
	user := GetUserFromDB(db, cookie)
	userID := user.ID
	_, err := db.Exec("INSERT INTO Comments (Comment, CommentatorId, PostId) VALUES (?,?,?)", comment, userID, postID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//InsertLikeDislikeIntoDB ..
func InsertLikeDislikeIntoDB(db *sql.DB, like, dislike string, id int, cookie *http.Cookie, req, freq, sreq string) {
	var yes bool
	var no bool
	if len(like) != 0 {
		yes = true
	} else if len(dislike) != 0 {
		no = true
	}
	if !yes && !no {
		return
	}
	user := GetUserFromDB(db, cookie)
	userID := user.ID
	do := IsPreviousLikesDisliked(db, yes, no, id, userID, freq, sreq)
	if do {
		_, err := db.Exec(req, yes, no, id, userID)
		if err != nil {
			log.Println(err, "--------")
			return
		}
		return
	}
	return
}

/*

	InsertNewUserIntoDB после того как проверили на уникальность
	имя , логин и пароль добавляем пользователя в БД

	InsertCookieIntoDB получаем айди пользователя по логину и затем
	добавляем в БД куки со значением(сами куки) время истечения и номер пользователя

	InsertPostIntoDB получаем айди пользователя из куки и присваеваем посту
	название , содержимое и номер пользователя который оставил пост

	InsertCategoriesToDB добавляем категории в БД по названию категории и
	номеру поста под которым находится категория

	InsertCommentToDB получаем ID поста под которым надо оставить коммент
	и добавляем в БД само значем значение коммента и кто оставил

	InsertLikeDislikeIntoDB смотрим значение кнопок лайка и дизлайка
	в зависимости от клика добавлям в БД по значение кнопки, но перед этим
	проверить предыдущие лайки пользователя под этим постом
*/
