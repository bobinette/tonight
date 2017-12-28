<template>
  <li class="list-group-item TaskRow flex-column align-items-start" @click="showLog" v-click-outside="hideLog">
    {{ task.title }}
    <textarea
      v-if="logInputVisible"
      v-autosize="log"
      v-focus="logInputVisible"
      v-model="log"
      @keydown.enter="addLog"
      placeholder="Add a new step..."
      rows="1"
    >
    </textarea>
  </li>
</template>

<script>
import ClickOutside from 'vue-click-outside';
import { focus } from 'vue-focus';

import { LOG_FOR_TASK } from '../state';

export default {
  name: 'task-row',
  props: {
    task: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      log: '',
      logInputVisible: false,
    };
  },
  methods: {
    showLog() {
      this.logInputVisible = true;
    },
    hideLog() {
      this.logInputVisible = false;
    },
    addLog(evt) {
      if (evt.shiftKey) {
        return;
      }
      evt.preventDefault();

      this.$store
        .dispatch({
          type: LOG_FOR_TASK,
          taskId: this.task.id,
          log: this.log,
        })
        .then(() => {
          this.log = '';
        })
        .catch();
    },
  },
  // Directives
  directives: {
    ClickOutside,
    focus,
  },
};
</script>

<style lang="scss">
.TaskRow {
  cursor: pointer;
}
</style>
