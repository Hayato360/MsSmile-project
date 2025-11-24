<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Save, ArrowLeft, Calendar, FileText, Syringe, FlaskConical } from 'lucide-vue-next'
import api from '../services/api'
import router from '../router'

const route = useRoute()
const patient = ref(null)
const visits = ref([])
const loading = ref(true)
const activeTab = ref('visits') // visits, medical, vaccination, lab

const visitForm = ref({
  VisitDate: new Date().toISOString().split('T')[0],
  GestationalAge: '',
  Weight: '',
  BloodPressure: '',
  HeightFundus: '',
  FetalHeartSound: 'ปกติ',
  FetalMovement: 'ปกติ',
  UrineProtein: 'ไม่พบ',
  UrineSugar: 'ไม่พบ',
  Swelling: 'ไม่บวม',
  MedicalDiagnosis: '',
})

const medicalHistoryForm = ref({
  ChronicDiseases: '',
  GeneticDiseases: '',
  DrugAllergies: '',
  FamilyHistoryHT: false,
  FamilyHistoryCongenital: false,
  OtherFamilyHistory: '',
})

const vaccinationForm = ref({
  VaccineTypeID: 1,
  IsPreviouslyVaccinated: false,
  PreviousDoses: 0,
  LastPreviousDateYear: null,
  Dose1DateDuringPreg: null,
  Dose2DateDuringPreg: null,
  Remarks: '',
})

const labResultForm = ref({
  TestDate: new Date().toISOString().split('T')[0],
  Hct: '',
  Hb: '',
  HbTyping: '',
  OtherRemarks: '',
})

const pregnancyId = computed(() => {
  if (!patient.value?.Pregnancies || patient.value.Pregnancies.length === 0) return null
  return patient.value.Pregnancies[patient.value.Pregnancies.length - 1].ID
})

const pregnancyData = computed(() => {
  if (!patient.value?.Pregnancies || patient.value.Pregnancies.length === 0) return null
  return patient.value.Pregnancies[patient.value.Pregnancies.length - 1]
})

const calculateGA = () => {
  if (!pregnancyData.value?.LMP || !visitForm.value.VisitDate) return

  const lmp = new Date(pregnancyData.value.LMP)
  const visitDate = new Date(visitForm.value.VisitDate)
  const diffTime = Math.abs(visitDate - lmp)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  const weeks = Math.floor(diffDays / 7)

  visitForm.value.GestationalAge = weeks
}

const onVisitDateChange = () => {
  calculateGA()
}

