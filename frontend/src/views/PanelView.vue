<template>
  <div class="panel-bg">
    <div class="panel">
      <!-- Sidebar Navigation -->
      <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <div class="sidebar-header" @click="toggleSidebar" role="button" :title="sidebarCollapsed ? 'Rozwiń menu' : 'Zwiń menu'">
          <span class="logo" v-show="!sidebarCollapsed">TRAINING BUDDY</span>
          <div class="user-role" :class="store.user?.role" v-show="!sidebarCollapsed">
            {{ store.user?.role === 'trainer' ? 'TRENER' : 'UCZEŃ' }}
          </div>
        </div>

        <nav class="sidebar-nav">
          <button @click="activeView = 'dashboard'" :class="{ active: activeView === 'dashboard' }">
            <span class="nav-icon">🏠</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Dashboard</span>
          </button>
          <button @click="goToProfile" :class="{ active: activeView === 'profile' }">
            <span class="nav-icon">👤</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Profil</span>
          </button>
          <button @click="toggleCalculator" :class="{ active: activeView === 'calculator' }">
            <span class="nav-icon">🧮</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Kalkulator</span>
          </button>
          <button @click="activeView = 'diet-plans'" :class="{ active: activeView === 'diet-plans' }">
            <span class="nav-icon">🍽️</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Diety</span>
          </button>
          <!-- ZMIENIONY PRZYCISK - teraz to są PLANY TRENINGOWE -->
          <button @click="activeView = 'training-plans'" :class="{ active: activeView === 'training-plans' }">
            <span class="nav-icon">🏋️</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Plany Treningowe</span>
          </button>
          <button @click="activeView = 'media'" :class="{ active: activeView === 'media' }">
            <span class="nav-icon">🎥</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Media</span>
          </button>
          <button @click="activeView = 'reports'" :class="{ active: activeView === 'reports' }">
            <span class="nav-icon">📊</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Raporty</span>
          </button>
          <button @click="activeView = 'chat'" :class="{ active: activeView === 'chat' }">
            <span class="nav-icon">💬</span>
            <span class="nav-text" v-show="!sidebarCollapsed">Chat</span>
          </button>
          <button @click="activeView = 'about'" :class="{ active: activeView === 'about' }">
            <span class="nav-icon">ℹ️</span>
            <span class="nav-text" v-show="!sidebarCollapsed">O nas</span>
          </button>
        </nav>

        <div class="sidebar-footer">
          <button @click="logout" class="logout-btn" :title="sidebarCollapsed ? 'Wyloguj' : ''">
            🚪
          </button>
        </div>
      </aside>

      <!-- Main Content -->
      <main class="main-content">
        <!-- O nas -->
        <AboutView v-if="activeView === 'about'" />

        <!-- Profil -->
        <template v-else-if="activeView === 'profile'">
          <TrainerProfileView
              v-if="store.user?.role === 'trainer'"
          />
          <UserProfileView
              v-else
          />
        </template>

        <!-- Kalkulator -->
        <CalorieDashboard
            v-else-if="activeView === 'calculator'"
        />

        <!-- Plany Żywieniowe -->
        <DietPlansView v-else-if="activeView === 'diet-plans'" />

        <!-- PLANY TRENINGOWE - zmieniony widok -->
        <TrainingPlansView v-else-if="activeView === 'training-plans'" />

        <!-- Media -->
        <MediaView v-else-if="activeView === 'media'" />

        <!-- Raporty -->
        <ReportsView v-else-if="activeView === 'reports'" />

        <!-- Chat -->
        <ChatView v-else-if="activeView === 'chat'" />

        <!-- Dashboardy -->
        <template v-else>
          <TrainerDashboard v-if="store.user?.role === 'trainer'" />
          <TraineeDashboard  v-else />
        </template>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import TrainerProfileView from "@/views/TrainerProfileView.vue";
