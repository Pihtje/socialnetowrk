package exec

import (
	"fmt"
	"forum/customstructs"
)

type Notification struct {
	NotificationId string `json:"notificationId"`
	Sender         any    `json:"sender"`
	Target         any    `json:"target"`
	Desc           string `json:"desc"` //"followRequest", "groupInvite", "unFollowRequest", "leaveGroup", "deleteGroup"
	SeenByTarget   string `json:"seenByTarget"`
	SeenBySender   string `json:"seenBySender"`
	Status         string `json:"status"`          //"pending", "accepted", "rejected"
	Value          any    `json:"value,omitempty"` // group id
}

// Routes incoming notifications based on the "Desc" field value
func RouteIncomingNotifications(n Notification, c *Client) {
	n.Sender = fmt.Sprint(n.Sender)
	n.Target = fmt.Sprint(n.Target)
	n.Value = fmt.Sprint(n.Value)

	switch n.Desc {
	case "followRequest":
		handleFollowRequests(n, c)
	case "unFollowRequest":
		handleUnfollowRequests(n, c)
	case "groupInvite": //for when user is invited to a group by another user in that group
		handleGroupInvite(n, c)
	case "groupJoinRequest": //for when user wants to join a group
		handleGroupJoinRequest(n, c)
	case "leaveGroup":
		handleLeavingGroup(n, c)
	case "deleteGroup", "groupDeleted":
		handleGroupDeletion(n, c)
	case "event":
		handleEventResponse(n, c)
	default:
		fmt.Println("unsupported notification desc: ", n.Desc)
	}
}

// Handler for sending and accepting follow requests
func handleFollowRequests(n Notification, c *Client) {
	fmt.Printf("new follower request for %v from %v\n", n.Target, n.Sender)
	switch n.Status {
	case "pending":
		targetUser, err := GetUser("users.user_id", fmt.Sprint(n.Target))
		if err != nil {
			fmt.Println("wsNotifications.go handleFollowRequests: targetUser does not exist: ", n.Target)
			return
		}

		if targetUser[0].Visibility == "public" {
			n.Status = "accepted"
			n.SeenByTarget = "1"
			insertNotificationIntoDb(n)
			handleFollowRequests(n, c)
			return
		}

		insertNotificationIntoDb(n)

		msg := NewMessage(n.Sender.(string), n.Target.(string), "getAllNotifications", n)
		c.Server.mu.Lock()
		fmt.Println("handleFollowRequests pending c.Server.Clients : ", c.Server.Clients)
		for i := range c.Server.Clients {
			fmt.Println("handleFollowRequests pending i.UserID: ", i.UserID)
			if fmt.Sprint(i.UserID) == n.Target.(string) || fmt.Sprint(i.UserID) == n.Sender.(string) {
				i.Recieve <- msg
			}
		}
		c.Server.mu.Unlock()
	case "accepted", "rejected":
		updateNotificationsTable(n)

		if n.SeenBySender == "1" && n.SeenByTarget == "1" {
			fmt.Println("handleFollowRequests: both sender and target have seen notification")
		} else {
			if n.Status == "accepted" {
				AddFollower(n.Target.(string), n.Sender.(string))
			}

			if n.Status == "rejected" {
				removeNotificationFromDb("notification_id", n.NotificationId)
			}
		}

		msg1 := NewMessage(n.Sender.(string), n.Target.(string), n.Desc, n)
		msg2 := NewMessage(n.Sender.(string), n.Target.(string), "newUser", []any{})
		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			fmt.Println("i: ", i.UserID)
			if fmt.Sprint(i.UserID) == n.Target.(string) || fmt.Sprint(i.UserID) == n.Sender.(string) {
				fmt.Println("sending to client")
				i.Recieve <- msg1
				i.Misc <- msg2
			}
		}
		c.Server.mu.Unlock()
	default:
		fmt.Println("wsNotifications.go: unknown notification status: ", n.Status)
		return
	}
}

