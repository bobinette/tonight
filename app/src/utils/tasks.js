export const completion = task => {
  if (!task.log) {
    return 0;
  }

  let c = 0;
  task.log.forEach(log => {
    if (log.completion > c) {
      c = log.completion;
    }
  });
  return c;
};

export const isDone = task =>
  task.log && task.log.findIndex(l => l.completion === 100) >= 0;

export const isWontDo = task =>
  task.log && task.log.findIndex(l => l.type === 'WONT_DO') >= 0;

export const isPending = task => !isDone(task) && !isWontDo(task);

export const isWorkedOn = task => {
  if (!task.log || !isPending(task)) {
    return false;
  }

  const lastWorkflowStep = task.log
    .slice() // Copy the list so reverse does not mutate the state
    .reverse()
    .find(l => l.type === 'START' || l.type === 'PAUSE');
  return lastWorkflowStep && lastWorkflowStep.type === 'START';
};

const postponeRegex = /postponed until (\d{4}-\d{2}-\d{2})/;

export const isPostponed = task => {
  const d = postponedUntil(task);
  return d && Date.now() - d < 0;
};

export const postponedUntil = task => {
  if (!task.log) {
    return null;
  }

  const lastWorkflowStep = task.log
    .slice() // Copy the list so reverse does not mutate the state
    .reverse()
    .find(l => l.type === 'POSTPONE');
  if (!lastWorkflowStep) {
    return null;
  }

  const match = lastWorkflowStep.description.match(postponeRegex);
  if (!match) {
    return null;
  }

  return Date.parse(match[1]);
};
