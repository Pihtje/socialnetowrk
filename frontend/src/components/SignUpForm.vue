<script setup>
import { RouterLink } from "vue-router";
import axios from "axios";
import { selfInfo } from "../main";
</script>
<script>
export default {
  data() {
    return {
      email: "",
      username: "",
      password: "",
      repeatPassword: "",
      firstName: "",
      lastName: "",
      birthDay: "",
      aboutMe: "",
      visibility: "private",
      err: "",
    };
  },
  validations() {
    return {
      email: {
        required,
        email,
        isUnique: withAsync(
          withMessage("Oops, Email already in use", async (value) => {
            if (value === "") {
              return true;
            }
            const resp = await axios.post("isUnique", { Email: value });
            return resp.data;
          })
        ),
      },
      password: {
        required,
        minLength: minLength(8),
        maxLength: maxLength(50),
        printableChars,
      },
      passwordConfirm: {
        required,
        sameAs: sameAs(this.password),
      },
      firstName: {
        required,
        maxLength: maxLength(50),
        alpha,
      },
      lastName: {
        required,
        maxLength: maxLength(50),
        alpha,
      },
      birthDay: {
        required,
        between: between(0, 120),
        integer,
      },
      username: {
        minLength: minLength(3),
        maxLength: maxLength(20),
        alphaNum,
        isUnique: withAsync(
          withMessage("Oops, Nickname already taken", async (value) => {
            if (value === "") {
              return true;
            }
            const resp = await axios.post("isUnique", { username: value });
            return resp.data;
          })
        ),
      },
      aboutMe: {
        maxLength: maxLength(500),
        alphaNum,
        printableChars,
      },
    };
  },
  methods: {
    async handleSignupSubmit(event) {
      event.preventDefault();
      const signupData = new FormData();
      signupData.append("email", this.email);
      signupData.append("password", this.password);
      signupData.append("repeatPassword", this.repeatPassword);
      signupData.append("firstName", this.firstName);
      signupData.append("lastName", this.lastName);
      signupData.append("dateOfBirth", this.birthDay);

      signupData.append("nickname", this.username);
      signupData.append("aboutMe", this.aboutMe);

      signupData.append("visibility", this.visibility);

      const avatarImage = event.target.avatarImage.files[0];
      if (avatarImage) {
        signupData.append("avatarImage", avatarImage, avatarImage.name);
      }

      try {
        const response = await axios.post(
          "http://localhost:8000/signup/",
          signupData,
          {
            withCredentials: true,
          }
        );
        // Handle successful login response
        console.log(response.data);
        if (response.data === "Registration successful!") {
          selfInfo();
          this.$router.push({ name: "feed" });
        } else {
          this.err = response.data;
        }
      } catch (error) {
        // Handle error
        this.err = error;
        console.error(error);
      }
    },
  },
};
</script>
<template>
  <div class="formDiv">
    
    <form class="formContent" @submit="handleSignupSubmit">
      
      <div class="formFields"> 
        <div class="requiredInfo">
        <label>Email:</label><br />
        <input
          class="formField"
          v-model="email"
          type="text"
          name="email"
          required
        /><br /><br />

        <label for="password">Password:</label><br />
        <input
          class="formField"
          v-model="password"
          type="password"
          name="password"
          required
        /><br /><br />

        <label for="password">Repeat Password:</label><br />
        <input
          class="formField"
          v-model="repeatPassword"
          type="password"
          name="repeatPassword"
          required
        /><br /><br />

        <label for="firstName">First Name:</label><br />
        <input
          class="formField"
          v-model="firstName"
          type="text"
          name="firstName"
          required
        /><br /><br />

        <label for="lastName">Last Name:</label><br />
        <input
          class="formField"
          v-model="lastName"
          type="text"
          name="lastName"
          required
        /><br /><br />

        <label for="birthday">Date of Birth:</label><br />
        <input
          class="formButton"
          v-model="birthDay"
          type="date"
          name="birthday"
          required
        /><br /><br />
      </div>
      <div class="optionalInfo">
        <h3>Optional Info</h3>
        <p>This info can be changed later from your profile page.</p>
        <br />

        <label for="username">Nickname:</label><br />
        <input
          class="formField"
          v-model="username"
          type="text"
          name="username"
        /><br /><br />

        <label for="aboutMe">About me:</label><br />
        <input
          class="formField"
          v-model="aboutMe"
          type="textarea"
          name="aboutMe"
        /><br /><br><br>

        <label for="avatarImage">Avatar image</label><br />
        <input
          class="formField"
          type="file"
          id="avatarImage"
          name="avatarImage"
          accept=".jpg,.jpeg,.png,.gif"
        /><br /><br />

        <label for="visibility">Account privacy:</label><br />
        <select class="formButton" name="visibility" id="visibility" v-model="visibility">
          <option value="public">Public</option>
          <option value="private">Private</option>
        </select>

        
      </div></div>
      <button class="formContent formButton" type="submit">Submit</button>
    </form>
    
    <p v-if="err" class="error">{{ err }}</p>

    <RouterLink class="formContent" to="/"
      >Already a member? Click here to login!</RouterLink
    >
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
.formDiv {
  width:fit-content ;
}

.formFields{
  margin-top: 1rem;
  width:fit-content;
  display: flex;
  flex-direction: row;
}
.optionalInfo{
  margin-top: 0;
  margin-left: 2rem;
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
</style>
