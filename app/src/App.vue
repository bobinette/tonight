<template>
  <div id="app">
    <Navbar></Navbar>
    <router-view/>
    <Notifications></Notifications>
  </div>
</template>

<script>
import Navbar from '@/components/nav-bar/Navbar';
import Notifications from '@/modules/notifications/Notifications';

import { LOAD_USER } from '@/modules/user/state';

export default {
  name: 'app',
  components: {
    Navbar,
    Notifications,
  },
  mounted() {
    this.$store.dispatch({ type: LOAD_USER }).catch();
  },
};
</script>
