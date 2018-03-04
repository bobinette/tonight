import axios from 'axios';

import { TASK_CREATED } from './events';

export const createTask = ({ commit }, { content }) =>
  axios
    .post('http://127.0.0.1:9090/api/tasks', { content })
    .then(response => {
      const task = response.data;
      commit({ type: TASK_CREATED, task });
      return task;
    })
    .catch(err => {
      console.log(err);
      throw err;
    });
