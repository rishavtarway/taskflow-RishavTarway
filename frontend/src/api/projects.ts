import { client } from './client'

export interface Project {
  id: string
  name: string
  description?: string
  owner_id: string
  created_at: string
  tasks?: Task[]
}

export interface CreateProjectInput {
  name: string
  description?: string
}

export interface UpdateProjectInput {
  name?: string
  description?: string
}

export const projectsAPI = {
  list: async () => {
    const response = await client.get<{ projects: Project[] }>('/projects')
    return response.data
  },

  get: async (id: string) => {
    const response = await client.get<Project>(`/projects/${id}`)
    return response.data
  },

  create: async (data: CreateProjectInput) => {
    const response = await client.post<Project>('/projects', data)
    return response.data
  },

  update: async (id: string, data: UpdateProjectInput) => {
    const response = await client.patch<Project>(`/projects/${id}`, data)
    return response.data
  },

  delete: async (id: string) => {
    await client.delete(`/projects/${id}`)
  },
}