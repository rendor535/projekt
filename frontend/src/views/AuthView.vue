<template>
  <div class="auth-page">
    <template v-if="route.path === '/'">
      <div class="background" />
      <div class="overlay" />
    </template>

    <aside class="sidebar">
      <AuthSidebar :currentForm="currentForm" @switch="currentForm = $event" />
    </aside>

    <main class="form-area">
      <div class="form-card">
        <AuthForm :type="currentForm" />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
const route = useRoute()
import { ref } from 'vue'
import AuthSidebar from '@/components/AuthSidebar.vue'
import AuthForm from '@/components/AuthForm.vue'

const currentForm = ref<'login' | 'register'>('login')
</script>

<style scoped>
/* RESET SCROLLA */
html, body {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

/* STRONA */
.auth-page {
  position: relative;
  width: 100vw;
  height: 100vh;
  display: flex;
  overflow: hidden;
}

/* TŁO */
.background,
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  width: 100%;
}

.background {
  background-image: url('@/assets/bg.jpg');
  background-size: cover;
  background-position: center;
  z-index: -2;
}

.overlay {
  backdrop-filter: blur(2px);
  background-color: rgba(255, 255, 255, 0.03);
  z-index: -1;
}

/* OBSZAR FORMULARZA */
.form-area {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* WYSPA */
.form-card {
  background-color: #ffffff;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 400px;
}
</style>
