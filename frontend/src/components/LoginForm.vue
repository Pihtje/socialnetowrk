<script setup>
import { RouterLink } from 'vue-router';
import axios from 'axios';
import {selfInfo} from "../main.js"
</script>

<template>
  <div class="formDiv">
    <br>
   <form class="formContent" >
    <label>Email:</label><br>
    <input class="formField" type="text" 
          v-model="email"
					autocomplete="email"
          /><br><br>
    <label >Password:</label><br>
    <input  @keyup.enter="handleLoginSubmit" class="formField" v-model="password" type="password"/><br><br>
  </form> 
    <p v-if="wrongCredentials" class="formContent" style="color: red;">{{ this.wrongCredentials }}</p>
    <br>
    <button @click="handleLoginSubmit" class="formButton formContent" type="button">Login</button>  <br><br>
  <RouterLink class="formContent" to="/signup">No account? Click here to register!</RouterLink>
  </div>
</template>
<script>
export default {
  data() {
    return {
      wrongCredentials: false,
      email: '',
      password: ''
    };
  },
  methods: {
    async handleLoginSubmit() {
      const loginData = {
        Email: this.email,
        Password: this.password
      };
      
      try {
        const response = await axios.post('http://localhost:8000/login/', loginData, {
          withCredentials: true
        });
        console.log(response.data);
        if (response.data === "login successful") {
          this.wrongCredentials = false
					this.$router.push({name: 'feed'})  
        } else {
          this.wrongCredentials = response.data;
        }
        selfInfo()
      } catch (error) {
        console.error(error);
      }
    }
  },
  async created(){
    const u = await selfInfo();
    if (u.Id) {
      console.log("guests only");
      this.$router.push({name: 'feed'});
    }
  },
}
</script>
