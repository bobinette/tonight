import axios from 'axios';
import Cookies from 'js-cookie';

// Cookie
export const LOAD_COOKIE = 'LOAD_COOKIE';
export const COOKIE_LOADED = 'COOKIE_LOADED';

// User
export const LOAD_USER = 'LOAD_USER';
export const USER_LOADED = 'USER_LOADED';

export const plugins = [
  store =>
    store.subscribe(mutation => {
      if (mutation.type !== COOKIE_LOADED) {
        return;
      }

      store.dispatch({ type: LOAD_USER }).catch(); // No need to handle the error here
    }),
];

export default {
  state: {
    user: {
      loaded: false,
      id: 0,
      name: '',
    },
    cookie: {
      loaded: false,
      token: null,
    },
  },
  getters: {
    accessToken({ cookie }) {
      return cookie.token;
    },
    username({ user }) {
      return user.name;
    },
  },
  mutations: {
    [COOKIE_LOADED]: (state, { token }) => {
      state.cookie = {
        loaded: true,
        token,
      };
    },
    [USER_LOADED]: (state, { id, name }) => {
      state.user = {
        loaded: true,
        id,
        name,
      };
    },
  },
  actions: {
    [LOAD_COOKIE]: context => {
      const token = Cookies.get('access_token');
      context.commit(COOKIE_LOADED, { token });
    },
    [LOAD_USER]: context =>
      axios
        .get('http://127.0.0.1:9090/api/me')
        .then(response => {
          const { id, name } = response.data;
          context.commit(USER_LOADED, { id, name });
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
  },
};
