export const isDone = task =>
  task.log.findIndex(l => l.completion === 100) >= 0;
