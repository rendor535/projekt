F<template>
  <div class="profile-wrapper">
    <div class="profile-view">
      <!-- Profile Card with Avatar and Basic Info -->
      <div class="profile-card">
        <div class="profile-header">
          <div class="avatar-container">
            <div class="profile-avatar" :style="{ backgroundColor: avatarColor }">
              <span class="avatar-initials">{{ userInitials }}</span>
            </div>
            <button @click="generateNewAvatar" class="refresh-avatar" title="Nowy kolor awatara">
              🎨
            </button>
          </div>
          <div class="profile-basic">
            <h2>{{ currentUser }}</h2>
            <p class="profile-subtitle">Twój profil zdrowia</p>
            <div class="profile-date">{{ formattedProfileDate }}</div>
          </div>
        </div>

        <div v-if="profile" class="profile-details">
          <div class="detail-row">
            <span class="detail-icon">🎂</span>
            <span class="detail-label">Wiek:</span>
            <span class="detail-value">{{ profile.age }} lat</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">📏</span>
            <span class="detail-label">Wzrost:</span>
            <span class="detail-value">{{ profile.height }} cm</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">⚖️</span>
            <span class="detail-label">Waga:</span>
            <span class="detail-value">{{ profile.weight }} kg</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">{{ profile.gender === 'male' ? '👨' : '👩' }}</span>
            <span class="detail-label">Płeć:</span>
            <span class="detail-value">{{ getGenderText(profile.gender) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">🏃</span>
            <span class="detail-label">Aktywność:</span>
            <span class="detail-value">{{ getActivityText(profile.activity_level) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-icon">🎯</span>
            <span class="detail-label">Cel:</span>
            <span class="detail-value">{{ getGoalText(profile.goal) }}</span>
          </div>
        </div>

        <div v-else-if="!loading" class="no-profile">
          <div class="no-profile-content">
            <span class="no-profile-icon">📝</span>
            <h3>Brak profilu</h3>
            <p>Uzupełnij swoje dane, aby zobaczyć analizę zdrowia</p>
            <button @click="showEditForm = true" class="create-profile-btn">
              Utwórz profil
            </button>
          </div>
        </div>

        <div v-else class="loading-state">
          <div class="loading-spinner"></div>
          <p>Ładowanie danych profilu...</p>
        </div>

        <!-- Edit Profile Button -->
        <div v-if="profile" class="profile-actions">
          <button @click="openEditForm" class="edit-profile-btn">
            ✏️ Edytuj profil
          </button>
        </div>
      </div>

      <!-- Metrics Section -->
      <div v-if="profile" class="metrics-section">
        <!-- BMI Card -->
        <div class="metric-card bmi-card">
          <div class="metric-header">
            <h3>🏥 Kalkulator BMI</h3>
            <div class="bmi-value" :class="bmiCategory.class">
              {{ bmi }}
            </div>
          </div>
          <div class="bmi-details">
            <div class="bmi-category" :class="bmiCategory.class">
              {{ bmiCategory.label }}
            </div>
            <div class="bmi-range">
              Norma: 18.5 - 24.9
            </div>
            <div class="bmi-progress">
              <div class="bmi-bar">
                <div
                    class="bmi-indicator"
                    :style="{ left: bmiPosition + '%' }"
                ></div>
              </div>
              <div class="bmi-labels">
                <span>Niedowaga</span>
                <span>Norma</span>
                <span>Nadwaga</span>
                <span>Otyłość</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Daily Calories Card -->
        <div class="metric-card calories-card">
          <div class="metric-header">
            <h3>🔥 Dzienne Kalorie</h3>
            <div class="calories-value">
              {{ dailyCalories }} kcal
            </div>
          </div>
          <div class="calories-breakdown">
            <div class="calorie-item">
              <span class="calorie-label">BMR (podstawowa przemiana):</span>
              <span class="calorie-value">{{ bmr }} kcal</span>
            </div>
            <div class="calorie-item">
              <span class="calorie-label">Z aktywnością:</span>
              <span class="calorie-value">{{ dailyCalories }} kcal</span>
            </div>
            <div class="calorie-item">
              <span class="calorie-label">Poziom aktywności:</span>
              <span class="calorie-value">{{ getActivityText(profile.activity_level) }}</span>
            </div>
          </div>
        </div>

        <!-- Body Composition Card -->
        <div class="metric-card body-card">
          <div class="metric-header">
            <h3>📊 Szacowany Skład Ciała</h3>
            <span class="info-tooltip" title="Wartości szacunkowe na podstawie wzorów naukowych">ℹ️</span>
          </div>
          <div class="body-metrics">
            <div class="body-metric">
              <div class="metric-icon">💧</div>
              <div class="metric-info">
                <div class="metric-name">Woda</div>
                <div class="metric-val">{{ waterPercentage }}%</div>
              </div>
            </div>
            <div class="body-metric">
              <div class="metric-icon">🦴</div>
              <div class="metric-info">
                <div class="metric-name">Kości</div>
                <div class="metric-val">{{ boneWeight }} kg</div>
              </div>
            </div>
            <div class="body-metric">
              <div class="metric-icon">💪</div>
              <div class="metric-info">
                <div class="metric-name">Mięśnie</div>
                <div class="metric-val">{{ musclePercentage }}%</div>
              </div>
            </div>
            <div class="body-metric">
              <div class="metric-icon">🟡</div>
              <div class="metric-info">
                <div class="metric-name">Tłuszcz</div>
                <div class="metric-val">{{ bodyFatPercentage }}%</div>
              </div>
            </div>
          </div>
          <div class="composition-note">
            <small>
              <strong>Jak obliczamy skład ciała:</strong><br>
              • <strong>Woda:</strong> Na podstawie płci i wieku (mężczyźni ~60%, kobiety ~55%)<br>
              • <strong>Kości:</strong> Około 15% masy ciała<br>
              • <strong>Mięśnie:</strong> Wzór oparty na płci i BMI<br>
              • <strong>Tłuszcz:</strong> Wzór Deurenberg uwzględniający BMI, wiek i płeć
            </small>
          </div>
        </div>

        <!-- Ideal Weight Card -->
        <div class="metric-card ideal-card">
          <div class="metric-header">
            <h3>🎯 Idealna Waga</h3>
          </div>
          <div class="ideal-weights">
            <div class="ideal-item">
              <span class="ideal-method">Robinson:</span>
              <span class="ideal-value">{{ idealWeight.robinson }} kg</span>
            </div>
            <div class="ideal-item">
              <span class="ideal-method">Miller:</span>
              <span class="ideal-value">{{ idealWeight.miller }} kg</span>
            </div>
            <div class="ideal-item">
              <span class="ideal-method">Devine:</span>
              <span class="ideal-value">{{ idealWeight.devine }} kg</span>
            </div>
            <div class="weight-difference" :class="weightDifferenceClass">
              {{ weightDifferenceText }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Profile Modal -->
    <div v-if="showEditForm" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ profile ? 'Edytuj profil' : 'Utwórz profil' }}</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>

        <form @submit.prevent="saveProfile" class="profile-form">
          <div class="form-group">
            <label>Wiek:</label>
            <input v-model.number="editForm.age" type="number" min="10" max="120" required />
          </div>

          <div class="form-group">
            <label>Wzrost (cm):</label>
            <input v-model.number="editForm.height" type="number" min="100" max="250" required />
          </div>

          <div class="form-group">
            <label>Waga (kg):</label>
            <input v-model.number="editForm.weight" type="number" min="30" max="300" step="0.1" required />
          </div>

          <div class="form-group">
            <label>Płeć:</label>
            <select v-model="editForm.gender" required>
              <option value="">Wybierz płeć</option>
              <option value="male">Mężczyzna</option>
              <option value="female">Kobieta</option>
            </select>
          </div>

          <div class="form-group">
            <label>Poziom aktywności:</label>
            <select v-model.number="editForm.activity_level" required>
              <option value="">Wybierz poziom</option>
              <option value="1.2">Siedzący tryb życia</option>
              <option value="1.38">Lekka aktywność (1-3 dni/tydzień)</option>
              <option value="1.55">Umiarkowana aktywność (3-5 dni/tydzień)</option>
              <option value="1.73">Wysoka aktywność (6-7 dni/tydzień)</option>
              <option value="1.9">Bardzo wysoka aktywność (2x dziennie)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Cel:</label>
            <select v-model="editForm.goal" required>
              <option value="">Wybierz cel</option>
              <option value="reduction">Utrata wagi</option>
              <option value="utrzymanie_wagi">Utrzymanie wagi</option>
              <option value="przyrost_masy">Przyrost masy</option>
              <option value="budowanie_miesni">Budowanie mięśni</option>
              <option value="rekompozycja_ciala">Rekompozycja ciała</option>
            </select>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeModal" class="cancel-btn">Anuluj</button>
            <button type="submit" :disabled="saving" class="save-btn">
              {{ saving ? 'Zapisywanie...' : 'Zapisz' }}
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
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const currentUser = ref('')
const profile = ref(null)
const avatarColor = ref('#8de78d')
const loading = ref(true)
const saving = ref(false)
const showEditForm = ref(false)

const editForm = ref({
  age: '',
  height: '',
  weight: '',
  gender: '',
  activity_level: '',
  goal: ''
})

// Show profile creation date instead of current date
const formattedProfileDate = computed(() => {
  if (profile.value?.created_at) {
    return new Date(profile.value.created_at).toLocaleDateString('pl-PL', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
  return 'Brak daty utworzenia'
})

// Generate initials
const userInitials = computed(() => {
  return currentUser.value
      .split(' ')
      .map(n => n.charAt(0).toUpperCase())
      .join('')
      .substring(0,2) || 'U'
})

const generateNewAvatar = () => {
  const colors = ['#8de78d','#6ed46e','#5bc45b','#a0c4ff','#c4d9ff','#87ceeb','#ffc6a0','#ffe0c4','#ffb366','#ff9999','#ffb3b3','#ff8a80','#dda0dd','#e6e6fa','#d8bfd8','#98e8a0','#b3ffb3','#90ee90']
  avatarColor.value = colors[Math.floor(Math.random()*colors.length)]
}

const getGenderText = gender => ({ male:'Mężczyzna', female:'Kobieta' })[gender] || gender
const getActivityText = level => ({
  1.2:'Siedzący tryb',
  1.38:'Lekka aktywność',
  1.55:'Umiarkowana aktywność',
  1.73:'Wysoka aktywność',
  1.9:'Bardzo wysoka aktywność'
}[level] || `Poziom ${level}`)

const getGoalText = goal => ({
  // Angielskie klucze
  reduction: 'Utrata wagi',
  recomposition: 'Rekompozycja ciała', 
  mass: 'Przyrost masy',
  
  // Klucze z podkreślnikami (prawdopodobnie z bazy danych)
  'utrata_wagi': 'Utrata wagi',
  'utrzymanie_wagi': 'Utrzymanie wagi',
  'przyrost_masy': 'Przyrost masy',
  'budowanie_miesni': 'Budowanie mięśni',
  'rekompozycja_ciala': 'Rekompozycja ciała',
  
  // Polskie teksty (z formularza)
  'Utrata wagi': 'Utrata wagi',
  'Utrzymanie wagi': 'Utrzymanie wagi', 
  'Przyrost masy': 'Przyrost masy',
  'Budowanie mięśni': 'Budowanie mięśni',
  'Rekompozycja': 'Rekompozycja ciała'
}[goal] || goal || 'Nie ustawiono')

// Calculations...
const bmi = computed(() => {
  if (!profile.value) return 0
  const h = profile.value.height/100
  return Math.round((profile.value.weight/(h*h))*10)/10
})
const bmiCategory = computed(() => {
  const v=bmi.value
  if(v<18.5) return {label:'Niedowaga',class:'underweight'}
  if(v<25)   return {label:'Waga prawidłowa',class:'normal'}
  if(v<30)   return {label:'Nadwaga',class:'overweight'}
  return {label:'Otyłość',class:'obese'}
})
const bmiPosition = computed(() => {
  const v=bmi.value
  if(v<18.5) return (v/18.5)*25
  if(v<25)   return 25+((v-18.5)/6.5)*25
  if(v<30)   return 50+((v-25)/5)*25
  return Math.min(75+((v-30)/10)*25,100)
})
const bmr = computed(() => {
  if(!profile.value) return 0
  const {weight,height,age,gender}=profile.value
  return Math.round(gender==='male'
      ?10*weight+6.25*height-5*age+5
      :10*weight+6.25*height-5*age-161)
})
const dailyCalories = computed(() =>
    profile.value?Math.round(bmr.value*profile.value.activity_level):0
)
const waterPercentage = computed(() => {
  if(!profile.value) return 0
  const {gender,age,weight,height}=profile.value
  const liters=gender==='male'
      ?2.447+0.09156*weight+0.1074*height-0.09*age
      :-2.097+0.1069*height+0.2466*weight
  return Math.round((liters/weight)*100)
})
const boneWeight = computed(() =>
    profile.value?Math.round(profile.value.weight*(profile.value.gender==='male'?0.15:0.12)*10)/10:0
)
const musclePercentage = computed(() => {
  if(!profile.value) return 0
  const {gender,age}=profile.value, v=bmi.value
  let base=gender==='male'?45:38,
      ageAdj=Math.max(0,(age-25)*0.3),
      bmiAdj=v<18.5?-5:v>30?-8:v>25?-2:0
  return Math.max(25,Math.round(base-ageAdj+bmiAdj))
})
const bodyFatPercentage = computed(() => {
  if(!profile.value) return 0
  const {gender,age}=profile.value, v=bmi.value
  const genderF=gender==='male'?1:0
  return Math.max(5,Math.round(1.2*v+0.23*age-10.8*genderF-5.4))
})
const idealWeight = computed(() => {
  if(!profile.value) return {robinson:0,miller:0,devine:0}
  const {height,gender}=profile.value, inch=height/2.54
  let rob, mill, dev
  if(gender==='male'){
    rob=52+1.9*(inch-60)
    mill=56.2+1.41*(inch-60)
    dev=50+2.3*(inch-60)
  } else {
    rob=49+1.7*(inch-60)
    mill=53.1+1.36*(inch-60)
    dev=45.5+2.3*(inch-60)
  }
  return {
    robinson:Math.round(rob*10)/10,
    miller:  Math.round(mill*10)/10,
    devine:  Math.round(dev*10)/10
  }
})
const weightDifferenceText = computed(()=>{
  if(!profile.value) return ''
  const avg=(idealWeight.value.robinson+idealWeight.value.miller+idealWeight.value.devine)/3
  const diff=Math.round((profile.value.weight-avg)*10)/10
  if(Math.abs(diff)<2) return 'Idealna waga! 🎉'
  return diff>0?`${diff} kg powyżej idealnej`:`${Math.abs(diff)} kg poniżej idealnej`
})
const weightDifferenceClass = computed(()=>{
  if(!profile.value) return ''
  const avg=(idealWeight.value.robinson+idealWeight.value.miller+idealWeight.value.devine)/3
  const diff=profile.value.weight-avg
  if(Math.abs(diff)<2) return 'ideal'
  return diff>0?'above':'below'
})

const closeModal = () => showEditForm.value=false
// const openEditForm = () => {
//   editForm.value = profile.value ? {...profile.value} : {age:'',height:'',weight:'',gender:'',activity_level:'',goal:''}
//   showEditForm.value = true
// }
const openEditForm = () => {
  if (profile.value) {
    const { age, height, weight, gender, activity_level, goal } = profile.value
    editForm.value = { age, height, weight, gender, activity_level, goal }
  } else {
    editForm.value = { age: '', height: '', weight: '', gender: '', activity_level: '', goal: '' }
  }
  showEditForm.value = true
}
const saveProfile = async () => {
  saving.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) throw new Error('Brak tokenu autoryzacji')

    const payload = {
      age:            editForm.value.age,
      height:         editForm.value.height,
      weight:         editForm.value.weight,
      gender:         editForm.value.gender,
      activity_level: editForm.value.activity_level,
      goal:           editForm.value.goal
    }

    const res = await axios.post('/api/profile', payload, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (res.status === 200 || res.status === 201) {
      await loadProfile()
      closeModal()
      alert('Profil został zapisany!')
    }
  } catch (error) {
    let msg = error.response?.data?.error
        || error.response?.data?.message
        || error.message
        || 'Błąd zapisu profilu'
    alert(msg)
  } finally {
    saving.value = false
  }
}

const loadProfile = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      console.error('No token found')
      profile.value = null
      loading.value = false
      return
    }

    const res = await axios.get('/api/profile', {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    console.log('Profile data loaded:', res.data)
    profile.value = res.data.profile || null

  } catch (error) {
    console.error('Error loading profile:', error)
    profile.value = null
  }
  loading.value = false
}

const loadCurrentUser = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      currentUser.value = 'Użytkownik'
      return
    }

    // Try to get from auth store first
    if (authStore.user?.username) {
      currentUser.value = authStore.user.username
      return
    }

    // Fallback to API call
    const res = await axios.get('/api/me', {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    console.log('User data loaded:', res.data)
    currentUser.value = res.data.username || res.data.name || 'Użytkownik'

  } catch (error) {
    console.error('Error loading user:', error)
    currentUser.value = 'Użytkownik'
  }
}

onMounted(async () => {
  generateNewAvatar()
  await loadCurrentUser()
  await loadProfile()
})
</script>

<style scoped>
.profile-wrapper {
  padding: 1rem;
  background: rgba(139, 231, 139, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 0;
  min-height: 100vh;
  width: 100%;
  overflow: hidden;
}

.profile-view {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 1.5rem;
  width: 100%;
  min-height: calc(100vh - 2rem);
  max-height: calc(100vh - 2rem);
}

/* Profile Card */
.profile-card {
  background: rgba(255, 255, 255, 0.95);
  padding: 1.5rem;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  height: fit-content;
  max-height: 100%;
  overflow-y: auto;
}

.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 1.5rem;
}

.avatar-container {
  position: relative;
  margin-bottom: 1rem;
}

.profile-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px solid #8de78d;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.avatar-initials {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.refresh-avatar {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: none;
  background: #8de78d;
  cursor: pointer;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.refresh-avatar:hover {
  background: #6ed46e;
  transform: scale(1.1);
}

.profile-basic h2 {
  margin: 0;
  color: #2d5016;
  font-size: 1.4rem;
  text-align: center;
}

.profile-subtitle {
  color: #666;
  margin: 0.5rem 0 0 0;
  text-align: center;
  font-size: 0.85rem;
}

.profile-date {
  color: #999;
  font-size: 0.75rem;
  text-align: center;
  margin-top: 0.5rem;
}

.profile-details {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
  margin-bottom: 1.2rem;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.6rem;
  background: rgba(248, 252, 248, 0.8);
  border-radius: 10px;
  border: 1px solid #e8f5e8;
}

.detail-icon {
  font-size: 1rem;
  width: 20px;
  text-align: center;
}

.detail-label {
  font-weight: 500;
  color: #2d5016;
  min-width: 60px;
  font-size: 0.85rem;
}

.detail-value {
  color: #333;
  font-weight: 600;
  margin-left: auto;
  font-size: 0.85rem;
}

.no-profile {
  text-align: center;
  padding: 1.5rem;
}

.no-profile-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.8rem;
}

.no-profile-icon {
  font-size: 2.5rem;
}

.no-profile h3 {
  margin: 0;
  color: #2d5016;
  font-size: 1.1rem;
}

.no-profile p {
  color: #666;
  margin: 0;
  font-size: 0.9rem;
}

.create-profile-btn,
.edit-profile-btn {
  background: linear-gradient(135deg, #8de78d, #6ed46e);
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
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

.profile-actions {
  margin-top: 1rem;
  text-align: center;
}

.loading-state {
  text-align: center;
  padding: 1.5rem;
}

.loading-spinner {
  width: 30px;
  height: 30px;
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

/* Metrics Section - Fixed to fill screen properly */
.metrics-section {
  display: grid;
  grid-template-columns: 3fr 3fr;
  gap: 1.2rem;
  height: 100%;
  padding-right: 0.5rem;
  align-content: start;
}

.metrics-section::-webkit-scrollbar {
  width: 6px;
}

.metrics-section::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.metrics-section::-webkit-scrollbar-thumb {
  background: #8de78d;
  border-radius: 3px;
}

.metric-card {
  background: rgba(255, 255, 255, 0.95);
  padding: 1.2rem;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  height: fit-content;
  min-height: 200px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.2rem;
}

.metric-header h3 {
  margin: 0;
  color: #2d5016;
  font-size: 1rem;
}

/* BMI Card */
.bmi-value {
  font-size: 1.6rem;
  font-weight: 700;
  padding: 0.4rem 0.8rem;
  border-radius: 12px;
}

.bmi-value.underweight { background: #e3f2fd; color: #1976d2; }
.bmi-value.normal { background: #e8f5e8; color: #2e7d32; }
.bmi-value.overweight { background: #fff3e0; color: #f57c00; }
.bmi-value.obese { background: #ffebee; color: #d32f2f; }

.bmi-category {
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 0.4rem;
}

.bmi-category.underweight { color: #1976d2; }
.bmi-category.normal { color: #2e7d32; }
.bmi-category.overweight { color: #f57c00; }
.bmi-category.obese { color: #d32f2f; }

.bmi-range {
  color: #666;
  font-size: 0.8rem;
  margin-bottom: 0.8rem;
}

.bmi-progress {
  margin-top: 0.8rem;
}

.bmi-bar {
  height: 6px;
  background: linear-gradient(to right, #2196f3 25%, #4caf50 25% 50%, #ff9800 50% 75%, #f44336 75%);
  border-radius: 3px;
  position: relative;
  margin-bottom: 0.4rem;
}

.bmi-indicator {
  position: absolute;
  top: -3px;
  width: 12px;
  height: 12px;
  background: #333;
  border-radius: 50%;
  transform: translateX(-50%);
  transition: left 0.3s ease;
}

.bmi-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.7rem;
  color: #666;
}

/* Calories Card */
.calories-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: #ff6b35;
  background: #fff3f0;
  padding: 0.4rem 0.8rem;
  border-radius: 12px;
}

.calories-breakdown {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.calorie-item {
  display: flex;
  justify-content: space-between;
  padding: 0.4rem;
  background: #fafafa;
  border-radius: 8px;
  font-size: 0.85rem;
}

.calorie-label {
  color: #666;
}

.calorie-value {
  font-weight: 600;
  color: #333;
}

/* Body Composition Card */
.body-metrics {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.8rem;
}

.body-metric {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.8rem;
  background: #fafafa;
  border-radius: 12px;
  border: 1px solid #eee;
}

.metric-icon {
  font-size: 1.2rem;
}

.metric-info {
  flex: 1;
}

.metric-name {
  font-size: 0.8rem;
  color: #666;
  margin-bottom: 0.2rem;
}

.metric-val {
  font-weight: 600;
  color: #333;
  font-size: 0.95rem;
}

/* Fixed composition note text color */
.composition-note {
  margin-top: 1rem;
  padding: 0.8rem;
  background: rgba(248, 252, 248, 0.8);
  border-radius: 8px;
  border-left: 4px solid #8de78d;
}

.composition-note small {
  color: #2d5016 !important;
  line-height: 1.4;
}

.composition-note strong {
  color: #1a3309 !important;
  font-weight: 600;
}

/* Ideal Weight Card */
.ideal-weights {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.ideal-item {
  display: flex;
  justify-content: space-between;
  padding: 0.6rem;
  background: #fafafa;
  border-radius: 8px;
  font-size: 0.85rem;
}

.ideal-method {
  color: #666;
  font-weight: 500;
}

.ideal-value {
  font-weight: 600;
  color: #333;
}

.weight-difference {
  margin-top: 0.8rem;
  padding: 0.6rem;
  border-radius: 8px;
  text-align: center;
  font-weight: 600;
  font-size: 0.85rem;
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

.modal-content {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  max-width: 450px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.2rem;
}

.modal-header h3 {
  margin: 0;
  color: #2d5016;
  font-size: 1.2rem;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.3rem;
  cursor: pointer;
  color: #666;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
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
.form-group select {
  padding: 0.6rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 0.9rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #8de78d;
  box-shadow: 0 0 0 3px rgba(141, 231, 141, 0.2);
}

.form-actions {
  display: flex;
  gap: 0.8rem;
  margin-top: 0.8rem;
}

.cancel-btn {
  flex: 1;
  padding: 0.6rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9rem;
}

.save-btn {
  flex: 1;
  padding: 0.6rem;
  background: linear-gradient(135deg, #8de78d, #6ed46e);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Info tooltip */
.info-tooltip {
  font-size: 0.9rem;
  color: #666;
  cursor: help;
}

/* Responsive Design */
@media (max-width: 1024px) {
  .profile-view {
    grid-template-columns: 1fr;
    gap: 1.2rem;
    min-height: auto;
    max-height: none;
  }

  .profile-card {
    position: static;
  }

  .metrics-section {
    grid-template-columns: 1fr;
    max-height: none;
  }

  .profile-wrapper {
    min-height: auto;
    height: auto;
  }
}

@media (max-width: 768px) {
  .profile-wrapper {
    padding: 1rem;
    min-height: auto;
    height: auto;
  }

  .body-metrics {
    grid-template-columns: 1fr;
  }
}
</style>
