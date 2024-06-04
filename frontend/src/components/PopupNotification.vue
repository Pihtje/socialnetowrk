<template>
    <div class="notification">
        {{ notification.desc }} ASSDFSDFSDFSDFSDFSDFSDF
    </div>
  </template>
  
  <script>
  export default {
    props: ['notificationMessage'],
  };
  </script>
  
  <style scoped>
  .notification {
    position: fixed;
    top: 10px;
    right: 10px;
    padding: 10px;
    background-color: #333;
    color: #fff;
    border-radius: 5px;
    box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    animation: slide-in 0.5s forwards, slide-out 0.5s 2s forwards;
  }
  
  @keyframes slide-in {
    0% {
      transform: translateX(100%);
      opacity: 0;
    }
    100% {
      transform: translateX(0);
      opacity: 1;
    }
  }
  
  @keyframes slide-out {
    0% {
      transform: translateX(0);
      opacity: 1;
    }
    100% {
      transform: translateX(100%);
      opacity: 0;
    }
  }
  </style>