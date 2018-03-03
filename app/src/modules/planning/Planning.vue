<template>
  <div>
    <TaskList :tasks="planning.tasks" v-if="!!planning">
      <div slot="header" class="PlanningHeader flex flex-space-between">
        <strong>{{ planning.tasks.length}} {{ plural("task", planning.tasks.length)}}</strong>
        <span>
          <i class="fa fa-clock-o"></i>{{ q }}{{ strict }}{{ duration }}
        </span>
        <button class="btn btn-link" @click="dismiss">Dismiss</button>
      </div>
      <li class="list-group-item progress">
        <div class="progress-bar progress-bar-small" role="progressbar" :style='{width: `${completion}%`}'></div>
      </li>
    </TaskList>
    <div class="card EmptyPlanning" v-else>
      <h5>You currently have no planning</h5>
      <div>Start a new planning by entering below how long you want to work:</div>
      <input
        type="text"
        class="form-control"
        placeholder="Duration e.g. 1h..."
        v-model="planningDuration"
        @keydown.13="startPlanning"
      >
    </div>
  </div>
</template>

<script>
import { formatDuration, plural } from '@/utils/formats';

import TaskList from '@/components/task-list/TaskList';

import { completion } from '@/utils/tasks';

import { START_PLANNING, DISMISS_PLANNING } from './events';

export default {
  data() {
    return {
      planningDuration: '',
    };
  },
  computed: {
    planning() {
      return this.$store.state.planning.planning;
    },
    q() {
      return this.planning.q !== '' ? `${this.planning.q} for ` : '';
    },
    strict() {
      return this.planning.strict ? '!' : '';
    },
    duration() {
      return formatDuration(this.planning.duration);
    },
    completion() {
      const c =
        100 *
        this.planning.tasks.reduce(
          (acc, task) => acc + completion(task) / 100,
          0,
        );
      return Math.round(c / this.planning.tasks.length);
    },
  },
  methods: {
    startPlanning() {
      this.$store
        .dispatch({
          type: START_PLANNING,
          input: this.planningDuration,
        })
        .then(() => {
          this.planningDuration = '';
        })
        .catch();
    },
    plural,
    dismiss() {
      this.$store.dispatch({ type: DISMISS_PLANNING }).catch();
    },
  },
  components: {
    TaskList,
  },
};
</script>

<style lang="scss" scoped>
@import 'style/_variables';

#planning {
  margin-bottom: 2rem;
  padding-left: inherit;
  padding-right: inherit;
}

.card {
  padding: 0.75rem 1.25rem;
}

ul.list-group {
  padding-left: 0;
  padding-right: 0;
}

.EmptyPlanning {
  text-align: center;

  input {
    width: 400px;
    margin-left: auto;
    margin-right: auto;
  }
}

.fa {
  margin-right: 0.2rem;
}

button.btn.btn-link {
  color: black;
}

.progress {
  border-radius: 0;

  &.list-group-item {
    padding: 0;
  }

  .progress-bar {
    border-color: $brand-primary;

    &.progress-bar-small {
      height: 8px;
    }
  }
}

.PlanningHeader {
  width: 100%;
}
</style>
