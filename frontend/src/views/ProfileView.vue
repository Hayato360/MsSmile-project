<script setup>
import { ref, onMounted } from 'vue'
import { UserCircle, Save, Edit, User, FileHeart } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const loading = ref(true)
const saving = ref(false)

// Edit states
const isEditingPersonal = ref(false)
const isEditingMedical = ref(false)
const isEditingHusband = ref(false)

// Forms
const personalForm = ref({
  full_name: '',
  birth_date: '',
  age: '',
  citizen_id: '',
  phone_number: '',
  email: '',
})

const medicalForm = ref({
  ChronicDiseases: '',
  HeartDisease: false,
  Thyroid: false,
  OtherDiseases: '',
  SurgeryHistory: '',
  OtherSurgery: '',
  GeneticDiseases: '',
  DrugAllergies: '',
  FamilyHistoryHT: false,
  FamilyHistoryDiabetes: false,
  FamilyHistoryThalassemia: false,
  FamilyHistoryCongenital: false,
  OtherFamilyHistory: '',
  ContraceptionBeforeMethod: '',
  ContraceptionBeforeDuration: '',
  ContraceptionLastMethod: '',
  ContraceptionLastDuration: '',
  MenstrualCycle: 28,
  MenstrualDuration: 5,
  MenstrualCondition: 'ปกติ',
})

const husbandForm = ref({
  full_name: '',
  age: '',
  citizen_id: '',
  phone_number: '',
  email: '',
})

const initForms = () => {
  const user = authStore.user

  // Personal
  if (user) {
    personalForm.value = {
      full_name: user.full_name || '',
      birth_date: user.birth_date ? user.birth_date.split('T')[0] : '',
      age: user.age || '',
      citizen_id: user.citizen_id || '',
      phone_number: user.phone_number || '',
      email: user.email || '',
    }

    // Medical
    if (user.MedicalHistories && user.MedicalHistories.length > 0) {
      const history = user.MedicalHistories[0]
      medicalForm.value = { ...history }
    }

    // Husband
    if (user.Husband) {
      husbandForm.value = {
        full_name: user.Husband.FullName || '',
        age: user.Husband.Age || '',
        citizen_id: user.Husband.CitizenID || '',
        phone_number: user.Husband.PhoneNumber || '',
        email: user.Husband.Email || '',
      }
    }
  }
}

onMounted(async () => {
  await authStore.fetchMe()
  initForms()
  loading.value = false
})

// Personal Actions
const savePersonal = async () => {
  saving.value = true
  try {
    await api.put('/profile/personal', personalForm.value)
    await authStore.fetchMe()
    initForms()
    isEditingPersonal.value = false
    alert('บันทึกข้อมูลส่วนตัวสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก')
  } finally {
    saving.value = false
  }
}

const cancelEditPersonal = () => {
  initForms()
  isEditingPersonal.value = false
}

// Medical Actions
const saveMedical = async () => {
  saving.value = true
  try {
    await api.put('/profile/medical-history', medicalForm.value)
    await authStore.fetchMe()
    initForms()
    isEditingMedical.value = false
    alert('บันทึกประวัติสุขภาพสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก')
  } finally {
    saving.value = false
  }
}

const cancelEditMedical = () => {
  initForms()
  isEditingMedical.value = false
}

// Husband Actions
const saveHusband = async () => {
  saving.value = true
  try {
    await api.put('/profile/husband', husbandForm.value)
    await authStore.fetchMe()
    initForms()
    isEditingHusband.value = false
    alert('บันทึกข้อมูลคู่สมรสสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก')
  } finally {
    saving.value = false
  }
}

const cancelEditHusband = () => {
  initForms()
  isEditingHusband.value = false
}
</script>

