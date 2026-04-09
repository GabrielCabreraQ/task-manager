// store/toast.js
import { defineStore } from "pinia";

export const useToastStore = defineStore("toast", {
  state: () => ({
    toasts: [],
  }),

  actions: {
    show(message, type = "info", duration = 3500) {
      const id = Date.now() + Math.random();
      this.toasts.push({ id, message, type });
      setTimeout(() => this.remove(id), duration);
    },

    success(message) {
      this.show(message, "is-success");
    },

    error(message) {
      this.show(message, "is-error", 5000);
    },

    info(message) {
      this.show(message, "is-info");
    },

    remove(id) {
      this.toasts = this.toasts.filter((t) => t.id !== id);
    },
  },
});