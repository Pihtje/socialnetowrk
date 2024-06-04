<script setup>
import axios from "axios";
import NavBar from "../components/NavBar.vue";
import Post from "../components/Post.vue";
import Chat from "../components/Chat.vue";
import CreateGroupEventView from "./CreateGroupEventView.vue";
import { sendWebSocketMessage } from "../webSocket/webSocket";
import { selfInfo } from "../main.js";
import { mapState } from "vuex";
import EmojiPicker from "../components/EmojiPicker.vue";
</script>
<script>
export default {
  components: {
    Post,
    NavBar,
    Chat,
    CreateGroupEventView,
  },
  computed: {
    groupId() {
      return this.$route.params.groupId;
    },

    userResponseStatusForEvent() {
      return (eventId) => {
        const responseObj = this.eventResponses.find(response => response.event_id == eventId && response.user.Id == this.selfId); // Updated the userId reference
        return responseObj ? responseObj.response : null;
      }
    },

    ...mapState(["allMessages", "allNotifications"]),
    formattedMessageTime() {
      return function (message) {
        const date = new Date(message.DateTime);
        const day = date.getDate().toString().padStart(2, "0");
        const month = (date.getMonth() + 1).toString().padStart(2, "0");
        const year = date.getFullYear();
        const hours = date.getHours().toString().padStart(2, "0");
        const minutes = date.getMinutes().toString().padStart(2, "0");
        return `${day}.${month}.${year} ${hours}:${minutes}`;
      };
    },
  },
  data() {
    return {
      selfData: {},
      showModal: false,
      groupInfo: {},
      groupPosts: [],
      showChat: false,
      socket: Object,
      groupEvents: [],
      eventResponses: [],
      userIsMember: true,
      showEmojiPicker: false,
    };
  },
  async created() {
    selfInfo().then((result) => {
      this.selfData = result;
    });
    try {
      const response = await axios.get(
        `http://localhost:8000/group/${this.groupId}`,
        {
          withCredentials: true,
        }
      );
      this.groupInfo = response.data[0];
      if (response.data == "user does not have access to group ") {
        this.userIsMember = false;
      } 
      if (response.data == `group with the id of ${this.groupId} does not exist`) {
        this.userIsMember = false;
      }
      console.log(this.groupInfo);
      this.groupPosts = this.groupInfo.Posts;
    } catch (error) {
      console.error("Error loading userInfo:", error);
    }

    //eventide kuvamiseks, mson endpointid
    try {
      const eventResponse = await axios.get(
        `http://localhost:8000/createGroupEvent/${this.groupId}`,
        {
          withCredentials: true,
        }
      );
      this.groupEvents = eventResponse.data;

      console.log("SingleGroupView.vue this.groupEvents: ", this.groupEvents);
    } catch (error) {
      console.error("Error loading group events:", error);
    }

    try {
      this.groupEvents.forEach(async (e) => {
        let res = await axios.get(
          `http://localhost:8000/getUserEventStatus/${e.event_id}`,
          {
            withCredentials: true,
          }
        );
        this.eventResponses = res.data;
      });
    } catch (error) {
    }
  },
  methods: {
    openEmojiPicker() {
      this.showEmojiPicker = !this.showEmojiPicker;
    },
    insertEmoji(emoji) {
      // Insert the selected emoji into the input field
      this.newMessage += emoji;
    },
    closeEmojiPicker() {
      this.showEmojiPicker = false;
    },
    sendGroupDM() {
      console.log("sending DM...");
      const timestamp = new Date().toISOString();

      console.log(
        this.selfId,
        this.newMessage,
        timestamp
      );

      if (this.selfId /* && this.ChatTarget.Id */ && this.newMessage) {
        const message = {
          sender: this.selfId + "",
          target: "group",
          msg: this.newMessage,
          serverObject: true,
          groupId: this.groupInfo.GroupId + "",
          DateTime: timestamp,
        };

        // Use the sendWebSocketMessage function to send the message through WebSocket
        sendWebSocketMessage(message, "message");
        this.newMessage = "";
      }
    },
    inviteMember(value) {
      let notif = {
        sender: this.selfId + "",
        target: value,
        desc: "groupInvite",
        seenByTarget: "0",
        seenBySender: "0",
        status: "pending",
        value: this.groupInfo.GroupId,
      };

      sendWebSocketMessage(notif, "notification");
      console.log("Person invited!");
    },
    requestJoinGroup() {
      let notif = {
        sender: this.selfId + "",
        target: "group",
        desc: "groupJoinRequest",
        seenByTarget: "0",
        seenBySender: "0",
        status: "pending",
        value: this.groupId,
      };

      sendWebSocketMessage(notif, "notification");
      console.log("Request sent!");
    },
    deleteGroup() {
      let notif = {
        sender: this.selfId + "",
        target: "server",
        desc: "deleteGroup",
        seenByTarget: "0",
        seenBySender: "1",
        status: "pending",
        value: this.groupInfo.GroupId + "",
      };

      sendWebSocketMessage(notif, "notification");
      console.log("group deleted");
      this.$router.push({ name: "Groups" });
    },
    leaveGroup() {
      let notif = {
        sender: this.selfId + "",
        target: "group",
        desc: "leaveGroup",
        seenByTarget: "0",
        seenBySender: "1",
        status: "pending",
        value: this.groupInfo.GroupId + "",
      };

      sendWebSocketMessage(notif, "notification");
      console.log("group left");
      this.$router.push({ name: "Groups" });
    },

    toggleChat() {
      console.log("toggling chat");
      // this.showChat = !this.showChat;
      if (this.showChat == false) {
        this.showChat = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showChat = false;
      }
    },
    async handleNewPostSubmit(event) {
      event.preventDefault();
      // console.log(event)
      const formData = new FormData();
      formData.append("Title", event.target.postTitle.value);
      formData.append("Body", event.target.postBody.value);
      formData.append("Category", "1");
      formData.append("ShareWith", " ");
      formData.append("ImageUrl", this.postImage)
      formData.append("GroupId", this.groupId);
      console.log("Form data:", formData);
      const files = event.target.postImages.files;
      for (let i = 0; i < files.length; i++) {
        formData.append("postImages", files[i]);
      }

      try {
        const response = await axios.post(
          "http://localhost:8000/newPost/",
          formData,
          {
            withCredentials: true,
          }
        );

        const data = response.data;
        console.log("handleNewPostSubmit: ", data);
        location.reload();
      } catch (error) {
        console.error("Failed to fetch:", error);
      }
    },
    OpenMembers() {
      if (this.showModal == false) {
        this.showModal = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showModal = false;
      }
    },
    
    respondToEvent(eventId, responseStatus) {
      let notific = this.allNotifications.find((n) => n.value.event_id == eventId);

      let notif = {
        notificationId: notific.notificationId,
        sender: notific.sender,
        target: notific.target.Id,
        desc: "event",
        seenByTarget: "1",
        seenBySender: "0",
        status: responseStatus,
        value: eventId,
      };

      sendWebSocketMessage(notif, "notification");
      //removeNotification(notif, "notification")
      console.log("Invite accepted!");
    },
  },
};
</script>
<template>
  <div class="singleGroupView" id="app">
    <div v-if="!this.userIsMember" class="groupContainer">
      <h1>You are sadly not a member of this group...</h1>
      <button class="formButton" @click="requestJoinGroup">Join up!</button>
    </div>
    <div v-else class="group-groupChat">
      <div class="groupContainer">
        
        <div class="groupInfoContainer">
          <div class="group-info">
            <h1 class="groupTitle">{{ this.groupInfo.GroupTitle }}</h1>
            <br />
            <h2 class="group-creator">
              <span v-if="this.groupInfo.Creator"
                >Group creator: {{ this.groupInfo.Creator[0].FirstName }}</span
              >
            </h2>
            <br />
            <p>Group Description:</p>
            <h3 class="groupDesc">{{ this.groupInfo.GroupDesc }}</h3>
            <br />
            <div
              v-if="
                this.groupInfo.Creator &&
                this.selfData.Id == this.groupInfo.Creator[0].Id
              "
            >
              <button class="formButton" @click="deleteGroup">Delete Group</button>
            </div>
            <div v-else><button class="formButton" @click="leaveGroup">Leave Group</button></div>
          </div>
          <div @click="OpenMembers" class="membersContainer">
            <div class="members">
              <span>Members</span>
              <a class="memberAmount">{{
                this.groupInfo.Members?.length || "1"
              }}</a>
              <div v-if="showModal" class="membersList">
                <div v-for="member in this.groupInfo.Members" :key="member.Id">
                  <span>{{member.Nickname?.String || member.FirstName + ' ' + member.LastName}}</span>
                </div>
              </div>
            </div>
          </div>
          <div>
            <h2>Invite Members</h2>

            <select class="formField" v-model="selected" v-if="selfData">
              <option
                v-for="person in this.selfData.Followers"
                :key="person.Id"
                :value="person.Id"
              >
                {{ person.Nickname.String || person.FirstName + ' ' + person.LastName }}
              </option>
            </select>
            <button style="margin-left: 1rem;" class="formButton" @click="inviteMember(selected)">Invite Member</button>
          </div>
          <div class="createPostContainer">
            <form class="createPost" @submit="handleNewPostSubmit">
              <h1>Create a new post!</h1>
              <label for="postTitle">Title:</label>
              <br />
              <br />
              <input class="postForm" type="text" name="postTitle" required />
              <br />
              <label for="postBody">Text body</label>
              <br />
              <textarea class="postForm" name="postBody" required></textarea>
              <br />
              <br />
              <input
                @change="onFileSelected"
                class="postForm"
                type="file"
                id="postImages"
                multiple
                accept=".jpg,.jpeg,.png,.gif"
                name="postImages"
              />
              <br />
              <button class="postForm" type="submit">Create new post!</button>
            </form>
            <div class="groupEventsContainer">
              <CreateGroupEventView :groupId="groupId" />
              <h2>Group Events</h2>
              <ul>
                <li v-for="event in groupEvents" :key="event.event_id">
                  <h3>Event title: {{ event.title }}</h3>
                  <p>Event description: {{ event.description }}</p>
                  <p>Date & Time: {{ event.day_time }}</p>
                  <div :key="event.event_id">
                    <p>Attending users:</p>
                    <div v-if="eventResponses">
                      <ul>
                        <li v-for="response in eventResponses" :key="response.response_id">
                          {{ response.user.Nickname.String || response.user.FirstName + ' ' + response.user.LastName }} - {{ response.response }}
                        </li>
                      </ul>
                    </div>
                    <div v-else>
                      <p>Loading responses...</p>
                    </div>
                    <div v-if="userResponseStatusForEvent(event.event_id) === 'Pending'">
                      <button class="formButton" style="margin: 5px;" @click="respondToEvent(event.event_id, 'accepted')">Attending</button>
                      <button class="formButton" style="margin: 5px;" @click="respondToEvent(event.event_id, 'rejected')">Not attending</button>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <br />
        </div>
      <div class="chatBox">
        <div class="chat_header">
          <div class="imgcontent">
            <link
              rel="stylesheet"
              href="https://cdn.lineicons.com/1.0.1/LineIcons.min.css"
            />
            <h3>GroupChat</h3>
          </div>
        </div>
        <div class="messagesContainer">
          <div
            class="messageBox"
            v-for="message in this.allMessages"
            :key="message.messageId"
          >
            <div
              class="message"
               :class="message.sender.Id != this.selfId ? 'incoming' : 'sent'"

              v-if="message.GroupId == this.groupId"
            >
            <!-- <div class="message"> -->
              {{ message.msg }}<br /><span class="time">{{
                formattedMessageTime(message)
              }}</span>
            </div>
          </div>
        </div>
        <div class="messageInput">
          <input
            v-model="newMessage"
            type="text"
            name="msg"
            id="msg"
            cols="30"
            rows="10"
            maxlength="255"
            placeholder="Enter message content, 1 to 255 characters"
            autocomplete="off"
            @keyup.enter="sendGroupDM" 
          />
          <button class="formButton" @click="openEmojiPicker">😀</button>

