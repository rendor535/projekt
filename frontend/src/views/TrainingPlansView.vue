<template>
  <div class="training-plans-view">
    <div class="page-header">
      <div class="header-content">
        <h1>🏋️ Plany Treningowe</h1>
        <p v-if="store.user?.role === 'trainer'">
          Twórz i zarządzaj planami treningowymi dla swoich podopiecznych
        </p>
        <p v-else>
          Twoje plany treningowe od trenera
        </p>
      </div>
      <button
          v-if="store.user?.role === 'trainer'"
          @click="showCreateModal = true"
          class="btn-create-plan"
      >
        ➕ Nowy Plan
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>Ładowanie planów treningowych...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <span class="error-icon">⚠️</span>
      <p>{{ error }}</p>
      <button @click="fetchPlans" class="btn-retry">Spróbuj ponownie</button>
    </div>

    <!-- Main Content -->
    <div v-else class="plans-content">
      <!-- Dashboard Cards -->
      <div class="dashboard-cards">
        <div class="dashboard-card">
          <div class="card-icon">📊</div>
          <div class="card-content">
            <h3>{{ plans.length }}</h3>
            <p>Aktywnych planów</p>
          </div>
        </div>
        <div class="dashboard-card">
          <div class="card-icon">📅</div>
          <div class="card-content">
            <h3>{{ upcomingWorkouts }}</h3>
            <p>Nadchodzących treningów</p>
          </div>
        </div>
        <div class="dashboard-card">
          <div class="card-icon">👥</div>
          <div class="card-content">
            <h3>{{ store.user?.role === 'trainer' ? trainees.length : (myTrainer ? 1 : 0) }}</h3>
            <p>{{ store.user?.role === 'trainer' ? 'Podopiecznych' : 'Trener' }}</p>
          </div>
        </div>
      </div>

      <!-- Calendar View -->
      <div class="calendar-section">
        <div class="section-header">
          <h3 class="section-title">
            📅 Kalendarz Treningów
          </h3>
          <div class="calendar-controls">
            <button @click="previousMonth" class="calendar-nav-btn">‹</button>
            <h4 class="calendar-month">{{ currentMonthYear }}</h4>
            <button @click="nextMonth" class="calendar-nav-btn">›</button>
          </div>
        </div>

        <div class="calendar-container">
          <div class="calendar-weekdays">
            <div class="weekday" v-for="day in ['Pon', 'Wt', 'Śr', 'Czw', 'Pt', 'Sob', 'Nd']" :key="day">
              {{ day }}
            </div>
          </div>
          <div class="calendar-grid">
            <div v-for="day in calendarDays" :key="day.date"
                 class="calendar-day"
                 :class="{
                   'has-events': day.events.length > 0,
                   'today': day.isToday,
                   'other-month': !day.isCurrentMonth
                 }">
              <div class="day-number">{{ day.dayNumber }}</div>
              <div class="events-container">
                <div v-for="event in day.events.slice(0, 2)" :key="event.id"
                     class="calendar-event"
                     @click="viewPlanDetails(event.extendedProps.planId)">
                  {{ event.title }}
                </div>
                <div v-if="day.events.length > 2" class="more-events">
                  +{{ day.events.length - 2 }} więcej
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Plans List -->
      <div class="plans-section">
        <div class="section-header">
          <h3 class="section-title">
            📋 Lista Planów
          </h3>
          <div class="filter-controls">
            <select v-model="filterStatus" class="filter-select">
              <option value="">Wszystkie</option>
              <option value="active">Aktywne</option>
              <option value="completed">Zakończone</option>
            </select>
          </div>
        </div>

        <div v-if="filteredPlans.length === 0" class="empty-state">
          <div class="empty-illustration">
            <div class="empty-icon">🏋️‍♂️</div>
            <h3>{{ plans.length === 0 ? 'Brak planów' : 'Brak planów spełniających kryteria' }}</h3>
            <p v-if="store.user?.role === 'trainer' && plans.length === 0">
              Utwórz pierwszy plan treningowy dla swoich podopiecznych
            </p>
            <p v-else-if="plans.length === 0">
              Skontaktuj się ze swoim trenerem aby otrzymać plan treningowy
            </p>
            <p v-else>
              Spróbuj zmienić filtry aby zobaczyć więcej planów
            </p>
            <button
                v-if="store.user?.role === 'trainer' && plans.length === 0"
                @click="showCreateModal = true"
                class="btn-create-first">
              Utwórz pierwszy plan
            </button>
          </div>
        </div>

        <div v-else class="plans-grid">
          <div v-for="plan in filteredPlans" :key="plan.id" class="plan-card">
            <div class="plan-header">
              <div class="plan-title-section">
                <h4 class="plan-name">{{ plan.name }}</h4>
                <div class="plan-status" :class="getPlanStatus(plan)">
                  {{ getPlanStatusText(plan) }}
                </div>
              </div>
              <div class="plan-actions">
                <button @click="viewPlanDetails(plan.id)" class="btn-action btn-view" title="Zobacz szczegóły">
                  👁️
                </button>
                <button
                    v-if="store.user?.role === 'trainer'"
                    @click="editPlan(plan.id)"
                    class="btn-action btn-edit"
                    title="Edytuj plan"
                >
                  ✏️
                </button>
                <button
                    v-if="store.user?.role === 'trainer'"
                    @click="deletePlan(plan.id)"
                    class="btn-action btn-danger"
                    title="Usuń plan"
                >
                  🗑️
                </button>
              </div>
            </div>

            <div class="plan-info">
              <p class="plan-description" v-if="plan.description">{{ plan.description }}</p>

              <div class="plan-meta">
                <div class="meta-item">
                  <span class="meta-icon">👤</span>
                  <span>{{ store.user?.role === 'trainer' ? plan.client_name : plan.trainer_name }}</span>
                </div>

                <div class="meta-item" v-if="plan.start_date">
                  <span class="meta-icon">📅</span>
                  <span>{{ formatDateRange(plan.start_date, plan.end_date) }}</span>
                </div>

                <div class="meta-item">
                  <span class="meta-icon">🏋️</span>
                  <span>{{ getExerciseCount(plan) }} ćwiczeń</span>
                </div>

                <div class="meta-item">
                  <span class="meta-icon">📊</span>
                  <span>{{ getProgressText(plan) }}</span>
                </div>
              </div>

              <div class="plan-progress">
                <div class="progress-bar">
                  <div class="progress-fill" :style="{
                    width: getProgressPercentage(plan) + '%'
                  }"></div>
                </div>
                <span class="progress-text">{{ getProgressPercentage(plan) }}% ukończone</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Plan Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal-content modal-large" @click.stop>
        <div class="modal-header">
          <h3>🏋️ Nowy Plan Treningowy</h3>
          <button @click="closeCreateModal" class="btn-close">✕</button>
        </div>

        <!-- Progress Steps -->
        <div class="creation-steps">
          <div class="step" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">
            <div class="step-number">1</div>
            <span>Podstawowe info</span>
          </div>
          <div class="step" :class="{ active: currentStep >= 2, completed: currentStep > 2 }">
            <div class="step-number">2</div>
            <span>Terminy</span>
          </div>
          <div class="step" :class="{ active: currentStep >= 3, completed: currentStep > 3 }">
            <div class="step-number">3</div>
            <span>Dni treningowe</span>
          </div>
          <div class="step" :class="{ active: currentStep >= 4 }">
            <div class="step-number">4</div>
            <span>Ćwiczenia</span>
          </div>
        </div>

        <!-- Step 1: Basic Info -->
        <div v-if="currentStep === 1" class="step-content">
          <h4>📝 Podstawowe informacje</h4>

          <div class="form-group">
            <label>Nazwa planu:</label>
            <input v-model="newPlan.name" type="text" required class="form-input"
                   placeholder="np. Plan Redukcyjny - Maj 2025">
          </div>

          <div class="form-group">
            <label>Opis planu:</label>
            <textarea v-model="newPlan.description" class="form-textarea"
                      placeholder="Opisz cel i szczegóły planu treningowego..."></textarea>
          </div>

          <div class="form-group">
            <label>Podopieczny:</label>
            <select v-model="newPlan.client_id" required class="form-select">
              <option value="">Wybierz podopiecznego</option>
              <option v-for="trainee in trainees" :key="trainee.id" :value="trainee.id">
                {{ trainee.username }} {{ trainee.full_name ? `(${trainee.full_name})` : '' }}
              </option>
            </select>
          </div>
        </div>

        <!-- Step 2: Dates -->
        <div v-if="currentStep === 2" class="step-content">
          <h4>📅 Terminy planu</h4>

          <div class="form-row">
            <div class="form-group">
              <label>Data rozpoczęcia:</label>
              <input v-model="newPlan.start_date" type="date" required class="form-input"
                     :min="todayDate" @change="updateAvailableDays">
            </div>

            <div class="form-group">
              <label>Data zakończenia:</label>
              <input v-model="newPlan.end_date" type="date" class="form-input"
                     :min="newPlan.start_date || todayDate">
            </div>
          </div>

          <div v-if="planDuration" class="duration-info">
            <span class="duration-badge">
              📊 Długość planu: {{ planDuration }} dni
            </span>
          </div>
        </div>

        <!-- Step 3: Training Days Selection -->
        <div v-if="currentStep === 3" class="step-content">
          <h4>📅 Wybierz dni treningowe</h4>
          <p class="step-description">
            Wybierz które dni tygodnia będą treningowe w ramach tego planu
          </p>

          <div class="weekdays-selector">
            <div v-for="day in dayOptions" :key="day"
                 class="weekday-option"
                 :class="{ selected: isWeekdaySelected(day) }"
                 @click="toggleWeekday(day)">
              <div class="weekday-icon">{{ getWeekdayIcon(day) }}</div>
              <div class="weekday-name">{{ day }}</div>
              <div v-if="isWeekdaySelected(day)" class="weekday-check">✓</div>
            </div>
          </div>

          <div v-if="selectedWeekdays.length > 0" class="selected-summary">
            <strong>Wybrane dni:</strong> {{ selectedWeekdays.join(', ') }}
            <br>
            <small>{{ selectedWeekdays.length }} dni treningowych w tygodniu</small>
          </div>
        </div>

        <!-- Step 4: Exercises -->
        <div v-if="currentStep === 4" class="step-content">
          <h4>🏋️ Ćwiczenia dla każdego dnia</h4>

          <div v-if="newPlan.training_days.length === 0" class="no-days-selected">
            <p>⚠️ Nie wybrano żadnych dni treningowych</p>
            <button @click="currentStep = 3" class="btn-back">Wróć do wyboru dni</button>
          </div>

          <div v-else class="training-days-exercises">
            <div v-for="(day, dayIndex) in newPlan.training_days" :key="dayIndex" class="training-day-exercises">
              <div class="day-header-fancy">
                <div class="day-icon">{{ getWeekdayIcon(day.day_of_week) }}</div>
                <h5>{{ day.day_of_week }}</h5>
                <div class="exercises-count">{{ day.exercises.length }} ćwiczeń</div>
              </div>

              <div class="exercises-container">
                <div v-for="(exercise, exerciseIndex) in day.exercises" :key="exerciseIndex" class="exercise-card">
                  <div class="exercise-main">
                    <select v-model="exercise.exercise_id" class="form-select">
                      <option value="">Wybierz ćwiczenie</option>
                      <option v-for="ex in exercises" :key="ex.id" :value="ex.id">
                        {{ ex.name }}
                      </option>
                    </select>

                    <div class="exercise-params-grid">
                      <div class="param-group">
                        <label>Serie</label>
                        <input v-model.number="exercise.sets" type="number" min="1" max="10" class="param-input">
                      </div>
                      <div class="param-group">
                        <label>Powtórzenia</label>
                        <input v-model.number="exercise.reps" type="number" min="1" max="50" class="param-input">
                      </div>
                      <div class="param-group">
                        <label>Ciężar (kg)</label>
                        <input v-model.number="exercise.weight" type="number" min="0" step="0.5" class="param-input">
                      </div>
                    </div>
                  </div>

                  <div class="exercise-instructions">
                    <textarea v-model="exercise.instructions"
                              placeholder="Dodatkowe instrukcje dla tego ćwiczenia..."
                              class="instructions-textarea"></textarea>
                  </div>

                  <button @click="removeExercise(dayIndex, exerciseIndex)" class="btn-remove-exercise">
                    🗑️ Usuń
                  </button>
                </div>

                <button @click="addExercise(dayIndex)" class="btn-add-exercise">
                  ➕ Dodaj ćwiczenie
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation -->
        <div class="modal-navigation">
          <button v-if="currentStep > 1" @click="currentStep--" class="btn-nav btn-prev">
            ← Poprzedni krok
          </button>

          <button v-if="currentStep < 4" @click="nextStep" class="btn-nav btn-next" :disabled="!canProceedToNextStep">
            Następny krok →
          </button>

          <button v-if="currentStep === 4" @click="createPlan" class="btn-nav btn-create" :disabled="creating || !canCreatePlan">
            {{ creating ? 'Tworzenie...' : '🏋️ Utwórz Plan' }}
          </button>

          <button @click="closeCreateModal" class="btn-nav btn-cancel">
            Anuluj
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Plan Modal -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content modal-large" @click.stop>
        <div class="modal-header">
          <h3>✏️ Edytuj Plan Treningowy</h3>
          <button @click="closeEditModal" class="btn-close">✕</button>
        </div>

        <!-- Progress Steps -->
        <div class="creation-steps">
          <div class="step" :class="{ active: editCurrentStep >= 1, completed: editCurrentStep > 1 }">
            <div class="step-number">1</div>
            <span>Podstawowe info</span>
          </div>
          <div class="step" :class="{ active: editCurrentStep >= 2, completed: editCurrentStep > 2 }">
            <div class="step-number">2</div>
            <span>Terminy</span>
          </div>
          <div class="step" :class="{ active: editCurrentStep >= 3, completed: editCurrentStep > 3 }">
            <div class="step-number">3</div>
            <span>Dni treningowe</span>
          </div>
          <div class="step" :class="{ active: editCurrentStep >= 4 }">
            <div class="step-number">4</div>
            <span>Ćwiczenia</span>
          </div>
        </div>

        <!-- Step 1: Basic Info -->
        <div v-if="editCurrentStep === 1" class="step-content">
          <h4>📝 Podstawowe informacje</h4>

          <div class="form-group">
            <label>Nazwa planu:</label>
            <input v-model="editForm.name" type="text" required class="form-input"
                   placeholder="np. Plan Redukcyjny - Maj 2025">
          </div>

          <div class="form-group">
            <label>Opis planu:</label>
            <textarea v-model="editForm.description" class="form-textarea"
                      placeholder="Opisz cel i szczegóły planu treningowego..."></textarea>
          </div>

          <div class="form-group">
            <label>Podopieczny:</label>
            <select v-model="editForm.client_id" required class="form-select">
              <option value="">Wybierz podopiecznego</option>
              <option v-for="trainee in trainees" :key="trainee.id" :value="trainee.id">
                {{ trainee.username }} {{ trainee.full_name ? `(${trainee.full_name})` : '' }}
              </option>
            </select>
          </div>
        </div>

        <!-- Step 2: Dates -->
        <div v-if="editCurrentStep === 2" class="step-content">
          <h4>📅 Terminy planu</h4>

          <div class="form-row">
            <div class="form-group">
              <label>Data rozpoczęcia:</label>
              <input v-model="editForm.start_date" type="date" required class="form-input"
                     :min="todayDate" @change="updateEditAvailableDays">
            </div>

            <div class="form-group">
              <label>Data zakończenia:</label>
              <input v-model="editForm.end_date" type="date" class="form-input"
                     :min="editForm.start_date || todayDate">
            </div>
          </div>

          <div v-if="editPlanDuration" class="duration-info">
            <span class="duration-badge">
              📊 Długość planu: {{ editPlanDuration }} dni
            </span>
          </div>
        </div>

        <!-- Step 3: Training Days Selection -->
        <div v-if="editCurrentStep === 3" class="step-content">
          <h4>📅 Wybierz dni treningowe</h4>
          <p class="step-description">
            Wybierz które dni tygodnia będą treningowe w ramach tego planu
          </p>

          <div class="weekdays-selector">
            <div v-for="day in dayOptions" :key="day"
                 class="weekday-option"
                 :class="{ selected: isEditWeekdaySelected(day) }"
                 @click="toggleEditWeekday(day)">
              <div class="weekday-icon">{{ getWeekdayIcon(day) }}</div>
              <div class="weekday-name">{{ day }}</div>
              <div v-if="isEditWeekdaySelected(day)" class="weekday-check">✓</div>
            </div>
          </div>

          <div v-if="editSelectedWeekdays.length > 0" class="selected-summary">
            <strong>Wybrane dni:</strong> {{ editSelectedWeekdays.join(', ') }}
            <br>
            <small>{{ editSelectedWeekdays.length }} dni treningowych w tygodniu</small>
          </div>
        </div>

        <!-- Step 4: Exercises -->
        <div v-if="editCurrentStep === 4" class="step-content">
          <h4>🏋️ Ćwiczenia dla każdego dnia</h4>

          <div v-if="editForm.training_days.length === 0" class="no-days-selected">
            <p>⚠️ Nie wybrano żadnych dni treningowych</p>
            <button @click="editCurrentStep = 3" class="btn-back">Wróć do wyboru dni</button>
          </div>

          <div v-else class="training-days-exercises">
            <div v-for="(day, dayIndex) in editForm.training_days" :key="dayIndex" class="training-day-exercises">
              <div class="day-header-fancy">
                <div class="day-icon">{{ getWeekdayIcon(day.day_of_week) }}</div>
                <h5>{{ day.day_of_week }}</h5>
                <div class="exercises-count">{{ day.exercises.length }} ćwiczeń</div>
              </div>

              <div class="exercises-container">
                <div v-for="(exercise, exerciseIndex) in day.exercises" :key="exerciseIndex" class="exercise-card">
                  <div class="exercise-main">
                    <select v-model="exercise.exercise_id" class="form-select">
                      <option value="">Wybierz ćwiczenie</option>
                      <option v-for="ex in exercises" :key="ex.id" :value="ex.id">
                        {{ ex.name }}
                      </option>
                    </select>

                    <div class="exercise-params-grid">
                      <div class="param-group">
                        <label>Serie</label>
                        <input v-model.number="exercise.sets" type="number" min="1" max="10" class="param-input">
                      </div>
                      <div class="param-group">
                        <label>Powtórzenia</label>
                        <input v-model.number="exercise.reps" type="number" min="1" max="50" class="param-input">
                      </div>
                      <div class="param-group">
                        <label>Ciężar (kg)</label>
                        <input v-model.number="exercise.weight" type="number" min="0" step="0.5" class="param-input">
                      </div>
                    </div>
                  </div>

                  <div class="exercise-instructions">
                    <textarea v-model="exercise.instructions"
                              placeholder="Dodatkowe instrukcje dla tego ćwiczenia..."
                              class="instructions-textarea"></textarea>
                  </div>

                  <button @click="removeEditExercise(dayIndex, exerciseIndex)" class="btn-remove-exercise">
                    🗑️ Usuń
                  </button>
                </div>

                <button @click="addEditExercise(dayIndex)" class="btn-add-exercise">
                  ➕ Dodaj ćwiczenie
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation -->
        <div class="modal-navigation">
          <button v-if="editCurrentStep > 1" @click="editCurrentStep--" class="btn-nav btn-prev">
            ← Poprzedni krok
          </button>

          <button v-if="editCurrentStep < 4" @click="nextEditStep" class="btn-nav btn-next" :disabled="!canProceedToNextEditStep">
            Następny krok →
          </button>

          <button v-if="editCurrentStep === 4" @click="updatePlan" class="btn-nav btn-update" :disabled="updating || !canUpdatePlan">
            {{ updating ? 'Aktualizowanie...' : '💾 Zaktualizuj Plan' }}
          </button>

          <button @click="closeEditModal" class="btn-nav btn-cancel">
            Anuluj
          </button>
        </div>
      </div>
    </div>

    <!-- Plan Details Modal -->
    <div v-if="showDetailsModal" class="modal-overlay" @click="closeDetailsModal">
      <div class="modal-content modal-xl" @click.stop>
        <div class="modal-header">
          <div class="plan-header-info">
            <h3>{{ selectedPlan?.name }}</h3>
            <div class="plan-badges">
              <span class="badge" :class="getPlanStatus(selectedPlan)">
                {{ getPlanStatusText(selectedPlan) }}
              </span>
              <span class="badge badge-info">
                {{ getTotalExercisesInPlan(selectedPlan) }} ćwiczeń
              </span>
            </div>
          </div>
          <button @click="closeDetailsModal" class="btn-close">✕</button>
        </div>

        <div v-if="selectedPlan" class="plan-details-enhanced">
          <div class="plan-overview">
            <div class="overview-card">
              <h4>📋 Informacje ogólne</h4>
              <p><strong>Opis:</strong> {{ selectedPlan.description || 'Brak opisu' }}</p>
              <p><strong>Okres:</strong> {{ formatDateRange(selectedPlan.start_date, selectedPlan.end_date) }}</p>
              <p><strong>{{ store.user?.role === 'trainer' ? 'Podopieczny' : 'Trener' }}:</strong>
                {{ store.user?.role === 'trainer' ? selectedPlan.client_name : selectedPlan.trainer_name }}</p>
            </div>

            <div class="overview-card">
              <h4>📊 Postęp</h4>
              <div class="progress-circle">
                <div class="circle-progress" :style="{ '--progress': getProgressPercentage(selectedPlan) }">
                  <span>{{ getProgressPercentage(selectedPlan) }}%</span>
                </div>
              </div>
              <p>{{ getProgressText(selectedPlan) }}</p>
            </div>
          </div>

          <div v-if="selectedPlan.days && selectedPlan.days.length > 0" class="training-schedule">
            <h4>🗓️ Harmonogram treningów</h4>
            <div class="schedule-grid">
              <div v-for="day in selectedPlan.days" :key="day.id" class="schedule-day">
                <div class="day-header-detailed">
                  <div class="day-icon-large">{{ getWeekdayIcon(day.day_of_week) }}</div>
                  <div class="day-info">
                    <h5>{{ day.day_of_week }}</h5>
                    <span class="exercises-count">{{ day.exercises?.length || 0 }} ćwiczeń</span>
                  </div>
                </div>

                <div v-if="day.exercises && day.exercises.length > 0" class="day-exercises">
                  <div v-for="(exercise, index) in day.exercises" :key="exercise.id" class="exercise-item-detailed">
                    <div class="exercise-number">{{ index + 1 }}</div>
                    <div class="exercise-content">
                      <div class="exercise-name-detailed">{{ exercise.exercise_name }}</div>
                      <div class="exercise-specs">
                        <span class="spec-item">{{ exercise.sets }}x{{ exercise.reps }}</span>
                        <span v-if="exercise.weight > 0" class="spec-item">{{ exercise.weight }}kg</span>
                      </div>
                      <div v-if="exercise.instructions" class="exercise-notes">
                        📝 {{ exercise.instructions }}
                      </div>
                    </div>
                  </div>
                </div>
                <div v-else class="no-exercises-day">
                  <span class="empty-day-icon">😴</span>
                  <p>Dzień odpoczynku</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

