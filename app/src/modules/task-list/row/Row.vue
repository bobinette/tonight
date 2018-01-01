<template>
  <li class="list-group-item TaskRow flex-column align-items-start" v-click-outside="hideAll">
    <div v-if="!editMode" class="w-100" @click="showLog">
      <div class="flex flex-align-center flex-space-between w-100">
        <span>
          {{ task.title }}
          <span class="badge badge-pill badge-danger Priority">{{ priority }}</span>
        </span>
        <span class="Actions">
          <button class="btn btn-link" @click="switchToEditMode">
            <i class="fa fa-pencil"></i>
          </button>
        </span>
      </div>
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
    </div>
    <div v-else class="w-100">
      <textarea
        v-autosize="raw"
        v-focus="editMode"
        v-model="raw"
        @keydown.enter="edit"
        placeholder="Edit the task..."
        rows="1"
      >
      </textarea>
    </div>
  </li>
</template>

<script>
import ClickOutside from 'vue-click-outside';
import { focus } from 'vue-focus';

import { formatRaw } from '@/utils/formats';

import { LOG_FOR_TASK, UPDATE_TASK } from '../state';

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
      editMode: false,
      raw: '',
    };
  },
  computed: {
    priority() {
      return this.task.priority > 0 ? this.task.priority.toString() : '';
    },
  },
  methods: {
    showLog() {
      this.logInputVisible = true;
    },
    hideAll() {
      this.logInputVisible = false;
      this.editMode = false;
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
    edit(evt) {
      if (evt.shiftKey) {
        return;
      }
      evt.preventDefault();

      this.$store
        .dispatch({
          type: UPDATE_TASK,
          taskId: this.task.id,
          content: this.raw,
        })
        .then(() => {
          this.editMode = false;
        })
        .catch();
    },
    switchToEditMode(evt) {
      evt.preventDefault();
      evt.stopPropagation(); // Do not open the log input

      this.raw = formatRaw(this.task);
      this.editMode = true;
    },
  },
  // Directives
  directives: {
    ClickOutside,
    focus,
  },
};
</script>

<style lang="scss" scoped>
@import 'style/_variables';

.TaskRow {
  cursor: pointer;
}

.badge {
  margin-left: 0.5rem;
}

.Actions {
  .btn.btn-link {
    color: $gray-light;

    &:hover {
      color: $body-color;
    }
  }
}
</style>
