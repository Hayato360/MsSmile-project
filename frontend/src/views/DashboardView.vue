<script setup>
import { ref, computed, onMounted } from 'vue'
import { Smile, Calendar, Footprints, ArrowRight, Baby } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const kickCount = ref(0)
const loading = ref(true)

const pregnancyData = computed(() => {
  if (!authStore.user?.Pregnancies || authStore.user.Pregnancies.length === 0) return null
  return authStore.user.Pregnancies.find((p) => p.status === 'Active') || null
})

const formatTime = (dateString) => {
  if (!dateString) return '09:00 น.'
  const date = new Date(dateString)
  return date.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) + ' น.'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  // Buddhist Era (BE) year = AD + 543
  const options = { day: 'numeric', month: 'short', year: 'numeric' }
  const formatted = date.toLocaleDateString('th-TH', options)
  // Customize if needed to match "25 พ.ย. 68" format precisely
  return formatted
}

const calculateGA = () => {
  if (!pregnancyData.value?.LMP) return '-'
  const lmp = new Date(pregnancyData.value.LMP)
  const now = new Date()
  const diffTime = Math.abs(now - lmp)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  const weeks = Math.floor(diffDays / 7)
  return weeks
}

// Mock Appointment Data (Replace with real API if available)
const nextAppointment = ref({})

onMounted(async () => {
  if (authStore.pregnancyId) {
    try {
      const response = await api.get(`/kick-counts/pregnancy/${authStore.pregnancyId}`)
      const records = response.data || []
      const today = new Date().toISOString().split('T')[0]
      const todayRecord = records.find((r) => r.RecordDate?.startsWith(today))
      if (todayRecord) {
        kickCount.value =
          (todayRecord.KickCountMorning || 0) +
          (todayRecord.KickCountAfternoon || 0) +
          (todayRecord.KickCountEvening || 0)
      }
    } catch (error) {
      console.error('Error fetching kick count:', error)
    }
  }
  loading.value = false
})
</script>

<template>
  <div class="dashboard-view">
    <!-- Header Section -->
    <header class="page-header">
      <div class="welcome-text">
        <h1>สวัสดี, {{ authStore.user?.full_name || 'คุณแม่สมมติ' }}</h1>
        <p>ขอให้วันนี้เป็นวันที่สดใสและแข็งแรง</p>
      </div>
      <div v-if="pregnancyData" class="due-date-badge">
        <Calendar size="16" />
        กำหนดคลอด: {{ formatDate(pregnancyData.EDC) }}
      </div>
    </header>

    <div v-if="!pregnancyData" class="empty-state-card">
      <div class="empty-content">
        <div class="empty-icon-wrapper">
          <Baby class="empty-icon" size="48" />
        </div>
        <h3>ยินดีต้อนรับสู่ Pregnanzy!</h3>
        <p>เริ่มต้นการเดินทางที่แสนพิเศษของคุณกับเรา</p>
        <div class="info-box">
          <p>ขณะนี้คุณยังไม่มีข้อมูลการฝากครรภ์ในระบบ</p>
          <p class="sub-text">กรุณาติดต่อแพทย์เพื่อสร้างสมุดฝากครรภ์ดิจิทัลของคุณ</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else class="content-wrapper">
      <!-- Stats Cards -->
      <div class="stats-grid">
        <!-- GA Card -->
        <div class="stat-card green-card">
          <div class="icon-circle green-icon">
             <Smile size="32" />
          </div>
          <p class="card-label">อายุครรภ์ปัจจุบัน</p>
          <h3 class="card-value">{{ calculateGA() }} <span class="unit">สัปดาห์</span></h3>
        </div>

        <!-- Appointment Card -->
        <div class="stat-card orange-card">
          <div class="icon-circle orange-icon">
            <Calendar size="32" />
          </div>
          <p class="card-label">นัดหมายครั้งถัดไป</p>
          <div v-if="nextAppointment.appointment_date || nextAppointment.AppointmentDate">
            <h3 class="card-value">{{ formatDate(nextAppointment.appointment_date || nextAppointment.AppointmentDate) }}</h3>
            <p class="card-subtext">{{ nextAppointment.time || formatTime(nextAppointment.appointment_date || nextAppointment.AppointmentDate) }}</p>
          </div>
          <div v-else>
            <h3 class="card-value text-gray-400">-</h3>
             <p class="card-subtext">ไม่มีการนัดหมาย</p>
          </div>
        </div>

        <!-- Kick Count Card -->
        <div class="stat-card blue-card">
          <div class="icon-circle blue-icon">
            <Footprints size="32" />
          </div>
          <p class="card-label">ลูกดิ้นวันนี้</p>
          <h3 class="card-value">{{ kickCount }} <span class="unit">ครั้ง</span></h3>
          <p class="status-badge success">
            <span class="dot"></span> เกณฑ์ปกติ
          </p>
        </div>
      </div>

      <!-- Growth Banner -->
      <div class="growth-banner">
        <div class="banner-overlay"></div>
        <div class="banner-content">
          <h4>การเติบโตของลูกน้อย</h4>
          <p class="week-desc">สัปดาห์ที่ {{ calculateGA() }}: ลูกน้อยมีขนาดประมาณผลแก้วมังกร เริ่มมีการฝึกหายใจและลืมตาได้แล้ว</p>
          <button class="btn-more">
            ดูเพิ่มเติม
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-view {
  max-width: 1200px;
  margin: 0 auto;
  font-family: 'Prompt', sans-serif; /* Ensuring compatible font */
}

/* Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  background: transparent;
}

.welcome-text h1 {
  font-size: 1.8rem;
  color: #3f6212; /* Dark Green */
  margin-bottom: 0.5rem;
  font-weight: 700;
}

.welcome-text p {
  color: #64748b;
  font-size: 1rem;
}

.due-date-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: white;
  padding: 0.5rem 1rem;
  border-radius: 999px;
  border: 1px solid #e2e8f0;
  color: #65a30d;
  font-size: 0.875rem;
  font-weight: 500;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: white;
  border-radius: 1.5rem;
  padding: 2rem 1.5rem;
  text-align: center;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  align-items: center;
  border-top: 4px solid transparent;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.green-card { border-top-color: #84cc16; }
.orange-card { border-top-color: #f97316; }
.blue-card { border-top-color: #3b82f6; }

.icon-circle {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.green-icon { background-color: #ecfccb; color: #84cc16; }
.orange-icon { background-color: #ffedd5; color: #f97316; }
.blue-icon { background-color: #dbeafe; color: #3b82f6; }

.card-label {
  color: #64748b;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.card-value {
  font-size: 2rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.5rem;
  line-height: 1.2;
}

.unit {
  font-size: 1rem;
  color: #94a3b8;
  font-weight: 500;
}

.card-subtext {
  color: #64748b;
  font-size: 0.95rem;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}

.status-badge.success {
  background-color: #dcfce7;
  color: #166534;
}

.dot {
  width: 6px;
  height: 6px;
  background-color: currentColor;
  border-radius: 50%;
}

/* Growth Banner */
.growth-banner {
  position: relative;
  background-image: url('https://images.unsplash.com/photo-1518173946687-a4c8892bbd9f?q=80&w=2574&auto=format&fit=crop'); /* Forest background */
  background-size: cover;
  background-position: center;
  border-radius: 2rem;
  overflow: hidden;
  height: 280px;
  display: flex;
  align-items: flex-end;
  color: white;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0) 60%);
}

.banner-content {
  position: relative;
  z-index: 10;
  padding: 2.5rem;
  width: 100%;
}

.banner-content h4 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.week-desc {
  font-size: 1rem;
  opacity: 0.9;
  max-width: 600px;
  margin-bottom: 1.5rem;
}

.btn-more {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 9999px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-more:hover {
  background: #2563eb;
}

/* Empty State */
.empty-state-card {
  background: white;
  border-radius: 20px;
  padding: 4rem 2rem;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  margin-top: 1rem;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
.empty-content { max-width: 480px; padding: 2rem; }
.empty-icon-wrapper {
  width: 100px; height: 100px; background: #e0f7fa; border-radius: 50%;
  display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; color: #00838f;
}
.empty-content h3 { font-size: 1.5rem; margin-bottom: 0.5rem; }
.info-box { background: #f8f9fa; padding: 1.5rem; border-radius: 12px; margin-top: 2rem; }
</style>
