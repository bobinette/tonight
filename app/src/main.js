// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import Vuex from 'vuex';
import VueAutosize from 'vue-autosize';

import router from './router';

// Start by loading my custom CSS
import '@/style/base.scss';

// Add font-awesome icons
import 'font-awesome/scss/font-awesome.scss';

import App from './App';

import store from './store';

Vue.config.productionTip = false;

Vue.use(Vuex);
Vue.use(VueAutosize);

// Define custom modifiers
Vue.config.keyCodes.esc = 27;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App },

  store: new Vuex.Store(store),
});
