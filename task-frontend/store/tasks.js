// store/tasks.js
import { defineStore } from "pinia";

// ─── Helpers para normalizar la respuesta del backend ────────────────────────
// El backend Go puede devolver: array directo, { data: [] }, { tasks: [] }, etc.
function normalizeList(raw) {
  if (!raw) return [];
  if (Array.isArray(raw)) return raw;
  for (const key of ["data", "tasks", "items", "results"]) {
    if (Array.isArray(raw[key])) return raw[key];
  }
  return [];
}

function normalizeItem(raw) {
  if (!raw) return raw;
  if (raw.data && typeof raw.data === "object" && !Array.isArray(raw.data)) {
    return raw.data;
  }
  return raw;
}

// ─── Store ───────────────────────────────────────────────────────────────────
export const useTasksStore = defineStore("tasks", {
  state: () => ({
    tasks: [],
    loading: false,
    error: null,
    filter: "all",
    search: "",
    activeTag: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
    },
  }),

  getters: {
    filteredTasks(state) {
      if (!Array.isArray(state.tasks)) return [];
      let result = [...state.tasks];

      if (state.filter === "pending") {
        result = result.filter((t) => !t.completed);
      } else if (state.filter === "completed") {
        result = result.filter((t) => t.completed);
      }

      if (state.activeTag) {
        result = result.filter(
          (t) => Array.isArray(t.tags) && t.tags.includes(state.activeTag)
        );
      }

      if (state.search.trim()) {
        const q = state.search.toLowerCase();
        result = result.filter(
          (t) =>
            (t.title || "").toLowerCase().includes(q) ||
            (t.description || "").toLowerCase().includes(q)
        );
      }

      return result;
    },

    paginatedTasks(state) {
      const filtered = this.filteredTasks;
      const start = (state.pagination.page - 1) * state.pagination.limit;
      return filtered.slice(start, start + state.pagination.limit);
    },

    totalPages(state) {
      return Math.max(1, Math.ceil(this.filteredTasks.length / state.pagination.limit));
    },

    allTags(state) {
      if (!Array.isArray(state.tasks)) return [];
      const tagSet = new Set();
      state.tasks.forEach((t) => {
        if (Array.isArray(t.tags)) t.tags.forEach((tag) => tagSet.add(tag));
      });
      return [...tagSet].sort();
    },

    stats(state) {
      const list = Array.isArray(state.tasks) ? state.tasks : [];
      const total = list.length;
      const completed = list.filter((t) => t.completed).length;
      const pending = total - completed;
      return { total, completed, pending };
    },
  },

  actions: {
    async fetchTasks() {
      this.loading = true;
      this.error = null;
      try {
        const { $api } = useNuxtApp();
        const { data } = await $api.get("/tasks");
        this.tasks = normalizeList(data);
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    async createTask(payload) {
      const { $api } = useNuxtApp();
      const { data } = await $api.post("/tasks", payload);
      const task = normalizeItem(data);
      if (!Array.isArray(this.tasks)) this.tasks = [];
      this.tasks.unshift(task);
      return task;
    },

    async updateTask(id, payload) {
      const { $api } = useNuxtApp();
      const { data } = await $api.put(`/tasks/${id}`, payload);
      const task = normalizeItem(data);
      const idx = this.tasks.findIndex((t) => t._id === id || t.id === id);
      if (idx !== -1) this.tasks[idx] = task;
      return task;
    },

    async markCompleted(id) {
      const { $api } = useNuxtApp();
      const { data } = await $api.put(`/tasks/${id}/complete`);
      const task = normalizeItem(data);
      const idx = this.tasks.findIndex((t) => t._id === id || t.id === id);
      if (idx !== -1) this.tasks[idx] = task;
      return task;
    },

    async deleteTask(id) {
      const { $api } = useNuxtApp();
      await $api.delete(`/tasks/${id}`);
      this.tasks = this.tasks.filter((t) => t._id !== id && t.id !== id);
    },

    async fetchByTag(tag) {
      this.loading = true;
      this.error = null;
      try {
        const { $api } = useNuxtApp();
        const { data } = await $api.get(`/tasks/tag/${encodeURIComponent(tag)}`);
        this.tasks = normalizeList(data);
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    setFilter(f) {
      this.filter = f;
      this.pagination.page = 1;
    },

    setSearch(q) {
      this.search = q;
      this.pagination.page = 1;
    },

    setActiveTag(tag) {
      this.activeTag = this.activeTag === tag ? null : tag;
      this.pagination.page = 1;
    },

    setPage(p) {
      this.pagination.page = p;
    },

    setLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page = 1;
    },
  },
});