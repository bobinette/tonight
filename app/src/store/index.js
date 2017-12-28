import planningState, {
  plugins as planningPlugins,
} from '@/modules/planning/state';
import userState, { plugins as userPlugins } from '@/modules/user/state';
import taskListState, {
  plugins as taskListPlugins,
} from '@/modules/task-list/state';

const store = {
  modules: {
    planning: planningState,
    user: userState,
    tasks: taskListState,
  },

  // Plugins cannot be defined in modules
  plugins: [...planningPlugins, ...userPlugins, ...taskListPlugins],

  // Raises errors when the state is mutated outside a mutation function
  // Expensive =< do not run in production
  strict: process.NODE_ENV !== 'production',
};

export default store;
