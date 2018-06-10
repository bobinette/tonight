import Vue from 'vue';
import Vuex from 'vuex';

import NavBar from './NavBar';

Vue.use(Vuex);

test('render navbar', () => {
  const Constructor = Vue.extend(NavBar);
  const store = new Vuex.Store({
    getters: {
      userid: () => 1,
      username: () => 'test',
    },
  });
  const vm = new Constructor({ store }).$mount();

  expect(vm.$el).toMatchSnapshot();
});
