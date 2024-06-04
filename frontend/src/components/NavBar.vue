<script setup>
import { RouterLink } from "vue-router";
import axios from "axios";
import { closeWebSocket } from "../webSocket/webSocket";
import { sendWebSocketMessage } from "../webSocket/webSocket";
import { mapState } from "vuex";
import { selfInfo } from "../main.js";

</script>
<script>
export default {
  data() {
    return {
      selfData: {},
      isDropdownOpen: false,
      
      showMumm: false,
    };
  },
  created(){
    selfInfo().then((result) => {
      this.selfData = result
    });
  },
  computed: {
    ...mapState(["allNotifications"]),
  },
  props: {
    // Define a prop to receive the information
    showChat: Boolean, // Adjust the type as needed
  },
  methods: {
    acceptNotification(notificationId, desc, senderId, targetId, groupId){
      let notif = {
        notificationId,
        sender: senderId,
        target: targetId,
        desc: desc,
        seenByTarget: "1",
        seenBySender: "0",
        status: "accepted",
        value: groupId,
      };

      sendWebSocketMessage(notif, "notification");
      //removeNotification(notif, "notification")
      //this.showMummMethod(false)
      console.log("Invite accepted!");
    },

    declineNotification(notificationId, desc, senderId, targetId, groupId){
      let notif = {
        notificationId,
        sender: senderId,
        target: targetId,
        desc: desc,
        seenByTarget: "1",
        seenBySender: "0",
        status: "rejected",
        value: groupId,
      };

      sendWebSocketMessage(notif, "notification");
      //removeNotification(notif, "notification")
      //this.showMummMethod(false)
      console.log("Invite rejected");
    },
    //for dismissing the notification that user's request has been accepted - Karl
    dismissNotification(notificationId, desc, senderId, targetId, groupId){
      let notif = {
        notificationId,
        sender: senderId,
        target: targetId,
        desc: desc,
        seenByTarget: "1",
        seenBySender: "1",
        status: "accepted",
        value: groupId,
      };

      sendWebSocketMessage(notif, "notification");
      //removeNotification(notif, "notification")
      //this.showMummMethod(false)
      console.log("Notification dismissed");
    },
    toggleDropdown() {
      this.isDropdownOpen = !this.isDropdownOpen;
    },
    toggleChat() {
      console.log("chatbuttonpressed");
      this.$emit("toggleChat");
    },
    profileClicked() {
      selfInfo().then(data => {
        let userId = data.Id;
        this.$router.push({ name: "profile", params: { userId } });
      });
    },
    logout() {
      fetch("http://localhost:8000/logout/", {
        credentials: "include",
      })
        .then((response) => {
          return response.json();
        })
        .then(async (data) => {
          if (data === "You have been logged out") {
            console.log("logout successful!");
            closeWebSocket();
            this.$router.push({ name: "login" });
          } else {
            console.log("Server response from logout.go: ", data);
            this.$router.push({ name: "login" });
          }
        })
        .catch((error) => {
          console.error("Failed to fetch:", error);
        });
    },
    goToGroup(groupId){
      console.log("GroupId", groupId)
      this.$router.push({ name: 'SingleGroup', params: { groupId } });  
    },

    showMummMethod(value) {
      //console.log("showMummMethod: ", value);
      this.showMumm = value;
    },
  },
  
  watch:{
    allNotifications (data) {
      //console.log("watcher data:", data);
      this.showMummMethod(false);
    }
  },
};
</script>
<template>
  <div class="navbar">
    <meta name="viewport" content="width=device-width" />
    <link
      rel="stylesheet"
      href="https://cdn.lineicons.com/1.0.1/LineIcons.min.css"
    />
    <div class="colum">
      <div class="chat">
        <button
          id="chat"
          class="lni lni-bubble size-md"
          @click="toggleChat"
        ></button>
      </div>
      <div class="bell">
        <button
          id="bell"
          class="lni lni-alarm size-md"
          @click="toggleDropdown"
        > 
        <span v-if="showMumm" class="mumm"></span>
        </button>
        <div v-show="isDropdownOpen" class="dropdown-content">
        <!-- <div v-if="isDropdownOpen" class="dropdown-content"> -->
          <!-- <div v-if="allNotifications" v-for="notification in allNotifications" :key="notification.notificationId"> -->
          <div v-for="notification in allNotifications" :key="notification.notificationId">

            <template v-if="notification.desc === 'followRequest' && notification.target.Id === selfData.Id && notification.seenByTarget === '0'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.sender.FirstName }} sent you a follow request.</div>
                <button
                  class="accept"
                  @click="acceptNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id)"
                >
                Accept
                </button>
                <button
                  class="decline"
                  @click="declineNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id)"
                >
                Reject
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'followRequest' && notification.target.Id != selfData.Id && notification.seenBySender === '0' && notification.seenByTarget === '1' && notification.status === 'accepted'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.target.Nickname.String || notification.target.FirstName + ' ' + notification.target.LastName }} has accepted you as a follower.</div>
                <button
                  class="dismiss"
                  @click="dismissNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                  >
                  Dismiss
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'groupInvite' && notification.target.Id === selfData.Id && notification.seenByTarget === '0'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.sender.Nickname.String || notification.sender.FirstName + ' ' + notification.sender.LastName }} Has invited you to join their group: {{ notification.value.GroupTitle }}</div>
                <button
                  class="accept"
                  @click="acceptNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                >
                  Accept
                </button>
                <button
                  class="decline"
                  @click="declineNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                >
                  Reject
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'groupInvite' && notification.target.Id != selfData.Id && notification.seenBySender === '0' && notification.status === 'accepted'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.target.Nickname.String || notification.target.FirstName + ' ' + notification.target.LastName }} Has accepted your invite to join the group: {{ notification.value.GroupTitle }}</div>
                <button
                  class="dismiss"
                  @click="dismissNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                >
                  Dismiss
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'groupJoinRequest' && notification.sender.Id !== selfData.Id && notification.seenByTarget === '0'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.sender.Nickname.String || notification.sender.FirstName + ' ' + notification.sender.LastName }} Has requested to join the group "{{ notification.value.GroupTitle }}"</div>
                <button
                  class="accept"
                  @click="acceptNotification(notification.notificationId, notification.desc, notification.sender.Id, selfData.Id, notification.value.GroupId)"
                >
                  Accept
                </button>
                <button
                  class="decline"
                  @click="declineNotification(notification.notificationId, notification.desc, notification.sender.Id, selfData.Id, notification.value.GroupId)"
                >
                  Reject
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'groupJoinRequest' && notification.sender.Id === selfData.Id && notification.seenByTarget === '1' && notification.status === 'accepted'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>{{ notification.value.Creator[0].Nickname.String || notification.value.Creator[0].FirstName + ' ' + notification.value.Creator[0].LastName }} has accepted your request to join the group "{{ notification.value.GroupTitle }}"</div>
                <button
                  class="dismiss"
                  @click="dismissNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                >
                  Dismiss
                </button>
              </div>
            </template>

            <template v-if="notification.desc === 'groupDeleted' && notification.seenByTarget === '0'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>Group {{ notification.value }} has been deleted.</div>
                <button
                  class="dismiss"
                  @click="dismissNotification(notification.notificationId, notification.desc, notification.sender, notification.target.Id, notification.value)"
                >
                  Dismiss
                </button>
              </div>
            </template>

            <template v-if="notification.desc == 'event' && notification.target.Id === selfData.Id && notification.seenByTarget === '0'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>You have been invited to join an event: {{ notification.value.title }}</div>
                <button
                  class="accept"
                  @click="acceptNotification(notification.notificationId, notification.desc, notification.sender, notification.target.Id, notification.value.event_id)"
                >
                  Accept
                </button>
                <button
                  class="middle"
                  @click="declineNotification(notification.notificationId, notification.desc, notification.sender, notification.target.Id, notification.value.event_id)"
                >
                  Reject
                </button>
                <button
                  class="decline"
                  @click="goToGroup(notification.value.group.GroupId)"
                >
                  Go to group
                </button>
              </div>
            </template>

            <template v-if="notification.desc == 'event' && notification.target.Id === selfData.Id && notification.seenByTarget === '1' && notification.status === 'accepted'">
              {{ showMummMethod(true) }}
              <div class="notification">
                <div>You have successfully been registered for the event: {{ notification.value.title }}!</div>
                <button
                  class="dismiss"
                  @click="dismissNotification(notification.notificationId, notification.desc, notification.sender.Id, notification.target.Id, notification.value.GroupId)"
                >
                Dismiss
                </button>
              </div>
            </template>

          </div>
          <p v-show="!showMumm">No notifications sorry. :D</p>
        </div>
      </div>
      <div class="home">
        <RouterLink 
          to="/feed" 
          class="lni-home size-md" 
          id="home">
        
        </RouterLink>
      </div>

      <div class="post">
        <RouterLink
          to="/createPost"
          class="lni-write size-md"
          id="post">

        </RouterLink>
      </div>
      <div class="group">
        <RouterLink
          to="/Groups"
          id="group"
          class="lni lni-users size-md"
        ></RouterLink>
      </div>
      <div class="user">
        <span @click="profileClicked" class="lni-user size-md" id="user"></span>
      </div>
      <div class="logout">
        <button
          @click="logout"
          to="/"
          class="lni-exit size-md"
          id="logout"
        ></button>
      </div>
    </div>
  </div>
