<script setup>
import NavBar from "../components/NavBar.vue";
import Post from "../components/Post.vue";
import Groups from "../components/Group.vue";
import axios from "axios";
import Chat from "../components/Chat.vue"
import { selfInfo } from "../main.js"
import { sendWebSocketMessage } from "../webSocket/webSocket";
import { mapState } from "vuex";
</script>
<script>
export default {
  data() {
    return {
      userInfo: {},
      allPosts: [],
      userPosts: [],
      userDataUpdated: false,
      groups: [],
      showModal: false,
      showModal2: false,
      Comments: [],
      avatarImageUrl: "",
      showChat: false,
      showEditProfile: false,
      ClientUser: false,
      selfData: {},
      showProfileNotifications: false,
      isUserFollowing: false,
    };
  },
  components: {
    NavBar,
    Chat,
  },
  watch: {
    userDataUpdated(newValue) {
      if (newValue) {
        // Trigger a hot reload
        location.reload();
      }
    },
 
    $route(to, from) {
      if (to.path !== from.path) {
        location.reload();
      }
    },
  },
  computed: {
    ...mapState(["allNotifications"]),
    userId() {
      return this.$route.params.userId;
    },
  },
  
  async created() {
    selfInfo().then((result) => {
      this.selfData = result
    });
    (async () => {
      try {
        const response = await axios.get(`http://localhost:8000/userInfo/${this.userId}`, {
          withCredentials: true,
        });
        this.userInfo = response.data[0];
        this.Comments = this.userInfo.Comments;
        console.log(this.userId);
        this.isUserFollowing = this.isFollowing(this.selfData.Id, this.userInfo.Followers);
        this.avatarImageUrl = "http://localhost:8000/backend/" + this.userInfo.ImageUrl.String;
        console.log("all notifics", this.allNotifications)
        if (this.selfId == this.userInfo.Id) {
          this.ClientUser = true;
        }
      } catch (error) {
        console.error("Error loading userInfo:", error);
      }
    })();

  },
  methods: {
    isFollowing(userId, followers) {
      if ( followers == null){
        return false
      }
      for (const follower of followers) {
        if (follower.Id === userId) {
          return true;
        }
      }
      return false;
    },
 
    followRequest() {
      let notif = {
        sender: this.selfData.Id,
        target: this.userId,
        desc: "followRequest",
        seenByTarget: "0",
        seenBySender: "0",
        status: "pending",
        value: "0",
      };

      sendWebSocketMessage(notif, "notification");
    },
    unFollowRequest() {
      let notif = {
        sender: this.selfData.Id,
        target: this.userId,
        desc: "unFollowRequest",
        seenByTarget: "0",
        seenBySender: "1",
        status: "pending",
        value: "0",
      };

      sendWebSocketMessage(notif, "notification");
      this.isUserFollowing = this.isFollowing(this.selfData.Id, this.userInfo.Followers);

      this.userDataUpdated = true;
    },
    toggleChat() {
      console.log("toggling chat")
      if (this.showChat == false) {
        this.showChat = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showChat = false;
      }
    },
    OpenProfileEditor() {
      if (this.showEditProfile == false) {
        this.showEditProfile = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showEditProfile = false;
      }
    },
    OpenFollowers() {
      if (this.showModal == false) {
        this.showModal = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showModal = false;
      }
    },
    OpenFollowing() {
      if (this.showModal2 == false) {
        this.showModal2 = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showModal2 = false;
      }
    },
    async handleUserdataUpdate(event) {
      event.preventDefault();

      const userUpdateData = new FormData();
      userUpdateData.append("nickname", event.target.nickname.value);
      userUpdateData.append("aboutMe", event.target.aboutMe.value);
      userUpdateData.append('visibility', event.target.visibility.value);

      const avatarImage = event.target.avatarImage.files[0];
      if (avatarImage) {
        userUpdateData.append("avatarImage", avatarImage, avatarImage.name);
      }

      console.log("FrontPage.js userUpdateData:");

      for (const pair of userUpdateData.entries()) {
        console.log(`${pair[0]}, ${pair[1]}`);
      }

      (async () => {
        try {
          const response = await axios.post("http://localhost:8000/userInfo/", userUpdateData, {
            withCredentials: true,
          });

          console.log("FrontPage.js handleUserdataUpdate: ", response.data);
          this.userDataUpdated = true;
        } catch (error) {
          console.error("Error updating user data:", error);
        }
      })();
    },
  },
};
</script>
<template>
  <div class="profilePageContainer">
    <div class="profileContainer">
      <div class="userContainer">
        <!-- User Avatar -->
        <img class="user-avatar" :src="avatarImageUrl" alt="User Avatar" @error="
          $event.target.src =
          'https://wallpapers-clan.com/wp-content/uploads/2022/08/default-pfp-1.jpg'
          " />
        <!-- User Name and Bio -->
        <div class="user-info">
          <h1 class="user-name">
            {{ userInfo?.FirstName || "FirstName" }}
            {{ userInfo?.LastName || "LastName" }}
          </h1>
          <h3 class="user-Nickname">
            {{ userInfo.Nickname?.String || "Nickname" }}
          </h3>
          <br />
          <p class="user-bio">
            Born: {{ userInfo.DateOfBirth || "Users Birthday." }}
          </p>
          <p class="user-bio">
            {{ userInfo.AboutMe?.String || "About me section." }}
          </p>
          <p v-if="isUserFollowing == false && ClientUser == false && userInfo.visibility == 'private'" class="user-bio">
            This user's profile is private and you are not following this user so you are not authorised to view this users profile. 
          </p>
        </div>
      </div>
      <br />
      <button v-if="ClientUser" @click="OpenProfileEditor" class="formButton">Edit profile</button>
      <div v-if="showEditProfile" class="editUserData">
        <form @submit="handleUserdataUpdate" method="POST" class="userdataUpdate">
          <h2>Optional data update</h2>

          <div class="formgroup">
            <label >
              Privacy option:
              <br/>
              <select id="visibility" class="formField" name="visibilty" v-model="selected"  >
            <option   selected disabled value="">Select your privacy</option>
            <option value="private">Private</option>
            <option value="public">Public</option>
          </select>
            </label>
            <br>
            <label>
              Nickname:
              <br />
              <input class="formField" type="text" id="nickname" name="nickname" placeholder="Enter nickname"
                maxLength="25" autocomplete="off" />
            </label>
          </div>
          <br />
          <div class="formgroup">
            <label>
              About me:
              <br />
              <input class="formField" type="textarea" id="aboutMe" name="aboutMe" placeholder="Enter aboutMe"
                maxLength="255" autocomplete="off" />
            </label>
          </div>

          <div class="formgroup">
            <br />
            <label  for="avatarImage">
              Avatar image
              <br />
              <input class="formField" type="file" id="avatarImage" name="avatarImage" accept=".jpg,.jpeg,.png,.gif" />
            </label>
          </div>
          <br />
          <button class="formButton" type="submit">Update info</button>
        </form>
      </div>
      <div class="followsContainer">
        <div class="follow">
          <span>Followers</span>
          <a @click="OpenFollowers" class="followAmount">{{ userInfo.Followers?.length || "0" }}</a>
          <div v-if="showModal" class="followersList">
            <div v-for="follower in userInfo.Followers" :key="follower.Id">
              <span>{{ follower.FirstName }}{{ follower.LastName }}</span>
            </div>
          </div>
        </div>
        <div v-if="ClientUser == false">
          <button class="formButton" v-if="isUserFollowing" @click="unFollowRequest">Unfollow</button>
          <button class="formButton" v-else @click="followRequest">Follow user</button>
        </div>
        <div class="follow">
          <span>Following</span>
          <a @click="OpenFollowing" class="followAmount">{{ userInfo.Following?.length || "0" }}</a>
          <div v-if="showModal2" class="followingList">
            <div v-for="following in userInfo.Following" :key="following.Id">
              <span>{{ following.FirstName }}{{ following.LastName }}</span>
            </div>
          </div>
        </div>
        <br />
      </div>
    </div>

      <div class="user-posts">
        <h2>Posts</h2>
        <Post v-for="post in this.userInfo.Posts" 
        :key="post.Id" 
        :Username="this.userInfo.Nickname.String"
        :avatarImageUrl="avatarImageUrl" 
        :Title="post.Title" 
        :PostTime="post.PostTime" :Body="post.Body"
          :mediaUrl="post.mediaUrl" :Comments="post.Comments" :Likes="post.Likes" :Id="post.Id"
          :Creator="this.userInfo" />
      </div>

    <NavBar :showChat="showChat" @toggleChat="toggleChat" />
    <Chat :showChat="showChat" @toggleChat="toggleChat" />
  </div>
