<script >
import { mapState } from "vuex";
import { sendWebSocketMessage } from "../webSocket/webSocket";
import EmojiPicker from "./EmojiPicker.vue";
export default {
  data() {
    return {
      newMessage: "",
      showEmojiPicker: false,
    };
  },
  components: {
    'emoji-picker': EmojiPicker, // Register the emoji picker component
  },
  props: {
    ChatTarget: {},
    socket: Object,
  },
  computed: {
    ...mapState(["allMessages"]),
    filteredMessages() {
    return this.allMessages.filter(message => message.groupId === "0");
    },
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
  methods: {
    openEmojiPicker() {
      //this.showEmojiPicker = true;
      this.showEmojiPicker = !this.showEmojiPicker;
    },
    insertEmoji(emoji) {
      // Insert the selected emoji into the input field
      this.newMessage += emoji;
    },
    closeEmojiPicker() {
      this.showEmojiPicker = false;
    },
    close() {
      this.$emit("close");
    },
    sendDM() {
      console.log("sending DM...");
      const timestamp = new Date().toISOString();

      console.log(this.userData.Id, this.ChatTarget.Id, this.newMessage, timestamp);

      if (this.userData.Id && this.ChatTarget.Id && this.newMessage) {
        const message = {
          sender: this.userData.Id + "",
          target: this.ChatTarget.Id + "",
          msg: this.newMessage,
          serverObject: true,
          groupId: 0 + "",
          DateTime: timestamp,
        };

        // Use the sendWebSocketMessage function to send the message through WebSocket
        sendWebSocketMessage(message, "message");
        this.newMessage = "";
        this.$nextTick(() => {
      const messageContainer = this.$refs.messageContainer;
      messageContainer.scrollTop = messageContainer.scrollHeight;
    });
      }
    },
  },
};
</script>
<template>
  <div class="chatBox">
    <div class="chat_header">
      <div class="imgcontent">
        <link
          rel="stylesheet"
          href="https://cdn.lineicons.com/1.0.1/LineIcons.min.css"
        />
        <div class="imgBx">
          <img
            class="user-avatar"
            :src="BackendImg + ChatTarget.ImageUrl.String"
            alt="User Avatar"
            @error="
              $event.target.src =
                'https://wallpapers-clan.com/wp-content/uploads/2022/08/default-pfp-1.jpg'
            "
          />
        </div>
        <h3>
          {{ ChatTarget.FirstName }}<br /><span>{{
            ChatTarget.OnlineStatus
          }}</span>
        </h3>
      </div>
      <div class="back">
        <button
          @click="close()"
          class="lni lni-close size-sm close-btn"
        ></button>
      </div>
    </div>
    <div class="messagesContainer">
      <div
      ref="messageContainer"
        class="messageBox"
        v-for="message in this.allMessages"
        :key="message.messageId"
      >
        <div
          class="message"
          :class="{
            incoming: message.sender.Id == ChatTarget.Id && message.target != 'group',
            sent: message.sender.Id != ChatTarget.Id && message.target != 'group',
          }"
          v-if="
            message.sender.Id == ChatTarget.Id  && message.target != 'group'|| message.target.Id == ChatTarget.Id && message.target != 'group'
          "
        >
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
        @keyup.enter="sendDM" 
      />
      <button class="formButton" @click="openEmojiPicker">😀</button>

<!-- Display the emoji picker when it's open -->
<emoji-picker
  v-if="showEmojiPicker"
  @emoji-selected="insertEmoji"
  @close-picker="closeEmojiPicker"
/>
      <button class="formButton" @click="sendDM">Send</button>
    </div>
  </div>
</template>
<style scoped>
.emoji-picker {
  position: absolute;
  background: white;
  border: 1px solid #ccc;
  border-radius: 5px;
  padding: 10px;
  z-index: 100;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  max-height: 150px;
  overflow-y: auto;
}

button {
  display: block;
  margin-top: 10px;
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
  pointer-events: auto;
  align-self: flex-end;
  border-radius: 1rem;
  width: 21rem;
  max-height: 30rem;
  background: rgb(0, 0, 0);
  z-index: 10000;
  margin: 0 0.5rem;
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

.chat_header .imgcontent {
  display: flex;
  align-items: center;
  gap: 5px;
}

.chat_header .imgcontent .imgBx {
  position: relative;
  width: 35px;
  height: 35px;
  border-radius: 50%;
  overflow: hidden;
}

.chat_header .imgcontent .imgBx img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  object-fit: cover;
}

.chat_header .imgcontent h3 {
  color: #00bd7e;
  font-size: 1em;
  font-weight: 500;
  line-height: 1.1em;
}

.chat_header .imgcontent h3 span {
  font-size: 0.7em;
  font-weight: 400;
}

.messagesContainer {
  overflow-x: hidden;
  overflow-y: auto;
  position: relative;
  padding: 20px;
  height: 20rem;
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
