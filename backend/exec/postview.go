package exec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// Handler for "/post/" endpoint.
//
// GET - retrieves data for a single post from the database.
//
// POST - handles comment submissions.
func PostView(w http.ResponseWriter, r *http.Request) {
	// var userdata customstructs.PageTemplate

	postID := path.Base(r.URL.Path)

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("PostView err 1:", err)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	userdata := GetUserByCookie(cookie.Value)
	if userdata == 0 {
		fmt.Println("PostView err 2: incorrect cookie value")
		json.NewEncoder(w).Encode("incorrect cookie value")
		return
	}

	post, err := GetPost(postID, userdata)
	if err != nil {
		fmt.Println("postview.go GetPost err: ", err)
		msg, _ := json.Marshal(err.Error())
		w.Write(msg)
		return
	}

	switch r.Method {
	case "POST":
		if userdata == 0 {
			fmt.Println("Logged in only")
			msg, _ := json.Marshal("You need to be logged in to add a comment!")
			w.Write(msg)
			return
		}

		var comment struct {
			CommentBody string
		}

		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			fmt.Println("Failed to parse multipart form: ", err)
			msg, _ := json.Marshal(err.Error())
			w.Write(msg)
			return
		}

		comment.CommentBody = r.FormValue("newCommentBody")

		if len(strings.TrimSpace(comment.CommentBody)) == 0 {
			msg, _ := json.Marshal("You cannot post an empty comment!")
			w.Write(msg)
			return
		}

		if len(comment.CommentBody) > 255 {
			msg, _ := json.Marshal("Your comment cannot be longer than 255 characters!")
			w.Write(msg)
			return
		}

		commentId, err := Comment(post[0].Id, userdata, comment.CommentBody)
		if err != nil {
			fmt.Println("postView.go PostView: Error adding comment: ", err)
			msg, _ := json.Marshal(err.Error())
			w.Write(msg)
			return
		}

		files, ok := r.MultipartForm.File["commentImages"]
		if ok && len(files) > 0 {
			savedPaths, err := SaveCommentImages(files, commentId[0])
			if err != nil {
				fmt.Println("Error saving comment images: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			for _, path := range savedPaths {
				_, err := SaveCommentImagePath(commentId[0], path)
				if err != nil {
					fmt.Println("Error saving image path to database:", err)
					msg, _ := json.Marshal(err.Error())
					w.Write(msg)
					return
				}
			}
		}
		msg, _ := json.Marshal("Comment and images added!")
		w.Write(msg)
	case "GET":
		userdataJSON, _ := json.Marshal(userdata)
		w.Write(userdataJSON)
	}
}
