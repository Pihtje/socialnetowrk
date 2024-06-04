package exec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Handler for "/signup/" endpoint.
// Creates a new user based on FormData sent by the frontend in a POST request.
// Only POST requests are supported.
func SignupPage(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("in SignupPage")

	var (
		errMsg string
		valid  bool = true
		reg    bool
		email  string
	)

	if r.Method == "POST" {
		email = r.FormValue("email")
		password := r.FormValue("password")
		repeatPassword := r.FormValue("repeatPassword")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		dateOfBirth := r.FormValue("dateOfBirth")
		username := r.FormValue("nickname")
		aboutMe := r.FormValue("aboutMe")

		visibility := r.FormValue("visibility")

		errMsg, valid = ValidateRegister(email, password, repeatPassword, dateOfBirth)

		if valid {
			avatarImage, avatarHeader, err := r.FormFile("avatarImage")

			var imagePath string

			if err != nil && err != http.ErrMissingFile {
				err = fmt.Errorf("error reading avatar")
				json.NewEncoder(w).Encode(err.Error())
				return
			} else if err == nil {
				imagePath, err = SaveUserAvatar(avatarImage, avatarHeader)
				if err != nil {
					err = fmt.Errorf("error saving avatar")
					json.NewEncoder(w).Encode(err.Error())
					return
				}
			}

			err = Signup(email, password, firstName, lastName, dateOfBirth, imagePath, username, aboutMe, visibility)
			if err != nil {
				fmt.Println("signup.go error: ", err)
				json.NewEncoder(w).Encode(err.Error())
				return
			}

			fmt.Println("SignupPage 6")

			reg = true
		}
	}

	if !valid {
		errMsg, _ := json.Marshal(errMsg)
		w.Write(errMsg)
	} else if reg {
		StartSession(w, r, email)
		msg, _ := json.Marshal("Registration successful!")
		w.Write(msg)
	}
}

// Checks if the registration data is valid.
func ValidateRegister(email, password, repeatpass, dateOfBirth string) (string, bool) {
	sl1 := strings.Split(email, "@")
	if len(sl1) < 2 {
		return "Invalid email address", false
	}

	sl2 := strings.Split(sl1[1], ".")
	if len(sl2) < 2 {
		return "Invalid email address", false
	}

	userData, _ := GetUserFromDb("", "")

	if email == userData[0].Email {
		return "Username or email already taken!", false
	}

	if len(password) < 8 {
		return "Your password must be at least 8 characters long!", false
	}

	if password != repeatpass {
		return "Passwords do not match!", false
	}

	t, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		fmt.Println("signup.go ValidateRegister err 1:", err)
		return "Invalid date format!", false
	}

	if t.After(time.Now()) {
		fmt.Println("signup.go ValidateRegister err 2: birth date is in the future")
		return "Birth date cannot be in the future!", false
	}

	return "", true
}

// Inserts a new user into the database.
func Signup(email, password, first_name, last_name, date_of_birth, avatarPath, nickname, about_me, visibility string) error {
	fmt.Println("Signup func")
	statement, err := Db.Prepare("INSERT INTO users (email, password, first_name, last_name, date_of_birth, online_status, role_id) VALUES (?,?,?,?,?,?,?);")
	if err != nil {
		fmt.Println("signup.go Signup err 1: ", err)
		return err
	}
	defer statement.Close()

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		fmt.Println("signup.go Signup err 2: ", err)
		return err
	}

	res, err := statement.Exec(email, encryptedPass, first_name, last_name, date_of_birth, "offline", 1)
	if err != nil {
		fmt.Println("signup.go Signup err 3: ", err)
		return err
	}

	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("signup.go Signup err 4: ", err)
		return err
	}

	err = RemoveExistingUserAvatar(lastInsertedID)
	if err != nil {
		fmt.Println("Error removing existing avatar:", err)
	}

	var newAvatarPath string

	if avatarPath != "" {
		oldFullPath := avatarPath
		newAvatarPath = strings.Replace(avatarPath, "user_pic_", fmt.Sprintf("user_pic_%d_", lastInsertedID), 1)

		err = RenameImage(oldFullPath, newAvatarPath)
		if err != nil {
			fmt.Println("signup.go Signup err 5: ", err)
			return err
		}

		imgStatement, err := Db.Prepare("INSERT INTO images (user_id, image_URL) VALUES (?,?);")
		if err != nil {
			fmt.Println("signup.go Signup err 6: ", err)
			return err
		}
		defer imgStatement.Close()

		_, err = imgStatement.Exec(lastInsertedID, newAvatarPath)
		if err != nil {
			fmt.Println("signup.go Signup err 7: ", err)
			return err
		}
	}
	UpdateOptionalUserdata(newAvatarPath, nickname, about_me, fmt.Sprint(lastInsertedID))
	UpdateUserVisibility(int(lastInsertedID), visibility)
	return nil
}
