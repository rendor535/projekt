<template>
  <div class="media-view">
    <div class="page-header">
      <div class="header-content">
        <h1>🎥 Media</h1>
        <p v-if="store.user?.role === 'trainer'">
          Udostępniaj materiały wideo i pliki swoim podopiecznym
        </p>
        <p v-else>
          Materiały udostępnione przez trenera
        </p>
      </div>
    </div>

    <!-- Panel trenera -->
    <div v-if="store.user?.role === 'trainer'" class="trainer-panel">
      <div class="panel-header">
        <h2>Panel trenera</h2>
        <button @click="showAddForm = !showAddForm" class="btn-primary">
          {{ showAddForm ? 'Anuluj' : 'Dodaj nowy film' }}
        </button>
      </div>

      <!-- Formularz dodawania filmu -->
      <div v-if="showAddForm" class="add-form">
        <h3>Przypisz film YouTube</h3>

        <!-- Wybór typu dodawania -->
        <div class="form-group">
          <div class="video-type-selector">
            <label class="radio-option">
              <input type="radio" v-model="videoType" value="predefined" />
              <span>📚 Wybierz z gotowych filmów</span>
            </label>
            <label class="radio-option">
              <input type="radio" v-model="videoType" value="custom" />
              <span>➕ Dodaj własny link YouTube</span>
            </label>
          </div>
        </div>

        <form @submit.prevent="addMedia">
          <!-- Sekcja gotowych filmów -->
          <div v-if="videoType === 'predefined'" class="predefined-section">
            <div class="form-group">
              <label>Wybierz gotowy film:</label>
              <div class="search-box">
                <input
                    v-model="searchQuery"
                    type="text"
                    placeholder="Szuka/optj ćwiczenia..."
                    class="search-input"
                >
              </div>

              <div v-if="loading" class="loading-small">Ładowanie filmów...</div>

              <div v-else class="video-selector">
                <div
                    v-for="video in filteredExerciseVideos"
                    :key="video.id"
                    class="video-option"
                    :class="{ active: selectedExerciseVideo?.id === video.id }"
                    @click="selectExerciseVideo(video)"
                >
                  <div class="video-preview-small">
                    <iframe
                        :src="getEmbedUrl(video.youtube_url)"
                        frameborder="0"
                        class="mini-iframe"
                    ></iframe>
                  </div>
                  <div class="video-details">
                    <h4>{{ video.name }}</h4>
                    <span class="video-date">{{ formatDate(video.created_at) }}</span>
                  </div>
                  <div class="selection-indicator">
                    <span v-if="selectedExerciseVideo?.id === video.id">✅</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Sekcja własnego linku -->
          <div v-if="videoType === 'custom'" class="custom-section">
            <div class="form-group">
              <label>Nazwa filmu:</label>
              <input v-model="newMedia.file_name" type="text" required placeholder="Wprowadź nazwę filmu">
            </div>

            <div class="form-group">
              <label>Link YouTube:</label>
              <input v-model="newMedia.file_url" type="url" required placeholder="https://www.youtube.com/watch?v=...">
            </div>

            <div class="form-group">
              <label>Opis:</label>
              <textarea v-model="newMedia.description" placeholder="Opcjonalny opis filmu"></textarea>
            </div>
          </div>

          <!-- Wspólna sekcja przypisania -->
          <div class="form-group">
            <label>Przypisz do użytkownika:</label>
            <select v-model="selectedRecipientId" required>
              <option value="">Wybierz użytkownika</option>
              <option v-for="user in users" :key="user.id" :value="user.id">
                {{ user.username }}
              </option>
            </select>
          </div>

          <div class="form-actions">
            <button
                type="submit"
                class="btn-success"
                :disabled="loading || !canSubmit"
            >
              {{ loading ? 'Dodawanie...' : 'Przypisz film' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Lista przesłanych mediów przez trenera -->
      <div class="uploaded-media">
        <h3>Twoje przesłane filmy</h3>
        <div v-if="uploadedMedia.length === 0" class="no-media">
          Nie przesłałeś jeszcze żadnych filmów.
        </div>
        <div v-else class="media-grid">
          <div v-for="media in uploadedMedia" :key="media.id" class="media-card trainer">
            <div class="media-header">
              <h4>{{ media.file_name }}</h4>
              <button @click="deleteMedia(media.id)" class="btn-danger-small">
                🗑️ Usuń
              </button>
            </div>
            <div class="media-info">
              <p><strong>Przypisano do:</strong> {{ media.recipient_username }}</p>
              <p><strong>Status:</strong>
                <span :class="media.is_read ? 'status-read' : 'status-unread'">
                  {{ media.is_read ? 'Obejrzane' : 'Nieobejrzane' }}
                </span>
              </p>
              <p><strong>Data dodania:</strong> {{ formatDate(media.created_at) }}</p>
              <p v-if="media.description"><strong>Opis:</strong> {{ media.description }}</p>
            </div>
            <div class="youtube-preview">
              <iframe
                  :src="media.file_url"
                  frameborder="0"
                  allowfullscreen
                  class="youtube-iframe"
              ></iframe>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Panel uczestnika -->
    <div v-else class="participant-panel">
      <div class="panel-header">
        <h2>Twoje filmy od trenera</h2>
        <span v-if="unreadCount > 0" class="unread-badge">
          {{ unreadCount }} nowych
        </span>
      </div>

      <div v-if="sharedMedia.length === 0" class="no-media">
        Trener nie przesłał Ci jeszcze żadnych filmów.
      </div>

      <div v-else class="media-grid">
        <div v-for="media in sharedMedia" :key="media.id" class="media-card participant">
          <div class="media-header">
            <h4>{{ media.file_name }}</h4>
            <span v-if="!media.is_read" class="new-badge">NOWY</span>
          </div>

          <div class="media-info">
            <p><strong>Od:</strong> {{ media.uploader_username }}</p>
            <p><strong>Data:</strong> {{ formatDate(media.created_at) }}</p>
            <p v-if="media.description"><strong>Opis:</strong> {{ media.description }}</p>
          </div>

          <div class="youtube-player">
            <iframe
                :src="media.file_url"
                frameborder="0"
                allowfullscreen
                class="youtube-iframe"
                @load="markAsRead(media.id)"
            ></iframe>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const store = useAuthStore()
const loading = ref(false)
const showAddForm = ref(false)

// Nowe dane dla wyboru typu filmu
const videoType = ref('predefined') // 'predefined' lub 'custom'
const exerciseVideos = ref([])
const selectedExerciseVideo = ref(null)
const searchQuery = ref('')
const selectedRecipientId = ref('')

// Istniejące dane
const users = ref([])
const sharedMedia = ref([])
const uploadedMedia = ref([])
const unreadCount = ref(0)

// Formularz nowego filmu
const newMedia = ref({
  file_name: '',
  file_url: '',
  file_type: 'video',
  description: '',
  recipient_id: ''
})

// Computed properties
const filteredExerciseVideos = computed(() => {
  if (!searchQuery.value) return exerciseVideos.value

  const query = searchQuery.value.toLowerCase()
  return exerciseVideos.value.filter(video =>
      video.name.toLowerCase().includes(query)
  )
})

const canSubmit = computed(() => {
  if (!selectedRecipientId.value) return false

  if (videoType.value === 'predefined') {
    return selectedExerciseVideo.value !== null
  } else {
    return newMedia.value.file_name && newMedia.value.file_url
  }
})

// Funkcje pomocnicze
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getEmbedUrl = (url) => {
  const videoId = extractYouTubeID(url)
  return videoId ? `https://www.youtube.com/embed/${videoId}` : url
}

const extractYouTubeID = (url) => {
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/,
    /youtube\.com\/v\/([^&\n?#]+)/
  ]

  for (const pattern of patterns) {
    const matches = url.match(pattern)
    if (matches && matches[1]) {
      return matches[1]
    }
  }
  return null
}

const selectExerciseVideo = (video) => {
  selectedExerciseVideo.value = video
}

const resetForm = () => {
  newMedia.value = {
    file_name: '',
    file_url: '',
    file_type: 'video',
    description: '',
    recipient_id: ''
  }
  selectedExerciseVideo.value = null
  selectedRecipientId.value = ''
  searchQuery.value = ''
  videoType.value = 'predefined'
  showAddForm.value = false
}

// API request helper
const apiRequest = async (url, options = {}) => {
  const token = store.token || localStorage.getItem('token')

  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  }

  const finalOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers
    }
  }

  const response = await fetch(`http://188.245.239.124:8080/api${url}`, finalOptions)

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return response.json()
}

