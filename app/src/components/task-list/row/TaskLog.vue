<template>
  <ul class="progress-tracker progress-tracker--text progress-tracker--vertical">
    <li v-for="log in logs" class="progress-step is-complete ">
      <span class="progress-marker" :class="markerClass(log)">
        <i class="ProgressIcon" :class="markerIcon(log)"></i>
      </span>
      <span class="progress-text">
        <small class="text-muted"><em>{{ formatDate(log.createdAt) }}</em></small>
        <div v-if="log.description" v-html="markdown(log.description)"></div>
        <div v-else><em class="text-muted">(No description for this step)</em></div>
      </span>
    </li>
    <li class="progress-step" v-if="canAddLog">
      <span class="progress-marker bg-success no-bottom-padding"><i class="fa fa-plus"></i></span>
      <span class="progress-text">
        <textarea
          v-autosize="newLog"
          v-model="newLog"
          @keydown.enter="addLog"
          placeholder="Type here to add a new step..."
          rows="1"
        >
        </textarea>
      </span>
    </li>
  </ul>
</template>

<script>
import moment from 'moment';
import remark from 'remark';
import html from 'remark-html';

export default {
  props: {
    logs: {
      type: Array,
      required: true,
    },
    canAddLog: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      newLog: '',
    };
  },
  methods: {
    addLog(evt) {
      if (evt.shiftKey || !this.newLog) {
        return;
      }
      evt.preventDefault();

      this.$emit('addLog', this.newLog);
    },
    formatDate(date) {
      const deadline = moment(date);
      return deadline.fromNow();
    },
    markdown(text) {
      const md = remark()
        .use(html, { sanitize: true })
        .processSync(text);
      return md.contents;
    },
    markerIcon(log) {
      if (log.type === 'PROGRESS' && log.completion === 100) {
        return ['fa fa-check'];
      }

      return {
        PROGRESS: ['inner-circle'],
        COMMENT: ['fa fa-comment'], // Not used yet
        START: ['fa fa-flag-checkered'],
        PAUSE: ['fa fa-coffee'],
        WONT_DO: ['fa fa-times'],
        DURATION: ['fa fa-clock-o marker-padding-top'],
      }[log.type];
    },
    markerClass(log) {
      if (log.type === 'PROGRESS' && log.completion === 100) {
        return ['success-bg', 'no-bottom-padding'];
      }

      if (log.type === 'WONT_DO') {
        return ['warning-bg'];
      }

      if (log.type === 'PROGRESS') {
        return ['no-bottom-padding'];
      }

      return [];
    },
  },
};
</script>

<style lang="scss" scoped>
@import 'style/_variables';

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

.marker-padding-top {
  padding-top: 2px;
}

textarea {
  background-color: transparent;

  width: 100%;
  max-height: 250px;
  border: none;

  &:focus {
    outline: none;
  }
}
</style>
