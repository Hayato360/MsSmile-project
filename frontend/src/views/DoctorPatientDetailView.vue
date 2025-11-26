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
const isEditingMedicalHistory = ref(false)

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
  MenstrualCycle: '',
  MenstrualDuration: '',
  MenstrualCondition: 'สม่ำเสมอ',
})

const previousPregnancies = ref([])
const previousPregnancyForm = ref({
  PregnancyNo: 1,
  DeliveryDate: '',
  GestationalAge: '',
  DeliveryMethod: 'Normal',
  BirthWeight: '',
  Sex: 'Male',
  DeliveryPlace: '',
  Complications: '',
  ChildStatus: 'Healthy',
})

const vaccinationForm = ref({
  VaccineTypeID: 1,
  IsPreviouslyVaccinated: false,
  PreviousDoses: 0,
  LastPreviousDateYear: null,
  Dose1DateDuringPreg: null,
  Dose2DateDuringPreg: null,
  Remarks: '',
  IsPreviouslyVaccinated: false, // Duplicate removed
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

    // Load previous pregnancies
    const prevPregRes = await api.get(`/doctor/patient/${patientId}/previous-pregnancies`)
    previousPregnancies.value = prevPregRes.data.data || []

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
    isEditingMedicalHistory.value = false
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

const createPreviousPregnancy = async () => {
  try {
    await api.post('/doctor/previous-pregnancy', {
      PregnantWomanID: parseInt(route.params.id),
      pregnancy_no: parseInt(previousPregnancyForm.value.PregnancyNo),
      delivery_date: new Date(previousPregnancyForm.value.DeliveryDate).toISOString(),
      gestational_age: parseInt(previousPregnancyForm.value.GestationalAge),
      delivery_method: previousPregnancyForm.value.DeliveryMethod,
      birth_weight: parseFloat(previousPregnancyForm.value.BirthWeight),
      sex: previousPregnancyForm.value.Sex,
      delivery_place: previousPregnancyForm.value.DeliveryPlace,
      complications: previousPregnancyForm.value.Complications,
      child_status: previousPregnancyForm.value.ChildStatus,
    })
    alert('บันทึกข้อมูลครรภ์ในอดีตสำเร็จ')

    // Refresh list
    const prevPregRes = await api.get(`/doctor/patient/${route.params.id}/previous-pregnancies`)
    previousPregnancies.value = prevPregRes.data.data || []

    // Reset form (optional, maybe keep for next entry but increment No)
    previousPregnancyForm.value.PregnancyNo++
    previousPregnancyForm.value.DeliveryDate = ''
    previousPregnancyForm.value.GestationalAge = ''
    previousPregnancyForm.value.BirthWeight = ''
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
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

    // Load previous pregnancies
    const prevPregRes = await api.get(`/doctor/patient/${patientId}/previous-pregnancies`)
    previousPregnancies.value = prevPregRes.data.data || []

    calculateGA()
  } catch (error) {
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
})

const goBack = () => router.push('/doctor/dashboard')
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  if (date.getFullYear() < 1900) return '-'
  return date.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
const calculateAge = (birthDate) => {
  if (!birthDate) return '-'
  const birth = new Date(birthDate)
  const today = new Date()
  let age = today.getFullYear() - birth.getFullYear()
  const m = today.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) {
    age--
  }
  return age
}
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
          <div><span class="label">อายุ:</span> {{ calculateAge(patient.birth_date) }} ปี</div>
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
        <button :class="{ active: activeTab === 'prev_preg' }" @click="activeTab = 'prev_preg'">
          <FileText size="18" /> ประวัติครรภ์อดีต
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

        <!-- View Mode -->
        <div v-if="!isEditingMedicalHistory">
          <div class="flex justify-between items-center mb-4">
            <h3>ประวัติการเจ็บป่วย</h3>
            <button @click="isEditingMedicalHistory = true" class="btn-edit">แก้ไขข้อมูล</button>
          </div>

          <div class="view-grid">
            <div class="view-item">
              <span class="label">โรคประจำตัว:</span>
              <span class="value">{{ medicalHistoryForm.ChronicDiseases || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">โรคหัวใจ:</span>
              <span class="value">{{ medicalHistoryForm.HeartDisease ? 'มี' : 'ไม่มี' }}</span>
            </div>
            <div class="view-item">
              <span class="label">โรคไทรอยด์:</span>
              <span class="value">{{ medicalHistoryForm.Thyroid ? 'มี' : 'ไม่มี' }}</span>
            </div>
            <div class="view-item">
              <span class="label">โรคอื่นๆ:</span>
              <span class="value">{{ medicalHistoryForm.OtherDiseases || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ประวัติผ่าตัด:</span>
              <span class="value">{{ medicalHistoryForm.SurgeryHistory || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ผ่าตัดอื่นๆ:</span>
              <span class="value">{{ medicalHistoryForm.OtherSurgery || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">โรคทางพันธุกรรม:</span>
              <span class="value">{{ medicalHistoryForm.GeneticDiseases || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">แพ้ยา:</span>
              <span class="value">{{ medicalHistoryForm.DrugAllergies || '-' }}</span>
            </div>
          </div>

          <h3 class="mt-6">ประวัติครอบครัว</h3>
          <div class="view-grid">
            <div class="view-item">
              <span class="label">ความดันโลหิตสูง:</span>
              <span class="value">{{ medicalHistoryForm.FamilyHistoryHT ? 'มี' : 'ไม่มี' }}</span>
            </div>
            <div class="view-item">
              <span class="label">เบาหวาน:</span>
              <span class="value">{{
                medicalHistoryForm.FamilyHistoryDiabetes ? 'มี' : 'ไม่มี'
              }}</span>
            </div>
            <div class="view-item">
              <span class="label">โลหิตจาง (ธาลัสซีเมีย):</span>
              <span class="value">{{
                medicalHistoryForm.FamilyHistoryThalassemia ? 'มี' : 'ไม่มี'
              }}</span>
            </div>
            <div class="view-item">
              <span class="label">พิการแต่กำเนิด:</span>
              <span class="value">{{
                medicalHistoryForm.FamilyHistoryCongenital ? 'มี' : 'ไม่มี'
              }}</span>
            </div>
            <div class="view-item col-span-2">
              <span class="label">ประวัติอื่นๆ:</span>
              <span class="value">{{ medicalHistoryForm.OtherFamilyHistory || '-' }}</span>
            </div>
          </div>

          <h3 class="mt-6">ประวัติการคุมกำเนิด & ประจำเดือน</h3>
          <div class="view-grid">
            <div class="view-item">
              <span class="label">คุมกำเนิดก่อนตั้งครรภ์ (วิธี):</span>
              <span class="value">{{ medicalHistoryForm.ContraceptionBeforeMethod || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ระยะเวลา:</span>
              <span class="value">{{ medicalHistoryForm.ContraceptionBeforeDuration || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">คุมกำเนิดครั้งหลังสุด (วิธี):</span>
              <span class="value">{{ medicalHistoryForm.ContraceptionLastMethod || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ระยะเวลา:</span>
              <span class="value">{{ medicalHistoryForm.ContraceptionLastDuration || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ประจำเดือนมาทุก (วัน):</span>
              <span class="value">{{ medicalHistoryForm.MenstrualCycle || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">นานครั้งละ (วัน):</span>
              <span class="value">{{ medicalHistoryForm.MenstrualDuration || '-' }}</span>
            </div>
            <div class="view-item">
              <span class="label">ความสม่ำเสมอ:</span>
              <span class="value">{{ medicalHistoryForm.MenstrualCondition || '-' }}</span>
            </div>
          </div>
        </div>

        <!-- Edit Mode -->
        <form v-else @submit.prevent="saveMedicalHistory">
          <h3>ประวัติการเจ็บป่วย</h3>
          <div class="form-grid">
            <div>
              <label>โรคประจำตัว (ระบุ)</label>
              <input type="text" v-model="medicalHistoryForm.ChronicDiseases" />
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.HeartDisease" />
                โรคหัวใจ
              </label>
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.Thyroid" />
                โรคไทรอยด์
              </label>
            </div>
            <div>
              <label>โรคอื่นๆ</label>
              <input type="text" v-model="medicalHistoryForm.OtherDiseases" />
            </div>
            <div>
              <label>ประวัติผ่าตัด</label>
              <input type="text" v-model="medicalHistoryForm.SurgeryHistory" />
            </div>
            <div>
              <label>ผ่าตัดอื่นๆ</label>
              <input type="text" v-model="medicalHistoryForm.OtherSurgery" />
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
          </div>

          <h3>ประวัติครอบครัว</h3>
          <div class="form-grid">
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryHT" />
                ความดันโลหิตสูง
              </label>
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryDiabetes" />
                เบาหวาน
              </label>
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryThalassemia" />
                โลหิตจาง (ธาลัสซีเมีย)
              </label>
            </div>
            <div>
              <label>
                <input type="checkbox" v-model="medicalHistoryForm.FamilyHistoryCongenital" />
                พิการแต่กำเนิด
              </label>
            </div>
            <div style="grid-column: 1 / -1">
              <label>ประวัติอื่นๆ</label>
              <textarea v-model="medicalHistoryForm.OtherFamilyHistory" rows="2"></textarea>
            </div>
          </div>

          <h3>ประวัติการคุมกำเนิด & ประจำเดือน</h3>
          <div class="form-grid">
            <div>
              <label>คุมกำเนิดก่อนตั้งครรภ์ (วิธี)</label>
              <input type="text" v-model="medicalHistoryForm.ContraceptionBeforeMethod" />
            </div>
            <div>
              <label>ระยะเวลา</label>
              <input type="text" v-model="medicalHistoryForm.ContraceptionBeforeDuration" />
            </div>
            <div>
              <label>คุมกำเนิดครั้งหลังสุด (วิธี)</label>
              <input type="text" v-model="medicalHistoryForm.ContraceptionLastMethod" />
            </div>
            <div>
              <label>ระยะเวลา</label>
              <input type="text" v-model="medicalHistoryForm.ContraceptionLastDuration" />
            </div>
            <div>
              <label>ประจำเดือนมาทุก (วัน)</label>
              <input type="number" v-model.number="medicalHistoryForm.MenstrualCycle" />
            </div>
            <div>
              <label>นานครั้งละ (วัน)</label>
              <input type="number" v-model.number="medicalHistoryForm.MenstrualDuration" />
            </div>
            <div>
              <label>ความสม่ำเสมอ</label>
              <select v-model="medicalHistoryForm.MenstrualCondition">
                <option>สม่ำเสมอ</option>
                <option>ไม่สม่ำเสมอ</option>
              </select>
            </div>
          </div>

          <div class="flex gap-2 mt-4">
            <button type="submit" class="btn-save">
              <Save size="18" />
              บันทึกประวัติสุขภาพ
            </button>
            <button type="button" @click="isEditingMedicalHistory = false" class="btn-cancel">
              ยกเลิก
            </button>
          </div>
        </form>
      </div>

      <!-- Previous Pregnancies Tab -->
      <div v-if="activeTab === 'prev_preg'" class="card">
        <h2>ประวัติการตั้งครรภ์ในอดีต</h2>

        <!-- List -->
        <div class="table-responsive mb-4">
          <table class="data-table">
            <thead>
              <tr>
                <th>ครรภ์ที่</th>
                <th>วันคลอด</th>
                <th>GA (wk)</th>
                <th>วิธีคลอด</th>
                <th>นน. (kg)</th>
                <th>เพศ</th>
                <th>สถานที่</th>
                <th>ภาวะแทรกซ้อน</th>
                <th>สถานะเด็ก</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="preg in previousPregnancies" :key="preg.ID">
                <td>{{ preg.pregnancy_no }}</td>
                <td>{{ formatDate(preg.delivery_date) }}</td>
                <td>{{ preg.gestational_age }}</td>
                <td>{{ preg.delivery_method }}</td>
                <td>{{ preg.birth_weight }}</td>
                <td>{{ preg.sex }}</td>
                <td>{{ preg.delivery_place }}</td>
                <td>{{ preg.complications }}</td>
                <td>{{ preg.child_status }}</td>
              </tr>
              <tr v-if="previousPregnancies.length === 0">
                <td colspan="9" class="text-center">ไม่มีข้อมูล</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Add Form -->
        <h3>เพิ่มข้อมูลครรภ์ในอดีต</h3>
        <form @submit.prevent="createPreviousPregnancy">
          <div class="form-grid">
            <div>
              <label>ครรภ์ที่</label>
              <input type="number" v-model.number="previousPregnancyForm.PregnancyNo" required />
            </div>
            <div>
              <label>วันที่คลอด/แท้ง</label>
              <input type="date" v-model="previousPregnancyForm.DeliveryDate" required />
            </div>
            <div>
              <label>อายุครรภ์ (สัปดาห์)</label>
              <input type="number" v-model.number="previousPregnancyForm.GestationalAge" required />
            </div>
            <div>
              <label>วิธีคลอด/แท้ง</label>
              <select v-model="previousPregnancyForm.DeliveryMethod">
                <option>Normal</option>
                <option>C-Section</option>
                <option>Vacuum</option>
                <option>Abortion</option>
              </select>
            </div>
            <div>
              <label>น้ำหนักทารก (kg)</label>
              <input
                type="number"
                step="0.1"
                v-model.number="previousPregnancyForm.BirthWeight"
                required
              />
            </div>
            <div>
              <label>เพศ</label>
              <select v-model="previousPregnancyForm.Sex">
                <option>Male</option>
                <option>Female</option>
              </select>
            </div>
            <div>
              <label>สถานที่คลอด</label>
              <input type="text" v-model="previousPregnancyForm.DeliveryPlace" />
            </div>
            <div>
              <label>ภาวะแทรกซ้อน</label>
              <input type="text" v-model="previousPregnancyForm.Complications" />
            </div>
            <div>
              <label>สถานะเด็ก</label>
              <select v-model="previousPregnancyForm.ChildStatus">
                <option>Healthy</option>
                <option>Deceased</option>
                <option>Sick</option>
              </select>
            </div>
          </div>
          <button type="submit" class="btn-save">
            <Save size="18" />
            เพิ่มข้อมูล
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
  color: white;
  border-color: var(--color-primary);
}
.no-pregnancy {
  text-align: center;
  padding: 2rem;
  background: #fff7ed;
  border-radius: 0.5rem;
}
.warning-text {
  color: #c2410c;
  font-weight: 600;
  margin-bottom: 1rem;
}
.btn-create {
  padding: 0.75rem 1.5rem;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
}
/* View Mode Styles */
.view-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
  background: #f9fafb;
  padding: 1.5rem;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
}

.view-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.view-item .label {
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
}

.view-item .value {
  color: #1e293b;
  font-size: 1rem;
  font-weight: 600;
}

.btn-edit {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: white;
  border: 1px solid #d1d5db;
  color: #374151;
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-edit:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.btn-cancel {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: white;
  border: 1px solid #d1d5db;
  color: #374151;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
}

.btn-cancel:hover {
  background: #f3f4f6;
}

/* Form Styles Update */
.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.form-grid label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  font-size: 0.875rem;
  color: #374151;
}

.form-grid input[type='text'],
.form-grid input[type='number'],
.form-grid input[type='date'],
.form-grid select,
.form-grid textarea {
  width: 100%;
  padding: 0.625rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

.form-grid input:focus,
.form-grid select:focus,
.form-grid textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(163, 230, 53, 0.1);
}

.btn-save,
.btn-create {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn-save:hover,
.btn-create:hover {
  background: #65a30d; /* Darker green */
}

/* Table Styles */
.table-responsive {
  overflow-x: auto;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.data-table th {
  background-color: #f8fafc;
  color: #475569;
  font-weight: 600;
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
  white-space: nowrap;
}

.data-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #e5e7eb;
  color: #1e293b;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr:hover {
  background-color: #f1f5f9;
}

.text-center {
  text-align: center;
}
</style>
