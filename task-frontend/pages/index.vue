<template>
  <div>
    <!-- Page Header -->
    <div class="page-header" style="margin-top:2rem">
      <div class="page-eyebrow"> GESTIÓN DE TAREAS</div>
      <h1 class="page-title">
        Mis tareas
      </h1>
      <p class="page-subtitle">
        Organiza tu trabajo de forma eficiente y sin complicaciones. Crea, edita y gestiona tus tareas fácilmente.
      </p>
    </div>

    <!-- Action Row -->
    <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:1.5rem;gap:1rem;flex-wrap:wrap">
      <div style="display:flex;gap:0.5rem;align-items:center">
        <button class="button btn-primary" @click="openCreate">
          <i class="fas fa-plus" style="margin-right:6px"></i>
          Nueva tarea
        </button>
        <button
          class="button btn-ghost"
          :class="{ 'is-loading': tasksStore.loading }"
          @click="tasksStore.fetchTasks()"
          title="Recargar"
        >
          <i class="fas fa-sync-alt"></i>
        </button>
      </div>

      <!-- Error banner -->
      <div v-if="tasksStore.error" style="color:#ff4757;font-size:0.8rem;font-family:'Space Mono',monospace">
        <i class="fas fa-exclamation-circle"></i>
        {{ tasksStore.error }}
      </div>
    </div>

    <!-- Tags Panel -->
    <TagsSidebar
      :tags="tasksStore.allTags"
      :active-tag="tasksStore.activeTag"
      @select="tasksStore.setActiveTag($event)"
      @clear="tasksStore.setActiveTag(null)"
    />

    <!-- Filter Bar -->
    <div class="filter-bar">
      <span class="filter-label">Vista:</span>
      <button
        v-for="f in filters"
        :key="f.value"
        class="filter-btn"
        :class="{ 'is-active': tasksStore.filter === f.value }"
        @click="tasksStore.setFilter(f.value)"
      >
        {{ f.label }}
      </button>
      <div class="filter-divider"></div>
      <div class="search-wrapper">
        <i class="fas fa-search search-icon"></i>
        <input
          :value="tasksStore.search"
          type="text"
          placeholder="Buscar tareas…"
          @input="tasksStore.setSearch($event.target.value)"
        />
      </div>
    </div>

    <!-- Loader -->
    <div v-if="tasksStore.loading && !tasksStore.tasks.length" class="tf-loader">
      <div class="spinner"></div>
      <span class="loader-text">CARGANDO TAREAS…</span>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="!tasksStore.loading && tasksStore.filteredTasks.length === 0"
      class="empty-state"
    >
      <div class="empty-icon"><i class="fas fa-inbox"></i></div>
      <p class="empty-title">
        {{ tasksStore.tasks.length === 0 ? "No hay tareas aún" : "Sin resultados" }}
      </p>
      <p class="empty-sub">
        {{ tasksStore.tasks.length === 0 ? "Crea tu primera tarea usando el botón de arriba." : "Intenta ajustar los filtros." }}
      </p>
    </div>

    <!-- Task List -->
    <TransitionGroup v-else name="fade" tag="div" style="display:flex;flex-direction:column;gap:0.75rem">
      <TaskCard
        v-for="task in tasksStore.paginatedTasks"
        :key="task._id || task.id"
        :task="task"
        @complete="handleComplete"
        @edit="openEdit"
        @delete="openDelete"
        @filter-tag="tasksStore.setActiveTag($event)"
      />
    </TransitionGroup>

    <!-- ── Pagination + Page Size ───────────────────────────────────────── -->
    <div v-if="tasksStore.filteredTasks.length > 0" class="pagination-bar">

      <!-- Info -->
      <span class="pagination-info">
        {{ rangeStart }}–{{ rangeEnd }} de {{ tasksStore.filteredTasks.length }}
      </span>

      <!-- Page size selector -->
      <div class="page-size-selector">
        <span class="page-size-label">Por página:</span>
        <div class="page-size-options">
          <button
            v-for="size in pageSizes"
            :key="size.value"
            class="size-btn"
            :class="{ 'is-active': tasksStore.pagination.limit === size.value }"
            @click="tasksStore.setLimit(size.value)"
          >
            {{ size.label }}
          </button>
        </div>
      </div>

      <!-- Page navigation -->
      <div v-if="tasksStore.totalPages > 1" class="tf-pagination" style="margin-top:0">
        <!-- Prev -->
        <button
          class="page-btn"
          :disabled="tasksStore.pagination.page <= 1"
          @click="tasksStore.setPage(tasksStore.pagination.page - 1)"
        >
          <i class="fas fa-chevron-left" style="font-size:0.7rem"></i>
        </button>

        <!-- Pages con ellipsis -->
        <template v-for="item in pageItems" :key="item.key">
          <span v-if="item.ellipsis" class="page-ellipsis">…</span>
          <button
            v-else
            class="page-btn"
            :class="{ 'is-active': tasksStore.pagination.page === item.page }"
            @click="tasksStore.setPage(item.page)"
          >
            {{ item.page }}
          </button>
        </template>

        <!-- Next -->
        <button
          class="page-btn"
          :disabled="tasksStore.pagination.page >= tasksStore.totalPages"
          @click="tasksStore.setPage(tasksStore.pagination.page + 1)"
        >
          <i class="fas fa-chevron-right" style="font-size:0.7rem"></i>
        </button>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <TaskModal
      :is-open="showModal"
      :task="editingTask"
      @close="closeModal"
      @submit="handleSubmit"
    />

    <!-- Confirm Delete Modal -->
    <ConfirmModal
      :is-open="showConfirm"
      title="Eliminar tarea"
      :message="`¿Eliminar &quot;${deletingTask?.title}&quot;? Esta acción no se puede deshacer.`"
      confirm-label="Eliminar"
      :loading="deleting"
      @confirm="handleDelete"
      @cancel="showConfirm = false"
    />
  </div>
