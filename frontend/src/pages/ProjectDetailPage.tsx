import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2 } from 'lucide-react'
import { projectsAPI } from '../api/projects'
import { tasksAPI, type Task, type CreateTaskInput, type UpdateTaskInput } from '../api/tasks'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Textarea } from '../components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../components/ui/dialog'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const taskSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  description: z.string().optional(),
  priority: z.string().optional(),
  status: z.string().optional(),
  due_date: z.string().optional(),
})

type TaskForm = z.infer<typeof taskSchema>

function TaskCard({
  task,
  onUpdate,
  onDelete,
}: {
  task: Task
  onUpdate: (id: string, data: UpdateTaskInput) => void
  onDelete: (id: string) => void
}) {
  const [isEditing, setIsEditing] = useState(false)
  const queryClient = useQueryClient()

  const { register, handleSubmit } = useForm<TaskForm>({
    resolver: zodResolver(taskSchema),
    defaultValues: {
      title: task.title,
      description: task.description || '',
      priority: task.priority,
      status: task.status,
      due_date: task.due_date || '',
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTaskInput }) => tasksAPI.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', task.project_id] })
      setIsEditing(false)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: tasksAPI.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', task.project_id] })
    },
  })

  const onSubmit = (data: TaskForm) => {
    updateMutation.mutate({ id: task.id, data })
  }

  const statusOptions = [
    { value: 'todo', label: 'To Do' },
    { value: 'in_progress', label: 'In Progress' },
    { value: 'done', label: 'Done' },
  ]

  const priorityOptions = [
    { value: 'low', label: 'Low' },
    { value: 'medium', label: 'Medium' },
    { value: 'high', label: 'High' },
  ]

  const priorityColors = {
    low: 'bg-green-100 text-green-800',
    medium: 'bg-yellow-100 text-yellow-800',
    high: 'bg-red-100 text-red-800',
  }

  if (isEditing) {
    return (
      <form onSubmit={handleSubmit(onSubmit)} className="rounded-lg border bg-card p-4 space-y-4">
        <Input {...register('title')} placeholder="Task title" />
        <Textarea {...register('description')} placeholder="Description" />
        <div className="flex gap-2">
          <Select {...register('status')} defaultValue={task.status}>
            <SelectTrigger className="w-32">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {statusOptions.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select {...register('priority')} defaultValue={task.priority}>
            <SelectTrigger className="w-28">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {priorityOptions.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input type="date" {...register('due_date')} className="w-36" defaultValue={task.due_date} />
        </div>
        <div className="flex gap-2">
          <Button type="submit" size="sm" disabled={updateMutation.isPending}>
            {updateMutation.isPending ? 'Saving...' : 'Save'}
          </Button>
          <Button type="button" variant="outline" size="sm" onClick={() => setIsEditing(false)}>
            Cancel
          </Button>
        </div>
      </form>
    )
  }

  return (
    <div className="rounded-lg border bg-card p-4">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h4 className="font-medium">{task.title}</h4>
          {task.description && (
            <p className="mt-1 text-sm text-muted-foreground">{task.description}</p>
          )}
          <div className="mt-2 flex items-center gap-2">
            <span className={`rounded px-2 py-0.5 text-xs font-medium ${priorityColors[task.priority]}`}>
              {task.priority}
            </span>
            {task.due_date && (
              <span className="text-xs text-muted-foreground">Due: {task.due_date}</span>
            )}
          </div>
        </div>
        <Select
          value={task.status}
          onValueChange={(value) => onUpdate(task.id, { status: value })}
        >
          <SelectTrigger className="w-32">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {statusOptions.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="mt-3 flex gap-2">
        <Button variant="outline" size="sm" onClick={() => setIsEditing(true)}>
          Edit
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onDelete(task.id)}
          disabled={deleteMutation.isPending}
        >
          <Trash2 className="h-4 w-4 text-destructive" />
        </Button>
      </div>
    </div>
  )
}

function CreateTaskDialog({ projectId }: { projectId: string }) {
  const [open, setOpen] = useState(false)
  const queryClient = useQueryClient()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<TaskForm>({
    resolver: zodResolver(taskSchema),
    defaultValues: { title: '', priority: 'medium' },
  })

  const createMutation = useMutation({
    mutationFn: (data: CreateTaskInput) => tasksAPI.create(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
      setOpen(false)
      reset()
    },
  })

  const onSubmit = (data: TaskForm) => {
    createMutation.mutate({
      title: data.title,
      description: data.description,
      priority: data.priority || 'medium',
      due_date: data.due_date,
    })
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Add Task
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Task</DialogTitle>
          <DialogDescription>Add a new task to this project.</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <label htmlFor="title" className="text-sm font-medium">
                Title
              </label>
              <Input id="title" {...register('title')} />
              {errors.title && (
                <p className="text-sm text-destructive">{errors.title.message}</p>
              )}
            </div>
            <div className="grid gap-2">
              <label htmlFor="description" className="text-sm font-medium">
                Description
              </label>
              <Textarea id="description" {...register('description')} />
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div className="grid gap-2">
                <label htmlFor="priority" className="text-sm font-medium">
                  Priority
                </label>
                <Select {...register('priority')} defaultValue="medium">
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="low">Low</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="high">High</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <label htmlFor="due_date" className="text-sm font-medium">
                  Due Date
                </label>
                <Input type="date" id="due_date" {...register('due_date')} />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Creating...' : 'Create Task'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: project, isLoading, error } = useQuery({
    queryKey: ['project', id],
    queryFn: () => projectsAPI.get(id!),
    enabled: !!id,
  })

  const deleteProjectMutation = useMutation({
    mutationFn: projectsAPI.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      navigate('/projects')
    },
  })

  const updateTaskMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTaskInput }) => tasksAPI.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', id] })
    },
  })

  const deleteTaskMutation = useMutation({
    mutationFn: tasksAPI.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', id] })
    },
  })

  const handleTaskUpdate = (taskId: string, data: UpdateTaskInput) => {
    updateTaskMutation.mutate({ id: taskId, data })
  }

  const handleTaskDelete = (taskId: string) => {
    if (confirm('Delete this task?')) {
      deleteTaskMutation.mutate(taskId)
    }
  }

  const handleDeleteProject = () => {
    if (confirm('Delete this project and all its tasks?')) {
      deleteProjectMutation.mutate(id!)
    }
  }

  if (isLoading) {
    return (
      <div className="container py-8">
        <div className="h-64 rounded-lg border bg-muted animate-pulse" />
      </div>
    )
  }

  if (error || !project) {
    return (
      <div className="container py-8">
        <div className="rounded-md bg-destructive/10 p-4 text-destructive">
          Failed to load project.
        </div>
      </div>
    )
  }

  const tasks = project.tasks || []
  const todoTasks = tasks.filter((t) => t.status === 'todo')
  const inProgressTasks = tasks.filter((t) => t.status === 'in_progress')
  const doneTasks = tasks.filter((t) => t.status === 'done')

  return (
    <div className="container py-8">
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to="/projects">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold">{project.name}</h1>
            {project.description && (
              <p className="text-muted-foreground">{project.description}</p>
            )}
          </div>
        </div>
        <div className="flex gap-2">
          <CreateTaskDialog projectId={project.id} />
          <Button variant="destructive" onClick={handleDeleteProject}>
            Delete Project
          </Button>
        </div>
      </div>

      {tasks.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No tasks yet. Add one to get started!</p>
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-3">
          <div>
            <h2 className="mb-4 text-sm font-semibold text-muted-foreground">
              To Do ({todoTasks.length})
            </h2>
            <div className="space-y-3">
              {todoTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onUpdate={handleTaskUpdate}
                  onDelete={handleTaskDelete}
                />
              ))}
            </div>
          </div>
          <div>
            <h2 className="mb-4 text-sm font-semibold text-muted-foreground">
              In Progress ({inProgressTasks.length})
            </h2>
            <div className="space-y-3">
              {inProgressTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onUpdate={handleTaskUpdate}
                  onDelete={handleTaskDelete}
                />
              ))}
            </div>
          </div>
          <div>
            <h2 className="mb-4 text-sm font-semibold text-muted-foreground">
              Done ({doneTasks.length})
            </h2>
            <div className="space-y-3">
              {doneTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onUpdate={handleTaskUpdate}
                  onDelete={handleTaskDelete}
                />
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}