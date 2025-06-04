<template>
  <div class="dashboard-container">
    <div class="trainee-dashboard">
      <div class="dashboard-header">
        <h2>💪 Witaj {{ userNickname }}!</h2>
        <p>Kontynuuj swoją podróż fitness i śledź postępy.</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Ładowanie danych...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <h3>❌ Wystąpił błąd</h3>
        <p>{{ error }}</p>
        <button @click="retryLoad" class="retry-btn">🔄 Spróbuj ponownie</button>
      </div>

      <!-- Trainer Search Section (shown when no trainer assigned) -->
      <div v-else-if="!hasApprovedTrainer" class="trainer-search-section">
        <div class="search-card">
          <h3>🔍 Znajdź swojego trenera</h3>
          <p>Wyszukaj i wyślij prośbę o współpracę do trenera</p>

          <div class="search-input-group">
            <input
                v-model="trainerSearchQuery"
                @input="debouncedSearch"
                @keyup.enter="searchTrainers"
                type="text"
                placeholder="Wpisz imię lub nazwisko trenera..."
                class="trainer-search-input"
                :disabled="searchLoading"
            />
            <button @click="searchTrainers" class="search-btn" :disabled="searchLoading || trainerSearchQuery.length < 2">
              {{ searchLoading ? '🔄' : '🔍' }} Szukaj
            </button>
          </div>

          <!-- Search Results -->
          <div v-if="searchResults.length > 0" class="search-results">
            <h4>Wyniki wyszukiwania:</h4>
            <div class="trainer-list">
              <div
                  v-for="trainer in searchResults"
                  :key="trainer.id"
                  class="trainer-result-card"
              >
                <div class="trainer-result-info">
                  <h5>{{ trainer.first_name + ' ' + trainer.last_name }}</h5>
                  <p>{{ trainer.specializations?.join(', ') || 'Brak specjalizacji' }}</p>
                  <span class="trainer-experience">{{ trainer.experience || 0 }} lat doświadczenia</span>
                  <div class="trainer-status" :class="{ online: trainer.is_online }">
                    {{ trainer.is_online ? '🟢 Online' : '⚫ Offline' }}
                  </div>
                </div>
                <button
                    @click="sendTrainerRequest(trainer)"
                    :disabled="hasPendingRequest(trainer.id) || requestLoading === trainer.id"
                    class="request-btn"
                    :class="{ 'pending': hasPendingRequest(trainer.id) }"
                >
                  <span v-if="requestLoading === trainer.id">🔄 Wysyłanie...</span>
                  <span v-else-if="hasPendingRequest(trainer.id)">⏳ Oczekuje...</span>
                  <span v-else>✉️ Wyślij prośbę</span>
                </button>
              </div>
            </div>
          </div>

          <!-- No Results -->
          <div v-else-if="searchAttempted && !searchLoading" class="no-results">
            <p>Nie znaleziono trenerów spełniających kryteria.</p>
          </div>

          <!-- Pending Requests -->
          <div v-if="pendingRequests.length > 0" class="pending-requests">
            <h4>Wysłane prośby:</h4>
            <div class="request-list">
              <div
                  v-for="request in pendingRequests"
                  :key="request.id"
                  class="request-item"
                  :class="{ ['status-' + request.status]: true }"
              >
                <div class="request-info">
                  <span class="trainer-name">{{ request.trainer_username }}</span>
                  <span class="request-status">
                    <span v-if="request.status === 'pending'">⏳ Oczekuje na odpowiedź</span>
                    <span v-else-if="request.status === 'approved'">✅ Zaakceptowana</span>
                    <span v-else-if="request.status === 'rejected'">❌ Odrzucona</span>
                    <span v-else>{{ request.status }}</span>
                  </span>
                  <span class="request-date">{{ formatDate(new Date(request.created_at)) }}</span>
                </div>
                <button
                    v-if="request.status === 'pending'"
                    @click="cancelRequest(request.id)"
                    class="cancel-btn"
                    :disabled="cancelLoading === request.id"
                >
                  {{ cancelLoading === request.id ? '🔄' : '❌' }} Anuluj
                </button>
              </div>
            </div>
          </div>

          <!-- Tips Section -->
          <div class="tips-section">
            <h4>💡 Wskazówki</h4>
            <ul>
              <li>Wyszukaj trenera po imieniu, nazwisku lub specjalizacji</li>
              <li>Sprawdź profil trenera przed wysłaniem prośby</li>
              <li>Możesz wysłać prośby do maksymalnie 3 trenerów jednocześnie</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Main Dashboard Content (shown only when trainer is approved) -->
      <div v-else class="dashboard-content">
        <!-- Trainer Profile Section -->
        <div class="trainer-profile-section">
          <div class="section-header">
            <h3>👨‍🏫 Profil Twojego Trenera</h3>
            <button @click="toggleProfileExpanded" class="expand-profile-btn">
              {{ showExpandedProfile ? '▲ Zwiń profil' : '▼ Rozwiń profil' }}
            </button>
          </div>

          <div class="trainer-profile-card" v-if="trainerProfile">
            <div class="trainer-profile-header">
              <div class="trainer-avatar" :style="{ backgroundColor: trainerAvatarColor }">
                <span class="trainer-initials">{{ trainerInitials }}</span>
              </div>

              <div class="trainer-basic-info">
                <h4>{{ trainerProfile.first_name }} {{ trainerProfile.last_name }}</h4>
                <div class="trainer-title">🏆 Twój Trener Personalny</div>
                <div class="trainer-experience-level">
                  <span class="experience-badge" :class="experienceLevel.class">
                    {{ experienceLevel.label }}
                  </span>
                  <span class="experience-years">{{ trainerProfile.experience }} {{ getYearsText(trainerProfile.experience) }}</span>
                </div>
                <div class="trainer-status" :class="{ online: approvedTrainer?.is_online }">
                  {{ approvedTrainer?.is_online ? '🟢 Online' : '⚫ Offline' }}
                </div>
              </div>

              <div class="trainer-contact-actions">
                <button
                    @click="openChatWithTrainer"
                    class="action-btn chat-btn"
                    :class="{ 'has-notifications': unreadMessages > 0 }"
                >
                  💬 Chat
                  <span v-if="unreadMessages > 0" class="notification-badge messages-badge">{{ unreadMessages }}</span>
                </button>
              </div>
            </div>

            <!-- Expanded Profile Details -->
            <div v-if="showExpandedProfile" class="expanded-profile-details">
              <div class="profile-details-grid">
                <!-- Basic Information -->
                <div class="detail-section">
                  <h4>📋 Podstawowe informacje</h4>
                  <div class="profile-details">
                    <div class="detail-item">
                      <span class="detail-label">💰 Cena za godzinę:</span>
                      <span class="detail-value">{{ trainerProfile.price_per_hour || 'Nie podano' }} zł</span>
                    </div>
                    <div class="detail-item">
                      <span class="detail-label">📍 Lokalizacja:</span>
                      <span class="detail-value">{{ trainerProfile.location || 'Nie podano' }}</span>
                    </div>
                    <div class="detail-item" v-if="trainerProfile.contact_email">
                      <span class="detail-label">📧 Email:</span>
                      <span class="detail-value">{{ trainerProfile.contact_email }}</span>
                    </div>
                    <div class="detail-item" v-if="trainerProfile.contact_phone">
                      <span class="detail-label">📱 Telefon:</span>
                      <span class="detail-value">{{ trainerProfile.contact_phone }}</span>
                    </div>
                  </div>
                </div>

                <!-- Specializations -->
                <div class="detail-section">
                  <h4>🎯 Specjalizacje</h4>
                  <div class="tags-container">
                    <span v-for="spec in trainerProfile.specializations" :key="spec" class="tag specialization-tag">
                      {{ spec }}
                    </span>
                    <span v-if="!trainerProfile.specializations || trainerProfile.specializations.length === 0" class="no-data">
                      Brak specjalizacji
                    </span>
                  </div>
                </div>

                <!-- Certifications -->
                <div class="detail-section">
                  <h4>📜 Certyfikaty</h4>
                  <div class="tags-container">
                    <span v-for="cert in trainerProfile.certifications" :key="cert" class="tag certification-tag">
                      {{ cert }}
                    </span>
                    <span v-if="!trainerProfile.certifications || trainerProfile.certifications.length === 0" class="no-data">
                      Brak certyfikatów
                    </span>
                  </div>
                </div>

                <!-- Languages -->
                <div class="detail-section">
                  <h4>🌍 Języki</h4>
                  <div class="tags-container">
                    <span v-for="lang in trainerProfile.languages" :key="lang" class="tag language-tag">
                      {{ lang }}
                    </span>
                    <span v-if="!trainerProfile.languages || trainerProfile.languages.length === 0" class="no-data">
                      Brak języków
                    </span>
                  </div>
                </div>
              </div>

              <!-- Bio Section -->
              <div v-if="trainerProfile.bio" class="bio-section">
                <h4>📝 O trenerze</h4>
                <div class="bio-text">
                  {{ trainerProfile.bio }}
                </div>
              </div>

              <!-- Stats Section -->
              <div class="trainer-stats-section">
                <h4>📊 Statystyki trenera</h4>
                <div class="calculations">
                  <div class="calc-item">
                    <span class="calc-label">🏋️ Lata doświadczenia:</span>
                    <span class="calc-value">{{ trainerProfile.experience || 0 }}</span>
                  </div>
                  <div class="calc-item">
                    <span class="calc-label">🎯 Liczba specjalizacji:</span>
                    <span class="calc-value">{{ trainerProfile.specializations?.length || 0 }}</span>
                  </div>
                  <div class="calc-item">
                    <span class="calc-label">📜 Liczba certyfikatów:</span>
                    <span class="calc-value">{{ trainerProfile.certifications?.length || 0 }}</span>
                  </div>
                  <div class="calc-item">
                    <span class="calc-label">📅 Trener od:</span>
                    <span class="calc-value">{{ formatProfileDate(trainerProfile.created_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Loading trainer profile -->
          <div v-else-if="trainerProfileLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <p>Ładowanie profilu trenera...</p>
          </div>

          <!-- No trainer profile -->
          <div v-else class="no-profile">
            <div class="trainer-basic-card">
              <div class="trainer-avatar" :style="{ backgroundColor: '#2d5016' }">
                <span class="trainer-initials">{{ approvedTrainer?.trainer_username?.substring(0, 2).toUpperCase() || 'TR' }}</span>
              </div>
              <div class="trainer-basic-info">
                <h4>{{ approvedTrainer?.trainer_username || 'Nieznany trener' }}</h4>
                <div class="trainer-title">🏆 Twój Trener</div>
                <div class="trainer-specialization">{{ approvedTrainer?.trainer_specialization || 'Trener personalny' }}</div>
                <div class="trainer-status" :class="{ online: approvedTrainer?.is_online }">
                  {{ approvedTrainer?.is_online ? '🟢 Online' : '⚫ Offline' }}
                </div>
              </div>
              <div class="trainer-contact-actions">
                <button @click="openChatWithTrainer" class="action-btn chat-btn" :class="{ 'has-notifications': unreadMessages > 0 }">
                  💬 Chat
                  <span v-if="unreadMessages > 0" class="notification-badge messages-badge">{{ unreadMessages }}</span>
                </button>
              </div>
            </div>
            <p>👤 Trener nie uzupełnił jeszcze swojego profilu lub profil nie został załadowany</p>
            <button @click="loadTrainerProfile" class="retry-btn">🔄 Spróbuj ponownie załadować profil</button>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="actions-section">
          <div class="section-header">
            <h3>🎯 Szybkie Akcje</h3>
          </div>
          <div class="trainee-actions">
            <button @click="checkTrainingPlan" class="action-btn plan-btn">
              🏋️‍♂️ Plan Treningowy
              <div class="action-status">Sprawdź swój plan</div>
            </button>

            <button @click="navigateToReports" class="action-btn reports-btn">
              📊 Prześlij Raport
              <div class="action-status">Wyślij postępy</div>
            </button>

            <button @click="checkSharedVideos" class="action-btn media-btn">
              🎥 Pliki Wideo
              <span v-if="newVideos > 0" class="notification-badge media-badge">{{ newVideos }}</span>
              <div class="action-status">Materiały od trenera</div>
            </button>

            <button @click="checkDiet" class="action-btn diet-btn">
              🥗 Twoja Dieta
              <div class="action-status">Sprawdź dietę</div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Success/Error Notifications -->
    <div v-if="systemNotification" class="system-notification" :class="systemNotification.type">
      <span>{{ systemNotification.message }}</span>
      <button @click="systemNotification = null" class="notification-close">×</button>
    </div>
  </div>
</template>

<script setup>
import {ref, onMounted, computed, onUnmounted, watch, inject} from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const switchToView = inject('switchToView', null)
const safeSwitchToView = (viewName) => {
  if (switchToView && typeof switchToView === 'function') {
    switchToView(viewName)
  } else {
    const routeMap = {
      'training-plans': 'training-plans',
      'user-media': '/user-media',
      'diet-plans': '/diet-plans',
      'messages': '/messages',
      'trainer-profile': '/trainer-profile',
      'reports': '/reports'
    }
    if (routeMap[viewName]) {
      router.push(routeMap[viewName])
    }
  }
}

const store = useAuthStore()
const router = useRouter()

// Loading and error states
const loading = ref(true)
const error = ref(null)
const searchLoading = ref(false)
const requestLoading = ref(null)
const cancelLoading = ref(null)
const searchAttempted = ref(false)
const trainerProfileLoading = ref(false)

// User data
const userNickname = ref('rcepce')
const hasApprovedTrainer = ref(false)
const approvedTrainer = ref(null)
const unreadMessages = ref(0)
const newVideos = ref(0)

// Trainer profile data
const trainerProfile = ref(null)
const showExpandedProfile = ref(false)
const trainerAvatarColor = ref('#2d5016')

// Debug mode
const debugMode = ref(false)

// Trainer search functionality
const trainerSearchQuery = ref('')
const searchResults = ref([])
const pendingRequests = ref([])

// System notification
const systemNotification = ref(null)

// Computed properties
const trainerInitials = computed(() => {
  if (trainerProfile.value && trainerProfile.value.first_name && trainerProfile.value.last_name) {
    return `${trainerProfile.value.first_name.charAt(0)}${trainerProfile.value.last_name.charAt(0)}`.toUpperCase()
  }
  if (approvedTrainer.value?.trainer_username) {
    return approvedTrainer.value.trainer_username.substring(0, 2).toUpperCase()
  }
  return 'TR'
})

const experienceLevel = computed(() => {
  const exp = trainerProfile.value?.experience || 0
  if (exp < 1) return { label: 'Początkujący', class: 'beginner' }
  if (exp < 3) return { label: 'Junior', class: 'junior' }
  if (exp < 5) return { label: 'Doświadczony', class: 'experienced' }
  if (exp < 10) return { label: 'Senior', class: 'senior' }
  return { label: 'Ekspert', class: 'expert' }
})

// API base URL
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://188.245.239.124:8080/api'

// API utility function
const apiRequest = async (url, options = {}) => {
  const token = store.token
  console.log(`🌐 API Request: ${API_BASE_URL}${url}`)

  const response = await fetch(`${API_BASE_URL}${url}`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : '',
      ...options.headers
    },
    ...options
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Network error' }))
    throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
  }

  const data = await response.json()
  console.log(`📥 API Response for ${url}:`, data)
  return data
}

// Notification utility
const showSystemNotification = (message, type = 'success') => {
  systemNotification.value = { message, type }
  setTimeout(() => {
    systemNotification.value = null
  }, 5000)
}

// Generate random avatar color
const generateAvatarColor = () => {
  const colors = ['#2d5016', '#1976d2', '#d32f2f', '#f57c00', '#7b1fa2', '#00695c']
  trainerAvatarColor.value = colors[Math.floor(Math.random() * colors.length)]
}

// Debounced search
let searchTimeout = null
const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    if (trainerSearchQuery.value.length >= 2) {
      searchTrainers()
    }
  }, 500)
}

