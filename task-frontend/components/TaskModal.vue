<template>
  <div class="modal" :class="{ 'is-active': isOpen }">
    <div class="modal-background" @click="close"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">
          {{ isEditing ? "Editar tarea" : "Nueva tarea" }}
        </p>
        <button class="delete" aria-label="close" @click="close"></button>
      </header>

      <section class="modal-card-body">
        <div class="field">
          <label class="tf-label">Título *</label>
          <input
            v-model="form.title"
            class="input tf-input"
            type="text"
            placeholder="Nombre de la tarea…"
            @keydown.enter="submit"
          />
          <p v-if="errors.title" class="help is-danger" style="font-size:0.75rem;margin-top:4px">
            {{ errors.title }}
          </p>
        </div>

        <div class="field">
          <label class="tf-label">Descripción</label>
          <textarea
            v-model="form.description"
            class="textarea tf-input tf-textarea"
            placeholder="Detalles opcionales…"
            rows="3"
          ></textarea>
        </div>

        <div class="field">
          <label class="tf-label">Etiquetas</label>
          <TagInput v-model="form.tags" />
          <p class="help" style="color:#555;font-size:0.72rem;margin-top:6px">
            Presiona <kbd style="background:#2e2e2e;color:#888;padding:1px 5px;border-radius:3px;font-size:0.65rem">Enter</kbd>
            o <kbd style="background:#2e2e2e;color:#888;padding:1px 5px;border-radius:3px;font-size:0.65rem">,</kbd> para añadir
          </p>
        </div>
      </section>

      <footer class="modal-card-foot">
        <button class="button btn-ghost" @click="close">Cancelar</button>
        <button
          class="button btn-primary"
          :class="{ 'is-loading': loading }"
          @click="submit"
        >
          {{ isEditing ? "Guardar cambios" : "Crear tarea" }}
        </button>
      </footer>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  isOpen:   { type: Boolean, default: false },
  task:     { type: Object, default: null },
});
const emit = defineEmits(["close", "submit"]);

const loading = ref(false);
const errors  = ref({});
const isEditing = computed(() => !!props.task);

const form = reactive({
  title: "",
  description: "",
  tags: [],
});

watch(
  () => props.isOpen,
  (val) => {
    if (val) {
      if (props.task) {
        form.title       = props.task.title || "";
        form.description = props.task.description || "";
        form.tags        = [...(props.task.tags || [])];
      } else {
        form.title       = "";
        form.description = "";
        form.tags        = [];
      }
      errors.value = {};
    }
  }
);

function validate() {
  errors.value = {};
  if (!form.title.trim()) errors.value.title = "El título es obligatorio";
  return Object.keys(errors.value).length === 0;
}

async function submit() {
  if (!validate()) return;
  loading.value = true;
  try {
    await emit("submit", {
      title:       form.title.trim(),
      description: form.description.trim(),
      tags:        form.tags,
    });
  } finally {
    loading.value = false;
  }
}

function close() {
  emit("close");
}
</script>