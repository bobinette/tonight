<template>
  <div class="container">
    <div v-if="userLoggedIn">
      <Planning class="col-md-12" id="planning"></Planning>
      <TaskList class="col-md-12"></TaskList>
      <NewTaskInput></NewTaskInput>
    </div>
    <div v-if="displayEmptyState" class="text-align-center">
      <button class="btn btn-link" @click="login">Login</button>
      or
      <button class="btn btn-link" @click="login">sign up</button>
      (both via Google)
      to start using Tonight
    </div>
  </div>
</template>

<script>
import TaskList from '@/modules/task-list/TaskList';
import Planning from '@/modules/planning/Planning';
import NewTaskInput from '@/modules/new-task/NewTask';

import { LOGIN } from '@/modules/user/state';

export default {
  components: {
    NewTaskInput,
    Planning,
    TaskList,
  },
  computed: {
    userLoggedIn() {
      return this.$store.state.user.user.id !== 0;
    },
    displayEmptyState() {
      return (
        this.$store.state.user.user.loaded &&
        this.$store.state.user.user.id === 0
      );
    },
  },
  methods: {
    login() {
      this.$store.dispatch({ type: LOGIN }).catch();
    },
  },
};
</script>

<style lang="scss">

</style>
