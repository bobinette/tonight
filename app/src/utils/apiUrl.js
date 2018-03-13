// URLs
let u = 'https://tonight.bobi.space';
if (process.env.NODE_ENV !== 'production') {
  u = 'http://127.0.0.1:9090';
}

const url = u;
export default url;
