<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBabyKickingStore } from '../stores/babyKicking'
import { useAuthStore } from '../stores/auth'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
} from 'chart.js'
import { Save, Baby } from 'lucide-vue-next'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const store = useBabyKickingStore()
const authStore = useAuthStore()

// Form Data
const selectedDate = ref(new Date().toISOString().split('T')[0])
const morningKicks = ref(0)
const noonKicks = ref(0)
const eveningKicks = ref(0)

// Initialize form with today's data if exists
const initForm = () => {
  const record = store.kicks.find((k) => k.CountDate.split('T')[0] === selectedDate.value)
  if (record) {
    morningKicks.value = record.KickCountMorning
    noonKicks.value = record.KickCountLunch
    eveningKicks.value = record.KickCountEvening
  } else {
    morningKicks.value = 0
    noonKicks.value = 0
    eveningKicks.value = 0
  }
}

onMounted(async () => {
  await store.fetchKicks()
  initForm()
})

const saveKicks = async () => {
  await store.addKickRecord({
    date: selectedDate.value,
    morning: parseInt(morningKicks.value),
    noon: parseInt(noonKicks.value),
    evening: parseInt(eveningKicks.value),
  })
  alert('บันทึกข้อมูลเรียบร้อยแล้ว')
}

// Chart Data
const chartData = computed(() => {
  const labels = store.weeklyData.map((d) => {
    const date = new Date(d.CountDate)
    return `${date.getDate()} พ.ย.` // Simplified Thai month for demo
  })

  return {
    labels,
    datasets: [
      {
        label: 'เช้า',
        backgroundColor: '#A3E635', // Primary
        data: store.weeklyData.map((d) => d.KickCountMorning),
        stack: 'Stack 0',
      },
      {
        label: 'กลางวัน',
        backgroundColor: '#65a30d', // Darker Green
        data: store.weeklyData.map((d) => d.KickCountLunch),
        stack: 'Stack 0',
      },
      {
        label: 'เย็น',
        backgroundColor: '#365314', // Darkest Green
        data: store.weeklyData.map((d) => d.KickCountEvening),
        stack: 'Stack 0',
      },
    ],
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      stacked: true,
      grid: { display: false },
    },
    y: {
      stacked: true,
      beginAtZero: true,
    },
  },
  plugins: {
    legend: {
      position: 'bottom',
    },
  },
}

// Helper for Thai Date
const formatThaiDate = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear() + 543
  const month = date.toLocaleDateString('th-TH', { month: 'long' })
  const day = date.getDate()
  return `${day} ${month} ${year}`
}

const getTotal = (record) =>
  record.KickCountMorning + record.KickCountLunch + record.KickCountEvening
</script>

