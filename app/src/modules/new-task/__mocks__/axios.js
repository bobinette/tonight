let response = null;
let _call = null;

const post = (url, body) => {
  _call = { url, body };

  return new Promise(resolve => {
    resolve(response);
  });
};

const setResponse = res => {
  response = res;
};

const call = () => _call;

export default {
  post,
  setResponse,
  call,
};
