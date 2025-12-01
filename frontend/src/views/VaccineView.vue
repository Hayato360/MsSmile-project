<script setup>
import { ref, onMounted } from 'vue'
import { Shield, Syringe, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const vaccinations = ref([])
const loading = ref(true)

onMounted(async () => {
  if (!authStore.user?.ID) {
    loading.value = false
    return
  }

  try {
    const response = await api.get(`/vaccinations/pregnant-woman/${authStore.user.ID}`)
    vaccinations.value = response.data || []
  } catch (error) {
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  if (date.getFullYear() < 1900) return '-'
  return date.toLocaleDateString('th-TH', {
    year: '2-digit',
    month: 'short',
    day: 'numeric',
  })
}

const formatYear = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    if (date.getFullYear() < 1900) return '-'
    return date.getFullYear() + 543 // Thai Year
}
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

    <div v-else class="table-container">
      <table class="vaccine-table">
        <thead>
          <tr>
            <th rowspan="2" class="col-vaccine">วัคซีน</th>
            <th colspan="3" class="col-history">ประวัติการได้รับก่อนตั้งครรภ์</th>
            <th colspan="4" class="col-current">ในระหว่างการตั้งครรภ์นี้</th>
          </tr>
          <tr>
            <th>เคยฉีด</th>
            <th>จำนวน (ครั้ง)</th>
            <th>ครั้งสุดท้าย (ปี)</th>
            <th>เข็มที่ 1</th>
            <th>เข็มที่ 2</th>
            <th>เข็มที่ 3</th>
            <th>หมายเหตุ</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="vaccine in vaccinations" :key="vaccine.ID">
            <td class="vaccine-name">{{ vaccine.VaccineType?.Name || 'ไม่ระบุ' }}</td>
            
            <!-- History -->
            <td class="text-center">
                <span v-if="vaccine.IsHistoryUnknown">ไม่ทราบ</span>
                <span v-else-if="vaccine.IsPreviouslyVaccinated">เคย</span>
                <span v-else>ไม่เคย</span>
            </td>
            <td class="text-center">{{ vaccine.IsPreviouslyVaccinated ? vaccine.PreviousDoses : '-' }}</td>
            <td class="text-center">{{ vaccine.IsPreviouslyVaccinated ? formatYear(vaccine.LastPreviousDateYear) : '-' }}</td>

            <!-- Current Pregnancy -->
            <td class="text-center">{{ formatDate(vaccine.Dose1DateDuringPreg) }}</td>
            <td class="text-center">{{ formatDate(vaccine.Dose2DateDuringPreg) }}</td>
            <td class="text-center">{{ formatDate(vaccine.Dose3DateDuringPreg) }}</td>
            
            <td class="text-center">
                <span v-if="vaccine.ReasonForNotVaccinating" class="text-danger">{{ vaccine.ReasonForNotVaccinating }}</span>
                <span v-else>{{ vaccine.Remarks || '-' }}</span>
            </td>
          </tr>
          <tr v-if="vaccinations.length === 0">
            <td colspan="8" class="text-center py-4 text-muted">ยังไม่มีข้อมูลการฉีดวัคซีน</td>
          </tr>
        </tbody>
      </table>
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

.table-container {
  overflow-x: auto;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  background: white;
}

.vaccine-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 800px;
}

.vaccine-table th,
.vaccine-table td {
  border: 1px solid #e5e7eb;
  padding: 0.75rem;
  font-size: 0.95rem;
}

.vaccine-table th {
  background-color: #f0f9ff; /* Light Blue */
  color: #1e3a8a;
  font-weight: 600;
  text-align: center;
}

.col-history {
  background-color: #fce7f3 !important; /* Light Pink */
  color: #831843 !important;
}

.col-current {
  background-color: #fef3c7 !important; /* Light Yellow */
  color: #92400e !important;
}

.vaccine-name {
  font-weight: 500;
  color: #1f2937;
  background-color: #f8fafc;
}

.text-center {
  text-align: center;
}

.text-muted {
  color: #9ca3af;
}

.text-danger {
    color: #ef4444;
}

.py-4 {
    padding-top: 1rem;
    padding-bottom: 1rem;
}
</style>
