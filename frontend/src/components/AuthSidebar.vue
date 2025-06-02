<template>
  <div class="sidebar">
    <div class="logo">
      <span>🏋️</span>
    </div>

    <div class="sidebar-center">
      <button :class="{ active: currentForm === 'login' }" @click="$emit('switch', 'login')">Zaloguj się</button>
      <button :class="{ active: currentForm === 'register' }" @click="$emit('switch', 'register')">Zarejestruj się</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  currentForm: 'login' | 'register'
}>()
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  width: 220px;
  height: 100vh;
  padding: 2rem 1rem;

  /* 🔹 prawdziwy efekt blur */
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  background-color: rgba(255, 255, 255, 0.15);

  /* opcjonalny cien i delikatna ramka */
  box-shadow: inset -1px 0 0 rgba(0, 0, 0, 0.05);
  border-right: 1px solid rgba(255, 255, 255, 0.2);
}
.sidebar-center {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  flex-grow: 1;
}
button {
  background-color: #b7f5b7;
  border: none;
  padding: 1rem 2rem;
  margin: 0.5rem 0;
  font-weight: bold;
  font-size: 1rem;
  white-space: nowrap; /* ❗ ZAPOBIEGA ZAWIJANIU TEKSTU */
  cursor: pointer;
  border-radius: 8px;
  width: 100%;
  max-width: 180px;
  transition: background-color 0.2s ease;
}

button.active {
  background-color: #8de78d;
}
.logo {
  text-align: center;
  font-size: 8.8rem;
  font-weight: 700;
  font-style: italic;
  color: #398dff;
  margin-bottom: 3rem;
  line-height: 1.2;
  font-family: 'Arial Black', cursive;
  text-shadow: 0 0 6px #ffffff, 0 0 6px #ffffff;
}
</style>