<!-- Display the emoji picker when it's open -->
<emoji-picker
  v-if="showEmojiPicker"
  @emoji-selected="insertEmoji"
  @close-picker="closeEmojiPicker"
/>
          <button class="formButton" @click="sendGroupDM">Send</button>
        </div>
      </div>
    </div>

    <Post
      v-for="post in groupPosts"
      :key="post.Id"
      :Username="post.Username"
      :PostTime="post.PostTime"
      :Avatar="post.Avatar"
      :Title="post.Title"
      :Body="post.Body"
      :ImageUrl="post.ImageUrl"
      :Comments="post.Comments"
      :Likes="post.Likes"
      :Id="post.Id"
      :Creator="post.Creator"
      @author-clicked="goToProfile"
    />
    <NavBar :showChat="showChat" @toggleChat="toggleChat" />
    <Chat :showChat="showChat" @toggleChat="toggleChat" />
  </div>
</template>
<style scoped>
.singleGroupView {
  flex-direction: column;
}
.group-groupChat {
  width: 100%;
  display: flex;
  margin-bottom: 5rem;
  flex-direction: row;
}
.groupContainer {
  position: relative;
  top: 10%;
  background-color: #101010;
  width: 80%;
  border-radius: 1rem;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  padding: 2rem;
  margin: 2.5rem;
}
.createPostContainer {
  display: flex;
}
.membersContainer {
  display: flex;
  flex-direction: row;
  box-shadow: 0.2rem 0.2rem #1a1a1a4b;
}
.members span {
  color: #00bd7e;
  background-color: #006644;
  border: 0;
  margin-left: 5rem;
  padding: 0.5rem;
  border-top-left-radius: 0.5rem;
  border-bottom-left-radius: 0.5rem;
  cursor: pointer;
}
.memberAmount {
  color: #1b1b1b;
  background-color: #00bd7e;
  border: 0;
  padding: 0.6rem;
  border-top-right-radius: 0.5rem;
  border-bottom-right-radius: 0.5rem;
  cursor: pointer;
}

