import axios from 'axios';

import apiUrl from '@/utils/apiUrl';

import { TASK_CREATED } from './events';

export const createTask = ({ commit }, { content }) =>
  axios
    .post(`${apiUrl}/api/tasks`, { content })
    .then(response => {
      const task = response.data;
      commit({ type: TASK_CREATED, task });
      return task;
    })
    .catch(err => {
      console.log(err);
      throw err;
    });
