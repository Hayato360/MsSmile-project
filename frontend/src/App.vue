<script setup>
import { onMounted, computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import DoctorSidebar from './components/DoctorSidebar.vue'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const authStore = useAuthStore()

onMounted(() => {
  authStore.init()
})

// Check if we should show sidebar
const showSidebar = computed(() => {
  return route.meta.requiresAuth !== false && !route.meta.public
})

const isDoctor = computed(() => {
  return authStore.role === 'doctor'
})
</script>

<template>
  <div class="app-container">
    <DoctorSidebar v-if="showSidebar && isDoctor" />
    <Sidebar v-else-if="showSidebar && !isDoctor" />
    <main class="main-content" :class="{ 'no-sidebar': !showSidebar }">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  width: 100%;
}

.main-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  padding: 2rem;
  background-color: var(--color-bg);
  min-height: 100vh;
}

.main-content.no-sidebar {
  margin-left: 0;
}
</style>
