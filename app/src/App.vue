<template>
  <div id="app">
    <Navbar></Navbar>
    <div class="container">
      <TaskList class="col-md-12"></TaskList>
    </div>
  </div>
</template>

<script>
// Modules
import TaskList from '@/modules/task-list/TaskList';
import { LOAD_COOKIE } from '@/modules/user/state';

// Components
import Navbar from '@/components/navbar/Navbar';

export default {
  name: 'app',
  components: {
    Navbar,
    TaskList,
  },
  mounted() {
    this.$store.dispatch({ type: LOAD_COOKIE });
  },
};
</script>

<style lang="scss">
@import 'style/base.scss';
</style>
