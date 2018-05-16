<template>
  <div class="NewTaskInput">
    <div class="NewTaskTextArea" v-if="isOpen" v-click-outside="close">
      <AutosuggestTextarea
        :autofocus="isOpen"
        :value="newTaskContent"
        placeholder="Create a new task..."
        @input="updateContent"
        @keydown.enter="createTask"
        @keydown.esc="close"
        rows="5"
        ref="textarea"
      >
      </AutosuggestTextarea>
      <small class="text-muted NewTaskTextAreaHelp">
        Press enter to create <i class="fa fa-level-down fa-rotate-90"></i>
      </small>
    </div>
    <button class="btn btn-success NewTaskButton" @click="open">
      <i class="fa fa-plus"></i>
    </button>
  </div>
</template>

<script>
import ClickOutside from 'vue-click-outside';

import AutosuggestTextarea from '@/components/autosuggest-textarea/AutosuggestTextarea';

import { CREATE_TASK } from './events';

export default {
  data() {
    return {
      isOpen: false,
      newTaskContent: '',
    };
  },

  methods: {
    close() {
      this.isOpen = false;
    },
    open() {
      this.isOpen = true;
    },
    updateContent(value) {
      this.newTaskContent = value;
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
    AutosuggestTextarea,
  },

  directives: {
    ClickOutside,
  },
};
</script>

<style scoped lang="scss">
@import 'style/_variables';

.NewTaskInput {
  text-align: right;

  width: 33%;
  position: fixed;
  right: 1rem;
  bottom: 1rem;
}

.NewTaskButton {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;

  cursor: pointer;
}

.NewTaskTextAreaHelp {
  margin-right: 0.3rem;
}

.NewTaskTextArea {
  background: white;

  border: $input-btn-border-width solid $input-border-color;
  border-radius: $input-border-radius;

  padding: 0.2rem;
  margin-bottom: 0.2rem;

  &:focus-within {
    border-color: $input-border-focus;
  }

  textarea {
    background: transparent;
    border: none;
    max-height: 50vh;

    width: 100%;
    max-height: 250px;
    border: none;

    &:focus {
      outline: none;
    }
  }
}
</style>
