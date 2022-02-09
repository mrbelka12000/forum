package database

import (
	"database/sql"
	forum "forum/internal"
	"log"
	"net/http"
	"time"
)

//RemoveExpiredCookie ..
func RemoveExpiredCookie(db *sql.DB, create chan int) {
	for {
		_, err := db.Exec("DELETE FROM Cookie WHERE Expires <?", time.Now())
		if err != nil {
			create <- 0
		}
		time.Sleep(10 * time.Second)
	}
}

//IsUserInSession ..
func IsUserInSession(db *sql.DB, cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	var val string
	err := db.QueryRow("SELECT Value FROM Cookie WHERE Value = ?", cookie.Value).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false
	}
	if val == "" {
		return false
	}
	return true
}

//DeleteCookieFromDB ..
func DeleteCookieFromDB(db *sql.DB, cookie *http.Cookie) {
	_, err := db.Exec("DELETE FROM Cookie WHERE Value = ?", cookie.Value)
	if err != nil {
		log.Println(err)
		return
	}
}

//GetUserFromDB ..
func GetUserFromDB(db *sql.DB, cookie *http.Cookie) forum.User {
	var user forum.User

	row := db.QueryRow("SELECT SignInUser.Id, SignInUser.Name FROM SignInUser JOIN Cookie on UserId = id WHERE Value = ? ", cookie.Value)
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		log.Println(err)
		return user
	}

	if err := row.Err(); err != nil {
		return user
	}

	return user
}

//IsPreviousLikesDisliked ..
func IsPreviousLikesDisliked(db *sql.DB, yes, no bool, postID, userID int, firstReq, secondReq string) bool {
	var like, dislike bool
	row := db.QueryRow(firstReq, postID, userID)
	if err := row.Scan(&like, &dislike); err != nil {
		return true
	}
	_, err := db.Exec(secondReq, postID, userID)
	if err != nil {
		return false
	}
	if yes && like {
		return false
	} else if yes && !like {
		return true
	} else if no && dislike {
		return false
	} else if no && !dislike {
		return true
	}
	return true
}

/*

	RemoveExpiredCookie гоурутина которая каждые 30 секунд проверяет куки из БД
	если время жизни куки уже подошло к концу то удаляем ее из БД

	IsPreviousLikesDislikedPost

	DeleteCookieFromDB удаляем куки с БД если пользователь сам захотел выйти с сайта

	GetUserFromDB получаем айди пользователя исходя из его куки

	IsPreviousLikesDislikedPost проверяем последнее значение лайка и дизлайка
	пользователя из БД и делаем что то типа switch case

*/
