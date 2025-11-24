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

const getStatusColor = (status) => {
  if (status === 'ฉีดแล้ว') return 'success'
  if (status === 'โปรดฉีด') return 'warning'
  return 'default'
}
</script>

<template>
  <div class="vaccine-view">
    <header class="page-header">
      <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
      <p v-if="authStore.gestationalAge">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</p>
      <p v-else class="text-muted">ยังไม่มีข้อมูลครรภ์</p>
    </header>

    <div v-if="loading" class="loading">กำลังโหลดข้อมูล...</div>

    <div v-else-if="vaccinations.length === 0" class="empty-state">
      <Syringe size="64" class="empty-icon" />
      <h3>ยังไม่มีข้อมูลการฉีดวัคซีน</h3>
      <p>กรุณาติดต่อแพทย์เพื่อบันทึกข้อมูลการฉีดวัคซีน</p>
    </div>

    <div v-else class="vaccine-cards">
      <div v-for="vaccine in vaccinations" :key="vaccine.ID" class="card vaccine-card">
        <div class="vaccine-header">
          <Shield size="24" />
          <h3>{{ vaccine.vaccine_name }}</h3>
        </div>
        <div class="vaccine-body">
          <div class="info-row">
            <span class="label">ประวัติการฉีดวัคซีน:</span>
            <span class="value">{{ vaccine.vaccination_history || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">เข้ารับวัคซีน (ครรภ์นี้):</span>
            <span class="value">{{ vaccine.current_pregnancy_date || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">สถานะ:</span>
            <span :class="['badge', getStatusColor(vaccine.status)]">{{ vaccine.status }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vaccine-view {
  max-width: 1200px;
  margin: 0 auto;
}
.page-header {
  margin-bottom: 2rem;
}
.page-header h2 {
  margin: 0;
}
.page-header p {
  color: var(--color-text-light);
  margin: 0.5rem 0 0;
}
.loading,
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 0.5rem;
}
.empty-icon {
  color: var(--color-text-light);
  margin-bottom: 1rem;
}
.empty-state h3 {
  margin: 0 0 0.5rem;
}
.empty-state p {
  color: var(--color-text-light);
  margin: 0;
}
.vaccine-cards {
  display: grid;
  gap: 1.5rem;
}
.vaccine-card {
  padding: 0;
  overflow: hidden;
}
.vaccine-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}
.vaccine-header h3 {
  margin: 0;
}
.vaccine-body {
  padding: 1.5rem;
}
.info-row {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--color-border);
}
.info-row:last-child {
  border-bottom: none;
}
.label {
  color: var(--color-text-light);
}
.value {
  font-weight: 500;
}
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
}
.badge.success {
  background: #dcfce7;
  color: #166534;
}
.badge.warning {
  background: #fef3c7;
  color: #92400e;
}
.text-muted {
  color: var(--color-text-light);
  font-style: italic;
}
</style>
