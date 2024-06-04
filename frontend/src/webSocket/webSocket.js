import {store} from "./store";
import { useEventBus } from './EventBus';

let socket = null;

export function setupWebSocket() {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    console.log("Setting up socket...");
    socket = new WebSocket("ws://localhost:8000/chat/");

    socket.onopen = () => {
      console.log("WebSocket connection successful!");
    };

    socket.onmessage = (msg) => {
      let msgJson = JSON.parse(msg.data);
      const bus = useEventBus();
      bus.value.emit('show-notification', msgJson);
      switch (msgJson.sender) {
        case 'server':
          switch (msgJson.msg) {
            case 'allUsers':
              store.commit('setAllUsers', msgJson.ServerObject);
              break;
            case 'getAllDMs':
              store.commit('setAllMessages', msgJson.ServerObject);
              break;
            case 'getAllNotifications':
              console.log("ws setting all notifications");
              store.commit('setAllNotifications', msgJson.ServerObject);
              break;
            default:
              console.log('Unknown server message: ', msgJson);
              break;
          }
          break;
        default:
          console.log('default message: ', msgJson);
          break;
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error: ", error);
    };
  }

  return socket;
}
export function closeWebSocket() {
  console.log("Closing websocket!")
  socket.close()
}
export function getWebSocket() {
  return socket;
}

export function sendWebSocketMessage(message, messageType) {
  console.log("webSocket.js sendWebSocketMessage message: ", message);
  if (socket && socket.readyState === WebSocket.OPEN) {
    var messageWrapper = {
      messageType: messageType,
      body: message,
    };
    socket.send(JSON.stringify(messageWrapper));
    console.log("WebSocket message sent")
  }
}