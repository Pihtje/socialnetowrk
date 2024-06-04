package main

import (
	"encoding/json"
	"fmt"
	"forum/db/sqliteDb"
	"forum/exec"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqliteDb.Initialise()
	sqliteDb.ApplyMigrationsUp()
	sqliteDb.PopulateUsers()
	sqliteDb.InitToPosts()
	sqliteDb.CreateGroups()

	//http.Handle("/", http.FileServer(http.Dir("../frontend/dist")))
	http.HandleFunc("/signup/", enableCors(exec.Authentication(exec.SignupPage, 0)))
	http.HandleFunc("/login/", enableCors(exec.Authentication(exec.Login, 0)))
	http.HandleFunc("/logout/", enableCors(exec.Authentication(exec.Logout, 1)))
	http.HandleFunc("/newPost/", enableCors(exec.Authentication(exec.CreatePost, 1)))
	http.HandleFunc("/post/", enableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { exec.PostView(w, r) })))
	http.HandleFunc("/group/", enableCors(exec.Authentication(exec.ServeGroup, 1)))
	http.HandleFunc("/userInfo/", enableCors(exec.Authentication(exec.ServeUserInfo, 1)))
	http.HandleFunc("/allPosts/", enableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { exec.ServeFilteredPosts(w, r) })))
	http.HandleFunc("/createGroupEvent/", enableCors(exec.Authentication(exec.CreateGroupEventHandler, 1)))
	http.HandleFunc("/getUserEventStatus/", enableCors(exec.Authentication(exec.GetUserEventStatusHandler, 1)))

	// For serving image files.
	http.HandleFunc("/backend/images/", func(w http.ResponseWriter, r *http.Request) {
		imagePath := r.URL.Path[len("/backend/images/"):]

		folderPaths := map[string]string{
			"comment_images": filepath.Join("images", "comment_images"),
			"post_images":    filepath.Join("images", "post_images"),
			"user_images":    filepath.Join("images", "user_images"),
		}

		parts := strings.Split(imagePath, "/")
		folderName := parts[0]

		folderPath, ok := folderPaths[folderName]
		if !ok {
			err := fmt.Errorf("invalid image folder")
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		fullImagePath := filepath.Join(folderPath, strings.Join(parts[1:], "/"))

		file, err := os.Open(fullImagePath)
		if err != nil {
			err = fmt.Errorf("image not found")
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		defer file.Close()

		contentType := getContentType(fullImagePath)
		if contentType == "" {
			err = fmt.Errorf("unsupported image type")
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		fmt.Println(contentType)
		w.Header().Set("Content-Type", contentType)

		_, err = io.Copy(w, file)
		if err != nil {
			err = fmt.Errorf("error serving image")
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	})

	go exec.WebsocketServer.Run()

	http.HandleFunc("/chat/", enableCors(exec.Authentication(exec.ServeWs, 1)))

	fmt.Println("http://localhost:8000/")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("main.go: ", err)
		return
	}
}

// https://www.stackhawk.com/blog/golang-cors-guide-what-it-is-and-how-to-enable-it/
// some of the header options can be removed
func enableCors(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		fmt.Println("enableCors origin: ", origin)

		if origin == "http://localhost:8080" || origin == "http://localhost:5173" || origin == "http://127.0.0.1:8080" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "Authorisation")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("credentials", "same-origin")
		w.Header().Add("credentials", "include")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

// For getting the image extension.
func getContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}
