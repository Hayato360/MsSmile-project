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

const getGestationalWeeks = () => {
  if (!pregnancyData.value?.LMP) return 0
  const lmp = new Date(pregnancyData.value.LMP)
  const now = new Date()
  const diffTime = Math.abs(now - lmp)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return Math.floor(diffDays / 7)
}

const weeklyThemes = [
  {
    range: [0, 12],
    gradient: 'linear-gradient(135deg, #a8e063 0%, #56ab2f 100%)',
    title: 'ไตรมาสที่ 1: การเริ่มต้น',
    description: 'ช่วงเวลาสำคัญของการสร้างอวัยวะต่างๆ ของลูกน้อย เซลล์กำลังแบ่งตัวอย่างรวดเร็วเพื่อสร้างรากฐานของชีวิต'
  },
  {
    range: [13, 16],
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    title: 'การเติบโตของลูกน้อย (ไตรมาสที่ 1)',
    description: 'ลูกเริ่มมีการปรับองค์ประกอบร่างกาย เริ่มมีการทำงานของระบบต่างๆ และเริ่มขยับตัวได้แล้ว!'
  },
  {
    range: [17, 26],
    gradient: 'linear-gradient(135deg, #00b09b 0%, #96c93d 100%)',
    title: 'ไตรมาสที่ 2: การเติบโตอย่างรวดเร็ว',
    description: 'ลูกน้อยเริ่มขยับตัวแรงขึ้นจนคุณแม่รู้สึกได้ อวัยวะต่างๆ พัฒนาสมบูรณ์ขึ้นเรื่อยๆ'
  },
  {
    range: [27, 40],
    gradient: 'linear-gradient(135deg, #f2994a 0%, #f2c94c 100%)',
    title: 'ไตรมาสที่ 3: โค้งสุดท้าย',
    description: 'ลูกน้อยกำลังสะสมไขมันและเตรียมความพร้อมสำหรับการลืมตาดูโลกในอีกไม่ช้า'
  }
]

const currentWeekInfo = computed(() => {
  const weeks = getGestationalWeeks()
  
  // Specific week overrides (optional, can add more detailed weekly data here)
  const specificWeeks = {
    15: {
      title: 'สัปดาห์ที่ 15: ลูกน้อยเริ่มรับรู้แสง',
      description: 'แม้เปลือกตาจะยังปิดอยู่ แต่ลูกเริ่มรับรู้แสงสว่างได้แล้วนะ! ขนอ่อนๆ เริ่มขึ้นปกคลุมร่างกาย',
      gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)'
    }
  }

  if (specificWeeks[weeks]) {
    return specificWeeks[weeks]
  }

  // Fallback to range themes
  const theme = weeklyThemes.find(t => weeks >= t.range[0] && weeks <= t.range[1])
  
  return theme || {
    title: `สัปดาห์ที่ ${weeks}`,
    description: 'พัฒนาการของลูกน้อยกำลังดำเนินไปอย่างต่อเนื่อง',
    gradient: 'linear-gradient(135deg, #56ab2f 0%, #a8e063 100%)'
  }
})

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
          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Lilly"
          alt="Avatar"
          class="header-avatar"
          width="60"
          height="60"
        />
        <div>
          <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
          <p>อายุครรภ์: {{ calculateGA() }}</p>
        </div>
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

    <div v-if="pregnancyData" class="card banner" :style="{ background: currentWeekInfo.gradient }">
      <div class="banner-content">
        <h3>{{ currentWeekInfo.title }}</h3>
        <p>{{ currentWeekInfo.description }}</p>
        <button class="btn-primary">อ่านเพิ่มเติม</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ... existing styles ... */
.banner {
  /* Default gradient as fallback */
  background: linear-gradient(135deg, #56ab2f 0%, #a8e063 100%);
  color: white;
  transition: background 0.5s ease;
}
</style>