// Trainer search methods
const searchTrainers = async () => {
  if (trainerSearchQuery.value.length < 2) {
    searchResults.value = []
    searchAttempted.value = false
    return
  }

  searchLoading.value = true
  searchAttempted.value = true

  try {
    const data = await apiRequest(`/trainers/search?q=${encodeURIComponent(trainerSearchQuery.value)}`)
    searchResults.value = data || []
  } catch (err) {
    console.error('🚨 Search trainers error:', err)
    searchResults.value = []
    showSystemNotification('Błąd podczas wyszukiwania trenerów: ' + err.message, 'error')
  } finally {
    searchLoading.value = false
  }
}

const sendTrainerRequest = async (selectedTrainer) => {
  if (pendingRequests.value.length >= 3) {
    showSystemNotification('Możesz wysłać maksymalnie 3 prośby jednocześnie', 'error')
    return
  }

  requestLoading.value = selectedTrainer.id

  try {
    const spec = selectedTrainer.specialization ||
        (selectedTrainer.specializations && selectedTrainer.specializations.length > 0
            ? selectedTrainer.specializations.join(', ')
            : 'nie podano');
    const displayName = selectedTrainer.username ||
        ((selectedTrainer.first_name || '') + ' ' + (selectedTrainer.last_name || '')).trim();

    const requestData = {
      trainer_id: selectedTrainer.id,
      message: `Witaj ${displayName}! Chciałbym/chciałabym rozpocząć współpracę z Tobą jako moim trenerem personalnym!!!`
    }

    await apiRequest('/trainer-requests', {
      method: 'POST',
      body: JSON.stringify(requestData)
    })

    await loadPendingRequests()
    showSystemNotification(`Prośba do ${displayName} została wysłana!`, 'success')

  } catch (err) {
    console.error('🚨 Send trainer request error:', err)
    showSystemNotification('Błąd podczas wysyłania prośby: ' + err.message, 'error')
  } finally {
    requestLoading.value = null
  }
}

