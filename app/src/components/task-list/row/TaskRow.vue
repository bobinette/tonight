<template>
  <li class="list-group-item TaskRow flex-column align-items-start" :class="{ highlight: isWorkedOn }">
    <div v-if="!editMode" class="w-100">
      <div class="flex flex-align-center flex-space-between w-100 RowHeader" @click.stop="open">
        <span class="flex flex-align-center w-100">
          <h6>{{ task.id }}. {{ task.title }}</h6>
          <span class="badge badge-pill badge-danger RowPriority">{{ priority }}</span>
          <div class="Actions PriorityActions" v-if="isPending">
            <button class="btn btn-link btn-sm" @click.stop="incrementPriority" v-if="isPending">
              <i class="fa fa-plus"></i>
            </button>
            /
            <button class="btn btn-link btn-sm" @click.stop="decrementPriority" v-if="isPending">
              <i class="fa fa-minus"></i>
            </button>
          </div>
        </span>
        <span class="flex flex-align-center Actions">
          <span class="badge badge-pill badge-warning" v-if="isPostponed">Postponed: {{ postponedUntil }}</span>
          <button class="btn btn-link" @click.stop="deleteTask" v-if="isPending">
            <i class="fa fa-trash"></i>
          </button>
          <button class="btn btn-link" @click.stop="switchToEditMode" v-if="isPending">
            <i class="fa fa-pencil"></i>
          </button>
          <span class="badge badge-pill badge-success" v-if="isDone">Done</span>
          <span class="badge badge-pill badge-warning" v-if="isWontDo">Won't do</span>
          <button class="btn btn-link" @click.stop="open">
            <i class="fa" :class="{'fa-chevron-down': !isOpen, 'fa-chevron-up': isOpen}"></i>
          </button>
        </span>
      </div>
      <div class="flex flex-align-center flex-space-between">
        <div>
          <span
            v-for="tag in task.tags"
            :key="tag"
            class="dropdown Tag"
            :class="{show: openTag === tag}"
          >
            <button
              class="btn btn-link"
              :class="{ 'custom-tag': customTag(tag) }"
              :style="tagColourStyle(tag)"
              @click.stop="openTagColourInput(tag)"
            >
              #{{ tag }}
            </button>
            <div
              v-click-outside="hideTagColourInput"
              class="dropdown-menu TagColour"
              @click.stop="() => {}"
            >
              <div class="flex flex-align-center">
                <div
                  class="color-preview"
                  :style="lastValidColour ? {'background-color': lastValidColour} : {}"
                ></div>
                <input
                  class="form-control"
                  :class="{ danger: !colourIsValid }"
                  type="text"
                  v-model="tagColour"
                  @keydown.enter="customizeColour"
                  @keydown.esc="hideTagColourInput"
                >
              </div>
            </div>
          </span>
        </div>
        <div class="flex flex-align-center flex-space-between">
          <div class="text-muted small RowDetail">
            <i class="fa fa-clock-o"></i>
            <span v-if="task.duration">
              {{ task.duration }}
            </span>
            <span v-if="task.duration && displayableLogs && displayableLogs.length">&#9679;</span>
            <span v-if="displayableLogs && displayableLogs.length > 0">
              <span v-if="displayableLogs.length === 1">1 comment</span>
              <span v-else>{{ displayableLogs.length }} comments</span>
            </span>
          </div>
        </div>
      </div>
      <div v-if="isOpen" class="Smaller">
        <div v-if="task.description" class="TaskRow__Details">
          <strong>Description: </strong>
          <span v-html="markdown(task.description)"></span>
        </div>
        <div v-if="task.dependencies && task.dependencies.length" class="TaskRow__Details">
          <strong>Needs: </strong>
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
        <TaskLog
          :logs="displayableLogs"
          :canAddLog="isPending"
          @addLog="addLog"
        ></TaskLog>
        <div class="flex flex-align-center w-100 TaskRow__Details TaskRow__Details__Date">
          <div class="text-muted RowDetail">
            created {{ formatDate(task.createdAt) }}
            <span v-if="task.createdAt !== task.updatedAt">
              &#9679;
              updated {{ formatDate(task.updatedAt) }}
            </span>
          </div>
          <div class="text-muted RowDetail" v-if="task.deadline">
            <i class="fa fa-calendar"></i>
            <em>{{ formattedDeadline }}</em>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="w-100 TaskRow__Edit">
      <AutosuggestTextarea
        :autofocus="editMode"
        :value="raw"
        placeholder="Edit the task..."
        @input="input"
        @keydown.enter="edit"
        @keydown.esc="editMode = false"
        rows="1"
        :disabled="busy"
      >
      </AutosuggestTextarea>
      <span>
        <i class="fa fa-circle-o-notch fa-spin" v-if="busy"></i>
      </span>
    </div>
  </li>
</template>

<script>
import ClickOutside from 'vue-click-outside';
import { focus } from 'vue-focus';

import moment from 'moment';
import remark from 'remark';
import html from 'remark-html';

import AutosuggestTextarea from '@/components/autosuggest-textarea/AutosuggestTextarea';
import { CUSTOMIZE_COLOUR } from '@/modules/user/state';

import { formatRaw } from '@/utils/formats';
import {
  isPending,
  isPostponed,
  isDone,
  isWontDo,
  isWorkedOn,
  postponedUntil,
} from '@/utils/tasks';

import {
  LOG_FOR_TASK,
  UPDATE_TASK,
  DELETE_TASK,
} from '@/modules/task-list/state';

