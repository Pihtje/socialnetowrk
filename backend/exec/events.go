package exec

import (
	"database/sql"
	"encoding/json"
	"fmt"
	customstructs "forum/customstructs"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Creates a group event in the database if event details are valid.
func CreateGroupEvent(event customstructs.GroupEvent) (int, error) {
	fmt.Println("CreateGroupEvent start, event: ", event)
	if event.Title == "" || event.Description == "" {
		return 0, fmt.Errorf("title and Description cannot be empty")
	}

	if event.DayTime.Before(time.Now()) {
		return 0, fmt.Errorf("event time cannot be in the past")
	}

	result, err := Db.Exec("INSERT INTO group_events (group_id, title, description, day_time) VALUES (?, ?, ?, ?)", event.Group.(customstructs.GroupData).GroupId, event.Title, event.Description, event.DayTime)
	if err != nil {
		return 0, fmt.Errorf("failed to create group event: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %v", err)
	}

	fmt.Println("CreateGroupEvent end")

	return int(id), nil
}

// Handler for "/group/" endpoint.
//
// Method: POST - Creates a new group event.
//
// Method: GET - Fetches events associated with the group_id. group_id needs to be specified in the request URL.
func CreateGroupEventHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("events.go CreateGroupEventHandler: session cookie err: ", err)
		msg, _ := json.Marshal("events.go CreateGroupEventHandler: session cookie missing")
		w.Write(msg)
		return
	}

	userId := GetUserByCookie(cookie.Value)
	if userId == 0 {
		err := fmt.Errorf("incorrect session cookie value")
		fmt.Println("events.go CreateGroupEventHandler: ", err)
		msg, _ := json.Marshal(err.Error())
		w.Write(msg)
		return
	}

	switch r.Method {
	case "GET":
		groupId := path.Base(r.URL.Path)
		_, err := strconv.Atoi(groupId)
		if err != nil {
			err = fmt.Errorf("invalid groupId: %v", groupId)
			fmt.Println(err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		groupMembers, err := GetGroupMembersFromDb(groupId)
		if err != nil {
			err = fmt.Errorf("groupId not found in db: %v", groupId)
			fmt.Println(err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		isUserPartOfGroup := false
		for _, v := range groupMembers {
			if v.Id == userId {
				isUserPartOfGroup = true
				break
			}
		}

		if !isUserPartOfGroup {
			err := fmt.Errorf("user %v is not a member of group %v", userId, groupId)
			fmt.Println(err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		events, err := GetEventFromDb("group_id", groupId)
		if err != nil {
			fmt.Println("events.go CreateGroupEventHandler err 4: ", err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		for i, v := range events {
			attendingUsers, err := GetAttendingUsers(v.EventID)
			if err != nil {
				fmt.Println("events.go CreateGroupEventHandler err 5: ", err)
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			events[i].AttendingUsers = append(events[i].AttendingUsers, attendingUsers...)
		}

		json.NewEncoder(w).Encode(events)
		return
	case "POST":
		var event customstructs.GroupEvent
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			fmt.Println("events.go CreateGroupEventHandler err 6: ", err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		group, err := GetGroupFromDb("group_id", event.Group.(string))
		if err != nil {
			fmt.Println("events.go GetEventFromDb err 7: ", err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		event.Group = group[0]

		eventID, err := CreateGroupEvent(event)
		if err != nil {
			fmt.Println("Unable to create event: ", err)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		jsonResponse := struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			EventID int    `json:"event_id"`
		}{
			Status:  "success",
			Message: "Event created successfully",
			EventID: eventID,
		}

		event.EventID = fmt.Sprint(eventID)

		go EventTimeout(event)
		handleGroupEventNotifications(event)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResponse)
	}
}

// Records a user's response (Attending/Not attending/Pending) to a specific event in the database.
// Returns error on failure, <nil> on success
func RespondToEvent(response customstructs.EventResponse) error {
	fmt.Println("RespondToEvent start")
	if response.Response != "Attending" && response.Response != "Not attending" /* && response.Response != "Pending" */ {
		return fmt.Errorf("invalid response value")
	}

	fmt.Println("RespondToEvent response: ", response)
	_, err := Db.Exec("UPDATE group_event_responses SET response = ? WHERE (event_id, user_id) = (?, ?);", response.Response, response.EventID, response.User)
	if err != nil {
		return fmt.Errorf("failed to respond to event: %v", err)
	}

	fmt.Println("RespondToEvent end")
	return nil
}

// Gets the user's response to an event, eventId must be specified in the request URL.
func GetUserEventStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := fmt.Errorf("only GET requests are allowed")
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("events.go GetUserEventStatusHandler: session cookie err: ", err)
		json.NewEncoder(w).Encode("session cookie missing")
		return
	}

	userId := GetUserByCookie(cookie.Value)
	if userId == 0 {
		err := fmt.Errorf("incorrect session cookie value")
		fmt.Println("events.go GetUserEventStatusHandler: ", err)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	eventId := path.Base(r.URL.Path)
	intEventId, err := strconv.Atoi(eventId)
	if err != nil {
		fmt.Println("events.go GetUserEventStatusHandler err 1: ", err)
		err = fmt.Errorf("invalid event ID")
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	var returnData []customstructs.EventResponse
	isUserPartOfGroup := false

	rows, err := Db.Query("SELECT * FROM group_event_responses WHERE event_id = ?", intEventId)
	if err != nil {
		fmt.Println("events.go GetUserEventStatusHandler err 2: ", err)
		json.NewEncoder(w).Encode("user is not signed up to the event")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r customstructs.EventResponse
		var id string
		rows.Scan(&r.ResponseId, &r.EventID, &id, &r.Response)
		r.User = id
		returnData = append(returnData, r)
		if r.User.(string) == fmt.Sprint(userId) {
			isUserPartOfGroup = true
		}
	}

	if !isUserPartOfGroup {
		json.NewEncoder(w).Encode("user does not have access to this event")
	}

	for i, v := range returnData {
		user, err := GetUser("users.user_id", v.User.(string))
		if err != nil {
			fmt.Println("events.go GetUserEventStatusHandler err 3: ", err)
			json.NewEncoder(w).Encode("unexpected error")
			return
		}
		user = ProcessUserInfo(user, -1)
		returnData[i].User = user[0]
	}

	json.NewEncoder(w).Encode(returnData)
}

// Retrieves event data from the database based on the condition string.
func GetEventFromDb(condition, value string) ([]customstructs.GroupEvent, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if condition == "" {
		rows, err = Db.Query("SELECT * FROM group_events;")
	} else {
		rows, err = Db.Query("SELECT * FROM group_events WHERE "+condition+"=$1;", value)
	}

	if err != nil {
		fmt.Println("events.go GetEventFromDb err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	var returnData []customstructs.GroupEvent

	for rows.Next() {
		var t customstructs.GroupEvent
		var gId string
		err = rows.Scan(&t.EventID, &gId, &t.Title, &t.Description, &t.DayTime)
		if err != nil {
			fmt.Println("events.go GetEventFromDb err 2: ", err)
			return nil, err
		}

		group, err := GetGroupFromDb("group_id", gId)
		if err != nil {
			fmt.Println("events.go GetEventFromDb err 3: ", err)
			return nil, err
		}
		t.Group = group[0]

		returnData = append(returnData, t)
	}

	return returnData, nil
}

// Retrieves the attending users for an event from the database.
func GetAttendingUsers(eventId string) ([]customstructs.UserData, error) {
	var (
		rows *sql.Rows
		err  error
	)

	query := `
	SELECT 
		user_id
	FROM 
		group_event_responses
	WHERE
		event_id = ?
		AND
		response="Attending"
	`

	rows, err = Db.Query(query, eventId)
	if err != nil {
		fmt.Println("events.go GetAttendingUsers err 1: ", err)
		return nil, err
	}
	defer rows.Close()

	var returnData []customstructs.UserData

	var userIdArr []int

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			fmt.Println("events.go GetEventFromDb err 2: ", err)
			return nil, err
		}

		userIdArr = append(userIdArr, i)
	}

	for _, v := range userIdArr {
		user, err := GetUser("users.user_id", fmt.Sprint(v))
		if err != nil {
			fmt.Println("events.go GetEventFromDb err 3: ", err)
			return nil, err
		}

		ProcessUserInfo(user, -1)

		returnData = append(returnData, user...)
	}

	return returnData, nil
}

// Removes the event(s) from the db based on the condition string.
func DeleteEvent(condition, value string) int {
	res, err := Db.Exec("DELETE FROM group_events WHERE "+condition+" = $1;", value)
	if err != nil {
		fmt.Println("DeleteEvent err 1: ", err)
		return 0
	}

	i, err := res.RowsAffected()
	if err != nil {
		fmt.Println("DeleteEvent err 2: ", err)
		return 0
	}

	fmt.Printf("DeleteEvent: %v events deleted\n", i)

	return int(i)
}

// Calls DeleteEvent when the event has expired.
func EventTimeout(event customstructs.GroupEvent) {
	sleepTime := time.Until(event.DayTime)
	time.Sleep(sleepTime)
	if DeleteEvent("event_id", event.EventID) != 0 {
		handleGroupEventNotifications(event)
	}
}
