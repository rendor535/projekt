<template>
  <div class="trainer-profile-wrapper">
    <div class="trainer-profile-view">
      <!-- Trainer Profile Card -->
      <div class="trainer-profile-card">
        <div class="trainer-header">
          <div class="trainer-avatar-container">
            <div class="trainer-avatar" :style="{ backgroundColor: avatarColor }">
              <span class="trainer-initials">{{ trainerInitials }}</span>
            </div>
            <button @click="generateNewAvatar" class="refresh-avatar" title="Nowy kolor awatara">
              🎨
            </button>
          </div>
          <div class="trainer-basic">
            <h2>{{ trainerFullName }}</h2>
            <p class="trainer-subtitle">Profil trenera</p>
            <div class="trainer-date">{{ formattedProfileDate }}</div>
          </div>
        </div>

        <div v-if="trainerProfile" class="trainer-details">
          <div class="detail-section">
            <h4>📋 Podstawowe informacje</h4>
            <div class="detail-row">
              <span class="detail-icon">👤</span>
              <span class="detail-label">Imię i nazwisko:</span>
              <span class="detail-value">{{ trainerProfile.first_name }} {{ trainerProfile.last_name }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-icon">🏆</span>
              <span class="detail-label">Doświadczenie:</span>
              <span class="detail-value">{{ trainerProfile.experience }} {{ getYearsText(trainerProfile.experience) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-icon">💰</span>
              <span class="detail-label">Cena za godzinę:</span>
              <span class="detail-value">{{ trainerProfile.price_per_hour }} zł</span>
            </div>
            <div class="detail-row">
              <span class="detail-icon">📍</span>
              <span class="detail-label">Lokalizacja:</span>
              <span class="detail-value">{{ trainerProfile.location || 'Nie podano' }}</span>
            </div>
          </div>
          <div class="detail-grid">
          <div class="detail-section">
            <h4>🎯 Specjalizacje</h4>
            <div class="tags-container">
              <span v-for="spec in trainerProfile.specializations" :key="spec" class="tag specialization-tag">
                {{ spec }}
              </span>
            </div>
          </div>

          <div class="detail-section">
            <h4>📜 Certyfikaty</h4>
            <div class="tags-container">
              <span v-for="cert in trainerProfile.certifications" :key="cert" class="tag certification-tag">
                {{ cert }}
              </span>
            </div>
          </div>

          <div class="detail-section">
            <h4>🌍 Języki</h4>
            <div class="tags-container">
              <span v-for="lang in trainerProfile.languages" :key="lang" class="tag language-tag">
                {{ lang }}
              </span>
            </div>
          </div>
          </div>
          <div class="detail-section">
            <h4>📧 Kontakt</h4>
            <div class="detail-row" v-if="trainerProfile.contact_email">
              <span class="detail-icon">📧</span>
              <span class="detail-label">Email:</span>
              <span class="detail-value">{{ trainerProfile.contact_email }}</span>
            </div>
            <div class="detail-row" v-if="trainerProfile.contact_phone">
              <span class="detail-icon">📱</span>
              <span class="detail-label">Telefon:</span>
              <span class="detail-value">{{ trainerProfile.contact_phone }}</span>
            </div>
          </div>

          <div class="detail-section">
            <h4>📝 O mnie</h4>
            <div class="bio-text">
              {{ trainerProfile.bio }}
            </div>
          </div>
        </div>

        <div v-else-if="!loading" class="no-trainer-profile">
          <div class="no-profile-content">
            <span class="no-profile-icon">🏋️‍♂️</span>
            <h3>Brak profilu trenera</h3>
            <p>Uzupełnij swoje dane jako trener, aby klienci mogli Cię znaleźć</p>
            <button @click="showEditForm = true" class="create-profile-btn">
              Utwórz profil trenera
            </button>
          </div>
        </div>

        <div v-else class="loading-state">
          <div class="loading-spinner"></div>
          <p>Ładowanie profilu trenera...</p>
        </div>

        <!-- Edit Profile Button -->
        <div v-if="trainerProfile" class="trainer-actions">
          <button @click="openEditForm" class="edit-profile-btn">
            ✏️ Edytuj profil trenera
          </button>
        </div>
      </div>

      <!-- Stats Section -->
      <div v-if="trainerProfile" class="trainer-stats-section">
        <div class="stat-card">
          <div class="stat-header">
            <h3>📊 Statystyki</h3>
          </div>
          <div class="stat-item">
            <span class="stat-label">Lata doświadczenia:</span>
            <span class="stat-value">{{ trainerProfile.experience }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Liczba specjalizacji:</span>
            <span class="stat-value">{{ trainerProfile.specializations.length }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Liczba certyfikatów:</span>
            <span class="stat-value">{{ trainerProfile.certifications.length }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Status:</span>
            <span class="stat-value" :class="{ 'active': trainerProfile.is_active, 'inactive': !trainerProfile.is_active }">
              {{ trainerProfile.is_active ? 'Aktywny' : 'Nieaktywny' }}
            </span>
          </div>
        </div>

        <div class="experience-card">
          <div class="stat-header">
            <h3>🎓 Poziom doświadczenia</h3>
          </div>
          <div class="experience-level">
            <div class="level-indicator" :class="experienceLevel.class">
              {{ experienceLevel.label }}
            </div>
            <div class="level-description">
              {{ experienceLevel.description }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Trainer Profile Modal -->
    <div v-if="showEditForm" class="modal-overlay" @click="closeModal">
      <div class="modal-content trainer-modal" @click.stop>
        <div class="modal-header">
          <h3>{{ trainerProfile ? 'Edytuj profil trenera' : 'Utwórz profil trenera' }}</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>

        <form @submit.prevent="saveTrainerProfile" class="trainer-form">
          <div class="form-row">
            <div class="form-group">
              <label>Imię:</label>
              <input v-model="editForm.first_name" type="text" required />
            </div>
            <div class="form-group">
              <label>Nazwisko:</label>
              <input v-model="editForm.last_name" type="text" required />
            </div>
          </div>

          <div class="form-group">
            <label>O mnie (biografia):</label>
            <textarea v-model="editForm.bio" rows="4" required placeholder="Opisz swoje doświadczenie, podejście do treningu, filozofię..."></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Lata doświadczenia:</label>
              <input v-model.number="editForm.experience" type="number" min="0" max="50" required />
            </div>
            <div class="form-group">
              <label>Cena za godzinę (zł):</label>
              <input v-model.number="editForm.price_per_hour" type="number" min="0" step="0.01" required />
            </div>
          </div>

          <div class="form-group">
            <label>Specjalizacje (oddziel przecinkami):</label>
            <input v-model="specializationsText" type="text" placeholder="np. Siłownia, Cardio, Yoga, Pilates, CrossFit" />
          </div>

          <div class="form-group">
            <label>Certyfikaty (oddziel przecinkami):</label>
            <input v-model="certificationsText" type="text" placeholder="np. Personal Trainer Level 2, Yoga Instructor, First Aid" />
          </div>

          <div class="form-group">
            <label>Języki (oddziel przecinkami):</label>
            <input v-model="languagesText" type="text" placeholder="np. Polski, Angielski, Niemiecki" />
          </div>

          <div class="form-group">
            <label>Lokalizacja:</label>
            <input v-model="editForm.location" type="text" placeholder="np. Warszawa, Kraków, Online" />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Email kontaktowy:</label>
              <input v-model="editForm.contact_email" type="email" />
            </div>
            <div class="form-group">
              <label>Telefon kontaktowy:</label>
              <input v-model="editForm.contact_phone" type="tel" />
            </div>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeModal" class="cancel-btn">Anuluj</button>
            <button type="submit" :disabled="saving" class="save-btn">
              {{ saving ? 'Zapisywanie...' : 'Zapisz profil trenera' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const trainerProfile = ref(null)
const avatarColor = ref('#8de78d')
const loading = ref(true)
const saving = ref(false)
const showEditForm = ref(false)

const editForm = ref({
  first_name: '',
  last_name: '',
  bio: '',
  specializations: [],
  experience: 0,
  certifications: [],
  price_per_hour: 0,
  languages: [],
  contact_email: '',
  contact_phone: '',
  location: ''
})

const specializationsText = ref('')
const certificationsText = ref('')
const languagesText = ref('')

// Computed properties
const trainerFullName = computed(() => {
  if (trainerProfile.value) {
    return `${trainerProfile.value.first_name} ${trainerProfile.value.last_name}`
  }
  return 'Nowy trener'
})

const trainerInitials = computed(() => {
  if (trainerProfile.value) {
    return `${trainerProfile.value.first_name.charAt(0)}${trainerProfile.value.last_name.charAt(0)}`.toUpperCase()
  }
  return 'NT'
})

const formattedProfileDate = computed(() => {
  if (trainerProfile.value?.created_at) {
    try {
      const date = new Date(trainerProfile.value.created_at)
      return `Trener od ${date.toLocaleDateString('pl-PL', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })}`
    } catch (e) {
      return 'Nowy trener'
    }
  }
  return 'Nowy trener'
})

const experienceLevel = computed(() => {
  const exp = trainerProfile.value?.experience || 0
  if (exp < 1) return { label: 'Początkujący', class: 'beginner', description: 'Świeżo po kursach, pełen energii' }
  if (exp < 3) return { label: 'Junior', class: 'junior', description: 'Nabiera doświadczenia, rozwija umiejętności' }
  if (exp < 5) return { label: 'Średni', class: 'intermediate', description: 'Pewny siebie, z solidną bazą wiedzy' }
  if (exp < 10) return { label: 'Doświadczony', class: 'experienced', description: 'Ekspert w swojej dziedzinie' }
  return { label: 'Master', class: 'master', description: 'Guru fitnessu z ogromnym doświadczeniem' }
})

// Helper functions
const getYearsText = (years) => {
  if (years === 1) return 'rok'
  if (years >= 2 && years <= 4) return 'lata'
  return 'lat'
}

const generateNewAvatar = () => {
  const colors = ['#8de78d','#6ed46e','#5bc45b','#a0c4ff','#c4d9ff','#87ceeb','#ffc6a0','#ffe0c4','#ffb366','#ff9999','#ffb3b3','#ff8a80','#dda0dd','#e6e6fa','#d8bfd8','#98e8a0','#b3ffb3','#90ee90']
  avatarColor.value = colors[Math.floor(Math.random()*colors.length)]
}

const closeModal = () => {
  showEditForm.value = false
  // Reset text fields
  specializationsText.value = ''
  certificationsText.value = ''
  languagesText.value = ''
}

const openEditForm = () => {
  if (trainerProfile.value) {
    editForm.value = { ...trainerProfile.value }
    specializationsText.value = trainerProfile.value.specializations.join(', ')
    certificationsText.value = trainerProfile.value.certifications.join(', ')
    languagesText.value = trainerProfile.value.languages.join(', ')
  } else {
    editForm.value = {
      first_name: '',
      last_name: '',
      bio: '',
      specializations: [],
      experience: 0,
      certifications: [],
      price_per_hour: 0,
      languages: [],
      contact_email: '',
      contact_phone: '',
      location: ''
    }
  }
  showEditForm.value = true
}

const saveTrainerProfile = async () => {
  saving.value = true
  try {
    // Convert text fields to arrays
    editForm.value.specializations = specializationsText.value.split(',').map(s => s.trim()).filter(s => s)
    editForm.value.certifications = certificationsText.value.split(',').map(s => s.trim()).filter(s => s)
    editForm.value.languages = languagesText.value.split(',').map(s => s.trim()).filter(s => s)

    console.log('Saving trainer profile:', editForm.value)

    const token = localStorage.getItem('token')
    if (!token) {
      throw new Error('Brak tokenu autoryzacji')
    }

    const response = await axios.post('/api/trainer-profile', editForm.value, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    console.log('Trainer profile save response:', response.data)

    if (response.status === 200) {
      await loadTrainerProfile()
      closeModal()
      alert('Profil trenera został zapisany!')
    }
  } catch (error) {
    console.error('Trainer profile save error:', error)

    let errorMessage = 'Błąd zapisu profilu trenera'
    if (error.response?.status === 401) {
      errorMessage = 'Błąd autoryzacji. Zaloguj się ponownie.'
    } else if (error.response?.data?.error) {
      errorMessage = error.response.data.error
    }

    alert(errorMessage)
  }
  saving.value = false
}

const loadTrainerProfile = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      console.error('No token found')
      trainerProfile.value = null
      loading.value = false
      return
    }

    const response = await axios.get('/api/trainer-profile', {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    console.log('Trainer profile data loaded:', response.data)
    trainerProfile.value = response.data.profile || null

  } catch (error) {
    console.error('Error loading trainer profile:', error)
    if (error.response?.status === 401) {
      console.error('Unauthorized - token may be invalid')
    }
    trainerProfile.value = null
  }
  loading.value = false
}

onMounted(async () => {
  console.log('TrainerProfileView mounted')
  generateNewAvatar()
  await loadTrainerProfile()
})
</script>

<style scoped>
.trainer-profile-wrapper {
  padding: 1rem;
  background: rgba(139, 231, 139, 0.1);
  backdrop-filter: blur(10px);
  min-height: 100vh;
  width: 100%;
}

.trainer-profile-view {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 1.5rem;
  width: 100%;
  min-height: calc(100vh - 2rem);
}

/* Trainer Profile Card */
.trainer-profile-card {
  background: rgba(255, 255, 255, 0.95);
  padding: 1.5rem;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  height: fit-content;
  overflow-y: auto;
}

.trainer-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px solid #e8f5e8;
}

.trainer-avatar-container {
  position: relative;
}

.trainer-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px solid #8de78d;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.trainer-initials {
  font-size: 1.8rem;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.refresh-avatar {
  position: absolute;
  bottom: -5px;
  right: -5px;
  width: 25px;
  height: 25px;
  border-radius: 50%;
  border: none;
  background: #8de78d;
  cursor: pointer;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.trainer-basic h2 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
  font-size: 1.6rem;
}

.trainer-subtitle {
  color: #666;
  margin: 0 0 0.5rem 0;
  font-size: 0.9rem;
}

.trainer-date {
  color: #999;
  font-size: 0.8rem;
}

/* Detail Sections */
.detail-section {
  margin-bottom: 1.5rem;
}

.detail-section h4 {
  margin: 0 0 0.8rem 0;
  color: #2d5016;
  font-size: 1rem;
  border-bottom: 1px solid #e8f5e8;
  padding-bottom: 0.3rem;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.6rem;
  background: rgba(248, 252, 248, 0.8);
  border-radius: 8px;
  margin-bottom: 0.5rem;
}

.detail-icon {
  font-size: 1rem;
  width: 20px;
  text-align: center;
}

.detail-label {
  font-weight: 500;
  color: #2d5016;
  min-width: 120px;
  font-size: 0.85rem;
}

.detail-value {
  color: #333;
  font-weight: 600;
  font-size: 0.85rem;
}

/* Tags */
.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.3rem 0.8rem;
  border-radius: 20px;
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
  background: #fff3e0;
  color: #f57c00;
}

/* Bio */
.bio-text {
  background: rgba(248, 252, 248, 0.8);
  padding: 1rem;
  border-radius: 8px;
  line-height: 1.6;
  color: #333;
  font-size: 0.9rem;
  border-left: 4px solid #8de78d;
}

/* Stats Section */
.trainer-stats-section {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.stat-card, .experience-card {
  background: rgba(255, 255, 255, 0.95);
  padding: 1.2rem;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.stat-header h3 {
  margin: 0 0 1rem 0;
  color: #2d5016;
  font-size: 1rem;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: 0.6rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-label {
  color: #666;
  font-size: 0.85rem;
}

.stat-value {
  font-weight: 600;
  color: #333;
  font-size: 0.85rem;
}

.stat-value.active {
  color: #2e7d32;
}

.stat-value.inactive {
  color: #d32f2f;
}

/* Experience Level */
.experience-level {
  text-align: center;
}

.level-indicator {
  padding: 0.8rem 1.5rem;
  border-radius: 25px;
  font-weight: 700;
  font-size: 1.1rem;
  margin-bottom: 0.8rem;
}

.level-indicator.beginner { background: #e3f2fd; color: #1976d2; }
.level-indicator.junior { background: #e8f5e8; color: #2e7d32; }
.level-indicator.intermediate { background: #fff3e0; color: #f57c00; }
.level-indicator.experienced { background: #f3e5f5; color: #7b1fa2; }
.level-indicator.master { background: #ffebee; color: #d32f2f; }

.level-description {
  color: #666;
  font-size: 0.8rem;
  line-height: 1.4;
}

/* Modal Styles */
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
  z-index: 1000;
}

.trainer-modal {
  max-width: 600px;
  width: 95%;
  max-height: 90vh;
}

.modal-content {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}
.detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}
@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
.modal-header h3 {
  margin: 0;
  color: #2d5016;
  font-size: 1.3rem;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
}

/* Form Styles */
.trainer-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.form-group label {
  font-weight: 500;
  color: #2d5016;
  font-size: 0.9rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.9rem;
  font-family: inherit;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #8de78d;
  box-shadow: 0 0 0 3px rgba(141, 231, 141, 0.2);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.cancel-btn, .save-btn {
  flex: 1;
  padding: 0.8rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.cancel-btn {
  border: 1px solid #ddd;
  background: white;
  color: #666;
}

.save-btn {
  border: none;
  background: linear-gradient(135deg, #8de78d, #6ed46e);
  color: white;
}

.save-btn:hover {
  background: linear-gradient(135deg, #6ed46e, #5bc45b);
  transform: translateY(-1px);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* No Profile State */
.no-trainer-profile {
  text-align: center;
  padding: 2rem;
}

.no-profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.no-profile-icon {
  font-size: 3rem;
}

.no-trainer-profile h3 {
  margin: 0;
  color: #2d5016;
  font-size: 1.2rem;
}

.no-trainer-profile p {
  color: #666;
  margin: 0;
  font-size: 0.95rem;
}

.create-profile-btn,
.edit-profile-btn {
  background: linear-gradient(135deg, #8de78d, #6ed46e);
  color: white;
  border: none;
  padding: 0.8rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(141, 231, 141, 0.4);
  font-size: 0.9rem;
}

.create-profile-btn:hover,
.edit-profile-btn:hover {
  background: linear-gradient(135deg, #6ed46e, #5bc45b);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(141, 231, 141, 0.5);
}

.trainer-actions {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 2px solid #e8f5e8;
  text-align: center;
}

/* Loading state */
.loading-state {
  text-align: center;
  padding: 2rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #8de78d;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Responsive Design */
@media (max-width: 1024px) {
  .trainer-profile-view {
    grid-template-columns: 1fr;
    gap: 1.2rem;
  }

  .trainer-header {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .trainer-profile-wrapper {
    padding: 0.5rem;
  }

  .trainer-modal {
    width: 98%;
    margin: 1rem;
  }

  .tags-container {
    justify-content: center;
  }
}
</style>
