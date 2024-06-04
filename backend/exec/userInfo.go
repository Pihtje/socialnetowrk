package exec

import (
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Handler for "/userInfo/" endpoint.
//
// GET - retrieves data for the currently signed in user, based on the session cookie.
//
// POST - updates optional userdata for the currently signed in user.
func ServeUserInfo(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("userinfo.go ServeUserInfo: session cookie err: ", err)
		msg, _ := json.Marshal("userinfo.go ServeUserInfo: session cookie missing")
		w.Write(msg)
		return
	}

	userId := GetUserByCookie(cookie.Value)
	if userId == 0 {
		err := fmt.Errorf("incorrect cookie value")
		fmt.Println("userinfo.go ServeUserInfo: ", err)
		msg, _ := json.Marshal(err.Error())
		w.Write(msg)
		return
	}

	switch r.Method {
	case "GET":
		userIdFromPath := path.Base(r.URL.Path)
		var userInfo []customstructs.UserData

		switch userIdFromPath {
		case "self":
			userInfo, err := GetUser("users.user_id", fmt.Sprint(userId))
			if err != nil {
				fmt.Println("userInfo.go ServeUserInfo err 9: ", err.Error())
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			userInfo = ProcessUserInfo(userInfo, userId)

			msg, _ := json.Marshal(userInfo)
			w.Write(msg)
			return
		default:
			userIdFromPathInt, err := strconv.Atoi(userIdFromPath)
			if err != nil {
				fmt.Println("userInfo.go ServeUserInfo: userIdFromPath from URL missing or malformed, getting all users")
				userInfo, err = GetUser("", "")
				if err != nil {
					fmt.Println("userInfo.go ServeUserInfo err 7: ", err.Error())
					msg, _ := json.Marshal(err.Error())
					w.Write(msg)
					return
				}
			} else {
				userInfo, err = GetUser("users.user_id", fmt.Sprint(userIdFromPathInt))
				if err != nil {
					fmt.Println("userInfo.go ServeUserInfo err 8: ", err.Error())
					msg, _ := json.Marshal(err.Error())
					w.Write(msg)
					return
				}
			}
		}

		ProcessUserInfo(userInfo, userId)

		msg, _ := json.Marshal(userInfo)
		w.Write(msg)
	case "POST":
		nickname := r.FormValue("nickname")
		fmt.Println("nickname 1:", nickname)
		aboutMe := r.FormValue("aboutMe")
		userVisibility := r.FormValue("visibility")

		var avatarPath string
		var newAvatarPath string

		avatarImage, avatarHeader, err := r.FormFile("avatarImage")
		if err != nil && err != http.ErrMissingFile {
			err = fmt.Errorf("error reading avatar")
			json.NewEncoder(w).Encode(err.Error())
			return
		} else if err == nil {
			avatarPath, err = SaveUserAvatar(avatarImage, avatarHeader)
			if err != nil {
				err = fmt.Errorf("error saving avatar")
				json.NewEncoder(w).Encode(err.Error())
				return
			}
		}

		if avatarPath != "" {
			oldFullPath := avatarPath

			fmt.Println("oldFullPath", oldFullPath)

			newAvatarPath = strings.Replace(avatarPath, "user_pic_", fmt.Sprintf("user_pic_%v_", userId), 1)

			fmt.Println("newFullPath", newAvatarPath)

			err = RenameImage(oldFullPath, newAvatarPath)
			if err != nil {
				fmt.Println("userInfo.go ServeUserInfo err 2: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			imgStatement, err := Db.Prepare("INSERT INTO images (user_id, image_URL) VALUES (?,?);")
			if err != nil {
				fmt.Println("userInfo.go ServeUserInfo err 3: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}
			defer imgStatement.Close()

			_, err = imgStatement.Exec(userId, newAvatarPath)
			if err != nil {
				fmt.Println("userInfo.go ServeUserInfo err 4: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}
		}

		UpdateOptionalUserdata(newAvatarPath, nickname, aboutMe, fmt.Sprint(userId))

		if userVisibility != "" {
			UpdateUserVisibility(userId, userVisibility)
		}

		msg, _ := json.Marshal("optional info updated!")
		w.Write(msg)
	default:
		fmt.Println("userInfo.go ServeUserInfo: unknown method: ", r.Method)
		msg, _ := json.Marshal("userInfo.go ServeUserInfo: unknown method: " + r.Method)
		w.Write(msg)
	}
}

// Updates optional userdata in the database for the user.
func UpdateOptionalUserdata(imageURL, nickname, aboutMe, userId string) {
	fmt.Println("updating optional userdata...")
	fmt.Println("imageURL:", imageURL)
	if len(strings.TrimSpace(imageURL)+strings.TrimSpace(nickname)+strings.TrimSpace(aboutMe)) == 0 {
		fmt.Println("userdata.go UpdateOptionalUserdata: nothing to update!")
		return
	}

	_, err := GetUserFromDb("user_id", userId)
	if err != nil {
		fmt.Println("userdata.go UpdateOptionalUserdata err 1: ", err)
		return
	}

	statementString := "UPDATE optional_userdata SET "

	s1, s2, s3 := strings.TrimSpace(imageURL), strings.TrimSpace(nickname), strings.TrimSpace(aboutMe)

	if s1 != "" {
		s1 = "image_URL = \"" + s1 + "\""
		if s2 != "" || s3 != "" {
			s1 += ","
		}
	}

	if s2 != "" {
		s2 = "nickname = \"" + s2 + "\""
		if s3 != "" {
			s2 += ","
		}
	}

	if s3 != "" {
		s3 = "about_me = \"" + s3 + "\""
	}

	statementString = statementString + s1 + s2 + s3 + " WHERE user_id = ?"

	statement, err := Db.Prepare(statementString)
	if err != nil {
		fmt.Println("userdata.go UpdateOptionalUserdata err 2: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(userId)
	if err != nil {
		fmt.Println("userdata.go UpdateOptionalUserdata err 3: ", err)
		return
	}

	fmt.Println("userinfo.go UpdateOptionalUserinfo: update successful!")
}

// Changes user account visibility.
// Available visibility values: "public", "private".
func UpdateUserVisibility(userId int, userVisibility string) {
	fmt.Println("userInfo.go UpdateUserVisibility: updating user visibility...")
	statement, err := Db.Prepare("UPDATE user_visibility SET visibility=? WHERE user_id=?")
	if err != nil {
		fmt.Println("userInfo.go UpdateUserVisibility: err 1: ", err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec(userVisibility, userId)
	if err != nil {
		fmt.Println("userInfo.go UpdateUserVisibility: err 2: ", err)
		return
	}
	fmt.Println("userInfo.go UpdateUserVisibility: done")
}

// Adds posts and comments to userinfo, filters visible fields based on the user's access level.
// Available access levels: "none", "public", "follower", "full".
func ProcessUserInfo(userInfo []customstructs.UserData, userId int) []customstructs.UserData {
	for i1, v1 := range userInfo {
		userHasAccess := "none"
		f, _ := GetFollowersFromDb([]customstructs.UserData{v1})
		userInfo[i1] = f[0]

		if v1.Id == userId {
			userHasAccess = "full"
		} else {
			for _, v2 := range userInfo[i1].Followers {
				if v2.Id == userId {
					userHasAccess = "follower"
				}
			}
		}

		if userInfo[i1].Visibility == "public" && userHasAccess == "none" {
			userHasAccess = "public"
		}

		if userHasAccess == "none" {

			minUserInfo := customstructs.UserData{}

			minUserInfo.Id = userInfo[i1].Id
			minUserInfo.Nickname = userInfo[i1].Nickname
			minUserInfo.ImageUrl = userInfo[i1].ImageUrl
			minUserInfo.Visibility = userInfo[i1].Visibility
			minUserInfo.RoleId = userInfo[i1].RoleId
			minUserInfo.FirstName = userInfo[i1].FirstName
			minUserInfo.LastName = userInfo[i1].LastName

			userInfo[i1] = minUserInfo
			continue
		}

		userInfo, err := GetFollowingFromDb(userInfo)
		if err != nil {
			fmt.Println("userInfo.go ProcessUserInfo err 2: ", err)
		}

		userPosts, err := GetPostsFromDb("user_id", v1.Id)
		if err != nil {
			fmt.Println("userInfo.go ProcessUserInfo err 3: ", err)
		}

		userComments, err := GetCommentsFromDb("user_id", v1.Id)
		if err != nil {
			fmt.Println("userInfo.go ProcessUserInfo err 4: ", err)
		}
		userInfo[i1].Comments = userComments

		switch userHasAccess {
		case "full":
			userInfo[i1].Posts = userPosts
		case "follower":
			publicPosts, err := GetPublicPosts(userPosts)
			if err != nil {
				fmt.Println("userInfo.go ProcessUserInfo err 5: ", err)
			}
			userInfo[i1].Posts = append(userInfo[i1].Posts, publicPosts...)

			followingPosts, err := GetFollowingPosts(userId, userPosts)
			if err != nil {
				fmt.Println("userInfo.go ProcessUserInfo err 6: ", err)
			}
			userInfo[i1].Posts = append(userInfo[i1].Posts, followingPosts...)

			allSharedPosts, err := GetSharedPosts(userId)
			if err != nil {
				fmt.Println("userInfo.go ProcessUserInfo err 7: ", err)
			} else {
				userInfo[i1].Posts = append(userInfo[i1].Posts, allSharedPosts...)
			}
		case "public":
			//fmt.Println("userHasAccess: ", userHasAccess)

			publicPosts, err := GetPublicPosts(userPosts)
			if err != nil {
				fmt.Println("userInfo.go ProcessUserInfo err 8: ", err)
			}
			userInfo[i1].Posts = append(userInfo[i1].Posts, publicPosts...)
		default:
			err := fmt.Errorf("userinfo.go ServeUserInfo: unknown value in userHasAccess: %v", userHasAccess)
			fmt.Println(err)
		}
	}
	return userInfo
}
