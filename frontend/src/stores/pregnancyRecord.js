import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import { useAuthStore } from './auth'

export const usePregnancyRecordStore = defineStore('pregnancyRecord', () => {
  const visits = ref([])
  const authStore = useAuthStore()

  async function fetchVisits() {
    if (!authStore.pregnancyId) return
    try {
      const response = await api.get(`/antenatal-visits/pregnancy/${authStore.pregnancyId}`)
      visits.value = response.data
    } catch (error) {
      console.error('Error fetching visits:', error)
    }
  }

  async function addVisit(data) {
    if (!authStore.pregnancyId) {
      alert('ไม่พบข้อมูลการตั้งครรภ์')
      return
    }

    try {
      const payload = {
        ...data,
        PregnancyID: authStore.pregnancyId,
      }
      await api.post('/antenatal-visits', payload)
      await fetchVisits()
    } catch (error) {
      console.error('Error adding visit:', error)
      alert('เกิดข้อผิดพลาดในการบันทึกข้อมูล')
    }
  }

  return { visits, fetchVisits, addVisit }
})
