package exec

import (
	"database/sql"
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"sort"
)

// Handler for "/allPosts/" endpoint.
// Retrieves all posts that the user is allowed to see from the database.
func ServeFilteredPosts(w http.ResponseWriter, r *http.Request) {
	var returnData []customstructs.PostsData
	var allPosts []customstructs.PostsData

	allPosts, err := GetPostsFromDb("group_id", "0")
	if err != nil {
		msg, _ := json.Marshal(err.Error())
		w.Write(msg)
		return
	}

	publicPosts, err := GetPublicPosts(allPosts)
	if err != nil {
		fmt.Println("getData2.go ServeFilteredPosts err 1: ", err.Error())
	} else {
		returnData = append(returnData, publicPosts...)
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("getData2.go ServeFilteredPosts err 2: ", err)
		json.NewEncoder(w).Encode("logged in only")
		return
	}

	userId := GetUserByCookie(cookie.Value)
	if userId == 0 {
		fmt.Println("getData2.go ServeFilteredPosts err 6: ", err)
		json.NewEncoder(w).Encode("incorrect cookie value")
		return
	}

	followingPosts, err := GetFollowingPosts(userId, allPosts)
	if err != nil {
		fmt.Println("getData2.go ServeFilteredPosts err 3: ", err.Error())
	} else {
		returnData = append(returnData, followingPosts...)
	}

	sharedPosts, err := GetSharedPosts(userId)
	if err != nil {
		fmt.Println("getData2.go ServeFilteredPosts err 4: ", err.Error())
	} else {
		returnData = append(returnData, sharedPosts...)
	}

	ownPosts, err := GetOwnPosts(userId, allPosts)
	if err != nil {
		fmt.Println("getData2.go ServeFilteredPosts err 5: ", err.Error())
	} else {
		returnData = append(returnData, ownPosts...)
	}

	for i, v := range returnData {
		p, _ := GetPost(fmt.Sprint(v.Id), userId)
		returnData[i] = p[0]
	}

	returnData = SortData(returnData)

	msg, _ := json.Marshal(returnData)
	w.Write(msg)
}

// Retrieves all followers of the user from the database. Returns an empty array if the user doesn't have any followers.
func GetFollowersFromDb(ud []customstructs.UserData) ([]customstructs.UserData, error) {
	var (
		rows *sql.Rows
		err  error
	)

	followersIntArr := []int{}

	if len(ud) == 0 {
		err = fmt.Errorf("a user with that id does not exist in the database")
		return nil, err
	}

	rows, err = Db.Query("SELECT follower_id FROM followers WHERE user_id = $1", ud[0].Id)
	if err != nil {
		fmt.Println("GetFollowersFromDb err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			fmt.Println("GetFollowersFromDb scan err 1: ", err)
			return nil, err
		}
		followersIntArr = append(followersIntArr, i)
	}

	followers := []customstructs.UserData{}

	if len(followersIntArr) != 0 {
		for _, v := range followersIntArr {
			f, err := GetUser("users.user_id", fmt.Sprint(v))
			if err != nil {
				fmt.Println("getData2.go GetFollowersFromDb err 2: ", err)
				return nil, err
			}
			followers = append(followers, f[0])
		}
	}

	ud[0].Followers = followers

	return ud, nil
}

// Retrieves all the users that the currently logged in user is a follower of. Returns an empty array if the user isn't following any other users.
func GetFollowingFromDb(ud []customstructs.UserData) ([]customstructs.UserData, error) {
	var (
		rows *sql.Rows
		err  error
	)

	followingIntArr := []int{}

	if len(ud) == 0 {
		err = fmt.Errorf("a user with that id does not exist in the database")
		return nil, err
	}

	rows, err = Db.Query("SELECT user_id FROM followers WHERE follower_id = $1", ud[0].Id)
	if err != nil {
		fmt.Println("GetFollowingFromDb err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			fmt.Println("GetFollowingFromDb scan err 1: ", err)
			return nil, err
		}
		followingIntArr = append(followingIntArr, i)
	}

	followingUsers := []customstructs.UserData{}

	if len(followingIntArr) != 0 {
		for _, v := range followingIntArr {
			fu, err := GetUser("users.user_id", fmt.Sprint(v))
			if err != nil {
				fmt.Println("getData2.go GetFollowingFromDb err 2: ", err)
				return nil, err
			}
			followingUsers = append(followingUsers, fu[0])
		}
	}

	ud[0].Following = followingUsers

	return ud, nil
}

// Retrieves all publicly accessible posts from the database.
func GetPublicPosts(data []customstructs.PostsData) ([]customstructs.PostsData, error) {
	if len(data) == 0 {
		err := fmt.Errorf("getData2.go GetPublicPosts: allPosts is empty")
		return nil, err
	}

	returnData := []customstructs.PostsData{}

	for _, v := range data {
		if v.CategoryId == 1 {
			returnData = append(returnData, v)
		}
	}

	return returnData, nil
}

