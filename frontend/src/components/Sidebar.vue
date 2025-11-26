<script setup>
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { LayoutDashboard, FileText, Syringe, Baby, FileHeart, User, LogOut } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'



const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()



const menuItems = [
  { name: 'ภาพรวม (Dashboard)', path: '/', icon: LayoutDashboard },
  { name: 'บันทึกการตรวจครรภ์', path: '/pregnancy-record', icon: FileText },
  { name: 'ประวัติวัคซีน', path: '/vaccine', icon: Syringe },
  { name: 'นับลูกดิ้น (Kick Count)', path: '/baby-kicking', icon: Baby },
  { name: 'ประวัติสุขภาพ', path: '/health-history', icon: FileHeart },
  { name: 'ข้อมูลส่วนตัว', path: '/profile', icon: User },
]

const handleLogout = () => {
  if (confirm('ต้องการออกจากระบบหรือไม่?')) {
    authStore.logout()
  }
}
</script>

<template>
  <aside class="sidebar">
    <div class="logo-container">
      <div class="logo-icon">
        <!-- Simple Leaf Icon SVG -->
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="text-primary"
        >
          <path
            d="M2 22c1.25-.47 2.87-1.03 4.5-2 3.5-2.1 6.5-5.6 7.5-10 .5-2.2-.5-4.5-2.5-5.5C9.5 3.5 7.5 4 6 5.5c-1.5 1.5-2 3.5-1.5 5.5 1 4.4 4.5 7.4 10 7.5 2 .05 4.05-.55 5.5-2"
          />
        </svg>
      </div>
      <div class="logo-text">
        <h1 class="app-name">Healthcare Center</h1>
        <span class="app-subtitle">สมุดฝากครรภ์ดิจิทัล</span>
      </div>
    </div>

    <div class="user-profile">
      <div class="avatar">
        <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix" alt="User Avatar" />
        <div class="status-dot"></div>
      </div>
      <div class="user-info">
        <h3>{{ authStore.user?.full_name || 'คุณแม่สมมติ' }}</h3>
        <p v-if="authStore.gestationalAge">อายุครรภ์: {{ authStore.gestationalAge }} สัปดาห์</p>
        <p v-else class="text-muted">ยังไม่มีข้อมูลครรภ์</p>
      </div>
    </div>

    <nav class="menu">
      <RouterLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="menu-item"
        :class="{ active: route.path === item.path }"
      >
        <component :is="item.icon" class="menu-icon" />
        <span>{{ item.name }}</span>
      </RouterLink>

      <button @click="handleLogout" class="menu-item logout-btn">
        <LogOut class="menu-icon" />
        <span>ออกจากระบบ</span>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  background-color: var(--sidebar-bg);
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  z-index: 10;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
}

.logo-icon {
  color: var(--color-primary);
}

.logo-text h1 {
  font-size: 1.25rem;
  color: var(--color-text);
  line-height: 1.2;
}

.app-subtitle {
  font-size: 0.75rem;
  color: var(--color-text-light);
}

.user-profile {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid var(--color-border);
}

.avatar {
  position: relative;
  width: 80px;
  height: 80px;
  margin-bottom: 1rem;
}

.avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid #f0fdf4;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.status-dot {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 12px;
  height: 12px;
  background-color: var(--color-primary);
  border-radius: 50%;
  border: 2px solid white;
}

.user-info h3 {
  font-size: 1rem;
  margin-bottom: 0.25rem;
}

.user-info p {
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  text-decoration: none;
  color: #4b5563; /* Gray-600 */
  border-radius: 0.5rem;
  transition: all 0.2s;
  font-weight: 500;
}

.menu-item:hover {
  background-color: #f9fafb;
  color: var(--color-text);
}

.menu-item.active {
  background-color: var(--color-sidebar-active);
  color: var(--color-text);
  border-left: 4px solid var(--color-primary);
}

.menu-icon {
  width: 20px;
  height: 20px;
}

.logout-btn {
  margin-top: auto;
  border: none;
  background: none;
  cursor: pointer;
  width: 100%;
  font-size: 1rem;
  color: #dc2626;
}

.logout-btn:hover {
  background-color: #fef2f2;
  color: #991b1b;
}

.text-muted {
  color: var(--color-text-light);
  font-size: 0.875rem;
  font-style: italic;
}
</style>