// API calls
const fetchExerciseVideos = async () => {
  try {
    console.log('🎬 Rozpoczynam pobieranie filmów...')
    console.log('🎬 Token:', store.token || localStorage.getItem('token'))
    console.log('🎬 User role:', store.user?.role)

    const data = await apiRequest('/media/exercise-videos')
    console.log('🎬 Otrzymane dane:', data)
    console.log('🎬 Typ danych:', typeof data)
    console.log('🎬 Czy to tablica?:', Array.isArray(data))

    exerciseVideos.value = (data || []).map(video => ({
      ...video,
      selectedRecipient: ''
    }))

    console.log('🎬 Przetworzone filmy:', exerciseVideos.value)
  } catch (error) {
    console.error('❌ Błąd pobierania filmów:', error)
    console.error('❌ Error details:', error.message)
    alert('Błąd podczas pobierania biblioteki filmów: ' + error.message)
  }
}

const fetchUsers = async () => {
  try {
    const data = await apiRequest('/media/users')
    users.value = data
  } catch (error) {
    console.error('Error fetching users:', error)
  }
}

const fetchSharedMedia = async () => {
  try {
    const data = await apiRequest('/media/shared')
    sharedMedia.value = data.media || []
    unreadCount.value = data.unread_count || 0
  } catch (error) {
    console.error('Error fetching shared media:', error)
  }
}

