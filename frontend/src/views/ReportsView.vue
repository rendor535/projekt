<template>
  <div class="reports-view">
    <!-- WIDOK DLA TRENERA -->
    <div v-if="isTrainer" class="trainer-view">
      <!-- Header -->
      <div class="page-header">
        <div class="header-content">
          <h1>📊 Raporty podopiecznych</h1>
          <p>Przeglądaj raporty i statystyki swoich podopiecznych</p>
        </div>
      </div>

      <!-- Główny układ trenera - dwie kolumny -->
      <div class="trainer-layout">
        <!-- Lewa kolumna - Lista podopiecznych -->
        <div class="trainees-section">
          <div class="section-header">
            <h2 class="section-title">👥 Podopieczni</h2>
          </div>

          <div v-if="loadingTrainerStats" class="loading">Ładowanie podopiecznych...</div>

          <div v-else class="trainees-list">
            <div
                v-for="trainee in (trainerStats?.trainee_stats || [])"
                :key="trainee.trainee_id"
                class="trainee-card"
                :class="{ active: selectedTraineeId === trainee.trainee_id }"
                @click="selectTrainee(trainee)"
            >
              <div class="trainee-avatar">{{ trainee.trainee_name?.charAt(0)?.toUpperCase() || '?' }}</div>
              <div class="trainee-info">
                <h4>{{ trainee.trainee_name || 'Nieznany' }}</h4>
                <div class="trainee-stats">
                  <span class="stat-item">📊 {{ trainee.report_count || 0 }} raportów</span>
                  <span class="stat-item" v-if="trainee.latest_weight">
                    ⚖️ {{ trainee.latest_weight.toFixed(1) }}kg
                  </span>
                  <span class="stat-item" v-if="trainee.last_report_date">
                    📅 {{ formatDate(trainee.last_report_date) }}
                  </span>
                </div>
              </div>
              <div class="select-arrow">→</div>
            </div>

            <div v-if="!(trainerStats?.trainee_stats?.length)" class="empty-state">
              <p>Nie masz jeszcze podopiecznych.</p>
            </div>
          </div>
        </div>

        <!-- Prawa kolumna - Raporty lub Statystyki -->
        <div class="main-content">
          <!-- Lista raportów (domyślny widok) -->
          <div v-if="!selectedTraineeId" class="reports-list">
            <div class="section-header">
              <h2 class="section-title">📋 Wszystkie raporty</h2>
            </div>

            <div v-if="loadingAllReports" class="loading">Ładowanie raportów...</div>

            <div v-else class="reports-grid">
              <div
                  v-for="report in (allTraineeReports || [])"
                  :key="report.id"
                  class="report-card clickable"
                  @click="showReportDetails(report)"
              >
                <div class="report-header">
                  <div class="report-date">{{ formatDate(report.created_at) }}</div>
                  <div class="report-author">{{ getTraineeName(report.client_id) }}</div>
                </div>
                <div class="report-summary">
                  <div class="summary-row">
                    <span v-if="report.weight">⚖️ {{ report.weight.toFixed(1) }}kg</span>
                    <span v-if="report.body_fat">🧮 {{ report.body_fat.toFixed(1) }}%</span>
                  </div>
                  <p v-if="report.notes" class="report-notes">{{ report.notes.substring(0, 100) }}...</p>
                </div>
                <div class="click-hint">Kliknij aby zobaczyć szczegóły</div>
              </div>

              <div v-if="!(allTraineeReports?.length)" class="empty-state">
                <p>Brak raportów od podopiecznych.</p>
              </div>
            </div>
          </div>

          <!-- Statystyki wybranego podopiecznego -->
          <div v-else class="trainee-details">
            <div class="section-header">
              <h2 class="section-title">📈 {{ selectedTrainee?.trainee_name || 'Podopieczny' }} - Statystyki</h2>
              <button class="btn-back" @click="deselectTrainee">← Powrót do listy</button>
            </div>

            <div v-if="loadingTraineeData" class="loading">Ładowanie danych podopiecznego...</div>

            <div v-else>
              <!-- Dashboard wybranego podopiecznego -->
              <div class="dashboard-cards">
                <div class="dashboard-card">
                  <div class="card-icon">📊</div>
                  <div class="card-content">
                    <h3>{{ selectedTrainee?.report_count || 0 }}</h3>
                    <p>Wszystkich raportów</p>
                  </div>
                </div>
                <div class="dashboard-card">
                  <div class="card-icon">⚖️</div>
                  <div class="card-content">
                    <h3>{{ selectedTrainee?.latest_weight?.toFixed(1) || 'Brak' }}</h3>
                    <p>Aktualna waga (kg)</p>
                  </div>
                </div>
                <div class="dashboard-card">
                  <div class="card-icon">📈</div>
                  <div class="card-content">
                    <h3>{{ getWeightChange(selectedTrainee) }}</h3>
                    <p>Zmiana wagi (kg)</p>
                  </div>
                </div>
              </div>

              <!-- Enhanced Charts Section -->
              <div class="charts-section" v-if="enhancedChartData">
                <div class="section-header">
                  <h3 class="section-title">📈 Wykresy postępów</h3>
                </div>

                <!-- Chart Type Buttons -->
                <div class="chart-type-selector">
                  <button
                      v-for="chartType in availableChartTypes"
                      :key="chartType.key"
                      @click="activeChart = chartType.key"
                      class="chart-type-btn"
                      :class="{ active: activeChart === chartType.key }"
                      :disabled="!chartType.hasData"
                  >
                    <span class="chart-icon">{{ chartType.icon }}</span>
                    <span class="chart-label">{{ chartType.label }}</span>
                    <span v-if="chartType.hasData" class="data-count">({{ chartType.dataCount }})</span>
                  </button>
                </div>

                <!-- Chart Container -->
                <div class="enhanced-chart-container">
                  <div class="chart-header">
                    <h4>{{ getCurrentChartTitle() }}</h4>
                    <div class="chart-stats">
                      <span v-if="getCurrentChartStats().latest">Najnowsze: {{ getCurrentChartStats().latest }}</span>
                      <span v-if="getCurrentChartStats().change" :class="getChangeClass(getCurrentChartStats().change)">
                        Zmiana: {{ getCurrentChartStats().change }}
                      </span>
                    </div>
                  </div>
                  <canvas ref="trainerChart" width="800" height="400"></canvas>
                </div>
              </div>

              <!-- Ostatnie raporty wybranego podopiecznego -->
              <div class="recent-reports">
                <h3 class="section-title">Ostatnie raporty</h3>
                <div class="reports-timeline">
                  <div
                      v-for="report in (traineeReports || []).slice(0, 5)"
                      :key="report.id"
                      class="timeline-item"
                      @click="showReportDetails(report)"
                  >
                    <div class="timeline-date">{{ formatDate(report.created_at) }}</div>
                    <div class="timeline-content">
                      <div class="timeline-data">
                        <span v-if="report.weight">⚖️ {{ report.weight.toFixed(1) }}kg</span>
                        <span v-if="report.body_fat">🧮 {{ report.body_fat.toFixed(1) }}%</span>
                      </div>
                      <p v-if="report.notes">{{ report.notes }}</p>
                    </div>
                  </div>

                  <div v-if="!(traineeReports?.length)" class="empty-state">
                    <p>Brak raportów dla tego podopiecznego.</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- WIDOK DLA UCZNIA -->
    <div v-else class="student-view">
      <!-- Header -->
      <div class="page-header">
        <div class="header-content">
          <h1>📊 Moje raporty postępów</h1>
          <p>Dodawaj raporty i śledź swoje postępy</p>
        </div>
        <div class="header-actions">
          <button class="btn-add-report" @click="showCreateModal = true">
            <span class="btn-icon">➕</span>
            Dodaj raport
          </button>
        </div>
      </div>

      <div v-if="loadingUserData" class="loading">Ładowanie danych...</div>

      <div v-else>
        <!-- Moje statystyki -->
        <div class="my-stats">
          <div class="dashboard-cards">
            <div class="dashboard-card">
              <div class="card-icon">📊</div>
              <div class="card-content">
                <h3>{{ userStats?.total_reports || 0 }}</h3>
                <p>Wszystkich raportów</p>
              </div>
            </div>
            <div class="dashboard-card">
              <div class="card-icon">⚖️</div>
              <div class="card-content">
                <h3>{{ userStats?.latest_weight ? userStats.latest_weight.toFixed(1) : 'Brak' }}</h3>
                <p>Aktualna waga (kg)</p>
              </div>
            </div>
            <div class="dashboard-card">
              <div class="card-icon">📈</div>
              <div class="card-content">
                <h3>{{ weightChangeDisplay }}</h3>
                <p>Zmiana wagi (kg)</p>
              </div>
            </div>
            <div class="dashboard-card">
              <div class="card-icon">🎯</div>
              <div class="card-content">
                <h3>{{ userStats?.avg_body_fat ? userStats.avg_body_fat.toFixed(1) : 'Brak' }}%</h3>
                <p>Średnia tkanka tłuszczowa</p>
              </div>
            </div>
          </div>

          <!-- Enhanced Student Charts -->
          <div class="charts-section" v-if="enhancedChartData">
            <div class="section-header">
              <h2 class="section-title">📈 Twoje postępy</h2>
            </div>

            <!-- Chart Type Buttons -->
            <div class="chart-type-selector">
              <button
                  v-for="chartType in availableChartTypes"
                  :key="chartType.key"
                  @click="activeChart = chartType.key"
                  class="chart-type-btn"
                  :class="{ active: activeChart === chartType.key }"
                  :disabled="!chartType.hasData"
              >
                <span class="chart-icon">{{ chartType.icon }}</span>
                <span class="chart-label">{{ chartType.label }}</span>
                <span v-if="chartType.hasData" class="data-count">({{ chartType.dataCount }})</span>
              </button>
            </div>

            <!-- Chart Container -->
            <div class="enhanced-chart-container">
              <div class="chart-header">
                <h4>{{ getCurrentChartTitle() }}</h4>
                <div class="chart-stats">
                  <span v-if="getCurrentChartStats().latest">Najnowsze: {{ getCurrentChartStats().latest }}</span>
                  <span v-if="getCurrentChartStats().change" :class="getChangeClass(getCurrentChartStats().change)">
                    Zmiana: {{ getCurrentChartStats().change }}
                  </span>
                </div>
              </div>
              <canvas ref="studentChart" width="800" height="400"></canvas>
            </div>
          </div>

          <!-- Moja historia raportów -->
          <div class="my-reports-history">
            <h2 class="section-title">Historia twoich raportów</h2>
            <div class="reports-timeline">
              <div
                  v-for="report in (myReports || [])"
                  :key="report.id"
                  class="timeline-item"
                  @click="showReportDetails(report)"
              >
                <div class="timeline-date">{{ formatDate(report.created_at) }}</div>
                <div class="timeline-content">
                  <div class="timeline-data">
                    <span v-if="report.weight">⚖️ {{ report.weight.toFixed(1) }}kg</span>
                    <span v-if="report.body_fat">🧮 {{ report.body_fat.toFixed(1) }}%</span>
                  </div>
                  <p v-if="report.notes">{{ report.notes }}</p>
                  <small class="added-by">
                    {{ report.trainer_username !== 'Self-report' ? `Dodane przez: ${report.trainer_username}` : 'Raport własny' }}
                  </small>
                </div>
              </div>

              <div v-if="!(myReports?.length)" class="empty-state">
                <p>Nie masz jeszcze żadnych raportów. Dodaj swój pierwszy raport!</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal dodawania raportu (tylko dla ucznia) -->
    <div v-if="showCreateModal && !isTrainer" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>➕ Dodaj nowy raport postępów</h3>
          <button class="btn-close" @click="showCreateModal = false">✕</button>
        </div>

        <div class="modal-body">
          <div class="form-row">
            <div class="form-group">
              <label>Waga (kg) *</label>
              <input
                  type="number"
                  step="0.1"
                  class="form-input"
                  v-model.number="newReport.weight"
                  placeholder="np. 75.5"
                  required
              >
            </div>
            <div class="form-group">
              <label>Tkanka tłuszczowa (%)</label>
              <input
                  type="number"
                  step="0.1"
                  class="form-input"
                  v-model.number="newReport.body_fat"
                  placeholder="np. 15.2"
              >
            </div>
          </div>

          <div class="measurements-form">
            <h4>Pomiary ciała (opcjonalne)</h4>
            <div class="measurements-grid">
              <div class="form-group">
                <label>Klatka piersiowa (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.klatka"
                    placeholder="np. 100"
                >
              </div>
              <div class="form-group">
                <label>Pas (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.pas"
                    placeholder="np. 80"
                >
              </div>
              <div class="form-group">
                <label>Biodra (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.biodra"
                    placeholder="np. 95"
                >
              </div>
              <div class="form-group">
                <label>Biceps (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.biceps"
                    placeholder="np. 35"
                >
              </div>
              <div class="form-group">
                <label>Udo (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.udo"
                    placeholder="np. 60"
                >
              </div>
              <div class="form-group">
                <label>Łydka (cm)</label>
                <input
                    type="number"
                    step="0.1"
                    class="form-input"
                    v-model.number="measurements.lydka"
                    placeholder="np. 40"
                >
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>Notatki</label>
            <textarea
                class="form-textarea"
                v-model="newReport.notes"
                placeholder="Dodatkowe uwagi o postępach, samopoczuciu, treningach..."
            ></textarea>
          </div>
        </div>

        <div class="modal-navigation">
          <button class="btn-cancel" @click="showCreateModal = false">Anuluj</button>
          <button
              class="btn-create"
              @click="createReport"
              :disabled="!newReport.weight || loading"
          >
            {{ loading ? 'Dodawanie...' : 'Dodaj raport' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal szczegółów raportu -->
    <div v-if="selectedReport" class="modal-overlay" @click="selectedReport = null">
      <div class="modal-content modal-large" @click.stop>
        <div class="modal-header">
          <h3>📊 Szczegóły raportu z {{ formatDate(selectedReport.created_at) }}</h3>
          <button class="btn-close" @click="selectedReport = null">✕</button>
        </div>

        <div class="modal-body">
          <div class="report-details">
            <div class="detail-section">
              <h4>Podstawowe dane</h4>
              <div class="detail-grid">
                <div class="detail-item" v-if="selectedReport.weight">
                  <span class="detail-label">Waga:</span>
                  <span class="detail-value">{{ selectedReport.weight.toFixed(1) }} kg</span>
                </div>
                <div class="detail-item" v-if="selectedReport.body_fat">
                  <span class="detail-label">Tkanka tłuszczowa:</span>
                  <span class="detail-value">{{ selectedReport.body_fat.toFixed(1) }}%</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">Data dodania:</span>
                  <span class="detail-value">{{ formatDateTime(selectedReport.created_at) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">Dodane przez:</span>
                  <span class="detail-value">{{ selectedReport.trainer_username !== 'Self-report' ? selectedReport.trainer_username : 'Raport własny' }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section" v-if="selectedReport.measurements && selectedReport.measurements !== '{}' && selectedReport.measurements !== ''">
              <h4>Pomiary ciała</h4>
              <div class="measurements-grid">
                <div
                    v-for="(value, key) in parseMeasurements(selectedReport.measurements)"
                    :key="key"
                    class="measurement-item"
                >
                  <span class="measurement-label">{{ translateMeasurement(key) }}:</span>
                  <span class="measurement-value">{{ value }} cm</span>
                </div>
              </div>
            </div>

            <div class="detail-section" v-if="selectedReport.notes">
              <h4>Notatki</h4>
              <div class="notes-content">
                {{ selectedReport.notes }}
              </div>
            </div>
          </div>
        </div>

        <div class="modal-navigation">
          <button class="btn-cancel" @click="selectedReport = null">Zamknij</button>
        </div>
      </div>
    </div>

    <!-- Error Toast -->
    <div v-if="errorMessage" class="error-toast" @click="errorMessage = ''">
      ❌ {{ errorMessage }}
    </div>

    <!-- Success Toast -->
    <div v-if="successMessage" class="success-toast" @click="successMessage = ''">
      ✅ {{ successMessage }}
    </div>
  </div>
</template>

<script setup>
import {ref, computed, onMounted, nextTick, watch, onUnmounted} from 'vue'
import { useAuthStore } from '@/stores/auth'

const store = useAuthStore()

// Stan komponentu
const showCreateModal = ref(false)
const activeChart = ref('weight')
const loading = ref(false)
const selectedReport = ref(null)
const errorMessage = ref('')
const successMessage = ref('')

// Loading states
const loadingUserData = ref(false)
const loadingTrainerStats = ref(false)
const loadingAllReports = ref(false)
const loadingTraineeData = ref(false)

// Dane z backendu
const userStats = ref(null)
const trainerStats = ref(null)
const enhancedChartData = ref(null)
const myReports = ref([])
const traineeReports = ref([])
const allTraineeReports = ref([])

// Wybrany podopieczny (dla trenera)
const selectedTraineeId = ref(null)
const selectedTrainee = ref(null)

// Canvas refs
const trainerChart = ref(null)
const studentChart = ref(null)

// Nowy raport i pomiary
const newReport = ref({
  weight: null,
  body_fat: null,
  notes: ''
})

const measurements = ref({
  klatka: null,
  pas: null,
  biodra: null,
  biceps: null,
  udo: null,
  lydka: null
})

// Enhanced Chart Configuration
const chartTypes = {
  weight: {
    icon: '⚖️',
    label: 'Waga',
    unit: 'kg',
    color: '#3b82f6',
    gradientStart: 'rgba(59, 130, 246, 0.3)',
    gradientEnd: 'rgba(59, 130, 246, 0.05)'
  },
  body_fat: {
    icon: '🧮',
    label: 'Tkanka tłuszczowa',
    unit: '%',
    color: '#ef4444',
    gradientStart: 'rgba(239, 68, 68, 0.3)',
    gradientEnd: 'rgba(239, 68, 68, 0.05)'
  },
  klatka: {
    icon: '💪',
    label: 'Klatka piersiowa',
    unit: 'cm',
    color: '#8b5cf6',
    gradientStart: 'rgba(139, 92, 246, 0.3)',
    gradientEnd: 'rgba(139, 92, 246, 0.05)'
  },
  pas: {
    icon: '⭕',
    label: 'Pas',
    unit: 'cm',
    color: '#f59e0b',
    gradientStart: 'rgba(245, 158, 11, 0.3)',
    gradientEnd: 'rgba(245, 158, 11, 0.05)'
  },
  biodra: {
    icon: '🍑',
    label: 'Biodra',
    unit: 'cm',
    color: '#ec4899',
    gradientStart: 'rgba(236, 72, 153, 0.3)',
    gradientEnd: 'rgba(236, 72, 153, 0.05)'
  },
  biceps: {
    icon: '💪',
    label: 'Biceps',
    unit: 'cm',
    color: '#06b6d4',
    gradientStart: 'rgba(6, 182, 212, 0.3)',
    gradientEnd: 'rgba(6, 182, 212, 0.05)'
  },
  udo: {
    icon: '🦵',
    label: 'Udo',
    unit: 'cm',
    color: '#10b981',
    gradientStart: 'rgba(16, 185, 129, 0.3)',
    gradientEnd: 'rgba(16, 185, 129, 0.05)'
  },
  lydka: {
    icon: '🦵',
    label: 'Łydka',
    unit: 'cm',
    color: '#84cc16',
    gradientStart: 'rgba(132, 204, 22, 0.3)',
    gradientEnd: 'rgba(132, 204, 22, 0.05)'
  }
}

// Computed properties
const isTrainer = computed(() => store.user?.role === 'trainer')

const weightChangeDisplay = computed(() => {
  if (!userStats.value?.weight_change) return 'Brak'
  const change = userStats.value.weight_change
  return `${change > 0 ? '+' : ''}${change.toFixed(1)}`
})

const availableChartTypes = computed(() => {
  if (!enhancedChartData.value) return []

  return Object.keys(chartTypes).map(key => {
    const data = enhancedChartData.value[key] || []
    return {
      key,
      ...chartTypes[key],
      hasData: data.length > 0,
      dataCount: data.length
    }
  }).filter(chart => chart.hasData)
})

// Helper functions
const showError = (message) => {
  errorMessage.value = message
  setTimeout(() => errorMessage.value = '', 5000)
}

const showSuccess = (message) => {
  successMessage.value = message
  setTimeout(() => successMessage.value = '', 3000)
}

const handleApiError = (error, context) => {
  console.error(`Błąd ${context}:`, error)
  showError(`Błąd ${context}. Spróbuj ponownie.`)
}

const getCurrentChartTitle = () => {
  const chartType = chartTypes[activeChart.value]
  return chartType ? `${chartType.label} (${chartType.unit})` : 'Wykres'
}

const getCurrentChartStats = () => {
  if (!enhancedChartData.value || !enhancedChartData.value[activeChart.value]) {
    return { latest: null, change: null }
  }

  const data = enhancedChartData.value[activeChart.value]
  if (data.length < 1) return { latest: null, change: null }

  const latest = data[data.length - 1]
  const chartType = chartTypes[activeChart.value]
  const latestValue = `${latest.value.toFixed(1)} ${chartType.unit}`

  let change = null
  if (data.length >= 2) {
    const previous = data[data.length - 2]
    const diff = latest.value - previous.value
    const sign = diff > 0 ? '+' : ''
    change = `${sign}${diff.toFixed(1)} ${chartType.unit}`
  }

  return { latest: latestValue, change }
}

const getChangeClass = (change) => {
  if (!change) return ''
  const isPositive = change.startsWith('+')
  const isWeight = activeChart.value === 'weight'

  // For weight, negative change is good (weight loss)
  // For measurements, it depends on the goal
  if (isWeight) {
    return isPositive ? 'change-negative' : 'change-positive'
  } else {
    return isPositive ? 'change-positive' : 'change-negative'
  }
}

// API call helper with proper base URL
const apiCall = async (url, options = {}) => {
  if (!store.token) {
    throw new Error("No authentication token available. Please log in.");
  }

  const headers = {
    'Authorization': `Bearer ${store.token}`,
    ...options.headers,
  };

  const response = await fetch(`http://localhost:8080${url}`, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    store.logout();
    throw new Error("Session expired. Please log in again.");
  }

  return response;
}

// Enhanced Chart Data Fetching
const fetchEnhancedChartData = async (traineeId = null) => {
  try {
    const url = traineeId
        ? `/reports/trainer/trainee/${traineeId}`
        : '/reports/my'

    const response = await apiCall(url)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    const reports = await response.json()

    // Process reports into chart data
    const chartData = {}

    // Initialize all chart types
    Object.keys(chartTypes).forEach(key => {
      chartData[key] = []
    })

    // Process each report
    reports.forEach(report => {
      const date = new Date(report.created_at).toLocaleDateString('pl-PL')

      // Basic measurements
      if (report.weight) {
        chartData.weight.push({ date, value: report.weight })
      }
      if (report.body_fat) {
        chartData.body_fat.push({ date, value: report.body_fat })
      }

      // Body measurements from JSON
      if (report.measurements) {
        try {
          const measurements = typeof report.measurements === 'string'
              ? JSON.parse(report.measurements)
              : report.measurements

          Object.keys(measurements).forEach(key => {
            if (chartData[key] && measurements[key]) {
              chartData[key].push({ date, value: measurements[key] })
            }
          })
        } catch (e) {
          console.warn('Failed to parse measurements:', e)
        }
      }
    })

    // Sort all data by date
    Object.keys(chartData).forEach(key => {
      chartData[key].sort((a, b) => new Date(a.date) - new Date(b.date))
    })

    enhancedChartData.value = chartData

    // Set default active chart to first available
    const firstAvailable = Object.keys(chartData).find(key => chartData[key].length > 0)
    if (firstAvailable && !chartData[activeChart.value]?.length) {
      activeChart.value = firstAvailable
    }

    await nextTick()
    drawEnhancedChart()

  } catch (error) {
    handleApiError(error, 'pobierania danych wykresu')
    enhancedChartData.value = {}
  }
}

// Enhanced Area Chart Drawing
const drawEnhancedChart = () => {
  const canvas = isTrainer.value && selectedTraineeId.value
      ? trainerChart.value
      : studentChart.value

  if (!canvas || !enhancedChartData.value || !enhancedChartData.value[activeChart.value]) return

  const data = enhancedChartData.value[activeChart.value]
  if (!data.length) return

  const ctx = canvas.getContext('2d')
  const width = canvas.width
  const height = canvas.height
  const chartType = chartTypes[activeChart.value]

  // Clear canvas
  ctx.clearRect(0, 0, width, height)

  // Set up margins
  const margin = { top: 40, right: 40, bottom: 60, left: 60 }
  const chartWidth = width - margin.left - margin.right
  const chartHeight = height - margin.top - margin.bottom

  // Calculate scales
  const values = data.map(d => d.value)
  const minVal = Math.min(...values)
  const maxVal = Math.max(...values)
  const range = maxVal - minVal || 1
  const padding = range * 0.1 // 10% padding
  const yMin = minVal - padding
  const yMax = maxVal + padding
  const yRange = yMax - yMin

  // Create gradient
  const gradient = ctx.createLinearGradient(0, margin.top, 0, height - margin.bottom)
  gradient.addColorStop(0, chartType.gradientStart)
  gradient.addColorStop(1, chartType.gradientEnd)

  // Draw grid lines
  ctx.strokeStyle = '#e5e7eb'
  ctx.lineWidth = 1

  // Horizontal grid lines
  for (let i = 0; i <= 5; i++) {
    const y = margin.top + (chartHeight / 5) * i
    ctx.beginPath()
    ctx.moveTo(margin.left, y)
    ctx.lineTo(width - margin.right, y)
    ctx.stroke()

    // Y-axis labels
    const value = yMax - (yRange / 5) * i
    ctx.fillStyle = '#6b7280'
    ctx.font = '12px Arial'
    ctx.textAlign = 'right'
    ctx.fillText(value.toFixed(1), margin.left - 10, y + 4)
  }

  // Vertical grid lines (for dates)
  const xStep = chartWidth / Math.max(1, data.length - 1)
  for (let i = 0; i < data.length; i++) {
    const x = margin.left + xStep * i
    ctx.beginPath()
    ctx.moveTo(x, margin.top)
    ctx.lineTo(x, height - margin.bottom)
    ctx.stroke()
  }

  // Draw area fill
  ctx.fillStyle = gradient
  ctx.beginPath()

  // Start from bottom left
  ctx.moveTo(margin.left, height - margin.bottom)

  // Draw the line
  data.forEach((point, index) => {
    const x = margin.left + (index / Math.max(1, data.length - 1)) * chartWidth
    const y = margin.top + (yMax - point.value) / yRange * chartHeight

    if (index === 0) {
      ctx.lineTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })

  // Close the path to bottom right
  const lastX = margin.left + chartWidth
  ctx.lineTo(lastX, height - margin.bottom)
  ctx.closePath()
  ctx.fill()

  // Draw the line on top
  ctx.strokeStyle = chartType.color
  ctx.lineWidth = 3
  ctx.beginPath()

  data.forEach((point, index) => {
    const x = margin.left + (index / Math.max(1, data.length - 1)) * chartWidth
    const y = margin.top + (yMax - point.value) / yRange * chartHeight

    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()

  // Draw points
  ctx.fillStyle = chartType.color
  data.forEach((point, index) => {
    const x = margin.left + (index / Math.max(1, data.length - 1)) * chartWidth
    const y = margin.top + (yMax - point.value) / yRange * chartHeight

    ctx.beginPath()
    ctx.arc(x, y, 5, 0, 2 * Math.PI)
    ctx.fill()

    // Add white border to points
    ctx.strokeStyle = 'white'
    ctx.lineWidth = 2
    ctx.stroke()
    ctx.strokeStyle = chartType.color
    ctx.lineWidth = 3
  })

  // Draw X-axis labels (dates)
  ctx.fillStyle = '#6b7280'
  ctx.font = '11px Arial'
  ctx.textAlign = 'center'

  // Show max 8 date labels to avoid overcrowding
  const labelStep = Math.ceil(data.length / 8)
  data.forEach((point, index) => {
    if (index % labelStep === 0 || index === data.length - 1) {
      const x = margin.left + (index / Math.max(1, data.length - 1)) * chartWidth
      ctx.save()
      ctx.translate(x, height - margin.bottom + 20)
      ctx.rotate(-Math.PI / 4)
      ctx.fillText(point.date, 0, 0)
      ctx.restore()
    }
  })

  // Draw title
  ctx.fillStyle = '#1f2937'
  ctx.font = 'bold 16px Arial'
  ctx.textAlign = 'center'
  ctx.fillText(getCurrentChartTitle(), width / 2, 25)

  // Draw Y-axis label
  ctx.save()
  ctx.translate(20, height / 2)
  ctx.rotate(-Math.PI / 2)
  ctx.fillStyle = '#6b7280'
  ctx.font = '12px Arial'
  ctx.textAlign = 'center'
  ctx.fillText(chartType.unit, 0, 0)
  ctx.restore()
}

// API calls with proper error handling
const fetchUserStats = async () => {
  try {
    const response = await apiCall('/reports/my/stats')
    if (!response.ok) {
      const errorText = await response.text()
      console.error('Backend error:', errorText)
      throw new Error(`HTTP ${response.status}`)
    }
    const data = await response.json()
    userStats.value = data
    console.log('User stats loaded:', data)
  } catch (error) {
    handleApiError(error, 'pobierania statystyk')
    userStats.value = {
      total_reports: 0,
      latest_weight: 0,
      weight_change: 0,
      avg_body_fat: 0
    }
  }
}

const fetchTrainerStats = async () => {
  loadingTrainerStats.value = true
  try {
    const response = await apiCall('/reports/trainer/stats')
    if (!response.ok) {
      const errorText = await response.text()
      console.error('Backend error:', errorText)
      throw new Error(`HTTP ${response.status}`)
    }
    const data = await response.json()
    trainerStats.value = data
    console.log('Trainer stats loaded:', data)
    await fetchAllTraineeReports()
  } catch (error) {
    handleApiError(error, 'pobierania statystyk trenera')
    trainerStats.value = {
      trainee_stats: [],
      total_trainees: 0,
      total_reports: 0
    }
  } finally {
    loadingTrainerStats.value = false
  }
}

const fetchAllTraineeReports = async () => {
  if (!trainerStats.value?.trainee_stats?.length) {
    console.log('No trainee stats available')
    allTraineeReports.value = []
    return
  }

  loadingAllReports.value = true
  try {
    const allReports = []

    for (const trainee of trainerStats.value.trainee_stats) {
      try {
        const response = await apiCall(`/reports/trainer/trainee/${trainee.trainee_id}`)

        if (response.ok) {
          const reports = await response.json()

          if (Array.isArray(reports)) {
            allReports.push(...reports)
          } else if (reports && typeof reports === 'object' && reports.data && Array.isArray(reports.data)) {
            allReports.push(...reports.data)
          } else {
            console.warn(`Invalid reports data for trainee ${trainee.trainee_id}:`, reports)
          }
        } else {
          console.warn(`Failed to fetch reports for trainee ${trainee.trainee_id}:`, response.status)
        }
      } catch (error) {
        console.error(`Error fetching reports for trainee ${trainee.trainee_id}:`, error)
      }
    }

    allTraineeReports.value = allReports.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    console.log('All trainee reports loaded:', allTraineeReports.value.length)

  } catch (error) {
    handleApiError(error, 'pobierania wszystkich raportów')
    allTraineeReports.value = []
  } finally {
    loadingAllReports.value = false
  }
}

const fetchMyReports = async () => {
  try {
    const response = await apiCall('/reports/my')
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }
    myReports.value = await response.json()
  } catch (error) {
    handleApiError(error, 'pobierania raportów')
    myReports.value = []
  }
}

const fetchTraineeReports = async (traineeId) => {
  try {
    const response = await apiCall(`/reports/trainer/trainee/${traineeId}`)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }
    traineeReports.value = await response.json()
  } catch (error) {
    handleApiError(error, 'pobierania raportów podopiecznego')
    traineeReports.value = []
  }
}

// Metody
const selectTrainee = async (trainee) => {
  selectedTraineeId.value = trainee.trainee_id
  selectedTrainee.value = trainee
  loadingTraineeData.value = true

  try {
    await Promise.all([
      fetchTraineeReports(trainee.trainee_id),
      fetchEnhancedChartData(trainee.trainee_id)
    ])
  } finally {
    loadingTraineeData.value = false
  }
}

const deselectTrainee = () => {
  selectedTraineeId.value = null
  selectedTrainee.value = null
  traineeReports.value = []
  enhancedChartData.value = null
}

const showReportDetails = (report) => {
  selectedReport.value = report
}

const createReport = async () => {
  if (!newReport.value.weight) {
    showError('Waga jest wymagana')
    return
  }

  loading.value = true
  try {
    const filteredMeasurements = {}
    Object.keys(measurements.value).forEach(key => {
      if (measurements.value[key] !== null && measurements.value[key] !== '') {
        filteredMeasurements[key] = measurements.value[key]
      }
    })

    const reportData = {
      weight: parseFloat(newReport.value.weight),
      body_fat: newReport.value.body_fat ? parseFloat(newReport.value.body_fat) : null,
      notes: newReport.value.notes || '',
      measurements: Object.keys(filteredMeasurements).length > 0
          ? JSON.stringify(filteredMeasurements)
          : null
    };

    console.log('Sending report data:', reportData)

    const response = await apiCall('/reports/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(reportData)
    })

    console.log('Response status:', response.status)

    if (!response.ok) {
      const errorText = await response.text()
      console.error('Backend error response:', errorText)
      throw new Error(`HTTP ${response.status}: ${errorText}`)
    }

    const result = await response.json()
    console.log('Success result:', result)

    showCreateModal.value = false
    newReport.value = {
      weight: null,
      body_fat: null,
      notes: ''
    }
    measurements.value = {
      klatka: null,
      pas: null,
      biodra: null,
      biceps: null,
      udo: null,
      lydka: null
    }

    showSuccess('Raport został dodany pomyślnie!')

    await Promise.all([
      fetchMyReports(),
      fetchUserStats(),
      fetchEnhancedChartData()
    ])
  } catch (error) {
    console.error('Full error in createReport:', error)
    showError(`Błąd tworzenia raportu: ${error.message}`)
  } finally {
    loading.value = false
  }
}

const getTraineeName = (clientId) => {
  const trainee = trainerStats.value?.trainee_stats?.find(t => t.trainee_id === clientId)
  return trainee?.trainee_name || 'Nieznany'
}

const getWeightChange = (trainee) => {
  if (!trainee?.weight_change) return 'Brak'
  const change = trainee.weight_change
  return `${change > 0 ? '+' : ''}${change.toFixed(1)}`
}

const parseMeasurements = (measurementsStr) => {
  try {
    return JSON.parse(measurementsStr || '{}')
  } catch {
    return {}
  }
}

const translateMeasurement = (key) => {
  const translations = {
    klatka: 'Klatka piersiowa',
    pas: 'Pas',
    biodra: 'Biodra',
    biceps: 'Biceps',
    udo: 'Udo',
    lydka: 'Łydka'
  }
  return translations[key] || key
}

const formatDate = (dateString) => {
  if (!dateString) return 'Brak daty'
  try {
    return new Date(dateString).toLocaleDateString('pl-PL')
  } catch {
    return 'Nieprawidłowa data'
  }
}

const formatDateTime = (dateString) => {
  if (!dateString) return 'Brak daty'
  try {
    return new Date(dateString).toLocaleString('pl-PL')
  } catch {
    return 'Nieprawidłowa data'
  }
}

// Watch for chart changes
watch(() => activeChart.value, () => {
  nextTick(() => {
    drawEnhancedChart()
  })
})

// Inicjalizacja
onMounted(async () => {
  console.log('Component mounted, user role:', store.user?.role)
  if (isTrainer.value) {
    await fetchTrainerStats()
  } else {
    loadingUserData.value = true
    try {
      await Promise.all([
        fetchUserStats(),
        fetchMyReports(),
        fetchEnhancedChartData()
      ])
    } finally {
      loadingUserData.value = false
    }
  }
})
</script>

<style scoped>
/* Existing styles remain the same */
.reports-view {
  background: white;
  max-width: 95%;
  margin: 2rem;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  border-radius: 10px;
}

/* Enhanced Chart Styles */
.chart-type-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.chart-type-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  background: white;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
  color: #4a5568;
  min-width: 120px;
  justify-content: center;
}

.chart-type-btn:hover:not(:disabled) {
  border-color: #4f46e5;
  color: #4f46e5;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
}

.chart-type-btn.active {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  border-color: #4f46e5;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 70, 229, 0.3);
}

.chart-type-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: #f7fafc;
  color: #a0aec0;
}

.chart-icon {
  font-size: 18px;
}

.chart-label {
  font-weight: 600;
}

.data-count {
  font-size: 12px;
  opacity: 0.8;
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
}

.chart-type-btn.active .data-count {
  background: rgba(255, 255, 255, 0.3);
}

.enhanced-chart-container {
  background: white;
  border-radius: 16px;
  padding: 24px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.05);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e2e8f0;
}

.chart-header h4 {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
}

.chart-stats {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color:black;
}

.chart-stats span {
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 8px;
  font-weight: 500;
}

.change-positive {
  background: #d1fae5 !important;
  color: #065f46 !important;
}

.change-negative {
  background: #fee2e2 !important;
  color: #991b1b !important;
}

/* Loading states */
.loading {
  text-align: center;
  padding: 40px;
  color: #666;
  font-style: italic;
}

/* Error and Success toasts */
.error-toast, .success-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 8px;
  color: white;
  z-index: 1000;
  cursor: pointer;
  animation: slideInRight 0.3s ease;
}

.error-toast {
  background-color: #e53e3e;
}

.success-toast {
  background-color: #38a169;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* Empty states */
.empty-state {
  text-align: center;
  padding: 40px;
  color: #666;
  background: #f7fafc;
  border-radius: 12px;
  border: 2px dashed #cbd5e0;
}

/* Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px 0;
  border-bottom: 1px solid #e2e8f0;
}

.header-content h1 {
  margin: 0 0 8px 0;
  color: #2d3748;
  font-size: 28px;
  font-weight: bold;
}

.header-content p {
  margin: 0;
  color: #718096;
  font-size: 16px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.btn-add-report {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 14px 0 rgba(79, 70, 229, 0.3);
}

.btn-add-report:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px 0 rgba(79, 70, 229, 0.4);
}

/* Trainer Layout */
.trainer-layout {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 30px;
  align-items: start;
}

.trainees-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

.main-content {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

/* Section Headers */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e2e8f0;
}

.section-title {
  margin: 0;
  color: #2d3748;
  font-size: 20px;
  font-weight: 600;
}

.btn-back {
  background: #e2e8f0;
  color: #4a5568;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-back:hover {
  background: #cbd5e0;
}

/* Trainee Cards */
.trainees-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.trainee-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  border: 2px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
  background: #f7fafc;
}

.trainee-card:hover {
  border-color: #4f46e5;
  background: #edf2f7;
  transform: translateY(-1px);
}

.trainee-card.active {
  border-color: #4f46e5;
  background: linear-gradient(135deg, #eef2ff, #e0e7ff);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
}

.trainee-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 20px;
}

.trainee-info {
  flex: 1;
}

.trainee-info h4 {
  margin: 0 0 8px 0;
  color: #2d3748;
  font-size: 16px;
  font-weight: 600;
}

.trainee-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item {
  font-size: 13px;
  color: #718096;
}

.select-arrow {
  color: #a0aec0;
  font-size: 18px;
  transition: transform 0.2s;
}

.trainee-card:hover .select-arrow {
  transform: translateX(4px);
}

/* Dashboard Cards */
.dashboard-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.dashboard-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #f7fafc, #edf2f7);
  padding: 20px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  transition: transform 0.2s;
}

.dashboard-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.card-icon {
  font-size: 24px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-content h3 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: bold;
  color: #2d3748;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  color: #718096;
}

/* Reports Grid */
.reports-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.report-card {
  background: #f7fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.2s;
}

.report-card.clickable {
  cursor: pointer;
}
/* Report Card Hover and Interactions */
.report-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border-color: #4f46e5;
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.report-date {
  font-weight: 600;
  color: #4a5568;
}

.report-author {
  font-size: 14px;
  color: #718096;
}

.report-summary {
  margin-bottom: 8px;
}

.summary-row {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}

.summary-row span {
  font-size: 14px;
  color: #4a5568;
  font-weight: 500;
}

.report-notes {
  font-size: 14px;
  color: #718096;
  margin: 8px 0 0 0;
  line-height: 1.4;
}

.click-hint {
  text-align: center;
  font-size: 12px;
  color: #a0aec0;
  margin-top: 8px;
  font-style: italic;
}

/* Charts Section */
.charts-section {
  margin: 30px 0;
}

/* Timeline Styles */
.reports-timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.timeline-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: #f7fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.timeline-item:hover {
  background: #edf2f7;
  border-color: #4f46e5;
  transform: translateY(-1px);
}

.timeline-date {
  flex-shrink: 0;
  width: 100px;
  font-weight: 600;
  color: #4a5568;
  font-size: 14px;
}

.timeline-content {
  flex: 1;
}

.timeline-data {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}

.timeline-data span {
  font-size: 14px;
  color: #4a5568;
  font-weight: 500;
}

.timeline-content p {
  margin: 8px 0 0 0;
  color: #718096;
  font-size: 14px;
  line-height: 1.4;
}

.added-by {
  color: #a0aec0;
  font-size: 12px;
  font-style: italic;
  margin-top: 4px;
  display: block;
}

/* Recent Reports Section */
.recent-reports {
  margin-top: 30px;
}
.no-scroll {
  overflow: hidden !important;
  height: 100vh !important;
}
/* Modals */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;

}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-large {
  max-width: 700px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  color: #2d3748;
  font-size: 20px;
  font-weight: 600;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #a0aec0;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: color 0.2s;
}

.btn-close:hover {
  color: #718096;
}

.modal-body {
  padding: 24px;
}

.modal-navigation {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid #e2e8f0;
}

/* Form Elements */
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 20px;
}

.measurements-form {
  margin: 20px 0;
}

.measurements-form h4 {
  margin: 0 0 16px 0;
  color: #2d3748;
  font-size: 16px;
  font-weight: 600;
}

.measurements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #2d3748;
  font-size: 14px;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 12px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #4f46e5;
}

.form-textarea {
  min-height: 80px;
  resize: vertical;
}

/* Buttons */
.btn-cancel, .btn-create {
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-cancel {
  background: #e2e8f0;
  color: #4a5568;
}

.btn-cancel:hover {
  background: #cbd5e0;
}

.btn-create {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
}

.btn-create:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.btn-create:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Report Details */
.report-details {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section h4 {
  margin: 0 0 16px 0;
  color: #2d3748;
  font-size: 16px;
  font-weight: 600;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  background: #f7fafc;
  border-radius: 8px;
}

.detail-label {
  font-weight: 600;
  color: #4a5568;
}

.detail-value {
  color: #2d3748;
}

.measurement-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f7fafc;
  border-radius: 6px;
  font-size: 14px;
}

.measurement-label {
  font-weight: 500;
  color: #4a5568;
  text-transform: capitalize;
}

.measurement-value {
  color: #2d3748;
  font-weight: 600;
}

.notes-content {
  background: #f7fafc;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  line-height: 1.6;
  color: #4a5568;
}

/* Student View Specific */
.student-view {
  width: 100%;
}

.my-stats {
  width: 100%;
}

.my-reports-history {
  margin-top: 30px;
}

/* Trainer View Specific */
.trainer-view {
  width: 100%;
}

.trainee-details {
  width: 100%;
}

.reports-list {
  width: 100%;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .trainer-layout {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .trainees-section {
    order: 2;
  }

  .main-content {
    order: 1;
  }
}

@media (max-width: 768px) {
  .reports-view {
    margin: 1rem;
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .header-actions {
    justify-content: center;
  }

  .dashboard-cards {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 16px;
  }

  .reports-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .measurements-grid {
    grid-template-columns: 1fr;
  }

  .modal-navigation {
    flex-direction: column;
    gap: 8px;
  }

  .chart-type-selector {
    justify-content: center;
  }

  .chart-type-btn {
    min-width: 100px;
    font-size: 14px;
  }

  .timeline-item {
    flex-direction: column;
    gap: 8px;
  }

  .timeline-date {
    width: auto;
    font-size: 12px;
  }

  .trainer-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .chart-type-selector {
    flex-direction: column;
    align-items: stretch;
  }

  .chart-type-btn {
    min-width: auto;
    width: 100%;
  }

  .dashboard-cards {
    grid-template-columns: 1fr;
  }

  .modal-content {
    margin: 10px;
    max-width: calc(100vw - 20px);
  }
}
</style>