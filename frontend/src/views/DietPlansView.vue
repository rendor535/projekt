<template>
  <div class="diet-plans-view">
    <!-- Header Section -->
    <div class="page-header">
      <div class="header-content">
        <h1>🥗 Plany Żywieniowe</h1>
        <p v-if="store.user?.role === 'trainer'">
          Twórz i zarządzaj planami żywieniowymi dla swoich podopiecznych
        </p>
        <p v-else>
          Twoje plany żywieniowe od trenera
        </p>
      </div>

      <!-- Create Diet Plan Button (Trainer Only) -->
      <button
          v-if="store.user?.role === 'trainer'"
          @click="showCreateModal = true"
          class="create-btn"
      >
        <span class="btn-icon">➕</span>
        Stwórz plan
      </button>
    </div>

    <!-- Filters Section -->
    <div class="filters-section">
      <div class="filter-group">
        <label class="filter-label">
          <input
              type="checkbox"
              v-model="showActiveOnly"
              @change="loadDietPlans"
          />
          Pokaż tylko aktywne plany
        </label>
      </div>

      <div v-if="store.user?.role === 'trainer'" class="filter-group">
        <label for="clientFilter">Filtruj po podopiecznym:</label>
        <select id="clientFilter" v-model="selectedClientFilter" @change="loadDietPlans">
          <option value="">Wszyscy podopieczni</option>
          <option v-for="trainee in trainees" :key="trainee.id" :value="trainee.id">
            {{ trainee.username }}
          </option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Ładowanie planów żywieniowych...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <h3>❌ Wystąpił błąd</h3>
      <p>{{ error }}</p>
      <button @click="loadDietPlans" class="retry-btn">🔄 Spróbuj ponownie</button>
    </div>

    <!-- Diet Plans List -->
    <div v-else class="diet-plans-container">
      <!-- Empty State -->
      <div v-if="dietPlans.length === 0" class="empty-state">
        <div class="empty-icon">🍽️</div>
        <h3>{{ store.user?.role === 'trainer' ? 'Brak utworzonych planów' : 'Brak przypisanych planów' }}</h3>
        <p v-if="store.user?.role === 'trainer'">
          Stwórz pierwszy plan żywieniowy dla swoich podopiecznych
        </p>
        <p v-else>
          Skontaktuj się z trenerem, aby otrzymać plan żywieniowy
        </p>
      </div>

      <!-- Diet Plans Grid -->
      <div v-else class="diet-plans-grid">
        <div
            v-for="plan in dietPlans"
            :key="plan.id"
            class="diet-plan-card"
            :class="{ active: plan.is_active }"
        >
          <!-- Plan Header -->
          <div class="plan-header">
            <div class="plan-info">
              <h3>{{ plan.name }}</h3>
              <div class="plan-meta">
                <span class="calories">🔥 {{ plan.calories || 0 }} kcal/dzień</span>
                <span class="status" :class="plan.is_active ? 'active' : 'inactive'">
                  {{ plan.is_active ? '✅ Aktywny' : '⏸️ Nieaktywny' }}
                </span>
              </div>
              <div v-if="store.user?.role === 'trainer'" class="client-info">
                👤 {{ plan.client_username || plan.client_full_name || `ID: ${plan.client_id}` }}
              </div>
              <div v-else class="trainer-info">
                👨‍💼 {{ plan.trainer_username || plan.trainer_full_name || `ID: ${plan.trainer_id}` }}
              </div>
            </div>

            <div class="plan-actions">
              <button @click="viewPlan(plan)" class="action-btn view-btn">
                👁️ Zobacz
              </button>
              <button
                  v-if="store.user?.role === 'trainer'"
                  @click="editPlan(plan)"
                  class="action-btn edit-btn"
              >
                ✏️ Edytuj
              </button>
              <button
                  v-if="store.user?.role === 'trainee'"
                  @click="activatePlan(plan)"
                  class="action-btn toggle-btn"
                  :disabled="processingPlan === plan.id || plan.is_active"
              >
                {{ processingPlan === plan.id ? '⏳' : '▶️' }}
                {{ plan.is_active ? 'Aktywny' : 'Aktywuj' }}
              </button>
              <button
                  v-if="store.user?.role === 'trainer'"
                  @click="deletePlan(plan)"
                  class="action-btn delete-btn"
                  :disabled="plan.is_active"
              >
                🗑️ Usuń
              </button>
            </div>
          </div>

          <!-- Plan Description -->
          <div v-if="plan.description" class="plan-description">
            <p>{{ plan.description }}</p>
          </div>

          <!-- Plan Duration -->
          <div class="plan-duration">
            <span class="duration-label">📅 Okres:</span>
            <span v-if="plan.start_date && plan.end_date">
              {{ formatDate(plan.start_date) }} - {{ formatDate(plan.end_date) }}
            </span>
            <span v-else-if="plan.start_date">
              Od {{ formatDate(plan.start_date) }}
            </span>
            <span v-else>Nieokreślony</span>
          </div>

          <!-- Meals Preview -->
          <div class="meals-preview">
            <h4>📋 Posiłki ({{ plan.meals_count || 0 }})</h4>
            <div class="meals-summary">
              <span class="meals-info">Kliknij "Zobacz" aby zobaczyć szczegóły posiłków</span>
            </div>
          </div>

          <!-- Plan Footer -->
          <div class="plan-footer">
            <span class="created-date">
              📆 Utworzony: {{ formatDate(plan.created_at) }}
            </span>
            <div v-if="store.user?.role === 'trainer'" class="plan-controls">
              <button
                  @click="duplicatePlan(plan)"
                  class="control-btn"
                  title="Duplikuj plan"
              >
                📋
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Diet Plan Modal -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModals">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditModal ? '✏️ Edytuj Plan Żywieniowy' : '➕ Nowy Plan Żywieniowy' }}</h3>
          <button @click="closeModals" class="close-btn">❌</button>
        </div>

        <form @submit.prevent="submitDietPlan" class="diet-form">
          <!-- Basic Information -->
          <div class="form-section">
            <h4>📝 Podstawowe informacje</h4>

            <div class="form-group">
              <label for="planName">Nazwa planu: *</label>
              <input
                  id="planName"
                  v-model="planForm.name"
                  type="text"
                  required
                  placeholder="np. Plan redukcyjny - 1800 kcal"
                  maxlength="255"
              />
            </div>

            <div class="form-group">
              <label for="planDescription">Opis planu:</label>
              <textarea
                  id="planDescription"
                  v-model="planForm.description"
                  rows="3"
                  placeholder="Opisz cel i założenia planu żywieniowego..."
              ></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="clientSelect">Podopieczny: *</label>
                <select id="clientSelect" v-model="planForm.client_id" required>
                  <option value="">Wybierz podopiecznego</option>
                  <option
                      v-for="trainee in trainees"
                      :key="trainee.id"
                      :value="trainee.id"
                  >
                    {{ trainee.username }}
                  </option>
                </select>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="startDate">Data rozpoczęcia:</label>
                <input
                    id="startDate"
                    v-model="planForm.start_date"
                    type="date"
                />
              </div>

              <div class="form-group">
                <label for="endDate">Data zakończenia:</label>
                <input
                    id="endDate"
                    v-model="planForm.end_date"
                    type="date"
                    :min="planForm.start_date"
                />
              </div>
            </div>
          </div>

          <!-- Meals Configuration -->
          <div class="form-section">
            <div class="section-header">
              <h4>🍽️ Konfiguracja posiłków</h4>
              <button
                  type="button"
                  @click="addMeal"
                  class="add-meal-btn"
                  :disabled="planForm.meals.length >= 6"
              >
                ➕ Dodaj posiłek
              </button>
            </div>

            <div v-if="planForm.meals.length === 0" class="no-meals-message">
              <p>Brak skonfigurowanych posiłków. Dodaj pierwszy posiłek.</p>
            </div>

            <div v-else class="meals-configuration">
              <div
                  v-for="(meal, index) in planForm.meals"
                  :key="index"
                  class="meal-config"
              >
                <div class="meal-header">
                  <h5>Posiłek {{ index + 1 }}</h5>
                  <button
                      type="button"
                      @click="removeMeal(index)"
                      class="remove-meal-btn"
                      title="Usuń posiłek"
                  >
                    🗑️
                  </button>
                </div>

                <div class="meal-form">
                  <div class="form-group">
                    <label>Nazwa posiłku: *</label>
                    <input
                        v-model="meal.name"
                        type="text"
                        required
                        placeholder="np. Śniadanie"
                    />
                  </div>

                  <!-- Products in meal -->
                  <div class="products-section">
                    <h6>Produkty w posiłku:</h6>
                    <div class="product-search">
                      <input
                          v-model="productSearchTerms[index]"
                          type="text"
                          placeholder="Szukaj produktu..."
                          @input="searchProducts(index)"
                      />
                      <div v-if="productSearchResults[index]?.length" class="search-results">
                        <div
                            v-for="product in productSearchResults[index]"
                            :key="product.id"
                            class="search-result-item"
                            @click="addProductToMeal(meal, product, index)"
                        >
                          {{ product.name }} ({{ product.calories }}kcal/100g)
                        </div>
                      </div>
                    </div>

                    <div v-if="meal.products?.length" class="meal-products">
                      <div
                          v-for="(product, pIndex) in meal.products"
                          :key="pIndex"
                          class="meal-product"
                      >
                        <span class="product-info">{{ product.product_name || `Produkt ID: ${product.product_id}` }}</span>
                        <input
                            v-model.number="product.amount_grams"
                            type="number"
                            min="1"
                            step="1"
                            placeholder="Gramy"
                            @input="updateMealNutrition(meal)"
                        />
                        <span>g</span>
                        <button
                            type="button"
                            @click="removeProductFromMeal(meal, pIndex)"
                            class="remove-product-btn"
                        >
                          ❌
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="form-actions">
            <button type="button" @click="closeModals" class="cancel-btn">
              Anuluj
            </button>
            <button
                type="submit"
                class="submit-btn"
                :disabled="formLoading || !isFormValid"
            >
              {{ formLoading ? '⏳ Zapisywanie...' : (showEditModal ? '💾 Zapisz zmiany' : '✅ Utwórz plan') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- View Diet Plan Modal -->
    <div v-if="showViewModal && selectedPlan" class="modal-overlay" @click="showViewModal = false">
      <div class="modal-content view-modal" @click.stop>
        <div class="modal-header">
          <h3>👁️ {{ selectedPlan.name }}</h3>
          <button @click="showViewModal = false" class="close-btn">❌</button>
        </div>

        <div class="plan-details">
          <!-- Plan Info -->
          <div class="detail-section">
            <h4>📋 Informacje o planie</h4>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="label">Kalorie dzienne:</span>
                <span class="value">🔥 {{ selectedPlan.calories }} kcal</span>
              </div>
              <div class="detail-item">
                <span class="label">Status:</span>
                <span class="value" :class="selectedPlan.is_active ? 'active' : 'inactive'">
                  {{ selectedPlan.is_active ? '✅ Aktywny' : '⏸️ Nieaktywny' }}
                </span>
              </div>
              <div v-if="store.user?.role === 'trainer'" class="detail-item">
                <span class="label">Podopieczny:</span>
                <span class="value">👤 {{ selectedPlan.client_username }}</span>
              </div>
              <div v-else class="detail-item">
                <span class="label">Trener:</span>
                <span class="value">👨‍💼 {{ selectedPlan.trainer_username }}</span>
              </div>
              <div class="detail-item">
                <span class="label">Okres:</span>
                <span class="value">
                  📅 {{ selectedPlan.start_date ? formatDate(selectedPlan.start_date) : 'Brak' }} -
                  {{ selectedPlan.end_date ? formatDate(selectedPlan.end_date) : 'Bez końca' }}
                </span>
              </div>
            </div>
            <div v-if="selectedPlan.description" class="description">
              <p>{{ selectedPlan.description }}</p>
            </div>
          </div>

          <!-- Meals Details -->
          <div class="detail-section">
            <h4>🍽️ Szczegóły posiłków</h4>
            <div v-if="selectedPlanMeals?.length" class="meals-detail">
              <div
                  v-for="(meal, index) in selectedPlanMeals"
                  :key="meal.id"
                  class="meal-detail"
              >
                <div class="meal-header">
                  <h5>{{ index + 1 }}. {{ meal.name }}</h5>
                  <span class="meal-calories">{{ meal.calories }} kcal</span>
                </div>
                <div class="meal-macros">
                  <div class="macro">
                    <span class="macro-label">Białko:</span>
                    <span class="macro-value">{{ meal.proteins || 0 }}g</span>
                  </div>
                  <div class="macro">
                    <span class="macro-label">Tłuszcze:</span>
                    <span class="macro-value">{{ meal.fats || 0 }}g</span>
                  </div>
                  <div class="macro">
                    <span class="macro-label">Węglowodany:</span>
                    <span class="macro-value">{{ meal.carbs || 0 }}g</span>
                  </div>
                </div>

                <!-- Products in meal -->
                <div v-if="meal.products?.length" class="meal-products-view">
                  <h6>Produkty:</h6>
                  <div v-for="product in meal.products" :key="product.id" class="product-view">
                    <div class="product-info">
                      <span class="product-name">{{ product.product_name }}</span>
                      <span class="product-amount">{{ product.amount_grams }}g</span>
                    </div>
                    <div class="product-nutrition">
                      <span>{{ Math.round(product.total_calories) }}kcal</span>
                      <span>B: {{ Math.round(product.total_protein) }}g</span>
                      <span>T: {{ Math.round(product.total_fat) }}g</span>
                      <span>W: {{ Math.round(product.total_carbs) }}g</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="no-meals">
              <p>Brak szczegółowych informacji o posiłkach.</p>
            </div>
          </div>
        </div>

        <div class="view-actions">
          <button @click="showViewModal = false" class="close-view-btn">
            Zamknij
          </button>
          <button
              v-if="store.user?.role === 'trainer'"
              @click="editPlanFromView"
              class="edit-from-view-btn"
          >
            ✏️ Edytuj plan
          </button>
        </div>
      </div>
    </div>

    <!-- Success/Error Notifications -->
    <div v-if="notification" class="notification" :class="notification.type">
      <span>{{ notification.message }}</span>
      <button @click="notification = null" class="notification-close">×</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

const store = useAuthStore()

// API Configuration
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// State variables
const loading = ref(true)
const error = ref(null)
const dietPlans = ref([])
const trainees = ref([])
const notification = ref(null)
const processingPlan = ref(null)
const formLoading = ref(false)

// Filter states
const showActiveOnly = ref(false)
const selectedClientFilter = ref('')

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const selectedPlan = ref(null)
const selectedPlanMeals = ref([])
const editingPlan = ref(null)

// Product search
const productSearchTerms = ref({})
const productSearchResults = ref({})
const allProducts = ref([])

// Form data
const planForm = ref({
  name: '',
  description: '',
  client_id: '',
  start_date: '',
  end_date: '',
  meals: []
})

// Default meal structure for new backend format
const defaultMeal = () => ({
  name: '',
  order: 1,
  products: []
})

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

// Notification utility
const showNotification = (message, type = 'success') => {
  notification.value = { message, type }
  setTimeout(() => {
    notification.value = null
  }, 5000)
}

// Computed properties
const isFormValid = computed(() => {
  return planForm.value.name &&
      planForm.value.client_id &&
      planForm.value.meals.length > 0 &&
      planForm.value.meals.every(meal =>
          meal.name &&
          meal.products &&
          meal.products.length > 0 &&
          meal.products.every(product => product.amount_grams > 0)
      )
})

// Data loading functions
const loadDietPlans = async () => {
  loading.value = true
  error.value = null

  try {
    let endpoint
    let params = new URLSearchParams()

    if (store.user?.role === 'trainer') {
      endpoint = '/diet-plans/trainer'
      if (selectedClientFilter.value) {
        params.append('client_id', selectedClientFilter.value)
      }
    } else {
      endpoint = '/diet-plans/trainee'
    }

    if (showActiveOnly.value) {
      params.append('is_active', 'true')
    }

    const queryString = params.toString()
    const fullEndpoint = queryString ? `${endpoint}?${queryString}` : endpoint

    const data = await apiRequest(fullEndpoint)
    dietPlans.value = data || []
  } catch (err) {
    console.error('Error loading diet plans:', err)
    error.value = `Nie udało się załadować planów żywieniowych: ${err.message}`
  } finally {
    loading.value = false
  }
}

const loadTrainees = async () => {
  if (store.user?.role !== 'trainer') return

  try {
    const data = await apiRequest('/trainees')
    trainees.value = data || []
  } catch (err) {
    console.error('Error loading trainees:', err)
  }
}

const loadProducts = async () => {
  try {
    const data = await apiRequest('/products')
    allProducts.value = data || []
  } catch (err) {
    console.error('Error loading products:', err)
  }
}

// Product search functionality
const searchProducts = (mealIndex) => {
  const term = productSearchTerms.value[mealIndex]
  if (!term || term.length < 2) {
    productSearchResults.value[mealIndex] = []
    return
  }

  const filtered = allProducts.value.filter(product =>
      product.name.toLowerCase().includes(term.toLowerCase())
  ).slice(0, 10)

  productSearchResults.value[mealIndex] = filtered
}

const addProductToMeal = (meal, product, mealIndex) => {
  if (!meal.products) {
    meal.products = []
  }

  // Check if product already exists in meal
  const existingProduct = meal.products.find(p => p.product_id === product.id)
  if (existingProduct) {
    showNotification('Produkt już jest dodany do tego posiłku', 'warning')
    return
  }

  meal.products.push({
    product_id: product.id,
    product_name: product.name,
    amount_grams: 100,
    calories_per_100g: product.calories,
    protein_per_100g: product.protein,
    fat_per_100g: product.fat,
    carbs_per_100g: product.carbs
  })

  productSearchTerms.value[mealIndex] = ''
  productSearchResults.value[mealIndex] = []
  updateMealNutrition(meal)
}

const removeProductFromMeal = (meal, productIndex) => {
  meal.products.splice(productIndex, 1)
  updateMealNutrition(meal)
}

const updateMealNutrition = (meal) => {
  if (!meal.products) return

  let totalCalories = 0
  let totalProtein = 0
  let totalFat = 0
  let totalCarbs = 0

  meal.products.forEach(product => {
    const factor = (product.amount_grams || 0) / 100
    totalCalories += (product.calories_per_100g || 0) * factor
    totalProtein += (product.protein_per_100g || 0) * factor
    totalFat += (product.fat_per_100g || 0) * factor
    totalCarbs += (product.carbs_per_100g || 0) * factor
  })

  meal.calculated_calories = Math.round(totalCalories)
  meal.calculated_protein = Math.round(totalProtein)
  meal.calculated_fat = Math.round(totalFat)
  meal.calculated_carbs = Math.round(totalCarbs)
}

// Plan management functions
const viewPlan = async (plan) => {
  try {
    selectedPlan.value = plan
    // Load detailed meals for the plan
    const mealsData = await apiRequest(`/diet-plans/${plan.id}/details`)
    selectedPlanMeals.value = mealsData.meals || []
    showViewModal.value = true
  } catch (err) {
    console.error('Error loading plan details:', err)
    showNotification('Nie udało się załadować szczegółów planu', 'error')
  }
}

const editPlan = async (plan) => {
  try {
    selectedPlan.value = plan
    const mealsData = await apiRequest(`/diet-plans/${plan.id}/details`)

    editingPlan.value = plan
    planForm.value = {
      name: plan.name,
      description: plan.description || '',
      client_id: plan.client_id,
      start_date: plan.start_date ? plan.start_date.split('T')[0] : '',
      end_date: plan.end_date ? plan.end_date.split('T')[0] : '',
      meals: (mealsData.meals || []).map(meal => ({
        name: meal.name,
        order: meal.meal_order,
        products: meal.products || []
      }))
    }

    showEditModal.value = true
  } catch (err) {
    console.error('Error loading plan for editing:', err)
    showNotification('Nie udało się załadować planu do edycji', 'error')
  }
}

const editPlanFromView = () => {
  if (selectedPlan.value) {
    showViewModal.value = false
    editPlan(selectedPlan.value)
  }
}

const activatePlan = async (plan) => {
  processingPlan.value = plan.id

  try {
    await apiRequest(`/diet-plans/${plan.id}/activate`, {
      method: 'PUT'
    })

    // Refresh the plans list
    await loadDietPlans()
    showNotification('Plan został aktywowany', 'success')
  } catch (err) {
    console.error('Error activating plan:', err)
    showNotification('Nie udało się aktywować planu', 'error')
  } finally {
    processingPlan.value = null
  }
}

const duplicatePlan = async (plan) => {
  try {
    const mealsData = await apiRequest(`/diet-plans/${plan.id}/details`)

    planForm.value = {
      name: `${plan.name} (kopia)`,
      description: plan.description || '',
      client_id: '',
      start_date: '',
      end_date: '',
      meals: (mealsData.meals || []).map(meal => ({
        name: meal.name,
        order: meal.meal_order,
        products: meal.products || []
      }))
    }

    showCreateModal.value = true
  } catch (err) {
    console.error('Error duplicating plan:', err)
    showNotification('Nie udało się zduplikować planu', 'error')
  }
}

const deletePlan = async (plan) => {
  if (!confirm(`Czy na pewno chcesz usunąć plan "${plan.name}"? Ta operacja jest nieodwracalna.`)) {
    return
  }

  try {
    await apiRequest(`/diet-plans/${plan.id}`, {
      method: 'DELETE'
    })

    dietPlans.value = dietPlans.value.filter(p => p.id !== plan.id)
    showNotification('Plan został usunięty', 'success')
  } catch (err) {
    console.error('Error deleting plan:', err)
    showNotification('Nie udało się usunąć planu', 'error')
  }
}

// Form management functions
const addMeal = () => {
  const newMeal = defaultMeal()
  newMeal.order = planForm.value.meals.length + 1
  planForm.value.meals.push(newMeal)
}

const removeMeal = (index) => {
  planForm.value.meals.splice(index, 1)
  // Reorder remaining meals
  planForm.value.meals.forEach((meal, idx) => {
    meal.order = idx + 1
  })
}

const resetForm = () => {
  planForm.value = {
    name: '',
    description: '',
    client_id: '',
    start_date: '',
    end_date: '',
    meals: []
  }
  editingPlan.value = null
  productSearchTerms.value = {}
  productSearchResults.value = {}
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  resetForm()
}
const toISOStringIfDate = (dateStr) => {
  return dateStr ? new Date(dateStr).toISOString() : null
}

const submitDietPlan = async () => {
  if (!isFormValid.value) {
    showNotification('Wypełnij wszystkie wymagane pola', 'error')
    return
  }

  formLoading.value = true

  try {
    // Update meal nutritions before submitting
    planForm.value.meals.forEach(meal => updateMealNutrition(meal))

    const planData = {
      name: planForm.value.name,
      description: planForm.value.description,
      client_id: parseInt(planForm.value.client_id),
      start_date: toISOStringIfDate(planForm.value.start_date),
      end_date: toISOStringIfDate(planForm.value.end_date),
      meals: planForm.value.meals.map((meal, index) => ({
        name: meal.name,
        order: index + 1,
        products: meal.products.map(product => ({
          product_id: product.product_id,
          amount_grams: product.amount_grams
        }))
      }))
    }

    if (showEditModal.value && editingPlan.value) {
      // Update existing plan
      await apiRequest(`/diet-plans/${editingPlan.value.id}`, {
        method: 'PUT',
        body: JSON.stringify(planData)
      })
      showNotification('Plan został zaktualizowany', 'success')
    } else {
      // Create new plan
      await apiRequest('/diet-plans', {
        method: 'POST',
        body: JSON.stringify(planData)
      })
      showNotification('Plan został utworzony', 'success')
    }

    closeModals()
    await loadDietPlans()

  } catch (err) {
    console.error('Error submitting diet plan:', err)
    showNotification(`Nie udało się zapisać planu: ${err.message}`, 'error')
  } finally {
    formLoading.value = false
  }
}

// Utility functions
const formatDate = (dateString) => {
  if (!dateString) return 'Brak'
  const date = new Date(dateString)
  return date.toLocaleDateString('pl-PL')
}

// Initialize component
onMounted(async () => {
  await Promise.all([
    loadDietPlans(),
    loadTrainees(),
    loadProducts()
  ])

  // Add one default meal when creating new plan
  if (store.user?.role === 'trainer') {
    planForm.value.meals.push(defaultMeal())
  }
})

// Watch for changes in trainees to set default client
watch(trainees, (newTrainees) => {
  if (newTrainees.length === 1 && !planForm.value.client_id) {
    planForm.value.client_id = newTrainees[0].id
  }
}, { immediate: true })
</script>

<style scoped>
/* Zachowuję wszystkie istniejące style i dodaję nowe dla nowych funkcji */

.diet-plans-view {
  background: white;
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);
  max-width: 95%;
  width: 100%;
  color: #181818;
  overflow-y: auto;
  max-height: calc(100vh - 4rem);
  margin-left: 1.5rem;
  margin-top: 1.5rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f0f0f0;
}

.header-content h1 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
  font-size: 2rem;
}