const fetchUploadedMedia = async () => {
  try {
    const data = await apiRequest('/media/uploaded')
    uploadedMedia.value = data || []
  } catch (error) {
    console.error('Error fetching uploaded media:', error)
  }
}

const addMedia = async () => {
  loading.value = true
  try {
    if (videoType.value === 'predefined') {
      // Przypisz gotowy film
      await apiRequest('/media/assign-exercise-video', {
        method: 'POST',
        body: JSON.stringify({
          exercise_video_id: selectedExerciseVideo.value.id,
          recipient_id: selectedRecipientId.value
        })
      })
    } else {
      // Dodaj własny film
      await apiRequest('/media/upload', {
        method: 'POST',
        body: JSON.stringify({
          ...newMedia.value,
          recipient_id: selectedRecipientId.value
        })
      })
    }

    // Reset formularza
    resetForm()

    // Odśwież listę
    await fetchUploadedMedia()

    alert('Film został pomyślnie przypisany!')
  } catch (error) {
    console.error('Error adding media:', error)
    alert('Błąd podczas dodawania filmu. Sprawdź dane i spróbuj ponownie.')
  } finally {
    loading.value = false
  }
}

const deleteMedia = async (mediaId) => {
  if (!confirm('Czy na pewno chcesz usunąć ten film?')) return

  try {
    await apiRequest(`/media/${mediaId}`, {
      method: 'DELETE'
    })
    await fetchUploadedMedia()
    alert('Film został usunięty!')
  } catch (error) {
    console.error('Error deleting media:', error)
    alert('Błąd podczas usuwania filmu.')
  }
}

const markAsRead = async (mediaId) => {
  try {
    await apiRequest(`/media/${mediaId}/read`, {
      method: 'PATCH'
    })

    // Aktualizacja lokalna = UX bez opóźnienia
    const mediaItem = sharedMedia.value.find(m => m.id === mediaId)
    if (mediaItem) {
      mediaItem.is_read = true
    }

    unreadCount.value = sharedMedia.value.filter(m => !m.is_read).length
  } catch (error) {
    console.error('Error marking as read:', error)
  }
}

// Lifecycle
onMounted(async () => {
  if (store.user?.role === 'trainer') {
    await fetchUsers()
    await fetchUploadedMedia()
    await fetchExerciseVideos() // Pobierz gotowe filmy
  } else {
    await fetchSharedMedia()
  }
})
</script>

<style scoped>
.media-view {
  margin-left: 1.5rem;
  margin-top: 1.5rem;
  background: white;
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);
  max-width: 95%;
  width: 100%;
  color: #181818;
  overflow-y: auto;
  max-height: calc(100vh - 4rem);
}

.page-header {
  text-align: center;
  margin-bottom: 3rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #f0f0f0;
}

.header-content h1 {
  color: #2d5016;
  margin: 0 0 1rem 0;
  font-size: 2.5rem;
}

