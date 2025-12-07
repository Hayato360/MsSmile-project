<script setup>
import { onMounted, ref } from 'vue'
import { useHealthHistoryStore } from '../stores/healthHistory'
import { useAuthStore } from '../stores/auth'
import { FileText } from 'lucide-vue-next'

const store = useHealthHistoryStore()
const authStore = useAuthStore()
const isLoading = ref(true)

onMounted(async () => {
  try {
    await Promise.all([store.fetchHealthHistory(), store.fetchPreviousPregnancies()])
  } finally {
    isLoading.value = false
  }
})

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <div class="health-history-view">
    <header class="page-header">
      <div class="user-summary">
        <img
          src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
          alt="Avatar"
          class="header-avatar"
        />
        <div>
          <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
          <p v-if="authStore.gestationalAge">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</p>
          <p v-else class="text-muted">ยังไม่มีข้อมูลครรภ์</p>
        </div>
      </div>
    </header>

    <div v-if="isLoading" class="empty-state">
      <h3>กำลังโหลดข้อมูล...</h3>
    </div>

    <div v-else-if="store.history" class="content-grid">
      <!-- Section 1: Contraception & Menstruation -->
      <div class="card section-card">
        <h3>ประวัติสุขภาพของหญิงตั้งครรภ์และครอบครัว</h3>

        <div class="info-grid">
          <div class="info-group">
            <label>ก่อนตั้งครรภ์คุมกำเนิดวิธี:</label>
            <span>{{ store.history.ContraceptionBeforeMethod || '-' }}</span>
          </div>
          <div class="info-group">
            <label>เป็นระยะเวลา:</label>
            <span>{{ store.history.ContraceptionBeforeDuration || '-' }}</span>
          </div>

          <div class="info-group">
            <label>หยุดคุมกำเนิดครั้งหลังสุด:</label>
            <span>{{ store.history.ContraceptionLastMethod || '-' }}</span>
          </div>
          <div class="info-group">
            <label>ระยะเวลา:</label>
            <span>{{ store.history.ContraceptionLastDuration || '-' }}</span>
          </div>

          <div class="info-group full-width">
            <label>ประวัติประจำเดือนมาทุก:</label>
            <span>{{ store.history.MenstrualCycle || '-' }} วัน</span>
            <label class="ml-4">นานครั้งละ:</label>
            <span>{{ store.history.MenstrualDuration || '-' }} วัน</span>
            <span
              class="ml-4 badge"
              :class="store.history.MenstrualCondition === 'สม่ำเสมอ' ? 'success' : 'warning'"
            >
              {{ store.history.MenstrualCondition || '-' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Section 2: Previous Pregnancies -->
      <div class="card section-card">
        <h3>ประวัติการตั้งครรภ์</h3>
        <div class="table-responsive">
          <table class="data-table">
            <thead>
              <tr>
                <th>ครรภ์ที่</th>
                <th>ว/ด/ป คลอด/แท้ง</th>
                <th>อายุครรภ์ (สัปดาห์)</th>
                <th>วิธีคลอด/แท้ง</th>
                <th>น้ำหนักทารก (กก.)</th>
                <th>เพศ</th>
                <th>สถานที่คลอด/แท้ง</th>
                <th>ภาวะแทรกซ้อน</th>
                <th>สภาพทารกปัจจุบัน</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="preg in store.previousPregnancies" :key="preg.ID">
                <td>{{ preg.pregnancy_no }}</td>
                <td>{{ formatDate(preg.delivery_date) }}</td>
                <td>{{ preg.gestational_age }}</td>
                <td>{{ preg.delivery_method }}</td>
                <td>{{ preg.birth_weight }}</td>
                <td>{{ preg.sex }}</td>
                <td>{{ preg.delivery_place }}</td>
                <td>{{ preg.complications || '-' }}</td>
                <td>{{ preg.child_status }}</td>
              </tr>
              <tr v-if="store.previousPregnancies.length === 0">
                <td colspan="9" class="text-center text-muted">ไม่มีประวัติการตั้งครรภ์ในอดีต</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Section 3: Medical History Checkboxes -->
      <div class="card section-card">
        <h3>ประวัติเจ็บป่วยของหญิงตั้งครรภ์</h3>
        <div class="info-grid">
          <div class="info-group full-width">
            <label>โรคประจำตัว:</label>
            <span>{{ store.history.ChronicDiseases || '-' }}</span>
          </div>
          <div class="info-group full-width">
            <label>อื่นๆ:</label>
            <span>{{ store.history.OtherDiseases || '-' }}</span>
          </div>
        </div>

        <div class="info-row mt-4">
          <label>ประวัติผ่าตัดคลอด:</label>
          <span>{{ store.history.SurgeryHistory || '-' }}</span>
        </div>
        <div class="info-row">
          <label>อื่นๆ:</label>
          <span>{{ store.history.OtherSurgery || '-' }}</span>
        </div>

        <div class="info-row mt-2">
          <label>ประวัติแพ้ยา:</label>
          <span class="text-danger">{{ store.history.DrugAllergies || 'ปฏิเสธ' }}</span>
        </div>
      </div>

      <!-- Section 4: Family History -->
      <div class="card section-card">
        <h3>ประวัติการเจ็บป่วยของบุคคลในครอบครัว</h3>
        <div class="info-grid">
          <div class="info-group full-width">
            <label>ประวัติครอบครัว:</label>
            <span>{{ store.history.OtherFamilyHistory || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <FileText size="64" class="empty-icon" />
      <h3>ยังไม่มีข้อมูลประวัติสุขภาพ</h3>
      <p>กรุณาติดต่อแพทย์เพื่อบันทึกข้อมูลประวัติสุขภาพ โรคประจำตัว และผลตรวจทางห้องปฏิบัติการ</p>
    </div>
  </div>
</template>

<style scoped>
.health-history-view {
  max-width: 1200px;
  margin: 0 auto;
  padding-bottom: 2rem;
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
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.section-card {
  padding: 1.5rem;
}

.section-card h3 {
  color: var(--color-primary);
  border-bottom: 2px solid var(--color-background);
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-group {
  display: flex;
  gap: 0.5rem;
  align-items: baseline;
}

.info-group label {
  font-weight: 600;
  color: var(--color-text-light);
  white-space: nowrap;
}

.full-width {
  grid-column: 1 / -1;
}

.ml-4 {
  margin-left: 1rem;
}

.table-responsive {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.data-table th,
.data-table td {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  text-align: center;
}

.data-table th {
  background-color: var(--color-background);
  font-weight: 600;
}

.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.checkbox-item input[type='checkbox'] {
  accent-color: var(--color-primary);
  width: 1.2rem;
  height: 1.2rem;
}

.info-row {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.info-row label {
  font-weight: 600;
}

.text-danger {
  color: #ef4444;
  font-weight: 600;
}

.badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

.badge.success {
  background-color: #dcfce7;
  color: #166534;
}

.badge.warning {
  background-color: #fee2e2;
  color: #991b1b;
}

.text-center {
  text-align: center;
}

.text-muted {
  color: var(--color-text-light);
}

.empty-state {
  text-align: center;
  padding: 4rem;
}
</style>
