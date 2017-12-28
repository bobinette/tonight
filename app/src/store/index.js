import userState, { plugins as userPlugins } from '@/modules/user/state';
import taskListState, {
  plugins as taskListPlugins,
} from '@/modules/task-list/state';

const store = {
  modules: {
    user: userState,
    tasks: taskListState,
  },

  // Plugins cannot be defined in modules
  plugins: [...userPlugins, ...taskListPlugins],

  // Raises errors when the state is mutated outside a mutation function
  // Expensive =< do not run in production
  strict: process.NODE_ENV !== 'production',
};

export default store;