</template>

<script setup>
import { useTasksStore } from "~/store/tasks";
import { useToastStore  } from "~/store/toast";

const tasksStore = useTasksStore();
const toastStore = useToastStore();

// ─── Fetch on mount ───────────────────────────────────────────────────────
onMounted(() => tasksStore.fetchTasks());

// ─── Filters ─────────────────────────────────────────────────────────────
const filters = [
  { label: "Todas",       value: "all"       },
  { label: "Pendientes",  value: "pending"   },
  { label: "Completadas", value: "completed" },
];

// ─── Page size options ────────────────────────────────────────────────────
const pageSizes = [
  { label: "1",   value: 1   },
  { label: "5",   value: 5   },
  { label: "10",  value: 10  },
  { label: "20",  value: 20  },
  { label: "50",  value: 50  },
  { label: "100", value: 100 },
  { label: "200", value: 200 },
  { label: "500", value: 500 },
];

// ─── Pagination range info ────────────────────────────────────────────────
const rangeStart = computed(() => {
  const { page, limit } = tasksStore.pagination;
  return Math.min((page - 1) * limit + 1, tasksStore.filteredTasks.length);
});

const rangeEnd = computed(() => {
  const { page, limit } = tasksStore.pagination;
  return Math.min(page * limit, tasksStore.filteredTasks.length);
});

// ─── Page items con ellipsis ──────────────────────────────────────────────
const pageItems = computed(() => {
  const total   = tasksStore.totalPages;
  const current = tasksStore.pagination.page;
  const items   = [];

  if (total <= 7) {
    for (let i = 1; i <= total; i++) items.push({ key: i, page: i });
    return items;
  }

  // Siempre mostrar primera página
  items.push({ key: 1, page: 1 });

  if (current > 3) items.push({ key: "e1", ellipsis: true });

  // Páginas alrededor de la actual
  const start = Math.max(2, current - 1);
  const end   = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) items.push({ key: i, page: i });

  if (current < total - 2) items.push({ key: "e2", ellipsis: true });

  // Siempre mostrar última página
  items.push({ key: total, page: total });

  return items;
});

// ─── Modal state ─────────────────────────────────────────────────────────
const showModal   = ref(false);
const editingTask = ref(null);

function openCreate() {
  editingTask.value = null;
  showModal.value   = true;
}

function openEdit(task) {
  editingTask.value = task;
  showModal.value   = true;
}

function closeModal() {
  showModal.value   = false;
  editingTask.value = null;
}

// ─── Create / Update ─────────────────────────────────────────────────────
async function handleSubmit(payload) {
  try {
    if (editingTask.value) {
      const id = editingTask.value._id || editingTask.value.id;
      await tasksStore.updateTask(id, payload);
      toastStore.success("Tarea actualizada correctamente");
    } else {
      await tasksStore.createTask(payload);
      toastStore.success("Tarea creada correctamente");
    }
    closeModal();
  } catch (e) {
    toastStore.error(e.message);
  }
}

// ─── Complete ─────────────────────────────────────────────────────────────
async function handleComplete(task) {
  try {
    const id = task._id || task.id;
    await tasksStore.markCompleted(id);
    toastStore.success("¡Tarea completada!");
  } catch (e) {
    toastStore.error(e.message);
  }
}

// ─── Delete ──────────────────────────────────────────────────────────────
const showConfirm  = ref(false);
const deletingTask = ref(null);
const deleting     = ref(false);

function openDelete(task) {
  deletingTask.value = task;
  showConfirm.value  = true;
}

async function handleDelete() {
  deleting.value = true;
  try {
    const id = deletingTask.value._id || deletingTask.value.id;
    await tasksStore.deleteTask(id);
    toastStore.success("Tarea eliminada");
    showConfirm.value  = false;
    deletingTask.value = null;
  } catch (e) {
    toastStore.error(e.message);
  } finally {
    deleting.value = false;
  }
}

// ─── SEO ──────────────────────────────────────────────────────────────────
useHead({ title: "TaskManagment — Mis Tareas" });
</script>

<style scoped>
/* ── Pagination bar ───────────────────────────────────────────────────── */
.pagination-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #1e1e1e;
}

.pagination-info {
  font-family: "Space Mono", monospace;
  font-size: 0.68rem;
  color: #555;
  letter-spacing: 0.06em;
}

/* ── Page size selector ───────────────────────────────────────────────── */
.page-size-selector {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  flex-wrap: wrap;
  justify-content: center;
}

.page-size-label {
  font-family: "Space Mono", monospace;
  font-size: 0.65rem;
  color: #555;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  white-space: nowrap;
}

.page-size-options {
  display: flex;
  gap: 0.3rem;
  flex-wrap: wrap;
  justify-content: center;
}

.size-btn {
  background: transparent;
  border: 1px solid #2e2e2e;
  color: #666;
  font-family: "Space Mono", monospace;
  font-size: 0.68rem;
  padding: 4px 9px;
  border-radius: 3px;
  cursor: pointer;
  transition: all 0.15s ease;
  letter-spacing: 0.03em;
  line-height: 1;

  &:hover {
    border-color: #b5f542;
    color: #b5f542;
  }

  &.is-active {
    background: #b5f542;
    color: #0d0d0d;
    border-color: #b5f542;
    font-weight: 700;
  }
}

/* ── Ellipsis ─────────────────────────────────────────────────────────── */
.page-ellipsis {
  font-family: "Space Mono", monospace;
  font-size: 0.75rem;
  color: #444;
  padding: 0 2px;
  line-height: 34px;
}
</style>