export default {
  setup() {
    const store = useAuthStore()

    // ========== REACTIVE DATA ==========
    const loading = ref(false)
    const error = ref('')
    const plans = ref([])
    const calendarEvents = ref([])
    const exercises = ref([])
    const trainees = ref([])
    const myTrainer = ref(null)
    const currentDate = ref(new Date())
    const filterStatus = ref('')

    // Create Plan Modal
    const showCreateModal = ref(false)
    const creating = ref(false)
    const currentStep = ref(1)
    const selectedWeekdays = ref([])
    const newPlan = ref({
      name: '',
      description: '',
      client_id: '',
      start_date: '',
      end_date: '',
      training_days: []
    })

    // Edit Plan Modal
    const showEditModal = ref(false)
    const editingPlan = ref(null)
    const updating = ref(false)
    const editCurrentStep = ref(1)
    const editSelectedWeekdays = ref([])
    const editForm = ref({
      id: '',
      name: '',
      description: '',
      client_id: '',
      start_date: '',
      end_date: '',
      training_days: []
    })

    // Plan Details Modal
    const showDetailsModal = ref(false)
    const selectedPlan = ref(null)

    // ========== CONSTANTS ==========
    const dayOptions = [
      'poniedziałek', 'wtorek', 'środa', 'czwartek',
      'piątek', 'sobota', 'niedziela'
    ]

    // ========== COMPUTED PROPERTIES ==========
    const todayDate = computed(() => {
      return new Date().toISOString().split('T')[0]
    })

    const currentMonthYear = computed(() => {
      return currentDate.value.toLocaleDateString('pl-PL', {
        month: 'long',
        year: 'numeric'
      })
    })

    const planDuration = computed(() => {
      if (!newPlan.value.start_date || !newPlan.value.end_date) return null
      const start = new Date(newPlan.value.start_date)
      const end = new Date(newPlan.value.end_date)
      const diffTime = Math.abs(end - start)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    })

    const editPlanDuration = computed(() => {
      if (!editForm.value.start_date || !editForm.value.end_date) return null
      const start = new Date(editForm.value.start_date)
      const end = new Date(editForm.value.end_date)
      const diffTime = Math.abs(end - start)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    })

    const filteredPlans = computed(() => {
      if (!filterStatus.value) return plans.value

      return plans.value.filter(plan => {
        const now = new Date()
        const startDate = new Date(plan.start_date)
        const endDate = plan.end_date ? new Date(plan.end_date) : null

        if (filterStatus.value === 'active') {
          return now >= startDate && (!endDate || now <= endDate)
        } else if (filterStatus.value === 'completed') {
          return endDate && now > endDate
        }
        return true
      })
    })

    const upcomingWorkouts = computed(() => {
      return calendarEvents.value.filter(event => {
        const eventDate = new Date(event.start)
        return eventDate >= new Date()
      }).length
    })

    const canProceedToNextStep = computed(() => {
      switch (currentStep.value) {
        case 1:
          return newPlan.value.name && newPlan.value.client_id
        case 2:
          return newPlan.value.start_date
        case 3:
          return selectedWeekdays.value.length > 0
        default:
          return true
      }
    })

    const canProceedToNextEditStep = computed(() => {
      switch (editCurrentStep.value) {
        case 1:
          return editForm.value.name && editForm.value.client_id
        case 2:
          return editForm.value.start_date
        case 3:
          return editSelectedWeekdays.value.length > 0
        default:
          return true
      }
    })

    const canCreatePlan = computed(() => {
      return newPlan.value.training_days.every(day =>
          day.exercises.length > 0 &&
          day.exercises.every(ex => ex.exercise_id && ex.sets && ex.reps)
      )
    })

    const canUpdatePlan = computed(() => {
      return editForm.value.training_days.every(day =>
          day.exercises.length > 0 &&
          day.exercises.every(ex => ex.exercise_id && ex.sets && ex.reps)
      )
    })

    const calendarDays = computed(() => {
      const year = currentDate.value.getFullYear()
      const month = currentDate.value.getMonth()
      const firstDay = new Date(year, month, 1)
      const lastDay = new Date(year, month + 1, 0)
      const daysInMonth = lastDay.getDate()
      const startingDayOfWeek = firstDay.getDay() === 0 ? 6 : firstDay.getDay() - 1

      const days = []
      const today = new Date()

      // Add days from previous month
      for (let i = startingDayOfWeek - 1; i >= 0; i--) {
        const date = new Date(year, month, -i)
        days.push({
          date: date.toISOString().split('T')[0],
          dayNumber: date.getDate(),
          events: [],
          isToday: false,
          isCurrentMonth: false
        })
      }

      // Add days from current month
      for (let day = 1; day <= daysInMonth; day++) {
        const date = new Date(year, month, day)
        const dateStr = date.toISOString().split('T')[0]
        const isToday = dateStr === today.toISOString().split('T')[0]

        days.push({
          date: dateStr,
          dayNumber: day,
          events: calendarEvents.value.filter(event =>
              event.start && event.start.startsWith(dateStr)
          ),
          isToday,
          isCurrentMonth: true
        })
      }

      return days
    })

    // ========== API FUNCTIONS ==========
    const makeApiCall = async (url, options = {}) => {
      const token = store.token
      if (!token) {
        throw new Error('Brak tokenu autoryzacji')
      }

      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://188.245.239.124:8080/api'}${url}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
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

    // ========== DATA FETCHING ==========
    const fetchPlans = async () => {
      loading.value = true
      error.value = ''

      try {
        const endpoint = store.user?.role === 'trainer' ? '/training-plans/trainer' : '/training-plans/trainee'
        const result = await makeApiCall(endpoint)
        plans.value = result || []
      } catch (err) {
        console.error('❌ ERROR in fetchPlans:', err)
        error.value = err.message || 'Błąd ładowania'
        plans.value = []
      } finally {
        loading.value = false
      }
    }

    const fetchCalendarEvents = async () => {
      try {
        const events = await makeApiCall('/training-plans/calendar')
        calendarEvents.value = events || []
      } catch (err) {
        console.error('❌ Calendar events ERROR:', err)
        calendarEvents.value = []
      }
    }

    const fetchExercises = async () => {
      try {
        const exercisesList = await makeApiCall('/training-plans/exercises')
        exercises.value = exercisesList || []
      } catch (err) {
        console.error('❌ Exercises ERROR:', err)
        exercises.value = []
      }
    }

    const fetchTrainees = async () => {
      if (store.user?.role !== 'trainer') return

      try {
        const traineesList = await makeApiCall('/training-plans/trainees')
        trainees.value = traineesList || []
      } catch (err) {
        console.error('❌ Trainees ERROR:', err)
        trainees.value = []
      }
    }

    const fetchMyTrainer = async () => {
      if (store.user?.role === 'trainee') {
        try {
          const trainer = await makeApiCall('/trainer-requests/my-trainer')
          myTrainer.value = trainer
        } catch (err) {
          myTrainer.value = null
        }
      }
    }

    // ========== EXERCISE COUNT FUNCTIONS ==========
    const getExerciseCount = (plan) => {
      if (!plan || !plan.training_days) return 0

      return plan.training_days.reduce((total, day) => {
        return total + (day.exercises?.length || 0)
      }, 0)
    }

    const getTotalExercisesInPlan = (plan) => {
      if (!plan || !plan.days) return 0

      return plan.days.reduce((total, day) => {
        return total + (day.exercises?.length || 0)
      }, 0)
    }

    // ========== CREATE PLAN FUNCTIONS ==========
    const createPlan = async () => {
      creating.value = true

      try {
        if (!newPlan.value.name || !newPlan.value.client_id || !newPlan.value.start_date) {
          throw new Error('Wypełnij wszystkie wymagane pola')
        }

        if (newPlan.value.training_days.length === 0) {
          throw new Error('Dodaj przynajmniej jeden dzień treningowy')
        }

        await makeApiCall('/training-plans', {
          method: 'POST',
          body: JSON.stringify(newPlan.value)
        })

        await fetchPlans()
        closeCreateModal()

      } catch (err) {
        error.value = err.message || 'Błąd podczas tworzenia planu'
        console.error('Error creating plan:', err)
      } finally {
        creating.value = false
      }
    }

    const closeCreateModal = () => {
      showCreateModal.value = false
      currentStep.value = 1
      selectedWeekdays.value = []
      newPlan.value = {
        name: '',
        description: '',
        client_id: '',
        start_date: '',
        end_date: '',
        training_days: []
      }
    }

    const nextStep = () => {
      if (canProceedToNextStep.value) {
        currentStep.value++
        if (currentStep.value === 4) {
          updateTrainingDaysFromWeekdays()
        }
      }
    }

    const updateAvailableDays = () => {
      selectedWeekdays.value = []
      newPlan.value.training_days = []
    }

    const isWeekdaySelected = (weekday) => {
      return selectedWeekdays.value.includes(weekday)
    }

    const toggleWeekday = (weekday) => {
      const index = selectedWeekdays.value.indexOf(weekday)
      if (index > -1) {
        selectedWeekdays.value.splice(index, 1)
      } else {
        selectedWeekdays.value.push(weekday)
      }
    }

    const updateTrainingDaysFromWeekdays = () => {
      newPlan.value.training_days = selectedWeekdays.value.map(weekday => ({
        day_of_week: weekday,
        exercises: []
      }))
    }

    const addExercise = (dayIndex) => {
      newPlan.value.training_days[dayIndex].exercises.push({
        exercise_id: '',
        sets: 3,
        reps: 10,
        weight: 0,
        instructions: ''
      })
    }

    const removeExercise = (dayIndex, exerciseIndex) => {
      newPlan.value.training_days[dayIndex].exercises.splice(exerciseIndex, 1)
    }

    // ========== EDIT PLAN FUNCTIONS ==========
    const editPlan = async (planId) => {
      try {
        const planDetails = await makeApiCall(`/training-plans/${planId}/details`)

        editingPlan.value = planDetails
        editForm.value = {
          id: planDetails.id,
          name: planDetails.name,
          description: planDetails.description || '',
          client_id: planDetails.client_id,
          start_date: planDetails.start_date,
          end_date: planDetails.end_date || '',
          training_days: planDetails.days ? planDetails.days.map(day => ({
            id: day.id,
            day_of_week: day.day_of_week,
            exercises: day.exercises ? day.exercises.map(ex => ({
              id: ex.id,
              exercise_id: ex.exercise_id,
              exercise_name: ex.exercise_name,
              sets: ex.sets,
              reps: ex.reps,
              weight: ex.weight || 0,
              instructions: ex.instructions || ''
            })) : []
          })) : []
        }

        editSelectedWeekdays.value = planDetails.days ?
            planDetails.days.map(day => day.day_of_week) : []

        showEditModal.value = true
        editCurrentStep.value = 1

      } catch (err) {
        error.value = err.message || 'Błąd podczas pobierania danych planu'
        console.error('Error loading plan for edit:', err)
      }
    }

    const updatePlan = async () => {
      updating.value = true

      try {
        if (!editForm.value.name || !editForm.value.client_id || !editForm.value.start_date) {
          throw new Error('Wypełnij wszystkie wymagane pola')
        }

        if (editForm.value.training_days.length === 0) {
          throw new Error('Dodaj przynajmniej jeden dzień treningowy')
        }

        await makeApiCall(`/training-plans/${editForm.value.id}`, {
          method: 'PUT',
          body: JSON.stringify(editForm.value)
        })

        await fetchPlans()
        closeEditModal()

      } catch (err) {
        error.value = err.message || 'Błąd podczas aktualizacji planu'
        console.error('Error updating plan:', err)
      } finally {
        updating.value = false
      }
    }

    const closeEditModal = () => {
      showEditModal.value = false
      editCurrentStep.value = 1
      editSelectedWeekdays.value = []
      editingPlan.value = null
      editForm.value = {
        id: '',
        name: '',
        description: '',
        client_id: '',
        start_date: '',
        end_date: '',
        training_days: []
      }
    }

    const nextEditStep = () => {
      if (canProceedToNextEditStep.value) {
        editCurrentStep.value++
        if (editCurrentStep.value === 4) {
          updateEditTrainingDaysFromWeekdays()
        }
      }
    }

    const updateEditAvailableDays = () => {
      editSelectedWeekdays.value = []
      editForm.value.training_days = []
    }

    const isEditWeekdaySelected = (weekday) => {
      return editSelectedWeekdays.value.includes(weekday)
    }

    const toggleEditWeekday = (weekday) => {
      const index = editSelectedWeekdays.value.indexOf(weekday)
      if (index > -1) {
        editSelectedWeekdays.value.splice(index, 1)
      } else {
        editSelectedWeekdays.value.push(weekday)
      }
    }

    const updateEditTrainingDaysFromWeekdays = () => {
      const existingDays = editForm.value.training_days.reduce((acc, day) => {
        acc[day.day_of_week] = day
        return acc
      }, {})

      editForm.value.training_days = editSelectedWeekdays.value.map(weekday => {
        return existingDays[weekday] || {
          day_of_week: weekday,
          exercises: []
        }
      })
    }

    const addEditExercise = (dayIndex) => {
      editForm.value.training_days[dayIndex].exercises.push({
        exercise_id: '',
        sets: 3,
        reps: 10,
        weight: 0,
        instructions: ''
      })
    }

    const removeEditExercise = (dayIndex, exerciseIndex) => {
      editForm.value.training_days[dayIndex].exercises.splice(exerciseIndex, 1)
    }

    // ========== OTHER FUNCTIONS ==========
    const deletePlan = async (planId) => {
      if (!confirm('Czy na pewno chcesz usunąć ten plan?')) return

      try {
        await makeApiCall(`/training-plans/${planId}`, {
          method: 'DELETE'
        })

        if (editingPlan.value?.id === planId) {
          closeEditModal()
        }

        await fetchPlans()
      } catch (err) {
        error.value = err.message || 'Błąd podczas usuwania planu'
      }
    }

    const viewPlanDetails = async (planId) => {
      try {
        const planDetails = await makeApiCall(`/training-plans/${planId}/details`)
        selectedPlan.value = planDetails
        showDetailsModal.value = true
      } catch (err) {
        error.value = err.message || 'Błąd podczas pobierania szczegółów'
      }
    }

    const closeDetailsModal = () => {
      showDetailsModal.value = false
      selectedPlan.value = null
    }

    const previousMonth = () => {
      currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
    }

    const nextMonth = () => {
      currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
    }

    const formatDate = (dateString) => {
      if (!dateString) return ''
      return new Date(dateString).toLocaleDateString('pl-PL')
    }

    const formatDateRange = (startDate, endDate) => {
      const start = formatDate(startDate)
      if (!endDate) return `od ${start}`
      const end = formatDate(endDate)
      return `${start} - ${end}`
    }

    const getPlanStatus = (plan) => {
      if (!plan) return 'unknown'

      const now = new Date()
      const startDate = new Date(plan.start_date)
      const endDate = plan.end_date ? new Date(plan.end_date) : null

      if (now < startDate) return 'upcoming'
      if (endDate && now > endDate) return 'completed'
      return 'active'
    }

    const getPlanStatusText = (plan) => {
      const status = getPlanStatus(plan)
      switch (status) {
        case 'upcoming': return 'Nadchodzący'
        case 'active': return 'Aktywny'
        case 'completed': return 'Zakończony'
        default: return 'Nieznany'
      }
    }

    const getProgressPercentage = (plan) => {
      if (!plan) return 0

      const now = new Date()
      const startDate = new Date(plan.start_date)
      const endDate = plan.end_date ? new Date(plan.end_date) : null

      if (!endDate || now < startDate) return 0
      if (now > endDate) return 100

      const totalDuration = endDate - startDate
      const elapsed = now - startDate
      return Math.round((elapsed / totalDuration) * 100)
    }

    const getProgressText = (plan) => {
      const percentage = getProgressPercentage(plan)
      if (percentage === 0) return 'Nie rozpoczęty'
      if (percentage === 100) return 'Zakończony'
      return `W trakcie (${percentage}%)`
    }

    const getWeekdayIcon = (weekday) => {
      const icons = {
        'poniedziałek': '🏃‍♂️',
        'wtorek': '💪',
        'środa': '🏋️‍♂️',
        'czwartek': '🤸‍♂️',
        'piątek': '🚴‍♂️',
        'sobota': '🏊‍♂️',
        'niedziela': '🧘‍♂️'
      }
      return icons[weekday] || '🏋️'
    }

    // ========== LIFECYCLE ==========
    onMounted(async () => {
      await fetchPlans()
      await fetchCalendarEvents()
      await fetchMyTrainer()
      await fetchExercises()
      await fetchTrainees()
    })

    // ========== RETURN ==========
    return {
      // Store
      store,

      // Reactive data
      loading,
      error,
      plans,
      calendarEvents,
      exercises,
      trainees,
      myTrainer,
      currentDate,
      filterStatus,

      // Create modal
      showCreateModal,
      creating,
      currentStep,
      selectedWeekdays,
      newPlan,

      // Edit modal
      showEditModal,
      editingPlan,
      updating,
      editCurrentStep,
      editSelectedWeekdays,
      editForm,

      // Details modal
      showDetailsModal,
      selectedPlan,

      // Constants
      dayOptions,

      // Computed
      todayDate,
      currentMonthYear,
      planDuration,
      editPlanDuration,
      filteredPlans,
      upcomingWorkouts,
      canProceedToNextStep,
      canProceedToNextEditStep,
      canCreatePlan,
      canUpdatePlan,
      calendarDays,

      // Functions
      fetchPlans,
      createPlan,
      editPlan,
      updatePlan,
      deletePlan,
      viewPlanDetails,
      closeCreateModal,
      closeEditModal,
      closeDetailsModal,
      nextStep,
      nextEditStep,
      updateAvailableDays,
      updateEditAvailableDays,
      isWeekdaySelected,
      isEditWeekdaySelected,
      toggleWeekday,
      toggleEditWeekday,
      updateTrainingDaysFromWeekdays,
      updateEditTrainingDaysFromWeekdays,
      addExercise,
      addEditExercise,
      removeExercise,
      removeEditExercise,
      getExerciseCount,
      getTotalExercisesInPlan,
      previousMonth,
      nextMonth,
      formatDate,
      formatDateRange,
      getPlanStatus,
      getPlanStatusText,
      getProgressPercentage,
      getProgressText,
      getWeekdayIcon
    }
  }
}
</script>
<style scoped>
.training-plans-view {
  margin-left: 1.5rem;
  margin-top: 1.5rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 2rem;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  max-width: 95%;
  width: 100%;
  color: #2d3748;
  overflow-y: auto;
  max-height: calc(100vh - 4rem);
}
.training-plans-view::-webkit-scrollbar {
  width: 6px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 2rem;
  background: white;
  border-radius: 20px;
  box-shadow: 0 8px 25px rgba(0,0,0,0.1);
}

.header-content h1 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
  font-size: 2.5rem;
  font-weight: 700;
}