// Retrieves all posts made by users that the current user is following.
func GetFollowingPosts(userId int, data []customstructs.PostsData) ([]customstructs.PostsData, error) {
	var returnData []customstructs.PostsData

	fmt.Println("GetFollowingPosts userId:", userId)

	userInfo, err := GetUser("users.user_id", fmt.Sprint(userId))
	if err != nil {
		fmt.Println("getData2.go GetFollowingPosts err 1: ", err)
		return nil, err
	}

	userInfo = ProcessUserInfo(userInfo, userId)

	for _, v := range userInfo[0].Following {
		var loopData []customstructs.PostsData
		for _, a := range data {
			if a.CategoryId == 2 && a.Creator.Id == v.Id {
				loopData = append(loopData, a)
			}
		}

		returnData = append(returnData, loopData...)
	}

	return returnData, nil
}

// Sorts the data array based on the post Id.
func SortData(data []customstructs.PostsData) []customstructs.PostsData {
	fmt.Println("Sorting data...")
	sort.Slice(data, func(i, j int) bool {
		a, b := data[i], data[j]
		return a.Id < b.Id
	})
	return data
}

// Retrieves all posts that have been shared with the user.
func GetSharedPosts(userId int) ([]customstructs.PostsData, error) {
	Rows, err := Db.Query("SELECT posts.* from posts, shared_posts WHERE posts.post_id=shared_posts.post_id AND shared_posts.user_id=?;", userId)
	if err != nil {
		fmt.Println("getData2.go GetSharedPosts scan err 1: ", err)
		return nil, err
	}
	defer Rows.Close()

	var returnData []customstructs.PostsData

	for Rows.Next() {
		var tmpData customstructs.PostsData
		err := Rows.Scan(&tmpData.Id, &tmpData.Creator.Id, &tmpData.PostTime, &tmpData.Title, &tmpData.Body, &tmpData.CategoryId, &tmpData.GroupId)
		if err != nil {
			fmt.Println("GetData2.go GetSharedPosts: ", err)
		}
		returnData = append(returnData, tmpData)
	}

	return returnData, nil
}

