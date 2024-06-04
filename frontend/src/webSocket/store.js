// store.js
import {createStore} from 'vuex';

const store = createStore({
  state() {
    return {
      allUsers: [],
      allMessages: [], // Initialize allMessages as an empty array
      allNotifications: [],
      latestMessageTimestamps: {},
    };
  },
    mutations: {
      setAllUsers(state, users) {
        state.allUsers = users;
      },
      setAllMessages(state, messages) {
        state.allMessages = messages;
      },
      setAllNotifications(state, notifications) {
        state.allNotifications = notifications;
      },
    },
    actions: {
      // Your actions here
    },
  });
  
export { store };