/* Create Event Styles */
.createEventContainer {
  display: flex;
  flex-direction: column ;
  background-color: #101010;
  padding: 1rem;
  border-radius: 1rem;
  /* width: 30rem; */
  /* box-shadow: 0rem 0rem 0.5rem rgba(0, 0, 0, 0.2); */
  margin-top: 1rem;
}
.postForm {
  margin-left: 1rem;
}
.createEventContainer label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.createEventContainer input[type="text"],
.createEventContainer textarea,
.createEventContainer select {
  width: 100%;
  padding: 0.5rem;
  border-radius: 0.5rem;
  border: 1px solid #ddd;
  font-size: 1rem;
  margin-bottom: 1rem;
}



.createEventContainer button:hover {
  background-color: #0056b3;
}

.close-btn {
  background-color: #080808;
  border: 0;
  margin-left: 0;
  margin-right: 0;
  margin-top: 2px;
  color: #00bd7e;
  transition: width 0.5s, background-color 0.5s, border-radius 0.5s;
}

.chatBox {
    position: fixed;
    border-radius: 1rem;
    right: 0;
    bottom: 0%;
    width: 21%;
    height: 80%;
    background: rgb(0, 0, 0);
    z-index: 10000;
}

.chat_header {
  border-bottom: 0.1rem;
  border-bottom-color: #00bd7e;
  border-bottom-style: solid;
  position: relative;
  width: 100%;
  padding: 10px;
  background: #080808;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top-left-radius: 1rem;
  border-top-right-radius: 1rem;
}

.messagesContainer {
    overflow-x: hidden;
    overflow-y: auto;
    position: relative;
    padding: 20px;
    height: 80%;
}

.message {
  margin-top: 1rem;
  position: relative;
  padding: 10px;
  background: #00bd7e;
  border-radius: 10px;
}
.sent::before {
  content: "";
  position: absolute;

  top: 0;
  right: -10px;
  border: 10px solid transparent;
  border-top: 10px solid #00bd7e;
}
.incoming::before {
  content: "";
  position: absolute;

  top: 0;
  left: -10px;
  border: 10px solid transparent;
  border-top: 10px solid #00bd7e;
}

.time {
  position: relative;
  display: block;
  font-size: 0.7em;
  width: 100%;
  text-align: end;
}

.formButton {
  width: fit-content;
}
.messageInput {
  background-color: #080808;
  width: 100%;
  height: 5rem;
  padding: 5px 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #00bd7e;
}

input {
  position: relative;
  background: #0a0a0a;
  padding: 5px 10px;
  width: 15rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #080808;
  color: #00bd7e;
  border: 1px solid #00bd7e;
  border-radius: 0.5rem;
}
</style>