// Inserts a user into the followers list of another user.
func AddFollower(userId, followerId string) error {
	err := RemoveFollower(userId, followerId)
	if err != nil {
		fmt.Println("GetData2.go AddFollower err 1: ", err)
	}

	statement, err := Db.Prepare("INSERT INTO followers (user_id, follower_id) VALUES (?, ?);")
	if err != nil {
		fmt.Println("GetData2.go AddFollower err 2: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userId, followerId)
	if err != nil {
		fmt.Println("GetData2.go AddFollower err 3: ", err)
		return err
	}
	fmt.Println("Follower added successfully!")
	return nil
}

// Removes a user from the followers list of another user.
func RemoveFollower(userId, followerId string) error {
	statement, err := Db.Prepare("DELETE FROM followers WHERE (user_id, follower_id) = (?, ?);")
	if err != nil {
		fmt.Println("GetData2.go RemoveFollower err 1: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userId, followerId)
	if err != nil {
		fmt.Println("GetData2.go RemoveFollower err 2: ", err)
		return err
	}
	fmt.Println("Follower removed successfully!")
	return nil
}

// Shares a specific post with a specific user.
func SharePost(postId, userId int) error {
	statement, err := Db.Prepare("INSERT INTO shared_posts (post_id, user_id) VALUES (?, ?);")
	if err != nil {
		fmt.Println("GetData2.go SharePost err 1: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(postId, userId)
	if err != nil {
		fmt.Println("GetData2.go AddFollower err 2: ", err)
		return err
	}
	fmt.Println("Post", postId, "shared with", userId, "successfully!")
	return nil
}

// Retrieves all posts that that the user themselves have made.
func GetOwnPosts(userId int, data []customstructs.PostsData) ([]customstructs.PostsData, error) {
	fmt.Println("Getting user's own posts...")

	if len(data) == 0 {
		err := fmt.Errorf("getData2.go GetPublicPosts: allPosts is empty")
		return nil, err
	}

	var returnData []customstructs.PostsData

	for _, v := range data {
		if v.Creator.Id == userId && v.CategoryId != 1 {
			returnData = append(returnData, v)
		}
	}

	return returnData, nil
}

// Retrieves image links from the database based on the condition string.
func GetImageUrlFromDb(condition, id string) ([]string, error) {
	rows, err := Db.Query("SELECT image_URL FROM images WHERE "+condition+"=$1", id)
	if err != nil {
		fmt.Println("getData2.go err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	var returnData []string
	for rows.Next() {
		var i string
		err := rows.Scan(&i)
		if err != nil {
			fmt.Println("getData2.go err 1: ", err)
			return nil, err
		}
		fmt.Println("getData2.go GetImageUrlFromDb i: ", i)
		returnData = append(returnData, i)
	}

	if len(returnData) == 0 {
		err := fmt.Errorf("image not found")
		return nil, err
	}

	return returnData, nil
}

// Gets one single post from the database with comments and images.
func GetPost(pId string, userId int) ([]customstructs.PostsData, error) {
	posts, err := GetPostsFromDb("post_id", pId)
	if err != nil {
		fmt.Println("postview.go GetPost() error getting posts from db: ", err)
		return nil, err
	}

	if len(posts) == 0 {
		fmt.Printf("\nA post with the id of %v does not exist in the database\n", pId)
		err := fmt.Errorf("a post with the id of %v does not exist in the database", pId)
		return nil, err
	}

	userData, err := GetUser("users.user_id", fmt.Sprint(posts[0].Creator.Id))
	if err != nil {
		fmt.Println("postview.go GetPost GetUser err: ", err)
		return nil, err
	}
	userData = ProcessUserInfo(userData, userId)
	posts[0].Creator = userData[0]

	comments := []customstructs.CommentData{}
	commentsFromDb, err := GetCommentsFromDb("post_id", pId)
	if err != nil {
		fmt.Println("postview.go GetPost GetCommentsFromDb err: ", err)
		return nil, err
	}
	comments = append(comments, commentsFromDb...)

	postImageUrl, err := GetImageUrlFromDb("post_id", pId)
	if err == nil {
		fmt.Println("GetImageUrlFromDb getting image...")
		posts[0].ImageUrl = postImageUrl[0]
	}

	posts[0].Comments = comments

	for index, c := range posts[0].Comments {
		userData, err := GetUser("users.user_id", fmt.Sprint(c.Creator.Id))
		if err != nil {
			fmt.Println("postview.go GetPost GetUserFromDb err: ", err)
			return nil, err
		}

		commentImageUrl, err := GetImageUrlFromDb("comment_id", fmt.Sprint(c.Id))
		if err == nil {
			fmt.Println("postview.go GetPost GetImageUrlFromDb: comment image found!")
			posts[0].Comments[index].ImageUrl = commentImageUrl[0]
		}

		posts[0].Comments[index].Creator = ProcessUserInfo(userData, userId)[0]
	}

	return posts, nil
}

// Retrieves user data from the database with added optional userdata.
func GetUser(condition, value string) ([]customstructs.UserData, error) {
	returnData := []customstructs.UserData{}

	var (
		rows *sql.Rows
		err  error
	)

	if condition == "" {
		rows, err = Db.Query("SELECT users.*, optional_userdata.*, user_visibility.visibility FROM users, optional_userdata, user_visibility WHERE users.user_id=optional_userdata.user_id AND users.user_id=user_visibility.user_id")
	} else {
		rows, err = Db.Query("SELECT users.*, optional_userdata.*, user_visibility.visibility FROM users, optional_userdata, user_visibility WHERE users.user_id=optional_userdata.user_id AND users.user_id=user_visibility.user_id AND "+condition+"=$1", value)
	}

	if err != nil {
		fmt.Println("getData2.go GetUser err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			tmpData customstructs.UserData
			trash   any
		)

		rows.Scan(&tmpData.Id, &tmpData.Email, &tmpData.Password, &tmpData.FirstName, &tmpData.LastName, &tmpData.DateOfBirth, &tmpData.OnlineStatus, &tmpData.RoleId, &trash, &tmpData.ImageUrl, &tmpData.Nickname, &tmpData.AboutMe, &tmpData.Visibility)
		tmpData.Password = ""
		returnData = append(returnData, tmpData)
	}

	return returnData, nil
}

// Retrieves optional userdata (nickname, about me, avatar image) from the database.
func GetOptionalUserdata(userId int) ([]customstructs.UserData, error) {
	rows, err := Db.Query("SELECT image_URL, nickname, about_me FROM optional_userdata WHERE user_id=$1", userId)
	if err != nil {
		fmt.Println("getData2.go GetOptionalUserdata err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	var returnData []customstructs.UserData

	for rows.Next() {
		var tmpData customstructs.UserData
		rows.Scan(&tmpData.ImageUrl, &tmpData.Nickname, &tmpData.AboutMe)
		returnData = append(returnData, tmpData)
	}

	return returnData, nil
}

// Removes a post from the database.
func RemovePost(condition, value string) {
	statement, err := Db.Prepare("DELETE FROM posts WHERE " + condition + " = ?")
	if err != nil {
		fmt.Println("getData2.go RemovePost err 1: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(value)
	if err != nil {
		fmt.Println("getData2.go RemovePost err 2: ", err)
		return
	}

	fmt.Printf("post %v removed\n", value)
}
