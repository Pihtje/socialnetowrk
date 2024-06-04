<template>
  <p>notif: {{ notif }}</p>
  <div :key="notif.value.event_id">
    <h3>Event title: {{ notif.value.title }}</h3>
    <p>Event description: {{ notif.value.description }}</p>
    <p>Date & Time: {{ notif.value.day_time }}</p>
    <div :key="notif.value.event_id">
        <p>Attending users:</p>
        <div v-if="eventResponses">
          <ul>
            <li v-for="response in eventResponses" :key="response.response_id">
              {{ response.user.Email }} - {{ response.response }}
            </li>
          </ul>
        </div>
        <div v-else>
          <p>Loading responses...</p>
        </div>
      <div v-if="userResponseStatusForEvent(notif?.value.event_id) === 'pending'">
        <button @click="acceptEventInvite()">Attending</button>
        <button @click="rejectEventInvite()">Not attending</button>
      </div>
    </div>
  </div>
</template>

<script>
import { sendWebSocketMessage } from "../webSocket/webSocket";
import axios from "axios";

export default {
  props: {
    notif: {
      type: Object,
      required: true,
    }
  },
  data() {
    return {
      notif: {},
    }
  },
  methods: {
    acceptEventInvite(){
      let n = {
        notificationId: this.notif.notificationId,
        sender: "group",
        target: this.notificationId.target.Id,
        desc: "event",
        seenByTarget: "1",
        seenBySender: "0",
        status: "accepted",
        value: this.notif.value.event_id,
      };

      sendWebSocketMessage(n, "notification");
      console.log("event invite accepted!");
    },
    rejectEventInvite(){
      let n = {
        notificationId: this.notif.notificationId,
        sender: "group",
        target: this.notificationId.target.Id,
        desc: "event",
        seenByTarget: "1",
        seenBySender: "0",
        status: "rejected",
        value: this.notif.value.event_id,
      };

      sendWebSocketMessage(n, "notification");
      console.log("event invite rejected");
    },
  },
  computed: {
    async eventResponses() {
      try {
        let response = await axios.get(
        `http://localhost:8000/getUserEventStatus/${notif.value.event_id}`, {
          withCredentials: true,
        });
        return response.data
      } catch (error) {
        console.log("Error fetching event responses:", error);
        return null
      } 
    },
    userResponseStatusForEvent() {
      try {
        this.eventResponses
      } catch (error) {
        console.log("userResponseStatusForEvent err:", error);
        return null
      }
      
      const responseObj = this.eventResponses.find(response => response.EventID === eventId && response.UserID === notif.target.Id); // Updated the userId reference
      return responseObj ? responseObj.Response : null;
    },
  }
};
</script>

<style scoped>
.eventContainer {
  background-color: #101010;
  padding-top: 0.5rem;
  padding-bottom: 1rem;
  border-radius: 1rem;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  width: 35rem;
  margin-top: 2.5rem;
  position: relative;
  align-self: center;
}

.eventContainer ul {
  list-style: none;
  padding: 0;
}

.eventContainer li {
  margin-bottom: 0.5rem;
}

.eventContainer button {
  background-color: #00bd7e;
  border: none;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:hover {
    background-color: #008d65;
  }
}
</style>