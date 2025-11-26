<script setup>
import { ref, computed, onMounted } from 'vue'
import { Smile, Calendar, Activity, Baby } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()
const kickCount = ref(0)
const loading = ref(true)

const pregnancyData = computed(() => {
  if (!authStore.user?.Pregnancies || authStore.user.Pregnancies.length === 0) return null
  return authStore.user.Pregnancies.find((p) => p.status === 'Active') || null
})

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('th-TH', { day: 'numeric', month: 'long', year: 'numeric' })
}

const calculateGA = () => {
  if (!pregnancyData.value?.LMP) return '-'
  const lmp = new Date(pregnancyData.value.LMP)
  const now = new Date()
  const diffTime = Math.abs(now - lmp)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  const weeks = Math.floor(diffDays / 7)
  return `${weeks} สัปดาห์`
}

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
    <header class="page-header">
      <div class="user-summary">
        <img
          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
          alt="Avatar"
          class="header-avatar"
        />
        <div>
          <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
          <p>อายุครรภ์: {{ calculateGA() }}</p>
        </div>
      </div>
    </header>

    <div v-if="!pregnancyData" class="empty-state-card">
      <div class="empty-icon-container">
        <Baby class="empty-icon" />
      </div>
      <h3>ไม่พบข้อมูลการตั้งครรภ์ปัจจุบัน</h3>
      <p>คุณยังไม่มีข้อมูลการฝากครรภ์ในระบบ หรือการตั้งครรภ์ได้สิ้นสุดลงแล้ว</p>
      <p class="sub-text">กรุณาติดต่อแพทย์เพื่อสร้างข้อมูลการตั้งครรภ์ใหม่</p>
    </div>

    <div v-else class="summary-cards">
      <div class="card summary-card green">
        <Smile size="32" />
        <div>
          <h3>{{ calculateGA() }}</h3>
          <p>อายุครรภ์วันนี้</p>
        </div>
      </div>

      <div class="card summary-card orange">
        <Calendar size="32" />
        <div>
          <h3>{{ formatDate(pregnancyData.EDC) }}</h3>
          <p>วันนัดคลอดไป</p>
        </div>
      </div>

      <div class="card summary-card blue">
        <Activity size="32" />
        <div>
          <h3>{{ kickCount }} ครั้ง</h3>
          <p>ลูกเตะวันนี้</p>
        </div>
      </div>
    </div>

    <div v-if="pregnancyData" class="card banner">
      <div class="banner-content">
        <h3>การเดินโลของลูกน้อย</h3>
        <p>
          สัปดาห์ที่ {{ calculateGA() }} ลูกเริ่มมีการปรับองค์ประกอบก่อนเหมือนกับจริง
          เริ่มมีการทำงานได้แล้วด้วยเอ้อ!
        </p>
        <button class="btn-primary">อ่านเพิ่มเติ่ม</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 1.5rem;
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
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.summary-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  border-left: 4px solid;
}
.summary-card.green {
  border-left-color: var(--color-primary);
}
.summary-card.orange {
  border-left-color: #f97316;
}
.summary-card.blue {
  border-left-color: #3b82f6;
}
.summary-card h3 {
  font-size: 1.5rem;
  margin: 0;
}
.summary-card p {
  color: var(--color-text-light);
  margin: 0.25rem 0 0;
}
.empty-state-card {
  text-align: center;
  padding: 3rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
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
.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}
</style>
