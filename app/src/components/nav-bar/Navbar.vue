<template>
  <nav id="navbar" class="navbar navbar-row navbar-top navbar-inverse bg-primary">
    <div class="container">
      <div class="row flex-align-center">
        <h1 class="col-md-6 offset-md-3 NavBar__Title">
          <img :src="logo" height="36" width="36">
          Tonight
        </h1>
        <span class="col-md-3 NavBar__Username" v-if="userid !== 0">
          <span>{{ username }}</span>
          <button class="btn btn-link white" @click="logout" >
            <i class="fa fa-sign-out"></i>
          </button>
        </span>
        <span  class="col-md-3 Username" v-else>
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

<style lang="scss" scoped>
#navbar {
  color: white;
  margin-bottom: 1rem;

  flex-direction: row;
  align-items: center;
  justify-content: center;

  .NavBar__Title {
    text-align: center;
  }

  .NavBar__Username {
    text-align: right;
  }
}
</style>
