package exec

import (
	"fmt"
	"time"
)

type DirectMessage struct {
	MessageId    string `json:"messageId"`
	Sender       any    `json:"sender,omitempty"`
	Target       any    `json:"target"`
	Msg          string `json:"msg"`
	ServerObject any
	SeenByTarget bool
	DateTime     time.Time
	GroupId      string //0 is normal dm, anything else is chat message in a specific group
}

// Routes incoming direct messages based on the "Target" and "Msg" fields.
func HandleIncomingDMs(dm DirectMessage, c *Client) {
	fmt.Println("wsMessages.go dm: ", dm)
	switch dm.Target {
	case "server":
		switch dm.Msg {
		case "logout":
			fmt.Println("handleIncomingDMs: logout")
			c.Server.Unregister <- c
			return
		case "getAllDMs":
			c.Recieve <- dm
			return
		}

		fmt.Println("chat.go handleIncomingDMs: unknown server dm: ", dm)
	default:
		fmt.Println("wsMessages.go default dm")
		switch dm.GroupId {
		case "0":
			fmt.Println("wsMessages.go default dm group 0")
			if !doesUserExist(dm.Target.(string)) {
				fmt.Printf("User with the id of %v does not exist!\n", dm.Target)
				return
			}

			insertMsgIntoDb(&dm)

			fmt.Println("Sending msg from front to target's reader...")

			dm.Msg = "getAllDMs"

			c.Server.mu.Lock()
			for v := range c.Server.Clients {
				if fmt.Sprint(v.UserID) == dm.Target || fmt.Sprint(v.UserID) == dm.Sender {
					v.Recieve <- dm
				}
			}
			c.Server.mu.Unlock()
		default:
			fmt.Println("wsMessages.go default dm group msg")
			group, err := GetGroupFromDb("group_id", dm.GroupId)
			if err != nil || len(group) == 0 {
				fmt.Printf("\nws handleIncomingDMs: group %v not found in db", dm.GroupId)
				return
			}

			members, err := GetGroupMembersFromDb(fmt.Sprint(group[0].GroupId))
			if err != nil {
				fmt.Printf("\ngroup %v does not have any members", group[0].GroupId)
				return
			}

			isUserPartOfGroup := false
			for _, v := range members {
				if fmt.Sprint(v.Id) == fmt.Sprint(dm.Sender) {
					isUserPartOfGroup = true
					break
				}
			}

			if !isUserPartOfGroup {
				err := fmt.Errorf("user %v is not a member of group %v", dm.Sender, group[0].GroupId)
				fmt.Println(err)
				return
			}

			insertMsgIntoDb(&dm)

			c.Server.mu.Lock()
			for _, v1 := range members {
			innerLoop:
				for v2 := range c.Server.Clients {
					if v2.UserID == v1.Id {
						v2.Recieve <- NewMessage("", dm.Target.(string), "getAllDMs", true)
						fmt.Println("dm sent to", v2.UserID)
						break innerLoop
					}
				}
			}
			c.Server.mu.Unlock()
		}
		fmt.Println("done sending msg from handleIncomingMessages")
	}
}

// Creates a new DirectMessage struct.
func NewMessage(sender, target, msg string, serverObject interface{}) DirectMessage {
	return DirectMessage{
		Sender:       sender,
		Target:       target,
		Msg:          msg,
		ServerObject: serverObject,
		SeenByTarget: false,
		DateTime:     time.Now(),
	}
}

// Inserts a DirectMessage into the "direct_messages" table.
func insertMsgIntoDb(msg *DirectMessage) {
	fmt.Println("chat.go insertMsgIntoDb: inserting message into db...")
	statement, err := Db.Prepare("INSERT INTO direct_messages (sender_id, target_id, message, message_seen, message_datetime, group_id) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		fmt.Println("insertMsgIntoDb: Error adding direct message to db: ", err)
		return
	}
	defer statement.Close()
	statement.Exec(msg.Sender, msg.Target, msg.Msg, 0, time.Now(), msg.GroupId)
	fmt.Println("chat.go insertMsgIntoDb: done")
}

// Removes a DirectMessage from the "direct_messages" table.
func deleteMessageFromDb(condition, value string) {
	fmt.Println("chat.go deleteMessageFromDb start")
	statement, err := Db.Prepare("DELETE FROM direct_messages WHERE " + condition + " = $1;")
	if err != nil {
		fmt.Println("deleteMessageFromDb: Error deleting message from db: ", err)
		return
	}
	defer statement.Close()
	statement.Exec(value)
	fmt.Println("chat.go deleteMessageFromDb: done")
}

// Replaces the DirectMessage.Sender and DirectMessage.Target string values with UserData values.
func ProcessAllDirectMessages(allMessages []DirectMessage) []DirectMessage {
	for i, v := range allMessages {
		sUser, err := GetUser("users.user_id", fmt.Sprint(v.Sender))
		if err == nil && len(sUser) != 0 {
			allMessages[i].Sender = sUser[0]
		}
		tUser, err := GetUser("users.user_id", fmt.Sprint(v.Target))
		if err == nil && len(tUser) != 0 {
			allMessages[i].Target = tUser[0]
		}
	}
	return allMessages
}
