package customstructs

import (
	"database/sql"
	"time"
)

type UserData struct {
	Id           int
	Email        string
	Password     string
	FirstName    string
	LastName     string
	DateOfBirth  string
	OnlineStatus string
	RoleId       int
	Followers    []UserData
	Following    []UserData

	ImageUrl sql.NullString
	Nickname sql.NullString
	AboutMe  sql.NullString

	Visibility string // "public" - everyone can see details, "private" - only followers can see details
	Posts      []PostsData
	Comments   []CommentData
}

type CommentData struct {
	Id          int
	PostId      int
	CommentTime string
	Body        string
	// Likes       int
	// Dislikes    int
	ImageUrl string
	Creator  UserData
}

type Session struct {
	Session  string
	UserID   int
	Username string // in db as "email"
	RoleID   int
	Datetime string
}

type PostsData struct {
	Id         int
	PostTime   string
	Title      string
	Body       string
	Comments   []CommentData
	CategoryId int
	//CatTitle   string
	//Likes      int
	//Dislikes   int
	GroupId  int
	ImageUrl string
	Creator  UserData
}

type GroupData struct {
	GroupId    int
	Creator    []UserData
	GroupTitle string
	GroupDesc  string
	Members    []UserData
	Posts      []PostsData
}

type GroupEvent struct {
	EventID        string     `json:"event_id"`
	Group          any        `json:"group"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	DayTime        time.Time  `json:"day_time"`
	AttendingUsers []UserData `json:"attending_users"`
}

type EventResponse struct {
	ResponseId string `json:"response_id"`
	EventID    string `json:"event_id"`
	User       any    `json:"user"`
	Response   string `json:"response"` // "Attending" or "Not attending" or "Pending"
}
