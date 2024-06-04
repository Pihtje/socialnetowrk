<script setup>
import NavBar from "../components/NavBar.vue";
import axios from "axios";
import { selfInfo } from "../main.js";
</script>
<script>
export default {
  components: {
    NavBar,
  },
  data() {
    return {
      postImage: null,
      err: '',
      selectedCategory: 1,
      selfData: {},
      sharedUsers: [],
      sharedUserNames: [],
      selectedUser: 0,
    };
  },
  created(){
    selfInfo().then((result) => {
      this.selfData = result
    });
  },
  methods: {
    async handleNewPostSubmit(event) {
      event.preventDefault();
      const formData = new FormData();
      formData.append("Title", event.target.postTitle.value);
      formData.append("Body", event.target.postBody.value);
      formData.append("Category", this.selectedCategory);

      console.log("shared users:", event.target);

      formData.append("ShareWith", this.sharedUsers);

      formData.append("GroupId", 0);
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
        this.err = data;

        if (response.data === "New post added!") {
          this.$router.push({ name: "feed" });
        }
      } catch (error) {
        console.error("Failed to fetch:", error);
      }
    },
    addUser() {
      if (!this.selectedUser || this.sharedUsers.indexOf(this.selectedUser) !== -1) {
        return
      } 
      this.sharedUsers.push(this.selectedUser);

      let user = this.selfData.Followers.filter(f => f.Id === this.selectedUser)[0]
      console.log("user:", user);
      this.sharedUserNames.push(user.Nickname.String || user.FirstName + ' ' + user.LastName);
    },
    removeUser() {
      this.sharedUsers = this.sharedUsers.filter((u) => u !== this.selectedUser);

      this.sharedUserNames = [];
      this.sharedUsers.forEach((id)=> {
        let user = this.selfData.Followers.filter(f => f.Id === id)[0];
        let name = user.Nickname.String || user.FirstName + ' ' + user.LastName;
        console.log("name:", name);
        this.sharedUserNames.push(name);
      })
      
    }
  },
};
</script>
<template>
  <div class="createPostPageContainer" id="app">
    <NavBar />
    <div class="createPostContainer">
      <form class="createPost" @submit="handleNewPostSubmit" autocomplete="off">
        <h1>Creating a new post!</h1>
        <br>
        <label for="postTitle">Title:</label>
        <input class="postForm" type="text" name="postTitle" required/>
        <br />
        <label>Choose a privacy option:</label>
        <select class="postForm" v-model="selectedCategory" name="postCategories" required>
          <option value="1">Public</option>
          <option value="2">Private</option>
          <option value="3">Shared</option>
        </select>
        <div class="createPost" v-if="selectedCategory == 3">
          <select class="postForm" v-model="selectedUser">
            <option name="sharedUsers" id="sharedUsers" v-for="user in this.selfData.Followers" :value="user.Id">{{ user.Nickname.String || user.FirstName + ' ' + user.LastName }}</option>
          </select>
          <br>
          <div>
            <button type="button" class="formButton" @click="addUser">Add user</button>
            <button type="button" class="formButton" @click="removeUser">Remove user</button>
          </div>
          <br>
          <p class="green">Currently sharing with:</p>
          <br>
          <p v-if="this.sharedUserNames.length !== 0">{{ this.sharedUserNames.join(', ') }}</p>
          <p v-else>No users selected!</p>
        </div>
        
        <br />
        <label for="postBody">Post body</label>
        <textarea class="postForm" name="postBody" required rows="5"></textarea>
        <br />
        <input
          @change="onFileSelected"
          class="postForm"
          type="file"
          id="postImages"
          accept=".jpg,.jpeg,.png,.gif"
          name="postImages"
        />
        <br />
        <button class="formButton" type="submit">Create new post!</button>
      </form>
      <p v-if="err" class="error">{{ err }}</p>
    </div>
  </div>
</template>
<style scoped>
.postForm {
  margin-top: 1rem;
  background-color: #1b1b1b;
}
.createPostPageContainer {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.createPostContainer {
  align-self: center;
  background-color: #101010;
  border-radius: 1rem;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  width: 50rem;
  padding: 2rem 0rem 2rem 2rem;
  margin: auto;
}
.createPost {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.error {
  font-size: 0.8rem;
  width: fit-content;
  text-align: center;
  color: red;
  margin-top: 1rem;
  padding: 0.5rem;
  border-radius: 1rem;
  border: 1px solid red;
}

button[type=button] {
  display: inline-block;
  margin: 0 5px;
}

button[type=submit] {
  width: fit-content;
}

input[type=text], input[type=file], textarea {
  width: 50%;
}

</style>
