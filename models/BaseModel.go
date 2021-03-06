package models

// BaseModel ...
type BaseModel struct {
	IsAuthorized bool
}

// PostListModel ...
type PostListModel struct {
	BaseModel
	EmptyMsg bool
	Exist    bool
	Posts    []Post
}

// EditPostModel ...
type EditPostModel struct {
	BaseModel
	Post Post
}
