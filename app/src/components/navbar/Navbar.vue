<template>
  <nav id="navbar" class="navbar navbar-row navbar-top navbar-inverse bg-primary">
    <div class="container">
      <div class="row flex-align-center">
        <h1 class="col-md-8 offset-md-2">
          <img :src="logo" height="36" width="36">
          Tonight
        </h1>
        <span class="col-md-2 Username">
          <button class="btn btn-link white" @click="logout" v-if="userid !== 0">
            <i class="fa fa-sign-out"></i>
          </button>
          <span>{{ username }}</span>
        </span>
      </div>
    </div>
  </nav>
</template>

<script>
import { LOGOUT } from '@/modules/user/state';
import logo from '@/assets/logo-i.png';

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
    logout() {
      this.$store
        .dispatch({ type: LOGOUT })
        .then(() => {
          this.$router.push({ path: '/login' });
        })
        .catch();
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

  h1 {
    // margin: 0;
    text-align: center;
  }

  .Username {
    text-align: right;
  }
}
</style>
