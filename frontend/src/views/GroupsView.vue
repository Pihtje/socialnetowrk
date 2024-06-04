<script setup>
import NavBar from "../components/NavBar.vue";
import Group from "../components/Group.vue";
import axios from "axios";
import Chat from "../components/Chat.vue"
import Events from "../components/Events.vue";
import { selfInfo } from "../main.js";
</script>
<script>
export default {
  data() {
    return {
      groups: [],
      selfData: {},
      showChat: false,
    };
  },
  components:{
    Chat,
    NavBar,
    Group,
    Events,
  },
  

  async created() {
    selfInfo().then((result) => {
      this.selfData = result
});
    try {
      const response = await axios.get("http://localhost:8000/group/", {
        withCredentials: true,
      });
      console.log("all notifics", this.allNotifications)
      this.groups = response.data;
    } catch (error) {
      console.error("Error loading posts:", error);
    }
  },
  methods: {
   
    toggleChat() {
      console.log("Allnotifications",this.allNotifications)
      console.log("toggling chat")
      if (this.showChat == false) {
        this.showChat = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showChat = false;
      }
    },
    async handleNewGroupSubmit(event) {
      event.preventDefault();

      const formData = new FormData();
      formData.append("title", event.target.groupTitle.value);
      formData.append("body", event.target.groupDesc.value);
      console.log("Form data:", formData);

      try {
        const response = await axios.post(
          "http://localhost:8000/group/",
          formData,
          {
            withCredentials: true,
          }
        );
        try {
      const response = await axios.get("http://localhost:8000/group/", {
        withCredentials: true,
      });

      this.groups = response.data;
    } catch (error) {
      console.error("Error loading posts:", error);
    }
        const data = response.data;
        console.log("handleNewGroupSubmit: ", data);
      } catch (error) {
        console.error("Failed to fetch:", error);
      }
    },
  }
};
</script>

<template>
  <div id="app">
    <div class="createGroupContainer" >
      <h1 style="font-weight: bold;">Create a group!</h1>
      <br>
      <form class="createGroup" @submit="handleNewGroupSubmit">
        <label for="groupTitle">Title:</label>
        <input
          class="groupForm"
          type="text"
          name="groupTitle"
          required
        />
        <br />
        <label for="groupDesc">Group description</label>
        <textarea class="groupForm" name="groupDesc" rows="5" required></textarea>
        <br />
        <button class="formButton" type="submit">Create new Group!</button>
      </form>
    </div>
    <div class="groupsSection" v-if="groups">
      <Group
        v-for="group in groups"
        class="groupCard"
        :key="group.GroupId"
        :GroupId="group.GroupId"
        :Creator="group.Creator"
        :GroupTitle="group.GroupTitle"
        :GroupDesc="group.GroupDesc"
        :Members="group.Members"
        :Posts="group.Posts"
      />
    </div>
    <div v-else>
      <p>Loading groups...</p>
    </div>
    <NavBar :showChat="showChat" @toggleChat="toggleChat"/>
    <Chat :showChat="showChat" @toggleChat="toggleChat" />
  </div>
</template>

<style scoped>
.groupCard{
  padding: 1rem;
}

.groupsSection {
    margin-top: 6rem;
    height: auto;
    flex-direction: column;
    display: flex;
    width: 50%;
    align-items: flex-end;
}

.groupForm {
  margin-top: 1rem;
  
}

.createGroup {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.createGroupContainer {
    display: flex;
    flex-direction: column;
    background-color: #101010;
    border-radius: 1rem;
    box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
    width: 50%;
    padding: 2rem;
    margin-top: 6rem;
    margin-right: 2rem;
}

.formButton {
  width: fit-content;
}

input[type=text], input[type=file], textarea {
  width: 100%;
}
</style>
