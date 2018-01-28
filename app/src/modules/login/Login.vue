<template>
  <div id="username-input" class="container">
    Login:
    <input
      type="text"
      class="form-control"
      placeholder="Username..."
      v-model="username"
      @keydown.enter="login"
    >
  </div>
</template>

<script>
import { LOGIN } from '@/modules/user/state';

export default {
  data() {
    return {
      username: '',
    };
  },
  methods: {
    login() {
      this.$store
        .dispatch({ type: LOGIN, username: this.username })
        .then(() => {
          this.$router.push({ path: '/' });
        })
        .catch();
    },
  },
};
</script>

<style lang="scss">
#username-input {
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}
</style>
