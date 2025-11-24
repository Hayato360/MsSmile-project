import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { public: true },
    },
    // Patient Routes
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    {
      path: '/pregnancy-record',
      name: 'pregnancy-record',
      component: () => import('../views/PregnancyRecordView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    {
      path: '/vaccine',
      name: 'vaccine',
      component: () => import('../views/VaccineView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    {
      path: '/baby-kicking',
      name: 'baby-kicking',
      component: () => import('../views/BabyKickingView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    {
      path: '/health-history',
      name: 'health-history',
      component: () => import('../views/HealthHistoryView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true, role: 'pregnant' },
    },
    // Doctor Routes
    {
      path: '/doctor/dashboard',
      name: 'doctor-dashboard',
      component: () => import('../views/DoctorDashboardView.vue'),
      meta: { requiresAuth: true, role: 'doctor' },
    },
    {
      path: '/doctor/patient/:id',
      name: 'doctor-patient-detail',
      component: () => import('../views/DoctorPatientDetailView.vue'),
      meta: { requiresAuth: true, role: 'doctor' },
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  // If authenticated and trying to access login, redirect based on role
  if (to.path === '/login' && authStore.isAuthenticated) {
    if (authStore.role === 'doctor') {
      next('/doctor/dashboard')
    } else {
      next('/')
    }
    return
  }

  // Check role-based access
  if (to.meta.role && authStore.role !== to.meta.role) {
    // Redirect to appropriate dashboard if accessing wrong role's route
    if (authStore.role === 'doctor') {
      next('/doctor/dashboard')
    } else {
      next('/')
    }
    return
  }

  next()
})

export default router