<template>
  <div class="profile-view">
    <header class="page-header">
      <h2>ข้อมูลส่วนตัว</h2>
      <p>จัดการข้อมูลส่วนตัว ประวัติสุขภาพ และข้อมูลครอบครัว</p>
    </header>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else class="cards-container">
      <!-- Personal Information Card -->
      <div class="card">
        <div class="card-header">
          <User size="24" />
          <h3>ข้อมูลส่วนตัว</h3>
        </div>

        <div v-if="!isEditingPersonal" class="display-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">ชื่อ-นามสกุล</span>
              <span class="value">{{ personalForm.full_name }}</span>
            </div>
            <div class="info-item">
              <span class="label">วันเกิด</span>
              <span class="value">{{ personalForm.birth_date || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">อายุ</span>
              <span class="value">{{ personalForm.age }} ปี</span>
            </div>
            <div class="info-item">
              <span class="label">เลขบัตรประชาชน</span>
              <span class="value">{{ personalForm.citizen_id }}</span>
            </div>
            <div class="info-item">
              <span class="label">เบอร์โทรศัพท์</span>
              <span class="value">{{ personalForm.phone_number || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">อีเมล</span>
              <span class="value">{{ personalForm.email || '-' }}</span>
            </div>
          </div>
          <button @click="isEditingPersonal = true" class="btn-edit">
            <Edit size="18" />
            แก้ไขข้อมูล
          </button>
        </div>

        <form v-else @submit.prevent="savePersonal" class="form-content">
          <div class="form-grid">
            <div class="form-group">
              <label>ชื่อ-นามสกุล *</label>
              <input type="text" v-model="personalForm.full_name" required />
            </div>
            <div class="form-group">
              <label>วันเกิด</label>
              <input type="date" v-model="personalForm.birth_date" />
            </div>
            <div class="form-group">
              <label>เลขบัตรประชาชน *</label>
              <input type="text" v-model="personalForm.citizen_id" required maxlength="13" />
            </div>
            <div class="form-group">
              <label>เบอร์โทรศัพท์</label>
              <input type="tel" v-model="personalForm.phone_number" />
            </div>
            <div class="form-group full-width">
              <label>อีเมล</label>
              <input type="email" v-model="personalForm.email" />
            </div>
          </div>
          <div class="form-actions">
            <button type="button" @click="cancelEditPersonal" class="btn-cancel">ยกเลิก</button>
            <button type="submit" :disabled="saving" class="btn-save">
              <Save size="18" />
              บันทึก
            </button>
          </div>
        </form>
      </div>

      <!-- Medical History Card -->
      <div class="card">
        <div class="card-header">
          <FileHeart size="24" />
          <h3>ประวัติสุขภาพ</h3>
        </div>

        <div v-if="!isEditingMedical" class="display-content">
          <div class="info-grid">
            <div class="info-item full-width">
              <span class="label">โรคประจำตัว</span>
              <span class="value">{{ medicalForm.ChronicDiseases || 'ไม่มี' }}</span>
            </div>
            <div class="info-item full-width">
              <span class="label">ประวัติการผ่าตัด</span>
              <span class="value">{{ medicalForm.SurgeryHistory || 'ไม่มี' }}</span>
            </div>
            <div class="info-item full-width">
              <span class="label">ประวัติแพ้ยา/อาหาร</span>
              <span class="value">{{ medicalForm.DrugAllergies || 'ไม่มี' }}</span>
            </div>
            <div class="info-item full-width">
              <span class="label">โรคทางพันธุกรรม</span>
              <span class="value">{{ medicalForm.GeneticDiseases || 'ไม่มี' }}</span>
            </div>
          </div>
          <button @click="isEditingMedical = true" class="btn-edit">
            <Edit size="18" />
            แก้ไขข้อมูล
          </button>
        </div>

        <form v-else @submit.prevent="saveMedical" class="form-content">
          <div class="form-grid">
            <div class="form-group full-width">
              <label>โรคประจำตัว</label>
              <textarea v-model="medicalForm.ChronicDiseases" rows="2"></textarea>
            </div>
            <div class="form-group full-width">
              <label>ประวัติการผ่าตัด</label>
              <textarea v-model="medicalForm.SurgeryHistory" rows="2"></textarea>
            </div>
            <div class="form-group full-width">
              <label>ประวัติแพ้ยา/อาหาร</label>
              <textarea v-model="medicalForm.DrugAllergies" rows="2"></textarea>
            </div>
            <div class="form-group full-width">
              <label>โรคทางพันธุกรรม</label>
              <textarea v-model="medicalForm.GeneticDiseases" rows="2"></textarea>
            </div>
          </div>
          <div class="form-actions">
            <button type="button" @click="cancelEditMedical" class="btn-cancel">ยกเลิก</button>
            <button type="submit" :disabled="saving" class="btn-save">
              <Save size="18" />
              บันทึก
            </button>
          </div>
        </form>
      </div>

      <!-- Husband Information Card -->
      <div class="card">
        <div class="card-header">
          <UserCircle size="24" />
          <h3>ข้อมูลคู่สมรส</h3>
        </div>

        <div v-if="!isEditingHusband" class="display-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">ชื่อ-นามสกุล</span>
              <span class="value">{{ husbandForm.full_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">อายุ</span>
              <span class="value">{{ husbandForm.age ? husbandForm.age + ' ปี' : '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">เลขบัตรประชาชน</span>
              <span class="value">{{ husbandForm.citizen_id || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">เบอร์โทรศัพท์</span>
              <span class="value">{{ husbandForm.phone_number || '-' }}</span>
            </div>
          </div>
          <button @click="isEditingHusband = true" class="btn-edit">
            <Edit size="18" />
            แก้ไขข้อมูล
          </button>
        </div>

        <form v-else @submit.prevent="saveHusband" class="form-content">
          <div class="form-grid">
            <div class="form-group">
              <label>ชื่อ-นามสกุล *</label>
              <input type="text" v-model="husbandForm.full_name" required />
            </div>
            <div class="form-group">
              <label>อายุ *</label>
              <input type="number" v-model.number="husbandForm.age" required />
            </div>
            <div class="form-group">
              <label>เลขบัตรประชาชน *</label>
              <input type="text" v-model="husbandForm.citizen_id" required maxlength="13" />
            </div>
            <div class="form-group">
              <label>เบอร์โทรศัพท์</label>
              <input type="tel" v-model="husbandForm.phone_number" />
            </div>
            <div class="form-group full-width">
              <label>อีเมล</label>
              <input type="email" v-model="husbandForm.email" />
            </div>
          </div>
          <div class="form-actions">
            <button type="button" @click="cancelEditHusband" class="btn-cancel">ยกเลิก</button>
            <button type="submit" :disabled="saving" class="btn-save">
              <Save size="18" />
              บันทึก
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-view {
  max-width: 800px;
  margin: 0 auto;
  padding-bottom: 4rem;
}
.page-header {
  margin-bottom: 2rem;
}
.page-header h2 {
  margin: 0;
  color: var(--color-text);
}
.page-header p {
  color: var(--color-text-light);
  margin: 0.5rem 0 0;
}
.cards-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
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
  font-size: 1.25rem;
}
.form-content,
.display-content {
  padding: 2rem;
}
.form-grid,
.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.form-group,
.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.full-width {
  grid-column: 1 / -1;
}
.label {
  font-size: 0.875rem;
  color: var(--color-text-light);
}
.value {
  font-size: 1rem;
  color: var(--color-text);
  font-weight: 500;
}
.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text);
}
.form-group input,
.form-group textarea {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  font-size: 1rem;
}
.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--color-primary);
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
}
.btn-save {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
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
@media (max-width: 640px) {
  .form-grid,
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
