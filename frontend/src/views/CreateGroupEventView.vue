<template>
    <div class="createEventContainer">
      <h2>Create Event</h2>
      
      <label for="eventTitle">Event Title:</label>
      <input class="formField" type="text" id="eventTitle" v-model="eventTitle" placeholder="Event Title">
      
      <label for="eventDescription">Event Description:</label>
      <input class="formField" type="text" id="eventDescription" v-model="eventDescription" placeholder="Event Description">
      
      <label for="eventDateTime">Event Date & Time:</label>
      <input class="formField" type="datetime-local" id="eventDateTime" v-model="eventDateTime">
      
      
  
      <button class="formButton" @click="addEvent">Submit</button>
    </div>
  </template>
  
  
  <script>
  import axios from 'axios';
  
  export default {
    props: ['groupId'],
    
      data() {
          return {
              eventTitle: "",
              eventDescription: "",
              eventDateTime: "", // datetime for the event
              selectedMembers: [], // array of selected member IDs
              members: [], // This can be fetched from API or passed as prop
              userId: "",   
            //   groupId: ""
          };
      },
      methods: {
          async addEvent() {
              try {
              // Convert the user-selected datetime to the expected format for Go
              let eventDate = new Date(this.eventDateTime);
              let formattedDate = `${eventDate.getUTCFullYear()}-${String(eventDate.getUTCMonth() + 1).padStart(2, '0')}-${String(eventDate.getUTCDate()).padStart(2, '0')}T${String(eventDate.getUTCHours()).padStart(2, '0')}:${String(eventDate.getUTCMinutes()).padStart(2, '0')}:00Z`;
                  const response = await axios.post(
                      `http://localhost:8000/createGroupEvent/${this.groupId}`,
                      {
                          title: this.eventTitle,
                          description: this.eventDescription,
                          day_time: formattedDate,
                          members: this.selectedMembers,
                          userId: this.userId,
                          group: this.groupId
                      }, 
                        {
                            withCredentials: true
                        }
                  );
                  console.log("Event Added:", response.data);
              } catch (error) {
                  console.error("Error adding event:", error);
              }
          }
      },
      mounted() {
        console.log("Received groupId:", this.groupId);
          // Fetch or compute userId, groupId, and members list
      }
  };
  </script>
  <style scoped>
  .formField{
    margin-bottom: 1rem;
  }
  .createEventContainer input{
    border-color: #00bd7e;
    color: #00bd7e;
  border: none;
  border-radius: 0.5rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
  }
</style>
  
  
  
  
