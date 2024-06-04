package exec

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler for "/login/" endpoint.
// Calls StartSession if valid credentials are provided.
// Only "POST" method is supported.
func Login(w http.ResponseWriter, r *http.Request) {
	var (
		invalid bool
		reason  string
	)

	var userDataFromJson struct {
		Email    string
		Password string
	}

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userDataFromJson)

		if err != nil {
			fmt.Println("login.go Login(): error decoding json into struct: ", err)
			return
		}

		reason = ValidateLogin(userDataFromJson.Password, userDataFromJson.Email)

		if len(reason) == 0 {
			fmt.Println("Validate: true")
			StartSession(w, r, userDataFromJson.Email)
			msg, _ := json.Marshal("login successful")
			w.Write(msg)
			return
		} else {
			invalid = true
		}
	}

	if invalid {
		fmt.Println("Validate: false")
		msg, _ := json.Marshal(reason)
		w.Write(msg)
	}
}
