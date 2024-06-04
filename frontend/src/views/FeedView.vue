<script setup>
import NavBar from "../components/NavBar.vue";
import Post from "../components/Post.vue";
import axios from "axios";
import Chat from "../components/Chat.vue"
import { selfInfo } from "../main";
</script>
<script>
export default {
  components: {
    Post,
    NavBar,
    Chat
  },
  data() {
    return {
      posts: [],
      showChat: false,
    };
  },
  async created() {
    selfInfo().then(res => {
      if (!res.Id) {
        console.log("logged in only");
        this.$router.push({name: 'login'});
        return
      }
    });

    try {
      const response = await axios.get("http://localhost:8000/allPosts/", {
        withCredentials: true,
      });

      this.posts = response.data;
    } catch (error) {
      console.error("Error loading posts:", error);
    }
  },
  methods: {
    toggleChat() {
      console.log("toggling chat")
      // this.showChat = !this.showChat;
      if (this.showChat == false) {
        this.showChat = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showChat = false;
      }
    },
    goToProfile(userId){
      console.log(userId)
      this.$router.push({ name: 'profile', params: { userId } });  
    },
  },
};
</script>
<template>
  <div id="app">
    <div class="postDiv">
      <Post
        v-for="post in posts"
        :key="post.Id"
        :Username="post.Username"
        :PostTime="post.PostTime"
        :Avatar="post.Avatar"
        :Title="post.Title"
        :Body="post.Body"
        :ImageUrl="post.ImageUrl"
        :Comments="post.Comments"

        :Id="post.Id"
        :Creator="post.Creator"
        @author-clicked="goToProfile"
      />
    </div>
    <NavBar :showChat="showChat" @toggleChat="toggleChat"/>
    <Chat :showChat="showChat" @toggleChat="toggleChat" />
  </div>
</template>

<style scoped>
.postDiv {
  margin-top: 5rem !important;
  margin: auto;
  height: auto;
  flex-direction: column-reverse;
  display: flex;
  align-content: center;
}
</style>
