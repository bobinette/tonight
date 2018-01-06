export const isDone = task =>
  task.log && task.log.findIndex(l => l.completion === 100) >= 0;

export const isWontDo = task =>
  task.log && task.log.findIndex(l => l.type === 'WONT_DO') >= 0;

export const isPending = task => !isDone(task) && !isWontDo(task);

export const isWorkedOn = task => {
  if (!task.log || isDone(task) || isWontDo(task)) {
    return false;
  }

  // slice to prevent reversing the log array
  const lastWorkflowStep = task.log
    .slice()
    .reverse()
    .find(l => l.type === 'START' || l.type === 'PAUSE');
  return lastWorkflowStep.type === 'START';
};
