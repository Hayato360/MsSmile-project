<script setup>
import { ref, onMounted } from 'vue'
import { UserCircle, Save, Edit, User, FileHeart, Syringe } from 'lucide-vue-next'
import TagInput from '../components/TagInput.vue'
import VaccinationCard from '../components/VaccinationCard.vue'
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

// Medical History Existence Flags
const hasChronicDiseases = ref(false)
const hasSurgeryHistory = ref(false)
const hasDrugAllergies = ref(false)
const hasGeneticDiseases = ref(false)
const hasFamilyHistory = ref(false)

const husbandForm = ref({
  full_name: '',
  age: '',
  citizen_id: '',
  phone_number: '',
  email: '',
})

// Vaccination Data
const vaccineTypes = ref([])
const vaccinationForms = ref({})
const vaccinations = ref([])

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

      // Initialize Flags
      hasChronicDiseases.value = !!history.ChronicDiseases
      hasSurgeryHistory.value = !!history.SurgeryHistory
      hasDrugAllergies.value = !!history.DrugAllergies
      hasGeneticDiseases.value = !!history.GeneticDiseases
      hasFamilyHistory.value = !!history.OtherFamilyHistory
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
  
  // Load Vaccinations
  try {
      if (authStore.user) {
        const vTypeRes = await api.get('/vaccine-types')
        vaccineTypes.value = vTypeRes.data.data || []
        
        // Init empty forms
        vaccineTypes.value.forEach(type => {
            vaccinationForms.value[type.ID] = {
                VaccineTypeID: type.ID,
                IsPreviouslyVaccinated: false,
                PreviousDoses: 0,
                LastPreviousDateYear: null,
                Dose1DateDuringPreg: null,
                Dose2DateDuringPreg: null,
                Dose3DateDuringPreg: null,
                IsHistoryUnknown: false,
                ReasonForNotVaccinating: '',
                Remarks: '',
            }
        })

        const vacRes = await api.get(`/vaccinations/pregnant-woman/${authStore.user.ID}`)
        vaccinations.value = vacRes.data || []
        
        // Map data
        vaccinations.value.forEach(v => {
             if (vaccinationForms.value[v.VaccineTypeID]) {
                const vData = { ...v }
                if (vData.LastPreviousDateYear) vData.LastPreviousDateYear = vData.LastPreviousDateYear.split('T')[0]
                if (vData.Dose1DateDuringPreg) vData.Dose1DateDuringPreg = vData.Dose1DateDuringPreg.split('T')[0]
                if (vData.Dose2DateDuringPreg) vData.Dose2DateDuringPreg = vData.Dose2DateDuringPreg.split('T')[0]
                if (vData.Dose3DateDuringPreg) vData.Dose3DateDuringPreg = vData.Dose3DateDuringPreg.split('T')[0]
                vaccinationForms.value[v.VaccineTypeID] = { ...vaccinationForms.value[v.VaccineTypeID], ...vData }
             }
        })
      }
  } catch (e) {
      console.error("Error loading vaccinations", e)
  }

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
    // Clear fields if "No" is selected
    const payload = { ...medicalForm.value }
    if (!hasChronicDiseases.value) payload.ChronicDiseases = ''
    if (!hasSurgeryHistory.value) payload.SurgeryHistory = ''
    if (!hasDrugAllergies.value) payload.DrugAllergies = ''
    if (!hasGeneticDiseases.value) payload.GeneticDiseases = ''
    if (!hasFamilyHistory.value) payload.OtherFamilyHistory = ''

    await api.put('/profile/medical-history', payload)
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
const parseTags = (text) => {
  if (!text) return []
  // Split by newline or comma, then trim and filter empty
  return text.split(/[,\n]/).map(t => t.trim()).filter(t => t.length > 0)
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
          <div class="info-section">
            <h4>ประวัติการเจ็บป่วย</h4>
            <div class="info-grid">
              <div class="info-item full-width">
                <span class="label">โรคประจำตัว</span>
                <div class="tags">
                  <template v-if="medicalForm.ChronicDiseases">
                    <span 
                      v-for="(tag, index) in parseTags(medicalForm.ChronicDiseases)" 
                      :key="'chronic-'+index" 
                      class="tag"
                    >
                      {{ tag }}
                    </span>
                  </template>
                  <template v-if="medicalForm.OtherDiseases">
                    <span 
                      v-for="(tag, index) in parseTags(medicalForm.OtherDiseases)" 
                      :key="'other-'+index" 
                      class="tag"
                    >
                      {{ tag }}
                    </span>
                  </template>
                  <span
                    v-if="!medicalForm.ChronicDiseases && !medicalForm.OtherDiseases"
                    class="text-muted"
                    >-</span
                  >
                </div>
              </div>
              <div class="info-item full-width">
                <span class="label">ประวัติผ่าตัด</span>
                <span class="value">{{ medicalForm.SurgeryHistory || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="label">ประวัติแพ้ยา/อาหาร</span>
                <span class="value">{{ medicalForm.DrugAllergies || '-' }}</span>
              </div>
            </div>
          </div>

          <div class="divider"></div>

          <div class="info-section">
            <h4>ประวัติครอบครัว</h4>
            <div class="tags">
              <span v-if="medicalForm.OtherFamilyHistory" class="tag">{{
                medicalForm.OtherFamilyHistory
              }}</span>
              <span v-if="!medicalForm.OtherFamilyHistory" class="text-muted"
                >ไม่มีประวัติระบุ</span
              >
            </div>
          </div>

          <button @click="isEditingMedical = true" class="btn-edit">
            <Edit size="18" />
            แก้ไขข้อมูล
          </button>
        </div>

        <form v-else @submit.prevent="saveMedical" class="form-content">
          <!-- Illness History -->
          <div class="form-section">
            <h4 class="section-title">ประวัติการเจ็บป่วย</h4>
            <div class="form-grid">
              <div class="form-group full-width">
                <label>โรคประจำตัว</label>
                <div class="radio-group">
                  <label class="radio-label">
                    <input type="radio" :value="false" v-model="hasChronicDiseases" />
                    ไม่มี
                  </label>
                  <label class="radio-label">
                    <input type="radio" :value="true" v-model="hasChronicDiseases" />
                    มี
                  </label>
                </div>
                <TagInput
                  v-if="hasChronicDiseases"
                  v-model="medicalForm.ChronicDiseases"
                  placeholder="พิมพ์แล้วกด Enter เพื่อเพิ่มโรค"
                />
              </div>

              <div class="form-group full-width">
                <label>ประวัติการผ่าตัด</label>
                <div class="radio-group">
                  <label class="radio-label">
                    <input type="radio" :value="false" v-model="hasSurgeryHistory" />
                    ไม่มี
                  </label>
                  <label class="radio-label">
                    <input type="radio" :value="true" v-model="hasSurgeryHistory" />
                    มี
                  </label>
                </div>
                <TagInput
                  v-if="hasSurgeryHistory"
                  v-model="medicalForm.SurgeryHistory"
                  placeholder="พิมพ์แล้วกด Enter เพื่อเพิ่มประวัติ"
                />
              </div>

              <div class="form-group full-width">
                <label>ประวัติแพ้ยา/อาหาร</label>
                <div class="radio-group">
                  <label class="radio-label">
                    <input type="radio" :value="false" v-model="hasDrugAllergies" />
                    ไม่มี
                  </label>
                  <label class="radio-label">
                    <input type="radio" :value="true" v-model="hasDrugAllergies" />
                    มี
                  </label>
                </div>
                <TagInput
                  v-if="hasDrugAllergies"
                  v-model="medicalForm.DrugAllergies"
                  placeholder="พิมพ์แล้วกด Enter เพื่อเพิ่มประวัติแพ้"
                />
              </div>

              <div class="form-group full-width">
                <label>โรคทางพันธุกรรม</label>
                <div class="radio-group">
                  <label class="radio-label">
                    <input type="radio" :value="false" v-model="hasGeneticDiseases" />
                    ไม่มี
                  </label>
                  <label class="radio-label">
                    <input type="radio" :value="true" v-model="hasGeneticDiseases" />
                    มี
                  </label>
                </div>
                <TagInput
                  v-if="hasGeneticDiseases"
                  v-model="medicalForm.GeneticDiseases"
                  placeholder="พิมพ์แล้วกด Enter เพื่อเพิ่มโรค"
                />
              </div>
            </div>
          </div>

          <div class="divider"></div>

          <!-- Family History -->
          <div class="form-section">
            <h4 class="section-title">ประวัติครอบครัว</h4>
            <div class="form-group full-width mt-3">
              <label>ประวัติครอบครัว</label>
              <div class="radio-group">
                <label class="radio-label">
                  <input type="radio" :value="false" v-model="hasFamilyHistory" />
                  ไม่มี
                </label>
                <label class="radio-label">
                  <input type="radio" :value="true" v-model="hasFamilyHistory" />
                  มี
                </label>
              </div>
              <TagInput
                v-if="hasFamilyHistory"
                v-model="medicalForm.OtherFamilyHistory"
                placeholder="พิมพ์แล้วกด Enter เพื่อเพิ่มประวัติ"
              />
            </div>
          </div>

          <div class="divider"></div>

          <!-- Menstruation & Contraception -->
          <div class="form-section">
            <h4 class="section-title">ประวัติประจำเดือน & การคุมกำเนิด</h4>
            <div class="form-grid">
              <div class="form-group">
                <label>รอบประจำเดือน (วัน)</label>
                <input type="number" v-model.number="medicalForm.MenstrualCycle" placeholder="28" />
              </div>
              <div class="form-group">
                <label>จำนวนวันที่มีประจำเดือน</label>
                <input
                  type="number"
                  v-model.number="medicalForm.MenstrualDuration"
                  placeholder="5"
                />
              </div>
              <div class="form-group full-width">
                <label>ลักษณะประจำเดือน</label>
                <select v-model="medicalForm.MenstrualCondition">
                  <option value="ปกติ">ปกติ</option>
                  <option value="มาไม่สม่ำเสมอ">มาไม่สม่ำเสมอ</option>
                  <option value="ปวดท้องรุนแรง">ปวดท้องรุนแรง</option>
                </select>
              </div>

              <div class="form-group full-width">
                <label>การคุมกำเนิดก่อนตั้งครรภ์</label>
                <div class="input-group">
                  <input
                    type="text"
                    v-model="medicalForm.ContraceptionBeforeMethod"
                    placeholder="วิธีคุมกำเนิด"
                  />
                  <input
                    type="text"
                    v-model="medicalForm.ContraceptionBeforeDuration"
                    placeholder="ระยะเวลา"
                  />
                </div>
              </div>
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

      <!-- Vaccination Information Card -->
      <div class="card">
        <div class="card-header">
          <Syringe size="24" />
          <h3>ข้อมูลวัคซีน</h3>
        </div>
        <div class="display-content">
             <div class="vaccination-grid">
                <div v-for="type in vaccineTypes" :key="type.ID" class="vaccine-col">
                    <VaccinationCard 
                        :vaccine-type="type"
                        v-model="vaccinationForms[type.ID]"
                        :readonly="true"
                    />
                </div>
            </div>
             <p v-if="vaccineTypes.length === 0" class="text-center text-muted">ไม่พบข้อมูลวัคซีน</p>
        </div>
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

/* New Styles for Medical History Form */
.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-primary);
  margin-bottom: 1rem;
  border-left: 4px solid var(--color-primary);
  padding-left: 0.75rem;
}

.divider {
  height: 1px;
  background-color: var(--color-border);
  margin: 2rem 0;
}

.radio-group {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.95rem;
  color: var(--color-text);
}

.radio-label input[type='radio'] {
  width: 1.1rem;
  height: 1.1rem;
  accent-color: var(--color-primary);
  cursor: pointer;
}

.checkbox-group {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 0.5rem;
}

.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 1rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.95rem;
  color: var(--color-text);
}

.checkbox-label input[type='checkbox'] {
  width: 1.1rem;
  height: 1.1rem;
  accent-color: var(--color-primary);
  cursor: pointer;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  background-color: #ecfccb; /* Lime 100 */
  color: #3f6212; /* Lime 800 */
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.vaccination-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
}
.text-center { text-align: center; }
.text-muted { color: #6b7280; }

.text-muted {
  color: #9ca3af;
  font-style: italic;
}

.mt-2 {
  margin-top: 0.5rem;
}

.mt-3 {
  margin-top: 0.75rem;
}

.input-group {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

select {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
  font-size: 1rem;
  background-color: white;
  width: 100%;
}

select:focus {
  outline: none;
  border-color: var(--color-primary);
}
</style>
