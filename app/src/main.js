// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import Vuex from 'vuex';
import VueAutosize from 'vue-autosize';

// Add font-awesome icons
import 'font-awesome/scss/font-awesome.scss';

import App from './App';

import store from './store';

Vue.config.productionTip = false;
Vue.use(Vuex);
Vue.use(VueAutosize);

/* eslint-disable no-new */
new Vue({
  el: '#app',
  template: '<App/>',
  components: { App },

  store: new Vuex.Store(store),
});
