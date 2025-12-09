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

const vaccineTypes = ref([])
const vaccinations = ref([])

const vaccinationForm = ref({
  VaccineTypeID: '',
  IsPreviouslyVaccinated: false,
  PreviousDoses: 0,
  LastPreviousDateYear: null,
  Dose1DateDuringPreg: null,
  Dose2DateDuringPreg: null,
  Dose3DateDuringPreg: null,
  IsHistoryUnknown: false,
  ReasonForNotVaccinating: '',
  Remarks: '',
})

const vaccinationHistoryStatus = computed({
  get: () => {
    if (vaccinationForm.value.IsHistoryUnknown) return 'unknown'
    if (vaccinationForm.value.IsPreviouslyVaccinated) return 'previously'
    return 'never'
  },
  set: (val) => {
    if (val === 'unknown') {
      vaccinationForm.value.IsHistoryUnknown = true
      vaccinationForm.value.IsPreviouslyVaccinated = false
    } else if (val === 'previously') {
      vaccinationForm.value.IsHistoryUnknown = false
      vaccinationForm.value.IsPreviouslyVaccinated = true
    } else {
      vaccinationForm.value.IsHistoryUnknown = false
      vaccinationForm.value.IsPreviouslyVaccinated = false
    }
  },
})

const handleHistoryUnknownChange = (e) => {
  if (e.target.checked) {
    vaccinationForm.value.IsPreviouslyVaccinated = false
  }
}

const labResultForm = ref({
  TestDate: new Date().toISOString().split('T')[0],
  Hct: '',
  Hb: '',
  HbTyping: '',
  OtherRemarks: '',
})

const pregnancyId = computed(() => {
  if (!patient.value?.Pregnancies) return null
  const active = patient.value.Pregnancies.find((p) => p.status === 'Active')
  return active ? active.ID : null
})

const pregnancyData = computed(() => {
  if (!patient.value?.Pregnancies) return null
  return patient.value.Pregnancies.find((p) => p.status === 'Active')
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

    // Load vaccine types
    const vacTypesRes = await api.get('/vaccine-types')
    vaccineTypes.value = vacTypesRes.data.data || []
    if (vaccineTypes.value.length > 0) {
      vaccinationForm.value.VaccineTypeID = vaccineTypes.value[0].ID
    }

    // Load existing vaccinations
    const vacRes = await api.get(`/vaccinations/pregnant-woman/${patientId}`)
    vaccinations.value = vacRes.data || []

    calculateGA()
  } catch (error) {
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
})

const newPregnancyLMP = ref(new Date().toISOString().split('T')[0])
const newPregnancyEDC = ref('')

// Watch LMP to auto-calculate EDC
import { watch } from 'vue'
watch(newPregnancyLMP, (newValue) => {
  if (newValue) {
    const lmp = new Date(newValue)
    const edc = new Date(lmp)
    edc.setDate(edc.getDate() + 280) // +280 days
    newPregnancyEDC.value = edc.toISOString().split('T')[0]
  }
})

const createPregnancy = async () => {
  try {
    await api.post('/doctor/pregnancy', {
      PregnantWomanID: parseInt(route.params.id),
      PregnancyNo: (patient.value.Pregnancies?.length || 0) + 1,
      LMP: new Date(newPregnancyLMP.value).toISOString(),
      EDC: new Date(newPregnancyEDC.value).toISOString(),
      PrePregnancyWeight: 50, // Should be input, but keeping simple for now
      Height: 160,
      PrePregnancyBMI: 19.5,
    })

    alert('สร้างข้อมูลการตั้งครรภ์สำเร็จ')
    const patientRes = await api.get(`/doctor/patients/${route.params.id}`)
    patient.value = patientRes.data
  } catch (error) {
    console.error('Error:', error)
    alert(error.response?.data?.error || 'เกิดข้อผิดพลาด')
  }
}

