<template>
  <div id="app">
    <Navbar></Navbar>
    <router-view/>
    <Notifications></Notifications>
  </div>
</template>

<script>
// Modules
import Notifications from '@/modules/notifications/Notifications';

import { LOAD_USER } from '@/modules/user/state';

// Components
import Navbar from '@/components/navbar/Navbar';

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
