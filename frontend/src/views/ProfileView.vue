<script setup>
import { ref, onMounted } from 'vue'
import { UserCircle, Save, Edit } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const loading = ref(true)
const saving = ref(false)
const isEditing = ref(true)
const hasData = ref(false)

const husbandForm = ref({
  full_name: '',
  age: '',
  citizen_id: '',
  phone_number: '',
  email: '',
})

onMounted(async () => {
  // Load existing husband data if available
  if (authStore.user?.Husband) {
    husbandForm.value = {
      full_name: authStore.user.Husband.FullName || '',
      age: authStore.user.Husband.Age || '',
      citizen_id: authStore.user.Husband.CitizenID || '',
      phone_number: authStore.user.Husband.PhoneNumber || '',
      email: authStore.user.Husband.Email || '',
    }
    hasData.value = true
    isEditing.value = false
  }
  loading.value = false
})

const cancelEdit = () => {
  if (authStore.user?.Husband) {
    husbandForm.value = {
      full_name: authStore.user.Husband.FullName || '',
      age: authStore.user.Husband.Age || '',
      citizen_id: authStore.user.Husband.CitizenID || '',
      phone_number: authStore.user.Husband.PhoneNumber || '',
      email: authStore.user.Husband.Email || '',
    }
  }
  isEditing.value = false
}

const saveHusband = async () => {
  saving.value = true
  try {
    await api.put('/profile/husband', husbandForm.value)
    alert('บันทึกข้อมูลสำเร็จ')

    // Reload user data
    await authStore.fetchMe()
    hasData.value = true
    isEditing.value = false
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก')
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="profile-view">
    <header class="page-header">
      <h2>ข้อมูลส่วนตัว</h2>
      <p>กรอกข้อมูลครอบครัวของคุณ</p>
    </header>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else class="card">
      <div class="card-header">
        <UserCircle size="24" />
        <h3>ข้อมูลคู่สมรส</h3>
      </div>

      <div v-if="!isEditing && hasData" class="display-content">
        <div class="info-grid">
          <div class="info-item">
            <span class="label">ชื่อ-นามสกุล</span>
            <span class="value">{{ husbandForm.full_name }}</span>
          </div>
          <div class="info-item">
            <span class="label">อายุ</span>
            <span class="value">{{ husbandForm.age }} ปี</span>
          </div>
          <div class="info-item">
            <span class="label">เลขบัตรประชาชน</span>
            <span class="value">{{ husbandForm.citizen_id }}</span>
          </div>
          <div class="info-item">
            <span class="label">เบอร์โทรศัพท์</span>
            <span class="value">{{ husbandForm.phone_number || '-' }}</span>
          </div>
          <div class="info-item full-width">
            <span class="label">อีเมล</span>
            <span class="value">{{ husbandForm.email || '-' }}</span>
          </div>
        </div>
        <button @click="isEditing = true" class="btn-edit">
          <Edit size="18" />
          แก้ไขข้อมูล
        </button>
      </div>

      <form v-else @submit.prevent="saveHusband" class="form-content">
        <div class="form-grid">
          <div class="form-group">
            <label>ชื่อ-นามสกุล *</label>
            <input
              type="text"
              v-model="husbandForm.full_name"
              required
              placeholder="นาย สมชาย ใจดี"
            />
          </div>

          <div class="form-group">
            <label>อายุ *</label>
            <input type="number" v-model.number="husbandForm.age" required placeholder="30" />
          </div>

          <div class="form-group">
            <label>เลขบัตรประชาชน *</label>
            <input
              type="text"
              v-model="husbandForm.citizen_id"
              required
              placeholder="1-2345-67890-12-3"
              maxlength="17"
            />
          </div>

          <div class="form-group">
            <label>เบอร์โทรศัพท์</label>
            <input type="tel" v-model="husbandForm.phone_number" placeholder="081-234-5678" />
          </div>

          <div class="form-group full-width">
            <label>อีเมล</label>
            <input type="email" v-model="husbandForm.email" placeholder="example@email.com" />
          </div>
        </div>

        <div class="form-actions">
          <button v-if="hasData" type="button" @click="cancelEdit" class="btn-cancel">
            ยกเลิก
          </button>
          <button type="submit" :disabled="saving" class="btn-save">
            <Save size="18" />
            {{ saving ? 'กำลังบันทึก...' : 'บันทึกข้อมูล' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.profile-view {
  max-width: 800px;
  margin: 0 auto;
}
.page-header {
  margin-bottom: 2rem;
}
.page-header h2 {
  margin: 0;
}
.page-header p {
  color: var(--color-text-light);
  margin: 0.5rem 0 0;
}
.card {
  background: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--color-primary);
  color: white;
}
.card-header h3 {
  margin: 0;
}
.form-content {
  padding: 2rem;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.form-group.full-width {
  grid-column: 1 / -1;
}
.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text);
}
.form-group input {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  font-size: 1rem;
}
.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
}
.btn-save {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: var(--color-text);
  border: none;
  border-radius: 0.375rem;
  font-weight: 600;
  cursor: pointer;
}
.btn-save:hover:not(:disabled) {
  background: var(--color-primary-hover);
}
.btn-save:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.loading {
  text-align: center;
  padding: 3rem;
}

.display-content {
  padding: 2rem;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item .label {
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.info-item .value {
  font-size: 1rem;
  color: var(--color-text);
  font-weight: 500;
}

.btn-edit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1.5rem;
  background: white;
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
  border-radius: 0.375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-edit:hover {
  background: #f0fdf4;
}

.form-actions {
  display: flex;
  gap: 1rem;
}

.btn-cancel {
  flex: 1;
  padding: 0.75rem 1.5rem;
  background: white;
  color: #6b7280;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-cancel:hover {
  background: #f9fafb;
  color: #374151;
}

.btn-save {
  flex: 1;
  /* existing styles... */
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  /* width: 100%; remove width: 100% to allow flex sharing */
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: var(--color-text);
  border: none;
  border-radius: 0.375rem;
  font-weight: 600;
  cursor: pointer;
}
</style>