const cancelRequest = async (requestId) => {
  cancelLoading.value = requestId

  try {
    await apiRequest(`/trainer-requests/${requestId}`, {
      method: 'DELETE'
    })

    pendingRequests.value = pendingRequests.value.filter(req => req.id !== requestId)
    showSystemNotification('Prośba została anulowana', 'success')

  } catch (err) {
    console.error('🚨 Cancel request error:', err)
    showSystemNotification('Błąd podczas anulowania prośby: ' + err.message, 'error')
  } finally {
    cancelLoading.value = null
  }
}

const hasPendingRequest = (trainerId) => {
  return pendingRequests.value.some(req => req.trainer_id === trainerId && req.status === 'pending')
}

// Data loading functions
const loadPendingRequests = async () => {
  try {
    const data = await apiRequest('/trainer-requests/my')
    pendingRequests.value = data || []
  } catch (err) {
    console.error('🚨 Load pending requests error:', err)
  }
}

const loadApprovedTrainer = async () => {
  try {
    const data = await apiRequest('/trainer-requests/my-trainer')
    console.log('✅ Approved trainer data:', data)

    if (data) {
      approvedTrainer.value = data
      hasApprovedTrainer.value = true
    } else {
      hasApprovedTrainer.value = false
      approvedTrainer.value = null
    }
  } catch (err) {
    console.error('🚨 Load approved trainer error:', err)
    hasApprovedTrainer.value = false
    approvedTrainer.value = null
  }
}

