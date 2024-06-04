package exec

import (
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"time"

	uuid "github.com/gofrs/uuid"
)

// Creates a new session UUID string and adds it to the "session" cookie for authentication.
func StartSession(w http.ResponseWriter, r *http.Request, email string) {
	fmt.Println("starting session...")
	var user []customstructs.UserData
	var err error

	sessionTokenUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("session.go err 2: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionToken := sessionTokenUUID.String()

	user, err = GetUserFromDb("email", email)
	if err != nil {
		fmt.Println("session.go error 1: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userid := user[0].Id

	// err = RemoveSession("user_id", userid)
	// if err != nil {
	// 	fmt.Println("session.go error 2: ", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	addSessionToDB(sessionToken, user[0].Email, userid)

	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	fmt.Println("StartSession end")
}

// Adds session info to the database.
func addSessionToDB(session, email string, userid int) {
	fmt.Println("addSessionToDB start")
	statement, err := Db.Prepare("INSERT INTO sessions (session_id, user_id, email, datetime, role_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("addSessionToDB err 1: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(session, userid, email, time.Now(), 1)
	if err != nil {
		fmt.Println("addSessionToDB err 2: ", err)
		return
	}
	fmt.Println("addSessionToDB end")
}

// Removes expired sessions from the database.
func RemoveSession(condition string, value interface{}) error {
	_, err := Db.Exec("DELETE FROM sessions WHERE "+condition+" = $1", value)
	if err != nil {
		return err
	}
	return nil
}

// Checks permission level of the user.
func Authentication(fn http.HandlerFunc, auth int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			if auth == 0 {
				fn.ServeHTTP(w, r)
				return
			} else if auth == 1 {
				msg, _ := json.Marshal("You need to log in first!")
				w.Write(msg)
				return
			} else {
				fmt.Println("session.go Authentication: unknown auth value: ", auth)
				return
			}
		}

		session := GetUserBySessionID(cookie.Value)

		if len(fmt.Sprint(session.RoleID)) == 0 {
			fmt.Println("session.go Authentication: malformed session cookie: ", cookie.Value)
			msg, _ := json.Marshal("malformed session cookie")
			w.Write(msg)
			return
		}

		roleId := session.UserID

		if auth == 0 && roleId > auth {
			fmt.Println("Guests only")
			msg, _ := json.Marshal("Guests only!")
			w.Write(msg)
			return
		}

		if auth == 1 && roleId < auth {
			fmt.Println("Logged in only")
			msg, _ := json.Marshal("You need to log in first!")
			w.Write(msg)
			return
		}

		fn.ServeHTTP(w, r)
	}
}
