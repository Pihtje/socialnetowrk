<script>
export default {
  props: {
    GroupId: {
      type: Number,
      required: true,
    },
    Creator: {
      type: Array,
      required: true,
    },
    GroupTitle: {
      type: String,
      required: true,
    },
    GroupDesc: {
      type: String,
      required: true,
    },
    Members: {
      type: Array,
      required: true,
    },
    Posts: {
      type: Array,
      required: false,
    },
  },
  methods: {
    goToGroup(groupId){
      console.log("GroupId",groupId)
      this.$router.push({ name: 'SingleGroup', params: { groupId } });  
    },
  },
};
</script>

<template>
    <div @click="goToGroup(GroupId)" class="groupContainer">
      <h2 class="groupTitle">{{ GroupTitle }}</h2>
      <br />
      <p class="green">Created by:</p>
      <p>{{ Creator[0].Nickname.String || Creator[0].FirstName + ' ' + Creator[0].LastName }}</p>
      <br />
      <p class="green">Group description: </p>
      <p>{{ GroupDesc }}</p>
      <br />
      <div>
      <p class="green">Total Members:</p>
      <p> {{ Members.length }}</p>
      </div>
    </div>
</template>
<style scoped>
.groupContainer {
    background-color: #101010;
    border-radius: 1rem;
    box-shadow: 0rem 0rem 1rem 0.01rem #00bd7e;
    width: 100%;
    margin-bottom: 1.5rem;
    position: relative;
}
.groupTitle {
  font-weight: bold;
}
</style>
