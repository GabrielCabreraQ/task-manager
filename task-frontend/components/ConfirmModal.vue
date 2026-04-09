<template>
  <div class="modal" :class="{ 'is-active': isOpen }">
    <div class="modal-background" @click="$emit('cancel')"></div>
    <div class="modal-card" style="max-width:400px">
      <header class="modal-card-head">
        <p class="modal-card-title">{{ title }}</p>
        <button class="delete" @click="$emit('cancel')"></button>
      </header>
      <section class="modal-card-body">
        <p style="color:#888;font-size:0.9rem">{{ message }}</p>
      </section>
      <footer class="modal-card-foot">
        <button class="button btn-ghost" @click="$emit('cancel')">Cancelar</button>
        <button
          class="button"
          style="background:#ff4757;color:#fff;border:none;font-family:'Space Mono',monospace;font-size:0.78rem"
          :class="{ 'is-loading': loading }"
          @click="$emit('confirm')"
        >
          {{ confirmLabel }}
        </button>
      </footer>
    </div>
  </div>
</template>

<script setup>
defineProps({
  isOpen:       { type: Boolean, default: false },
  title:        { type: String, default: "Confirmar acción" },
  message:      { type: String, default: "¿Estás seguro?" },
  confirmLabel: { type: String, default: "Confirmar" },
  loading:      { type: Boolean, default: false },
});
defineEmits(["confirm", "cancel"]);
</script>