.header-content p {
  color: #666;
  margin: 0;
}

.create-btn {
  background: linear-gradient(135deg, #2d5016 0%, #4a7c3e 100%);
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 4px 15px rgba(45, 80, 22, 0.3);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(45, 80, 22, 0.4);
}

/* Filters Section */
.filters-section {
  display: flex;
  gap: 2rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 12px;
  align-items: center;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

/* Products Section */
.products-section {
  margin-top: 1rem;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
}

.product-search {
  position: relative;
  margin-bottom: 1rem;
}

.product-search input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #ddd;
  border-top: none;
  border-radius: 0 0 6px 6px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 10;
}

.search-result-item {
  padding: 0.75rem;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
}

.search-result-item:hover {
  background: #f8f9fa;
}

.meal-products {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.meal-product {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 6px;
}

.product-info {
  flex: 1;
  font-weight: 500;
}

.meal-product input {
  width: 80px;
  padding: 0.25rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  text-align: center;
}

.remove-product-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
}

/* Meal Products View */
.meal-products-view {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid #e9ecef;
}

.meal-products-view h6 {
  margin: 0 0 0.5rem 0;
  color: #666;
  font-size: 0.9rem;
}

.product-view {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 0.25rem;
}

.product-nutrition {
  display: flex;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: #666;
}

/* Loading and Error States */
.loading-state, .error-state {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #2d5016;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.retry-btn {
  background: #2d5016;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  margin-top: 1rem;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #666;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  color: #2d5016;
  margin-bottom: 1rem;
}

/* Diet Plans Grid */
.diet-plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 2rem;
}