</template>
<style>
.mumm{
  border: 5px solid #00bd7e;
  border-radius: 5px;
}
.dropdown-content {
  background-color: #252525;
  padding: 1rem;
  border-radius: 1rem;
  width: max-content;
  position: relative;
  left: -3rem;
}
.navbar {
  width: 98%;
  height: 60px;
  background-color: #101010;
  left: 1%;
  top: 1%;
  position: fixed;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  /* box-shadow: 0.1px 0.1px 15px 0.1px #273c75; */
  border-radius: 1rem;
  display: flex;
}

#chat, #bell, #logout, #user, #group {
  margin-left: 0;
  margin-top: 2px;
  color: #00bd7e;
  transition: width 0.5s, background-color 0.5s, border-radius 0.5s;
  cursor: pointer;
  background-color: #101010;
  border: 0;
  transition: 0.4s;
}

#group:hover, #user:hover, #logout:hover, #bell:hover, #home:hover, #post:hover, #chat:hover{
  color: #101010;
  background-color: #00bd7e;
  border-radius: 0.5rem;
}
.colum {
  display: flex;
  width: 100%;
  height: 40px;
  margin-top: 10px;
  margin-left: 10px;
}

.post, .home {
  margin-left: 50%;
  position: fixed;
  padding-top: 2px;
}

.home {
  right: 50%;
}

.chat {
  width: 80px;
  height: 34px;
  margin-left: 20px;
  margin-right: 1rem;
  display: flex;
}

.bell {
  width: 80px;
  height: 34px;
  margin-left: 2px;
  display: flex;
  align-items: center;
  flex-direction: column;
}
.logout {
  right: 0%;
  position: fixed;
  width: 80px;
  height: 34px;
  margin-left: 20px;
  display: flex;
}
.user {
  right: 5%;
  position: fixed;
  width: 80px;
  height: 34px;
  margin-left: 20px;
  display: flex;
}
.group {
  right: 10%;
  position: fixed;
  width: 80px;
  height: 34px;
  margin-left: 20px;
  display: flex;
}
.clicked {
  width: 100px;
  height: 34px;
  background-color: #feca57;
  border-radius: 10px;
  transition: background-color 0.3s, border-radius 0.5s, width 0.5s;
}


</style>
