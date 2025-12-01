<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserPlus } from 'lucide-vue-next'
import api from '../services/api'

const router = useRouter()

const formData = ref({
  full_name: '',
  email: '',
  username: '',
  password: '',
  phone_number: '',
  birth_date: '',
  hn: '',
  citizen_id: '',
})

const error = ref('')
const loading = ref(false)

const handleRegister = async () => {
  error.value = ''
  loading.value = true

  try {
    await api.post('/register', {
      full_name: formData.value.full_name,
      email: formData.value.email,
      username: formData.value.username,
      password: formData.value.password,
      phone_number: formData.value.phone_number,
      birth_date: new Date(formData.value.birth_date).toISOString(),
      hn: formData.value.hn,
      citizen_id: formData.value.citizen_id,
    })

    alert('ลงทะเบียนสำเร็จ! กรุณาเข้าสู่ระบบ')
    router.push('/login')
  } catch (err) {
    error.value = err.response?.data?.error || 'เกิดข้อผิดพลาดในการลงทะเบียน'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="logo-section">
          <h1>🌿 Pregnanzy </h1>
          <p>สมุดฝากครรภ์ดิจิทัล</p>
        </div>

        <form @submit.prevent="handleRegister" class="register-form">
          <h2>ลงทะเบียน</h2>

          <div class="form-row">
            <div class="form-group">
              <label>ชื่อ-นามสกุล *</label>
              <input type="text" v-model="formData.full_name" required />
            </div>
            <div class="form-group">
              <label>วันเกิด *</label>
              <input type="date" v-model="formData.birth_date" required />
            </div>
          </div>

          <div class="form-group">
            <label>อีเมล *</label>
            <input type="email" v-model="formData.email" required />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>เบอร์โทรศัพท์ *</label>
              <input type="tel" v-model="formData.phone_number" required />
            </div>
            <div class="form-group">
              <label>HN *</label>
              <input type="text" v-model="formData.hn" required />
            </div>
          </div>

          <div class="form-group">
            <label>เลขบัตรประชาชน *</label>
            <input type="text" v-model="formData.citizen_id" maxlength="13" required />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>ชื่อผู้ใช้ *</label>
              <input type="text" v-model="formData.username" required />
            </div>
            <div class="form-group">
              <label>รหัสผ่าน *</label>
              <input type="password" v-model="formData.password" required />
            </div>
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <button type="submit" class="btn-register" :disabled="loading">
            <UserPlus size="18" />
            {{ loading ? 'กำลังลงทะเบียน...' : 'ลงทะเบียน' }}
          </button>

          <div class="login-link">
            มีบัญชีแล้ว? <router-link to="/login">เข้าสู่ระบบ</router-link>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--color-bg) 0%, #e8f5e9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.register-container {
  width: 100%;
  max-width: 600px;
}

.register-card {
  background: var(--color-card-bg);
  border-radius: 1rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.logo-section {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-section h1 {
  font-size: 2rem;
  color: var(--color-text);
  margin-bottom: 0.5rem;
}

.logo-section p {
  color: var(--color-text-light);
  font-size: 0.875rem;
}

.register-form h2 {
  font-size: 1.5rem;
  color: var(--color-text);
  margin-bottom: 1.5rem;
  text-align: center;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 0.5rem;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(163, 230, 53, 0.1);
}

.error-message {
  background-color: #fee2e2;
  color: #991b1b;
  padding: 0.75rem;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}

.btn-register {
  width: 100%;
  padding: 0.875rem;
  background-color: var(--color-primary);
  color: var(--color-text);
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: background-color 0.2s;
  margin-top: 0.5rem;
}

.btn-register:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-register:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.login-link a {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 600;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>
