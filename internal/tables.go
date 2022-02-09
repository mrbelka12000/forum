package forum

// SingInUsers ..
type SingInUsers struct {
	ID       int
	Name     string
	Login    string
	Password string
}

// Post ..
type Post struct {
	ID              int
	Title           string
	Text            string
	Author          string
	Tags            []string
	Comments        []Comment
	CountOfLikes    int
	CountOfDisLikes int
	ErrorVal        Error
}

// Comment ..
type Comment struct {
	ID              int
	Text            string
	Author          string
	CountOfLikes    int
	CountOfDisLikes int
}

//Error ..
type Error struct {
	MSG string
	Err bool
}

//Homepage ..
type Homepage struct {
	Posts      []Post
	Categories []string
	Category   string
	InSession  bool
}

//UserProfile ..
type UserProfile struct {
	User         User
	CreatedPosts []Post
	LikedPosts   []Post
}

//UnionEnterPost ..
type UnionEnterPost struct {
	ErrorVal   Error
	Categories []string
}

type UnionEnterComment struct {
	ErrorVal Error
	PostId   int
	Comment  string
}

//CreatePost ..
type CreatePost struct {
	Title      string
	Text       string
	Categories []string
}

//Categories
type Categories struct {
	Tags      []string
	InSession bool
}
