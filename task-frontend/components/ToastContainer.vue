<template>
  <div class="tf-toast-container">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        class="tf-toast"
        :class="toast.type"
        @click="toastStore.remove(toast.id)"
      >
        <i :class="iconClass(toast.type)"></i>
        <span>{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { useToastStore } from "~/store/toast";

const toastStore = useToastStore();

function iconClass(type) {
  const map = {
    "is-success": "fas fa-check-circle text-accent",
    "is-error":   "fas fa-times-circle",
    "is-info":    "fas fa-info-circle",
  };
  return map[type] || "fas fa-info-circle";
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>