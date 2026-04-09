<template>
  <div class="tag-input-wrapper" @click="focusInput">
    <span
      v-for="(tag, i) in modelValue"
      :key="i"
      class="tag-badge"
    >
      #{{ tag }}
      <span class="tag-remove" @click.stop="removeTag(i)">
        <i class="fas fa-times" style="font-size:0.55rem"></i>
      </span>
    </span>
    <input
      ref="inputRef"
      v-model="inputVal"
      :placeholder="modelValue.length === 0 ? placeholder : ''"
      @keydown.enter.prevent="addTag"
      @keydown.tab.prevent="addTag"
      @keydown.backspace="onBackspace"
      @keydown.188.prevent="addTag"
    />
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Añadir tag y presionar Enter…' },
});
const emit = defineEmits(["update:modelValue"]);

const inputRef = ref(null);
const inputVal = ref("");

function focusInput() {
  inputRef.value?.focus();
}

function addTag() {
  const val = inputVal.value.trim().toLowerCase().replace(/[^a-z0-9-_]/g, "");
  if (val && !props.modelValue.includes(val)) {
    emit("update:modelValue", [...props.modelValue, val]);
  }
  inputVal.value = "";
}

function removeTag(idx) {
  const updated = [...props.modelValue];
  updated.splice(idx, 1);
  emit("update:modelValue", updated);
}

function onBackspace() {
  if (!inputVal.value && props.modelValue.length > 0) {
    removeTag(props.modelValue.length - 1);
  }
}
</script>