// Handler for removing followers
func handleUnfollowRequests(n Notification, c *Client) {
	fmt.Println("handleUnfollowRequests start")
	RemoveFollower(n.Target.(string), n.Sender.(string))
	msg := NewMessage(n.Sender.(string), n.Target.(string), "newUser", []any{})
	c.Server.mu.Lock()
	for i := range c.Server.Clients {
		fmt.Println("handleFollowRequests unFollowRequest i.UserID: ", i.UserID)
		if fmt.Sprint(i.UserID) == n.Target.(string) || fmt.Sprint(i.UserID) == n.Sender.(string) {
			i.Misc <- msg
		}
	}
	c.Server.mu.Unlock()
	fmt.Println("handleUnfollowRequests end")
}

// Inserts the notification into the "notifications" table.
func insertNotificationIntoDb(n Notification) {
	fmt.Println("wsNotifications.go insertNotificationIntoDb: inserting notification into db...")
	statement, err := Db.Prepare("INSERT INTO notifications (sender_id, target_id, description, seen_by_target, seen_by_sender, notification_status, notification_value) VALUES (?,?,?,?,?,?,?);")
	if err != nil {
		fmt.Println("insertNotificationIntoDb err 1: ", err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec(n.Sender.(string), n.Target.(string), n.Desc, n.SeenByTarget, n.SeenBySender, n.Status, n.Value.(string))
	if err != nil {
		fmt.Println("insertNotificationIntoDb err 2: ", err)
		return
	}
	fmt.Println("wsNotifications.go insertNotificationIntoDb: done")
}

// Removes the notification from the "notifications" table based on the condition string.
func removeNotificationFromDb(condition, value string) {
	fmt.Println("wsNotifications.go removeNotificationFromDb: removing notification from db...")
	statement, err := Db.Prepare("DELETE FROM notifications WHERE " + condition + " = ?;")
	if err != nil {
		fmt.Println("removeNotificationFromDb: Error removing direct message from db: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(value)
	if err != nil {
		fmt.Println("removeNotificationFromDb err 2: ", err)
		return
	}
	fmt.Println("wsNotifications.go removeNotificationFromDb: done")
}

// Updates the values of a notification already in the "notifications" table.
func updateNotificationsTable(n Notification) {
	fmt.Println("updating notifications table...")
	if n.SeenBySender == "1" && n.SeenByTarget == "1" {
		fmt.Println("updateNotificationsTable: notification seen by both sender and target, removing from db")
		removeNotificationFromDb("notification_id", fmt.Sprint(n.NotificationId))
		return
	}

	statement, err := Db.Prepare("UPDATE notifications SET (seen_by_target, seen_by_sender, notification_status) = (?, ?, ?) WHERE notification_id = ?")
	if err != nil {
		fmt.Println("updateNotificationsTable err 1: ", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(n.SeenByTarget, n.SeenBySender, n.Status, n.NotificationId)
	if err != nil {
		fmt.Println("updateNotificationsTable err 2: ", err)
		return
	}

	fmt.Println("done updating notifications table")
}

// Retrieves all notifications that involve the user from the database.
func getAllUserNotifications(userId int) []Notification {
	var returnData []Notification
	rows, err := Db.Query("SELECT * FROM notifications WHERE (target_id = $1 OR (sender_id = $1 AND description != $2))", userId, "event")
	if err != nil {
		fmt.Println("wsNotifications.go getAllUserNotifications error: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var n Notification
		err := rows.Scan(&n.NotificationId, &n.Sender, &n.Target, &n.Desc, &n.SeenByTarget, &n.SeenBySender, &n.Status, &n.Value)
		if err != nil {
			fmt.Println("wsNotifications.go getAllUserNotifications error: ", err)
			return returnData
		}

		n.Sender = fmt.Sprint(n.Sender)
		n.Target = fmt.Sprint(n.Target)
		n.Value = fmt.Sprint(n.Value)

		returnData = append(returnData, n)
	}
	return returnData
}

// Handler for when a member of a group sends an invite to another user
func handleGroupInvite(n Notification, c *Client) {
	if n.Value == "" {
		fmt.Println("handleGroupInvite: n.Value missing!")
		return
	}

	fmt.Printf("new groupInvite for %v from %v\n", n.Target, n.Sender)

	group, err := GetGroupFromDb("group_id", n.Value)
	if err != nil {
		fmt.Println("handleGroupInvite: unknown group id: ", n.Value)
		return
	}

	msg := NewMessage(n.Sender.(string), n.Target.(string), n.Desc, n)

	switch n.Status {
	case "pending":
		fmt.Println("handleGroupInvite: pending")

		groupMembers, err := GetGroupMembersFromDb(fmt.Sprint(group[0].GroupId))
		if err != nil {
			fmt.Println("handleGroupInvite: error getting group members: ", err)
			return
		}

		isSenderPartOfGroup := false

		for _, v := range groupMembers {
			if n.Target == fmt.Sprint(v.Id) {
				fmt.Printf("user %v is already a member of group %v\n", n.Target, n.Value)
				return
			}

			if n.Sender.(string) == fmt.Sprint(v.Id) {
				isSenderPartOfGroup = true
			}
		}

		if !isSenderPartOfGroup {
			fmt.Printf("user %v is not part of the group and cannot send group invites\n", n.Sender)
			return
		}

		insertNotificationIntoDb(n)

		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			if fmt.Sprint(i.UserID) == msg.Target {
				i.Recieve <- msg
			}
		}
		c.Server.mu.Unlock()
	case "accepted", "rejected":
		fmt.Println("handleGroupInvite: accepted, rejected")

		updateNotificationsTable(n)

		if n.SeenBySender == "1" && n.SeenByTarget == "1" {
			fmt.Println("handleGroupInvite: both sender and target have seen notification")
		} else {
			if n.Status == "accepted" {
				AddGroupMember(fmt.Sprint(group[0].GroupId), msg.Target.(string))
			}

			if n.Status == "rejected" {
				removeNotificationFromDb("notification_id", n.NotificationId)
			}
		}

		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			//fmt.Println(i.UserID)
			if fmt.Sprint(i.UserID) == msg.Sender || fmt.Sprint(i.UserID) == msg.Target {
				//fmt.Println("sending...")
				i.Recieve <- msg
			}
		}
		c.Server.mu.Unlock()
	default:
		fmt.Println("handleGroupInvite: unknown notification status: ", n.Status)
		return
	}
}

// Handles removing user from group.
func handleLeavingGroup(n Notification, c *Client) {
	fmt.Println("handleLeavingGroup start")

	if n.Value == "" {
		fmt.Println("handleLeavingGroup: n.Value missing!")
		return
	}

	group, err := GetGroupFromDb("group_id", n.Value)
	if err != nil {
		fmt.Println("handleLeavingGroup err 1: ", err)
		return
	}

	switch n.Target {
	case "group":
		//for user leaving the group themselves

		if n.Sender.(string) == fmt.Sprint(group[0].Creator[0].Id) {
			fmt.Println("The creator of a group cannot leave without deleting the group first!")
			return
		}

		RemoveGroupMember(n.Value.(string), n.Sender.(string))
		msg1 := NewMessage(n.Sender.(string), n.Target.(string), n.Desc, n)
		msg2 := NewMessage(n.Sender.(string), n.Target.(string), "getAllDMs", n)
		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			//fmt.Println("handleFollowRequests leaveGroup i.UserID: ", i.UserID)
			if fmt.Sprint(i.UserID) == n.Sender.(string) {
				i.Recieve <- msg1
				i.Recieve <- msg2
				break
			}
		}
		c.Server.mu.Unlock()
	default:
		fmt.Println("handleLeavingGroup: unknown n.Target:", n.Target)
	}
	fmt.Println("handleLeavingGroup end")
}

// Handler for deleting a group by the group creator.
func handleGroupDeletion(n Notification, c *Client) {
	if n.Value == "" {
		fmt.Println("handleGroupDeletion: n.Value missing!")
		return
	}

	if n.Desc == "groupDeleted" && n.SeenByTarget == "1" {
		removeNotificationFromDb("notification_id", n.NotificationId)
		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			if fmt.Sprint(i.UserID) == n.Target.(string) {
				i.Recieve <- NewMessage("server", fmt.Sprint(i.UserID), "getAllNotifications", true)
				break
			}
		}
		c.Server.mu.Unlock()
		return
	}

	group, err := GetGroupFromDb("group_id", n.Value)
	if err != nil {
		fmt.Println("handleGroupDeletion err 1: ", err)
		return
	}

	if fmt.Sprint(group[0].Creator[0].Id) != n.Sender.(string) {
		fmt.Printf("handleGroupDeletion: user %v does not have authority to delete group %v\n", n.Sender, n.Value)
		return
	}

	DeleteGroup(n.Value.(string))
}

// Checks if group id is valid before sending notifications to group members.
func handleGroupEventNotifications(e customstructs.GroupEvent) {
	fmt.Println("handleGroupEventNotifications start")
	if len(fmt.Sprint(e.Group.(customstructs.GroupData).GroupId)) == 0 {
		fmt.Println("handleGroupEventNotifications: incorrect groupId: ", e.Group.(customstructs.GroupData).GroupId)
		return
	}

	sendNotificationsToGroup(fmt.Sprint(e.Group.(customstructs.GroupData).GroupId))

	fmt.Println("handleGroupEventNotifications end")
}

// Sends notification to all group members.
func sendNotificationsToGroup(groupId string) {
	members, err := GetGroupMembersFromDb(groupId)
	if err != nil {
		fmt.Println("handleGroupEventNotifications err 1: ", err)
		return
	}

	WebsocketServer.mu.Lock()
	for _, v1 := range members {
		msg := NewMessage("group", fmt.Sprint(v1.Id), "getAllNotifications", "")
		for i := range WebsocketServer.Clients {
			if i.UserID == v1.Id {
				i.Recieve <- msg
				break
			}
		}
	}
	WebsocketServer.mu.Unlock()
}

// Adds structs to notification fields.
func ProcessAllNotifications(allNotifs []Notification) []Notification {
	fmt.Println("ProcessAllNotifications start")
	for i1, v1 := range allNotifs {

		sUser, err := GetUser("users.user_id", fmt.Sprint(v1.Sender))
		if err == nil && len(sUser) != 0 {
			allNotifs[i1].Sender = sUser[0]
		}

		tUser, err := GetUser("users.user_id", fmt.Sprint(v1.Target))
		if err == nil && len(tUser) != 0 {
			allNotifs[i1].Target = tUser[0]
		}

		switch v1.Desc {
		case "groupInvite":
			group, err := GetGroupFromDb("group_id", v1.Value)
			if err != nil {
				fmt.Println("ProcessAllNotifications err 4: ", err)
				return nil
			}
			allNotifs[i1].Value = group[0]
		case "groupJoinRequest":
			group, err := GetGroupFromDb("group_id", v1.Value)
			if err != nil {
				fmt.Println("ProcessAllNotifications err 4: ", err)
				return nil
			}
			allNotifs[i1].Value = group[0]
			allNotifs[i1].Target = "group"
		case "event":
			event, err := GetEventFromDb("event_id", v1.Value.(string))
			if err != nil {
				fmt.Println("ProcessAllNotifications err 6: ", err)
				return nil
			}
			allNotifs[i1].Value = event[0]
			if allNotifs[i1].SeenByTarget == "1" {
				allNotifs[i1].SeenBySender = "1"
			}
		default:
			fmt.Println("ProcessAllNotifications: unknown v1.Desc value:", v1.Desc, ", using default")
		}
	}
	fmt.Println("ProcessAllNotifications end")
	return allNotifs
}

// Handler for when a user wants to join a new group.
func handleGroupJoinRequest(n Notification, c *Client) {
	if n.Value == "" {
		fmt.Println("handleGroupJoinRequest: n.Value missing!")
		return
	}

	fmt.Printf("new groupJoinRequest for group %v from %v\n", n.Value, n.Sender)

	group, err := GetGroupFromDb("group_id", n.Value)
	if err != nil {
		fmt.Println("handleGroupJoinRequest: unknown group id: ", n.Value)
		return
	}

	msg := NewMessage(n.Sender.(string), n.Target.(string), "getAllNotifications", n)

	switch n.Status {
	case "pending":
		fmt.Println("handleGroupJoinRequest: pending")

		groupMembers, err := GetGroupMembersFromDb(fmt.Sprint(group[0].GroupId))
		if err != nil {
			fmt.Println("handleGroupJoinRequest: error getting group members: ", err)
			return
		}

		isSenderPartOfGroup := false

		for _, v := range groupMembers {
			if n.Target == fmt.Sprint(v.Id) {
				fmt.Printf("user %v is already a member of group %v\n", n.Target, n.Value)
				return
			}

			if n.Sender.(string) == fmt.Sprint(v.Id) {
				isSenderPartOfGroup = true
			}
		}

		if isSenderPartOfGroup {
			fmt.Printf("user %v is already a member of group %v\n", n.Sender, n.Value)
			return
		}

		if !isSenderPartOfGroup {
			fmt.Println("handleGroupJoinRequest group creator id: ", group[0].Creator[0].Id)
			n.Target = fmt.Sprint(group[0].Creator[0].Id)
			msg.Target = fmt.Sprint(group[0].Creator[0].Id)
		}

		insertNotificationIntoDb(n)

		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			if fmt.Sprint(i.UserID) == msg.Target {
				fmt.Println("msg before sending:", msg)
				i.Recieve <- msg
			}
		}
		c.Server.mu.Unlock()

	case "accepted", "rejected":
		fmt.Println("handleGroupJoinRequest: accepted, rejected")

		updateNotificationsTable(n)

		if n.SeenBySender == "1" && n.SeenByTarget == "1" {
			fmt.Println("handleGroupInvite: both sender and target have seen notification")
		} else {
			if n.Status == "accepted" {
				AddGroupMember(fmt.Sprint(group[0].GroupId), n.Sender.(string))
			}
			if n.Status == "rejected" {
				removeNotificationFromDb("notification_id", n.NotificationId)
			}
		}

		c.Server.mu.Lock()
		for i := range c.Server.Clients {
			if fmt.Sprint(i.UserID) == msg.Target || fmt.Sprint(i.UserID) == msg.Sender {
				i.Recieve <- msg
			}
		}
		c.Server.mu.Unlock()

	default:
		fmt.Println("handleGroupInvite: unknown notification status: ", n.Status)
		return
	}
}

// Handler for sending event responses.
func handleEventResponse(n Notification, c *Client) {
	fmt.Println("handleEventResponse start")

	responseString := ""
	switch n.Status {
	case "accepted":
		responseString = "Attending"
	case "rejected":
		responseString = "Not attending"
	default:
		fmt.Println("handleEventResponse unknown notification status:", n.Status)
		return
	}

	responseObject := customstructs.EventResponse{
		EventID:  n.Value.(string),
		User:     n.Target.(string),
		Response: responseString,
	}

	err := RespondToEvent(responseObject)
	if err != nil {
		fmt.Println("handleEventResponse err 1:", err)
		return
	}

	if responseString == "Not attending" {
		removeNotificationFromDb("notification_id", n.NotificationId)
	} else {
		updateNotificationsTable(n)
	}

	if n.SeenBySender == "1" && n.SeenByTarget == "1" {
		fmt.Println("handleEventResponse: both sender and target have seen notification")
	}

	msg := NewMessage(n.Sender.(string), n.Target.(string), n.Desc, n)

	c.Server.mu.Lock()
	for i := range c.Server.Clients {
		if fmt.Sprint(i.UserID) == n.Target.(string) {
			i.Recieve <- msg
			break
		}
	}
	c.Server.mu.Unlock()

	fmt.Println("handleEventResponse end")
}
