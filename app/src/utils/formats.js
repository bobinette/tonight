import moment from 'moment';

// Plural

export const plural = (word, n) => (n > 1 ? `${word}s` : word);

// Duration

const second = 1000000000; // 1M
const minute = 60 * second;
const hour = 60 * minute;

export const formatDuration = dur => {
  let d = dur;
  // Duration is expected in ns
  let formatted = '';
  const h = Math.floor(d / hour);
  if (h > 0) {
    formatted = `${h}h`;
  }
  d -= h * hour;

  const m = Math.floor(d / minute);
  if (m > 0) {
    formatted = `${formatted}${m}m`;
  }
  d -= m * minute;

  const s = Math.floor(d / second);
  if (s > 0) {
    formatted = `${formatted}${s}s`;
  }
  d -= s * minute;

  return formatted;
};

// Raw task

export const formatRaw = task => {
  let formatted = task.title;

  // Add priority
  formatted = `${Array(task.priority + 1).join('!')}${formatted}`;

  // Add the description
  if (task.description !== '') {
    formatted = `${formatted}: ${task.description}`;
  }

  // Add tags
  if (task.tags && task.tags.length > 0) {
    const tags = task.tags.map(tag => `#${tag}`);
    formatted = `${formatted} ${tags.join(' ')}`;
  }

  // Add duration
  if (task.duration !== '') {
    formatted = `${formatted} ~${task.duration}`;
  }

  // Add deadline
  if (task.deadline) {
    formatted = `${formatted} >${moment(task.deadline).format('YYYY-MM-DD')}`;
  }

  // Add dependencies
  if (task.dependencies && task.dependencies.length > 0) {
    const deps = task.dependencies.map(dep => dep.id.toString());
    formatted = `${formatted} needs:${deps.join(',')}`;
  }

  return formatted;
};
