<template>
  <div>
    <ul v-if="!!planning" class="list-group">
      <li class="list-group-item list-group-item-header flex flex-space-between">
        <strong>{{ planning.tasks.length}} {{ plural("task", planning.tasks.length)}}</strong>
        <span>
          <i class="fa fa-clock-o"></i>{{ q }}{{ strict }}{{ duration }}
        </span>
        <button class="btn btn-link" @click="dismiss">Dismiss</button>
      </li>
      <Row v-for="task in planning.tasks" :key="task.id" :task="task"></Row>
    </ul>
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

import Row from '@/modules/task-list/row/Row';

import { START_PLANNING, DISMISS_PLANNING } from './state';

export default {
  data() {
    return {
      planningDuration: '',
    };
  },
  computed: {
    planning() {
      return this.$store.getters.planning;
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
    Row,
  },
};
</script>

<style lang="scss" scoped>
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
</style>
