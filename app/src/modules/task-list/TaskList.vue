<template>
  <ul class="list-group">
    <!-- Rows -->
    <Row v-for="task in tasks" :key="task.id" :task="task"></Row>

    <!-- New task input -->
    <li class="list-group-item">
      <textarea
        ref="createTaskInput"
        v-autosize="newTaskContent"
        v-model="newTaskContent"
        @keydown.enter="createTask"
        placeholder="Create a new task..."
        rows="1"
      >
      </textarea>
    </li>
  </ul>
</template>

<script >
import Row from './row/Row';

import { CREATE_TASK, UPDATE_NEW_TASK_CONTENT } from './state';

export default {
  name: 'task-list',
  data() {
    return {
      newTaskContent: '',
    };
  },
  computed: {
    tasks() {
      return this.$store.getters.tasks;
    },
  },
  methods: {
    updateNewTaskContent(evt) {
      return this.$store.commit({
        type: UPDATE_NEW_TASK_CONTENT,
        content: evt.target.value,
      });
    },
    createTask(evt) {
      if (evt.shiftKey) {
        return;
      }

      evt.preventDefault();
      this.$store
        .dispatch({ type: CREATE_TASK, content: this.newTaskContent })
        .then(() => {
          this.newTaskContent = '';
        })
        .catch(err => console.log(err));
    },
  },
  components: {
    Row,
  },
};
</script>

<style lang="scss">
textarea {
  width: 100%;
  max-height: 250px;
  border: none;

  &:focus {
    outline: none;
  }
}
</style>
