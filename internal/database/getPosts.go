package database

import (
	"database/sql"
	forum "forum/internal"
	"log"
	"strconv"
)

//GetPostsByCategories ..
func GetPostsByCategories(db *sql.DB, tag string) ([]forum.Post, error) {
	var posts []forum.Post
	rows, err := db.Query(`SELECT Post.id, Post.Title , Post.Post , Post.AuthorId
	 FROM Categories JOIN Post ON PostId = Post.id AND Categories.Name =? ORDER BY Post.id DESC`, tag)
	if err != nil {
		log.Println(err)
		return posts, err
	}

	defer rows.Close()

	for rows.Next() {
		var postid, authorid int
		var title, text, authorname string
		rows.Scan(&postid, &title, &text, &authorid)
		tags := GetOnePostCategories(db, postid)
		if err := db.QueryRow("SELECT Name FROM  SignInUser WHERE id = ?", authorid).Scan(&authorname); err != nil {
			log.Println(err)
			return posts, err
		}
		post := forum.Post{
			ID:     postid,
			Title:  title,
			Text:   text,
			Author: authorname,
			Tags:   tags,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

//GetPostByUser ..
func GetPostByUser(db *sql.DB, user forum.User) ([]forum.Post, error) {
	var postArr []forum.Post
	rows, err := db.Query("SELECT id, Title FROM POST WHERE AuthorId = ?", user.ID)
	if err != nil {
		log.Println(err)
		return postArr, err
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			log.Println(err)
			return postArr, err
		}
		curPost := forum.Post{
			ID:    id,
			Title: title,
		}
		postArr = append(postArr, curPost)
	}
	return postArr, nil
}

//GetCommentsFromPost ..
func GetCommentsFromPost(db *sql.DB, postID int) ([]forum.Comment, error) {
	var comments []forum.Comment
	rows, err := db.Query(`SELECT Comments.id, Comment, SignInUser.Name FROM Comments
	JOIN SignInUser ON CommentatorId = SignInUser.Id 
	WHERE PostId= ? ORDER BY Comments.id DESC`, postID)
	if err != nil {
		log.Println(err)
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment, userName string
		var id, like, dislike int
		if err := rows.Scan(&id, &comment, &userName); err != nil {
			return comments, err
		}
		db.QueryRow(`SELECT SUM(Like),SUM(DisLike) FROM CommentRaiting WHERE CommentsId =? `, id).Scan(&like, &dislike) //err?

		Comment := forum.Comment{
			ID:              id,
			Text:            comment,
			Author:          userName,
			CountOfLikes:    like,
			CountOfDisLikes: dislike,
		}
		comments = append(comments, Comment)
	}
	return comments, err
}

//GetLikedPostsByUser ..
func GetLikedPostsByUser(db *sql.DB, user forum.User) ([]forum.Post, error) {
	var postArr []forum.Post
	rows, err := db.Query(`SELECT Post.id, Post.Title FROM PostRaiting 
	JOIN Post ON PostId=Post.id WHERE LikerId=? AND Like=true 
	ORDER BY Post.id DESC`, user.ID)
	if err != nil {
		log.Println(err)
		return postArr, err
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			log.Println(err)
			return postArr, err
		}
		curPost := forum.Post{
			ID:    id,
			Title: title,
		}
		postArr = append(postArr, curPost)
	}
	return postArr, nil
}

//GetOnePostCategories ..
func GetOnePostCategories(db *sql.DB, postid int) []string {
	var tags []string
	row, err := db.Query(`
	SELECT Name FROM Categories
	 WHERE PostId=?;`, postid)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer row.Close()
	for row.Next() {
		var catname string
		if err := row.Scan(&catname); err != nil {
			log.Println(err)
			return nil
		}
		tags = append(tags, catname)
	}
	if row.Err() != nil {
		return nil
	}
	return tags
}

//GetPosts ..
func GetPosts(db *sql.DB) ([]forum.Post, error) {
	var posts []forum.Post
	rows, err := db.Query(`
	SELECT  Post.id, Post.Title, Post.Post, SignInUser.Name 
	FROM Post  JOIN SignInUser ON AuthorId = SignInUser.id ORDER BY Post.id DESC `)
	if err != nil {
		log.Println(err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var postid int
		var title, text, author string
		if err := rows.Scan(&postid, &title, &text, &author); err != nil {
			log.Println(err)
			return posts, err
		}
		tags := GetOnePostCategories(db, postid)

		post := forum.Post{
			ID:     postid,
			Title:  title,
			Text:   text,
			Author: author,
			Tags:   tags,
		}
		posts = append(posts, post)

	}

	if rows.Err() != nil {
		return posts, err
	}
	return posts, nil
}

//GetOnePost ..
func GetOnePost(db *sql.DB, postid int) (forum.Post, error) {
	var post forum.Post
	row := db.QueryRow(`SELECT  Post.Title, Post.Post, SignInUser.Name FROM Post
	JOIN SignInUser ON AuthorId = SignInUser.id WHERE Post.id =?`, postid)
	var title, text, author string
	if err := row.Scan(&title, &text, &author); err != nil {
		return post, err
	}
	rows, err := db.Query(`
		SELECT Name FROM Categories
		 WHERE PostId=?;`, postid)
	if err != nil {
		return post, err
	}
	var tags []string
	for rows.Next() {
		var categname string
		if err := rows.Scan(&categname); err != nil {
			return post, err
		}
		tags = append(tags, categname)
	}
	result, err := GetCommentsFromPost(db, postid)
	if err != nil {
		log.Println(err)
		return post, err
	}
	countlikes, countdislikes, err := GetPostLikes(db, postid)
	if err != nil {
		log.Println(err)
		return post, err
	}
	post = forum.Post{
		ID:              postid,
		Title:           title,
		Text:            text,
		Author:          author,
		Tags:            tags,
		Comments:        result,
		CountOfLikes:    countlikes,
		CountOfDisLikes: countdislikes,
	}
	return post, nil
}

//IsPostInDB ..
func IsPostInDB(db *sql.DB, id string) (bool, int) {
	postid, err := strconv.Atoi(id)
	if err != nil {
		return false, postid
	}
	var val int
	err = db.QueryRow("SELECT id FROM Post WHERE id = ?", postid).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false, postid
	}
	if val == 0 {
		return false, postid
	}
	return true, postid
}

//GetPostLikes ..
func GetPostLikes(db *sql.DB, postID int) (int, int, error) {
	rows, err := db.Query("SELECT Like, DisLike FROM PostRaiting WHERE PostId = ?", postID)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	var like, dislike int
	for rows.Next() {
		var pos, neg int
		if err := rows.Scan(&pos, &neg); err != nil {
			return 0, 0, err
		}
		like += pos
		dislike += neg
	}
	return like, dislike, nil
}
