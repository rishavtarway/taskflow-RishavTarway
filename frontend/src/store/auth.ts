import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  name: string
  email: string
  created_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
  hydrate: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: (token, user) => {
        set({ token, user, isAuthenticated: true })
      },

      logout: () => {
        set({ token: null, user: null, isAuthenticated: false })
      },

      hydrate: () => {
        const auth = localStorage.getItem('taskflow_auth')
        if (auth) {
          try {
            const { token, user } = JSON.parse(auth)
            if (token && user) {
              set({ token, user, isAuthenticated: true })
            }
          } catch {
            localStorage.removeItem('taskflow_auth')
          }
        }
      },
    }),
    {
      name: 'taskflow_auth',
      partialize: (state) => ({ token: state.token, user: state.user }),
    }
  )
)