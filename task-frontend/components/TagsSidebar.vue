<template>
  <div v-if="tags.length > 0" class="tags-panel">
    <div class="tags-panel-header">
      <span class="tf-label" style="margin:0">Etiquetas</span>
      <button
        v-if="activeTag"
        class="btn-ghost button is-small"
        @click="$emit('clear')"
        style="font-size:0.65rem;padding:2px 8px"
      >
        limpiar
      </button>
    </div>
    <div class="tags-panel-list">
      <span
        v-for="tag in tags"
        :key="tag"
        class="tag-badge"
        :class="{ 'is-active': activeTag === tag }"
        @click="$emit('select', tag)"
      >
        #{{ tag }}
      </span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  tags:      { type: Array, default: () => [] },
  activeTag: { type: String, default: null },
});
defineEmits(["select", "clear"]);
</script>

<style scoped>
.tags-panel {
  background: #161616;
  border: 1px solid #2e2e2e;
  border-radius: 8px;
  padding: 1rem 1.25rem;
  margin-bottom: 1.5rem;
}
.tags-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}
.tags-panel-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}
</style>