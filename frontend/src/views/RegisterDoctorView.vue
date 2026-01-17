<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserPlus, Hospital, CreditCard, FileText, User, Mail, Phone } from 'lucide-vue-next'
import api from '../services/api'

const router = useRouter()

const formData = ref({
  full_name: '',
  email: '',
  phone_number: '',
  citizen_id: '',
  hospital_code: '',
  doctor_license_no: '',
  username: '',
  password: '',
})

const error = ref('')
const loading = ref(false)

const handleRegister = async () => {
  error.value = ''
  loading.value = true

  try {
    await api.post('/register/doctor', formData.value)

    alert('ลงทะเบียนสำเร็จ! กรุณาเข้าสู่ระบบ')
    router.push('/login')
  } catch (err) {
    error.value = err.response?.data?.error || 'เกิดข้อผิดพลาดในการลงทะเบียน'
  } finally {
    loading.value = false
  }
}

const fillTestCredentials = (docNumber) => {
  if (docNumber === 1) {
    formData.value.full_name = 'นพ.ทดสอบ หนึ่ง'
    formData.value.email = 'test1@hospital.th'
    formData.value.phone_number = '0812345678'
    formData.value.citizen_id = '1000000000001'
    formData.value.hospital_code = 'HOS001'
    formData.value.doctor_license_no = 'DOC00001'
  } else {
    formData.value.full_name = 'พญ.ทดสอบ สอง'
    formData.value.email = 'test2@hospital.th'
    formData.value.phone_number = '0898765432'
    formData.value.citizen_id = '2000000000002'
    formData.value.hospital_code = 'HOS001'
    formData.value.doctor_license_no = 'DOC00002'
  }
}
</script>

<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="logo-section">
          <h1>🏥 ลงทะเบียนแพทย์</h1>
          <p>ระบบสมุดฝากครรภ์ดิจิทัล</p>
        </div>

        <!-- Registration Form -->
        <form @submit.prevent="handleRegister" class="register-form">
          <h2>ลงทะเบียนแพทย์</h2>
          <p class="form-description">กรอกข้อมูลให้ครบถ้วนเพื่อลงทะเบียนเข้าใช้งานระบบ</p>

          <!-- ข้อมูลส่วนตัว -->
          <div class="section-title">📋 ข้อมูลส่วนตัว</div>

          <div class="form-group">
            <label>
              <User size="18" />
              ชื่อ-นามสกุล *
            </label>
            <input
              type="text"
              v-model="formData.full_name"
              placeholder="เช่น นพ.สมชาย ใจดี"
              required
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>
                <Mail size="18" />
                อีเมล *
              </label>
              <input
                type="email"
                v-model="formData.email"
                placeholder="doctor@hospital.th"
                required
              />
            </div>

            <div class="form-group">
              <label>
                <Phone size="18" />
                เบอร์โทรศัพท์ *
              </label>
              <input
                type="tel"
                v-model="formData.phone_number"
                placeholder="0812345678"
                required
              />
            </div>
          </div>

          <!-- ข้อมูลยืนยันตัวตน -->
          <div class="section-title">🏥 ข้อมูลยืนยันตัวตนแพทย์</div>

          <div class="form-group">
            <label>
              <CreditCard size="18" />
              เลขบัตรประชาชน *
            </label>
            <input
              type="text"
              v-model="formData.citizen_id"
              placeholder="เลขบัตรประชาชน 13 หลัก"
              maxlength="13"
              required
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>
                <Hospital size="18" />
                รหัสโรงพยาบาล *
              </label>
              <input
                type="text"
                v-model="formData.hospital_code"
                placeholder="เช่น HOS001"
                required
              />
            </div>

            <div class="form-group">
              <label>
                <FileText size="18" />
                รหัสใบอนุญาตประกอบวิชาชีพ *
              </label>
              <input
                type="text"
                v-model="formData.doctor_license_no"
                placeholder="เช่น DOC12345"
                required
              />
            </div>
          </div>

          <!-- ข้อมูลบัญชี -->
          <div class="section-title">🔐 ข้อมูลบัญชีผู้ใช้</div>

          <div class="form-row">
            <div class="form-group">
              <label>ชื่อผู้ใช้ *</label>
              <input type="text" v-model="formData.username" placeholder="ชื่อผู้ใช้" required />
            </div>

            <div class="form-group">
              <label>รหัสผ่าน *</label>
              <input
                type="password"
                v-model="formData.password"
                placeholder="รหัสผ่าน"
                required
              />
            </div>
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <button type="submit" class="btn-primary" :disabled="loading">
            <UserPlus size="18" />
            {{ loading ? 'กำลังลงทะเบียน...' : 'ลงทะเบียน' }}
          </button>

          <div class="demo-credentials">
            <p class="demo-title">🧪 <strong>ข้อมูลสำหรับทดสอบ</strong></p>
            <p class="demo-subtitle">คลิกเพื่อใส่ข้อมูลทดสอบอัตโนมัติ</p>

            <div class="demo-options">
              <button type="button" @click="fillTestCredentials(1)" class="demo-button">
                <strong>แพทย์ 1</strong><br />
                <small>นพ.ทดสอบ หนึ่ง</small>
              </button>

              <button type="button" @click="fillTestCredentials(2)" class="demo-button">
                <strong>แพทย์ 2</strong><br />
                <small>พญ.ทดสอบ สอง</small>
              </button>
            </div>
          </div>
        </form>

        <div class="login-link">
          มีบัญชีแล้ว? <router-link to="/login">เข้าสู่ระบบ</router-link><br />
          เป็นผู้ป่วย? <router-link to="/register">ลงทะเบียนผู้ป่วย</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
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