// Load trainer profile - NAPRAWIONE
const loadTrainerProfile = async () => {
  // ZMIANA: użyj 'id' zamiast 'trainer_id'
  if (!approvedTrainer.value?.id) {
    console.log('❌ No trainer id available for profile loading')
    trainerProfile.value = null
    trainerProfileLoading.value = false
    return
  }

  trainerProfileLoading.value = true
  console.log(`🔍 Loading trainer profile for trainer_id: ${approvedTrainer.value.id}`)

  try {
    // ZMIANA: użyj 'id' zamiast 'trainer_id'
    const profileData = await apiRequest(`/trainer-profile/${approvedTrainer.value.id}`)
    console.log('📊 Profile response:', profileData)

    if (profileData && profileData.trainer) {
      trainerProfile.value = profileData.trainer
    } else if (profileData && profileData.profile) {
      trainerProfile.value = profileData.profile
    } else {
      console.log('❌ No trainer profile data found')
      trainerProfile.value = null
    }

    if (trainerProfile.value) {
      console.log('✅ Trainer profile loaded successfully:', trainerProfile.value)
      generateAvatarColor()
    }

  } catch (err) {
    console.error('🚨 Load trainer profile error:', err)
    trainerProfile.value = null
  } finally {
    trainerProfileLoading.value = false
  }
}

