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