<template>
  <div class="baby-kicking-view">
    <header class="page-header">
      <div class="user-summary">
        <img
          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
          alt="Avatar"
          class="header-avatar"
        />
        <div>
          <h2>{{ authStore.user?.full_name || 'คุณแม่สมมติ' }}</h2>
          <p v-if="authStore.gestationalAge">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</p>
          <p v-else class="text-muted">ยังไม่มีข้อมูลครรภ์</p>
        </div>
      </div>
    </header>

    <div v-if="!authStore.pregnancyId" class="empty-state-card">
      <div class="empty-icon-container">
        <Baby class="empty-icon" />
      </div>
      <h3>ไม่พบข้อมูลการตั้งครรภ์ปัจจุบัน</h3>
      <p>การนับลูกดิ้นจะเริ่มเมื่อมีการฝากครรภ์และสถานะการตั้งครรภ์เป็น "กำลังตั้งครรภ์"</p>
      <p class="sub-text">หากคุณเพิ่งฝากครรภ์ กรุณาติดต่อเจ้าหน้าที่เพื่อบันทึกข้อมูล</p>
    </div>

    <div v-else-if="authStore.gestationalAge" class="content-grid">
      <!-- Chart Section -->
      <div class="card chart-card">
        <h3>สถิติลูกดิ้นย้อนหลัง</h3>
        <div class="chart-container">
          <Bar :data="chartData" :options="chartOptions" />
        </div>
      </div>

      <!-- Input Section -->
      <div class="card input-card">
        <h3>บันทึกข้อมูลวันนี้</h3>
        <div class="form-group">
          <label>วันที่</label>
          <input type="date" v-model="selectedDate" @change="initForm" />
        </div>

        <div class="kick-inputs">
          <div class="input-group">
            <label>เช้า</label>
            <input type="number" v-model="morningKicks" min="0" />
          </div>
          <div class="input-group">
            <label>กลางวัน</label>
            <input type="number" v-model="noonKicks" min="0" />
          </div>
          <div class="input-group">
            <label>เย็น</label>
            <input type="number" v-model="eveningKicks" min="0" />
          </div>
        </div>

        <button class="btn-primary full-width" @click="saveKicks">
          <Save size="18" />
          บันทึกข้อมูล
        </button>
        <p class="note">*หากรวมกันน้อยกว่า 10 ครั้ง ควรปรึกษาแพทย์</p>
      </div>

      <!-- History Table -->
      <div class="card history-card">
        <h3>ประวัติการนับลูกดิ้น</h3>
        <table class="data-table">
          <thead>
            <tr>
              <th>วันที่</th>
              <th>เช้า</th>
              <th>กลางวัน</th>
              <th>เย็น</th>
              <th>รวม</th>
              <th>สถานะ</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in store.kicks.slice().reverse()" :key="record.CountDate">
              <td>{{ formatThaiDate(record.CountDate) }}</td>
              <td>{{ record.KickCountMorning }}</td>
              <td>{{ record.KickCountLunch }}</td>
              <td>{{ record.KickCountEvening }}</td>
              <td class="font-bold text-primary">{{ getTotal(record) }}</td>
              <td>
                <span class="badge success" v-if="getTotal(record) >= 10">ปกติ</span>
                <span class="badge warning" v-else>ควรเฝ้าระวัง</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-else class="empty-state">
      <div class="empty-icon-wrapper">
        <Baby size="64" class="empty-icon" />
      </div>
      <h3>ยังไม่มีข้อมูลการตั้งครรภ์</h3>
      <p>กรุณาติดต่อแพทย์เพื่อบันทึกข้อมูลการฝากครรภ์ก่อนเริ่มนับลูกดิ้น</p>
    </div>
  </div>
</template>

<style scoped>
.baby-kicking-view {
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

.content-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

.chart-card {
  grid-column: 1 / -1;
  height: 400px;
}

.chart-container {
  height: 320px;
  margin-top: 1rem;
}

.input-card {
  grid-column: 1 / 2;
}

.history-card {
  grid-column: 2 / -1;
}

@media (max-width: 1024px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
  .input-card,
  .history-card {
    grid-column: 1 / -1;
  }
}

.form-group {
  margin-bottom: 1rem;
}

.kick-inputs {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.input-group {
  flex: 1;
}

.input-group label {
  display: block;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  color: var(--color-text-light);
}

.full-width {
  width: 100%;
  justify-content: center;
}

.note {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 0.5rem;
  text-align: center;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.data-table th,
.data-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.data-table th {
  font-weight: 600;
  color: var(--color-text);
}

.text-primary {
  color: var(--color-text); /* Use dark green for text visibility */
}

.font-bold {
  font-weight: 700;
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

.badge.warning {
  background-color: #fee2e2;
  color: #991b1b;
}

.text-muted {
  color: var(--color-text-light);
  font-style: italic;
}

.empty-state-card {
  text-align: center;
  padding: 3rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.empty-icon-container {
  background-color: #ecfccb; /* Lime 100 */
  padding: 1.5rem;
  border-radius: 50%;
  margin-bottom: 0.5rem;
}

.empty-icon {
  width: 48px;
  height: 48px;
  color: #65a30d; /* Lime 600 */
}

.empty-state-card h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.empty-state-card p {
  color: #6b7280;
  margin: 0;
  max-width: 500px;
}

.empty-state-card .sub-text {
  font-size: 0.875rem;
  color: #9ca3af;
}
</style>