const loadUnreadMessages = async () => {
  try {
    const data = await apiRequest('/messages/unread')
    unreadMessages.value = data.count || 0
  } catch (err) {
    console.error('🚨 Load unread messages error:', err)
  }
}

const loadSharedMedia = async () => {
  try {
    const data = await apiRequest('/media/shared')
    if (data) {
      newVideos.value = data.unread_count || 0
    }
  } catch (err) {
    console.error('🚨 Load shared media error:', err)
  }
}

// Button functionality implementations
const openChatWithTrainer = () => {
  if (approvedTrainer.value) {
    safeSwitchToView('messages')
  } else {
    showSystemNotification('Nie masz przypisanego trenera', 'error')
  }
}

const toggleProfileExpanded = () => {
  showExpandedProfile.value = !showExpandedProfile.value
}

const checkTrainingPlan = () => {
  safeSwitchToView('training-plans')
}

const navigateToReports = () => {
  safeSwitchToView('reports')
}

const checkSharedVideos = () => {
  safeSwitchToView('user-media')
}

const checkDiet = () => {
  safeSwitchToView('diet-plans')
}

// Utility functions
const formatDate = (date) => {
  const now = new Date()
  const diff = now - date
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)

  if (days > 0) return `${days} dni temu`
  if (hours > 0) return `${hours} godz. temu`
  return 'Przed chwilą'
}

const formatProfileDate = (dateString) => {
  if (!dateString) return 'Niedawno'
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('pl-PL', { year: 'numeric', month: 'long' })
  } catch {
    return 'Niedawno'
  }
}

const getYearsText = (years) => {
  if (years === 1) return 'rok'
  if (years < 5) return 'lata'
  return 'lat'
}

const retryLoad = async () => {
  error.value = null
  await loadDashboardData()
}

const loadDashboardData = async () => {
  loading.value = true
  error.value = null

  try {
    // Load user data from auth store
    if (store.user) {
      userNickname.value = store.user.nickname || store.user.username || 'rcepce'
    }

    // Load all dashboard data in parallel where possible
    await Promise.allSettled([
      loadApprovedTrainer(),
      loadPendingRequests()
    ])

    // If we have an approved trainer, load more data
    if (hasApprovedTrainer.value) {
      await Promise.allSettled([
        loadTrainerProfile(),
        loadUnreadMessages(),
        loadSharedMedia()
      ])
    }

  } catch (err) {
    console.error('🚨 Load dashboard data error:', err)
    error.value = 'Nie udało się załadować danych. Sprawdź połączenie internetowe.'
  } finally {
    loading.value = false
  }
}

// Auto-refresh functionality
let refreshInterval = null

const startAutoRefresh = () => {
  // Refresh every 5 minutes
  refreshInterval = setInterval(async () => {
    if (hasApprovedTrainer.value) {
      await Promise.allSettled([
        loadUnreadMessages(),
        loadSharedMedia()
      ])
    }
  }, 5 * 60 * 1000)
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

// Lifecycle hooks
onMounted(async () => {
  console.log('🚀 TraineeDashboard mounted')
  await loadDashboardData()
  startAutoRefresh()
  if (localStorage.getItem('messagesCleared') === 'true') {
    unreadMessages.value = 0
    localStorage.removeItem('messagesCleared')
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})

// Watch for approved trainer changes
watch(hasApprovedTrainer, async (newValue) => {
  if (newValue) {
    console.log('👀 Watching: approved trainer changed, loading profile...')
    await loadTrainerProfile()
  }
})
</script>
<style scoped>
  /* Stylizacja zgodna z panelem trenera */
.dashboard-container {
  padding: 2rem;
  min-height: 100vh;
}

.trainee-dashboard {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  width: 95%;
  margin: 0 auto;
  color: #181818;
  overflow-y: auto;
  position: relative;
}

.trainee-dashboard::-webkit-scrollbar {
  width: 12px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 2rem;
}

.dashboard-header h2 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
}

.dashboard-header p {
  color: #666;
  margin: 0;
}

/* Loading States */
.loading-state, .error-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #2d5016;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.retry-btn {
  padding: 0.75rem 1.5rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.retry-btn:hover {
  background: #c8e6c9;
}

/* Section Headers */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h3 {
  color: #2d5016;
  margin: 0;
}

/* Trainer Profile Section */
.trainer-profile-section {
  margin-bottom: 2rem;
}

.trainer-profile-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  transition: 0.2s;
}

.trainer-profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.trainer-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  font-weight: bold;
  color: white;
  flex-shrink: 0;
}

.trainer-basic-info {
  flex: 1;
}

.trainer-basic-info h4 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
  font-size: 1.3rem;
  font-weight: 700;
}

.trainer-title {
  color: #666;
  margin-bottom: 0.5rem;
}

