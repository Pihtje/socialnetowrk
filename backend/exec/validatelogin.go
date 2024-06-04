package exec

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Checks if the data provided through the login endpoint is valid.
func ValidateLogin(pword, email string) string {
	var userid int
	var err error

	if len(email) == 0 {
		return "Wrong email or password"
	}

	userid, err = checkEmail(email)
	if err != nil {
		return "Wrong email or password"
	}

	if !checkPassword(userid, pword) {
		return "Wrong email or password"
	}

	// if userid != 0 {
	// 	//return "User is already logged in"
	// 	fmt.Println("User is already logged in")
	// 	RemoveSession("email", email)
	// }

	return ""
}

// Compares the hashed password in the database with the one provided by the user.
func checkPassword(user_id int, password string) bool {
	statement, err := Db.Prepare(`SELECT password FROM users WHERE user_id = ?`)
	if err != nil {
		fmt.Println("Validatelogin.go: error in checkPassword: ", err)
		return false
	}
	defer statement.Close()

	var pass string
	statement.QueryRow(user_id).Scan(&pass)
	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)); err != nil {
		return false
	}
	return true
}

// Compares the email in the database with the one provided by the user.
func checkEmail(email string) (user_id int, err error) {
	statement, err := Db.Prepare(`SELECT user_id FROM users WHERE email = ?`)
	if err != nil {
		fmt.Println("Validatelogin.go: error in checkEmail: ", err)
		return 0, err
	}
	defer statement.Close()

	var user int
	err = statement.QueryRow(email).Scan(&user)
	if err == sql.ErrNoRows {
		fmt.Println("Validatelogin.go: cannot find email in db: ", email)
		return 0, err
	}
	return user, nil
}