const showEndPregnancyModal = ref(false)
const endPregnancyForm = ref({
  DeliveryDate: new Date().toISOString().split('T')[0],
  DeliveryMethod: 'Normal',
  BirthWeight: '',
  Sex: 'Male',
  DeliveryPlace: '',
  Complications: 'None',
  ChildStatus: 'Healthy',
  GestationalAge: '',
})

const endPregnancy = async () => {
  if (!pregnancyId.value) return

  try {
    await api.post(`/doctor/pregnancy/${pregnancyId.value}/end`, {
      delivery_date: new Date(endPregnancyForm.value.DeliveryDate).toISOString(),
      delivery_method: endPregnancyForm.value.DeliveryMethod,
      birth_weight: parseFloat(endPregnancyForm.value.BirthWeight),
      sex: endPregnancyForm.value.Sex,
      delivery_place: endPregnancyForm.value.DeliveryPlace,
      complications: endPregnancyForm.value.Complications,
      child_status: endPregnancyForm.value.ChildStatus,
      gestational_age: parseInt(endPregnancyForm.value.GestationalAge),
    })

    alert('จบการตั้งครรภ์เรียบร้อยแล้ว')
    showEndPregnancyModal.value = false

    // Refresh data
    const patientRes = await api.get(`/doctor/patients/${route.params.id}`)
    patient.value = patientRes.data

    const prevPregRes = await api.get(`/doctor/patient/${route.params.id}/previous-pregnancies`)
    previousPregnancies.value = prevPregRes.data.data || []
  } catch (error) {
    console.error('Error:', error)
    alert(error.response?.data?.error || 'เกิดข้อผิดพลาด')
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
      Dose3DateDuringPreg: vaccinationForm.value.Dose3DateDuringPreg
        ? new Date(vaccinationForm.value.Dose3DateDuringPreg).toISOString()
        : null,
      IsHistoryUnknown: vaccinationForm.value.IsHistoryUnknown,
      ReasonForNotVaccinating: vaccinationForm.value.ReasonForNotVaccinating,
    })
    alert('บันทึกข้อมูลวัคซีนสำเร็จ')

    // Refresh list
    const vacRes = await api.get(`/vaccinations/pregnant-woman/${route.params.id}`)
    vaccinations.value = vacRes.data || []

    // Reset form (optional, but keeping type selected is usually good)
    vaccinationForm.value.IsPreviouslyVaccinated = false
    vaccinationForm.value.PreviousDoses = 0
    vaccinationForm.value.LastPreviousDateYear = null
    vaccinationForm.value.Dose1DateDuringPreg = null
    vaccinationForm.value.Dose2DateDuringPreg = null
    vaccinationForm.value.Dose3DateDuringPreg = null
    vaccinationForm.value.IsHistoryUnknown = false
    vaccinationForm.value.ReasonForNotVaccinating = ''
    vaccinationForm.value.Remarks = ''
  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาด')
  }
}

// Lab Results Logic
const labResults = ref([])
const labResultFile = ref(null)

const handleFileUpload = (event) => {
  labResultFile.value = event.target.files[0]
}

const fetchLabResults = async () => {
  if (!pregnancyId.value) return
  try {
    const res = await api.get(`/doctor/pregnancy/${pregnancyId.value}/lab-results`)
    labResults.value = res.data.data || []
  } catch (error) {
    console.error('Error fetching lab results:', error)
  }
}

