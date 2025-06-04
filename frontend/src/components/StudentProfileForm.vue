<template>
  <div class="profile-form-wrapper">
    <h2>Uzupełnij swój profil</h2>
    <form @submit.prevent="submitProfile">
      <label>Wiek:
        <input v-model.number="profile.age" type="number" min="1" required />
      </label>
      <label>Wzrost (cm):
        <input v-model.number="profile.height" type="number" min="50" required />
      </label>
      <label>Waga (kg):
        <input v-model.number="profile.weight" type="number" min="10" required />
      </label>
      <label>Płeć:
        <select v-model="profile.gender" required>
          <option value="male">Mężczyzna</option>
          <option value="female">Kobieta</option>
        </select>
      </label>
      <label>Poziom aktywności:
        <select v-model.number="profile.activity_level" required>
          <option value="1.2">Siedzący tryb życia (brak aktywności)</option>
          <option value="1.375">Lekko aktywny (1-3 dni w tygodniu)</option>
          <option value="1.55">Umiarkowanie aktywny (3-5 dni w tygodniu)</option>
          <option value="1.725">Bardzo aktywny (6-7 dni w tygodniu)</option>
          <option value="1.9">Super aktywny (2x dziennie, intensywnie)</option>
        </select>
      </label>
      <label>Cel:
        <select v-model="profile.goal" required>
          <option value="reduction">Redukcja</option>
          <option value="recomposition">Rekompozycja</option>
          <option value="mass">Masa</option>
        </select>
      </label>
      <button type="submit">Zapisz profil</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const profile = ref({
  age: 18,
  height: 170,
  weight: 70,
  gender: 'male',
  activity_level: 1.55, // Dodano domyślną wartość
  goal: 'reduction'
})

const emit = defineEmits(['saved'])

async function submitProfile() {
  try {
    console.log('Profil do wysłania:', JSON.stringify(profile.value)); // <— tutaj
    const response = await axios.post('/api/profile', profile.value);
    console.log('Odpowiedź:', response.data);
    emit('saved', profile.value);
  } catch (e) {
    console.error('Nie udało się zapisać profilu', e);
    console.error('Szczegóły błędu:', (e as any).response?.data);
  }
}
</script>

<style scoped>
.profile-form-wrapper {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  color: #181818;
}
form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
label {
  display: flex;
  flex-direction: column;
  font-weight: bold;
}
input, select {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 6px;
}
button {
  padding: 0.75rem;
  background: #b7f5b7;
  border: none;
  border-radius: 6px;
  font-weight: bold;
  cursor: pointer;
}
</style>
