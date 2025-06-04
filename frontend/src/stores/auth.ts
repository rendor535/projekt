// Zaktualizowany auth store z dodatkowymi getterami dla training plans
import { defineStore } from 'pinia'
import axios from 'axios'
axios.defaults.baseURL = '/api'
const token = localStorage.getItem('token')
if(token) {
	axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}
export interface User {
    id: number;
    username: string;
    role: string;
}

export interface Profile {
    age: number;
    height: number;
    weight: number;
    gender: string;
    activity_level: number;
    goal: string;
}

export const useAuthStore = defineStore('auth', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        user: null as null | User,
        profile: null as null | Profile
    }),

    getters: {
        // Gettery dla kompatybilności z komponentami training plans
        isAuthenticated: (state): boolean => !!state.token && !!state.user,
        isTrainer: (state): boolean => state.user?.role === 'trainer',
        isTrainee: (state): boolean => state.user?.role === 'trainee',
        getUserRole: (state): string | null => state.user?.role || null,
        userName: (state): string => state.user?.username || '',
        userId: (state): number | null => state.user?.id || null
    },

    actions: {
        async login(data: { username: string; password: string }) {
            const res = await axios.post('/api/login', data)

            this.token = res.data.token
            this.user = res.data.user

            // zapisz token lokalnie i ustaw nagłówek globalnie
            localStorage.setItem('token', this.token)
            axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`

            // Pobierz szczegóły użytkownika
            await this.fetchUser()
        },

        async register(data: { username: string; password: string; role: string }) {
            const lowercaseRole = data.role.toLowerCase()
            await axios.post('/api/register', {
                username: data.username,
                password: data.password,
                role: lowercaseRole
            })
            await this.login({
                username: data.username,
                password: data.password,
            })
        },

        async fetchUser() {
            try {
                const res = await axios.get('/api/me')
                this.user = {
                    id: res.data.userID,
                    role: res.data.role,
                    username: res.data.username || 'Użytkownik'
                }
            } catch (error) {
                console.error('Error fetching user:', error)
            }
        },

        logout() {
            this.token = ''
            this.user = null
            this.profile = null
            localStorage.removeItem('token')
            delete axios.defaults.headers.common['Authorization']
        },

        async fetchProfile() {
            try {
                const res = await axios.get('/api/profile')
                this.profile = res.data.profile
            } catch (e) {
                console.error('Error fetching profile:', e)
                this.profile = null
            }
        },

        // Metoda inicjalizacji dla kompatybilności
        async initialize(): Promise<void> {
            if (this.token && !this.user) {
                // Ustaw nagłówek autoryzacji jeśli token istnieje
                axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
                await this.fetchUser()
            }
        }
    }
})
