package exec

import (
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// for upgrading the HTTP connection into a websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var WebsocketServer = NewServer()

type Client struct {
	Conn   *websocket.Conn
	UserID int
	Server *Server

	//channel for communication between users
	Recieve chan DirectMessage

	//for messages that don't get sent to the front
	Misc chan DirectMessage
}

// Shared between all logged in users
type Server struct {
	mu         sync.Mutex
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

type MessageRaw struct {
	Type       string          `json:"messageType"`
	MessageRaw json.RawMessage `json:"body"`
}

// Creates a new Server struct
func NewServer() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
	}
}

// Goroutine for writing data to the websocket connection. Unique per client.
func (c *Client) writer() {
	fmt.Println("Starting writer...")
	defer func() {
		fmt.Println("writer shutting down...")
		close(c.Misc)
		close(c.Recieve)

		c.Conn.Close()
		fmt.Println("writer: goodbye")
	}()

WriterLoop:
	for {
		select {
		case serverMessage := <-c.Misc:
			switch serverMessage.Msg {
			case "logout":
				//for when a user logs out

				fmt.Println("writer: logout")
				break WriterLoop
			case "newUser":
				//for when a user logs in

				fmt.Println("writer: newUser")

				user, err := GetUser("users.user_id", fmt.Sprint(c.UserID))
				if err != nil {
					fmt.Println("chat.go err 1: ", err)
					return
				}

				user = ProcessUserInfo(user, user[0].Id)

				tmpArr := user[0].Followers
				tmpArr = append(tmpArr, user[0].Following...)
				uniqueUsersMap := make(map[int]bool)
				for _, v := range tmpArr {
					uniqueUsersMap[v.Id] = true
				}

				uniqueUsersArr := []customstructs.UserData{}
				for k := range uniqueUsersMap {
					user, err := GetUser("users.user_id", fmt.Sprint(k))
					if err != nil {
						fmt.Println("chat.go err 2: ", err)
						return
					}
					user = ProcessUserInfo(user, c.UserID)
					uniqueUsersArr = append(uniqueUsersArr, user...)
				}

				serverMessageJson, _ := json.Marshal(NewMessage("server", "", "allUsers", uniqueUsersArr))
				c.Conn.WriteMessage(websocket.TextMessage, serverMessageJson)
			default:
				fmt.Println("writer: unknown server message: ", serverMessage.Msg)
			}

		case msg, ok := <-c.Recieve:
			if !ok {
				fmt.Println("Error reading message from Recieve channel!")
				c.Server.Unregister <- c
				continue WriterLoop
			}

			fmt.Printf("chat.go writer: Recieve channel message for %v: %v\n", c.UserID, msg.Msg)

			c.Server.mu.Lock()

			switch msg.Msg {
			case "getAllDMs":
				allUserMessages := getAllUserMessages(c.UserID)
				allGroupMessages := getAllGroupMessages(c.UserID)

				allMessages := append(allUserMessages, allGroupMessages...)

				allMessages = ProcessAllDirectMessages(allMessages)

				serverMessageJson, _ := json.Marshal(NewMessage("server", fmt.Sprint(c.UserID), "getAllDMs", allMessages))
				err := c.Conn.WriteMessage(websocket.TextMessage, serverMessageJson)
				if err != nil {
					fmt.Println("chat.go writer c.Conn.WriteMessage err 1: ", err)
				}
			case "getAllNotifications", "deleteGroup", "unFollowRequest", "groupInvite", "followRequest", "groupJoinRequest", "event", "leaveGroup", "kickFromGroup":
				allNotifications := getAllUserNotifications(c.UserID)

				allNotifications = ProcessAllNotifications(allNotifications) // []Notification, int

				serverMessageJson, _ := json.Marshal(NewMessage("server", fmt.Sprint(c.UserID), "getAllNotifications", allNotifications))
				err := c.Conn.WriteMessage(websocket.TextMessage, serverMessageJson)
				if err != nil {
					fmt.Println("chat.go writer c.Conn.WriteMessage err 2: ", err)
				}
			default:
				fmt.Println("writerLoop: unknown message in c.Recieve: ", msg.Msg)
			}

			c.Server.mu.Unlock()

			continue WriterLoop
		}
	}
}

// Goroutine for reading data from the websocket connection. Unique per client.
func (c *Client) reader() {
	defer func() {
		fmt.Println("reader: shutting down")
	}()
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			if err.Error()[:16] == "websocket: close" {
				fmt.Println("reader: Client closed the connection from the front.")
			} else {
				fmt.Println("reader: Error reading from frontend: ", err)
			}
			c.Server.Unregister <- c
			break
		}
		switch messageType {
		case websocket.TextMessage:
			var message MessageRaw
			json.Unmarshal(p, &message)
			handleIncomingMessages(message, c)
		default:
			fmt.Println("reader: Message type not supported: ", messageType)
		}
	}
}

