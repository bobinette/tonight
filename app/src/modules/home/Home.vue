<template>
  <div class="container">
    <Planning class="col-md-12" id="planning"></Planning>
    <TaskList class="col-md-12"></TaskList>
    <NewTaskInput></NewTaskInput>
  </div>
</template>

<script>
import TaskList from '@/modules/task-list/TaskList';
import Planning from '@/modules/planning/Planning';
import NewTaskInput from '@/modules/new-task/NewTask';

import { LOAD_FILTERS } from '@/modules/task-list/state';

export default {
  components: {
    NewTaskInput,
    Planning,
    TaskList,
  },
  mounted() {
    this.$store
      .dispatch({ type: LOAD_FILTERS, query: this.$route.query })
      .catch();
  },
};
</script>
