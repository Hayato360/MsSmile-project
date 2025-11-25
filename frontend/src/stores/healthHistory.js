import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import { useAuthStore } from './auth'

export const useHealthHistoryStore = defineStore('healthHistory', () => {
  const history = ref(null)
  const previousPregnancies = ref([])
  const authStore = useAuthStore()

  async function fetchHealthHistory() {
    if (!authStore.user?.ID) return
    try {
      const response = await api.get(`/medical-histories/pregnant-woman/${authStore.user.ID}`)
      history.value = response.data.data
    } catch (error) {
      console.error('Error fetching health history:', error)
    }
  }

  async function fetchPreviousPregnancies() {
    if (!authStore.user?.ID) return
    try {
      const response = await api.get(`/previous-pregnancies/pregnant-woman/${authStore.user.ID}`)
      previousPregnancies.value = response.data.data
    } catch (error) {
      console.error('Error fetching previous pregnancies:', error)
    }
  }

  async function updateHealthHistory(id, data) {
    try {
      await api.put(`/medical-histories/${id}`, data)
      await fetchHealthHistory()
    } catch (error) {
      console.error('Error updating health history:', error)
      alert('เกิดข้อผิดพลาดในการอัปเดตข้อมูล')
    }
  }

  return {
    history,
    previousPregnancies,
    fetchHealthHistory,
    fetchPreviousPregnancies,
    updateHealthHistory,
  }
})
