import axios from 'axios';

import { createTask } from './actions';
import { TASK_CREATED } from './events';

describe('create task', () => {
  it('should call the create route and trigger TASK CREATED', async () => {
    const task = { title: 'new task' };
    axios.setResponse({ data: task });

    const commit = jest.fn();
    const content = 'new task';

    const res = await createTask({ commit }, { content });
    expect(res).toEqual(task);
    expect(axios.call()).toEqual({
      url: 'http://127.0.0.1:9093/api/tasks',
      body: { content },
    });
    expect(commit).toHaveBeenCalledWith({ type: TASK_CREATED, task });
  });
});
