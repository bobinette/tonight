<template>
  <li class="list-group-item TaskRow flex-column align-items-start" v-click-outside="hideAll">
    <div v-if="!editMode" class="w-100">
      <div class="flex flex-align-center flex-space-between w-100 RowHeader" @click.stop="open">
        <span class="flex flex-align-center w-100">
          <h6>{{ task.title }}</h6>
          <span class="badge badge-pill badge-danger RowPriority">{{ priority }}</span>
        </span>
        <span class="flex flex-align-center Actions">
          <button class="btn btn-link" @click.stop="deleteTask">
            <i class="fa fa-trash"></i>
          </button>
          <button class="btn btn-link" @click.stop="switchToEditMode">
            <i class="fa fa-pencil"></i>
          </button>
          <button class="btn btn-link" @click.stop="open">
            <i class="fa" :class="{'fa-chevron-down': !isOpen, 'fa-chevron-up': isOpen}"></i>
          </button>
        </span>
      </div>
      <div>
        <span class="badge badge-primary Tag" v-for="tag in task.tags">#{{ tag }}</span>
      </div>
      <div v-if="isOpen">
        <div v-if="task.description" class="Smaller">
          <span v-html="task.description"></span>
        </div>
        <div v-if="task.dependencies && task.dependencies.length" class="Smaller">
          <ul>
            <li v-for="dep in task.dependencies">
              <i class="fa" :class="{
                'fa-check': dep.done,
                success: dep.done,
                'fa-times': !dep.done,
                danger: !dep.done,
              }"></i>
              {{ dep.title }}
            </li>
          </ul>
        </div>
        <div class="flex flex-align-center flex-space-between">
          <div class="flex flex-align-center flex-space-between">
            <div class="text-muted RowDetail" v-if="task.duration">
              <i class="fa fa-clock-o"></i>
              {{ task.duration }}
              &#9679;
              created {{ formatDate(task.createdAt) }}
              <!--
              <span v-if="task.createdAt !== task.updatedAt">
                &#9679;
                updated {{ formatDate(task.updatedAt) }}
              </span>
              -->
            </div>
            <div class="text-muted RowDetail" v-if="task.deadline">
              <i class="fa fa-calendar"></i>
              <em>{{ formattedDeadline }}</em>
            </div>
          </div>
        </div>
        <ul class="progress-tracker progress-tracker--text progress-tracker--vertical">
          <li v-for="log in task.log" class="progress-step is-complete ">
            <span class="progress-marker" :class="markerClass(log)">
              <i class="ProgressIcon" :class="markerIcon(log)"></i>
            </span>
            <span class="progress-text">
              <small class="text-muted"><em>{{ formatDate(log.createdAt) }}</em></small>
              <div v-if="log.description">{{ log.description }}</div>
              <div v-else><em class="text-muted">(No description for this step)</em></div>
            </span>
          </li>
          <li class="progress-step">
            <span class="progress-marker bg-success no-bottom-padding"><i class="fa fa-plus"></i></span>
            <span class="progress-text">
              <textarea
                v-autosize="log"
                v-model="log"
                @keydown.enter="addLog"
                placeholder="Add a new step..."
                rows="1"
              >
              </textarea>
            </span>
          </li>
        </ul>
      </div>
    </div>
    <div v-else class="w-100">
      <textarea
        v-autosize="raw"
        v-focus="editMode"
        v-model="raw"
        @keydown.enter="edit"
        @keydown.esc="editMode = false"
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

// import ProgressTracker from 'vue-bulma-progress-tracker';

import moment from 'moment';

import { formatRaw } from '@/utils/formats';

import { LOG_FOR_TASK, UPDATE_TASK, DELETE_TASK } from '../state';

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
      isOpen: false,
      raw: '',
    };
  },
  computed: {
    priority() {
      return this.task.priority > 0 ? this.task.priority.toString() : '';
    },
    formattedDeadline() {
      const deadline = moment(this.task.deadline);
      return deadline.format('YYYY-MM-DD');
    },
  },
  methods: {
    hideAll() {
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

      this.raw = formatRaw(this.task);
      this.editMode = true;
    },
    open(evt) {
      evt.preventDefault();

      this.isOpen = !this.isOpen;
    },
    formatDate(date) {
      const deadline = moment(date);
      return deadline.fromNow();
    },
    markerIcon(log) {
      if (log.completion === 100) {
        return ['fa fa-check'];
      }

      return {
        COMPLETION: ['inner-circle'],
        START: ['fa fa-flag-checkered'],
        PAUSE: ['fa fa-coffee'],
        WONT_DO: ['fa fa-times'],
      }[log.type];
    },
    markerClass(log) {
      if (log.completion === 100) {
        return ['success-bg'];
      }

      if (log.type === 'WONT_DO') {
        return ['warning-bg'];
      }

      if (log.type === 'COMPLETION') {
        return ['no-bottom-padding'];
      }

      return [];
    },
    deleteTask() {
      this.$store.dispatch({ type: DELETE_TASK, taskId: this.task.id }).catch();
    },
  },
  // Directives
  directives: {
    ClickOutside,
    focus,
    // ProgressTracker,
    // StepItem,
  },
};
</script>

<style lang="scss" scoped>
@import 'style/_variables';

.RowHeader {
  cursor: pointer;
}

.RowPriority {
  margin-left: 0.5rem;
}

.RowDetail:not(:first-child) {
  margin-left: 0.5rem;
}

.Smaller {
  font-size: 0.9rem;
  margin: 0.5rem 0;
}

.Tag:not(:last-child) {
  margin-right: 0.2rem;
}

.Tag.badge {
  // font-weight: normal;
  padding: 0.25rem 0.5rem;
}

.Actions {
  margin-left: 0.3rem;

  .btn.btn-link {
    color: lighten($gray-light, 20);

    &:hover {
      color: $body-color;
    }
  }
}

.ProgressIcon {
  margin: auto;
}

.inner-circle {
  background: $body-bg;
  width: 80%;
  height: 80%;
  border-radius: 50%;
}

.progress-tracker {
  margin-top: 1rem;

  .progress-marker {
    right: - $marker-size;

    &.no-bottom-padding {
      padding-bottom: 0;
    }

    &.success-bg {
      background-color: $brand-success;
    }

    &.warning-bg {
      background-color: $brand-warning;
    }
  }
}

.progress-step:not(:last-child)::after {
  right: - $marker-size - $marker-size-half;
}
</style>
