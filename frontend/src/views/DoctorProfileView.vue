<script setup>
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { Save, UserCog, Pencil, X } from 'lucide-vue-next'
import api from '../services/api'
import DoctorSidebar from '../components/DoctorSidebar.vue'

const authStore = useAuthStore()
const isEditing = ref(false)
const form = ref({
  full_name: '',
  phone_number: '',
  email: '',
})
const loading = ref(false)

const populateForm = (user) => {
  if (!user) return
  form.value = {
    full_name: user.full_name || user.FullName || '',
    phone_number: user.phone_number || user.PhoneNumber || '',
    email: user.email || user.Email || '',
  }
}

onMounted(() => {
  if (authStore.user) {
    populateForm(authStore.user)
  }
})

// Add watcher for page reload case where user data loads after mount
watch(
  () => authStore.user,
  (newUser) => {
    if (newUser) {
      populateForm(newUser)
    }
  },
  { immediate: true }
)

const toggleEdit = () => {
    if (isEditing.value) {
        // Canceling edit: reset form to current user data
        populateForm(authStore.user)
    }
    isEditing.value = !isEditing.value
}

const updateProfile = async () => {
  loading.value = true
  try {
    const response = await api.put('/doctor/profile', form.value)
    
    // Update local store
    authStore.user = { ...authStore.user, ...response.data.data }
    // Update localStorage if necessary or rely on store persistence
    localStorage.setItem('user', JSON.stringify(authStore.user))

    alert('บันทึกข้อมูลเรียบร้อยแล้ว')
    isEditing.value = false // Switch back to view mode
  } catch (error) {
    console.error('Error updating profile:', error)
    alert('เกิดข้อผิดพลาดในการบันทึกข้อมูล')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page-container">
    <DoctorSidebar />
    <main class="main-content">
      <header class="page-header">
        <h1>ข้อมูลส่วนตัว</h1>
        <p>จัดการข้อมูลส่วนตัวของคุณ</p>
      </header>

      <div class="content-wrapper">
        <div class="card profile-card">
          <div class="card-header">
            <div class="header-title">
                <UserCog size="24" />
                <h2>แอดมิน / แพทย์</h2>
            </div>
            <button v-if="!isEditing" @click="toggleEdit" class="btn-edit-toggle">
                <Pencil size="18" /> แก้ไขข้อมูล
            </button>
          </div>
          
          <!-- View Mode -->
          <div v-if="!isEditing" class="view-mode">
            <div class="info-group">
                <label>ชื่อ-นามสกุล</label>
                <p>{{ authStore.user?.full_name || authStore.user?.FullName || '-' }}</p>
            </div>
            <div class="info-group">
                <label>เบอร์โทรศัพท์</label>
                <p>{{ authStore.user?.phone_number || authStore.user?.PhoneNumber || '-' }}</p>
            </div>
            <div class="info-group">
                <label>อีเมล</label>
                <p>{{ authStore.user?.email || authStore.user?.Email || '-' }}</p>
            </div>
          </div>

          <!-- Edit Mode -->
          <form v-else @submit.prevent="updateProfile" class="form-layout">
             <div class="form-group">
              <label>ชื่อ-นามสกุล</label>
              <input type="text" v-model="form.full_name" required placeholder="ชื่อ-นามสกุล" />
            </div>

            <div class="form-group">
              <label>เบอร์โทรศัพท์</label>
              <input type="tel" v-model="form.phone_number" placeholder="08x-xxx-xxxx" />
            </div>

            <div class="form-group">
              <label>อีเมล</label>
              <input type="email" v-model="form.email" placeholder="doctor@example.com" />
            </div>

            <div class="form-actions">
              <button type="button" @click="toggleEdit" class="btn-cancel" :disabled="loading">
                ยกเลิก
              </button>
              <button type="submit" class="btn-save" :disabled="loading">
                <span v-if="loading">กำลังบันทึก...</span>
                <span v-else class="flex-center">
                    <Save size="18" /> บันทึกการเปลี่ยนแปลง
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.page-container {
  display: flex;
  min-height: 100vh;
  background-color: #f8fafc;
}

.main-content {
  margin-left: var(--sidebar-width);
  flex: 1;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 1.875rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.5rem;
}

.page-header p {
  color: #64748b;
}

.content-wrapper {
  max-width: 600px;
}

.card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e2e8f0;
}

.header-title {
    display: flex;
    align-items: center;
    gap: 1rem;
    color: #0ea5e9;
}

.card-header h2 {
    font-size: 1.25rem;
    color: #0f172a;
    font-weight: 600;
}

.view-mode {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.info-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
}

.info-group label {
    font-size: 0.875rem;
    color: #64748b;
    font-weight: 500;
}

.info-group p {
    font-size: 1.125rem;
    color: #1e293b;
    font-weight: 500;
}

.form-layout {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #334155;
}

.form-group input {
  padding: 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #0ea5e9;
  box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.1);
}

.form-actions {
    display: flex;
    gap: 1rem;
}

.btn-edit-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background-color: white;
    border: 1px solid #cbd5e1;
    border-radius: 0.5rem;
    color: #475569;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-edit-toggle:hover {
    background-color: #f1f5f9;
    color: #0ea5e9;
    border-color: #0ea5e9;
}

.btn-cancel {
    flex: 1;
    padding: 0.75rem;
    background-color: white;
    border: 1px solid #cbd5e1;
    border-radius: 0.5rem;
    color: #475569;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-cancel:hover {
    background-color: #f1f5f9;
    border-color: #94a3b8;
}

.btn-save {
  flex: 2;
  background-color: #0ea5e9;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  justify-content: center;
  align-items: center;
}

.btn-save:hover {
  background-color: #0284c7;
}

.btn-save:disabled {
  background-color: #94a3b8;
  cursor: not-allowed;
}

.flex-center {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}
</style>
