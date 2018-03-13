<template>
  <ul class="list-group">
    <!-- Header -->
    <li class="list-group-item list-group-item-header">
      <slot name="header"></slot>
    </li>

    <!-- Default slot: holds the extra lines before the rows -->
    <slot name="before"></slot>

    <!-- Rows -->
    <TaskRow v-for="task in tasks" :key="task.id" :task="task"></TaskRow>
  </ul>
</template>

<script>
import TaskRow from './row/TaskRow';

export default {
  name: 'task-list',
  props: {
    tasks: {
      type: Array,
      required: true,
    },
  },
  components: {
    TaskRow,
  },
};
</script>
