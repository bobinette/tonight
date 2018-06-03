<template>
  <div class="AutosuggestTextarea">
    <textarea
      v-autosize="value"
      v-focus="autofocus"
      :value="value"
      :placeholder="placeholder"
      ref="textarea"
      @keydown="keydown"
      @input="input"
      :disabled="disabled"
    >
    </textarea>
    <div
      v-if="showSuggestions && suggestions"
      class="AutosuggestTextarea__SuggestionsContainer"
      :style="{ top: suggestionsPos.top + 'px', left: suggestionsPos.left + 'px' }"
      ref="suggestionsContainer"
    >
      <ul class="AutosuggestTextarea__Suggestions">
        <li
          v-for="(suggestion, idx) in suggestions"
          class="AutosuggestTextarea__Suggestion"
          :key="suggestion"
          :class="{ 'AutosuggestTextarea__Suggestion--active': idx === selectedSuggestion }"
        >
          {{ suggestion }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { focus } from 'vue-focus';

import getCaretCoordinates from 'textarea-caret';

import debounce from 'lodash.debounce';

import axios from 'axios';

import apiUrl from '@/utils/apiUrl';

const TAG_REGEX = /^(?:\B)(#(?:\w|-|:)+)\b$/;

const circle = (idx, up, len) => {
  let newIdx = idx + (up ? 1 : -1);
  if (newIdx >= len) {
    newIdx = 0;
  } else if (newIdx < 0) {
    newIdx = len - 1;
  }
  return newIdx;
};

export default {
  props: {
    value: { type: String, required: true },
    placeholder: { type: String, default: '' },
    autofocus: { type: Boolean, default: false },
    disabled: { type: Boolean, default: false },
  },

  data() {
    return {
      hideSuggestions: false,
      suggestions: [],
      query: '',
      selectedSuggestion: -1,
      suggestionsPos: { top: 0, left: 0, tagStart: 0 },
      justSelected: false,
    };
  },

  computed: {
    showSuggestions() {
      return this.suggestions.length > 0 && !this.hideSuggestions;
    },
  },

  methods: {
    input(evt) {
      this.$emit('input', evt.target.value);
    },
    keydown(evt) {
      const { keyCode } = evt;
      switch (keyCode) {
        case 40: // ArrowDown
          if (this.showSuggestions) {
            this.selectedSuggestion = circle(
              this.selectedSuggestion,
              true,
              this.suggestions.length,
            );

            const self = this;
            const previousActive = self.$el.querySelector('.AutosuggestTextarea__Suggestion--active');
            Vue.nextTick(() => {
              const container = self.$el.querySelector('.AutosuggestTextarea__SuggestionsContainer');
              const active = self.$el.querySelector('.AutosuggestTextarea__Suggestion--active');

              if (container.scrollTop > active.offsetTop) {
                container.scrollTop = active.offsetTop;
              } else if (
                container.scrollTop + container.offsetHeight <
                active.offsetTop + active.offsetHeight
              ) {
                container.scrollTop += previousActive.offsetHeight;
              }
            });

            evt.preventDefault();
            return;
          }
          break;
        case 38: // ArrowUp
          if (this.showSuggestions) {
            this.selectedSuggestion = circle(
              this.selectedSuggestion,
              false,
              this.suggestions.length,
            );

            const self = this;
            Vue.nextTick(() => {
              const container = self.$el.querySelector('.AutosuggestTextarea__SuggestionsContainer');
              const active = self.$el.querySelector('.AutosuggestTextarea__Suggestion--active');

              if (container.scrollTop > active.offsetTop) {
                container.scrollTop = active.offsetTop;
              } else if (
                container.scrollTop + container.offsetHeight <
                active.offsetTop + active.offsetHeight
              ) {
                container.scrollTop = active.offsetTop;
              }
            });

            evt.preventDefault();
            return;
          }
          break;
        case 27: // Escape
          if (this.showSuggestions) {
            this.hideSuggestions = true;
            return;
          }
          break;
        case 13:
          if (
            this.showSuggestions &&
            this.selectedSuggestion >= 0 &&
            this.selectedSuggestion < this.suggestions.length
          ) {
            const suggestion = this.suggestions[this.selectedSuggestion];
            const value = `${this.value.substring(
              0,
              this.suggestionsPos.tagStart,
            )}#${suggestion}${this.value.substring(this.suggestionsPos.tagStart + this.query.length)}`;

            this.suggestions = [];
            this.justSelected = true;
            evt.preventDefault();

            this.$emit('input', value);
            return;
          }
          break;
        default:
          break;
      }

      this.$emit('keydown', evt);
    },
  },

  watch: {
    value: debounce(async function (value) {
      const tagStart = value.lastIndexOf(
        '#',
        this.$refs.textarea.selectionStart,
      );
      const self = this;

      if (this.justSelected) {
        this.justSelected = false;
        this.suggestions = [];
        Vue.nextTick(() => {
          self.hideSuggestions = false;
        });
        return;
      }

      if (tagStart < 0) {
        this.suggestions = [];
        Vue.nextTick(() => {
          self.hideSuggestions = false;
        });
        return;
      }

      const substr = value.substr(
        tagStart,
        this.$refs.textarea.selectionStart - tagStart,
      );
      const match = substr.match(TAG_REGEX);
      if (!match) {
        this.suggestions = [];
        Vue.nextTick(() => {
          self.hideSuggestions = false;
        });
        return;
      }

      const { top, left, height } = getCaretCoordinates(
        this.$refs.textarea,
        tagStart,
      );
      this.suggestionsPos = { top: top + height, left, tagStart };

      // tmp
      let suggestions = [];
      try {
        const res = await axios.get(`${apiUrl}/api/tags?q=${encodeURI(match[1].substr(1))}`);
        suggestions = res.data.tags;
      } catch (err) {
        console.log(err);
      }
      this.query = match[1];
      this.suggestions = suggestions;

      // Restore the hideSuggestions
      Vue.nextTick(() => {
        self.hideSuggestions = false;
      });
    }, 200),

    showSuggestions(newValue) {
      if (!newValue) {
        this.selectedSuggestion = -1;
      }
    },
  },

  directives: {
    focus,
  },
};
</script>

<style scoped lang="scss">
@import 'style/_variables';

.AutosuggestTextarea {
  position: relative;
  z-index: 1000;
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

.AutosuggestTextarea__SuggestionsContainer {
  position: absolute;
  background: white;

  border: 1px solid $input-border-color;
  border-radius: $border-radius;

  max-height: 115px;
  overflow-y: scroll;
}

.AutosuggestTextarea__Suggestions {
  text-align: left;
  list-style: none;
  margin: 0;
  padding: 0;
}

.AutosuggestTextarea__Suggestion {
  padding: 0.2rem 0.5rem;
  border-bottom: 1px solid $input-border-color;
  height: 32px;

  &:last-child {
    border-bottom: none;
  }

  &--active {
    background: lighten($brand-primary, 55);
  }
}
</style>
