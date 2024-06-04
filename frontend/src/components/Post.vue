<script>
import axios from "axios";

export default {
  props: {
    Id: {
      type: Number,
      Required: true,
    },
    AvatarImageUrlProp: {
      type: String,
      required: false,
    },
    PostTime: {
      type: String,
      required: true,
    },
    ImageUrl: {
      type: String,
      required: false,
    },
    Title: {
      type: String,
      deafult: null,
    },
    Body: {
      type: String,
      required: true,
    },
    Comments: {
      type: Array,
      required: true,
      default: null, 
      default: () => [],
      validator: (Comments) => {
        // Validate each comment in the array
        return Comments.every((Comment) => {
          return (
            typeof Comment.Id === "number" &&
            typeof Comment.PostId === "number" &&
            typeof Comment.UserId === "number" &&
            typeof Comment.Username === "string" &&
            typeof Comment.CommentTime === "string" &&
            typeof Comment.Body === "string" &&
            typeof Comment.Likes === "number" &&
            typeof Comment.Dislikes === "number" &&
            typeof Comment.ImageUrl === "string"  &&
            typeof Comment.Creator
          );
        });
      },
    },
    Categories: [],
    Likes: {
      type: Number,
      deafult: null,
    },
    Creator: {
      required: true,
    }
  },
  data() {
    return {
      userInfo: {},
      showModal: false,
      newComment: "",
      postAuthor: this.Creator.Nickname.String || this.Creator.FirstName + " " + this.Creator.LastName,
      postImageUrl: "",
      commentAvatarUrl: "",
      avatarImageUrl: "",
    };
  },
  async created() {
  if (this.ImageUrl && this.ImageUrl.trim() !== "") {
    this.fetchImageFromBackend();
  }
},
  computed: {
    commentsArray() {
    // If this.Comments is null, return an empty array
    return this.Comments || [];
    },
    formattedPostTime() {
      const date = new Date(this.PostTime);
      const day = date.getDate().toString().padStart(2, "0");
      const month = (date.getMonth() + 1).toString().padStart(2, "0");
      const year = date.getFullYear();
      const hours = date.getHours().toString().padStart(2, "0");
      const minutes = date.getMinutes().toString().padStart(2, "0");
      return `${day}.${month}.${year} ${hours}:${minutes}`;
    },
    formattedCommentTime() {
      return function (comment) {
        const date = new Date(comment.CommentTime);
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
    goToProfile(userId) {
      this.$emit('author-clicked', userId);
    },
    OpenComments() {
      if (this.showModal == false) {
        this.showModal = true; // Open the modal when the "Comments" button is clicked
      } else {
        this.showModal = false;
      }
    },
    async fetchUserInfo() {
      try {
        const response = await axios.get("http://localhost:8000/userInfo/self", {
          withCredentials: true,
        });
        this.userInfo = response.data[0];
        console.log("Userinfo", this.userInfo);
      } catch (error) {
        console.error("Error loading user info:", error);
      }
    },
    async fetchImageFromBackend() {
      // this.postImageUrl = "http://localhost:8000/backend/" + this.ImageUrl
      // this.avatarImageUrl = "http://localhost:8000/backend/" + this.Creator.ImageUrl.String
      // this.commentAvatarUrl = "http://localhost:8000/backend"
    },
    handleCommentSubmit() {
      const postId = this.Id; // Replace with your post ID
      const formData = new FormData();
      console.log("PostId", postId, "\n Formdata", formData);
      formData.append("newCommentBody", this.newComment); // Adjust field name
      this.fetchUserInfo()
        .then(() => {

        })
        .then(() => {
          const currentTime = new Date();
          const commentTime = currentTime.toISOString(); // Convert to ISO string format
          axios
            .post(`http://localhost:8000/post/${postId}`, formData, {
              withCredentials: true, // Allow sending credentials
            })
            .then((response) => response.data)
            .then((data) => {
              console.log("Comment submission response: ", data);
              console.log("userinfo", this.userInfo.Id);
              // Update the local comments array with the new comment
              const newComment = {
                Id: data.id, // Assuming the response data has the new comment's ID
                PostId: postId,
                UserId: this.userInfo.Id, // Assuming the response data has the user ID
                Username: this.userInfo.FirstName, // Assuming the response data has the username
                CommentTime: commentTime, // Assuming the response data has the comment time
                Body: this.newComment,
                ImageUrl: data.imageUrl,
                Creator: this.userInfo // Assuming the response data has the image URL
              };
              console.log(newComment);
              this.Comments.push(newComment);
              // this.$emit('refreshPosts');
              //console.log("this.Comments Capital C", this.Comments)
              this.newComment = ""; // Clear the input field
            })
            .catch((error) => {
              console.error("Error submitting comment: ", error);
            });
        })
        .catch((error) => {
          console.error("Error fetching user info:", error);
        });
    },
  },
};
</script>
<template>
  <div class="post-container">
    <!-- Post Author Info -->
    <div @click="goToProfile(this.Creator.Id)" class="post-author">
      <img
        :src="BackendImg + Creator.ImageUrl.String"
        
        class="author-avatar"
        @error="
          $event.target.src =
            'https://wallpapers-clan.com/wp-content/uploads/2022/08/default-pfp-1.jpg'
        "
      />
      <div class="header-texts">
        <span class="author-name">{{ Creator.Nickname?.String || this.postAuthor }}</span>
        <span class="post-time">{{ formattedPostTime }}</span>
      </div>
    </div>

    <!-- Post Content -->
    <div class="post-content">
      <h2>{{ Title }}</h2>
      <p>{{ Body }}</p>
      <div>
        <img :src="BackendImg + ImageUrl" alt="" class="post-media" />
      </div>
    </div>
    <!-- Post Interactions -->
    <div class="post-interactions">

      <div class="postButton" @click="OpenComments">
        <span>Comments</span>
      </div>
    </div>
    <!-- Use the Modal component to display the post as a popup -->
    <div class="commentsSection" v-if="showModal">
      <br />
      <!-- Loop through the comments array and display the comments -->
      <div v-for="comment in Comments" :key="comment.Id" class="comment">
        <img
          class="commentAvatar"
          :src="BackendImg + comment.Creator.ImageUrl.String"
          alt="pfp"
          @error="
            $event.target.src =
              'https://wallpapers-clan.com/wp-content/uploads/2022/08/default-pfp-1.jpg'
          "
        />
        <div class="comment-content-and-time">
          <div class="commentContent">
            <span class="comment-username">{{ comment.Creator.Nickname?.String || comment.Creator.FirstName +' '+ comment.Creator.LastName }}</span>
            <span class="comment-content">{{ comment.Body }}</span>
          </div>
          <span class="comment-time">{{ formattedCommentTime(comment) }}</span>
        </div>
      </div>
      <br />
      <div class="addComment">
        <input
          @keyup.enter="handleCommentSubmit"
          class="commentForm"
          v-model="newComment"
          placeholder="Add a comment"
          maxLength="255"
          minlength="2"
        />
        <button class="postButton" @click="handleCommentSubmit">
          Add Comment
        </button>
      </div>
    </div>
  </div>
</template>
<style scoped>
.commentsSection {
  border-style: solid;
  border-color: #00bd7e;
  border-width: 0.1rem 0px 0rem 0px;
}
.comment {
  display: flex;
  flex-direction: row;
  margin: 1rem;
}
.commentAvatar {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  margin-right: 8px;
  margin-left: 8px;
}
.comment-username {
  font-weight: bold;
}
.comment-content-and-time {
  margin-left: 1rem;
}
.commentContent {
  padding: 0.5rem;
  border-radius: 0.5rem;
  width: fit-content;
  display: flex;
  flex-direction: column;
  background-color: #1b1b1b;
}
.addComment {
  margin: 1rem;
  display: flex;
  flex-direction: row;
}
.post-container {
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

.post-content {
  display: flex;
  flex-direction: column;
  height: 75%;
  align-items: center;
  margin-bottom: 10px;
  padding: 1rem 0.5rem 1rem 0.5rem;
}

.post-author {
  display: flex;
  padding-top: 8px;
  padding-bottom: 1rem;
  border-style: solid;
  border-color: #00bd7e;
  border-width: 0px 0px 0.1rem 0px;
}

.author-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 8px;
  margin-left: 8px;
}
.header-texts {
  margin-left: 1rem;
  align-self: center;
  font-weight: bold;
  display: flex;
  flex-direction: column;
}
.author-name {
  margin-left: 1rem;
  align-self: center;
  font-weight: bold;
}
.post-time,
.comment-time {
  font-size: smaller;
  margin-left: 1rem;
  align-self: center;
}

.post-interactions {
  padding: 1rem;
  border-style: solid;
  border-color: #00bd7e;
  border-width: 0.1rem 0px 0rem 0px;
  /* display: flex;
  align-self: flex-end;
  justify-content: space-between;
  margin-top: 1rem; */
}

.post-media {
  max-width: 100%;
  height: 100%;
  border-radius: 8px;
  object-fit: cover;
}

.postButton {
  color: #00bd7e;
  background-color: #1b1b1b;
  border: 0;
  margin-left: 5rem;
  margin-right: 5rem;
  padding: 0.6rem;
  border-radius: 0.5rem;
  box-shadow: 0.2rem 0.2rem #1a1a1a4b;
  display: flex;
  cursor: pointer;
}

.postButton i {
  margin-right: 5px;
}

.postButton .liked {
  color: red;
}
</style>