.header-content p {
  color: #666;
  font-size: 1.1rem;
  margin: 0;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.panel-header h2 {
  color: #2d5016;
  margin: 0;
}

.unread-badge {
  background: #e74c3c;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: bold;
}

.btn-primary {
  background: #2d5016;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.3s;
}

.btn-primary:hover {
  background: #1e3510;
}

.btn-success {
  background: #27ae60;
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  transition: background 0.3s;
}

.btn-success:hover:not(:disabled) {
  background: #219a52;
}

.btn-success:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.btn-danger-small {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background 0.3s;
}

.btn-danger-small:hover {
  background: #c0392b;
}

.add-form {
  background: #f8f9fa;
  padding: 2rem;
  border-radius: 12px;
  margin-bottom: 3rem;
}

.add-form h3 {
  color: #2d5016;
  margin-bottom: 1.5rem;
}

/* Nowe style dla wyboru typu filmu */
.video-type-selector {
  display: flex;
  gap: 2rem;
  margin-bottom: 1.5rem;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.75rem 1rem;
  border: 2px solid #ddd;
  border-radius: 8px;
  transition: all 0.3s;
}

.radio-option:hover {
  border-color: #2d5016;
}

.radio-option input[type="radio"] {
  margin: 0;
  width: auto;
}

.radio-option input[type="radio"]:checked + span {
  color: #2d5016;
  font-weight: bold;
}

.search-box {
  margin-bottom: 1rem;
}

.search-input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  box-sizing: border-box;
}

.search-input:focus {
  outline: none;
  border-color: #2d5016;
}

.loading-small {
  text-align: center;
  padding: 1rem;
  color: #666;
  font-style: italic;
}

.video-selector {
  max-height: 300px;
  overflow-y: auto;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 0.5rem;
}

.video-option {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 0.5rem;
}

.video-option:hover {
  background: #f8f9fa;
  border-color: #ddd;
}

.video-option.active {
  background: #e7f3e7;
  border-color: #2d5016;
}

.video-preview-small {
  flex-shrink: 0;
  width: 120px;
  height: 68px;
}

.mini-iframe {
  width: 100%;
  height: 100%;
  border-radius: 4px;
}

.video-details {
  flex: 1;
}

.video-details h4 {
  margin: 0 0 0.25rem 0;
  color: #2d5016;
  font-size: 1rem;
}

.video-date {
  color: #666;
  font-size: 0.85rem;
}

.selection-indicator {
  flex-shrink: 0;
  width: 30px;
  text-align: center;
  font-size: 1.2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  color: #333;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #2d5016;
}

.form-group textarea {
  height: 100px;
  resize: vertical;
}

.form-actions {
  text-align: center;
  margin-top: 2rem;
}

.no-media {
  text-align: center;
  padding: 3rem;
  color: #666;
  font-size: 1.1rem;
  background: #f8f9fa;
  border-radius: 12px;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.media-card {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s;
}

.media-card:hover {
  border-color: #2d5016;
  box-shadow: 0 4px 16px rgba(45, 80, 22, 0.1);
}

.media-card.trainer {
  border-left: 4px solid #2d5016;
}

.media-card.participant {
  border-left: 4px solid #3498db;
}

.media-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.media-header h4 {
  color: #2d5016;
  margin: 0;
  font-size: 1.2rem;
}

.new-badge {
  background: #e74c3c;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: bold;
}

.media-info {
  margin-bottom: 1.5rem;
}

.media-info p {
  margin: 0.5rem 0;
  color: #555;
}

.status-read {
  color: #27ae60;
  font-weight: bold;
}

.status-unread {
  color: #e74c3c;
  font-weight: bold;
}

.youtube-player,
.youtube-preview {
  width: 100%;
  margin-top: 1rem;
}

.youtube-iframe {
  width: 100%;
  height: 250px;
  border-radius: 8px;
}

.uploaded-media h3 {
  color: #2d5016;
  margin-bottom: 1.5rem;
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

@media (max-width: 768px) {
  .video-type-selector {
    flex-direction: column;
    gap: 1rem;
  }

  .video-option {
    flex-direction: column;
    text-align: center;
  }

  .video-preview-small {
    width: 100%;
    height: 150px;
  }

  .media-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .panel-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .youtube-iframe {
    height: 200px;
  }
}
</style>
