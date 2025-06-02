<template>
  <div class="trainer-dashboard">
    <div class="dashboard-header">
      <h2>👨‍🏫 Witaj Trenerze!</h2>
      <p>Zarządzaj swoimi podopiecznymi i monitoruj ich postępy.</p>
    </div>

    <!-- Quick Stats -->
    <div class="stats-grid">
      <div class="stat-card">
        <span class="stat-number">{{ trainees.length }}</span>
        <span class="stat-label">👥 Podopieczni</span>
      </div>
      <div class="stat-card">
        <span class="stat-number">{{ activeTrainingPlans }}</span>
        <span class="stat-label">🏋️ Aktywne plany</span>
      </div>
      <div class="stat-card">
        <span class="stat-number">{{ pendingReports }}</span>
        <span class="stat-label">📊 Nowe raporty</span>
      </div>
      <div class="stat-card">
        <span class="stat-number">{{ unreadMessages }}</span>
        <span class="stat-label">💬 Wiadomości</span>
      </div>
    </div>

    <!-- Pending Requests -->
    <div v-if="pendingRequests.length > 0" class="requests-section">
      <div class="section-header">
        <h3>⏳ Oczekujące prośby o współpracę</h3>
      </div>
      <div class="requests-list">
        <div v-for="request in pendingRequests" :key="request.id" class="request-card">
          <div class="request-info">
            <h4>{{ request.client_username || 'Użytkownik' }} (#{{ request.client_id }})</h4>
            <p class="request-message">{{ request.message || 'Brak wiadomości' }}</p>
            <div class="request-date">Wysłano: {{ formatDate(new Date(request.created_at)) }}</div>
          </div>
          <div class="request-actions">
            <button @click="approveRequest(request.id)" class="approve-btn" :disabled="processingRequest === request.id">
              {{ processingRequest === request.id && actionType === 'approve' ? '⏳' : '✅' }} Akceptuj
            </button>
            <button @click="rejectRequest(request.id)" class="reject-btn" :disabled="processingRequest === request.id">
              {{ processingRequest === request.id && actionType === 'reject' ? '⏳' : '❌' }} Odrzuć
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Trainees List -->
    <div class="trainees-section">
      <div class="section-header">
        <h3>🎯 Twoi Podopieczni</h3>
      </div>

      <div v-if="loading" class="loading-state">
        <p>Ładowanie podopiecznych...</p>
      </div>

      <div v-else-if="trainees.length === 0" class="empty-state">
        <h3>Brak podopiecznych</h3>
        <p>Poczekaj na prośbę od klienta, aby rozpocząć pracę</p>
      </div>

      <div v-else class="trainees-list">
        <div v-for="trainee in trainees" :key="trainee.id" class="trainee-card">
          <div class="trainee-header">
            <div class="trainee-info">
              <h4>{{ trainee.username }}</h4>
              <div class="trainee-status" :class="trainee.is_online ? 'active' : 'inactive'">
                {{ trainee.is_online ? 'Aktywny' : 'Nieaktywny' }}
              </div>
              <div class="last-activity">
                Ostatnia aktywność: {{ formatDate(new Date(trainee.last_activity || Date.now())) }}
              </div>
            </div>
          </div>

          <!-- Action Buttons WITH NOTIFICATIONS -->
          <div class="trainee-actions">
            <button @click="viewTraineeProfile(trainee)" class="action-btn profile-btn">
              👤 Profil
            </button>

            <button @click="handleNotificationClick('chat', trainee)" class="action-btn chat-btn">
              💬 Chat
              <div v-if="trainee.unread_messages > 0" class="notification-badge messages-badge">
                {{ trainee.unread_messages }}
              </div>
            </button>

            <button @click="handleNotificationClick('reports', trainee)" class="action-btn reports-btn">
              📊 Raporty
              <div v-if="trainee.unread_reports > 0" class="notification-badge reports-badge">
                {{ trainee.unread_reports }}
              </div>
            </button>

            <button @click="handleNotificationClick('media', trainee)" class="action-btn media-btn">
              🎥 Media
              <div v-if="trainee.unread_media > 0" class="notification-badge media-badge">
                {{ trainee.unread_media }}
              </div>
            </button>

            <button @click="createTrainingPlan(trainee)" class="action-btn plan-btn">
              🏋️ Plan
            </button>

            <button @click="handleNotificationClick('diet-plans', trainee)" class="action-btn diet-btn">
              🍎 Dieta
              <div v-if="trainee.pending_diet_feedback > 0" class="notification-badge diet-badge">
                {{ trainee.pending_diet_feedback }}
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="activity-section" v-if="recentActivity.length > 0">
      <h3>📈 Ostatnia Aktywność</h3>
      <div class="activity-list">
        <div v-for="activity in recentActivity" :key="activity.id" class="activity-item">
          <span class="activity-icon">{{ activity.icon }}</span>
          <div class="activity-content">
            <div class="activity-text">{{ activity.text }}</div>
            <div class="activity-time">{{ formatDate(new Date(activity.timestamp)) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- User Profile Modal -->
    <div v-if="showProfileModal" class="modal-overlay" @click="closeProfileModal">
      <div class="profile-modal" @click.stop>
        <div class="modal-header">
          <h3>👤 Profil: {{ selectedTrainee?.username }}</h3>
          <button @click="closeProfileModal" class="close-btn">×</button>
        </div>
        <div class="modal-content">
          <div v-if="loadingProfile" class="loading-state">
            <p>Ładowanie profilu...</p>
          </div>
          <div v-else-if="traineeProfile" class="profile-content">
            <!-- Podstawowe dane -->
            <div class="profile-section">
              <h4>📊 Dane podstawowe</h4>
              <div class="profile-details">
                <div class="detail-item">
                  <span class="detail-label">🎂 Wiek:</span>
                  <span class="detail-value">{{ traineeProfile.age }} lat</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">📏 Wzrost:</span>
                  <span class="detail-value">{{ traineeProfile.height }} cm</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">⚖️ Waga:</span>
                  <span class="detail-value">{{ traineeProfile.weight }} kg</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">{{ traineeProfile.gender === 'male' ? '👨' : '👩' }} Płeć:</span>
                  <span class="detail-value">{{ getGenderText(traineeProfile.gender) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">🏃 Aktywność:</span>
                  <span class="detail-value">{{ getActivityText(traineeProfile.activity_level) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">🎯 Cel:</span>
                  <span class="detail-value">{{ getGoalText(traineeProfile.goal) }}</span>
                </div>
              </div>
            </div>

            <!-- BMI i kategoria -->
            <div class="profile-section">
              <h4>📈 BMI i Analiza Wagi</h4>
              <div class="bmi-section">
                <div class="bmi-display">
                  <div class="bmi-value" :class="getBMICategory(traineeProfile).class">
                    {{ calculateBMI(traineeProfile) }}
                  </div>
                  <div class="bmi-category">{{ getBMICategory(traineeProfile).label }}</div>
                </div>
                <div class="bmi-bar">
                  <div class="bmi-indicator" :style="{ left: getBMIPosition(traineeProfile) + '%' }"></div>
                  <div class="bmi-ranges">
                    <span class="range underweight">Niedowaga</span>
                    <span class="range normal">Norma</span>
                    <span class="range overweight">Nadwaga</span>
                    <span class="range obese">Otyłość</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Kalorie i metabolizm -->
            <div class="profile-section">
              <h4>🔥 Metabolizm i Kalorie</h4>
              <div class="calculations">
                <div class="calc-item">
                  <span class="calc-label">BMR (metabolizm podstawowy):</span>
                  <span class="calc-value">{{ calculateBMR(traineeProfile) }} kcal</span>
                </div>
                <div class="calc-item">
                  <span class="calc-label">Zapotrzebowanie kaloryczne:</span>
                  <span class="calc-value">{{ calculateDailyCalories(traineeProfile) }} kcal</span>
                </div>
              </div>
            </div>

            <!-- Skład ciała -->
            <div class="profile-section">
              <h4>🧬 Przewidywany Skład Ciała</h4>
              <div class="body-composition">
                <div class="composition-item">
                  <div class="comp-icon">💧</div>
                  <div class="comp-info">
                    <div class="comp-label">Woda</div>
                    <div class="comp-value">{{ calculateWaterPercentage(traineeProfile) }}%</div>
                    <div class="comp-bar">
                      <div class="comp-fill water" :style="{ width: calculateWaterPercentage(traineeProfile) + '%' }"></div>
                    </div>
                  </div>
                </div>

                <div class="composition-item">
                  <div class="comp-icon">💪</div>
                  <div class="comp-info">
                    <div class="comp-label">Mięśnie</div>
                    <div class="comp-value">{{ calculateMusclePercentage(traineeProfile) }}%</div>
                    <div class="comp-bar">
                      <div class="comp-fill muscle" :style="{ width: calculateMusclePercentage(traineeProfile) + '%' }"></div>
                    </div>
                  </div>
                </div>

                <div class="composition-item">
                  <div class="comp-icon">🟡</div>
                  <div class="comp-info">
                    <div class="comp-label">Tłuszcz</div>
                    <div class="comp-value">{{ calculateBodyFatPercentage(traineeProfile) }}%</div>
                    <div class="comp-bar">
                      <div class="comp-fill fat" :style="{ width: calculateBodyFatPercentage(traineeProfile) + '%' }"></div>
                    </div>
                  </div>
                </div>

                <div class="composition-item">
                  <div class="comp-icon">🦴</div>
                  <div class="comp-info">
                    <div class="comp-label">Kości</div>
                    <div class="comp-value">{{ calculateBoneWeight(traineeProfile) }} kg</div>
                    <div class="comp-bar">
                      <div class="comp-fill bone" :style="{ width: '15%' }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Idealna waga -->
            <div class="profile-section">
              <h4>🎯 Analiza Idealnej Wagi</h4>
              <div class="ideal-weight-section">
                <div class="weight-comparison">
                  <div class="current-weight">
                    <div class="weight-label">Obecna waga</div>
                    <div class="weight-value">{{ traineeProfile.weight }} kg</div>
                  </div>
                  <div class="weight-arrow">→</div>
                  <div class="ideal-weights">
                    <div class="ideal-item">
                      <span class="method">Robinson:</span>
                      <span class="value">{{ calculateIdealWeight(traineeProfile).robinson }} kg</span>
                    </div>
                    <div class="ideal-item">
                      <span class="method">Miller:</span>
                      <span class="value">{{ calculateIdealWeight(traineeProfile).miller }} kg</span>
                    </div>
                    <div class="ideal-item">
                      <span class="method">Devine:</span>
                      <span class="value">{{ calculateIdealWeight(traineeProfile).devine }} kg</span>
                    </div>
                  </div>
                </div>
                <div class="weight-difference" :class="getWeightDifferenceClass(traineeProfile)">
                  {{ getWeightDifferenceText(traineeProfile) }}
                </div>
              </div>
            </div>
          </div>
          <div v-else class="no-profile">
            <p>Użytkownik nie ma uzupełnionego profilu</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Notifications -->
    <div v-if="notification" class="notification" :class="notification.type">
      {{ notification.message }}
      <button @click="notification = null" class="notification-close">×</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted, inject } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const store = useAuthStore()
const router = useRouter()
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Inject the switchToView function from parent component (PanelView)
const switchToView = inject('switchToView')

// State variables
const loading = ref(true)
const trainees = ref([])
const pendingRequests = ref([])
const recentActivity = ref([])
const notification = ref(null)
const processingRequest = ref(null)
const actionType = ref('')
const pendingReports = ref(0)
const unreadMessages = ref(0)

// Profile modal state
const showProfileModal = ref(false)
const selectedTrainee = ref(null)
const traineeProfile = ref(null)
const loadingProfile = ref(false)

// API utility function
const apiRequest = async (url, options = {}) => {
  const token = store.token
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

  return response.json()
}

// Show notification
const showNotification = (message, type = 'success') => {
  notification.value = { message, type }
  setTimeout(() => {
    notification.value = null
  }, 5000)
}

// MAIN FUNCTION: Handle all navigation clicks
const handleNotificationClick = async (notificationType, trainee) => {
  try {
    console.log('Clicked notification:', notificationType, 'for trainee:', trainee.username)

    // Check if there are unread notifications to mark as read
    const hasUnreadNotifications = {
      'chat': trainee.unread_messages > 0,
      'reports': trainee.unread_reports > 0,
      'media': trainee.unread_media > 0,
      'diet-plans': trainee.pending_diet_feedback > 0
    }

    // Mark notifications as read FIRST (only if there are unread notifications)
    if (hasUnreadNotifications[notificationType]) {
      console.log(`Found ${getUnreadCount(trainee, notificationType)} unread ${notificationType} for ${trainee.username}`)
      await markNotificationsAsRead(notificationType, trainee.id)
      showNotification(`Oznaczono ${getNotificationDisplayName(notificationType)} jako przeczytane`, 'success')
    }

    if (switchToView) {
      // Store selected trainee info for context
      sessionStorage.setItem('selectedTraineeId', trainee.id.toString())
      sessionStorage.setItem('selectedTraineeName', trainee.username)

      // Map notification types to view names
      const viewMapping = {
        'chat': 'chat',
        'reports': 'reports',
        'media': 'media',
        'diet-plans': 'diet-plans'
      }

      const targetView = viewMapping[notificationType]
      if (targetView) {
        switchToView(targetView)
        showNotification(`Przechodzisz do ${getViewDisplayName(targetView)} dla ${trainee.username}`, 'success')
      }
    } else {
      console.error('switchToView function not available')
      showNotification('Błąd nawigacji - funkcja switchToView niedostępna', 'error')
    }
  } catch (error) {
    console.error('Error handling notification click:', error)
    showNotification('Wystąpił błąd podczas przetwarzania notyfikacji', 'error')
  }
}
const markNotificationsAsRead = async (notificationType, traineeId) => {
  try {
    const endpoints = {
      'chat': `/messages/mark-read?sender_id=${traineeId}`,
      'media': `/media/trainee/${traineeId}/mark-read`,
      'diet-plans': `/diet-plans/trainee/${traineeId}/mark-feedback-read`
    }

    const endpoint = endpoints[notificationType]
    if (endpoint) {
      console.log(`Marking ${notificationType} as read for trainee ${traineeId}`)
      await apiRequest(endpoint, { method: 'POST' })

      // Update local state immediately
      const trainee = trainees.value.find(t => t.id === traineeId)
      if (trainee) {
        switch(notificationType) {
          case 'chat':
            trainee.unread_messages = 0
            break
          case 'reports':
            trainee.unread_reports = 0
            // Update global pending reports count
            pendingReports.value = Math.max(0, pendingReports.value - trainee.unread_reports)
            break
          case 'media':
            trainee.unread_media = 0
            break
          case 'diet-plans':
            trainee.pending_diet_feedback = 0
            break
        }
      }

      // Recalculate total unread messages for chat
      if (notificationType === 'chat') {
        unreadMessages.value = trainees.value.reduce((sum, t) => sum + (t.unread_messages || 0), 0)
      }

      console.log(`Successfully marked ${notificationType} as read for trainee ${traineeId}`)
    }
  } catch (error) {
    console.error(`Error marking ${notificationType} as read:`, error)
    showNotification(`Nie udało się oznaczyć ${getNotificationDisplayName(notificationType)} jako przeczytane`, 'error')
  }
}

// Helper function to get display names for notifications
const getNotificationDisplayName = (notificationType) => {
  const displayNames = {
    'chat': 'wiadomości',
    'reports': 'raportów',
    'media': 'mediów',
    'diet-plans': 'planów żywieniowych'
  }
  return displayNames[notificationType] || notificationType
}

// Helper function to get unread count for specific notification type
const getUnreadCount = (trainee, notificationType) => {
  switch(notificationType) {
    case 'chat': return trainee.unread_messages || 0
    case 'reports': return trainee.unread_reports || 0
    case 'media': return trainee.unread_media || 0
    case 'diet-plans': return trainee.pending_diet_feedback || 0
    default: return 0
  }
}
// Helper function to get display names for views
const getViewDisplayName = (viewName) => {
  const displayNames = {
    'chat': 'czatu',
    'reports': 'raportów',
    'media': 'mediów',
    'diet-plans': 'planów żywieniowych'
  }
  return displayNames[viewName] || viewName
}

// Computed properties
const activeTrainingPlans = computed(() => {
  return trainees.value.filter(t => t.current_plan).length
})

// Format date for display
const formatDate = (date) => {
  const now = new Date()
  const diff = now - date
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)

  if (days > 0) return `${days} dni temu`
  if (hours > 0) return `${hours} godz. temu`
  return 'Przed chwilą'
}

// Helper functions for profile display
const getGenderText = (gender) => {
  return gender === 'male' ? 'Mężczyzna' : 'Kobieta'
}

const getActivityText = (level) => {
  console.log('Activity level received:', level, 'Type:', typeof level) // DEBUG

  // Konwertuj na liczbę jeśli to string
  const numLevel = typeof level === 'string' ? parseFloat(level) : level

  if (typeof numLevel !== 'number' || isNaN(numLevel)) {
    return `Nieznana (${level})`
  }

  // Standardowe wartości aktywności
  const activityLevels = [
    { value: 1.2, label: 'Bardzo niska' },
    { value: 1.375, label: 'Niska' },
    { value: 1.55, label: 'Umiarkowana' },
    { value: 1.725, label: 'Wysoka' },
    { value: 1.9, label: 'Bardzo wysoka' }
  ]

  // Znajdź najbliższą wartość
  let closest = activityLevels[0]
  let minDifference = Math.abs(numLevel - closest.value)

  for (const activity of activityLevels) {
    const difference = Math.abs(numLevel - activity.value)
    if (difference < minDifference) {
      minDifference = difference
      closest = activity
    }
  }

  return closest.label
}

const getGoalText = (goal) => {
  const goals = {
    'reduction': 'Utrata wagi',
    'recomposition': 'Rekompozycja ciała',
    'mass': 'Przyrost masy',
    'Utrata wagi': 'Utrata wagi',
    'Utrzymanie wagi': 'Utrzymanie wagi',
    'Przyrost masy': 'Przyrost masy',
    'Budowanie mięśni': 'Budowanie mięśni',
    'Rekompozycja': 'Rekompozycja ciała'
  }
  return goals[goal] || goal || 'Nie ustawiono'
}

// ========== ZAAWANSOWANE OBLICZENIA ==========

// BMI Calculations
const calculateBMI = (profile) => {
  if (!profile?.height || !profile?.weight) return 'N/A'
  const h = profile.height / 100
  return Math.round((profile.weight / (h * h)) * 10) / 10
}

const getBMICategory = (profile) => {
  const bmi = calculateBMI(profile)
  if (bmi === 'N/A') return { label: 'N/A', class: '' }

  if (bmi < 18.5) return { label: 'Niedowaga', class: 'underweight' }
  if (bmi < 25) return { label: 'Waga prawidłowa', class: 'normal' }
  if (bmi < 30) return { label: 'Nadwaga', class: 'overweight' }
  return { label: 'Otyłość', class: 'obese' }
}

const getBMIPosition = (profile) => {
  const bmi = calculateBMI(profile)
  if (bmi === 'N/A') return 0

  if (bmi < 18.5) return (bmi / 18.5) * 25
  if (bmi < 25) return 25 + ((bmi - 18.5) / 6.5) * 25
  if (bmi < 30) return 50 + ((bmi - 25) / 5) * 25
  return Math.min(75 + ((bmi - 30) / 10) * 25, 100)
}

// Metabolizm i kalorie
const calculateBMR = (profile) => {
  if (!profile?.weight || !profile?.height || !profile?.age || !profile?.gender) return 'N/A'

  const { weight, height, age, gender } = profile
  return Math.round(gender === 'male'
      ? 10 * weight + 6.25 * height - 5 * age + 5
      : 10 * weight + 6.25 * height - 5 * age - 161)
}

const calculateDailyCalories = (profile) => {
  const bmr = calculateBMR(profile)
  if (bmr === 'N/A' || !profile?.activity_level) return 'N/A'
  return Math.round(bmr * profile.activity_level)
}

// Skład ciała
const calculateWaterPercentage = (profile) => {
  if (!profile?.gender || !profile?.age || !profile?.weight || !profile?.height) return 0

  const { gender, age, weight, height } = profile
  const liters = gender === 'male'
      ? 2.447 + 0.09156 * weight + 0.1074 * height - 0.09 * age
      : -2.097 + 0.1069 * height + 0.2466 * weight
  return Math.round((liters / weight) * 100)
}

const calculateBoneWeight = (profile) => {
  if (!profile?.weight || !profile?.gender) return 0
  return Math.round(profile.weight * (profile.gender === 'male' ? 0.15 : 0.12) * 10) / 10
}

const calculateMusclePercentage = (profile) => {
  if (!profile?.gender || !profile?.age) return 0

  const { gender, age } = profile
  const bmi = calculateBMI(profile)
  if (bmi === 'N/A') return 0

  let base = gender === 'male' ? 45 : 38
  let ageAdj = Math.max(0, (age - 25) * 0.3)
  let bmiAdj = bmi < 18.5 ? -5 : bmi > 30 ? -8 : bmi > 25 ? -2 : 0

  return Math.max(25, Math.round(base - ageAdj + bmiAdj))
}

const calculateBodyFatPercentage = (profile) => {
  if (!profile?.gender || !profile?.age) return 0

  const { gender, age } = profile
  const bmi = calculateBMI(profile)
  if (bmi === 'N/A') return 0

  const genderF = gender === 'male' ? 1 : 0
  return Math.max(5, Math.round(1.2 * bmi + 0.23 * age - 10.8 * genderF - 5.4))
}

// Idealna waga
const calculateIdealWeight = (profile) => {
  if (!profile?.height || !profile?.gender) {
    return { robinson: 0, miller: 0, devine: 0 }
  }

  const { height, gender } = profile
  const inch = height / 2.54
  let rob, mill, dev

  if (gender === 'male') {
    rob = 52 + 1.9 * (inch - 60)
    mill = 56.2 + 1.41 * (inch - 60)
    dev = 50 + 2.3 * (inch - 60)
  } else {
    rob = 49 + 1.7 * (inch - 60)
    mill = 53.1 + 1.36 * (inch - 60)
    dev = 45.5 + 2.3 * (inch - 60)
  }

  return {
    robinson: Math.round(rob * 10) / 10,
    miller: Math.round(mill * 10) / 10,
    devine: Math.round(dev * 10) / 10
  }
}

const getWeightDifferenceText = (profile) => {
  if (!profile?.weight) return ''

  const idealWeights = calculateIdealWeight(profile)
  const avg = (idealWeights.robinson + idealWeights.miller + idealWeights.devine) / 3
  const diff = Math.round((profile.weight - avg) * 10) / 10

  if (Math.abs(diff) < 2) return 'Idealna waga! 🎉'
  return diff > 0 ? `${diff} kg powyżej idealnej` : `${Math.abs(diff)} kg poniżej idealnej`
}

const getWeightDifferenceClass = (profile) => {
  if (!profile?.weight) return ''

  const idealWeights = calculateIdealWeight(profile)
  const avg = (idealWeights.robinson + idealWeights.miller + idealWeights.devine) / 3
  const diff = profile.weight - avg

  if (Math.abs(diff) < 2) return 'ideal'
  return diff > 0 ? 'above' : 'below'
}

// ============== RESZTA KODU BEZ ZMIAN ==============

// Data fetching functions
const loadTrainees = async () => {
  try {
    const data = await apiRequest('/trainees')

    // Deduplikowanie podopiecznych na podstawie ID
    const uniqueTrainees = new Map()

    // Grupujemy podopiecznych według ID
    data.forEach(trainee => {
      if (uniqueTrainees.has(trainee.id)) {
        // Jeśli już mamy tego podopiecznego, dodajemy tylko plan treningowy
        const existing = uniqueTrainees.get(trainee.id)
        if (trainee.current_plan && !existing.training_plans.find(p => p.id === trainee.current_plan.id)) {
          existing.training_plans.push(trainee.current_plan)
        }
        // Zachowujemy najnowsze informacje o podopiecznym
        existing.last_updated = Math.max(existing.last_updated || 0, new Date(trainee.updated_at || 0).getTime())
      } else {
        // Nowy podopieczny
        uniqueTrainees.set(trainee.id, {
          ...trainee,
          training_plans: trainee.current_plan ? [trainee.current_plan] : [],
          last_updated: new Date(trainee.updated_at || 0).getTime()
        })
      }
    })

    // Konwertujemy Map z powrotem na array
    trainees.value = Array.from(uniqueTrainees.values())

    // Fetch detailed notification data for each trainee
    await Promise.all(trainees.value.map(async (trainee) => {
      try {
        // Get unread messages count
        const messageData = await apiRequest(`/messages/unread?sender_id=${trainee.id}`)
        trainee.unread_messages = messageData.count || 0

        // Get unread reports count
        const reportsData = await apiRequest(`/reports/trainer/trainee/${trainee.id}/unread`)
        trainee.unread_reports = reportsData.count || 0

        // Get unread media count
        const mediaData = await apiRequest(`/media/trainee/${trainee.id}/unread`)
        trainee.unread_media = mediaData.count || 0

        // Get pending diet feedback count
        const dietData = await apiRequest(`/diet-plans/trainee/${trainee.id}/pending-feedback`)
        trainee.pending_diet_feedback = dietData.count || 0

      } catch (error) {
        console.error(`Error fetching notifications for trainee ${trainee.id}:`, error)
        trainee.unread_messages = 0
        trainee.unread_reports = 0
        trainee.unread_media = 0
        trainee.pending_diet_feedback = 0
      }
    }))

    // Calculate total unread messages
    unreadMessages.value = trainees.value.reduce((sum, t) => sum + (t.unread_messages || 0), 0)

    console.log('Deduplicated trainees:', trainees.value.length) // DEBUG

  } catch (err) {
    console.error('Error fetching trainees:', err)
    showNotification(`Failed to fetch trainees: ${err.message}`, 'error')
  }
}

const loadPendingRequests = async () => {
  try {
    const data = await apiRequest('/trainer-requests/for-me')
    pendingRequests.value = data?.filter(req => req.status === 'pending') || []
  } catch (err) {
    console.error('Error fetching pending requests:', err)
    showNotification(`Failed to fetch pending requests: ${err.message}`, 'error')
  }
}

const loadRecentActivity = async () => {
  try {
    const activities = []

    pendingRequests.value.forEach((req, index) => {
      activities.push({
        id: `req-${req.id}`,
        icon: '📨',
        text: `Nowa prośba o współpracę od ${req.client_username || 'Użytkownik ' + req.client_id}`,
        timestamp: new Date(req.created_at)
      })
    })

    trainees.value.forEach(trainee => {
      if (trainee.unread_messages > 0) {
        activities.push({
          id: `msg-${trainee.id}`,
          icon: '💬',
          text: `${trainee.unread_messages} nowych wiadomości od ${trainee.username}`,
          timestamp: new Date(Date.now() - Math.random() * 24 * 60 * 60 * 1000)
        })
      }
      if (trainee.unread_reports > 0) {
        activities.push({
          id: `rep-${trainee.id}`,
          icon: '📊',
          text: `${trainee.unread_reports} nowych raportów od ${trainee.username}`,
          timestamp: new Date(Date.now() - Math.random() * 24 * 60 * 60 * 1000)
        })
      }
    })

    recentActivity.value = activities.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
  } catch (err) {
    console.error('Error creating activity feed:', err)
  }
}

const loadPendingReportsCount = async () => {
  try {
    const data = await apiRequest('/reports/trainer/pending')
    pendingReports.value = data.count || 0
  } catch (err) {
    console.error('Error fetching pending reports count:', err)
    pendingReports.value = 0
  }
}

// Profile modal functions
const viewTraineeProfile = async (trainee) => {
  selectedTrainee.value = trainee
  showProfileModal.value = true
  loadingProfile.value = true

  try {
    const data = await apiRequest(`/profile/trainee/${trainee.id}`)
    traineeProfile.value = data.profile
  } catch (error) {
    console.error('Error fetching trainee profile:', error)
    traineeProfile.value = null
    showNotification('Nie udało się pobrać profilu podopiecznego', 'error')
  } finally {
    loadingProfile.value = false
  }
}

const closeProfileModal = () => {
  showProfileModal.value = false
  selectedTrainee.value = null
  traineeProfile.value = null
}

// Request management functions
const approveRequest = async (requestId) => {
  processingRequest.value = requestId
  actionType.value = 'approve'

  try {
    await apiRequest(`/trainer-requests/${requestId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status: 'approved' })
    })

    pendingRequests.value = pendingRequests.value.filter(req => req.id !== requestId)
    await loadTrainees()
    await loadRecentActivity()

    showNotification('Prośba została zaakceptowana i dodano nowego podopiecznego', 'success')
  } catch (err) {
    console.error('Error approving request:', err)
    showNotification(`Failed to approve request: ${err.message}`, 'error')
  } finally {
    processingRequest.value = null
    actionType.value = ''
  }
}

const rejectRequest = async (requestId) => {
  processingRequest.value = requestId
  actionType.value = 'reject'

  try {
    await apiRequest(`/trainer-requests/${requestId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status: 'rejected' })
    })

    pendingRequests.value = pendingRequests.value.filter(req => req.id !== requestId)
    showNotification('Prośba została odrzucona', 'success')
  } catch (err) {
    console.error('Error rejecting request:', err)
    showNotification(`Failed to reject request: ${err.message}`, 'error')
  } finally {
    processingRequest.value = null
    actionType.value = ''
  }
}

// UPDATED: Trainee management functions to use switchToView instead of router.push
const createTrainingPlan = (trainee) => {
  // Store trainee info and switch to training plans view
  sessionStorage.setItem('selectedTraineeId', trainee.id.toString())
  sessionStorage.setItem('selectedTraineeName', trainee.username)
  sessionStorage.setItem('createNewPlan', 'true')

  if (switchToView) {
    switchToView('training-plans')
    showNotification(`Tworzysz nowy plan treningowy dla ${trainee.username}`, 'success')
  }
}

// Dashboard refresh functionality
let refreshInterval = null

const loadDashboardData = async () => {
  loading.value = true

  try {
    await Promise.all([
      loadTrainees(),
      loadPendingRequests(),
      loadPendingReportsCount()
    ])

    await loadRecentActivity()
  } catch (err) {
    console.error('Error loading dashboard data:', err)
    showNotification('Failed to load dashboard data', 'error')
  } finally {
    loading.value = false
  }
}

const startAutoRefresh = () => {
  refreshInterval = setInterval(loadDashboardData, 2 * 60 * 1000)
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
}

// Lifecycle hooks
onMounted(async () => {
  await loadDashboardData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.trainer-dashboard {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  width: 95%;
  height: 90%;
  margin: 2rem;
  color: #181818;
  overflow-y: auto;
  position: relative;
}

.trainer-dashboard::-webkit-scrollbar {
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

.requests-section {
  margin-bottom: 2rem;
}

.requests-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.request-card {
  background: #fffbef;
  border: 1px solid #ffe58f;
  border-radius: 8px;
  padding: 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.request-info {
  flex: 1;
}

.request-info h4 {
  margin: 0 0 0.5rem 0;
  color: #8c6200;
}

.request-message {
  color: #666;
  margin: 0.5rem 0;
  font-size: 0.9rem;
}

.request-date {
  font-size: 0.8rem;
  color: #888;
}

.request-actions {
  display: flex;
  gap: 0.5rem;
}

.approve-btn, .reject-btn {
  padding: 0.5rem 1rem;
  font-weight: 600;
  border-radius: 6px;
  cursor: pointer;
  transition: 0.2s;
  border: none;
}

.approve-btn {
  background: #e8f5e8;
  color: #2e7d32;
}

.approve-btn:hover {
  background: #c8e6c9;
}

.reject-btn {
  background: #ffebee;
  color: #d32f2f;
}

.reject-btn:hover {
  background: #ffcdd2;
}

.approve-btn:disabled, .reject-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.trainees-section {
  margin-bottom: 2rem;
}

.loading-state, .empty-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.trainees-list {
  display: grid;
  gap: 1rem;
}

.trainee-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  transition: 0.2s;
}

.trainee-card:hover {
  border-color: #b7f5b7;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.trainee-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.trainee-info h4 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.3rem;
  font-weight: 700;
}

.trainee-status {
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: bold;
  display: inline-block;
}

.trainee-status.active {
  background: #e8f5e8;
  color: #2e7d32;
}

.trainee-status.inactive {
  background: #ffebee;
  color: #d32f2f;
}

.last-activity {
  font-size: 0.8rem;
  color: #666;
  margin-top: 0.25rem;
}

/* Action Buttons WITH NOTIFICATION BADGES */
.trainee-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.5rem;
  margin-top: 1rem;
}

.action-btn {
  padding: 0.5rem 0.75rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  font-weight: 600;
  transition: 0.2s;
  text-align: center;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
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

.profile-btn {
  background: #e3f2fd;
  color: #1976d2;
}

.profile-btn:hover {
  background: #bbdefb;
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

/* Profile Modal */
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

.profile-modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
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

.modal-content {
  padding: 1.5rem;
}

.profile-section {
  margin-bottom: 2rem;
}

.profile-section h4 {
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
  background: #f8f9fa;
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

/* ========== NOWE STYLE DLA ZAAWANSOWANYCH OBLICZEŃ ========== */

/* BMI Section */
.bmi-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.bmi-display {
  text-align: center;
  margin-bottom: 1.5rem;
}

.bmi-value {
  font-size: 3rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.bmi-value.underweight { color: #1976d2; }
.bmi-value.normal { color: #2e7d32; }
.bmi-value.overweight { color: #f57c00; }
.bmi-value.obese { color: #d32f2f; }

.bmi-category {
  font-size: 1.125rem;
  font-weight: 600;
  color: #666;
}

.bmi-bar {
  position: relative;
  height: 20px;
  background: linear-gradient(90deg, #1976d2 0%, #1976d2 25%, #2e7d32 25%, #2e7d32 50%, #f57c00 50%, #f57c00 75%, #d32f2f 75%);
  border-radius: 10px;
  margin: 1rem 0;
}

.bmi-indicator {
  position: absolute;
  top: -5px;
  width: 10px;
  height: 30px;
  background: #333;
  border-radius: 5px;
  transform: translateX(-50%);
}

.bmi-ranges {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  margin-top: 0.5rem;
}

.range {
  font-weight: 600;
  text-align: center;
  flex: 1;
}
.range.normal { color: #ffffff; }
.range.overweight { color: #ffffff; }
.range.obese { color: #ffffff; }
.range.underweight { color: #ffffff; }

/* Body Composition */
.body-composition {
  display: grid;
  gap: 1rem;
}

.composition-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.comp-icon {
  font-size: 1.5rem;
  width: 48px;
  text-align: center;
}

.comp-info {
  flex: 1;
}

.comp-label {
  font-weight: 600;
  color: #2d5016;
  margin-bottom: 0.25rem;
}

.comp-value {
  font-size: 1.125rem;
  font-weight: bold;
  color: #333;
  margin-bottom: 0.5rem;
}

.comp-bar {
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.comp-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.comp-fill.water { background: linear-gradient(90deg, #1976d2, #42a5f5); }
.comp-fill.muscle { background: linear-gradient(90deg, #d32f2f, #f44336); }
.comp-fill.fat { background: linear-gradient(90deg, #ff9800, #ffc107); }
.comp-fill.bone { background: linear-gradient(90deg, #795548, #a1887f); }

/* Ideal Weight Section */
.ideal-weight-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.weight-comparison {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.current-weight {
  text-align: center;
}

.weight-label {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.25rem;
}

.weight-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #2d5016;
}

.weight-arrow {
  font-size: 1.5rem;
  color: #666;
  margin: 0 1rem;
}

.ideal-weights {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.ideal-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background: white;
  border-radius: 6px;
  border: 1px solid #e9ecef;
  min-width: 200px;
}

.method {
  font-weight: 600;
  color: #333;
}

.value {
  font-weight: bold;
  color: #2d5016;
}

.weight-difference {
  text-align: center;
  padding: 1rem;
  border-radius: 8px;
  font-weight: bold;
  font-size: 1.125rem;
}

.weight-difference.ideal {
  background: #e8f5e8;
  color: #2e7d32;
}

.weight-difference.above {
  background: #fff3e0;
  color: #f57c00;
}

.weight-difference.below {
  background: #e3f2fd;
  color: #1976d2;
}

.no-profile {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.activity-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
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
}

.activity-time {
  color: #666;
  font-size: 0.8rem;
  margin-left: auto;
}

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

@media (max-width: 768px) {
  .trainer-dashboard {
    padding: 1rem;
    margin: 1rem;
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

  .request-card {
    flex-direction: column;
    align-items: stretch;
  }

  .request-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
  }

  .profile-modal {
    width: 95%;
    margin: 1rem;
  }

  .weight-comparison {
    flex-direction: column;
    gap: 1rem;
  }

  .weight-arrow {
    transform: rotate(90deg);
  }

  .ideal-weights {
    width: 100%;
  }

  .body-composition {
    grid-template-columns: 1fr;
  }
}
</style>