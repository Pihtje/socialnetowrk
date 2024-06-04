package exec

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler for "/logout/" endpoint.
// Calls RemoveSession and deletes the session cookie.
func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout start")

	sessionToken, err := r.Cookie("session")
	if err != nil {
		msg, _ := json.Marshal("You're already logged out!")
		w.Write(msg)
		return
	}

	RemoveSession("session_id", sessionToken.Value)

	newCookie := http.Cookie{
		Name:     "session",
		Value:    sessionToken.Value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		MaxAge:   -1,
	}

	http.SetCookie(w, &newCookie)

	fmt.Println("Logout end")
	msg, _ := json.Marshal("You have been logged out")
	w.Write(msg)
}
