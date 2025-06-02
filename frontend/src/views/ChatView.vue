<template>
  <div class="chat-view">
    <!-- Header czatu -->
    <div class="chat-header">
      <h2>💬 Komunikacja</h2>
      <div class="connection-status" :class="{ connected: isConnected }">
        <div class="status-dot"></div>
        <span>{{ isConnected ? 'Online' : 'Offline' }}</span>
      </div>
    </div>

    <div class="chat-layout">
      <!-- Lewa kolumna - tylko dla trenera -->
      <div v-if="store.user?.role === 'trainer'" class="contacts-panel">
        <div class="contacts-header">
          <h3>🎯 Podopieczni</h3>
          <div class="contacts-count">{{ students.length }}</div>
        </div>

        <div class="contacts-list">
          <div
              v-for="student in students"
              :key="student.id"
              :class="['contact-item', { active: selectedStudentId === student.id }]"
              @click="selectStudent(student)"
          >
            <div class="contact-avatar">
              <span>{{ student.username.charAt(0).toUpperCase() }}</span>
              <div class="online-indicator" :class="{ online: student.is_online }"></div>
            </div>
            <div class="contact-info">
              <div class="contact-name">{{ student.username }}</div>
              <div class="contact-preview" v-if="student.lastMessage">
                {{ student.lastMessage.length > 25 ? student.lastMessage.substring(0, 25) + '...' : student.lastMessage }}
              </div>
              <div class="contact-status">{{ student.is_online ? 'Aktywny' : 'Nieaktywny' }}</div>
            </div>
            <div class="contact-meta">
              <div v-if="student.unreadCount > 0" class="unread-badge">
                {{ student.unreadCount }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Główny obszar czatu -->
      <div class="chat-container" :class="{ 'full-width': store.user?.role === 'trainee' }">
        <div v-if="(selectedStudent && store.user?.role === 'trainer') || (store.user?.role === 'trainee' && chatPartner)" class="chat-area">
          <!-- Chat Partner Info -->
          <div class="chat-partner-info">
            <div class="partner-avatar">
              <span>{{ chatPartner?.charAt(0).toUpperCase() }}</span>
              <div class="online-indicator" :class="{ online: isPartnerOnline }"></div>
            </div>
            <div class="partner-details">
              <h4>{{ chatPartner }}</h4>
              <span class="partner-status">{{ isPartnerOnline ? 'Aktywny teraz' : 'Nieaktywny' }}</span>
            </div>
            <div v-if="store.user?.role === 'trainee'" class="role-info">
              <span class="role-badge">Twój trener</span>
            </div>
          </div>

          <!-- Messages Area -->
          <div class="messages-area" ref="messagesContainer">
            <div v-if="loading" class="loading-state">
              <div class="loading-spinner"></div>
              <p>Ładowanie wiadomości...</p>
            </div>

            <div v-else-if="messages.length === 0" class="empty-state">
              <div class="empty-icon">💬</div>
              <h3>Rozpocznij konwersację</h3>
              <p v-if="store.user?.role === 'trainer'">Napisz pierwszą wiadomość do {{ chatPartner }}</p>
              <p v-else>Napisz pierwszą wiadomość do swojego trenera</p>
            </div>

            <div v-else class="messages-list">
              <div
                  v-for="message in messages"
                  :key="message.id"
                  :class="['message-item', { 'sent': message.sender_id === store.user.id }]"
              >
                <div class="message-bubble">
                  <p>{{ message.content }}</p>
                  <div class="message-footer">
                    <span class="message-time">{{ formatTime(message.created_at) }}</span>
                    <div v-if="message.sender_id === store.user.id" class="message-status">
                      <svg v-if="message.is_read" width="16" height="16" viewBox="0 0 24 24" fill="none" class="read-icon">
                        <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2"/>
                        <path d="m9 17 5-5" stroke="currentColor" stroke-width="2"/>
                      </svg>
                      <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" class="sent-icon">
                        <path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2"/>
                      </svg>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Message Input -->
          <div class="message-input-section">
            <div class="input-wrapper">
              <input
                  v-model="newMessage"
                  @keyup.enter="sendMessage"
                  :placeholder="store.user?.role === 'trainee' ? 'Napisz wiadomość do trenera...' : 'Wpisz wiadomość...'"
                  :disabled="sending"
                  ref="messageInput"
                  class="message-input"
              />
              <button
                  @click="sendMessage"
                  :disabled="!newMessage.trim() || sending"
                  class="send-btn"
              >
                <svg v-if="sending" width="20" height="20" viewBox="0 0 24 24" fill="none" class="spinner">
                  <path d="M21 12a9 9 0 11-6.219-8.56" stroke="currentColor" stroke-width="2"/>
                </svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none">
                  <path d="m22 2-7 20-4-9-9-4 20-7z" fill="currentColor"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="chat-empty-state">
          <div class="empty-content">
            <div class="empty-illustration">💬</div>
            <h3 v-if="store.user?.role === 'trainer'">Wybierz podopiecznego</h3>
            <h3 v-else>Ładowanie czatu...</h3>
            <p v-if="store.user?.role === 'trainer'">
              Kliknij na podopiecznego z listy po lewej stronie, aby rozpocząć rozmowę
            </p>
            <p v-else>
              Przygotowujemy czat z Twoim trenerem...
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const store = useAuthStore()

// Reactive data
const students = ref([])
const selectedStudent = ref(null)
const selectedStudentId = ref(null)
const messages = ref([])
const newMessage = ref('')
const loading = ref(false)
const sending = ref(false)
const chatPartner = ref('')
const isPartnerOnline = ref(false)
const messagesContainer = ref(null)
const messageInput = ref(null)
const isConnected = ref(true)

// Zmienne do śledzenia nowych wiadomości
const lastMessageId = ref(0)
const lastUnreadCheck = ref({})

// Interwały
let checkNewMessagesInterval = null
let checkUnreadInterval = null

// Pobierz listę uczniów (dla trenera) lub partnera czatu (dla ucznia)
const loadChatPartners = async () => {
  try {
    const response = await axios.get('/messages/partners')

    if (store.user?.role === 'trainer') {
      students.value = response.data.map(partner => ({
        ...partner,
        unreadCount: 0,
        lastMessage: ''
      }))

      // Załaduj dane tylko przy pierwszym ładowaniu
      for (let student of students.value) {
        await loadUnreadCount(student.id, student)
        await loadLastMessage(student.id, student)
      }
    } else {
      // Dla ucznia - znajdź trenera i automatycznie załaduj czat
      const trainer = response.data.find(p => p.role === 'trainer')
      if (trainer) {
        selectedStudent.value = trainer
        selectedStudentId.value = trainer.id
        chatPartner.value = trainer.username
        isPartnerOnline.value = trainer.is_online
        await loadConversation(trainer.id)

        // Automatyczne fokusowanie na pole tekstowe dla trainee
        await nextTick()
        if (messageInput.value) {
          messageInput.value.focus()
        }
      }
    }
    isConnected.value = true
  } catch (error) {
    console.error('Błąd przy ładowaniu partnerów czatu:', error)
    isConnected.value = false
  }
}

// Sprawdź czy są nowe wiadomości (tylko sprawdzenie, nie pełne ładowanie)
const checkForNewMessages = async () => {
  if (!selectedStudentId.value) return

  try {
    // Sprawdź tylko najnowszą wiadomość
    const response = await axios.get(`/messages/conversation/${selectedStudentId.value}?limit=1&after=${lastMessageId.value}`)

    if (response.data && response.data.length > 0) {
      const newMessages = response.data.filter(msg => msg.id > lastMessageId.value)

      if (newMessages.length > 0) {
        // Dodaj tylko nowe wiadomości
        for (let newMsg of newMessages) {
          messages.value.push(newMsg)
          lastMessageId.value = Math.max(lastMessageId.value, newMsg.id)
        }

        // Przewiń do dołu tylko jeśli dodano nowe wiadomości
        await nextTick()
        scrollToBottom()

        // Oznacz jako przeczytane jeśli to nie nasze wiadomości
        const unreadNewMessages = newMessages.filter(msg => msg.sender_id !== store.user.id && !msg.is_read)
        for (let msg of unreadNewMessages) {
          await markMessageAsRead(msg.id)
        }
      }
    }

    isConnected.value = true
  } catch (error) {
    console.error('Błąd przy sprawdzaniu nowych wiadomości:', error)
    isConnected.value = false
  }
}

// Sprawdź nieprzeczytane wiadomości (tylko dla trenera)
const checkUnreadMessages = async () => {
  if (store.user?.role !== 'trainer') return

  try {
    for (let student of students.value) {
      const response = await axios.get('/messages/unread', {
        params: { sender_id: student.id }
      })

      const currentCount = response.data.count || 0
      const lastCount = lastUnreadCheck.value[student.id] || 0

      if (currentCount !== lastCount) {
        student.unreadCount = currentCount
        lastUnreadCheck.value[student.id] = currentCount

        // Jeśli są nowe nieprzeczytane, zaktualizuj ostatnią wiadomość
        if (currentCount > lastCount) {
          await loadLastMessage(student.id, student)
        }
      }
    }
  } catch (error) {
    console.error('Błąd przy sprawdzaniu nieprzeczytanych wiadomości:', error)
  }
}

// Pobierz ostatnią wiadomość dla ucznia
const loadLastMessage = async (userId, student) => {
  try {
    const response = await axios.get(`/messages/conversation/${userId}?limit=1`)
    if (response.data && response.data.length > 0) {
      student.lastMessage = response.data[response.data.length - 1].content
    }
  } catch (error) {
    console.error('Błąd przy ładowaniu ostatniej wiadomości:', error)
  }
}

// Pobierz liczbę nieprzeczytanych wiadomości
const loadUnreadCount = async (userId, student = null) => {
  try {
    const response = await axios.get('/messages/unread', {
      params: { sender_id: userId }
    })

    if (student) {
      student.unreadCount = response.data.count || 0
      lastUnreadCheck.value[userId] = student.unreadCount
    }
  } catch (error) {
    console.error('Błąd przy ładowaniu nieprzeczytanych wiadomości:', error)
  }
}

// Wybierz ucznia (dla trenera)
const selectStudent = async (student) => {
  selectedStudent.value = student
  selectedStudentId.value = student.id
  chatPartner.value = student.username
  isPartnerOnline.value = student.is_online
  await axios.put(`/messages/conversation/${student.id}/read`)
  // Oznacz wiadomości jako przeczytane
  await markMessagesAsRead(student.id)

  // Załaduj konwersację
  await loadConversation(student.id)

  // Resetuj licznik nieprzeczytanych
  student.unreadCount = 0
  lastUnreadCheck.value[student.id] = 0

  // Fokus na pole tekstowe
  nextTick(() => {
    if (messageInput.value) {
      messageInput.value.focus()
    }
  })
}

// Załaduj konwersację z określonym użytkownikiem (pełne ładowanie)
const loadConversation = async (partnerId) => {
  loading.value = true
  try {
    const response = await axios.get(`/messages/conversation/${partnerId}`)
    messages.value = response.data || []

    // Ustaw ID ostatniej wiadomości
    if (messages.value.length > 0) {
      lastMessageId.value = Math.max(...messages.value.map(m => m.id))
    }

    // Przewiń do najnowszej wiadomości
    await nextTick()
    scrollToBottom()
    await axios.put(`/messages/conversation/${partnerId}/read`)
    localStorage.setItem('messagesCleared', 'true')
    isConnected.value = true
  } catch (error) {
    console.error('Błąd przy ładowaniu konwersacji:', error)
    isConnected.value = false
    messages.value = []
    lastMessageId.value = 0
  } finally {
    loading.value = false
  }
}

// Wyślij wiadomość
const sendMessage = async () => {
  if (!newMessage.value.trim() || sending.value) return

  const receiverId = store.user?.role === 'trainer' ? selectedStudentId.value : selectedStudent.value?.id
  if (!receiverId) return

  sending.value = true
  const messageText = newMessage.value.trim()

  try {
    const response = await axios.post('/messages', {
      receiver_id: receiverId,
      content: messageText
    })

    // Dodaj wiadomość do listy
    const newMsg = {
      id: response.data.id,
      sender_id: store.user.id,
      receiver_id: receiverId,
      content: messageText,
      created_at: new Date().toISOString(),
      is_read: false
    }

    messages.value.push(newMsg)
    lastMessageId.value = response.data.id

    // Wyczyść pole tekstowe
    newMessage.value = ''

    // Aktualizuj ostatnią wiadomość w liście uczniów (tylko dla trenera)
    if (store.user?.role === 'trainer') {
      const student = students.value.find(s => s.id === receiverId)
      if (student) {
        student.lastMessage = messageText
      }
    }

    // Przewiń do najnowszej wiadomości
    await nextTick()
    scrollToBottom()

    isConnected.value = true
  } catch (error) {
    console.error('Błąd przy wysyłaniu wiadomości:', error)
    isConnected.value = false
    alert('Nie udało się wysłać wiadomości')
  } finally {
    sending.value = false
  }
}

// Oznacz jedną wiadomość jako przeczytaną
const markMessageAsRead = async (messageId) => {
  try {
    await axios.put(`/messages/${messageId}/read`)

    // Aktualizuj lokalnie
    const message = messages.value.find(m => m.id === messageId)
    if (message) {
      message.is_read = true
    }
  } catch (error) {
    console.error('Błąd przy oznaczaniu wiadomości jako przeczytana:', error)
  }
}

// Oznacz wiadomości jako przeczytane (wszystkie od nadawcy)
const markMessagesAsRead = async (senderId) => {
  try {
    const unreadMessages = messages.value.filter(
        msg => msg.sender_id === senderId && !msg.is_read
    )

    for (let message of unreadMessages) {
      await markMessageAsRead(message.id)
    }
  } catch (error) {
    console.error('Błąd przy oznaczaniu wiadomości jako przeczytane:', error)
  }
}

// POPRAWIONA funkcja formatTime
const formatTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()

  // Sprawdź czy data jest prawidłowa
  if (isNaN(date.getTime())) {
    return 'Nieprawidłowa data'
  }

  // Oblicz różnicę w milisekundach (może być ujemna dla przyszłych dat)
  const diffTime = now - date

  // Jeśli wiadomość jest z przyszłości (błąd zegara/strefy czasowej)
  if (diffTime < 0) {
    return 'Teraz'
  }

  // Konwertuj na minuty, godziny, dni
  const diffMinutes = Math.floor(diffTime / (1000 * 60))
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60))
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

  // Logika formatowania
  if (diffMinutes < 1) {
    return 'Teraz'
  } else if (diffMinutes < 60) {
    return `${diffMinutes} min temu`
  } else if (diffHours < 24) {
    return `${diffHours}h temu`
  } else if (diffDays === 1) {
    return `Wczoraj ${date.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' })}`
  } else if (diffDays <= 7) {
    return date.toLocaleDateString('pl-PL', { weekday: 'short', hour: '2-digit', minute: '2-digit' })
  } else {
    return date.toLocaleDateString('pl-PL', { day: '2-digit', month: '2-digit' })
  }
}

