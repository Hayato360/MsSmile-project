import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import { useAuthStore } from './auth'

export const useBabyKickingStore = defineStore('babyKicking', () => {
  const kicks = ref([])
  const authStore = useAuthStore()

  const todayKicks = computed(() => {
    const today = new Date().toISOString().split('T')[0]
    return (
      kicks.value.find((k) => k.CountDate.split('T')[0] === today) || {
        CountDate: today,
        KickCountMorning: 0,
        KickCountLunch: 0,
        KickCountEvening: 0,
      }
    )
  })

  const weeklyData = computed(() => {
    // Return last 7 days (or available data)
    return kicks.value.slice(-7)
  })

  async function fetchKicks() {
    if (!authStore.pregnancyId) return
    try {
      const response = await api.get(`/kick-counts/pregnancy/${authStore.pregnancyId}`)
      kicks.value = response.data
    } catch (error) {
      console.error('Error fetching kicks:', error)
    }
  }

  async function addKickRecord(record) {
    if (!authStore.pregnancyId) {
      alert('ไม่พบข้อมูลการตั้งครรภ์')
      return
    }

    try {
      const payload = {
        PregnancyID: authStore.pregnancyId,
        CountDate: new Date(record.date).toISOString(),
        KickCountMorning: record.morning,
        KickCountLunch: record.noon,
        KickCountEvening: record.evening,
      }

      const response = await api.post('/kick-counts', payload)

      // Update local state
      const newRecord = response.data.data
      const index = kicks.value.findIndex((k) => k.CountDate.split('T')[0] === record.date)
      if (index !== -1) {
        kicks.value[index] = newRecord
      } else {
        kicks.value.push(newRecord)
      }
    } catch (error) {
      console.error('Error adding kick record:', error)
      alert('เกิดข้อผิดพลาดในการบันทึกข้อมูล')
    }
  }

  return { kicks, todayKicks, weeklyData, fetchKicks, addKickRecord }
})