</template>
<style scoped>
  label[for="avatarImage"] {
    cursor: pointer;
    background-color: #1b1b1b;
    color: #00bd7e;
    padding: 5px 10px;
    
    border-radius: 4px;
  }
  label[for="avatarImage"]:hover {
    cursor: pointer;
    background-color: #00bd7e;
    color: #1b1b1b;
    padding: 5px 10px;
    
    border-radius: 4px;
  }
  input[type="file"] {
    display: none;
  }


.followsContainer {
  display: flex;
  flex-direction: row;
  box-shadow: 0.2rem 0.2rem #1a1a1a4b;
  justify-content: center;
}
.follow{ 
  margin: 0 10rem;
  display: inline-flex;
  flex-direction: row;
}
.follow span {
  color: #00bd7e;
  background-color: #006644;
  border: 0;

  padding: 0.5rem;
  border-top-left-radius: 0.5rem;
  border-bottom-left-radius: 0.5rem;
  cursor: pointer;
}

.followAmount {
  color: #1b1b1b;
  background-color: #00bd7e;
  border: 0;
  padding: 0.5rem;
  border-top-right-radius: 0.5rem;
  border-bottom-right-radius: 0.5rem;
  cursor: pointer;
}

.profilePageContainer {
  margin-top: 5rem !important;
  margin: auto;
  display: flex;
  width: 100%;
  flex-direction: column;
}

.userContainer {
  display: flex;
  flex-direction: row;
}

.profileContainer {
  position: relative;
  top: 1%;
  background-color: #101010;
  width: 100%;
  border-radius: 1rem;
  box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
  padding: 2rem 0rem 2rem 2rem;
  margin: 2.5rem 0rem 2.5rem 0rem;
}

.user-info {
  margin-left: 1rem;
}

.user-posts {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.user-groups {
  margin-top: 5rem;
  display: flex;
  flex-direction: column;
}

.user-avatar {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  margin-right: 8px;
  margin-left: 8px;
}
</style>
