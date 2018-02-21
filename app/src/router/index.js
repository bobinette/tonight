import Vue from 'vue';
import Router from 'vue-router';
import Home from '@/modules/home/Home';
import Login from '@/modules/login/Login';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home,
    },
  ],
});
