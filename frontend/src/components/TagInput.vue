<script setup>
import { ref, watch, nextTick } from 'vue'
import { X } from 'lucide-vue-next'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const tags = ref([])
const inputValue = ref('')
const inputRef = ref(null)

// Parse initial string into tags
watch(
  () => props.modelValue,
  (newValue) => {
    // Only update if the joined string differs to avoid recursive loops during editing
    // or simply initialization. 
    // We compare current tags joined vs new value to decide if we need to sync from parent.
    const currentJoined = tags.value.join(',')
    if (newValue !== currentJoined) {
        if (!newValue) {
            tags.value = []
        } else {
            tags.value = newValue.split(',').map(t => t.trim()).filter(t => t)
        }
    }
  },
  { immediate: true }
)

const emitUpdate = () => {
  emit('update:modelValue', tags.value.join(','))
}

const addTag = () => {
  const val = inputValue.value.trim()
  if (val) {
    if (!tags.value.includes(val)) {
        tags.value.push(val)
        emitUpdate()
    }
    inputValue.value = ''
  }
}

const removeTag = (index) => {
  tags.value.splice(index, 1)
  emitUpdate()
}

const handleKeydown = (e) => {
  if (e.key === 'Enter' || e.key === ',') {
    e.preventDefault()
    addTag()
  } else if (e.key === 'Backspace' && inputValue.value === '' && tags.value.length > 0) {
    // Optional: remove last tag on backspace if input is empty
    removeTag(tags.value.length - 1)
  }
}

const handleBlur = () => {
    addTag()
}

const focusInput = () => {
    inputRef.value?.focus()
}
</script>

<template>
  <div class="tag-input-container" @click="focusInput">
    <div class="tags-wrapper">
      <span v-for="(tag, index) in tags" :key="index" class="tag-bubble">
        {{ tag }}
        <button type="button" @click.stop="removeTag(index)" class="remove-btn">
          <X size="14" />
        </button>
      </span>
      <input
        ref="inputRef"
        v-model="inputValue"
        type="text"
        class="tag-input-field"
        :placeholder="tags.length === 0 ? placeholder : ''"
        @keydown="handleKeydown"
        @blur="handleBlur"
      />
    </div>
  </div>
</template>

<style scoped>
.tag-input-container {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  background-color: white;
  cursor: text;
  min-height: 42px; /* standard input height */
  display: flex;
  align-items: center;
}

.tag-input-container:focus-within {
  border-color: var(--color-primary);
  outline: 1px solid var(--color-primary);
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  width: 100%;
}

.tag-bubble {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background-color: #ecfccb; /* Lime 100 */
  color: #3f6212; /* Lime 800 */
  padding: 0.125rem 0.5rem;
  padding-right: 0.25rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
  white-space: nowrap;
}

.remove-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  color: #3f6212;
  opacity: 0.6;
  border-radius: 50%;
}

.remove-btn:hover {
  opacity: 1;
  background-color: rgba(0,0,0,0.05);
}

.tag-input-field {
  flex: 1;
  min-width: 120px;
  border: none;
  outline: none;
  font-size: 1rem;
  padding: 0;
  margin: 0;
  background: transparent;
}

.tag-input-field::placeholder {
  color: #9ca3af;
}
</style>
