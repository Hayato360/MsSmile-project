import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import { useAuthStore } from './auth'

export const useVaccineStore = defineStore('vaccine', () => {
  const vaccinations = ref([])
  const authStore = useAuthStore()

  async function fetchVaccinations() {
    if (!authStore.user?.ID) return
    try {
      const response = await api.get(`/vaccinations/pregnant-woman/${authStore.user.ID}`)
      vaccinations.value = response.data
    } catch (error) {
      console.error('Error fetching vaccinations:', error)
    }
  }

  async function updateVaccination(id, data) {
    try {
      await api.put(`/vaccinations/${id}`, data)
      await fetchVaccinations()
    } catch (error) {
      console.error('Error updating vaccination:', error)
      alert('เกิดข้อผิดพลาดในการอัปเดตข้อมูล')
    }
  }

  return { vaccinations, fetchVaccinations, updateVaccination }
})
