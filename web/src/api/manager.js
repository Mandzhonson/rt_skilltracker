import apiClient from '../utils/apiClient.js';

export const managerAPI = {
  getEmployees: () => apiClient.get('/manager/employees'),
  listPlans: (params) => apiClient.get('/manager/plans', { params }),
  createPlan: (data) => apiClient.post('/manager/plans', data),
  getPlan: (planId) => apiClient.get(`/manager/plans/${planId}`),
  updatePlan: (planId, data) => apiClient.patch(`/manager/plans/${planId}`, data),
  deletePlan: (planId) => apiClient.delete(`/manager/plans/${planId}`),
  createTask: (planId, data) => apiClient.post(`/manager/plans/${planId}/tasks`, data),
  getTask: (taskId) => apiClient.get(`/manager/tasks/${taskId}`),
  updateTask: (taskId, data) => apiClient.patch(`/manager/tasks/${taskId}`, data),
  deleteTask: (taskId) => apiClient.delete(`/manager/tasks/${taskId}`),
};