<script setup>
import { computed } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
  title: {
    type: String,
    default: 'แจ้งเตือน',
  },
  message: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    default: 'success', // success, error, info, warning
    validator: (value) => ['success', 'error', 'info', 'warning'].includes(value),
  },
  confirmText: {
    type: String,
    default: 'ตกลง',
  },
})

const emit = defineEmits(['close'])

const iconClass = computed(() => {
  return {
    success: '✓',
    error: '✕',
    info: 'ℹ',
    warning: '⚠',
  }[props.type]
})

const iconColor = computed(() => {
  return {
    success: '#10b981',
    error: '#ef4444',
    info: '#3b82f6',
    warning: '#f59e0b',
  }[props.type]
})
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="modal-overlay" @click.self="emit('close')">
      <div class="modal-content">
        <div class="icon-wrapper" :style="{ backgroundColor: iconColor }">
          <span class="icon">{{ iconClass }}</span>
        </div>
        <h3 class="modal-title">{{ title }}</h3>
        <p class="modal-message">{{ message }}</p>
        <div class="modal-actions">
          <button @click="emit('close')" class="btn-confirm" :style="{ backgroundColor: iconColor }">
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 0.75rem;
  width: 90%;
  max-width: 400px;
  box-shadow:
    0 10px 25px -5px rgba(0, 0, 0, 0.1),
    0 8px 10px -6px rgba(0, 0, 0, 0.1);
  animation: slideDown 0.3s ease-out;
  text-align: center;
}

.icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  animation: scaleIn 0.4s ease-out;
}

.icon {
  font-size: 2rem;
  color: white;
  font-weight: bold;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.75rem;
}

.modal-message {
  color: #6b7280;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.modal-actions {
  display: flex;
  justify-content: center;
}

.btn-confirm {
  padding: 0.625rem 2rem;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 1rem;
}

.btn-confirm:hover {
  filter: brightness(0.9);
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.btn-confirm:active {
  transform: translateY(0);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}
</style>
