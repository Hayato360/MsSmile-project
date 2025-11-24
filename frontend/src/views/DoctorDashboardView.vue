<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Users, Search } from 'lucide-vue-next'
import api from '../services/api'

const router = useRouter()
const patients = ref([])
const searchQuery = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await api.get('/doctor/patients')
    patients.value = response.data
  } catch (error) {
    console.error('Error fetching patients:', error)
  } finally {
    loading.value = false
  }
})

const filteredPatients = computed(() => {
  if (!searchQuery.value) return patients.value
  return patients.value.filter(
    (p) =>
      p.full_name?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      p.hn?.toLowerCase().includes(searchQuery.value.toLowerCase()),
  )
})

const viewPatient = (patientId) => {
  router.push(`/doctor/patient/${patientId}`)
}

const calculateGA = (pregnancies) => {
  if (!pregnancies || pregnancies.length === 0) return '-'

  const pregnancy = pregnancies[pregnancies.length - 1]
  if (!pregnancy.LMP) return '-'

  const lmp = new Date(pregnancy.LMP)
  const today = new Date()
  const diffTime = Math.abs(today - lmp)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  const weeks = Math.floor(diffDays / 7)

  return `${weeks} สัปดาห์`
}
</script>

<template>
  <div class="doctor-dashboard">
    <header class="page-header">
      <h1>รายชื่อคนไข้ทั้งหมด</h1>
      <p>จัดการและติดตามคนไข้ของคุณ</p>
    </header>

    <div class="search-section">
      <div class="search-box">
        <Search size="20" />
        <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อหรือ HN..." />
      </div>
    </div>

    <div class="stats-cards">
      <div class="stat-card">
        <Users size="24" class="stat-icon" />
        <div>
          <p class="stat-label">คนไข้ทั้งหมด</p>
          <p class="stat-value">{{ patients.length }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else class="patients-table">
      <table>
        <thead>
          <tr>
            <th>HN</th>
            <th>ชื่อ-นามสกุล</th>
            <th>อายุ</th>
            <th>อายุครรภ์</th>
            <th>เบอร์โทร</th>
            <th>การดำเนินการ</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="patient in filteredPatients" :key="patient.ID">
            <td>{{ patient.hn }}</td>
            <td>{{ patient.full_name }}</td>
            <td>{{ patient.age }} ปี</td>
            <td>{{ calculateGA(patient.Pregnancies) }}</td>
            <td>{{ patient.phone_number }}</td>
            <td>
              <button @click="viewPatient(patient.ID)" class="btn-view">ดูรายละเอียด</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="filteredPatients.length === 0" class="no-data">ไม่พบข้อมูลคนไข้</div>
    </div>
  </div>
</template>

<style scoped>
.doctor-dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 1.875rem;
  color: var(--color-text);
  margin-bottom: 0.5rem;
}

.page-header p {
  color: var(--color-text-light);
}

.search-section {
  margin-bottom: 1.5rem;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  max-width: 400px;
}

.search-box input {
  border: none;
  outline: none;
  flex: 1;
  font-size: 1rem;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  color: var(--color-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--color-text-light);
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text);
}

.loading {
  text-align: center;
  padding: 3rem;
  color: var(--color-text-light);
}

.patients-table {
  background: white;
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background-color: #f9fafb;
}

th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--color-text);
  border-bottom: 1px solid var(--color-border);
}

td {
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

tbody tr:hover {
  background-color: #f9fafb;
}

.btn-view {
  padding: 0.5rem 1rem;
  background-color: var(--color-primary);
  color: var(--color-text);
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn-view:hover {
  background-color: var(--color-primary-hover);
}

.no-data {
  text-align: center;
  padding: 3rem;
  color: var(--color-text-light);
}
</style>
