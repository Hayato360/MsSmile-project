<script setup>
defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
  title: {
    type: String,
    default: 'ยืนยันการทำรายการ',
  },
  message: {
    type: String,
    default: 'คุณต้องการดำเนินการต่อหรือไม่?',
  },
  confirmText: {
    type: String,
    default: 'ตกลง',
  },
  cancelText: {
    type: String,
    default: 'ยกเลิก',
  },
  isDestructive: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['confirm', 'cancel'])
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="modal-overlay">
      <div class="modal-content">
        <h3 class="modal-title">{{ title }}</h3>
        <p class="modal-message">{{ message }}</p>
        <div class="modal-actions">
          <button @click="emit('cancel')" class="btn-cancel">{{ cancelText }}</button>
          <button
            @click="emit('confirm')"
            class="btn-confirm"
            :class="{ 'btn-destructive': isDestructive }"
          >
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
  padding: 1.5rem;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 400px;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  animation: fadeIn 0.2s ease-out;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.modal-message {
  color: #4b5563;
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.btn-cancel {
  padding: 0.5rem 1rem;
  background-color: white;
  border: 1px solid #d1d5db;
  color: #374151;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-cancel:hover {
  background-color: #f3f4f6;
}

.btn-confirm {
  padding: 0.5rem 1rem;
  background-color: var(--color-primary, #10b981);
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-confirm:hover {
  filter: brightness(0.9);
}

.btn-destructive {
  background-color: #ef4444;
}

.btn-destructive:hover {
  background-color: #dc2626;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
