<template>
  <div id="app">
    <Navbar></Navbar>
    <div class="container">
      <Planning class="col-md-12" id="planning"></Planning>

      <TaskList class="col-md-12"></TaskList>

      <NewTaskInput></NewTaskInput>
    </div>
    <Notifications></Notifications>
  </div>
</template>

<script>
// Modules
import Notifications from '@/modules/notifications/Notifications';
import TaskList from '@/modules/task-list/TaskList';
import Planning from '@/modules/planning/Planning';
import NewTaskInput from '@/modules/new-task/NewTask';

import { LOAD_FILTERS } from '@/modules/task-list/state';
import { LOAD_COOKIE } from '@/modules/user/state';

// Components
import Navbar from '@/components/navbar/Navbar';

export default {
  name: 'app',
  components: {
    Navbar,
    NewTaskInput,
    Notifications,
    Planning,
    TaskList,
  },
  mounted() {
    this.$store.dispatch({ type: LOAD_COOKIE });
    this.$store.dispatch({ type: LOAD_FILTERS, query: window.location.search });
  },
};
</script>

<style lang="scss">
#planning {
  margin-bottom: 2rem;
}
</style>