import UserProfileView from "@/views/UserProfileView.vue"
import TrainerDashboard from "@/components/TrainerDashboard.vue"
import TraineeDashboard from "@/components/TraineeDashboard.vue"
import AboutView from "@/views/AboutView.vue"
import DietPlansView from "@/views/DietPlansView.vue"
import MediaView from "@/views/MediaView.vue"
import ReportsView from "@/views/ReportsView.vue"
import ChatView from "@/views/ChatView.vue"
import {onMounted, provide, ref} from "vue"
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import CalorieDashboard from '@/components/CalorieDashboard.vue'
import TrainingPlansView from "@/views/TrainingPlansView.vue";

const store = useAuthStore()
const { profile, user, token, fetchUser, fetchProfile } = store
const router = useRouter()

// ZAKTUALIZOWANY TYP - dodano training-plans, usunięto exercises
const activeView = ref<'dashboard' | 'profile' | 'calculator' | 'about' | 'diet-plans' | 'training-plans' | 'media' | 'reports' | 'chat'>('dashboard')
const sidebarCollapsed = ref(false)

onMounted(async () => {
  if (!user && token) {
    await fetchUser()
  }

  if (user?.role === 'trainee') {
    await fetchProfile()
  }
})

const switchToView = (viewName: string) => {
  console.log('Switching to:', viewName) // DEBUG

  // ZAKTUALIZOWANE mapowanie nazw widoków:
  const viewMapping: Record<string, string> = {
    'training-plans': 'training-plans',  // Plany treningowe
    'user-media': 'media',
    'diet-plans': 'diet-plans',
    'messages': 'chat',
    'trainer-profile': 'profile',
    'reports': 'reports'
  }

  const targetView = viewMapping[viewName] || viewName
  activeView.value = targetView as any
}

provide('switchToView', switchToView)

function logout() {
  store.logout()
  router.push('/')
}

function goToProfile() {
  activeView.value = 'profile'
}

