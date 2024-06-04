package exec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Inserts a post into the database.
func Post(uId, cId int, title, body string, sharedUsers []int, groupId int) ([]int, error) {
	fmt.Println("Within Post function")
	statement, err := Db.Prepare("INSERT INTO posts (user_id, post_datetime, title, body, category_id, group_id) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Println("post.go Post err 1: ", err)
		return nil, err
	}
	defer statement.Close()
	res, err := statement.Exec(uId, time.Now(), title, body, cId, groupId)
	if err != nil {
		fmt.Println("post.go Post err 2: ", err)
		return nil, err
	}
	if cId == 3 {
		postId, _ := res.LastInsertId()
		for _, v := range sharedUsers {
			SharePost(int(postId), v)
		}
	}
	userId, err := res.LastInsertId()
	if err != nil {
		fmt.Println("post.go Post err 3: ", err)
		return nil, err
	}
	fmt.Println("Post added successfully!")
	return []int{int(userId)}, nil
}

// Handler for "/newPost/" endpoint.
// Creates a new post based on data provided by the user though the FormData sent with the POST request.
// Only POST method is supported.
func CreatePost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		fmt.Println("Logged in only")
		msg, _ := json.Marshal("You need log in first!")
		w.Write(msg)
		return
	}

	userdata := GetUserByCookie(cookie.Value)
	if userdata == 0 {
		fmt.Println("Logged in only")
		msg, _ := json.Marshal("incorrect cookie value")
		w.Write(msg)
		return
	}

	switch r.Method {
	case "POST":
		// if userdata.RoleID < 1 {
		// 	fmt.Println("Logged in only")
		// 	msg, _ := json.Marshal("You need log in first!")
		// 	w.Write(msg)
		// 	return
		// }

		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			fmt.Println("Error parsing form: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		title := r.FormValue("Title")
		body := r.FormValue("Body")
		categoryString := r.FormValue("Category")
		shareWithString := r.FormValue("ShareWith")
		groupIdString := r.FormValue("GroupId")

		fmt.Println("categoryString: ", categoryString, "shareWithString: ", shareWithString)

		category, err := strconv.Atoi(categoryString)
		if err != nil {
			fmt.Println("post.go CreatePost err 1: ", err)
			errMsg, _ := json.Marshal(err.Error())
			w.Write(errMsg)
			return
		}

		shareWithArr := strings.Split(shareWithString, ",")
		var sharedUsers []int
		for _, v := range shareWithArr {
			u, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("post.go CreatePost sharing err: ", err)
			}
			sharedUsers = append(sharedUsers, u)
		}

		var errMsg int

		if len([]int{category}) == 0 {
			errMsg = 3
		} else if len(body) < 10 || len(body) > 500 {
			errMsg = 2
		} else if len(title) < 5 || len(title) > 75 {
			errMsg = 1
		}

		groupId, err := strconv.Atoi(groupIdString)
		if err != nil {
			fmt.Println("groupId not found, post added to default group")
			groupId = 0
		}

		if errMsg == 0 {
			postId, err := Post(userdata, category, title, body, sharedUsers, groupId)
			if err != nil {
				fmt.Println("post.go Createpost err 3: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
			}
			files, ok := r.MultipartForm.File["postImages"]
			if ok && len(files) > 0 {
				savedPaths, err := SavePostImages(files, postId[0])
				if err != nil {
					fmt.Println("Error saving post images: ", err)
					return
				}
				for _, path := range savedPaths {
					fmt.Println("Image saved at:", path)
					_, err := SaveImagePath(postId[0], path)
					if err != nil {
						fmt.Println("Error saving image path to database: ", err)
					}
				}
			}
			msg, _ := json.Marshal("New post added!")
			w.Write(msg)
		} else {
			switch errMsg {
			case 1:
				msg, _ := json.Marshal("Title must be between 5 and 75 characters long")
				w.Write(msg)
			case 2:
				msg, _ := json.Marshal("Your post must be between 10 and 500 characters long")
				w.Write(msg)
			case 3:
				msg, _ := json.Marshal("Please choose at least one category")
				w.Write(msg)
			}
		}
	case "GET":
		fmt.Println("CreatePost recieved a get request!")
		msg, _ := json.Marshal("CreatePost recieved a get request!")
		w.Write(msg)
		return
	}
}
