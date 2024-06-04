package exec

import (
	"database/sql"
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"path"
)

// Handler for "/group/" endpoint.
//
// Method: GET - gets a specific group info from db, groupId comes from request URL.
//
// Method: POST - creates a new group based on info sent from the front as FormData.
func ServeGroup(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("groups.go ServeGroup err 2: ", err)
		msg, _ := json.Marshal("You need to log in first!")
		w.Write(msg)
		return
	}

	user := GetUserByCookie(cookie.Value)
	if user == 0 {
		fmt.Println("groups.go ServeGroup: RoleID is 0")
		msg, _ := json.Marshal("incorrect cookie value")
		w.Write(msg)
		return
	}

	switch r.Method {
	case "GET":
		groupIdString := path.Base(r.URL.Path)
		fmt.Println("groupIdString: ", groupIdString)

		var returnData []customstructs.GroupData

		if groupIdString == "group" {
			fmt.Println("groups.go ServeGroup: group ID not specified, getting all groups")

			returnData, err = GetGroupFromDb("", "")
			if err != nil {
				fmt.Println("groups.go ServeGroup err 3: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			for i, v := range returnData {
				returnData[i].Creator = ProcessUserInfo(v.Creator, user)
			}

		} else {
			returnData, err = GetGroupFromDb("group_id", groupIdString)
			if err != nil {
				fmt.Println("groups.go ServeGroup err 3: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			if len(returnData) == 0 {
				err := fmt.Errorf("group with the id of %v does not exist", groupIdString)
				json.NewEncoder(w).Encode(err.Error())
				return
			}

			returnData[0].Creator = ProcessUserInfo(returnData[0].Creator, user)
		}

		if len(returnData) == 0 {
			err := fmt.Errorf("groups not found")
			fmt.Println("groups.go ServeGroup err 7: ", err)
			msg, _ := json.Marshal(err.Error())
			w.Write(msg)
			return
		}

		for i, v := range returnData {
			groupMembers, err := GetGroupMembersFromDb(fmt.Sprint(v.GroupId))
			if err != nil {
				fmt.Println("groups.go ServeGroup err 4: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			for _, v2 := range groupMembers {
				returnData[i].Members = append(returnData[i].Members, ProcessUserInfo([]customstructs.UserData{v2}, user)...)
			}
		}

		if groupIdString != "group" {
			userHasAccess := false

			for _, v := range returnData[0].Members {
				if v.Id == user {
					userHasAccess = true
				}
			}

			if !userHasAccess {
				err := fmt.Errorf("user %v does not have access to group %v", user, returnData[0].GroupId)
				fmt.Println("groups.go ServeGroup err 4: ", err)
				err = fmt.Errorf("user does not have access to group ")
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}
		}

		for i, v := range returnData {
			groupPosts, err := GetPostsFromDb("group_id", v.GroupId)
			if err != nil {
				fmt.Println("groups.go ServeGroup err 5: ", err)
				msg, _ := json.Marshal(err.Error())
				w.Write(msg)
				return
			}

			for _, v2 := range groupPosts {
				p, err := GetPost(fmt.Sprint(v2.Id), user)
				if err == nil {
					returnData[i].Posts = append(returnData[i].Posts, p...)
				}
			}
		}

		msg, _ := json.Marshal(returnData)
		w.Write(msg)
	case "POST":
		groupTitle := r.FormValue("title")
		groupDesc := r.FormValue("body")

		groupId, err := CreateNewGroup(user, groupTitle, groupDesc)
		if err != nil {
			fmt.Println("groups.go ServeGroup err 6: ", err)
			msg, _ := json.Marshal(err.Error())
			w.Write(msg)
		}

		msg, _ := json.Marshal("group " + fmt.Sprint(groupId[0]) + " created!")
		w.Write(msg)
	}
}

// Retrieves group data from the database based on the condition string.
func GetGroupFromDb(condition string, value interface{}) ([]customstructs.GroupData, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if condition == "" {
		rows, err = Db.Query("SELECT * FROM group_index;")
	} else {
		rows, err = Db.Query("SELECT * FROM group_index WHERE "+condition+"=$1;", value)
	}

	if err != nil {
		fmt.Println("groups.go GetGroupFromDb err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	var returnData []customstructs.GroupData

	for rows.Next() {
		var t customstructs.GroupData
		var userId int
		err = rows.Scan(&t.GroupId, &userId, &t.GroupTitle, &t.GroupDesc)
		if err != nil {
			fmt.Println("groups.go GetGroupFromDb err 2: ", err)
			return nil, err
		}

		creator, err := GetUser("users.user_id", fmt.Sprint(userId))
		if err != nil {
			fmt.Println("groups.go GetGroupFromDb err 3: ", err)
			return nil, err
		}
		creator = ProcessUserInfo(creator, -1)
		t.Creator = creator

		returnData = append(returnData, t)
	}

	return returnData, nil
}

// Retrieves group members from the database based on the condition string.
func GetGroupMembersFromDb(groupId string) ([]customstructs.UserData, error) {
	var returnData []customstructs.UserData

	rows, err := Db.Query("SELECT user_id FROM group_members WHERE group_id=$1", groupId)
	if err != nil {
		fmt.Println("groups.go GetGroupFromDb err 1: ", err)
		return nil, err
	}

	var memberIds []int

	for rows.Next() {
		var t int
		err = rows.Scan(&t)
		if err != nil {
			fmt.Println("groups.go GetGroupFromDb err 2: ", err)
			return nil, err
		}
		memberIds = append(memberIds, t)
	}
	rows.Close()

	for _, v := range memberIds {
		userData, err := GetUser("users.user_id", fmt.Sprint(v))
		if err != nil {
			fmt.Println("groups.go GetGroupFromDb err 3: ", err)
			break
		}
		returnData = append(returnData, userData...)
	}

	return returnData, nil
}

// Inserts a new group into the database.
func CreateNewGroup(userId int, title, description string) ([]int, error) {
	statement, err := Db.Prepare("INSERT INTO group_index (user_id, group_title, group_description) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("groups.go CreateNewGroup err 1: ", err)
		return nil, err
	}
	defer statement.Close()

	res, err := statement.Exec(userId, title, description)
	if err != nil {
		fmt.Println("groups.go CreateNewGroup err 2: ", err)
		return nil, err
	}

	i, err := res.LastInsertId()
	if err != nil {
		fmt.Println("groups.go CreateNewGroup err 3: ", err)
		return nil, err
	}

	AddGroupMember(fmt.Sprint(i), fmt.Sprint(userId))

	return []int{int(i)}, nil
}

// Adds a specific user to the specific group's members list.
func AddGroupMember(groupId, userId string) {
	statement, err := Db.Prepare("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)")
	if err != nil {
		fmt.Println("groups.go AddGroupMember err 1: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(groupId, userId)
	if err != nil {
		fmt.Println("groups.go AddGroupMember err 2: ", err)
		return
	}
	fmt.Println(userId, "added to group", groupId)
}

// Removes a specific user from the specific group's members list.
func RemoveGroupMember(groupId, userId string) {
	statement, err := Db.Prepare("DELETE FROM group_members WHERE (group_id, user_id) = (?, ?)")
	if err != nil {
		fmt.Println("groups.go RemoveGroupMember err 1: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(groupId, userId)
	if err != nil {
		fmt.Println("groups.go RemoveGroupMember err 2: ", err)
		return
	}

	fmt.Println(userId, "removed from group", groupId)
}

// Removes a group and all of its entries from database.
func DeleteGroup(groupId string) {
	group, err := GetGroupFromDb("group_id", groupId)
	if err != nil {
		fmt.Println("DeleteGroup err 1:", err)
		return
	}

	members, err := GetGroupMembersFromDb(groupId)
	if err != nil {
		fmt.Println("DeleteGroup err 2:", err)
		return
	}

	//remove all posts from group
	RemovePost("group_id", groupId)

	//remove all group events
	DeleteEvent("group_id", groupId)

	//remove all dm-s in group
	deleteMessageFromDb("group_id", groupId)

	//remove all group notifications
	removeNotificationFromDb("notification_value", groupId)

	n := Notification{
		Sender:       "group",
		Desc:         "groupDeleted",
		SeenByTarget: "0",
		SeenBySender: "1",
		Status:       "accepted",
		Value:        group[0].GroupTitle,
	}

	for _, v := range members {
		n.Target = fmt.Sprint(v.Id)
		insertNotificationIntoDb(n)
	}

	sendNotificationsToGroup(groupId)

	//remove all members from group
	for _, v := range members {
		RemoveGroupMember(groupId, fmt.Sprint(v.Id))
	}

	_, err = Db.Exec("DELETE FROM group_index WHERE group_id = ?", groupId)
	if err != nil {
		fmt.Println("DeleteGroup err 3:", err)
		return
	}

	fmt.Printf("group %v deleted\n", groupId)
}