.header-content p {
  color: #718096;
  font-size: 1.1rem;
  margin: 0;

}

.btn-create-plan {
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  border: none;
  color: white;
  padding: 1rem 2rem;
  border-radius: 15px;
  font-weight: 600;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(45, 80, 22, 0.3);
}

.btn-create-plan:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(45, 80, 22, 0.4);
}

/* Dashboard Cards */
.dashboard-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  grid-column: 1 / -1; /* Rozciąga na pełną szerokość */
}

.dashboard-card {
  background: white;
  padding: 1.5rem;
  border-radius: 15px;
  max-height: 10rem;
  box-shadow: 0 4px 15px rgba(0,0,0,0.08);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.3s ease;
}

.dashboard-card:hover {
  transform: translateY(-2px);
}

.card-icon {
  font-size: 2rem;
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-content h3 {
  color: #2d5016;
  font-size: 2rem;
  font-weight: 700;
  margin: 0;
}

.card-content p {
  color: #718096;
  font-size: 0.9rem;
  margin: 0;
}

.plans-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 2rem;
  margin-top: 2rem;
}

.calendar-section, .plans-section {
  background: white;
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 8px 25px rgba(0,0,0,0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-title {
  color: #2d5016;
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.calendar-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.calendar-nav-btn {
  background: #2d5016;
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
}

.calendar-nav-btn:hover {
  background: #1e3309;
  transform: scale(1.1);
}

.calendar-month {
  color: #2d5016;
  margin: 0;
  text-transform: capitalize;
  font-weight: 600;
}

.calendar-container {
  background: #f8fafc;
  border-radius: 15px;
  padding: 1rem;
}

.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  margin-bottom: 1rem;
}

.weekday {
  padding: 0.5rem;
  text-align: center;
  font-weight: 600;
  color: #4a5568;
  font-size: 0.9rem;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  background: #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
}

.calendar-day {
  background: white;
  padding: 0.5rem;
  min-height: 100px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.calendar-day:hover {
  background: #f7fafc;
}

.calendar-day.today {
  background: linear-gradient(135deg, #e6ffcc 0%, #d4ffa6 100%);
}

.calendar-day.other-month {
  background: #f1f5f9;
  color: #a0aec0;
}

.calendar-day.has-events {
  border-left: 4px solid #2d5016;
}

.day-number {
  color: #2d5016;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.events-container {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.calendar-event {
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  color: white;
  padding: 0.2rem 0.4rem;
  border-radius: 6px;
  font-size: 0.75rem;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.calendar-event:hover {
  transform: scale(1.05);
}

.more-events {
  color: #718096;
  font-size: 0.7rem;
  text-align: center;
  padding: 0.2rem;
}

/* Filter Controls */
.filter-controls {
  display: flex;
  gap: 1rem;
}

.filter-select {
  background: white;
  border: 2px solid #e2e8f0;
  color: #2d5016;
  padding: 0.5rem 1rem;
  border-radius: 10px;
  font-size: 0.9rem;
  cursor: pointer;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
}

.empty-illustration {
  max-width: 400px;
  margin: 0 auto;
}

.empty-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.7;
}

.empty-state h3 {
  color: #2d5016;
  margin-bottom: 1rem;
  font-size: 1.5rem;
}

.empty-state p {
  color: #718096;
  margin-bottom: 2rem;
  line-height: 1.6;
}

.btn-create-first {
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  border: none;
  color: white;
  padding: 1rem 2rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

/* Plans Grid */
.plans-grid {
  display: grid;
  gap: 1.5rem;
}

.plan-card {
  background: white;
  border: 2px solid #f7fafc;
  border-radius: 20px;
  padding: 2rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}

.plan-card:hover {
  border-color: #2d5016;
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(45, 80, 22, 0.15);
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
}

.plan-title-section {
  flex: 1;
}

.plan-name {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
  font-size: 1.3rem;
  font-weight: 600;
}

.plan-badges {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge.active {
  background: #c6f6d5;
  color: #22543d;
}

.badge.upcoming {
  background: #bee3f8;
  color: #2a4365;
}

.badge.completed {
  background: #e2e8f0;
  color: #4a5568;
}

.badge.badge-info {
  background: #e6fffa;
  color: #234e52;
}

.plan-status {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.plan-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-action {
  background: #f7fafc;
  border: 2px solid #e2e8f0;
  color: #2d5016;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-action:hover {
  background: #edf2f7;
  transform: scale(1.1);
}

.btn-danger:hover {
  background: #fed7d7;
  border-color: #feb2b2;
  color: #e53e3e;
}

.plan-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #4a5568;
  font-size: 0.9rem;
}

.meta-icon {
  font-size: 1.1rem;
}

.plan-progress {
  margin-top: 1rem;
}

.progress-bar {
  background: #e2e8f0;
  border-radius: 10px;
  height: 8px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  background: linear-gradient(90deg, #2d5016 0%, #48bb78 100%);
  height: 100%;
  transition: width 0.3s ease;
}

.progress-text {
  color: #718096;
  font-size: 0.8rem;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 2rem;
  backdrop-filter: blur(5px);
}

.modal-content {
  background: white;
  border-radius: 25px;
  padding: 2rem;
  max-width: 800px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px rgba(0,0,0,0.25);
}

.modal-large {
  max-width: 900px;
}

.modal-xl {
  max-width: 1200px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f0f0f0;
}

.modal-header h3 {
  color: #2d5016;
  margin: 0;
  font-size: 1.8rem;
  font-weight: 600;
}

.plan-header-info {
  flex: 1;
}

.btn-close {
  background: #e53e3e;
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
}

.btn-close:hover {
  background: #c53030;
  transform: scale(1.1);
}

/* Creation Steps */
.creation-steps {
  display: flex;
  justify-content: space-between;
  margin-bottom: 2rem;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 15px;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
  padding: 1rem;
  transition: all 0.3s ease;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e2e8f0;
  color: #718096;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  transition: all 0.3s ease;
}

.step.active .step-number {
  background: #2d5016;
  color: white;
}

.step.completed .step-number {
  background: #48bb78;
  color: white;
}

.step span {
  font-size: 0.9rem;
  color: #718096;
  font-weight: 500;
}

.step.active span {
  color: #2d5016;
  font-weight: 600;
}

/* Step Content */
.step-content {
  margin-bottom: 2rem;
}

.step-content h4 {
  color: #2d5016;
  margin-bottom: 1rem;
  font-size: 1.3rem;
}

.step-description {
  color: #718096;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

/* Weekdays Selector */
.weekdays-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.weekday-option {
  background: white;
  border: 3px solid #e2e8f0;
  border-radius: 15px;
  padding: 1.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.weekday-option:hover {
  border-color: #cbd5e0;
  transform: translateY(-2px);
}

.weekday-option.selected {
  border-color: #2d5016;
  background: #f0fff4;
}

.weekday-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.weekday-name {
  color: #2d5016;
  font-weight: 600;
  text-transform: capitalize;
}

.weekday-check {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  background: #2d5016;
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
}

.selected-summary {
  background: #f0fff4;
  border: 2px solid #c6f6d5;
  border-radius: 12px;
  padding: 1rem;
  color: #22543d;
}

/* Duration Info */
.duration-info {
  margin-top: 1rem;
}

.duration-badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
}

/* Training Days Exercises */
.training-days-exercises {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.training-day-exercises {
  background: #f8fafc;
  border-radius: 20px;
  padding: 2rem;
}

.day-header-fancy {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  background: white;
  padding: 1rem;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}

.day-icon {
  font-size: 2rem;
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-header-fancy h5 {
  color: #2d5016;
  margin: 0;
  font-size: 1.3rem;
  font-weight: 600;
  flex: 1;
  text-transform: capitalize;
}

.exercises-count {
  background: #e6fffa;
  color: #234e52;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
}

.exercises-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.exercise-card {
  background: white;
  border: 2px solid #f7fafc;
  border-radius: 15px;
  padding: 1.5rem;
  transition: all 0.3s ease;
}

.exercise-card:hover {
  border-color: #e2e8f0;
  box-shadow: 0 4px 15px rgba(0,0,0,0.05);
}

.exercise-main {
  margin-bottom: 1rem;
}

.exercise-params-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-top: 1rem;
}

.param-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.param-group label {
  color: #4a5568;
  font-size: 0.9rem;
  font-weight: 600;
}

.param-input {
  background: white;
  border: 2px solid #e2e8f0;
  color: #2d5016;
  padding: 0.75rem;
  border-radius: 10px;
  font-size: 1rem;
  text-align: center;
  font-weight: 600;
}

.param-input:focus {
  outline: none;
  border-color: #2d5016;
  box-shadow: 0 0 10px rgba(45, 80, 22, 0.1);
}

.instructions-textarea {
  background: #f8fafc;
  border: 2px solid #e2e8f0;
  color: #4a5568;
  padding: 0.75rem;
  border-radius: 10px;
  font-size: 0.9rem;
  min-height: 80px;
  resize: vertical;
  width: 100%;
  font-family: inherit;
}

.btn-add-exercise {
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  border: none;
  color: white;
  padding: 1rem 2rem;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  align-self: flex-start;
}

.btn-add-exercise:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(45, 80, 22, 0.3);
}

.btn-remove-exercise {
  background: #e53e3e;
  border: none;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  align-self: flex-start;
  transition: all 0.3s ease;
}

.btn-remove-exercise:hover {
  background: #c53030;
  transform: scale(1.05);
}

/* Form Elements */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
.form-group label {
  color: #2d5016;
  font-weight: bold;
}

.form-help {
  color: #e53e3e;
  font-size: 0.9rem;
  font-style: italic;
}

.form-input, .form-select, .form-textarea {
  background: white;
  border: 2px solid #e2e8f0;
  color: #2d5016;
  padding: 0.75rem;
  border-radius: 10px;
  font-size: 1rem;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #2d5016;
  box-shadow: 0 0 10px rgba(45, 80, 22, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.training-days-section {
  background: #f9f9f9;
  border-radius: 15px;
  padding: 1.5rem;
}

.training-day {
  background: white;
  border-radius: 10px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.day-header {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 1rem;
}

.exercises-section h5 {
  color: #2d5016;
  margin-bottom: 1rem;
}

.exercise-item {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.exercise-params {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.5rem;
  margin: 0.5rem 0;
}

.form-input-small {
  background: white;
  border: 1px solid #e2e8f0;
  color: #2d5016;
  padding: 0.5rem;
  border-radius: 6px;
  font-size: 0.9rem;
}

.form-textarea-small {
  background: white;
  border: 1px solid #e2e8f0;
  color: #2d5016;
  padding: 0.5rem;
  border-radius: 6px;
  font-size: 0.9rem;
  min-height: 60px;
  resize: vertical;
}

.btn-add, .btn-remove {
  background: #2d5016;
  border: none;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
}

.btn-remove {
  background: #e53e3e;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.btn-cancel, .btn-submit {
  padding: 1rem 2rem;
  border-radius: 10px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel {
  background: #718096;
  border: none;
  color: white;
}

.btn-submit {
  background: #2d5016;
  border: none;
  color: white;
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.plan-details {
  color: #2d5016;
}

.plan-meta-details {
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #f0f0f0;
}

.training-days-details h4 {
  color: #2d5016;
  margin-bottom: 1rem;
}

.day-details {
  background: #f9f9f9;
  border-radius: 10px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.day-details h5 {
  color: #2d5016;
  margin-bottom: 1rem;
}

.exercises-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.exercise-details {
  background: white;
  border-radius: 8px;
  padding: 1rem;
}

.exercise-name {
  font-weight: bold;
  color: #2d5016;
  margin-bottom: 0.5rem;
}

.exercise-params {
  color: #666;
  margin-bottom: 0.5rem;
}

.exercise-instructions {
  color: #666;
  font-style: italic;
}

.no-exercises {
  color: #999;
  font-style: italic;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .training-plans-view {
    padding: 1rem;
    margin: 0.5rem;
  }

  .plans-content {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    gap: 1rem;
  }

  .calendar-grid {
    grid-template-columns: repeat(7, 1fr);
    gap: 1px;
  }

  .calendar-day {
    min-height: 60px;
    padding: 0.25rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .exercise-params {
    grid-template-columns: 1fr;
  }

  .modal-content {
    margin: 1rem;
    padding: 1rem;
  }
}
/* Brakujące style CSS */

/* Loading i Error States */
.loading-container, .error-container {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f0f0f0;
  border-top: 4px solid #2d5016;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-container {
  background: #fff5f5;
  border: 2px solid #feb2b2;
  border-radius: 15px;
}

.error-icon {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.btn-retry {
  background: #e53e3e;
  border: none;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 10px;
  cursor: pointer;
  margin-top: 1rem;
  transition: all 0.3s ease;
}

.btn-retry:hover {
  background: #c53030;
  transform: translateY(-2px);
}

/* Modal Navigation */
.modal-navigation {
  display: flex;
  gap: 1rem;
  justify-content: space-between;
  align-items: center;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 2px solid #f0f0f0;
}

.btn-nav {
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  font-size: 1rem;
}

.btn-prev {
  background: #718096;
  color: white;
}

.btn-prev:hover {
  background: #4a5568;
  transform: translateX(-3px);
}

.btn-next {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-next:hover {
  transform: translateX(3px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.btn-next:disabled {
  background: #e2e8f0;
  color: #a0aec0;
  cursor: not-allowed;
  transform: none;
}

.btn-create {
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  color: white;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(45, 80, 22, 0.3);
}

.btn-create:disabled {
  background: #e2e8f0;
  color: #a0aec0;
  cursor: not-allowed;
  transform: none;
}

.btn-cancel {
  background: transparent;
  color: #718096;
  border: 2px solid #e2e8f0;
}

.btn-cancel:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
}

.btn-back {
  background: #f7fafc;
  border: 2px solid #e2e8f0;
  color: #2d5016;
  padding: 0.75rem 1.5rem;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
}

.btn-back:hover {
  background: #edf2f7;
  transform: translateX(-3px);
}

/* No Days Selected */
.no-days-selected {
  text-align: center;
  padding: 3rem;
  background: #fff5f5;
  border: 2px solid #feb2b2;
  border-radius: 15px;
  color: #c53030;
}

.no-days-selected p {
  font-size: 1.1rem;
  margin-bottom: 1rem;
}

/* Plan Details Enhanced */
.plan-details-enhanced {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.plan-overview {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.overview-card {
  background: #f8fafc;
  border-radius: 15px;
  padding: 1.5rem;
}

.overview-card h4 {
  color: #2d5016;
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.overview-card p {
  color: #4a5568;
  margin-bottom: 0.5rem;
  line-height: 1.6;
}

/* Progress Circle */
.progress-circle {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 1rem 0;
}

.circle-progress {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: conic-gradient(#2d5016 calc(var(--progress) * 3.6deg), #e2e8f0 0deg);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.circle-progress::before {
  content: '';
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: white;
  position: absolute;
}

.circle-progress span {
  position: relative;
  z-index: 1;
  font-weight: 700;
  color: #2d5016;
  font-size: 1rem;
}

/* Training Schedule */
.training-schedule h4 {
  color: #2d5016;
  margin-bottom: 1.5rem;
  font-size: 1.3rem;
}

.schedule-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.schedule-day {
  background: #f8fafc;
  border-radius: 15px;
  padding: 1.5rem;
  border: 2px solid #e2e8f0;
  transition: all 0.3s ease;
}

.schedule-day:hover {
  border-color: #2d5016;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(45, 80, 22, 0.1);
}

.day-header-detailed {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
  padding: 1rem;
  background: white;
  border-radius: 10px;
}

.day-icon-large {
  font-size: 2.5rem;
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #2d5016 0%, #48bb78 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.day-info {
  flex: 1;
}

.day-info h5 {
  color: #2d5016;
  margin: 0 0 0.25rem 0;
  font-size: 1.2rem;
  font-weight: 600;
  text-transform: capitalize;
}

.day-exercises {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.exercise-item-detailed {
  display: flex;
  gap: 1rem;
  background: white;
  border-radius: 10px;
  padding: 1rem;
  transition: all 0.3s ease;
}

.exercise-item-detailed:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.exercise-number {
  background: #2d5016;
  color: white;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.exercise-content {
  flex: 1;
}

.exercise-name-detailed {
  color: #2d5016;
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

.exercise-specs {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.spec-item {
  background: #e6fffa;
  color: #234e52;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.8rem;
  font-weight: 600;
}

.exercise-notes {
  color: #718096;
  font-size: 0.9rem;
  font-style: italic;
  background: #f7fafc;
  padding: 0.5rem;
  border-radius: 6px;
  border-left: 3px solid #cbd5e0;
}

.no-exercises-day {
  text-align: center;
  padding: 2rem;
  color: #a0aec0;
  background: white;
  border-radius: 10px;
}

.empty-day-icon {
  font-size: 2rem;
  display: block;
  margin-bottom: 0.5rem;
  opacity: 0.7;
}

.no-exercises-day p {
  margin: 0;
  font-style: italic;
}

/* Exercise Instructions */
.exercise-instructions {
  margin-top: 1rem;
}

.instructions-textarea:focus {
  outline: none;
  border-color: #2d5016;
  box-shadow: 0 0 10px rgba(45, 80, 22, 0.1);
}

/* Additional Mobile Responsiveness */
@media (max-width: 768px) {
  .plan-overview {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .schedule-grid {
    grid-template-columns: 1fr;
  }

  .day-header-detailed {
    flex-direction: column;
    text-align: center;
    gap: 0.5rem;
  }

  .exercise-item-detailed {
    flex-direction: column;
    gap: 0.5rem;
  }

  .exercise-specs {
    flex-direction: column;
    gap: 0.5rem;
  }

  .weekdays-selector {
    grid-template-columns: repeat(2, 1fr);
  }

  .exercise-params-grid {
    grid-template-columns: 1fr;
  }

  .modal-navigation {
    flex-direction: column;
    gap: 1rem;
  }

  .btn-nav {
    width: 100%;
  }

  .dashboard-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .dashboard-cards {
    grid-template-columns: 1fr;
  }

  .weekdays-selector {
    grid-template-columns: 1fr;
  }

  .card-content h3 {
    font-size: 1.5rem;
  }

  .header-content h1 {
    font-size: 2rem;
  }
}
</style>