onMounted(async () => {
  const patientId = route.params.id
  try {
    const patientRes = await api.get(`/doctor/patients/${patientId}`)
    patient.value = patientRes.data

    const visitsRes = await api.get(`/doctor/patient/${patientId}/visits`)
    visits.value = visitsRes.data || []

    // Load medical history
    const historyRes = await api.get(`/doctor/patient/${patientId}/medical-history`)
    if (historyRes.data.data) {
      medicalHistoryForm.value = historyRes.data.data
    }

    calculateGA()
  } catch (error) {
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
})

const newPregnancyLMP = ref(new Date().toISOString().split('T')[0])

const createPregnancy = async () => {
  try {
    const lmpDate = new Date(newPregnancyLMP.value)
    const edcDate = new Date(lmpDate)
    edcDate.setMonth(edcDate.getMonth() + 9)
    edcDate.setDate(edcDate.getDate() + 7)

    await api.post('/doctor/pregnancy', {
      PregnantWomanID: parseInt(route.params.id),
      PregnancyNo: 1,
      LMP: lmpDate.toISOString(),
      EDC: edcDate.toISOString(),
      PrePregnancyWeight: 50,
      Height: 160,
      PrePregnancyBMI: 19.5,
    })

    alert('สร้างข้อมูลการตั้งครรภ์สำเร็จ')
    const patientRes = await api.get(`/doctor/patients/${route.params.id}`)
    patient.value = patientRes.data
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
}

const saveVisit = async () => {
  if (!pregnancyId.value) {
    alert('คนไข้ยังไม่มีข้อมูลการตั้งครรภ์')
    return
  }

  try {
    await api.post('/doctor/antenatal-visit', {
      ...visitForm.value,
      PregnancyID: pregnancyId.value,
      Weight: parseFloat(visitForm.value.Weight),
      GestationalAge: parseInt(visitForm.value.GestationalAge),
      HeightFundus: parseFloat(visitForm.value.HeightFundus),
      VisitDate: new Date(visitForm.value.VisitDate).toISOString(),
    })

    alert('บันทึกข้อมูลสำเร็จ')
    const visitsRes = await api.get(`/doctor/patient/${route.params.id}/visits`)
    visits.value = visitsRes.data || []
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก')
  }
}

const saveMedicalHistory = async () => {
  try {
    await api.post('/doctor/medical-history', {
      ...medicalHistoryForm.value,
      PregnantWomanID: parseInt(route.params.id),
    })
    alert('บันทึกประวัติสุขภาพสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
}

const saveVaccination = async () => {
  try {
    await api.post('/doctor/vaccination', {
      ...vaccinationForm.value,
      PregnantWomanID: parseInt(route.params.id),
      LastPreviousDateYear: vaccinationForm.value.LastPreviousDateYear
        ? new Date(vaccinationForm.value.LastPreviousDateYear).toISOString()
        : null,
      Dose1DateDuringPreg: vaccinationForm.value.Dose1DateDuringPreg
        ? new Date(vaccinationForm.value.Dose1DateDuringPreg).toISOString()
        : null,
      Dose2DateDuringPreg: vaccinationForm.value.Dose2DateDuringPreg
        ? new Date(vaccinationForm.value.Dose2DateDuringPreg).toISOString()
        : null,
    })
    alert('บันทึกข้อมูลวัคซีนสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
}

const saveLabResult = async () => {
  if (!pregnancyId.value) {
    alert('คนไข้ยังไม่มีข้อมูลการตั้งครรภ์')
    return
  }

  try {
    await api.post('/doctor/lab-result', {
      ...labResultForm.value,
      PregnancyID: pregnancyId.value,
      Hct: parseFloat(labResultForm.value.Hct),
      Hb: parseFloat(labResultForm.value.Hb),
      TestDate: new Date(labResultForm.value.TestDate).toISOString(),
    })
    alert('บันทึกผลแล็บสำเร็จ')
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
}

const goBack = () => router.push('/doctor/dashboard')
const formatDate = (dateString) => new Date(dateString).toLocaleDateString('th-TH')
</script>

<template>
  <div class="doctor-patient-detail">
    <button @click="goBack" class="btn-back">
      <ArrowLeft size="18" />
      กลับ
    </button>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else-if="patient">
      <div class="card">
        <h2>ข้อมูลคนไข้</h2>
        <div class="info-grid">
          <div><span class="label">ชื่อ:</span> {{ patient.full_name }}</div>
          <div><span class="label">HN:</span> {{ patient.hn }}</div>
          <div><span class="label">อายุ:</span> {{ patient.age }} ปี</div>
          <div><span class="label">เบอร์โทร:</span> {{ patient.phone_number }}</div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <button :class="{ active: activeTab === 'visits' }" @click="activeTab = 'visits'">
          <Calendar size="18" /> บันทึกการตรวจครรภ์
        </button>
        <button :class="{ active: activeTab === 'medical' }" @click="activeTab = 'medical'">
          <FileText size="18" /> ประวัติสุขภาพ
        </button>
        <button :class="{ active: activeTab === 'vaccination' }" @click="activeTab = 'vaccination'">
          <Syringe size="18" /> วัคซีน
        </button>
        <button :class="{ active: activeTab === 'lab' }" @click="activeTab = 'lab'">
          <FlaskConical size="18" /> ผลแล็บ
        </button>
      </div>

      <!-- Visits Tab -->
      <div v-if="activeTab === 'visits'" class="card">
        <h2>บันทึกการตรวจครรภ์</h2>

        <div v-if="!pregnancyId" class="no-pregnancy">
          <p class="warning-text">⚠️ คนไข้ยังไม่มีข้อมูลการตั้งครรภ์</p>
          <div class="create-form">
            <div class="form-group">
              <label>วันประจำเดือนหมดครั้งสุดท้าย (LMP)</label>
              <input type="date" v-model="newPregnancyLMP" />
            </div>
            <button @click="createPregnancy" class="btn-create">สร้างข้อมูลการตั้งครรภ์</button>
          </div>
        </div>

        <form v-else @submit.prevent="saveVisit">
          <div class="form-grid">
            <div>
              <label>วันที่ตรวจ *</label>
              <input
                type="date"
                v-model="visitForm.VisitDate"
                @change="onVisitDateChange"
                required
              />
            </div>
            <div>
              <label>อายุครรภ์ (สัปดาห์) - คำนวณอัตโนมัติ</label>
              <input
                type="number"
                v-model="visitForm.GestationalAge"
                readonly
                style="background: #f3f4f6"
              />
            </div>
            <div>
              <label>น้ำหนัก (kg) *</label>
              <input type="number" step="0.1" v-model="visitForm.Weight" required />
            </div>
            <div>
              <label>ความดัน *</label>
              <input type="text" v-model="visitForm.BloodPressure" placeholder="120/80" required />
            </div>
            <div>
              <label>ส่วนสูงมดลูก (cm)</label>
              <input type="number" step="0.1" v-model="visitForm.HeightFundus" placeholder="30" />
            </div>
            <div>
              <label>เสียงหัวใจทารก</label>
              <select v-model="visitForm.FetalHeartSound">
                <option>ปกติ</option>
                <option>ผิดปกติ</option>
                <option>ไม่พบ</option>
              </select>
            </div>
            <div>
              <label>การเคลื่อนไหวทารก</label>
              <select v-model="visitForm.FetalMovement">
                <option>ปกติ</option>
                <option>ผิดปกติ</option>
                <option>ไม่พบ</option>
              </select>
            </div>
            <div>
              <label>โปรตีนในปัสสาวะ</label>
              <select v-model="visitForm.UrineProtein">
                <option>ไม่พบ</option>
                <option>พบ (+)</option>
                <option>พบ (++)</option>
                <option>พบ (+++)</option>
              </select>
            </div>
            <div>
              <label>น้ำตาลในปัสสาวะ</label>
              <select v-model="visitForm.UrineSugar">
                <option>ไม่พบ</option>
                <option>พบ (+)</option>
                <option>พบ (++)</option>
                <option>พบ (+++)</option>
              </select>
            </div>
            <div>
              <label>การบวม</label>
              <select v-model="visitForm.Swelling">
                <option>ไม่บวม</option>
                <option>บวมเล็กน้อย</option>
                <option>บวมปานกลาง</option>
                <option>บวมมาก</option>
              </select>
            </div>
            <div style="grid-column: 1 / -1">
              <label>การวินิจฉัยทางการแพทย์</label>
              <textarea
                v-model="visitForm.MedicalDiagnosis"
                rows="3"
                placeholder="บันทึกการวินิจฉัยและข้อแนะนำ..."
              ></textarea>
            </div>
          </div>
          <button type="submit" class="btn-save">
            <Save size="18" />
            บันทึก
          </button>
        </form>
      </div>

      <!-- Medical History Tab -->
      <div v-if="activeTab === 'medical'" class="card">
        <h2>ประวัติสุขภาพ</h2>
        <form @submit.prevent="saveMedicalHistory">
          <div class="form-grid">
            <div>
              <label>โรคประจำตัว</label>
              <input
                type="text"
                v-model="medicalHistoryForm.ChronicDiseases"
                placeholder="เบาหวาน, ความดันโลหิตสูง"
              />
            </div>
            <div>
              <label>โรคทางพันธุกรรม</label>
              <input
                type="text"
                v-model="medicalHistoryForm.GeneticDiseases"
                placeholder="ธาลัสซีเมีย"
              />
            </div>
            <div>
              <label>แพ้ยา</label>
              <input
                type="text"
                v-model="medicalHistoryForm.DrugAllergies"
                placeholder="Penicillin, Aspirin"
              />
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryHT" />
                ประวัติความดันในครอบครัว
              </label>
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryCongenital" />
                ประวัติพิการแต่กำเนิดในครอบครัว
              </label>
            </div>
            <div style="grid-column: 1 / -1">
              <label>ประวัติอื่นๆ</label>
              <textarea v-model="medicalHistoryForm.OtherFamilyHistory" rows="3"></textarea>
            </div>
          </div>
          <button type="submit" class="btn-save">
            <Save size="18" />
            บันทึก
          </button>
        </form>
      </div>

      <!-- Vaccination Tab -->
      <div v-if="activeTab === 'vaccination'" class="card">
        <h2>บันทึกวัคซีน</h2>
        <form @submit.prevent="saveVaccination">
          <div class="form-grid">
            <div>
              <label>
                <input type="checkbox" v-model="vaccinationForm.IsPreviouslyVaccinated" />
                เคยฉีดวัคซีนมาก่อน
              </label>
            </div>
            <div>
              <label>จำนวนครั้งที่เคยฉีด</label>
              <input type="number" v-model.number="vaccinationForm.PreviousDoses" />
            </div>
            <div>
              <label>วันที่ฉีดครั้งสุดท้าย (ก่อนตั้งครรภ์)</label>
              <input type="date" v-model="vaccinationForm.LastPreviousDateYear" />
            </div>
            <div>
              <label>Dose 1 (ครรภ์นี้)</label>
              <input type="date" v-model="vaccinationForm.Dose1DateDuringPreg" />
            </div>
            <div>
              <label>Dose 2 (ครรภ์นี้)</label>
              <input type="date" v-model="vaccinationForm.Dose2DateDuringPreg" />
            </div>
            <div style="grid-column: 1 / -1">
              <label>หมายเหตุ</label>
              <textarea v-model="vaccinationForm.Remarks" rows="3"></textarea>
            </div>
          </div>
          <button type="submit" class="btn-save">
            <Save size="18" />
            บันทึก
          </button>
        </form>
      </div>

      <!-- Lab Result Tab -->
      <div v-if="activeTab === 'lab'" class="card">
        <h2>บันทึกผลแล็บ</h2>
        <div v-if="!pregnancyId" class="no-pregnancy">
          <p class="warning-text">⚠️ คนไข้ยังไม่มีข้อมูลการตั้งครรภ์</p>
        </div>
        <form v-else @submit.prevent="saveLabResult">
          <div class="form-grid">
            <div>
              <label>วันที่ตรวจ *</label>
              <input type="date" v-model="labResultForm.TestDate" required />
            </div>
            <div>
              <label>Hct *</label>
              <input
                type="number"
                step="0.1"
                v-model="labResultForm.Hct"
                required
                placeholder="40.5"
              />
            </div>
            <div>
              <label>Hb *</label>
              <input
                type="number"
                step="0.1"
                v-model="labResultForm.Hb"
                required
                placeholder="13.5"
              />
            </div>
            <div>
              <label>Hb Typing</label>
              <input type="text" v-model="labResultForm.HbTyping" placeholder="A+" />
            </div>
            <div style="grid-column: 1 / -1">
              <label>หมายเหตุ</label>
              <textarea v-model="labResultForm.OtherRemarks" rows="3"></textarea>
            </div>
          </div>
          <button type="submit" class="btn-save">
            <Save size="18" />
            บันทึก
          </button>
        </form>
      </div>

      <!-- Visit History -->
      <div v-if="activeTab === 'visits'" class="card">
        <h2>ประวัติการตรวจ</h2>
        <div v-if="visits.length === 0">ยังไม่มีบันทึก</div>
        <div v-else v-for="visit in visits" :key="visit.ID" class="visit-item">
          <div class="visit-header">
            <Calendar size="16" />
            {{ formatDate(visit.VisitDate) }}
            <span class="badge">GA: {{ visit.GestationalAge }} สัปดาห์</span>
          </div>
          <div>น้ำหนัก: {{ visit.Weight }} kg | ความดัน: {{ visit.BloodPressure }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
.btn-back {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 0.375rem;
  cursor: pointer;
  margin-bottom: 1.5rem;
}
.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
.label {
  color: #666;
  font-size: 0.875rem;
}
.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}
.tabs button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
}
.tabs button.active {
  background: var(--color-primary);
  color: var(--color-text);
  border-color: var(--color-primary);
}
.no-pregnancy {
  text-align: center;
  padding: 2rem;
  background: #fef3c7;
  border-radius: 0.5rem;
}
.warning-text {
  color: #92400e;
  font-weight: 500;
  margin-bottom: 1rem;
}
.btn-create {
  padding: 0.75rem 1.5rem;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 600;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}
.form-grid label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
}
.form-grid input[type='checkbox'] {
  width: auto;
  margin-right: 0.5rem;
}
.form-grid input,
.form-grid textarea,
.form-grid select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 0.375rem;
}
.btn-save {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: var(--color-text);
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 600;
}
.visit-item {
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 0.375rem;
  margin-top: 0.5rem;
}
.visit-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
}
.badge {
  margin-left: auto;
  background: var(--color-primary);
  color: var(--color-text);
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
}
.create-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  max-width: 300px;
  margin: 0 auto;
}
.create-form .form-group {
  width: 100%;
  text-align: left;
}
.create-form label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: #4b5563;
}
.create-form input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 0.375rem;
}
</style>
