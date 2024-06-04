<template>
    <div class="emoji-picker">
      <div class="emoji-grid">
        <div v-for="emoji in emojis" :key="emoji" @click="insertEmoji(emoji)">
          <span v-html="emoji"></span>
        </div>
      </div>
      <button @click="closeEmojiPicker">Close</button>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        emojis: ['😀', '😄', '😆', '😅', '😂', '🤣', '😊', '😍', '😋', '😎', '😇', '🙃'],
      };
    },
    methods: {
      insertEmoji(emoji) {
        // Insert the selected emoji into the chat input field or message
        this.$emit('emoji-selected', emoji);
        this.closeEmojiPicker();
      },
      closeEmojiPicker() {
        // Emit an event to close the emoji picker
        this.$emit('close-picker');
      },
    },
  };
  </script>
  
  <style scoped>
  .emoji-picker {
    position: absolute;
    background: white;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 10px;
    z-index: 100;
  }
  
  .emoji-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
    max-height: 150px;
    overflow-y: auto;
  }
  
  button {
    display: block;
    margin-top: 10px;
  }
  </style>