<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import axios from "axios";

let onlinePingInterval: ReturnType<typeof setInterval> | null = null

async function setOnline() {
  try {
    await axios.post('/me/online')
    console.log('[DEBUG] online ping sent')
  } catch (err) {
    console.error('❌ setOnline failed:', err)
  }
}

async function setOffline() {
  try {
    await axios.post('/me/offline', {}, {
      headers: {
        // `keepalive` jest obsługiwane tylko w `fetch`, więc to symboliczne
      }
    })
    console.log('[DEBUG] offline ping sent')
  } catch (err) {
    console.error('❌ setOffline failed:', err)
  }
}

function startOnlineTracking() {
  setOnline()
  onlinePingInterval = setInterval(setOnline, 5_000)

  window.addEventListener('beforeunload', setOffline)

  document.addEventListener('visibilitychange', () => {
    console.log('[DEBUG] visibility:', document.visibilityState)
    if (document.visibilityState === 'hidden') {
      setOffline()
    } else if (document.visibilityState === 'visible') {
      setOnline()
    }
  })
}

onMounted(() => {
  console.log('[DEBUG] App.vue mounted')

  const token = localStorage.getItem('token')
  console.log('[DEBUG] token przy starcie:', token)

  // 🔁 Jeśli token istnieje od razu, odpal od razu
  if (token) {
    startOnlineTracking()
  } else {
    // ⏳ Jeśli token pojawi się później (np. po logowaniu), ustaw listener
    const tokenCheck = setInterval(() => {
      const freshToken = localStorage.getItem('token')
      if (freshToken) {
        console.log('[DEBUG] token wykryty później:', freshToken)
        clearInterval(tokenCheck)
        startOnlineTracking()
      }
    }, 500)
  }
})
</script>


<style>
/*global*/
</style>
