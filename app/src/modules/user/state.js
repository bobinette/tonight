import axios from 'axios';

import { NOTIFICATION_FAILURE } from '@/modules/notifications/state';

// User
export const LOGIN = 'LOGIN';
export const LOGOUT = 'LOGOUT';
export const LOAD_USER = 'LOAD_USER';
export const USER_LOADED = 'USER_LOADED';

// Colours
export const CUSTOMIZE_COLOUR = 'CUSTOMIZE_COLOUR';

export const plugins = [];

export default {
  state: {
    user: {
      loaded: false,
      id: 0,
      name: '',
      tagColours: {},
    },
  },
  getters: {
    username({ user }) {
      return user.name;
    },
  },
  mutations: {
    [USER_LOADED]: (state, { id, name, tagColours }) => {
      state.user = {
        loaded: true,
        id,
        name,
        tagColours,
      };
    },
  },
  actions: {
    [LOAD_USER]: context =>
      axios
        .get('http://127.0.0.1:9090/api/me')
        .then(response => {
          context.commit(USER_LOADED, response.data);
          return response.data;
        })
        .catch(err => {
          console.log(err);
          throw err;
        }),
    [CUSTOMIZE_COLOUR]: (context, { tag, colour }) =>
      axios
        .post(`http://127.0.0.1:9090/api/tags/${tag}`, { colour })
        .then(response => {
          context.commit(USER_LOADED, response.data);
          return response.data;
        })
        .catch(err => {
          let message = err.message;
          if (err.response && err.response.data && err.response.data.error) {
            message = err.response.data.error;
          }

          context.dispatch({
            type: NOTIFICATION_FAILURE,
            text: `Error saving custom tag colour: ${message}`,
          });
          throw err;
        }),
    [LOGIN]: (context, { username }) =>
      axios
        .post('http://127.0.0.1:9090/api/login', { username })
        .then(() => {
          context.dispatch({ type: LOAD_USER });
        })
        .catch(err => {
          let message = err.message;
          if (err.response && err.response.data && err.response.data.error) {
            message = err.response.data.error;
          }

          context.dispatch({
            type: NOTIFICATION_FAILURE,
            text: `Error logging in: ${message}`,
          });
          throw err;
        }),
    [LOGOUT]: context =>
      axios
        .post('http://127.0.0.1:9090/api/logout')
        .then(() => {
          context.commit({
            type: USER_LOADED,
            id: 0,
            username: '',
            tagColours: {},
          });
        })
        .catch(err => {
          let message = err.message;
          if (err.response && err.response.data && err.response.data.error) {
            message = err.response.data.error;
          }

          context.dispatch({
            type: NOTIFICATION_FAILURE,
            text: `Error logging out: ${message}`,
          });
          throw err;
        }),
  },
};
