<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios'
import PopupNotification from './components/PopupNotification.vue'
import { useEventBus } from './webSocket/EventBus';
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div id="app">
      <router-view />
      <popup-notification v-if="notificationMessage" :notification="notificationMessage" />
  </div>
</template>
<script>

export default {
  setup () {
    axios.get('http://localhost:8000/')
      .then(response => {
        //handle response
        console.log("Axios response",response)
      })
  },
}
</script>
