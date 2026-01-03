<script setup>
import { ref, onMounted } from 'vue'
import { Shield, Syringe, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'
import VaccinationCard from '../components/VaccinationCard.vue'

const authStore = useAuthStore()
const vaccinations = ref([])
const loading = ref(true)

// Encryption of data state
const vaccineTypes = ref([])
const vaccinationForms = ref({})

onMounted(async () => {
  if (!authStore.user?.ID) {
    loading.value = false
    return
  }

  try {
     // 1. Fetch Vaccine Types
    const vTypeRes = await api.get('/vaccine-types')
    vaccineTypes.value = vTypeRes.data.data || []
    
    // 2. Initialize Forms
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
            Doses: []
        }
    })

    // 3. Fetch User Vaccinations
    const response = await api.get(`/vaccinations/pregnant-woman/${authStore.user.ID}`)
    vaccinations.value = response.data || []

    // 4. Map Data
    vaccinations.value.forEach(v => {
         if (vaccinationForms.value[v.VaccineTypeID]) {
            const vData = { ...v }
            // Format dates for inputs
            if (vData.LastPreviousDateYear) vData.LastPreviousDateYear = vData.LastPreviousDateYear.split('T')[0]
            if (vData.Dose1DateDuringPreg) vData.Dose1DateDuringPreg = vData.Dose1DateDuringPreg.split('T')[0]
            if (vData.Dose2DateDuringPreg) vData.Dose2DateDuringPreg = vData.Dose2DateDuringPreg.split('T')[0]
            if (vData.Dose3DateDuringPreg) vData.Dose3DateDuringPreg = vData.Dose3DateDuringPreg.split('T')[0]
            
            // Map Doses array dates
             if (vData.Doses && vData.Doses.length > 0) {
                 vData.Doses = vData.Doses.map(d => ({
                     ...d,
                     DoseDate: d.DoseDate ? d.DoseDate.split('T')[0] : ''
                 }))
             }

            vaccinationForms.value[v.VaccineTypeID] = { ...vaccinationForms.value[v.VaccineTypeID], ...vData }
         }
    })

  } catch (error) {
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
})

</script>

<template>
  <div class="vaccine-view">
    <header class="page-header">
      <h2>ประวัติวัคซีนหญิงตั้งครรภ์</h2>
      <div class="user-info">
        <span class="name">{{ authStore.user?.full_name || 'คุณแม่' }}</span>
        <span v-if="authStore.gestationalAge" class="ga">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</span>
      </div>
    </header>

    <div v-if="loading" class="loading">กำลังโหลดข้อมูล...</div>

    <div v-else class="content-container">
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
</template>

<style scoped>
.vaccine-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
  font-family: 'Sarabun', sans-serif;
}

.page-header {
  margin-bottom: 2rem;
  text-align: center;
}

.page-header h2 {
  color: #1e40af;
  margin-bottom: 0.5rem;
}

.user-info {
  font-size: 1.1rem;
  color: #4b5563;
}

.content-container {
    padding: 1rem;
}

.vaccination-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 1.5rem;
}

.text-center {
  text-align: center;
}

.text-muted {
  color: #9ca3af;
}
</style>
