import Vue from 'vue';
import Vuex from 'vuex';

import notificationState from '@/modules/notifications/state';

import planningState, {
  plugins as planningPlugins,
} from '@/modules/planning/state';

import taskListState, {
  plugins as taskListPlugins,
} from '@/modules/task-list/state';

import userState, { plugins as userPlugins } from '@/modules/user/state';

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    notifications: notificationState,
    planning: planningState,
    tasks: taskListState,
    user: userState,
  },

  // Plugins cannot be defined in modules
  plugins: [...planningPlugins, ...userPlugins, ...taskListPlugins],

  // Raises errors when the state is mutated outside a mutation function
  // Expensive => do not run in production
  strict: process.NODE_ENV !== 'production',
});

export default store;
