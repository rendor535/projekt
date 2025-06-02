<template>
  <div v-if="showSuccess" class="modal">
    <div class="modal-content">
      <p>✅ Pomyślnie zarejestrowano!</p>
      <button @click="goToPanel">OK</button>
    </div>
  </div>
  <form @submit.prevent="handleSubmit">
    <input v-model="form.username" type="text" placeholder="Nazwa użytkownika" required />

    <div class="password-wrapper">
      <input
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          placeholder="Hasło"
          required
      />
      <span class="toggle" @click="showPassword = !showPassword">
        {{ showPassword ? '🙈' : '👁️' }}
      </span>
    </div>

    <div v-if="props.type === 'register'" class="password-wrapper">
      <input
          :type="showPassword ? 'text' : 'password'"
          v-model="form.confirmPassword"
          placeholder="Powtórz hasło"
          required
      />
    </div>

    <div v-if="type === 'register'" class="role-buttons">
      <button
          type="button"
          :class="{ selected: form.role === 'Trainer' }"
          @click="form.role = 'Trainer'">
        Trener
      </button>
      <button
          type="button"
          :class="{ selected: form.role === 'Trainee' }"
          @click="form.role = 'Trainee'">
        Uczeń
      </button>
    </div>

    <button type="submit">{{ type === 'login' ? 'Zaloguj się' : 'Zarejestruj się' }}</button>
    <p v-if="error" class="error">{{ error }}</p>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import {useRouter }  from "vue-router";
const showSuccess = ref(false);
const router = useRouter()
const props = defineProps<{ type: 'login' | 'register' }>()
const form = ref({
  username: '',
  password: '',
  confirmPassword: '',
  role: ''
})
const error = ref('')
const showPassword = ref(false)

const store = useAuthStore()
function goToPanel() {
  if (showSuccess.value) {
    showSuccess.value = false
    router.push('/panel')
  }
}
async function handleSubmit() {
  error.value = ''
  try {
    if (props.type === 'register' && form.value.password !== form.value.confirmPassword) {
      error.value = 'Hasła nie są identyczne'
      return
    }

    if (props.type === 'login') {
      await store.login(form.value)
      router.push('/panel')
    } else {
      await store.register(form.value)
      showSuccess.value = true
      setTimeout(() => {
        router.push('/panel')
      }, 5000)
    }
  } catch (err: any) {
    if (err.response?.status === 401) {
      error.value = 'Nieprawidłowe dane logowania lub użytkownik nie istnieje.'
    } else {
      error.value = 'Wystąpił błąd. Spróbuj ponownie.'
    }
  }
}
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: white;
  padding: 2rem;
  border-radius: 8px;
  min-width: 300px;
}
input, button {
  padding: 0.75rem;
  font-size: 1rem;
}
button {
  background-color: #b7f5b7;
  border: none;
  cursor: pointer;
}
.error {
  color: red;
  font-size: 0.9rem;
}
.role-buttons {
  display: flex;
  gap: 1rem;
}
.role-buttons button {
  flex: 1;
  padding: 0.75rem;
  background-color: #e6ffe6;
  border: 1px solid #ccc;
  cursor: pointer;
  border-radius: 4px;
}
.role-buttons button.selected {
  background-color: #a9f3a9;
  border-color: #66cc66;
}

.password-wrapper {
  position: relative;
}
.password-wrapper input {
  width: 100%;
}
.toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  font-size: 1.2rem;
  user-select: none;
}
.modal {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  width: 100vw;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.modal-content {
  background: white;
  color: #181818;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
  text-align: center;
}

.modal-content p {
  font-size: 1.2rem;
  margin-bottom: 1rem;
}

.modal-content button {
  background-color: #b7f5b7;
  padding: 0.75rem 1.5rem;
  font-weight: bold;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

</style>
