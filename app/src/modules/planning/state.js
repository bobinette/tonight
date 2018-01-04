import axios from 'axios';

import { COOKIE_LOADED } from '@/modules/user/state';
import { TASK_UPDATED, TASK_DELETED } from '@/modules/task-list/state';

// Fetch
export const FETCH_PLANNING = 'FETCH_PLANNING';
export const PLANNING_RECEIVED = 'PLANNING_RECEIVED';

// Start
export const START_PLANNING = 'START_PLANNING';

// Dismiss
export const DISMISS_PLANNING = 'DISMISS_PLANNING';

// Plugins
export const plugins = [
  store =>
    store.subscribe(mutation => {
      const types = [COOKIE_LOADED, TASK_UPDATED, TASK_DELETED];
      if (!types.find(t => t === mutation.type)) {
        return;
      }

      store.dispatch({ type: FETCH_PLANNING }).catch(() => {});
    }),
];

// State module
export default {
  state: {
    planning: null,
    loading: false,
  },
  getters: {},
  mutations: {
    [PLANNING_RECEIVED]: (state, { planning }) => {
      state.planning = planning;
    },
  },
  actions: {
    [FETCH_PLANNING]: context =>
      axios
        .get('http://127.0.0.1:9090/api/planning')
        .then(response => {
          const planning = response.data;
          context.commit({ type: PLANNING_RECEIVED, planning });
          return planning;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [START_PLANNING]: (context, { input }) =>
      axios
        .post('http://127.0.0.1:9090/api/planning', { input })
        .then(response => {
          const planning = response.data;
          context.commit({ type: PLANNING_RECEIVED, planning });
          return planning;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [DISMISS_PLANNING]: context =>
      axios
        .delete('http://127.0.0.1:9090/api/planning')
        .then(() => {
          context.commit({ type: PLANNING_RECEIVED, planning: null });
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
  },
};
