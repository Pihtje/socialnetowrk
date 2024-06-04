package exec

import (
	"fmt"
	"time"
)

// Inserts comment data into the database.
func Comment(postId, userId int, body string) ([]int, error) {
	statement, err := Db.Prepare("INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (?, ?, ?, ?);")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(postId, userId, time.Now(), body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Comment added successfully to db!")

	returnData, err := res.LastInsertId()
	if err != nil {
		fmt.Println("comment.go Comment err: ", err)
		return nil, err
	}

	return []int{int(returnData)}, nil
}
