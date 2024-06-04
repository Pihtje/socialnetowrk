import './assets/main.css'
import axios from "axios";
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import {store} from './webSocket/store'; 
const app = createApp(App)
app.use(store);
// export async function webSocket() {

//     console.log("app.js: setting up socket...");
//     app.config.globalProperties.socket = new WebSocket("ws://localhost:8000/chat/");

//     this.socket.onopen = () => {
//       console.log("App.js: WS connection successful!");
//     };

//     this.socket.onmessage = (msg) => {
//       //console.log("App.js: new message: ", msg.data);
//       let msgJson = JSON.parse(msg.data);

//       console.log("msgJson: ", msgJson);

//       switch (msgJson.sender) {
//         case "server":
//           switch (msgJson.msg) {
//             case "allUsers":
//               //console.log("App.js: updating all users list", msgJson);
//               app.config.globalProperties.allUsers = msgJson.ServerObject.filter((user) => {
//                 user.lastSentMessage = new Date(
//                   "0001-01-01T00:00:00.000Z"
//                 ).toISOString();
//                 return user;
//               });
//               console.log("Chattable users:", allUsers);
//               //Re-renders the usersContainer when other users log on/off
//               // sortUsers();
//               // updateUsersContainer();

//               break;
//             case "getAllDMs":
//               console.log("Setting all DMs");
//               // setUpdateMessage(msgJson.ServerObject);
//               app.config.globalProperties.allMessages = msgJson.ServerObject;
//               console.log("All messages", allMessages);
//               // sortMessages();
//               // updateMessagesContainer();
//               // sortUsers();
//               // updateUsersContainer();
//               //updateHeader();
//               break;
//               case "getAllNotifications":
//               //console.log("App.js: updating notifications list", msgJson);
//               app.config.globalProperties.allNotifications = msgJson.ServerObject;
//               break;
//             default:
//               console.log("Unknown server message: ", msgJson);
//               break;
//           }

//           break;
//         default:
//           console.log("default message: ", msgJson);
//           break;
//       }
//     };

//     socket.onerror = (error) => {
//       console.log("App.js: Socket error: ", error);
//     };

//     return socket;
//   }
export async function selfInfo() {
  try {
    const response = await axios.get(
      "http://localhost:8000/userInfo/self",
      {
        withCredentials: true,
      }
    );
    console.log("Self info called \n SelfInfo:",response.data[0])
    if (response.data == "You need to log in first!") {
      router.push('/');
    }
    app.config.globalProperties.selfId = response.data[0].Id;
    app.config.globalProperties.userData = response.data[0];
    return response.data[0]
} catch (error) {
  console.error("Error getting selfInfo:", error);
} 
}

app.use(router)

app.config.globalProperties.BackendImg = "http://localhost:8000/backend/"
app.mount('#app')