.trainer-experience-level {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.experience-badge {
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: bold;
}

.experience-badge.beginner {
  background: #e3f2fd;
  color: #1976d2;
}

.experience-badge.junior {
  background: #e8f5e8;
  color: #2e7d32;
}

.experience-badge.experienced {
  background: #fff3e0;
  color: #f57c00;
}

.experience-badge.senior {
  background: #f3e5f5;
  color: #7b1fa2;
}

.experience-badge.expert {
  background: #ffebee;
  color: #d32f2f;
}

.trainer-status {
  font-size: 0.8rem;
  color: #666;
}

.trainer-status.online {
  color: #2e7d32;
  font-weight: bold;
}

.trainer-contact-actions {
  flex-shrink: 0;
}

.expand-profile-btn, .refresh-btn {
  padding: 0.5rem 1rem;
  background: #f8f9fa;
  color: #2d5016;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;
}

.expand-profile-btn:hover, .refresh-btn:hover {
  background: #e9ecef;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Expanded Profile Details */
.expanded-profile-details {
  border-top: 1px solid #e9ecef;
  padding-top: 1.5rem;
  margin-top: 1rem;
}

.profile-details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.detail-section {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.detail-section h4 {
  color: #2d5016;
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e8f5e8;
}

.profile-details {
  display: grid;
  gap: 0.75rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.detail-label {
  font-weight: 600;
  color: #333;
}

.detail-value {
  color: #666;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.3rem 0.6rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.specialization-tag {
  background: #e8f5e8;
  color: #2e7d32;
}

.certification-tag {
  background: #e3f2fd;
  color: #1976d2;
}

.language-tag {
  background: #f3e5f5;
  color: #7b1fa2;
}

.no-data {
  color: #999;
  font-style: italic;
}

.bio-section {
  margin-bottom: 1.5rem;
}

.bio-text {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  line-height: 1.6;
  color: #333;
}

.trainer-stats-section h4 {
  color: #2d5016;
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e8f5e8;
}

.calculations {
  display: grid;
  gap: 0.75rem;
}

.calc-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #e8f5e8;
  border-radius: 8px;
  border: 1px solid #c8e6c9;
}

.calc-label {
  font-weight: 600;
  color: #2d5016;
}

.calc-value {
  color: #2d5016;
  font-weight: bold;
}

.no-profile {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.trainer-basic-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
}

/* Actions Section */
.actions-section {
  margin-bottom: 2rem;
}

.trainee-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-btn {
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  transition: 0.2s;
  text-align: center;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  min-height: 80px;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.notification-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  min-width: 20px;
  height: 20px;
  border-radius: 50%;
  font-size: 0.7rem;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: 2px solid white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.messages-badge {
  background: linear-gradient(135deg, #FFD700, #FFA500);
}

.media-badge {
  background: linear-gradient(135deg, #FF69B4, #E91E63);
}

.chat-btn {
  background: #fff3e0;
  color: #f57c00;
}

.chat-btn:hover {
  background: #ffe0b2;
}

.reports-btn {
  background: #f3e5f5;
  color: #7b1fa2;
}

.reports-btn:hover {
  background: #e1bee7;
}

.media-btn {
  background: #fce4ec;
  color: #c2185b;
}

.media-btn:hover {
  background: #f8bbd9;
}

.plan-btn {
  background: #e0f2f1;
  color: #00695c;
}

.plan-btn:hover {
  background: #b2dfdb;
}

.diet-btn {
  background: #e8f5e8;
  color: #2e7d32;
}

.diet-btn:hover {
  background: #c8e6c9;
}

.action-status {
  font-size: 0.8rem;
  color: #666;
  text-align: center;
}

/* Notifications Section */
.notifications-section {
  margin-bottom: 2rem;
}

.notifications-list {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1.5rem;
}

.notifications-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.notification-item:hover {
  border-color: #2d5016;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.notification-item.unread {
  border-left: 4px solid #2d5016;
  background: #f0f8f0;
}

.notification-item.important {
  border-left: 4px solid #d32f2f;
}

.notification-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
  margin-top: 0.2rem;
}

.notification-content {
  flex: 1;
}

.notification-title {
  font-weight: 600;
  color: #2d5016;
  margin-bottom: 0.25rem;
}

.notification-message {
  color: #333;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.notification-time {
  font-size: 0.8rem;
  color: #666;
}

.notification-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 8px;
  height: 8px;
  background: #2d5016;
  border-radius: 50%;
}

.empty-notifications {
  text-align: center;
  padding: 3rem 1rem;
  color: #666;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

/* Trainer Search Styles */
.trainer-search-section {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 2rem;
}

.search-card h3 {
  color: #2d5016;
  margin-bottom: 1rem;
}

.search-input-group {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.trainer-search-input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

.search-btn {
  padding: 0.75rem 1.5rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.search-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.search-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.trainer-result-card {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.request-btn, .cancel-btn {
  padding: 0.5rem 1rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.request-btn:hover:not(:disabled), .cancel-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.request-btn:disabled, .cancel-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #2d5016;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

.close-btn:hover {
  color: #333;
}

.report-form {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2d5016;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #2d5016;
}

.form-group small {
  display: block;
  margin-top: 0.5rem;
  color: #666;
  font-size: 0.8rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.cancel-btn,
.submit-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
  border: none;
}

.cancel-btn {
  background: #f8f9fa;
  color: #666;
}

.cancel-btn:hover {
  background: #e9ecef;
}

.submit-btn {
  background: #e8f5e8;
  color: #2e7d32;
}

.submit-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* System Notification */
.system-notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  z-index: 1000;
  display: flex;
  align-items: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  animation: slideIn 0.3s ease;
}

.system-notification.success {
  background: #43a047;
}

.system-notification.error {
  background: #d32f2f;
}

.notification-close {
  background: none;
  border: none;
  color: white;
  margin-left: 12px;
  cursor: pointer;
  font-size: 1.2rem;
  font-weight: bold;
}

@keyframes slideIn {
  from {
    transform: translateY(100px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }

  .trainee-dashboard {
    padding: 1rem;
  }

  .section-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .trainee-actions {
    grid-template-columns: 1fr 1fr;
  }

  .trainer-profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}
/* Stylizacja zgodna z panelem trenera */
.dashboard-container {
  padding: 2rem;
  min-height: 100vh;
}

.trainee-dashboard {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  width: 95%;
  margin: 0 auto;
  color: #181818;
  overflow-y: auto;
  position: relative;
}

.trainee-dashboard::-webkit-scrollbar {
  width: 12px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 2rem;
}

.dashboard-header h2 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
}

.dashboard-header p {
  color: #666;
  margin: 0;
}

/* Loading States */
.loading-state, .error-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #2d5016;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.retry-btn {
  padding: 0.75rem 1.5rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.retry-btn:hover {
  background: #c8e6c9;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  text-align: center;
  border: 1px solid #e9ecef;
}

.stat-number {
  display: block;
  font-size: 2rem;
  font-weight: bold;
  color: #2d5016;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

/* Section Headers */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h3 {
  color: #2d5016;
  margin: 0;
}

/* Trainer Profile Section */
.trainer-profile-section {
  margin-bottom: 2rem;
}

.trainer-profile-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  transition: 0.2s;
}

.trainer-profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.trainer-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  font-weight: bold;
  color: white;
  flex-shrink: 0;
}

.trainer-basic-info {
  flex: 1;
}

.trainer-basic-info h4 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
  font-size: 1.3rem;
  font-weight: 700;
}

.trainer-title {
  color: #666;
  margin-bottom: 0.5rem;
}

.trainer-experience-level {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.experience-badge {
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: bold;
}

.experience-badge.beginner {
  background: #e3f2fd;
  color: #1976d2;
}

.experience-badge.junior {
  background: #e8f5e8;
  color: #2e7d32;
}

.experience-badge.experienced {
  background: #fff3e0;
  color: #f57c00;
}

.experience-badge.senior {
  background: #f3e5f5;
  color: #7b1fa2;
}

.experience-badge.expert {
  background: #ffebee;
  color: #d32f2f;
}

.trainer-status {
  font-size: 0.8rem;
  color: #666;
}

.trainer-status.online {
  color: #2e7d32;
  font-weight: bold;
}

.trainer-contact-actions {
  flex-shrink: 0;
}

.expand-profile-btn {
  padding: 0.5rem 1rem;
  background: #f8f9fa;
  color: #2d5016;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;
}

.expand-profile-btn:hover {
  background: #e9ecef;
}

/* Expanded Profile Details */
.expanded-profile-details {
  border-top: 1px solid #e9ecef;
  padding-top: 1.5rem;
  margin-top: 1rem;
}

.profile-details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.detail-section {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.detail-section h4 {
  color: #2d5016;
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e8f5e8;
}

.profile-details {
  display: grid;
  gap: 0.75rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.detail-label {
  font-weight: 600;
  color: #333;
}

.detail-value {
  color: #666;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.3rem 0.6rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.specialization-tag {
  background: #e8f5e8;
  color: #2e7d32;
}

.certification-tag {
  background: #e3f2fd;
  color: #1976d2;
}

.language-tag {
  background: #f3e5f5;
  color: #7b1fa2;
}

.no-data {
  color: #999;
  font-style: italic;
}

.bio-section {
  margin-bottom: 1.5rem;
}

.bio-text {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  line-height: 1.6;
  color: #333;
}

.trainer-stats-section h4 {
  color: #2d5016;
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e8f5e8;
}

.calculations {
  display: grid;
  gap: 0.75rem;
}

.calc-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #e8f5e8;
  border-radius: 8px;
  border: 1px solid #c8e6c9;
}

.calc-label {
  font-weight: 600;
  color: #2d5016;
}

.calc-value {
  color: #2d5016;
  font-weight: bold;
}

.no-profile {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.trainer-basic-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
}

/* Actions Section */
.actions-section {
  margin-bottom: 2rem;
}

.trainee-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-btn {
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  transition: 0.2s;
  text-align: center;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  min-height: 80px;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.notification-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  min-width: 20px;
  height: 20px;
  border-radius: 50%;
  font-size: 0.7rem;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: 2px solid white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.messages-badge {
  background: linear-gradient(135deg, #FFD700, #FFA500);
}

.reports-badge {
  background: linear-gradient(135deg, #4169E1, #1E3A8A);
}

.media-badge {
  background: linear-gradient(135deg, #FF69B4, #E91E63);
}

.diet-badge {
  background: linear-gradient(135deg, #32CD32, #228B22);
}

.chat-btn {
  background: #fff3e0;
  color: #f57c00;
}

.chat-btn:hover {
  background: #ffe0b2;
}

.reports-btn {
  background: #f3e5f5;
  color: #7b1fa2;
}

.reports-btn:hover {
  background: #e1bee7;
}

.media-btn {
  background: #fce4ec;
  color: #c2185b;
}

.media-btn:hover {
  background: #f8bbd9;
}

.plan-btn {
  background: #e0f2f1;
  color: #00695c;
}

.plan-btn:hover {
  background: #b2dfdb;
}

.diet-btn {
  background: #e8f5e8;
  color: #2e7d32;
}

.diet-btn:hover {
  background: #c8e6c9;
}

.action-status {
  font-size: 0.8rem;
  color: #666;
  text-align: center;
}

/* Progress Section */
.progress-section {
  margin-bottom: 2rem;
}

.progress-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.progress-card {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1.5rem;
}

.progress-card h4 {
  color: #2d5016;
  margin: 0 0 1rem 0;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.progress-label {
  color: #666;
  font-weight: 500;
}

.progress-percentage {
  color: #2d5016;
  font-weight: 600;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(135deg, #2e7d32, #43a047);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-time {
  color: #666;
  font-size: 0.9rem;
}

.workout-count {
  color: #2d5016;
  font-weight: 600;
  margin-bottom: 1rem;
}

.workout-dots {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.workout-dot {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 500;
  border: 1px solid #e9ecef;
  background: white;
  color: #666;
}

.workout-dot.completed {
  background: #2e7d32;
  color: white;
  border-color: #2e7d32;
}

.workout-dot.planned {
  background: #fff3e0;
  color: #f57c00;
  border-color: #f57c00;
}

.workout-dot.today {
  border-color: #2d5016;
  border-width: 2px;
}

.weight-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.current-weight {
  color: #2d5016;
  font-weight: 600;
  font-size: 1.2rem;
}

.weight-change {
  font-size: 0.9rem;
  font-weight: 500;
}

.weight-change.weight-up {
  color: #d32f2f;
}

.weight-change.weight-down {
  color: #2e7d32;
}

.weight-target {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 0.3rem;
}

.weight-remaining {
  color: #666;
  font-size: 0.9rem;
}

/* Activity Section */
.activity-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  margin-bottom: 2rem;
}

.activity-section h3 {
  color: #2d5016;
  margin: 0 0 1rem 0;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: white;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.activity-icon {
  font-size: 1.2rem;
}

.activity-text {
  color: #333;
  font-weight: 500;
  flex: 1;
}

.activity-time {
  color: #666;
  font-size: 0.8rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

/* Motivation Section */
.motivation-section {
  margin-bottom: 2rem;
}

.quote-card {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.quote-icon {
  font-size: 2rem;
  color: #2d5016;
}

.quote-content {
  flex: 1;
}

.quote-text {
  color: #2d5016;
  font-style: italic;
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
}

.quote-author {
  color: #666;
  font-weight: 500;
  font-size: 0.9rem;
}

.refresh-quote-btn {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.refresh-quote-btn:hover {
  background: #e9ecef;
  transform: rotate(180deg);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  color: #2d5016;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

.close-btn:hover {
  color: #333;
}

.report-form {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2d5016;
  font-weight: 500;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #2d5016;
}

.form-group small {
  display: block;
  margin-top: 0.5rem;
  color: #666;
  font-size: 0.8rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.cancel-btn,
.submit-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
  border: none;
}

.cancel-btn {
  background: #f8f9fa;
  color: #666;
}

.cancel-btn:hover {
  background: #e9ecef;
}

.submit-btn {
  background: #e8f5e8;
  color: #2e7d32;
}

.submit-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Notification */
.notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  z-index: 1000;
  display: flex;
  align-items: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  animation: slideIn 0.3s ease;
}

.notification.success {
  background: #43a047;
}

.notification.error {
  background: #d32f2f;
}

.notification-close {
  background: none;
  border: none;
  color: white;
  margin-left: 12px;
  cursor: pointer;
  font-size: 1.2rem;
  font-weight: bold;
}

@keyframes slideIn {
  from {
    transform: translateY(100px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* Trainer Search Styles */
.trainer-search-section {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 2rem;
}

.search-card h3 {
  color: #2d5016;
  margin-bottom: 1rem;
}

.search-input-group {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.trainer-search-input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

.search-btn {
  padding: 0.75rem 1.5rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.search-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.search-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.trainer-result-card {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.request-btn {
  padding: 0.5rem 1rem;
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.request-btn:hover:not(:disabled) {
  background: #c8e6c9;
}

.request-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Responsive Design */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }

  .trainee-dashboard {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: 1fr 1fr;
  }

  .section-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .trainee-actions {
    grid-template-columns: 1fr 1fr;
  }

  .trainer-profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .profile-details-grid {
    grid-template-columns: 1fr;
  }

  .progress-cards {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .trainee-actions {
    grid-template-columns: 1fr;
  }

  .trainer-avatar {
    width: 60px;
    height: 60px;
    font-size: 1.4rem;
  }
}

</style>
