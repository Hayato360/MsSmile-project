<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { LogIn } from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  const success = await authStore.login(username.value, password.value)

  loading.value = false

  if (success) {
    router.push('/')
  } else {
    error.value = 'ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง'
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="logo-section">
          <h1>🌿 Healthcare Center</h1>
          <p>สมุดฝากครรภ์ดิจิทัล</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form">
          <h2>เข้าสู่ระบบ</h2>

          <div class="form-group">
            <label>ชื่อผู้ใช้</label>
            <input type="text" v-model="username" placeholder="กรอกชื่อผู้ใช้" required />
          </div>

          <div class="form-group">
            <label>รหัสผ่าน</label>
            <input type="password" v-model="password" placeholder="กรอกรหัสผ่าน" required />
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <button type="submit" class="btn-login" :disabled="loading">
            <LogIn size="18" />
            {{ loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
          </button>

          <div class="register-link">
            ยังไม่มีบัญชี? <router-link to="/register">ลงทะเบียน</router-link>
          </div>

          <div class="demo-credentials">
            <p><strong>บัญชีทดสอบ:</strong></p>
            <p>Username: woman</p>
            <p>Password: 123456</p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--color-bg) 0%, #e8f5e9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
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

.login-form h2 {
  font-size: 1.5rem;
  color: var(--color-text);
  margin-bottom: 1.5rem;
  text-align: center;
}

.form-group {
  margin-bottom: 1.5rem;
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

.btn-login {
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
}

.btn-login:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.register-link a {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 600;
}

.register-link a:hover {
  text-decoration: underline;
}

.demo-credentials {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f0fdf4;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  color: var(--color-text-light);
  text-align: center;
}

.demo-credentials p {
  margin: 0.25rem 0;
}
</style>