import TaskLog from './TaskLog';

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
      busy: false,
      editMode: false,
      isOpen: false,
      raw: '',
      openTag: '',
      tagColour: '',
      lastValidColour: '',
    };
  },
  computed: {
    priority() {
      return this.task.priority > 0 ? this.task.priority.toString() : '';
    },
    isPending() {
      return isPending(this.task);
    },
    isDone() {
      return isDone(this.task);
    },
    isWontDo() {
      return isWontDo(this.task);
    },
    isWorkedOn() {
      return isWorkedOn(this.task);
    },
    isPostponed() {
      return isPostponed(this.task);
    },
    postponedUntil() {
      const d = postponedUntil(this.task);
      if (!d) {
        return null;
      }

      return moment(d).format('L');
    },
    formattedDeadline() {
      const deadline = moment(this.task.deadline);
      return deadline.format('YYYY-MM-DD');
    },
    colourIsValid() {
      return (
        !this.tagColour ||
        this.tagColour.match(/^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/)
      );
    },
    displayableLogs() {
      if (!this.task.log) {
        return [];
      }

      return this.task.log.filter(l => l.description);
    },
  },
  methods: {
    input(value) {
      this.raw = value;
    },
    addLog(log, { success, failure }) {
      this.$store
        .dispatch({
          type: LOG_FOR_TASK,
          taskId: this.task.id,
          log,
        })
        .then(() => {
          success();
        })
        .catch(() => {
          failure();
        });
    },
    edit(evt) {
      if (evt.shiftKey) {
        return;
      }
      this.busy = true;
      evt.preventDefault();

      this.$store
        .dispatch({
          type: UPDATE_TASK,
          taskId: this.task.id,
          content: this.raw,
        })
        .then(() => {
          this.busy = false;
          this.editMode = false;
        })
        .catch(() => {
          this.busy = false;
        });
    },
    incrementPriority() {
      const raw = `!${formatRaw(this.task)}`;
      this.$store
        .dispatch({
          type: UPDATE_TASK,
          taskId: this.task.id,
          content: raw,
        })
        .catch();
    },
    decrementPriority() {
      let raw = formatRaw(this.task);
      if (raw && raw[0] === '!') {
        raw = raw.substr(1);
      }
      this.$store
        .dispatch({
          type: UPDATE_TASK,
          taskId: this.task.id,
          content: raw,
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
    deleteTask() {
      const r = confirm(`Do you really want to delete the task:\n${this.task.title}`);
      if (r) {
        this.$store
          .dispatch({ type: DELETE_TASK, taskId: this.task.id })
          .catch();
      }
    },
    markdown(text) {
      const md = remark()
        .use(html, { sanitize: true })
        .processSync(text);
      return md.contents;
    },
    openTagColourInput(tag) {
      this.tagColour = '';
      this.lastValidColour = '';
      this.openTag = tag;

      const tagColour = this.$store.state.user.user.tagColours[tag];
      if (tagColour) {
        this.tagColour = tagColour;
      }
    },
    hideTagColourInput() {
      this.openTag = '';
      this.tagColour = '';
      this.lastValidColour = '';
    },
    customizeColour() {
      this.$store
        .dispatch({
          type: CUSTOMIZE_COLOUR,
          tag: this.openTag,
          colour: this.tagColour,
        })
        .catch();
    },
    customTag(tag) {
      return !!this.$store.state.user.user.tagColours[tag];
    },
    tagColourStyle(tag) {
      const tagColour = this.$store.state.user.user.tagColours[tag];
      return tagColour ? { 'background-color': tagColour } : {};
    },
  },
  watch: {
    tagColour() {
      if (this.tagColour.match(/^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/)) {
        this.lastValidColour = this.tagColour;
      }
    },
  },
  // Directives
  directives: {
    ClickOutside,
    focus,
    // ProgressTracker,
    // StepItem,
  },
  components: {
    AutosuggestTextarea,
    TaskLog,
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
}

.Tag .badge {
  // font-weight: normal;
  padding: 0.25rem 0.5rem;

  &:not(:last-child) {
    margin-right: 0.2rem;
  }
}

.Actions {
  margin-left: 0.3rem;

  > span,
  > button {
    margin: 0 0.1rem;
  }

  .btn.btn-link {
    color: lighten($gray-light, 20);

    &:hover {
      color: $body-color;
    }
  }
}

.PriorityActions {
  font-size: 0.8rem;

  button.btn {
    margin: 0;
    border: none;
  }

  .btn .fa {
    font-size: 0.8rem;
  }
}

.highlight {
  background-color: lighten($brand-primary, 55);
}

.Tag .btn.btn-link {
  cursor: pointer;
  font-size: 0.8rem;
  font-weight: bold;
}

.Tag .dropdown-menu {
  top: calc(100% + 10px);
}

.TagColour {
  padding: 0.5rem;
}

.color-preview {
  flex: 0 0 24px;
  height: 24px;
  border-radius: $input-border-radius/2;
  background-color: $brand-primary;
  margin-right: 0.2rem;
}

.custom-tag {
  color: $body-bg;
  padding: 0.2rem;
  border-radius: $border-radius;
}

.TaskRow__Details__Date {
  margin-top: 0.5rem;
}

.TaskRow__Details__Date {
  margin-top: 0.5rem;
  font-size: 0.8rem;
}

.TaskRow__Edit {
  display: flex;

  > .AutosuggestTextarea {
    flex: 1;
  }
}
</style>
