import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const role = ref(localStorage.getItem('role') || null)
  const pregnancyId = ref(null) // Store the current pregnancy ID
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  const gestationalAge = computed(() => {
    if (!user.value?.Pregnancies || user.value.Pregnancies.length === 0) return null

    const activePregnancy = user.value.Pregnancies.find((p) => p.status === 'Active')
    if (!activePregnancy || !activePregnancy.LMP) return null

    const lmp = new Date(activePregnancy.LMP)
    const today = new Date()
    const diffTime = Math.abs(today - lmp)
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    const weeks = Math.floor(diffDays / 7)

    return weeks
  })

  async function login(username, password) {
    try {
      const response = await api.post('/login', { username, password })
      const data = response.data.data

      token.value = data.token
      user.value = data.user
      role.value = data.role

      localStorage.setItem('token', data.token)
      localStorage.setItem('role', data.role)

      // Fetch user details to get pregnancy ID
      await fetchMe()

      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function fetchMe() {
    if (!token.value) return

    loading.value = true
    try {
      const response = await api.get('/me')
      const data = response.data.data
      user.value = data.user
      role.value = data.role

      // Update role in localStorage just in case
      localStorage.setItem('role', data.role)

      // Extract Pregnancy ID if available
      if (data.user.Pregnancies && data.user.Pregnancies.length > 0) {
        const activePregnancy = data.user.Pregnancies.find((p) => p.status === 'Active')
        pregnancyId.value = activePregnancy ? activePregnancy.ID : null
      }
    } catch (error) {
      console.error('Fetch user failed:', error)
      logout()
    } finally {
      loading.value = false
    }
  }

  // Initialize user data on app load
  async function init() {
    if (token.value && !user.value) {
      await fetchMe()
    }
  }

  function logout() {
    user.value = null
    token.value = null
    role.value = null
    pregnancyId.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    router.push('/login')
  }

  return {
    user,
    token,
    role,
    pregnancyId,
    isAuthenticated,
    gestationalAge,
    loading,
    login,
    fetchMe,
    init,
    logout,
  }
})
