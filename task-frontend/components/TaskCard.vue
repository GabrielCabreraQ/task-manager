<template>
  <div class="task-card" :class="{ 'is-completed': task.completed }">
    <div class="task-header">
      <h3 class="task-title">{{ task.title }}</h3>
      <div class="task-actions">
        <button
          v-if="!task.completed"
          class="btn-success-outline button is-small"
          title="Marcar como completada"
          @click="$emit('complete', task)"
        >
          <i class="fas fa-check"></i>
        </button>
        <button
          class="btn-ghost button is-small"
          title="Editar"
          @click="$emit('edit', task)"
        >
          <i class="fas fa-pen"></i>
        </button>
        <button
          class="btn-danger button is-small"
          title="Eliminar"
          @click="$emit('delete', task)"
        >
          <i class="fas fa-trash"></i>
        </button>
      </div>
    </div>

    <p v-if="task.description" class="task-description">
      {{ task.description }}
    </p>

    <div class="task-footer">
      <div class="task-tags">
        <span
          v-for="tag in task.tags"
          :key="tag"
          class="tag-badge"
          @click="$emit('filter-tag', tag)"
        >
          #{{ tag }}
        </span>
        <span v-if="!task.tags || task.tags.length === 0" class="text-dim" style="font-size:0.7rem">
          sin etiquetas
        </span>
      </div>

      <div style="display:flex;align-items:center;gap:0.75rem">
        <span :class="['status-badge', task.completed ? 'completed' : 'pending']">
          <span class="dot"></span>
          {{ task.completed ? "completada" : "pendiente" }}
        </span>
        <span class="task-meta">
          {{ formatDate(task.createdAt) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  task: { type: Object, required: true },
});
defineEmits(["complete", "edit", "delete", "filter-tag"]);

function formatDate(dateStr) {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  return d.toLocaleDateString("es-CL", { day: "2-digit", month: "short", year: "numeric" });
}
</script>