.diet-plan-card {
  background: #fff;
  border: 2px solid #e9ecef;
  border-radius: 16px;
  padding: 1.5rem;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.diet-plan-card:hover {
  border-color: #b7f5b7;
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.1);
}

.diet-plan-card.active {
  border-color: #2d5016;
  background: linear-gradient(135deg, #f8fff8 0%, #f0f8f0 100%);
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.plan-info h3 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
  font-size: 1.3rem;
}

.plan-meta {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  margin-bottom: 0.5rem;
}

.calories {
  color: #e74c3c;
  font-weight: 600;
}

.status.active {
  color: #27ae60;
  font-weight: 600;
}

.status.inactive {
  color: #7f8c8d;
}

.client-info, .trainer-info {
  color: #666;
  font-size: 0.9rem;
}

.plan-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.action-btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.view-btn {
  background: #3498db;
  color: white;
}

.view-btn:hover {
  background: #2980b9;
}

.edit-btn {
  background: #f39c12;
  color: white;
}

.edit-btn:hover {
  background: #e67e22;
}

.toggle-btn {
  background: #27ae60;
  color: white;
}

.toggle-btn:hover {
  background: #229954;
}

.toggle-btn:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.delete-btn {
  background: #e74c3c;
  color: white;
}

.delete-btn:hover {
  background: #c0392b;
}

.delete-btn:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.plan-description {
  margin-bottom: 1rem;
  color: #666;
  font-style: italic;
}

.plan-duration {
  margin-bottom: 1rem;
  color: #666;
}

.duration-label {
  font-weight: 600;
}

.meals-preview {
  margin-bottom: 1rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.meals-preview h4 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
}

.meals-summary {
  color: #666;
  font-style: italic;
}

.plan-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #e9ecef;
}

