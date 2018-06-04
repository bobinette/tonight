<template>
  <nav class="navbar navbar-row navbar-top navbar-inverse bg-primary white">
    <div class="container">
      <div class="row w-100 flex-align-center flex-space-between">
        <h5>
          <img :src="logo" height="24" width="24">
          Tonight
        </h5>
        <span v-if="userid !== 0">
          <span>{{ username }}</span>
          <button class="btn btn-link white" @click="logout" >
            <i class="fa fa-sign-out"></i>
          </button>
        </span>
        <span  v-else>
          <button class="btn btn-link white" @click="login" >
            <i class="fa fa-sign-in"></i>
          </button>
        </span>
      </div>
    </div>
  </nav>
</template>

<script>
import logo from '@/assets/logo-i.png';
import { LOGIN, LOGOUT } from '@/modules/user/state';

export default {
  name: 'navbar',
  data() {
    return {
      logo,
    };
  },
  computed: {
    userid() {
      return this.$store.state.user.user.id;
    },
    username() {
      return this.$store.getters.username;
    },
  },
  methods: {
    login() {
      this.$store.dispatch({ type: LOGIN }).catch();
    },
    logout() {
      this.$store.dispatch({ type: LOGOUT }).catch();
    },
  },
};
</script>
