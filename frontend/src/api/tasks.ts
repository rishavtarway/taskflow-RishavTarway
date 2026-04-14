import { client } from './client'

export interface Task {
  id: string
  title: string
  description?: string
  status: 'todo' | 'in_progress' | 'done'
  priority: 'low' | 'medium' | 'high'
  project_id: string
  assignee_id?: string
  due_date?: string
  created_at: string
  updated_at: string
}

export interface CreateTaskInput {
  title: string
  description?: string
  priority?: string
  assignee_id?: string
  due_date?: string
}

export interface UpdateTaskInput {
  title?: string
  description?: string
  status?: string
  priority?: string
  assignee_id?: string
  due_date?: string
}

export const tasksAPI = {
  list: async (projectId: string, status?: string, assignee?: string) => {
    const params = new URLSearchParams()
    if (status) params.append('status', status)
    if (assignee) params.append('assignee', assignee)
    const response = await client.get<{ tasks: Task[] }>(`/projects/${projectId}/tasks?${params}`)
    return response.data
  },

  create: async (projectId: string, data: CreateTaskInput) => {
    const response = await client.post<Task>(`/projects/${projectId}/tasks`, data)
    return response.data
  },

  update: async (id: string, data: UpdateTaskInput) => {
    const response = await client.patch<Task>(`/tasks/${id}`, data)
    return response.data
  },

  delete: async (id: string) => {
    await client.delete(`/tasks/${id}`)
  },
}