.created-date {
  color: #666;
  font-size: 0.9rem;
}

.plan-controls {
  display: flex;
  gap: 0.5rem;
}

.control-btn {
  background: none;
  border: 1px solid #ddd;
  padding: 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.control-btn:hover {
  background: #f8f9fa;
}

.control-btn.delete {
  border-color: #e74c3c;
  color: #e74c3c;
}

.control-btn.delete:hover {
  background: #e74c3c;
  color: white;
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
  border-radius: 16px;
  max-width: 90vw;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
}

.modal-content:not(.view-modal) {
  width: 800px;
}

.view-modal {
  width: 600px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
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
  padding: 0.5rem;
  border-radius: 6px;
  transition: background 0.3s ease;
}

.close-btn:hover {
  background: #f8f9fa;
}

/* Form Styles */
.diet-form {
  padding: 1.5rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h4 {
  margin: 0 0 1rem 0;
  color: #2d5016;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.add-meal-btn {
  background: #27ae60;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.add-meal-btn:hover {
  background: #229954;
}

.add-meal-btn:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #2d5016;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #2d5016;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.checkbox-custom {
  width: 20px;
  height: 20px;
  border: 2px solid #e9ecef;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.checkbox-label input[type="checkbox"] {
  display: none;
}

.checkbox-label input[type="checkbox"]:checked + .checkbox-custom {
  background: #2d5016;
  border-color: #2d5016;
}

.checkbox-label input[type="checkbox"]:checked + .checkbox-custom::after {
  content: '✓';
  color: white;
  font-weight: bold;
}

/* Meals Configuration */
.no-meals-message {
  text-align: center;
  padding: 2rem;
  color: #666;
  background: #f8f9fa;
  border-radius: 8px;
  border: 2px dashed #ddd;
}

.meals-configuration {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.meal-config {
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
  background: #f8f9fa;
}

.meal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.meal-header h5 {
  margin: 0;
  color: #2d5016;
}

.remove-meal-btn {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.remove-meal-btn:hover {
  background: #c0392b;
}

.meal-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.macros-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 1rem;
}

/* Form Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e9ecef;
}

.cancel-btn {
  background: #6c757d;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.cancel-btn:hover {
  background: #5a6268;
}

.submit-btn {
  background: linear-gradient(135deg, #2d5016 0%, #4a7c3e 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
}

.submit-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(45, 80, 22, 0.3);
}

.submit-btn:disabled {
  background: #95a5a6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* View Modal Specific Styles */
.plan-details {
  padding: 1.5rem;
}

.detail-section {
  margin-bottom: 2rem;
}

.detail-section h4 {
  margin: 0 0 1rem 0;
  color: #2d5016;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.detail-item .label {
  font-weight: 600;
  color: #666;
}

.detail-item .value {
  color: #2d5016;
  font-weight: 600;
}

.detail-item .value.active {
  color: #27ae60;
}

.detail-item .value.inactive {
  color: #7f8c8d;
}

.description {
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  font-style: italic;
  color: #666;
}

.meals-detail {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.meal-detail {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  background: #fff;
}

.meal-detail .meal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.meal-detail .meal-header h5 {
  margin: 0;
  color: #2d5016;
}

.meal-calories {
  color: #e74c3c;
  font-weight: 600;
}

.meal-macros {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.macro {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 6px;
  min-width: 80px;
}

.macro-label {
  font-size: 0.8rem;
  color: #666;
  margin-bottom: 0.25rem;
}

.macro-value {
  font-weight: 600;
  color: #2d5016;
}

.nutrition-summary {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.nutrition-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.nutrition-label {
  color: #666;
  font-weight: 600;
}

.nutrition-value {
  color: #2d5016;
  font-weight: 600;
}

.view-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1.5rem;
  border-top: 1px solid #e9ecef;
}

.close-view-btn {
  background: #6c757d;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.close-view-btn:hover {
  background: #5a6268;
}

.edit-from-view-btn {
  background: #f39c12;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.edit-from-view-btn:hover {
  background: #e67e22;
}

/* Notifications */
.notification {
  position: fixed;
  top: 2rem;
  right: 2rem;
  z-index: 1001;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  color: white;
  display: flex;
  align-items: center;
  gap: 1rem;
  min-width: 300px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.notification.success {
  background: #27ae60;
}

.notification.error {
  background: #e74c3c;
}

.notification.warning {
  background: #f39c12;
}

.notification-close {
  background: none;
  border: none;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  transition: background 0.3s ease;
}

.notification-close:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* Responsive Design */
@media (max-width: 768px) {
  .diet-plans-view {
    padding: 1rem;
    margin: 1rem;
  }

  .page-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .filters-section {
    flex-direction: column;
    gap: 1rem;
  }

  .diet-plans-grid {
    grid-template-columns: 1fr;
  }

  .plan-header {
    flex-direction: column;
    gap: 1rem;
  }

  .plan-actions {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .macros-row {
    grid-template-columns: 1fr;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .nutrition-summary {
    grid-template-columns: 1fr;
  }

  .modal-content {
    width: 95vw;
    margin: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }

  .view-actions {
    flex-direction: column;
  }
}
</style>