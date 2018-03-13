import * as actions from './actions';

import * as events from './events';

export default {
  mutations: {
    [events.TASK_CREATED]: () => {}, // Nothing to do
  },
  actions: {
    [events.CREATE_TASK]: actions.createTask,
  },
};
