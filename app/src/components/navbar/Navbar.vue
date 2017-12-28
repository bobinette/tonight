<template>
  <nav id="navbar" class="navbar navbar-row navbar-top navbar-inverse bg-primary">
    <div class="container">
      <div class="row flex-align-center">
        <h1 class="col-md-8 offset-md-2">
          <img :src="logo" height="36" width="36">
          Tonight
        </h1>
        <span class="col-md-2 Username">
          <span>{{ username }}</span>
        </span>
      </div>
    </div>
  </nav>
</template>

<script>
import logo from '@/assets/logo-i.png';

export default {
  name: 'navbar',
  data() {
    return {
      logo,
    };
  },
  computed: {
    username() {
      return this.$store.getters.username;
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
