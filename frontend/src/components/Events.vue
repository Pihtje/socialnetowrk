<template>
  <div class="eventContainer">
    <h2>Group Events</h2>
    <ul>
      <li v-for="event in events" :key="event.id">
        {{ event.name }} - {{ event.description }}
      </li>
    </ul>
    <button @click="showCreateEvent">Add New Event</button>
  </div>
</template>

<script>
export default {
  props: {
    events: {
      type: Array,
      required: true,
    }
  },
  methods: {
    showCreateEvent() {
      this.$router.push({ name: 'CreateGroupEventView' });
    }
  }
};
</script>

<style scoped>
.eventContainer {
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

.eventContainer ul {
  list-style: none;
  padding: 0;
}

.eventContainer li {
  margin-bottom: 0.5rem;
}

.eventContainer button {
  background-color: #00bd7e;
  border: none;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:hover {
    background-color: #008d65;
  }
}
</style>