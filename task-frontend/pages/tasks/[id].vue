<template>
  <div>
    <div v-if="loading" class="tf-loader" style="margin-top:4rem">
      <div class="spinner"></div>
      <span class="loader-text">CARGANDO TAREA…</span>
    </div>

    <div v-else-if="error" style="text-align:center;padding:4rem 2rem">
      <p style="color:#ff4757;font-family:'Space Mono',monospace;font-size:0.85rem">
        <i class="fas fa-exclamation-circle"></i> {{ error }}
      </p>
      <NuxtLink to="/" class="button btn-ghost" style="margin-top:1rem">
        ← Volver
      </NuxtLink>
    </div>

    <div v-else-if="task">
      <!-- Back -->
      <div style="margin-top:2rem;margin-bottom:1.5rem">
        <NuxtLink to="/" class="button btn-ghost">
          <i class="fas fa-arrow-left" style="margin-right:6px;font-size:0.75rem"></i>
          Volver a tareas
        </NuxtLink>
      </div>

      <!-- Header -->
      <div class="page-header">
        <div class="page-eyebrow">// DETALLE DE TAREA</div>
        <h1 class="page-title" style="font-size:clamp(1.5rem,3vw,2.2rem)">
          {{ task.title }}
        </h1>
      </div>

      <!-- Card -->
      <div style="background:#161616;border:1px solid #2e2e2e;border-radius:10px;padding:2rem;max-width:700px">

        <!-- Status + Date row -->
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:1.5rem;flex-wrap:wrap;gap:0.75rem">
          <span :class="['status-badge', task.completed ? 'completed' : 'pending']">
            <span class="dot"></span>
            {{ task.completed ? "Completada" : "Pendiente" }}
          </span>
          <span class="text-dim text-mono" style="font-size:0.7rem">
            Creada el {{ formatDate(task.createdAt) }}
          </span>
        </div>

        <!-- Description -->
        <div style="margin-bottom:1.5rem">
          <label class="tf-label">Descripción</label>
          <p style="color:#aaa;font-size:0.9rem;line-height:1.7">
            {{ task.description || "Sin descripción." }}
          </p>
        </div>

        <div class="divider"></div>

        <!-- Tags -->
        <div style="margin-bottom:1.5rem">
          <label class="tf-label">Etiquetas</label>
          <div style="display:flex;flex-wrap:wrap;gap:0.4rem;margin-top:0.5rem">
            <span v-for="tag in task.tags" :key="tag" class="tag-badge">#{{ tag }}</span>
            <span v-if="!task.tags || task.tags.length === 0" class="text-dim" style="font-size:0.8rem">
              Sin etiquetas
            </span>
          </div>
        </div>

        <div class="divider"></div>

        <!-- Actions -->
        <div style="display:flex;gap:0.75rem;flex-wrap:wrap">
          <button
            v-if="!task.completed"
            class="button btn-primary"
            :class="{ 'is-loading': completing }"
            @click="handleComplete"
          >
            <i class="fas fa-check" style="margin-right:6px"></i>
            Marcar completada
          </button>
          <button class="button btn-ghost" @click="openEdit">
            <i class="fas fa-pen" style="margin-right:6px;font-size:0.75rem"></i>
            Editar
          </button>
          <button class="button btn-danger" @click="showConfirm = true">
            <i class="fas fa-trash" style="margin-right:6px;font-size:0.75rem"></i>
            Eliminar
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <TaskModal
      :is-open="showModal"
      :task="task"
      @close="showModal = false"
      @submit="handleUpdate"
    />

    <!-- Confirm Delete -->
    <ConfirmModal
      :is-open="showConfirm"
      title="Eliminar tarea"
      :message="`¿Eliminar &quot;${task?.title}&quot;? Esta acción no se puede deshacer.`"
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

const route      = useRoute();
const router     = useRouter();
const tasksStore = useTasksStore();
const toastStore = useToastStore();
const { $api }   = useNuxtApp();

const task       = ref(null);
const loading    = ref(true);
const error      = ref(null);
const completing = ref(false);
const deleting   = ref(false);
const showModal  = ref(false);
const showConfirm = ref(false);

const id = computed(() => route.params.id);

onMounted(async () => {
  try {
    const { data } = await $api.get(`/tasks/${id.value}`);
    task.value = data;
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
});

function formatDate(dateStr) {
  if (!dateStr) return "";
  return new Date(dateStr).toLocaleDateString("es-CL", {
    day: "2-digit", month: "long", year: "numeric",
  });
}

async function handleComplete() {
  completing.value = true;
  try {
    const updated = await tasksStore.markCompleted(id.value);
    task.value = updated;
    toastStore.success("¡Tarea completada!");
  } catch (e) {
    toastStore.error(e.message);
  } finally {
    completing.value = false;
  }
}

function openEdit() {
  showModal.value = true;
}

async function handleUpdate(payload) {
  try {
    const updated = await tasksStore.updateTask(id.value, payload);
    task.value = updated;
    toastStore.success("Tarea actualizada");
    showModal.value = false;
  } catch (e) {
    toastStore.error(e.message);
  }
}

async function handleDelete() {
  deleting.value = true;
  try {
    await tasksStore.deleteTask(id.value);
    toastStore.success("Tarea eliminada");
    router.push("/");
  } catch (e) {
    toastStore.error(e.message);
    deleting.value = false;
  }
}

useHead(() => ({
  title: task.value ? `${task.value.title} — TaskFlow` : "TaskFlow",
}));
</script>