const saveLabResult = async () => {
  if (!pregnancyId.value) {
    alert('คนไข้ยังไม่มีข้อมูลการตั้งครรภ์')
    return
  }

  try {
    const formData = new FormData()
    formData.append('PregnancyID', pregnancyId.value)
    formData.append('TestDate', labResultForm.value.TestDate)
    formData.append('Hct', labResultForm.value.Hct)
    formData.append('Hb', labResultForm.value.Hb)
    formData.append('HbTyping', labResultForm.value.HbTyping)
    formData.append('OtherRemarks', labResultForm.value.OtherRemarks)
    
    if (labResultFile.value) {
      formData.append('File', labResultFile.value)
    }

    await api.post('/doctor/lab-result', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    
    alert('บันทึกผลแล็บสำเร็จ')
    
    // Reset form
    labResultForm.value = {
      TestDate: new Date().toISOString().split('T')[0],
      Hct: '',
      Hb: '',
      HbTyping: '',
      OtherRemarks: '',
    }
    labResultFile.value = null
    // Reset file input
    const fileInput = document.querySelector('input[type="file"]')
    if (fileInput) fileInput.value = ''

    await fetchLabResults()

  } catch (error) {
    console.error('Error:', error)
    alert('เกิดข้อผิดพลาดในการบันทึก: ' + (error.response?.data?.error || error.message))
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

    await fetchLabResults()
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
            <div class="form-group">
              <label>วันกำหนดคลอด (EDC)</label>
              <input type="date" v-model="newPregnancyEDC" />
            </div>
            <button @click="createPregnancy" class="btn-create">สร้างข้อมูลการตั้งครรภ์</button>
          </div>
        </div>

        <div v-else>
          <div class="flex justify-between items-center mb-4">
            <h3>บันทึกการตรวจ</h3>
            <button @click="showEndPregnancyModal = true" class="btn-end-pregnancy">
              จบการตั้งครรภ์
            </button>
          </div>
          <form @submit.prevent="saveVisit">
            <div class="visit-form-container">
              <div class="visit-form-columns">
                <!-- Left Column -->
                <div class="visit-form-col">
                  <div class="form-group">
                    <label>วันที่ตรวจ *</label>
                    <input
                      type="date"
                      v-model="visitForm.VisitDate"
                      @change="onVisitDateChange"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <label>อายุครรภ์ (สัปดาห์) - คำนวณอัตโนมัติ</label>
                    <input
                      type="number"
                      v-model="visitForm.GestationalAge"
                      readonly
                      style="background: #f3f4f6"
                    />
                  </div>
                  <div class="form-group">
                    <label>น้ำหนัก (kg) *</label>
                    <input type="number" step="0.1" v-model="visitForm.Weight" required />
                  </div>
                  <div class="form-group">
                    <label>ความดัน *</label>
                    <input
                      type="text"
                      v-model="visitForm.BloodPressure"
                      placeholder="120/80"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <label>ส่วนสูงมดลูก (cm)</label>
                    <input
                      type="number"
                      step="0.1"
                      v-model="visitForm.HeightFundus"
                      placeholder="30"
                    />
                  </div>
                </div>

                <!-- Right Column -->
                <div class="visit-form-col">
                  <div class="form-group">
                    <label>เสียงหัวใจทารก</label>
                    <select v-model="visitForm.FetalHeartSound">
                      <option>ปกติ</option>
                      <option>ผิดปกติ</option>
                      <option>ไม่พบ</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label>การเคลื่อนไหวทารก</label>
                    <select v-model="visitForm.FetalMovement">
                      <option>ปกติ</option>
                      <option>ผิดปกติ</option>
                      <option>ไม่พบ</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label>โปรตีนในปัสสาวะ</label>
                    <select v-model="visitForm.UrineProtein">
                      <option>ไม่พบ</option>
                      <option>พบ (+)</option>
                      <option>พบ (++)</option>
                      <option>พบ (+++)</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label>น้ำตาลในปัสสาวะ</label>
                    <select v-model="visitForm.UrineSugar">
                      <option>ไม่พบ</option>
                      <option>พบ (+)</option>
                      <option>พบ (++)</option>
                      <option>พบ (+++)</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label>การบวม</label>
                    <select v-model="visitForm.Swelling">
                      <option>ไม่บวม</option>
                      <option>บวมเล็กน้อย</option>
                      <option>บวมปานกลาง</option>
                      <option>บวมมาก</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Diagnosis Full Width -->
              <div class="form-group full-width mt-4">
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
      </div>

      <!-- Medical History Tab -->
      <div v-if="activeTab === 'medical'" class="card">
        <div class="card-header">
          <FileHeart size="24" />
          <h3>ประวัติสุขภาพ</h3>
        </div>

        <!-- View Mode -->
        <div v-if="!isEditingMedicalHistory">
          <div class="info-section">
            <h4>ประวัติการเจ็บป่วย</h4>
            <div class="info-grid">
              <div class="info-item full-width">
                <span class="label">โรคประจำตัว</span>
                <div class="tags">
                  <span v-if="medicalHistoryForm.ChronicDiseases" class="tag">{{
                    medicalHistoryForm.ChronicDiseases
                  }}</span>
                  <span v-if="medicalHistoryForm.OtherDiseases" class="tag">{{
                    medicalHistoryForm.OtherDiseases
                  }}</span>
                  <span
                    v-if="!medicalHistoryForm.ChronicDiseases && !medicalHistoryForm.OtherDiseases"
                    class="text-muted"
                    >-</span
                  >
                </div>
              </div>
              <div class="info-item full-width">
                <span class="label">ประวัติผ่าตัด</span>
                <span class="value">{{ medicalHistoryForm.SurgeryHistory || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="label">ประวัติแพ้ยา/อาหาร</span>
                <span class="value">{{ medicalHistoryForm.DrugAllergies || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="label">โรคทางพันธุกรรม</span>
                <span class="value">{{ medicalHistoryForm.GeneticDiseases || '-' }}</span>
              </div>
            </div>
          </div>

          <div class="divider"></div>

          <div class="info-section">
            <h4>ประวัติครอบครัว</h4>
            <div class="tags">
              <span v-if="medicalHistoryForm.OtherFamilyHistory" class="tag">{{
                medicalHistoryForm.OtherFamilyHistory
              }}</span>
              <span v-if="!medicalHistoryForm.OtherFamilyHistory" class="text-muted"
                >ไม่มีประวัติระบุ</span
              >
            </div>
          </div>

          <div class="divider"></div>

          <div class="info-section">
            <h4>ประวัติประจำเดือน & การคุมกำเนิด</h4>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">รอบประจำเดือน</span>
                <span class="value">{{
                  medicalHistoryForm.MenstrualCycle
                    ? medicalHistoryForm.MenstrualCycle + ' วัน'
                    : '-'
                }}</span>
              </div>
              <div class="info-item">
                <span class="label">จำนวนวันที่มี</span>
                <span class="value">{{
                  medicalHistoryForm.MenstrualDuration
                    ? medicalHistoryForm.MenstrualDuration + ' วัน'
                    : '-'
                }}</span>
              </div>
              <div class="info-item">
                <span class="label">ลักษณะประจำเดือน</span>
                <span class="value">{{ medicalHistoryForm.MenstrualCondition || '-' }}</span>
              </div>
              <div class="info-item full-width">
                <span class="label">การคุมกำเนิดก่อนตั้งครรภ์</span>
                <span class="value">
                  {{ medicalHistoryForm.ContraceptionBeforeMethod || '-' }}
                  <span v-if="medicalHistoryForm.ContraceptionBeforeDuration"
                    >({{ medicalHistoryForm.ContraceptionBeforeDuration }})</span
                  >
                </span>
              </div>
            </div>
          </div>

          <button @click="isEditingMedicalHistory = true" class="btn-edit mt-4">
            <Edit size="18" />
            แก้ไขข้อมูล
          </button>
        </div>

        <!-- Edit Mode -->
        <form v-else @submit.prevent="saveMedicalHistory" class="form-content">
          <!-- Illness History -->
          <div class="form-section">
            <h4 class="section-title">ประวัติการเจ็บป่วย</h4>
            <div class="form-grid">
              <div class="form-group full-width">
                <label>โรคประจำตัว</label>
                <textarea
                  v-model="medicalHistoryForm.ChronicDiseases"
                  rows="3"
                  placeholder="ระบุโรคประจำตัว (ถ้ามีหลายโรคให้ระบุทีละบรรทัด หรือคั่นด้วยจุลภาค)"
                ></textarea>
              </div>

              <div class="form-group">
                <label>ประวัติการผ่าตัด</label>
                <input
                  type="text"
                  v-model="medicalHistoryForm.SurgeryHistory"
                  placeholder="ระบุประวัติการผ่าตัด (ถ้ามี)"
                />
              </div>

              <div class="form-group">
                <label>ประวัติแพ้ยา/อาหาร</label>
                <input
                  type="text"
                  v-model="medicalHistoryForm.DrugAllergies"
                  placeholder="ระบุยาหรืออาหารที่แพ้ (ถ้ามี)"
                />
              </div>

              <div class="form-group">
                <label>โรคทางพันธุกรรม</label>
                <input
                  type="text"
                  v-model="medicalHistoryForm.GeneticDiseases"
                  placeholder="ระบุโรคทางพันธุกรรม (ถ้ามี)"
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
              <textarea
                v-model="medicalHistoryForm.OtherFamilyHistory"
                rows="3"
                placeholder="ระบุประวัติครอบครัว (ถ้ามีหลายโรคให้ระบุทีละบรรทัด หรือคั่นด้วยจุลภาค)"
              ></textarea>
            </div>
          </div>

          <div class="divider"></div>

          <!-- Menstruation & Contraception -->
          <div class="form-section">
            <h4 class="section-title">ประวัติประจำเดือน & การคุมกำเนิด</h4>
            <div class="form-grid">
              <div class="form-group">
                <label>รอบประจำเดือน (วัน)</label>
                <input
                  type="number"
                  v-model.number="medicalHistoryForm.MenstrualCycle"
                  placeholder="28"
                />
              </div>
              <div class="form-group">
                <label>จำนวนวันที่มีประจำเดือน</label>
                <input
                  type="number"
                  v-model.number="medicalHistoryForm.MenstrualDuration"
                  placeholder="5"
                />
              </div>
              <div class="form-group full-width">
                <label>ลักษณะประจำเดือน</label>
                <select v-model="medicalHistoryForm.MenstrualCondition">
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
                    v-model="medicalHistoryForm.ContraceptionBeforeMethod"
                    placeholder="วิธีคุมกำเนิด"
                  />
                  <input
                    type="text"
                    v-model="medicalHistoryForm.ContraceptionBeforeDuration"
                    placeholder="ระยะเวลา"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" @click="isEditingMedicalHistory = false" class="btn-cancel">
              ยกเลิก
            </button>
            <button type="submit" class="btn-save">
              <Save size="18" />
              บันทึก
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

        <!-- Vaccination List -->
        <div class="table-responsive mb-4">
          <table class="data-table">
            <thead>
              <tr>
                <th>ชนิดวัคซีน</th>
                <th>ประวัติ</th>
                <th>เข็มล่าสุด (ปี)</th>
                <th>Dose 1</th>
                <th>Dose 2</th>
                <th>Dose 3</th>
                <th>หมายเหตุ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="vac in vaccinations" :key="vac.ID">
                <td>{{ vac.VaccineType?.Name }}</td>
                <td>
                  <span v-if="vac.IsHistoryUnknown">ไม่ทราบประวัติ</span>
                  <span v-else-if="vac.IsPreviouslyVaccinated"
                    >เคยฉีด ({{ vac.PreviousDoses }} เข็ม)</span
                  >
                  <span v-else>ไม่เคยฉีด</span>
                </td>
                <td>
                  {{
                    vac.LastPreviousDateYear
                      ? formatDate(vac.LastPreviousDateYear).split(' ').pop()
                      : '-'
                  }}
                </td>
                <td>{{ formatDate(vac.Dose1DateDuringPreg) }}</td>
                <td>{{ formatDate(vac.Dose2DateDuringPreg) }}</td>
                <td>{{ formatDate(vac.Dose3DateDuringPreg) }}</td>
                <td>
                  <span v-if="vac.ReasonForNotVaccinating" class="text-danger">{{
                    vac.ReasonForNotVaccinating
                  }}</span>
                  <span v-else>{{ vac.Remarks || '-' }}</span>
                </td>
              </tr>
              <tr v-if="vaccinations.length === 0">
                <td colspan="7" class="text-center">ไม่มีข้อมูลการฉีดวัคซีน</td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3>เพิ่มข้อมูลวัคซีน</h3>
        <form @submit.prevent="saveVaccination">
          <div class="form-grid">
            <div style="grid-column: 1 / -1">
              <label>ชนิดวัคซีน</label>
              <select v-model="vaccinationForm.VaccineTypeID" required>
                <option v-for="vt in vaccineTypes" :key="vt.ID" :value="vt.ID">
                  {{ vt.Name }}
                </option>
              </select>
            </div>
            <div class="input-with-label-offset" style="grid-column: 1 / -1">
              <label>ประวัติการได้รับวัคซีน</label>
              <select v-model="vaccinationHistoryStatus">
                <option value="never">ไม่เคยฉีด</option>
                <option value="previously">เคยฉีดวัคซีนมาก่อน</option>
                <option value="unknown">ไม่ทราบประวัติ</option>
              </select>
            </div>
            <div>
              <label>จำนวนครั้งที่เคยฉีด</label>
              <input
                type="number"
                v-model.number="vaccinationForm.PreviousDoses"
                :disabled="!vaccinationForm.IsPreviouslyVaccinated"
              />
            </div>
            <div>
              <label>วันที่ฉีดครั้งสุดท้าย (ปี)</label>
              <input
                type="date"
                v-model="vaccinationForm.LastPreviousDateYear"
                :disabled="!vaccinationForm.IsPreviouslyVaccinated"
              />
            </div>
            <div>
              <label>Dose 1 (ครรภ์นี้)</label>
              <input type="date" v-model="vaccinationForm.Dose1DateDuringPreg" />
            </div>
            <div>
              <label>Dose 2 (ครรภ์นี้)</label>
              <input type="date" v-model="vaccinationForm.Dose2DateDuringPreg" />
            </div>
            <div>
              <label>Dose 3 (ครรภ์นี้)</label>
              <input type="date" v-model="vaccinationForm.Dose3DateDuringPreg" />
            </div>
            <div style="grid-column: 1 / -1">
              <label>เหตุผลที่ไม่ได้ฉีด (ถ้ามี)</label>
              <input
                type="text"
                v-model="vaccinationForm.ReasonForNotVaccinating"
                placeholder="เช่น แพ้วัคซีน, ปฏิเสธการรับวัคซีน"
              />
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
              <label>อัปโหลดเอกสาร (PDF)</label>
              <input type="file" ref="labResultFile" accept="application/pdf" @change="handleFileUpload" />
              <p class="text-sm text-gray-500 mt-1">รองรับไฟล์ PDF เท่านั้น</p>
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

        <!-- Lab Results History Table -->
        <h3 v-if="labResults && labResults.length > 0" class="mt-8 mb-4 border-t pt-4">ประวัติผลแล็บ</h3>
        <div class="table-responsive mb-4" v-if="labResults && labResults.length > 0">
          <table class="data-table">
            <thead>
              <tr>
                <th>วันที่ตรวจ</th>
                <th>Hct</th>
                <th>Hb</th>
                <th>Hb Typing</th>
                <th>ไฟล์แนบ</th>
                <th>หมายเหตุ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="result in labResults" :key="result.ID">
                <td>{{ formatDate(result.TestDate) }}</td>
                <td>{{ result.Hct }}</td>
                <td>{{ result.Hb }}</td>
                <td>{{ result.HbTyping }}</td>
                <td>
                  <a v-if="result.FilePath" 
                     :href="`http://localhost:8081/${result.FilePath}`" 
                     target="_blank" 
                     class="btn-view-pdf">
                    <FileText size="16" /> ดูเอกสาร
                  </a>
                  <span v-else>-</span>
                </td>
                <td>{{ result.OtherRemarks || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Visit History -->
      <div v-if="activeTab === 'visits'" class="card">
        <h2>ประวัติการตรวจ</h2>

        <div class="table-responsive">
          <table class="data-table">
            <thead>
              <tr>
                <th>วันที่</th>
                <th>อายุครรภ์</th>
                <th>น้ำหนัก (กก.)</th>
                <th>ความดัน</th>
                <th>เสียงหัวใจทารก</th>
                <th>การเคลื่อนไหว</th>
                <th>หมายเหตุ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="visit in visits" :key="visit.ID">
                <td>{{ formatDate(visit.VisitDate) }}</td>
                <td>{{ visit.GestationalAge }} สัปดาห์</td>
                <td>{{ visit.Weight }}</td>
                <td>{{ visit.BloodPressure }}</td>
                <td>{{ visit.FetalHeartSound }}</td>
                <td>{{ visit.FetalMovement }}</td>
                <td>{{ visit.MedicalDiagnosis || '-' }}</td>
              </tr>
              <tr v-if="visits.length === 0">
                <td colspan="7" class="text-center">ยังไม่มีบันทึก</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    <!-- End Pregnancy Modal -->
    <div v-if="showEndPregnancyModal" class="modal-overlay">
      <div class="modal-content">
        <h2>จบการตั้งครรภ์</h2>
        <form @submit.prevent="endPregnancy">
          <div class="form-grid">
            <div>
              <label>วันที่คลอด *</label>
              <input type="date" v-model="endPregnancyForm.DeliveryDate" required />
            </div>
            <div>
              <label>อายุครรภ์ (สัปดาห์) *</label>
              <input type="number" v-model="endPregnancyForm.GestationalAge" required />
            </div>
            <div>
              <label>วิธีการคลอด *</label>
              <select v-model="endPregnancyForm.DeliveryMethod">
                <option value="Normal">คลอดธรรมชาติ</option>
                <option value="C-Section">ผ่าตัดคลอด</option>
                <option value="Vacuum">เครื่องดูดสุญญากาศ</option>
              </select>
            </div>
            <div>
              <label>น้ำหนักแรกเกิด (กรัม) *</label>
              <input type="number" step="0.01" v-model="endPregnancyForm.BirthWeight" required />
            </div>
            <div>
              <label>เพศ *</label>
              <select v-model="endPregnancyForm.Sex">
                <option value="Male">ชาย</option>
                <option value="Female">หญิง</option>
              </select>
            </div>
            <div>
              <label>สถานที่คลอด</label>
              <input type="text" v-model="endPregnancyForm.DeliveryPlace" />
            </div>
            <div>
              <label>ภาวะแทรกซ้อน</label>
              <input type="text" v-model="endPregnancyForm.Complications" />
            </div>
            <div>
              <label>สถานะทารก</label>
              <select v-model="endPregnancyForm.ChildStatus">
                <option value="Healthy">แข็งแรง</option>
                <option value="Deceased">เสียชีวิต</option>
                <option value="Other">อื่นๆ</option>
              </select>
            </div>
          </div>
          <div class="flex gap-2 mt-4 justify-end">
            <button type="button" @click="showEndPregnancyModal = false" class="btn-cancel">
              ยกเลิก
            </button>
            <button type="submit" class="btn-save">บันทึก</button>
          </div>
        </form>
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
  font-weight: 500;
  transition: all 0.2s;
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
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-grid input:focus,
.form-grid select:focus,
.form-grid textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.btn-save {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn-save:hover {
  background: #059669;
}

/* End Pregnancy Button */
.btn-end-pregnancy {
  padding: 0.5rem 1rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.875rem;
  transition: background-color 0.2s;
}

.btn-end-pregnancy:hover {
  background-color: #dc2626;
}

/* Modal Styles */
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
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow:
    0 20px 25px -5px rgba(0, 0, 0, 0.1),
    0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.modal-content h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: #1f2937;
}

.btn-save:hover,
.btn-create:hover {
  background: #059669;
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

.mt-4 {
  margin-top: 1rem;
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

/* Visit Form 2-Column Layout */
.visit-form-columns {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 768px) {
  .visit-form-columns {
    grid-template-columns: 1fr 1fr;
  }
}

.visit-form-col {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  font-size: 0.875rem;
  color: #374151;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.625rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 1rem;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
  justify-content: flex-start;
}
</style>
