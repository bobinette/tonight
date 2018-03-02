// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import VueAutosize from 'vue-autosize';
import { sync } from 'vuex-router-sync';

// Start by loading my custom CSS
import '@/style/base.scss';

// Add font-awesome icons
import 'font-awesome/scss/font-awesome.scss';

import App from './App';

import router from './router';
import store from './store';

Vue.config.productionTip = false;

// Use additional vue plugins
Vue.use(VueAutosize);

// Define custom modifiers
Vue.config.keyCodes.esc = 27;

// Sync the router in the store to use the query string
// in vuex
sync(store, router);

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App },

  store,
});
