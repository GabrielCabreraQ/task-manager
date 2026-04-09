// store/tasks.js
import { defineStore } from "pinia";

// ─── Normaliza un item individual de respuesta ────────────────────────────────
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
    // Tareas de la página actual (devueltas por el backend)
    tasks: [],

    // Metadatos de paginación — vienen del backend (TaskList)
    pagination: {
      page:  1,
      limit: 10,
      total: 0,   // total real de documentos en MongoDB
    },

    // Filtros de UI (el backend no los soporta aún, se aplican localmente
    // sobre las tareas de la página actual)
    filter:    "all",   // 'all' | 'pending' | 'completed'
    search:    "",
    activeTag: null,

    loading: false,
    error:   null,
  }),

  getters: {
    // Tareas de la página actual filtradas en cliente
    // (solo afectan la visualización, no la paginación del servidor)
    visibleTasks(state) {
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

    // Total de páginas basado en el total real del backend
    totalPages(state) {
      if (!state.pagination.total || !state.pagination.limit) return 1;
      return Math.max(1, Math.ceil(state.pagination.total / state.pagination.limit));
    },

    // Rango visible: "11 – 20 de 47"
    rangeStart(state) {
      return Math.min(
        (state.pagination.page - 1) * state.pagination.limit + 1,
        state.pagination.total || 1
      );
    },

    rangeEnd(state) {
      return Math.min(
        state.pagination.page * state.pagination.limit,
        state.pagination.total || 0
      );
    },

    // Tags únicos de las tareas cargadas en esta página
    allTags(state) {
      if (!Array.isArray(state.tasks)) return [];
      const tagSet = new Set();
      state.tasks.forEach((t) => {
        if (Array.isArray(t.tags)) t.tags.forEach((tag) => tagSet.add(tag));
      });
      return [...tagSet].sort();
    },

    // Stats basados en el total real del backend
    stats(state) {
      const list = Array.isArray(state.tasks) ? state.tasks : [];
      const completedInPage = list.filter((t) => t.completed).length;
      return {
        total:     state.pagination.total || 0,
        completed: completedInPage,
        pending:   list.length - completedInPage,
      };
    },
  },

  actions: {
    // ── Fetch principal — llama al backend con page y limit reales ────────
    async fetchTasks() {
      this.loading = true;
      this.error   = null;
      try {
        const { $api } = useNuxtApp();
        const { data } = await $api.get("/tasks", {
          params: {
            page:  this.pagination.page,
            limit: this.pagination.limit,
          },
        });

        // El backend devuelve: { tasks: [...], total: N, page: N, limit: N }
        this.tasks              = Array.isArray(data.tasks) ? data.tasks : [];
        this.pagination.total   = data.total  ?? this.tasks.length;
        this.pagination.page    = data.page   ?? this.pagination.page;
        this.pagination.limit   = data.limit  ?? this.pagination.limit;
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    // ── CRUD ──────────────────────────────────────────────────────────────
    async createTask(payload) {
      const { $api } = useNuxtApp();
      const { data } = await $api.post("/tasks", payload);
      // Refrescar la página actual para reflejar el nuevo total
      await this.fetchTasks();
      return normalizeItem(data);
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
      // Si era la última tarea de la página, retroceder una página
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
          params: {
            page:  this.pagination.page,
            limit: this.pagination.limit,
          },
        });
        this.tasks            = Array.isArray(data.tasks) ? data.tasks : (Array.isArray(data) ? data : []);
        this.pagination.total = data.total ?? this.tasks.length;
      } catch (e) {
        this.error = e.message;
        this.tasks = [];
      } finally {
        this.loading = false;
      }
    },

    // ── Controles de paginación — cada cambio hace fetch al backend ───────
    setFilter(f) {
      this.filter = f;
      // Filtros de estado se aplican en cliente sobre la página actual
    },

    setSearch(q) {
      this.search = q;
      // Búsqueda se aplica en cliente sobre la página actual
    },

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