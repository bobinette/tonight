import axios from 'axios';
import qs from 'qs';

import Vue from 'vue';

import { COOKIE_LOADED } from '@/modules/user/state';

import { isDone } from '@/utils/tasks';

// Tasks
// -- List
export const UPDATE_Q = 'UPDATE_Q';
export const FETCH_TASKS = 'FETCH_TASKS';
export const TASK_FETCHING_STARTED = 'TASK_FETCHING_STARTED';
export const TASKS_RECEIVED = 'TASKS_RECEIVED';

export const UPDATE_STATUS_FILTER = 'UPDATE_STATUS_FILTER';

// -- Create
export const CREATE_TASK = 'CREATE_TASK';
export const TASK_CREATED = 'TASK_CREATED';

// -- Log
export const LOG_FOR_TASK = 'LOG_FOR_TASK';

// -- Update
export const TASK_UPDATED = 'TASK_UPDATED';
export const UPDATE_TASK = 'UPDATE_TASK';

// -- Delete
export const DELETE_TASK = 'DELETE_TASK';
export const TASK_DELETED = 'TASK_DELETED';

// Plugins
export const plugins = [
  store =>
    store.subscribe(mutation => {
      const types = [
        COOKIE_LOADED,
        TASK_UPDATED,
        TASK_DELETED,
        UPDATE_STATUS_FILTER,
      ];
      if (!types.find(t => t === mutation.type)) {
        return;
      }

      store.dispatch({ type: FETCH_TASKS }).catch(() => {});
    }),
];

// State module
export default {
  state: {
    q: '',
    statuses: [],
    tasks: [],
    loading: false,
  },
  getters: {},
  mutations: {
    // SEARCH / LIST
    [UPDATE_Q]: (state, { q }) => {
      state.q = q;
    },
    [UPDATE_STATUS_FILTER]: (state, { status }) => {
      const idx = state.statuses.findIndex(s => s === status);
      if (idx === -1) {
        state.statuses.push(status);
      } else {
        state.statuses.splice(idx, 1);
      }
    },
    [TASK_FETCHING_STARTED]: state => {
      state.loading = true;
    },
    [TASKS_RECEIVED]: (state, { tasks }) => {
      state.loading = false;
      state.tasks = tasks;
    },
    // CREATE
    [TASK_CREATED]: (state, { task }) => {
      state.tasks.push(task);
      state.newTaskContent = '';
    },
    // UPDATE
    [TASK_UPDATED]: (state, { task }) => {
      const idx = state.tasks.findIndex(t => task.id === t.id);
      if (idx === -1) {
        return;
      }

      // This does not work:
      // state.tasks[idx] = task;
      // I need to use Vue.set:
      Vue.set(state.tasks, idx, task);
    },
    [TASK_DELETED]: () => {}, // Nothing to do
  },
  actions: {
    [FETCH_TASKS]: context => {
      context.commit({ type: TASK_FETCHING_STARTED });
      const { statuses } = context.state;
      const q = encodeURIComponent(context.state.q);
      return axios
        .get(
          `http://127.0.0.1:9090/api/tasks?${qs.stringify(
            { q, statuses },
            { skipNulls: true, indices: false }
          )}`
        )
        .then(response => {
          const { tasks } = response.data;
          context.commit({ type: TASKS_RECEIVED, tasks });
          return tasks;
        })
        .catch(err => {
          console.log(err);
          throw err;
        });
    },
    [CREATE_TASK]: (context, { content }) =>
      axios
        .post('http://127.0.0.1:9090/api/tasks', { content })
        .then(response => {
          const task = response.data;
          context.commit({ type: TASK_CREATED, task });
          return task;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [LOG_FOR_TASK]: (context, { taskId, log }) =>
      axios
        .post(`http://127.0.0.1:9090/api/tasks/${taskId}/log`, { log })
        .then(response => {
          const task = response.data;
          context.commit({ type: TASK_UPDATED, task });
          return task;
        })
        .then(task => {
          if (isDone(task)) {
            context.dispatch({ type: FETCH_TASKS });
          }
          return task;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [UPDATE_TASK]: (context, { taskId, content }) =>
      axios
        .post(`http://127.0.0.1:9090/api/tasks/${taskId}`, { content })
        .then(response => {
          const updatedTask = response.data;
          context.commit({ type: TASK_UPDATED, task: updatedTask });
          return updatedTask;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [DELETE_TASK]: (context, { taskId }) =>
      axios
        .delete(`http://127.0.0.1:9090/api/tasks/${taskId}`)
        .then(() => {
          context.commit({ type: TASK_DELETED, taskId });
          return taskId;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
  },
};
