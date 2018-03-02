<template>
  <transition-group name="notifications" tag="ul" class="notifications">
    <li
      v-for="notification in notifications"
      :class="{alert: true, [`alert-${notification.kind}`]: true}"
      :key="notification.id"
    >
      {{ notification.text }}
    </li>
  </transition-group>
</template>

<script>
export default {
  computed: {
    notifications() {
      return this.$store.state.notifications.notifications;
    },
  },
};
</script>

<style lang="scss">
@import 'style/_variables';

.notifications {
  position: fixed;
  top: 20px;
  right: 20px;
}

.alert {
  width: 300px;
}

.notifications-enter-active,
.notifications-leave-active {
  transition: all 0.3s;
}
.notifications-enter,
.notifications-leave-to {
  opacity: 0;
}
</style>
