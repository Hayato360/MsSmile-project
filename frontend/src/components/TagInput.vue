<script setup>
import { ref, watch, onMounted } from 'vue'
import { X } from 'lucide-vue-next'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: 'Type and press Enter...',
  },
})

const emit = defineEmits(['update:modelValue'])

const tags = ref([])
const inputValue = ref('')

// Initialize tags from modelValue string
const initTags = () => {
  if (props.modelValue) {
    // Split by comma or newline
    tags.value = props.modelValue.split(/,|\n/).map(t => t.trim()).filter(t => t)
  } else {
    tags.value = []
  }
}

watch(() => props.modelValue, (newValue) => {
  // Only update if array serialization doesn't match new value (avoid cursor jump loops if we were syncing differently)
  // But here we sync one way mostly.
  const currentString = tags.value.join(',')
  if (newValue !== currentString) {
     initTags()
  }
})

onMounted(() => {
  initTags()
})

const addTag = () => {
  const val = inputValue.value.trim()
  if (val) {
    tags.value.push(val)
    inputValue.value = ''
    emitUpdate()
  }
}

const removeTag = (index) => {
  tags.value.splice(index, 1)
  emitUpdate()
}

const onBackspace = () => {
  if (inputValue.value === '' && tags.value.length > 0) {
    tags.value.pop()
    emitUpdate()
  }
}

const emitUpdate = () => {
  emit('update:modelValue', tags.value.join(','))
}
</script>

<template>
  <div class="tag-input-container">
    <div v-for="(tag, index) in tags" :key="index" class="tag-pill">
      {{ tag }}
      <button type="button" @click="removeTag(index)" class="tag-remove">
        <X size="14" />
      </button>
    </div>
    <input
      type="text"
      v-model="inputValue"
      @keydown.enter.prevent="addTag"
      @keydown.backspace="onBackspace"
      :placeholder="tags.length === 0 ? placeholder : ''"
      class="tag-input-field"
    />
  </div>
</template>

<style scoped>
.tag-input-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  background: white;
  min-height: 42px; /* Match standard input height roughly */
}

.tag-input-container:focus-within {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.tag-pill {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background-color: #e0f2fe;
  color: #0284c7;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.tag-remove {
  background: none;
  border: none;
  color: #0284c7;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
}

.tag-remove:hover {
  color: #0369a1;
}

.tag-input-field {
  border: none;
  outline: none;
  flex-grow: 1;
  min-width: 120px;
  font-size: 1rem;
  padding: 0;
  margin: 0;
}
</style>
