<script setup>
import ChatBox from "../components/ChatBox.vue";
import { getWebSocket, setupWebSocket } from "../webSocket/webSocket";
import { mapState, mapMutations, mapActions } from "vuex";
</script>
<script>

export default {
  data() {
    return {
      chatBoxes: [],
      allNotifications: [],
      socket: null,
    };
  },
  props: {
    showChat: Boolean,
  },
  async created() {
    const socket = getWebSocket();

    if (!socket || socket.readyState !== WebSocket.OPEN) {
      // WebSocket is not initialized or not open; initialize it.
      setupWebSocket();
    }
  },
  computed: {

    ...mapState(["allUsers", "allMessages"]),

    visibleChatBoxes() {
      return this.chatBoxes.filter((box) => box.visible);
    },
    sortedUsers() {
      return this.allUsers.slice().sort((a, b) => {
        const latestMessageA = this.getLatestMessageTimestamp(a.Id);
        const latestMessageB = this.getLatestMessageTimestamp(b.Id);
        return latestMessageB - latestMessageA;
      });
    },

  },

  methods: {
    getLatestMessageTimestamp(userId) {
  if (!this.allMessages) {
    return 0; // Return 0 if allMessages is null
  }

  const userMessages = this.allMessages.filter((message) => {
   
    return message.sender.Id == userId || message.target.Id == userId;
  });

  if (userMessages.length === 0) {
    return 0;
  }

  const latestTimestamp = Math.max(
    ...userMessages.map((message) => new Date(message.DateTime).getTime())
  );

  return latestTimestamp;
},
    getLatestMessage(userId) {
      if (!this.allMessages) {
        return 'No messages'; // Return a default value if allMessages is null
      }

      const userMessages = this.allMessages.filter(
        (message) =>
          message.sender.Id == userId || message.target.Id == userId
      );

      if (userMessages.length === 0) {
        return 'No messages';
      }

      // You can implement how you want to display the latest message here
      const latestMessage = userMessages.reduce((prev, current) =>
        prev.DateTime > current.DateTime ? prev : current
      );

      return latestMessage.msg; // Adjust this based on your message data structure
    },
    toggleChat() {
   
      this.$emit("toggleChat");
     
    },
    openChatBox(userInfo) {
      console.log("All users ", this.allUsers);
      console.log("All messages", this.allMessages);
   
      console.log("Opening chatbox", this.chatBoxes);
      this.chatBoxes.forEach((chatBox) => {
        if (chatBox.ChatTarget.FirstName === userInfo.FirstName) {
          this.chatBoxes = this.chatBoxes.filter(
            (item) => item.title !== userInfo.FirstName
          );
        }
      });
      this.chatBoxes.push({
        ChatTarget: userInfo,
        visible: true,
      });
    },
    closeChatBox(index) {
      this.chatBoxes.splice(index, 1);
    },
  },
};
</script>
<template>
  <div class="chat">
    <div v-if="showChat" class="container" id="chat-form-container">
      <header>
        <a href="#" class="logo">Chat</a>
        <button
          class="lni lni-close size-md close-btn"
          @click="toggleChat"
        ></button>
      </header>

      <div class="content">
        <div class="chatlist">
          <div
          v-for="chat in sortedUsers"
            :key="chat.Id"
            class="block"
            :class="{ unread: chat.Unread }"
            @click="openChatBox(chat)"
          >
            <div image class="imgbx">
              <img
                :src="BackendImg + chat.ImageUrl.String"
                @error="
                  $event.target.src =
                    'https://wallpapers-clan.com/wp-content/uploads/2022/08/default-pfp-1.jpg'
                "
              />
            </div>
            <div class="details">
              <div class="listHead">
                <h4>{{ chat.FirstName }} {{ chat.LastName }}</h4>
                <!-- <p class="time">{{ chat.lastSentMessage }}a</p> -->
              </div>
              <div class="message_p">
              
                <div>
                  Latest Message: {{ getLatestMessage(chat.Id) }}
      </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
 
      <ChatBox
      v-for="(chatBox, index) in chatBoxes"
      :key="index"
      class="chat-boxes"
      :class="{ visible: chatBox.visible }"
        :socket="socket"
        :allMessages="this.allMessages"
        :ChatTarget="chatBox.ChatTarget"
        :visible="chatBox.visible"
        @close="closeChatBox(index)"
      />
  </div>
</template>

<style scoped>
.chat{
  pointer-events: none;
  position: fixed;
  display: flex;
  flex-direction: row-reverse;
  justify-content: flex-start;
  width: auto;
  height: auto;
  bottom: 0;
  padding: 1rem 1rem 0 1rem;
  right: 0;
  overflow: visible;
}
.close-btn {
  background-color: #000000;
  border: 0;
  margin-left: 0;
  margin-right: 0;
  margin-top: 2px;
  color: #00bd7e;
  transition: width 0.5s, background-color 0.5s, border-radius 0.5s;
}
.container {
  pointer-events: auto;
  background-color: #000000;
  padding-top: 0.5rem;
  padding-bottom: 1rem;
  border-top-left-radius: 1rem;
  border-top-right-radius: 1rem;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  width: 30rem;
  height: 43rem;
}

header {
  border-bottom: 0.1rem;
  border-bottom-color: #00bd7e;
  border-bottom-style: solid;
  position: relative;
  background: 0;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

header .logo {
  color: #00bd7e;
  text-decoration: none;
  font-weight: 600;
  font-size: 20px;
}

.content {
  position: relative;
  display: flex;
  transition: 0.5s;
  overflow: hidden;
}

.content .data {
  position: relative;
  width: 100%;
  height: 510px;
}

.chatlist {
  position: relative;
  height: 100%;
  overflow-y: auto;
  width: 100%;
}

.chatlist .block {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  gap: 10px;
  padding: 15px 10px;
  display: flex;
  cursor: pointer;
}

.chatlist .block :hover {
  background: #0c0c0c;
  opacity: 0.7;
}

.chatlist .block .imgbx {
  position: relative;
  min-width: 45px;
  height: 45px;
  overflow: hidden;
  border-radius: 50%;
}

.chatlist .block .imgbx img {
  position: relative;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.chatlist .block .details {
  position: relative;
  width: 100%;
}

.chatlist .block .details .listHead {
  display: flex;
  justify-content: space-between;
}

.chatlist .block .details .listHead h4 {
  font-size: 1em;
  font-weight: 600;
  width: 100%;
  color: #00bd7e;
}

.chatlist .block .details .listHead .time {
  font-size: 0.75em;
  color: #9fc4b7;
}

.message_p {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chatlist .block .details p {
  color: #2c725b;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  font-size: 0.9em;
  /*overflow: hidden;*/
  text-overflow: ellipsis;
  -webkit-box-orient: vertical;
}

.chatlist .block.unread .details .listHead .time {
  color: #0b8f33;
}

.chatlist .block.unread .details p {
  color: #45a161;
  font-weight: 600;
}

.message_p b {
  background: #0b8f33;
  color: #000000;
  min-width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 0.75em;
}

::-webkit-scrollbar {
  width: 5px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #00bd7e;
}
</style>