/* Step Indicator */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 2rem;
  padding: 0 2rem;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e0e0e0;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  transition: all 0.3s;
}

.step.active .step-number,
.step.completed .step-number {
  background: var(--color-primary);
  color: var(--color-text);
}

.step-label {
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.step.active .step-label {
  color: var(--color-text);
  font-weight: 600;
}

.step-line {
  flex: 1;
  height: 2px;
  background: #e0e0e0;
  margin: 0 1rem;
  margin-bottom: 1.5rem;
  transition: all 0.3s;
}

.step-line.active {
  background: var(--color-primary);
}

/* Form */
.register-form h2 {
  font-size: 1.5rem;
  color: var(--color-text);
  margin-bottom: 0.5rem;
  text-align: center;
}

.form-description {
  text-align: center;
  color: var(--color-text-light);
  font-size: 0.875rem;
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--color-border);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
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

.verified-info {
  background: #f0fdf4;
  padding: 1rem;
  border-radius: 0.5rem;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.verified-info p {
  margin: 0.25rem 0;
}

.error-message {
  background-color: #fee2e2;
  color: #991b1b;
  padding: 0.75rem;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
  font-size: 0.875rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.875rem 1.5rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 0.2s;
}

.btn-primary {
  width: 100%;
  background-color: var(--color-primary);
  color: var(--color-text);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #e0e0e0;
  color: var(--color-text);
}

.btn-secondary:hover {
  background: #d0d0d0;
}

.button-group {
  display: flex;
  gap: 1rem;
}

.button-group .btn-primary,
.button-group .btn-secondary {
  flex: 1;
}

.login-link {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-text-light);
  line-height: 1.8;
}

.login-link a {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 600;
}

.login-link a:hover {
  text-decoration: underline;
}

.demo-credentials {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #f0fdf4;
  border: 2px dashed #86efac;
  border-radius: 0.5rem;
  text-align: center;
}

.demo-title {
  font-size: 0.875rem;
  color: #166534;
  margin: 0 0 0.25rem 0;
}

.demo-subtitle {
  font-size: 0.75rem;
  color: #16a34a;
  margin: 0 0 1rem 0;
}

.demo-options {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.demo-button {
  padding: 0.75rem;
  background: white;
  border: 1px solid #86efac;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
  line-height: 1.5;
  font-size: 0.75rem;
  color: #166534;
}

.demo-button:hover {
  background: #dcfce7;
  border-color: #22c55e;
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(34, 197, 94, 0.2);
}

.demo-button:active {
  transform: translateY(0);
}

.demo-button strong {
  display: block;
  margin-bottom: 0.25rem;
  color: #15803d;
}

.demo-button small {
  font-size: 0.65rem;
  color: #16a34a;
  font-family: monospace;
}

</style>