// Checks if the user is present in the "users" table
func doesUserExist(target string) bool {
	allUsers, err := GetUserFromDb("", "")
	if err != nil {
		fmt.Println("doesUserExist: Error getting all users: ", err)
	}
	for _, i := range allUsers {
		if fmt.Sprint(i.Id) == target {
			return true
		}
	}
	return false
}

// Goroutine for registering and removing clients to/from the Server. Only one server exists at a time.
// Activates the "reader" and "writer" goroutines for each new Client.
func (s *Server) Run() {
	defer func() {
		fmt.Println("chat.go server.Run: SHUTTING DOWN!")
	}()
	for {
		select {
		case client := <-s.Register:
			fmt.Println("chat.go Server.Run(): registering new client...")
			s.mu.Lock()

			s.Clients[client] = true
			go client.reader()
			go client.writer()

			client.Recieve <- NewMessage("server", "", "getAllDMs", []interface{}{})
			client.Recieve <- NewMessage("server", "", "getAllNotifications", []interface{}{})

			for c := range s.Clients {
				c.Misc <- NewMessage("server", "all", "newUser", []interface{}{})
			}

			//fmt.Println("registered users in server:", s.Clients)

			s.mu.Unlock()
			fmt.Println("chat.go Server.Run(): new client registered!")
		case client := <-s.Unregister:
			fmt.Println("chat.go Server.Run(): unregistering client...")
			s.mu.Lock()

			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				for c := range s.Clients {
					c.Misc <- NewMessage("server", "all", "newUser", []interface{}{})
				}
				client.Misc <- NewMessage("server", "all", "logout", []interface{}{})
			}
			s.mu.Unlock()
			fmt.Println("chat.go Server.Run(): Client has been unregistered!")
		}
	}
}

// Upgrades the HTTP connection into a websocket connection and sets up the client struct for real-time communication goroutines.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	server := WebsocketServer
	fmt.Println()
	fmt.Println("ServeWs called")
	fmt.Println()
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		fmt.Println("ServeWs cookie err: ", cookieErr)
		return
	}

	conn, UpgradeErr := upgrader.Upgrade(w, r, nil)
	if UpgradeErr != nil {
		fmt.Println("Upgrading err: ", UpgradeErr)
		return
	}

	userId := GetUserBySessionID(cookie.Value).UserID

	client := Client{
		Conn:    conn,
		UserID:  userId,
		Server:  server,
		Recieve: make(chan DirectMessage, 256),
		Misc:    make(chan DirectMessage, 256),
	}

	server.Register <- &client
	fmt.Println("serveWs: done")
}

// gets all sent and recieved messages by user from db
func getAllUserMessages(userId int) []DirectMessage {
	var returnData []DirectMessage
	rows, _ := Db.Query("SELECT * FROM direct_messages WHERE (sender_id = $1 OR target_id = $1) AND group_id = 0", userId)
	defer rows.Close()

	for rows.Next() {
		var m DirectMessage
		err := rows.Scan(&m.MessageId, &m.Sender, &m.Target, &m.Msg, &m.SeenByTarget, &m.DateTime, &m.GroupId)
		if err != nil {
			fmt.Println("chat.go getAllUserMessages error: ", err)
			return returnData
		}
		returnData = append(returnData, m)
	}
	return returnData
}

// Gets all messages for groups that the user is a member of from the database.
func getAllGroupMessages(userId int) []DirectMessage {
	var returnData []DirectMessage

	rows, _ := Db.Query("SELECT direct_messages.* FROM direct_messages, group_members WHERE direct_messages.group_id = group_members.group_id AND group_members.user_id = $1", userId)
	defer rows.Close()

	for rows.Next() {
		var m DirectMessage
		err := rows.Scan(&m.MessageId, &m.Sender, &m.Target, &m.Msg, &m.SeenByTarget, &m.DateTime, &m.GroupId)
		if err != nil {
			fmt.Println("chat.go getAllUserMessages error: ", err)
			return returnData
		}
		returnData = append(returnData, m)
	}
	return returnData
}

// Routes incoming messages from the "reader" goroutine.
func handleIncomingMessages(message MessageRaw, c *Client) {
	switch message.Type {
	case "message":
		fmt.Println("chat.go message.Type: message")
		var msg DirectMessage
		json.Unmarshal(message.MessageRaw, &msg)
		HandleIncomingDMs(msg, c)
	case "notification":
		fmt.Println("chat.go message.Type: notification")
		var notif Notification
		json.Unmarshal(message.MessageRaw, &notif)
		fmt.Println("chat.go handleIncomingMessages unmarshaled notification: ", notif)
		RouteIncomingNotifications(notif, c)
	default:
		fmt.Println("chat.go handleIncomingMessages: unknown message type: ", message.Type)
	}
}
