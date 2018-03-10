// URLs
let u = 'https://tonight.bobi.space';
if (
  window.location.hostname.indexOf('localhost') !== -1 ||
  window.location.hostname.indexOf('127.0.0.1') !== -1
) {
  u = 'http://127.0.0.1:9090';
}

const url = u;
export default url;
