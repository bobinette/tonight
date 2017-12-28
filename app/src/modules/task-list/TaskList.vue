<template>
  <ul class="list-group">
    <!-- Header -->
    <li class="list-group-item list-group-item-header ListHeader">
      <div class="col-md-1">
        <strong>{{ tasksLength }} {{ plural("task", tasksLength)}}</strong>
      </div>
      <div class="col-md-4 flex flex-align-center SearchInput">
        <span class="fa fa-search"></span>
        <input type="text" class="form-control" :value="q" @input="updateQ" @keydown.enter="search">
      </div>
      <div class="col-md-1 col-end">
        <i class="fa fa-circle-o-notch fa-spin" v-if="loading"></i>
      </div>
    </li>

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
import { plural } from '@/utils/formats';

import Row from './row/Row';

import {
  UPDATE_Q,
  FETCH_TASKS,
  CREATE_TASK,
  UPDATE_NEW_TASK_CONTENT,
} from './state';

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
    tasksLength() {
      return this.tasks && this.tasks.length ? this.tasks.length : 0;
    },
    q() {
      return this.$store.getters.q;
    },
    loading() {
      return this.$store.getters.loading;
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
    plural,
    updateQ(evt) {
      this.$store.commit({ type: UPDATE_Q, q: evt.target.value });
    },
    search() {
      this.$store.dispatch({ type: FETCH_TASKS }).catch();
    },
  },
  components: {
    Row,
  },
};
</script>

<style lang="scss">
@import 'style/_variables';

.SearchInput {
  background: $input-bg;

  border: $input-btn-border-width solid $input-border-color;
  border-radius: $input-border-radius;

  .fa {
    padding-left: $input-padding-x/2;
  }

  input {
    background: transparent;
    border: none;
    padding: $input-padding-y $input-padding-x/2;

    width: 100%;

    &:focus {
      box-shadow: none;
      outline: none;
    }
  }

  &:focus-within {
    border-color: $input-border-focus;
  }
}

textarea {
  width: 100%;
  max-height: 250px;
  border: none;

  &:focus {
    outline: none;
  }
}

.ListHeader {
  div:not(:last-child) {
    margin-right: 1rem;
  }

  .col-md-1,
  .col-md-4 {
    padding: 0;
  }

  .col-end {
    margin-left: auto;
    text-align: right;
  }
}
</style>
