package exec

import (
	"database/sql"
	"fmt"
	customstructs "forum/customstructs"
)

var (
	Db *sql.DB
)

// Retrieves data from the "users" table based on the condition string.
func GetUserFromDb(condition string, value interface{}) ([]customstructs.UserData, error) {
	var (
		rows       *sql.Rows
		err        error
		returnData []customstructs.UserData
	)
	switch condition {
	case "":
		rows, err = Db.Query("SELECT user_id, email, first_name, last_name, date_of_birth, online_status, role_id FROM users")
	default:
		rows, err = Db.Query("SELECT user_id, email, first_name, last_name, date_of_birth, online_status, role_id FROM users WHERE "+condition+" = $1", value)
	}

	if err != nil {
		fmt.Println("getData.go GetUserFromDb err 1: ", err)
		return nil, err
	}

	for rows.Next() {
		var t customstructs.UserData
		err = rows.Scan(&t.Id, &t.Email, &t.FirstName, &t.LastName, &t.DateOfBirth, &t.OnlineStatus, &t.RoleId)
		if err != nil {
			fmt.Println("getData.go GetUserFromDb err 2: ", err)
			return nil, err
		}
		returnData = append(returnData, t)
	}
	rows.Close()

	if len(returnData) == 0 {
		err := fmt.Errorf("getData.go GetUserFromDb: returndata is blank")
		return nil, err
	}

	return returnData, nil
}

// Retrieves data from the "posts" table based on the condition string.
func GetPostsFromDb(condition string, value interface{}) ([]customstructs.PostsData, error) {
	var (
		rows       *sql.Rows
		err        error
		returnData []customstructs.PostsData
	)

	switch condition {
	case "":
		rows, err = Db.Query("SELECT * FROM posts")
	default:
		rows, err = Db.Query("SELECT * FROM posts WHERE "+condition+" = $1", value)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t customstructs.PostsData
		err = rows.Scan(&t.Id, &t.Creator.Id, &t.PostTime, &t.Title, &t.Body, &t.CategoryId, &t.GroupId)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, t)
	}

	return returnData, nil
}

// Retrieves data from the "comments" table based on the condition string.
func GetCommentsFromDb(condition string, value interface{}) ([]customstructs.CommentData, error) {
	var (
		rows       *sql.Rows
		err        error
		returnData []customstructs.CommentData
	)

	switch condition {
	case "":
		rows, err = Db.Query("SELECT * FROM comments")
	default:
		rows, err = Db.Query("SELECT comment_id, post_id, user_id, comment_datetime, body FROM comments WHERE "+condition+" = $1", value)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t customstructs.CommentData
		err = rows.Scan(&t.Id, &t.PostId, &t.Creator.Id, &t.CommentTime, &t.Body)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, t)
	}
	return returnData, nil
}

// Retrieves data from the "sessions" table based on the session string. Used for user authentication.
func GetUserBySessionID(session string) customstructs.Session {
	var sessionInfo customstructs.Session
	rows, err := Db.Query("SELECT * FROM sessions WHERE session_id = $1", session)
	if err != nil {
		fmt.Println("GetUserBySessionID err:", err)
		return sessionInfo
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&sessionInfo.Session, &sessionInfo.UserID, &sessionInfo.Username, &sessionInfo.Datetime, &sessionInfo.RoleID)
	}

	return sessionInfo
}

// Retrieves the userId from the "sessions" table based on the session cookie value.
func GetUserByCookie(session string) int {
	sessionData := GetUserBySessionID(session)
	return sessionData.UserID
}
