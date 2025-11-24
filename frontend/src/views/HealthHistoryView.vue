<script setup>
import { ref } from 'vue'
import { FileText, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="health-history-view">
    <header class="page-header">
      <h2>{{ authStore.user?.full_name || 'คุณแม่' }}</h2>
      <p v-if="authStore.gestationalAge">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</p>
      <p v-else class="text-muted">ยังไม่มีข้อมูลครรภ์</p>
    </header>

    <div class="empty-state">
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
  max-width: 500px;
  margin: 0 auto;
}
.text-muted {
  color: var(--color-text-light);
  font-style: italic;
}
</style>