function toggleCalculator() {
  activeView.value = activeView.value === 'calculator' ? 'dashboard' : 'calculator'
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<!-- Style pozostają takie same jak masz -->
<style scoped>
.panel-bg {
  background: url('@/assets/bgdash.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  min-height: 100vh;
  width: 100vw;
  position: fixed;
  top: 0;
  left: 0;
  padding: 0;
  margin: 0;
}

.panel {
  width: 100%;
  height: 100vh;
  display: flex;
  border-radius: 0;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border: none;
  box-shadow: none;
}

.sidebar {
  width: 280px;
  background: rgba(151, 150, 161, 0.4);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width 0.3s ease;
  position: relative;
}

.sidebar.collapsed {
  width: 80px;
}

.sidebar-header {
  padding: 2rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  position: relative;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.sidebar-header:hover {
  background: rgba(255, 255, 255, 0.2);
}

.user-role {
  background: rgba(45, 80, 22, 0.9);
  color: #b7f5b7;
  font-weight: bold;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  font-family: "Roboto", sans-serif;
  text-align: center;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  letter-spacing: 1px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  box-shadow: 0 4px 15px rgba(45, 80, 22, 0.3);
  animation: roleGlow 3s ease-in-out infinite alternate;
}

@keyframes roleGlow {
  0% {
    background: rgba(45, 80, 22, 0.9);
    box-shadow: 0 4px 15px rgba(45, 80, 22, 0.3);
  }
  100% {
    background: rgba(45, 80, 22, 1);
    box-shadow: 0 6px 20px rgba(45, 80, 22, 0.5);
  }
}

.logo {
  font-family: "Dancing Script", cursive;
  font-size: 1.4rem;
  color: white;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  text-align: center;
  display: block;
  font-weight: 700;
}

.sidebar-nav {
  flex: 1;
  padding: 1.5rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  overflow-y: auto;
}

.sidebar-nav button {
  width: 100%;
  padding: 1rem 1.25rem;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(183, 245, 183, 0.3);
  color: white;
  font-weight: 600;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
  font-size: 0.95rem;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.sidebar-nav button:hover {
  background: rgba(183, 245, 183, 0.2);
  border-color: rgba(183, 245, 183, 0.6);
  transform: translateX(5px);
  box-shadow: 0 4px 15px rgba(183, 245, 183, 0.3);
}

.sidebar-nav button.active {
  background: rgba(183, 245, 183, 0.3);
  border-color: rgba(183, 245, 183, 0.8);
  color: #2d5016;
  font-weight: 700;
  transform: translateX(5px);
  box-shadow: 0 6px 20px rgba(183, 245, 183, 0.4);
}

.nav-icon {
  font-size: 1.2rem;
  min-width: 20px;
  text-align: center;
}

.nav-text {
  transition: opacity 0.3s ease;
}

.sidebar.collapsed .nav-text {
  opacity: 0;
  width: 0;
  overflow: hidden;
}

.sidebar.collapsed .sidebar-nav button {
  justify-content: center;
  padding: 1rem 0.5rem;
}

.sidebar-footer {
  padding: 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  justify-content: center;
}

.logout-btn {
  width: 60px;
  height: 60px;
  background: rgba(220, 53, 69, 0.8);
  border: 2px solid rgba(255, 255, 255, 0.3);
  color: white;
  font-size: 1.8rem;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  box-shadow: 0 4px 15px rgba(220, 53, 69, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-btn:hover {
  background: rgba(220, 53, 69, 1);
  transform: scale(1.1) rotate(-10deg);
  box-shadow: 0 6px 25px rgba(220, 53, 69, 0.5);
  border-color: rgba(255, 255, 255, 0.6);
}

.logout-btn:active {
  transform: scale(0.95) rotate(-10deg);
}

.main-content {
  flex: 1;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  position: relative;
}

/* Scrollbar styling */
.sidebar-nav::-webkit-scrollbar,
.main-content::-webkit-scrollbar {
  width: 8px;
}

.sidebar-nav::-webkit-scrollbar-track,
.main-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.sidebar-nav::-webkit-scrollbar-thumb,
.main-content::-webkit-scrollbar-thumb {
  background: rgba(183, 245, 183, 0.5);
  border-radius: 4px;
}

.sidebar-nav::-webkit-scrollbar-thumb:hover,
.main-content::-webkit-scrollbar-thumb:hover {
  background: rgba(183, 245, 183, 0.7);
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .panel-bg {
    padding: 0;
  }

  .panel {
    height: 100vh;
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    height: auto;
    max-height: 200px;
    flex-direction: row;
    overflow-x: auto;
    overflow-y: hidden;
  }

  .sidebar.collapsed {
    width: 100%;
    max-height: 80px;
  }

  .sidebar-header {
    flex-shrink: 0;
    padding: 1rem;
    border-right: 1px solid rgba(255, 255, 255, 0.2);
    border-bottom: none;
  }

  .sidebar-nav {
    flex-direction: row;
    padding: 1rem;
    gap: 0.5rem;
    min-width: max-content;
  }

  .sidebar-nav button {
    white-space: nowrap;
    min-width: 120px;
    flex-direction: column;
    gap: 0.25rem;
  }

  .nav-icon {
    font-size: 1rem;
  }

  .nav-text {
    font-size: 0.8rem;
  }

  .sidebar-footer {
    flex-shrink: 0;
    padding: 1rem;
    border-left: 1px solid rgba(255, 255, 255, 0.2);
    border-top: none;
  }

  .logout-btn {
    width: 50px;
    height: 50px;
    font-size: 1.5rem;
  }

  .logo {
    font-size: 1.1rem;
  }

  .user-role {
    font-size: 0.8rem;
    padding: 0.5rem 0.75rem;
    margin-bottom: 0.5rem;
  }
}
</style>