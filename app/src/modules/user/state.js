import axios from 'axios';

import apiUrl from '@/utils/apiUrl';

import { NOTIFICATION_FAILURE } from '@/modules/notifications/events';

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
    userid({ user }) {
      return user.id;
    },
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
        .get(`${apiUrl}/api/me`)
        .then(response => {
          context.commit(USER_LOADED, response.data);
          return response.data;
        })
        .catch(err => {
          if (err.response && err.response.status === 401) {
            context.commit(USER_LOADED, { id: 0, name: '', tagColours: {} });
            return;
          }

          throw err;
        }),
    [CUSTOMIZE_COLOUR]: (context, { tag, colour }) =>
      axios
        .post(`${apiUrl}/api/tags/${tag}`, { colour })
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
        .get(`${apiUrl}/api/oauth2/login?from=${encodeURI(window.location)}`, { username })
        .then(response => {
          window.location = response.data.url;
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
        .post(`${apiUrl}/api/logout`)
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
