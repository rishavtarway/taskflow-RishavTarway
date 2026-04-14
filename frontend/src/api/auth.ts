import { client } from './client'

export interface User {
  id: string
  name: string
  email: string
  created_at: string
}

interface AuthResponse {
  token: string
  user: User
}

export const authAPI = {
  register: async (name: string, email: string, password: string) => {
    const response = await client.post<AuthResponse>('/auth/register', { name, email, password })
    return response.data
  },

  login: async (email: string, password: string) => {
    const response = await client.post<AuthResponse>('/auth/login', { email, password })
    return response.data
  },
}