// Przewiń do dołu kontenera wiadomości
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// Obserwuj zmiany w selectedStudentId
watch(selectedStudentId, (newId) => {
  if (newId) {
    loadConversation(newId)
  }
})

// Załaduj dane przy montowaniu komponentu
onMounted(() => {
  loadChatPartners()

  // Sprawdzaj nowe wiadomości co 5 sekund (tylko sprawdzenie, nie pełne ładowanie)
  checkNewMessagesInterval = setInterval(checkForNewMessages, 5000)

  // Sprawdzaj nieprzeczytane wiadomości co 10 sekund (tylko dla trenera)
  if (store.user?.role === 'trainer') {
    checkUnreadInterval = setInterval(checkUnreadMessages, 10000)
  }
})

// Wyczyść interwały przy niszczeniu komponentu
onUnmounted(() => {
  if (checkNewMessagesInterval) {
    clearInterval(checkNewMessagesInterval)
  }
  if (checkUnreadInterval) {
    clearInterval(checkUnreadInterval)
  }
})
</script>

<style scoped>
.chat-view {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  width: 95%;
  height: 90%;
  margin: 2rem;
  color: #181818;
  display: flex;
  flex-direction: column;
}

.chat-header {
  text-align: center;
  margin-bottom: 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header h2 {
  color: #2d5016;
  margin: 0;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #666;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d32f2f;
}

.connection-status.connected .status-dot {
  background: #2e7d32;
}

.chat-layout {
  flex: 1;
  display: flex;
  gap: 1rem;
  min-height: 0;
}

/* Lewa kolumna - kontakty */
.contacts-panel {
  width: 300px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
}

.contacts-header {
  padding: 1rem;
  background: #e8f5e8;
  border-bottom: 1px solid #c8e6c9;
  border-radius: 8px 8px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.contacts-header h3 {
  color: #2d5016;
  margin: 0;
  font-size: 1rem;
}

.contacts-count {
  background: #2d5016;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: bold;
  min-width: 20px;
  text-align: center;
}

.contacts-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  margin-bottom: 0.5rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.contact-item:hover {
  border-color: #b7f5b7;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.contact-item.active {
  border-color: #2d5016;
  background: #f0f8f0;
  box-shadow: 0 2px 8px rgba(45,80,22,0.2);
}

.contact-avatar {
  position: relative;
  width: 40px;
  height: 40px;
  background: #2d5016;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.online-indicator {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 12px;
  height: 12px;
  background: #999;
  border: 2px solid white;
  border-radius: 50%;
}

.online-indicator.online {
  background: #2e7d32;
}

.contact-info {
  flex: 1;
  min-width: 0;
}

.contact-name {
  font-weight: 600;
  color: #2d5016;
  margin-bottom: 0.25rem;
}

.contact-preview {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.25rem;
}

.contact-status {
  font-size: 0.75rem;
  color: #999;
}

.contact-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.unread-badge {
  background: #d32f2f;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  padding: 0.125rem 0.375rem;
  border-radius: 12px;
  min-width: 20px;
  text-align: center;
}

/* Główny obszar czatu */
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  min-height: 0;
}

.chat-container.full-width {
  width: 100%;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-partner-info {
  padding: 1rem;
  background: #e8f5e8;
  border-bottom: 1px solid #c8e6c9;
  border-radius: 8px 8px 0 0;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.partner-avatar {
  position: relative;
  width: 48px;
  height: 48px;
  background: #2d5016;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.125rem;
  flex-shrink: 0;
}

.partner-details {
  flex: 1;
}

.partner-details h4 {
  color: #2d5016;
  margin: 0 0 0.25rem 0;
}

.partner-status {
  font-size: 0.875rem;
  color: #666;
}

.role-info {
  display: flex;
  align-items: center;
}

.role-badge {
  background: #2d5016;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

/* Messages Area */
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  background: white;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #2d5016;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: #666;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h3 {
  color: #2d5016;
  margin: 0 0 0.5rem 0;
}

.empty-state p {
  margin: 0;
}

/* Messages List */
.messages-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message-item {
  display: flex;
  max-width: 75%;
}

.message-item.sent {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-bubble {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 18px;
  padding: 0.75rem 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  max-width: 100%;
}

.message-item.sent .message-bubble {
  background: #e8f5e8;
  border-color: #c8e6c9;
}

.message-bubble p {
  margin: 0 0 0.5rem 0;
  line-height: 1.4;
  color: #333;
  word-wrap: break-word;
}

.message-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 0.75rem;
  margin-top: 0.25rem;
}

.message-time {
  color: #999;
}

.message-status {
  display: flex;
  align-items: center;
}

.read-icon {
  color: #2e7d32;
}

.sent-icon {
  color: #999;
}

/* Message Input */
.message-input-section {
  padding: 1rem;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
  border-radius: 0 0 8px 8px;
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 24px;
  padding: 0.5rem;
}

.input-wrapper:focus-within {
  border-color: #2d5016;
  box-shadow: 0 0 0 2px rgba(45,80,22,0.2);
}

.message-input {
  flex: 1;
  border: none;
  outline: none;
  padding: 0.75rem 1rem;
  font-size: 0.9375rem;
  color: #333;
  background: transparent;
}

.message-input::placeholder {
  color: #999;
}

.send-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: #2d5016;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: #1e3a0f;
  transform: translateY(-1px);
}

.send-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.spinner {
  animation: spin 1s linear infinite;
}

/* Empty Chat State */
.chat-empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 8px;
}

.empty-content {
  text-align: center;
  color: #666;
  padding: 2rem;
}

.empty-illustration {
  font-size: 4rem;
  margin-bottom: 1.5rem;
  opacity: 0.5;
}

.empty-content h3 {
  color: #2d5016;
  margin: 0 0 1rem 0;
}

.empty-content p {
  margin: 0;
  line-height: 1.5;
}

/* Scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #c8e6c9;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #b7f5b7;
}

/* Responsive */
@media (max-width: 768px) {
  .chat-view {
    padding: 1rem;
    margin: 1rem;
  }

  .chat-layout {
    flex-direction: column;
  }

  .contacts-panel {
    width: 100%;
    max-height: 250px;
  }

  .message-item {
    max-width: 90%;
  }
}
</style>