import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import axios from 'axios'
import './assets/main.css'
axios.defaults.baseURL = 'http://localhost:8080/api'
const token = localStorage.getItem('token')
if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}
const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
