// store/tasks.js
import { defineStore } from "pinia";

function taskId(task) {
  return task?._id || task?.id || null;
}

function normalizeItem(raw) {
  if (!raw) return raw;
  if (raw.data && typeof raw.data === "object" && !Array.isArray(raw.data)) {
    return raw.data;
  }
  return raw;
}

export const useTasksStore = defineStore("tasks", {
  state: () => ({
    tasks: [],
    pagination: {
      page:  1,
      limit: 10,
      total: 0,
    },
    // Contadores globales — se obtienen con un fetch sin filtro
    globalTotal:     0,
    globalCompleted: 0,
    globalPending:   0,

    filter:    "all",   // 'all' | 'pending' | 'completed' — se manda al backend
    search:    "",      // búsqueda local sobre la página actual
    activeTag: null,
    loading:   false,
    error:     null,
  }),

  getters: {
    // Filtro de búsqueda local (sobre la página actual)
    visibleTasks(state) {
      if (!Array.isArray(state.tasks)) return [];
      if (!state.search.trim() && !state.activeTag) return state.tasks;

      let result = [...state.tasks];

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

    totalPages(state) {
      if (!state.pagination.total || !state.pagination.limit) return 1;
      return Math.max(1, Math.ceil(state.pagination.total / state.pagination.limit));
    },

    rangeStart(state) {
      return state.pagination.total === 0
        ? 0
        : (state.pagination.page - 1) * state.pagination.limit + 1;
    },

    rangeEnd(state) {
      return Math.min(
        state.pagination.page * state.pagination.limit,
        state.pagination.total
      );
    },

    allTags(state) {
      if (!Array.isArray(state.tasks)) return [];
      const tagSet = new Set();
      state.tasks.forEach((t) => {
        if (Array.isArray(t.tags)) t.tags.forEach((tag) => tagSet.add(tag));
      });
      return [...tagSet].sort();
    },

    // Stats siempre globales — nunca cambian por página o filtro activo
    stats(state) {
      return {
        total:     state.globalTotal,
        completed: state.globalCompleted,
        pending:   state.globalPending,
      };
    },
  },

  actions: {
    // ── Fetch con filtro de estado enviado al backend ──────────────────────
    async fetchTasks() {
      this.loading = true;
      this.error   = null;
      try {
        const { $api } = useNuxtApp();

        const params = {
          page:  this.pagination.page,
          limit: this.pagination.limit,
        };

        // Mandar el filtro de estado al backend solo si no es "all"
        if (this.filter === "completed") params.completed = true;
        if (this.filter === "pending")   params.completed = false;

        const { data } = await $api.get("/tasks", { params });

        this.tasks            = Array.isArray(data.tasks) ? data.tasks : [];
        this.pagination.total = data.total ?? this.tasks.length;
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    // Fetch global sin filtros para los contadores del navbar
    async fetchGlobalStats() {
      try {
        const { $api } = useNuxtApp();
        const { data } = await $api.get("/tasks", {
          params: { page: 1, limit: 1 }, // solo necesitamos el total
        });
        this.globalTotal = data.total ?? 0;

        // Contar completadas
        const comp = await $api.get("/tasks", {
          params: { page: 1, limit: 1, completed: true },
        });
        this.globalCompleted = comp.data.total ?? 0;
        this.globalPending   = this.globalTotal - this.globalCompleted;
      } catch {
        // fallback silencioso
      }
    },

    // ── CRUD ──────────────────────────────────────────────────────────────
    async createTask(payload) {
      const { $api } = useNuxtApp();
      const { data } = await $api.post("/tasks", payload);
      await this.fetchTasks();
      await this.fetchGlobalStats();
      return normalizeItem(data);
    },

    async updateTask(id, payload) {
      const { $api } = useNuxtApp();
      const { data } = await $api.put(`/tasks/${id}`, payload);
      const task = normalizeItem(data);
      const idx = this.tasks.findIndex((t) => taskId(t) === id);
      if (idx !== -1) this.tasks[idx] = task;
      return task;
    },

    async markCompleted(id) {
      const { $api } = useNuxtApp();
      const { data } = await $api.put(`/tasks/${id}/complete`);
      const task = normalizeItem(data);
      const idx = this.tasks.findIndex((t) => taskId(t) === id);
      if (idx !== -1) this.tasks[idx] = task;
      // Actualizar stats globales localmente (evita un request extra)
      this.globalCompleted = Math.min(this.globalCompleted + 1, this.globalTotal);
      this.globalPending   = Math.max(this.globalPending   - 1, 0);
      // Si el filtro activo es "pending", refrescar lista
      if (this.filter === "pending") await this.fetchTasks();
      return task;
    },

    async deleteTask(id) {
      const { $api } = useNuxtApp();
      const task = this.tasks.find((t) => taskId(t) === id);
      await $api.delete(`/tasks/${id}`);
      if (task?.completed) {
        this.globalCompleted = Math.max(this.globalCompleted - 1, 0);
      } else {
        this.globalPending = Math.max(this.globalPending - 1, 0);
      }
      this.globalTotal = Math.max(this.globalTotal - 1, 0);
      if (this.tasks.length === 1 && this.pagination.page > 1) {
        this.pagination.page -= 1;
      }
      await this.fetchTasks();
    },

    async fetchByTag(tag) {
      this.loading = true;
      this.error   = null;
      try {
        const { $api } = useNuxtApp();
        const { data } = await $api.get(`/tasks/tag/${encodeURIComponent(tag)}`, {
          params: { page: this.pagination.page, limit: this.pagination.limit },
        });
        this.tasks            = Array.isArray(data.tasks) ? data.tasks
                              : Array.isArray(data)       ? data : [];
        this.pagination.total = data.total ?? this.tasks.length;
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    // ── Controles ─────────────────────────────────────────────────────────
    async setFilter(f) {
      this.filter           = f;
      this.pagination.page  = 1;
      await this.fetchTasks(); // el filtro va al backend
    },

    setSearch(q)  { this.search    = q; },

    setActiveTag(tag) {
      this.activeTag = this.activeTag === tag ? null : tag;
    },

    async setPage(p) {
      this.pagination.page = p;
      await this.fetchTasks();
    },

    async setLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page  = 1;
      await this.fetchTasks();
    },
  },
});