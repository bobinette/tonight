<template>
  <div id="app">
    <Navbar></Navbar>
    <div class="container">
      <Planning class="col-md-12" id="planning"></Planning>

      <TaskList class="col-md-12"></TaskList>

      <NewTaskInput></NewTaskInput>
    </div>
  </div>
</template>

<script>
// Modules
import TaskList from '@/modules/task-list/TaskList';
import Planning from '@/modules/planning/Planning';
import NewTaskInput from '@/modules/new-task/NewTask';

import { LOAD_COOKIE } from '@/modules/user/state';
import { LOAD_FILTERS } from '@/modules/task-list/state';

// Components
import Navbar from '@/components/navbar/Navbar';

export default {
  name: 'app',
  components: {
    Navbar,
    NewTaskInput,
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
