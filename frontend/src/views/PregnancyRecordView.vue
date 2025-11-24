<script setup>
import { ref, onMounted } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const records = ref([])
const loading = ref(true)

onMounted(async () => {
  if (!authStore.pregnancyId) {
    loading.value = false
    return
  }

  try {
    const response = await api.get(`/antenatal-visits/pregnancy/${authStore.pregnancyId}`)
    records.value = response.data || []
  } catch (error) {
    console.error('Error fetching visits:', error)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('th-TH', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

const calculateGA = (weeks) => {
  const w = Math.floor(weeks)
  const days = Math.round((weeks - w) * 7)
  return `${w}+${days}`
}
</script>

<template>
  <div class="pregnancy-record-view">
    <header class="page-header">
      <div class="user-summary">
        <img
          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
          alt="Avatar"
          class="header-avatar"
        />
        <div>
          <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
          <p>บันทึกการตรวจครรภ์</p>
        </div>
      </div>
    </header>

    <!-- Records Table -->
    <div class="card table-card">
      <div class="card-header">
        <h3>ประวัติการตรวจรายครั้ง</h3>
      </div>

      <div v-if="loading" class="loading">กำลังโหลดข้อมูล...</div>

      <div v-else-if="records.length === 0" class="no-data">ยังไม่มีบันทึกการตรวจครรภ์</div>

      <table v-else class="data-table">
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
          <tr v-for="record in records" :key="record.ID">
            <td>{{ formatDate(record.VisitDate) }}</td>
            <td>{{ record.GestationalAge }} สัปดาห์</td>
            <td>{{ record.Weight }}</td>
            <td>{{ record.BloodPressure }}</td>
            <td>{{ record.FetalHeartSound }}</td>
            <td>{{ record.FetalMovement }}</td>
            <td>{{ record.MedicalDiagnosis || '-' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.pregnancy-record-view {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.user-summary {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: 2px solid var(--color-primary);
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.info-card {
  padding: 1.5rem;
}

.label {
  font-size: 0.875rem;
  color: var(--color-text-light);
  margin-bottom: 0.5rem;
}

.info-card h3 {
  font-size: 1.25rem;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.tag.risk {
  background-color: #fee2e2;
  color: #991b1b;
}

.table-card {
  padding: 0;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-border);
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 1rem 1.5rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.data-table th {
  background-color: #f9fafb;
  font-weight: 600;
  color: var(--color-text-light);
  font-size: 0.875rem;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge.success {
  background-color: #dcfce7;
  color: #166534;
}

.loading,
.no-data {
  text-align: center;
  padding: 3rem;
  color: var(--color-text-light);
}
</style>
