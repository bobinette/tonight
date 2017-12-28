import axios from 'axios';

import { COOKIE_LOADED } from '@/modules/user/state';

import { isDone } from '@/utils/tasks';

// Tasks
// -- List
export const UPDATE_Q = 'UPDATE_Q';
export const FETCH_TASKS = 'FETCH_TASKS';
export const TASK_FETCHING_STARTED = 'TASK_FETCHING_STARTED';
export const TASKS_RECEIVED = 'TASKS_RECEIVED';

// -- Create
export const CREATE_TASK = 'CREATE_TASK';
export const TASK_CREATED = 'TASK_CREATED';

// -- Log
export const LOG_FOR_TASK = 'LOG_FOR_TASK';

// -- Update
export const UPDATE_TASK = 'UPDATE_TASK';

// Plugins
export const plugins = [
  store =>
    store.subscribe(mutation => {
      if (mutation.type !== COOKIE_LOADED) {
        return;
      }

      store.dispatch({ type: FETCH_TASKS }).catch();
    }),
];

// State module
export default {
  state: {
    q: '',
    tasks: [],
    loading: false,
  },
  getters: {
    tasks: ({ tasks }) => tasks,
    q: ({ q }) => q,
    loading: ({ loading }) => loading,
  },
  mutations: {
    // SEARCH / LIST
    [UPDATE_Q]: (state, { q }) => {
      state.q = q;
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
    [UPDATE_TASK]: (state, { task }) => {
      const idx = state.tasks.findIndex(t => task.id === t.id);
      if (idx === -1) {
        return;
      }

      state.tasks[idx] = task;
    },
  },
  actions: {
    [FETCH_TASKS]: context => {
      context.commit({ type: TASK_FETCHING_STARTED });
      const q = encodeURIComponent(context.state.q);
      return axios
        .get(`http://127.0.0.1:9090/api/tasks?q=${q}`)
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
          context.commit({ type: UPDATE_TASK, task });
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
  },
};
